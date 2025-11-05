"""
Servicio de lógica de negocio para anomalías
Orquesta las llamadas entre ParserService, ML Engine y BD
"""

from typing import List, Dict, Any, Optional
from app.models.model_manager import model_manager
import httpx
import logging
from datetime import datetime
from collections import defaultdict
import numpy as np
import pandas as pd
import warnings
import os
import sys

# ✅ SOLUCIÓN DEFINITIVA: Deshabilitar tqdm globalmente
os.environ['SHAP_SHOW_PROGRESS'] = 'false'

# ✅ Monkey patch tqdm antes de importar shap
class DummyTqdm:
    def __init__(self, *args, **kwargs):
        self.iterable = args[0] if args else []
    def __iter__(self):
        return iter(self.iterable)
    def __enter__(self):
        return self
    def __exit__(self, *args):
        pass
    def update(self, *args, **kwargs):
        pass
    def close(self):
        pass

sys.modules['tqdm'] = type(sys)('tqdm')
sys.modules['tqdm'].tqdm = DummyTqdm
sys.modules['tqdm.auto'] = type(sys)('tqdm.auto')
sys.modules['tqdm.auto'].tqdm = DummyTqdm

import shap

# ✅ Silenciar warnings
logging.getLogger('shap').setLevel(logging.ERROR)
warnings.filterwarnings('ignore', category=UserWarning, module='shap')
warnings.filterwarnings('ignore', category=FutureWarning, module='shap')

from app.config import config
from app.models.anomaly import (
    MLRequest, MLSensorRecord, MLResponse,
    AnomalyCreate, ParserResponse
)
from app.repositories.anomaly_repository import AnomalyRepository

logger = logging.getLogger(__name__)

# ========================================
# ✅ FUNCIONES HELPER HTM (de ml_engine)
# ========================================
def severity_from_likelihood(likelihood: float) -> str:
    """Determina severidad basada en likelihood (0-100)"""
    if likelihood < 40:
        return "baja"
    elif likelihood < 70:
        return "media"
    else:
        return "alta"


def description_from_severity(sev: str) -> str:
    """Genera descripción basada en severidad"""
    d = {
        "baja": "Funcionamiento dentro del rango esperado.",
        "media": "Comportamiento irregular detectado. Revisar condiciones operativas.",
        "alta": "Anomalía crítica detectada. Atención prioritaria requerida."
    }
    return d.get(sev, "Sin descripción disponible")

# ========================================

class AnomalyService:
    """Servicio de lógica de negocio para detección de anomalías"""
    
    def __init__(self, repository: AnomalyRepository):
        self.repository = repository
        self._shap_explainer = None  # Cache del explainer
    
    async def fetch_data_from_parser(
        self,
        activo_id: int,
        page: int = 1,
        limit: int = 1000,
    ) -> ParserResponse:
        """
        Obtiene datos del ParserService (iot-service)
        
        Args:
            activo_id: ID del activo
            page: Página de resultados
            limit: Cantidad de registros por página
            
        Returns:
            ParserResponse con los datos del activo
        """
        url = f"{config.iot_service_url}/lectura/{activo_id}/window"
        params = {"page": page, "limit": limit}
        
        logger.info(f"📡 Solicitando datos del ParserService: activo_id={activo_id}, limit={limit}")
        
        try:
            async with httpx.AsyncClient(timeout=config.iot_service_timeout) as client:
                response = await client.get(url, params=params)
                response.raise_for_status()
                
                data = response.json()
                parser_response = ParserResponse(**data)
                
                logger.info(
                    f"✅ Datos recibidos: {len(parser_response.sensores)} sensores, "
                    f"~{limit} registros solicitados"
                )
                
                return parser_response
                
        except httpx.HTTPStatusError as e:
            logger.error(f"❌ Error HTTP {e.response.status_code} del ParserService: {e}")
            raise
        except httpx.RequestError as e:
            logger.error(f"❌ Error de conexión con ParserService: {e}")
            raise
        except Exception as e:
            logger.error(f"❌ Error inesperado obteniendo datos: {e}")
            raise

    def transform_to_ml_format(
        self,
        parser_data: ParserResponse
    ) -> List[MLSensorRecord]:
        """
        Transforma datos del ParserService al formato del ML Engine
        
        El ParserService retorna datos por sensor:
        {
            "sensores": [
                {"id_sensor": "Temp_Agua", "datos": [{"tiempo": "2024-...", "valor": 22.3}, ...]},
                {"id_sensor": "Vibracion", "datos": [{"tiempo": "2024-...", "valor": 0.05}, ...]}
            ]
        }
        
        El ML Engine espera registros agrupados por timestamp:
        [
            {"timestamp": "2024-...", "Temp_Agua": 22.3, "Vibracion": 0.05, ...},
            ...
        ]
        """
        # Agrupar por timestamp
        records_by_time: Dict[str, Dict[str, float]] = defaultdict(dict)
        
        for sensor in parser_data.sensores:
            if sensor.datos:
                for dp in sensor.datos:
                    records_by_time[dp.tiempo][sensor.id_sensor] = dp.valor
            
        
        # Convertir a lista de MLSensorRecord
        ml_records = [
            {"timestamp": t, **vals}
            for t, vals in records_by_time.items()
            if len(vals) >= 3
        ]
        
        logger.info(
            f"📊 Transformados {len(ml_records)} registros agrupados por timestamp "
            f"desde {len(parser_data.sensores)} sensores"
        )
        logger.info(f"   ✓ {len(ml_records)} registros obtenidos")
        return ml_records
    
    def _get_feature_contribution(
        self,
        features: np.ndarray,
        feature_names: List[str],
        is_anomaly: bool
    ) -> Dict[str, Any]:
        """
        Calcula la variable con mayor contribución (mayor SHAP value absoluto) 
        para un registro anómalo usando SHAP
        
        Args:
            features: Array de features del registro [1, n_features]
            feature_names: Nombres de las features
            is_anomaly: Si el registro es anomalía o no
            
        Returns:
            Dict con variable más influyente y su contribución, o None si no hay anomalía
        """
        if not is_anomaly or model_manager.model is None:
            return {
                "most_influential_variable": None,
                "contribution_magnitude": None
            }
        
        try:
            # ✅ Crear explainer si no existe (cache)
            if self._shap_explainer is None:
                try:
                    # ✅ Usar LinearExplainer para modelos lineales
                    self._shap_explainer = shap.LinearExplainer(
                        model_manager.model,
                        model_manager.scaler_X.transform(model_manager.X_train) 
                        if hasattr(model_manager, 'X_train') and hasattr(model_manager.scaler_X, 'mean_') 
                        else features
                    )
                    logger.info("✅ SHAP LinearExplainer inicializado")
                except:
                    # ✅ Fallback: KernelExplainer (más lento pero universal)
                    background = shap.sample(features, min(100, len(features)))
                    self._shap_explainer = shap.KernelExplainer(
                        model_manager.model.predict,
                        background
                    )
                    logger.info("✅ SHAP KernelExplainer inicializado (fallback)")
        
            # ✅ PASO 1: Asegurar que features sea 2D [1, n_features]
            if features.ndim == 1:
                features = features.reshape(1, -1)
        
            # ✅ PASO 2: Escalar features si el scaler está disponible
            if hasattr(model_manager.scaler_X, 'mean_'):
                features_scaled = model_manager.scaler_X.transform(features)
            else:
                features_scaled = features
        
            # ✅ PASO 3: Calcular SHAP values para ESTE ÚNICO REGISTRO
            shap_values = self._shap_explainer.shap_values(features_scaled)
        
            # ✅ PASO 4: Manejar diferentes formatos de salida de SHAP
            # - Clasificación multiclase: shap_values es lista [[valores_clase1], [valores_clase2], ...]
            # - Regresión/Clasificación binaria: shap_values es array [1, n_features] o [n_features]
            if isinstance(shap_values, list):
                # Tomar primera clase (o última según el problema)
                shap_values_flat = np.array(shap_values[0])
            else:
                shap_values_flat = shap_values
        
            # ✅ PASO 5: Si es 2D [1, n_features], aplanar a [n_features]
            if shap_values_flat.ndim == 2:
                shap_values_flat = shap_values_flat.flatten()
        
            # ✅ PASO 6: Obtener magnitudes ABSOLUTAS (importancia sin considerar dirección)
            abs_shap_values = np.abs(shap_values_flat)
        
            # ✅ PASO 7: Encontrar el índice con MAYOR contribución absoluta
            max_idx = int(np.argmax(abs_shap_values))
            max_contribution = float(abs_shap_values[max_idx])
        
            # ✅ PASO 8: Validar índice y obtener nombre de variable
            if max_idx >= len(feature_names):
                logger.error(
                    f"❌ Índice {max_idx} fuera de rango. "
                    f"feature_names tiene {len(feature_names)} elementos, "
                    f"shap_values tiene {len(abs_shap_values)} valores"
                )
                variable_name = f"Feature_{max_idx}_OutOfRange"
            else:
                variable_name = feature_names[max_idx]
        
            return {
                "most_influential_variable": variable_name,
                "contribution_magnitude": round(max_contribution, 4)
            }
        
        except Exception as e:
            logger.warning(f"⚠️  Error calculando feature contribution: {e}", exc_info=True)
            return {
                "most_influential_variable": None,
                "contribution_magnitude": None
            }
    
    async def predict_anomalies(
        self,
        activo_id: int,
        ml_records: List[MLSensorRecord]
    ) -> MLResponse:
        """
        Predice anomalías utilizando el modelo integrado local
        
        Args:
            activo_id: ID del activo
            ml_records: Lista de registros en formato ML
            
        Returns:
            MLResponse con las predicciones
        """
        logger.info(
            f"🤖 Prediciendo anomalías para {len(ml_records)} registros "
            f"del activo {activo_id}"
        )
        
        try:
            from app.models.model_manager import model_manager
            
            # ✅ Constantes HTM
            ALPHA = 0.3  # Factor de suavizado EWMA
            K_ADAPT = 3.0  # Factor threshold adaptativo
            WINDOW_SIZE = 100  # Ventana para rolling threshold
            MIN_PERIODS = 50  # Mínimo de períodos para calcular threshold
            
            # Verificar si el modelo está disponible
            model_available = False
            if model_manager.model is None:
                logger.warning("⚠️  Modelo no disponible, usando valores de fallback")
                model_available = False
            else:
                model_available = True
                logger.info("✅ Modelo disponible para predicción")
            
            # ===================================
            # ✅ PREPARAR DATOS (VENTANA COMPLETA)
            # ===================================
            # Fase 2: Preprocesamiento
            logger.info("🔧 Fase 2: Preprocesamiento")
            X_list = []
            y_list = []
            feature_names = []
            timestamps = []  # ✅ Inicializar ANTES del loop
            
            for rec in ml_records:
                vals = [v for k, v in rec.items() if k != 'timestamp']
                if vals and not feature_names:  # Capturar nombres solo una vez
                    feature_names = [k for k in rec.keys() if k != 'timestamp']
                if vals:
                    X_list.append(vals)
                    y_list.append(np.mean(vals))
                    timestamps.append(rec['timestamp'])  # ✅ Agregar timestamp al mismo tiempo
            
            # Imprimir nombres de features
            logger.info(f"   📋 Features detectadas: {feature_names}")
            logger.info(f"   📋 Cantidad de features: {len(feature_names)}")
            
            # ✅ DIAGNÓSTICO: Verificar longitudes de X_list
            lengths = [len(x) for x in X_list]
            unique_lengths = set(lengths)
            
            if len(unique_lengths) > 1:
                logger.error(
                    f"❌ INCONSISTENCIA DETECTADA:\n"
                    f"   • Longitudes únicas: {unique_lengths}\n"
                    f"   • Distribución: {dict(zip(*np.unique(lengths, return_counts=True)))}\n"
                    f"   • feature_names tiene: {len(feature_names)} features"
                )
                
                # Mostrar ejemplos de registros con diferentes longitudes
                for length in unique_lengths:
                    idx = lengths.index(length)
                    logger.error(
                        f"   • Registro {idx} (timestamp={timestamps[idx]}) tiene {length} valores"
                    )
                    # Mostrar qué sensores tiene este registro
                    rec_sensors = [k for k in ml_records[idx].keys() if k != 'timestamp']
                    logger.error(f"     Sensores: {rec_sensors}")
                
                # ✅ SOLUCIÓN: Usar la longitud más común (la mayoría de registros)
                most_common_length = max(set(lengths), key=lengths.count)
                logger.warning(
                    f"⚠️  Longitud más común: {most_common_length} "
                    f"({lengths.count(most_common_length)} registros)"
                )
                
                # Actualizar feature_names a la longitud correcta
                for rec in ml_records:
                    vals = [v for k, v in rec.items() if k != 'timestamp']
                    if len(vals) == most_common_length:
                        feature_names = [k for k in rec.keys() if k != 'timestamp']
                        break
                
                logger.info(f"   📋 Features actualizadas: {feature_names}")
                logger.info(f"   📋 Cantidad de features: {len(feature_names)}")
                
                # Filtrar registros con la longitud correcta
                filtered_data = [
                    (x, ts, y) 
                    for x, ts, y in zip(X_list, timestamps, y_list) 
                    if len(x) == most_common_length
                ]
                
                if not filtered_data:
                    raise ValueError(
                        f"No hay registros válidos con {most_common_length} features. "
                        f"Longitudes encontradas: {unique_lengths}"
                    )
                
                # Desempaquetar datos filtrados
                X_list, timestamps, y_list = zip(*filtered_data)
                X_list = list(X_list)
                timestamps = list(timestamps)
                y_list = list(y_list)
                
                logger.warning(
                    f"⚠️  Filtrados {len(X_list)} registros válidos de {len(ml_records)} totales "
                    f"({len(ml_records) - len(X_list)} descartados por tener {list(unique_lengths - {most_common_length})} features)"
                )
            
            # ✅ CONVERSIÓN DIRECTA (sin padding)
            X = np.array(X_list)
            y = np.array(y_list)
            
            logger.info(f"   ✓ X: {X.shape}, y: {y.shape}, timestamps: {len(timestamps)}")
            timestamps = [rec['timestamp'] for rec in ml_records]
            # ✅ LOG para verificar dimensiones
            logger.info(f"📊 Ventana preparada: shape={X.shape} (debe ser [n_samples, n_sensores])")
            logger.info(f"📊 Features por registro: {X.shape[1]}")  # ✅ Usar X.shape[1] en lugar de max_feat
            
            # ✅ LOG DE MUESTRA DE PRIMER REGISTRO
            if len(ml_records) > 0:
                first_record = ml_records[0].copy()
                first_timestamp = first_record.pop('timestamp')
                logger.info(f"📝 Primer registro ({first_timestamp}):")
                for sensor_name, value in first_record.items():
                    logger.info(f"   • {sensor_name}: {value}")
            
            # ===================================
            # ✅ CALCULAR ESTADÍSTICAS DE LA VENTANA
            # ===================================
            mu_global = np.mean(X, axis=0)
            sigma_global = np.std(X, axis=0) + 1e-6

            logger.info(f"📊 Ventana: {len(X)} registros, {X.shape[1]} features")
            logger.info(f"📊 μ_global: {mu_global.mean():.4f}, σ_global: {sigma_global.mean():.4f}")
            
            results = []
            
            # ✅ Variables para cálculo HTM (inicializadas localmente, NO del modelo)
            err_ref = 1.0  # Se inicializa localmente y se actualiza con EWMA
            likelihood_value = 0.0  # Likelihood suavizado (EWMA)
            
            # ===================================
            # ✅ CALCULAR SCORES PARA CADA REGISTRO
            # ===================================
            for i in range(len(X)):
                try:
                    features = X[i:i+1]
                    actual_value = features[0][0] if len(features[0]) > 0 else 0.0
                    
                    if model_available:
                        # ===================================
                        # ✅ CON MODELO: Cálculo HTM completo
                        # ===================================
                        try:
                            # ✅ LOG de verificación detallado en el PRIMER registro
                            if i == 0:
                                logger.info(
                                    f"🔍 Análisis primer registro:\n"
                                    f"   • Features shape: {features.shape}\n"
                                    f"   • Features values: {features[0]}\n"
                                    f"   • Timestamp: {timestamps[i]}"
                                )
                            
                            # Escalar features si el scaler está entrenado
                            if hasattr(model_manager.scaler_X, 'mean_'):                                
                                features_scaled = model_manager.scaler_X.transform(features)
                            else:
                                features_scaled = features
                            
                            # Predecir con el modelo
                            prediction = model_manager.model.predict(features_scaled)[0]
                            
                            # ✅ Error normalizado (HTM-style)
                            y_true_normalized = (X[i] - mu_global) / sigma_global
                            err = np.sqrt(np.mean(y_true_normalized ** 2))
                            
                            # ✅ Error de referencia adaptativo (EWMA) - SE CALCULA LOCALMENTE
                            err_ref = (1 - ALPHA) * err_ref + ALPHA * err
                            
                            # ✅ AnomalyScore: error escalado 0-100
                            anomaly_score = np.clip((err / (err_ref + 1e-6)) * 50, 0, 100)
                            
                            # ✅ Likelihood suavizado (EWMA)
                            likelihood_value = (1 - ALPHA) * likelihood_value + ALPHA * anomaly_score
                        
                        except Exception as e:
                            #logger.warning(f"⚠️  Error en cálculo HTM punto {i}: {e}")
                            prediction = actual_value
                            anomaly_score = 20.0
                            likelihood_value = 20.0
                    
                    else:
                        # ===================================
                        # ✅ SIN MODELO: Fallback simple
                        # ===================================
                        prediction = actual_value + np.random.normal(0, 0.1)
                        
                        # Score basado en desviación de la media global
                        deviation = np.abs(X[i] - mu_global).mean()
                        anomaly_score = min((deviation / (sigma_global.mean() + 1e-6)) * 50, 100)
                        likelihood_value = anomaly_score
                    
                    # Normalizar a 0-100
                    anomaly_score = float(np.clip(anomaly_score, 0, 100))
                    likelihood_value = float(np.clip(likelihood_value, 0, 100))
                    
                    # ✅ Threshold se calculará después con rolling window
                    threshold = 50.0  # Valor temporal
                    
                    # ✅ Severidad basada en likelihood
                    severity = severity_from_likelihood(likelihood_value)
                    
                    # ✅ Descripción basada en severidad
                    description = description_from_severity(severity)
                    
                    # ✅ Crear resultado con is_anomaly como False por defecto
                    result = {
                        "timestamp": timestamps[i],
                        "is_anomaly": False,  # Se calculará después con threshold adaptativo
                        "AnomalyScore": round(anomaly_score, 2),
                        "AnomalyLikelihood": round(likelihood_value, 2),
                        "Threshold": threshold,  # Se actualizará
                        "Severity": severity,
                        "Description": description,
                        "ReconstructedValue": round(float(prediction), 4),
                        "most_influential_variable": None,  # ✅ Nuevo campo
                        "contribution_magnitude": None      # ✅ Nuevo campo
                    }
                    
                    results.append(result)
                    
                except Exception as e:
                    logger.warning(f"⚠️  Error procesando registro {i}: {e}")
                    continue
            
            logger.info(f"✅ Cálculo de scores completado para {len(results)} registros")
        
            # ===================================
            # ✅ CALCULAR THRESHOLD ADAPTATIVO, is_anomaly Y EXPLICABILIDAD
            # ===================================
            if results:
                df_results = pd.DataFrame(results)
                
                # Calcular rolling threshold sobre AnomalyLikelihood
                likelihood_series = df_results['AnomalyLikelihood']
                
                # Rolling mean + k * rolling std
                mean_rolling = likelihood_series.rolling(
                    window=WINDOW_SIZE, 
                    min_periods=MIN_PERIODS
                ).mean()
                
                std_rolling = likelihood_series.rolling(
                    window=WINDOW_SIZE, 
                    min_periods=MIN_PERIODS
                ).std()
                
                threshold_series = (mean_rolling + K_ADAPT * std_rolling).clip(lower=0)
                
                # Rellenar NaN en la ventana inicial con un valor por defecto
                # Puedes usar la media global o un percentil
                default_threshold = likelihood_series.quantile(0.75) if len(likelihood_series) > 0 else 50.0
                threshold_series = threshold_series.fillna(default_threshold)
                
                # ✅ Actualizar threshold y calcular is_anomaly como booleano
                df_results['Threshold'] = threshold_series.round(2)
                df_results['is_anomaly'] = (
                    df_results['AnomalyLikelihood'] > df_results['Threshold']
                )  # Ya es booleano por pandas
                
                # ===================================
                # ✅ MÉTRICAS DE ESTABILIDAD (Likelihood vs Threshold)
                # ===================================
                dif = np.abs(df_results['AnomalyLikelihood'] - df_results['Threshold'])
                mean_dif = np.mean(dif)
                std_dif = np.std(dif)
                cv_dif = std_dif / (mean_dif + 1e-6)
                
                
                # Interpretación del CV de diferencia
                if cv_dif < 0.5:
                    interpretacion = "Sistema ESTABLE - Detección consistente"
                    nivel_estabilidad = "🟢 ALTA"
                elif cv_dif < 1.0:
                    interpretacion = "Sistema MODERADO - Algunas fluctuaciones"
                    nivel_estabilidad = "🟡 MEDIA"
                else:
                    interpretacion = "Sistema INESTABLE - Alta variabilidad en detección"
                    nivel_estabilidad = "🔴 BAJA"
                
                # ✅ Actualizar severidad y descripción solo para anomalías (is_anomaly == True)
                # ✅ CALCULAR FEATURE CONTRIBUTION SOLO PARA ANOMALÍAS
                anomaly_mask = df_results['is_anomaly'] == True
                
                if anomaly_mask.any():
                    logger.info(f"🔍 Calculando explicabilidad para {anomaly_mask.sum()} anomalías")
                    
                    # ✅ LOG DE MUESTRA DE EXPLICABILIDAD (primeras 3 anomalías)
                    anomaly_indices = df_results[anomaly_mask].index[:3]
                    logger.info("=" * 70)
                    logger.info("🧠 MUESTRA DE EXPLICABILIDAD (primeras 3 anomalías)")
                    logger.info("=" * 70)
                    
                    for idx in df_results[anomaly_mask].index:
                        try:
                            # Obtener features del registro original
                            features = X[idx:idx+1]
                            
                            # Calcular contribución
                            contribution = self._get_feature_contribution(
                                features=features,
                                feature_names=feature_names,
                                is_anomaly=True
                            )
                            
                            # ✅ LOG DETALLADO PARA LAS PRIMERAS 3 ANOMALÍAS
                            if idx in anomaly_indices:
                                logger.info(
                                    f"  🚨 Anomalía [{df_results.loc[idx, 'timestamp']}]:\n"
                                    f"     • Variable más influyente: {contribution['most_influential_variable']}\n"
                                    f"     • Magnitud de contribución: {contribution['contribution_magnitude']}\n"
                                    f"     • AnomalyLikelihood: {df_results.loc[idx, 'AnomalyLikelihood']:.2f}\n"
                                    f"     • Threshold: {df_results.loc[idx, 'Threshold']:.2f}\n"
                                    f"     • Severidad: {df_results.loc[idx, 'Severity']}"
                                )
                            
                            # Actualizar DataFrame
                            df_results.loc[idx, 'most_influential_variable'] = contribution['most_influential_variable']
                            df_results.loc[idx, 'contribution_magnitude'] = contribution['contribution_magnitude']
                            
                            # Actualizar severidad y descripción
                            likelihood = df_results.loc[idx, 'AnomalyLikelihood']
                            df_results.loc[idx, 'Severity'] = severity_from_likelihood(likelihood)
                            
                            # ✅ Descripción enriquecida con variable influyente
                            base_desc = description_from_severity(df_results.loc[idx, 'Severity'])
                            if contribution['most_influential_variable']:
                                enriched_desc = (
                                    f"{base_desc} "
                                    f"Verificar: {contribution['most_influential_variable']}."
                                )
                                df_results.loc[idx, 'Description'] = enriched_desc
                            
                        except Exception as e:
                            logger.warning(f"⚠️  Error calculando contribución para índice {idx}: {e}")
                            continue
                
                # Convertir de vuelta a lista de dicts
                results = df_results.to_dict('records')
                
                # ✅ Contar anomalías usando booleanos
                anomaly_count = df_results['is_anomaly'].sum()
                
                logger.info(
                    f"📊 Threshold adaptativo calculado. "
                    f"Ventana: {WINDOW_SIZE}, k={K_ADAPT}, "
                    f"Threshold promedio: {threshold_series.mean():.2f}, "
                    f"Anomalías detectadas: {anomaly_count}"
                )
                
                # ✅ LOG PRINCIPAL: Métricas de estabilidad
                logger.info(
                    f"📈 ESTABILIDAD DEL SISTEMA ({nivel_estabilidad}):\n"
                    f"   • Media diferencia (Likelihood - Threshold): {mean_dif:.2f}\n"
                    f"   • Desviación estándar: {std_dif:.2f}\n"
                    f"   • Coeficiente de variación: {cv_dif:.4f}\n"
                    f"   ➜ {interpretacion}"
                )
                
                # ✅ LOG de advertencia si hay inestabilidad
                if cv_dif > 1.0:
                    logger.warning(
                        f"⚠️  INESTABILIDAD DETECTADA (CV={cv_dif:.4f}): "
                        f"El sistema está fluctuando entre anomalía y normalidad. "
                        f"Considerar:\n"
                        f"   1. Aumentar WINDOW_SIZE (actual: {WINDOW_SIZE})\n"
                        f"   2. Ajustar K_ADAPT (actual: {K_ADAPT})\n"
                        f"   3. Revisar calidad de datos de entrada"
                    )
                elif cv_dif > 0.5:
                    logger.info(
                        f"ℹ️  Estabilidad moderada (CV={cv_dif:.4f}). "
                        f"El sistema funciona correctamente pero hay algunas fluctuaciones normales."
                    )
                else:
                    logger.info(
                        f"✅ Excelente estabilidad (CV={cv_dif:.4f}). "
                        f"El sistema está detectando anomalías de manera consistente."
                    )
                # Distribución de severidades
                severity_counts = df_results['Severity'].value_counts()
                logger.info(f"  • Distribución de severidades:")
                for severity, count in severity_counts.items():
                    logger.info(f"    - {severity}: {count} ({count/len(df_results)*100:.1f}%)")
                
                # Top 3 variables más influyentes (si hay anomalías)
                if anomaly_mask.any():
                    anomalies_df = df_results[anomaly_mask]
                    top_vars = anomalies_df['most_influential_variable'].value_counts().head(3)
                    logger.info(f"  • Top 3 variables más influyentes en anomalías:")
                    for var, count in top_vars.items():
                        logger.info(f"    - {var}: {count} veces")
                
                logger.info("=" * 70)
                
                # ✅ Estadísticas adicionales de detección
                if anomaly_count > 0:
                    anomaly_indices = df_results[anomaly_mask].index
                    likelihood_anomalies = df_results.loc[anomaly_indices, 'AnomalyLikelihood']
                    threshold_anomalies = df_results.loc[anomaly_indices, 'Threshold']
                    
                    # Cuánto superan el threshold las anomalías
                    margin = (likelihood_anomalies - threshold_anomalies).mean()
                    
                    logger.info(
                        f"🚨 Análisis de anomalías detectadas:\n"
                        f"   • Total: {anomaly_count} de {len(df_results)} registros ({anomaly_count/len(df_results)*100:.1f}%)\n"
                        f"   • Margen promedio sobre threshold: {margin:.2f}\n"
                        f"   • Likelihood promedio en anomalías: {likelihood_anomalies.mean():.2f}\n"
                        f"   • Threshold promedio en anomalías: {threshold_anomalies.mean():.2f}"
                    )

            else:
                anomaly_count = 0
            
            # Crear respuesta en formato MLResponse
            from app.models.anomaly import MLResponse as MLResponseClass
            
            ml_response = {
                "pump": f"Activo_{activo_id}",
                "results": results
            }
            
            mode_text = "con modelo HTM" if model_available else "sin modelo (datos iniciales)"
            logger.info(
                f"✅ Predicción completada {mode_text}. "
                f"Resultados: {len(results)}, "
                f"Anomalías detectadas: {anomaly_count}"
            )
            
            return ml_response
                
        except Exception as e:
            logger.error(f"❌ Error en predicción de anomalías: {e}", exc_info=True)
            raise
    
    def save_ml_results(
        self,
        activo_id: int,
        ml_response: Dict[str, Any]
    ) -> Dict[str, int]:
        """
        Guarda los resultados de predicción en la base de datos
        
        Args:
            activo_id: ID del activo
            ml_response: Respuesta del modelo de ML (dict)
            
        Returns:
            Dict con estadísticas de guardado
        """
        anomalies_to_create = []
        
        results = ml_response.get("results", [])
        
        logger.info(f"💾 Preparando guardado de {len(results)} resultados para activo {activo_id}")
        
        for idx, result in enumerate(results):
            try:
                # ✅ Extraer TODOS los valores calculados por predict_anomalies
                anomaly_score = result.get("AnomalyScore")
                anomaly_likelihood = result.get("AnomalyLikelihood")
                threshold = result.get("Threshold")
                severidad = result.get("Severity")
                descripcion = result.get("Description")
                is_anomaly_value = result.get("is_anomaly")
               
                # Normalizar severidad a minúsculas para BD
                severidad = result.get("Severity", "baja").lower()
                
                # Parsear timestamp
                try:
                    timestamp_str = result.get("timestamp")
                    if not timestamp_str:
                        logger.warning(f"⚠️  Registro {idx} sin timestamp, saltando")
                        continue
                    
                    # Intentar varios formatos
                    timestamp = datetime.fromisoformat(
                        timestamp_str.replace('Z', '+00:00')
                    )
                except (ValueError, AttributeError) as e:
                    logger.warning(f"⚠️  Timestamp inválido en registro {idx}: '{timestamp_str}' - {e}")
                    # Fallback: usar timestamp actual
                    timestamp = datetime.utcnow()
                                
                # ✅ Convertir is_anomaly a booleano explícitamente
                is_anomaly_value = result.get("is_anomaly", False)
                if isinstance(is_anomaly_value, (int, float)):
                    is_anomaly_value = bool(is_anomaly_value)
                
                # ✅ Extraer nuevos campos de explicabilidad
                most_influential_variable = result.get("most_influential_variable")
                contribution_magnitude = result.get("contribution_magnitude")
                
                # Crear modelo de anomalía
                anomaly = AnomalyCreate(
                    activo_id=activo_id,
                    sensor_id=None,  # NULL para análisis a nivel de activo
                    timestamp=timestamp,
                    anomaly_score=float(anomaly_score),
                    anomaly_likelihood=float(anomaly_likelihood),
                    severidad=severidad,
                    descripcion=descripcion,
                    threshold=float(threshold),
                    is_anomaly=is_anomaly_value,
                    most_influential_variable=most_influential_variable,      # ✅ Nuevo
                    contribution_magnitude=float(contribution_magnitude) if contribution_magnitude else None  # ✅ Nuevo
                )
                
                anomalies_to_create.append(anomaly)
                
            except Exception as e:
                logger.error(f"❌ Error preparando anomalía: {e}", exc_info=True)
                continue
        
        # ✅ Validar que se crearon registros
        if not anomalies_to_create:
            logger.warning(
                f"⚠️  No se crearon registros válidos para guardar. "
                f"Total procesados: {len(results)}, válidos: 0"
            )
            return {
                "total": len(results),
                "saved": 0,
                "anomalies": 0,
                "failed": len(results)
            }
        
        # ✅ Guardar en batch
        try:
            saved_count = self.repository.bulk_create(anomalies_to_create)
            
            # ✅ Contar anomalías reales detectadas
            anomaly_count = sum(1 for a in anomalies_to_create if a.is_anomaly is True)
            
            logger.info(
                f"💾 ✅ Guardados {saved_count}/{len(results)} resultados ML para activo {activo_id}, "
                f"🚨 {anomaly_count} anomalías detectadas"
            )
            
            # ✅ Log de anomalías detectadas
            if anomaly_count > 0:
                logger.warning(f"🚨 ANOMALÍAS DETECTADAS: {anomaly_count} de {saved_count} registros")
                anomalies_only = [a for a in anomalies_to_create if a.is_anomaly]
                for i, anom in enumerate(anomalies_only[:5]):  # Mostrar primeras 5 anomalías
                    logger.warning(
                        f"  🚨 Anomalía {i+1}: {anom.timestamp.isoformat()} - "
                        f"Likelihood={anom.anomaly_likelihood:.2f} > Threshold={anom.threshold:.2f}, "
                        f"Severidad={anom.severidad}, "
                        f"Descripción={anom.descripcion}, "
                        f"Most Infl. Var={anom.most_influential_variable}, "
                        f"Contrib. Mag={anom.contribution_magnitude}"
                    )
            
            return {
                "total": len(results),
                "saved": saved_count,
                "anomalies": anomaly_count,
                "failed": len(results) - len(anomalies_to_create)
            }
            
        except Exception as e:
            logger.error(f"❌ Error guardando resultados en BD: {e}", exc_info=True)
            return {
                "total": len(results),
                "saved": 0,
                "anomalies": 0,
                "failed": len(results),
                "error": str(e)
            }

    async def detect_and_store_anomalies(
        self,
        activo_id: int,
        page: int = 1,
        limit: int = 1000
    ) -> Dict[str, Any]:
        """
        Flujo completo de detección de anomalías:
        1. Obtener datos del ParserService
        2. Transformar al formato ML
        3. Llamar al ML Engine
        4. Guardar resultados en BD
        
        Args:
            activo_id: ID del activo a analizar
            page: Página de datos a obtener
            limit: Cantidad de registros
            
        Returns:
            Dict con el resultado de la operación
        """
        logger.info(f"🔍 Iniciando detección de anomalías para activo {activo_id}")
        
        try:
            # 1. Obtener datos del ParserService
            parser_data = await self.fetch_data_from_parser(
                activo_id=activo_id,
                page=page,
                limit=limit
            )
            
            if not parser_data.sensores:
                logger.warning(f"⚠️  No hay datos disponibles para activo {activo_id}")
                return {
                    "activo_id": activo_id,
                    "status": "no_data",
                    "message": "No hay datos disponibles"
                }
            
            # 2. Transformar datos
            ml_records = self.transform_to_ml_format(parser_data)
            
            if not ml_records:
                logger.warning(f"⚠️  No se pudieron transformar los datos")
                return {
                    "activo_id": activo_id,
                    "status": "transform_error",
                    "message": "Error transformando datos"
                }
            
            # 3. Predecir anomalías
            ml_response = await self.predict_anomalies(activo_id, ml_records)
            
            # 4. Guardar resultados
            stats = self.save_ml_results(activo_id, ml_response)
            
            logger.info(f"✅ Detección completada para activo {activo_id}")
            
            return {
                "activo_id": activo_id,
                "status": "success",
                "data_points": len(ml_records),
                "results_saved": stats["saved"],
                "anomalies_detected": stats["anomalies"]
            }
            
        except Exception as e:
            logger.error(f"❌ Error en detección de anomalías: {e}")
            return {
                "activo_id": activo_id,
                "status": "error",
                "message": str(e)
            }
