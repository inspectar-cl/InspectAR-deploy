"""
Servicio de lógica de negocio para anomalías
Orquesta las llamadas entre ParserService, ML Engine y BD
"""

from typing import List, Dict, Any, Optional
import httpx
import logging
from datetime import datetime
from collections import defaultdict
import numpy as np
import pandas as pd

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
            feature_names = []  # Agregar lista para nombres
            
            for rec in ml_records:
                vals = [v for k, v in rec.items() if k != 'timestamp']
                if vals and not feature_names:  # Capturar nombres solo una vez
                    feature_names = [k for k in rec.keys() if k != 'timestamp']
                if vals:
                    X_list.append(vals)
                    y_list.append(np.mean(vals))
            
            # Imprimir nombres de features
            logger.info(f"   📋 Features detectadas: {feature_names}")
            logger.info(f"   📋 Cantidad de features: {len(feature_names)}")
            
            max_feat = max(len(x) for x in X_list)
            X = np.array([
                x + [0.0] * (max_feat - len(x)) if len(x) < max_feat else x[:max_feat]
                for x in X_list
            ])
            y = np.array(y_list)
            
            logger.info(f"   ✓ X: {X.shape}, y: {y.shape}")
            timestamps = [rec['timestamp'] for rec in ml_records]
            # ✅ LOG para verificar dimensiones
            logger.info(f"📊 Ventana preparada: shape={X.shape} (debe ser [n_samples, n_sensores])")
            logger.info(f"📊 Features por registro: {max_feat}")
            
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
                        "ReconstructedValue": round(float(prediction), 4)
                    }
                    
                    results.append(result)
                    
                except Exception as e:
                    logger.warning(f"⚠️  Error procesando registro {i}: {e}")
                    continue
            
            # ===================================
            # ✅ CALCULAR THRESHOLD ADAPTATIVO Y is_anomaly
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
                anomaly_mask = df_results['is_anomaly'] == True
                
                if anomaly_mask.any():
                    for idx in df_results[anomaly_mask].index:
                        likelihood = df_results.loc[idx, 'AnomalyLikelihood']
                        df_results.loc[idx, 'Severity'] = severity_from_likelihood(likelihood)
                        df_results.loc[idx, 'Description'] = description_from_severity(
                            df_results.loc[idx, 'Severity']
                        )
                
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
                
                # Crear modelo de anomalía
                anomaly = AnomalyCreate(
                    activo_id=activo_id,
                    sensor_id=None,  # NULL para análisis a nivel de activo
                    timestamp=timestamp,
                    anomaly_score=float(anomaly_score),  # Valor REAL del predict
                    anomaly_likelihood=float(anomaly_likelihood),  # Valor REAL del predict
                    severidad=severidad,  # Severidad REAL del predict
                    descripcion=descripcion,  # Descripción completa
                    threshold=float(threshold),  # Threshold REAL calculado con rolling
                    is_anomaly=is_anomaly_value  # is_anomaly REAL del predict
                )
                
                anomalies_to_create.append(anomaly)
                
            except Exception as e:
                logger.error(
                    f"❌ Error preparando anomalía para registro {idx} "
                    f"(timestamp: {result.get('timestamp', 'unknown')}): {e}",
                    exc_info=True
                )
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
            
            # ✅ Log de muestra de datos guardados (primeros 3)
            if saved_count > 0 and len(anomalies_to_create) > 0:
                sample_size = min(3, len(anomalies_to_create))
                logger.info(f"📝 Muestra de {sample_size} registros guardados:")
                for i, anomaly in enumerate(anomalies_to_create[:sample_size]):
                    logger.info(
                        f"  [{i+1}] {anomaly.timestamp.isoformat()} - "
                        f"is_anomaly={anomaly.is_anomaly}, "
                        f"Score={anomaly.anomaly_score:.2f}, "
                        f"Likelihood={anomaly.anomaly_likelihood:.2f}, "
                        f"Threshold={anomaly.threshold:.2f}, "
                        f"Severidad={anomaly.severidad}"
                    )
            
            # ✅ Log de anomalías detectadas
            if anomaly_count > 0:
                logger.warning(f"🚨 ANOMALÍAS DETECTADAS: {anomaly_count} de {saved_count} registros")
                anomalies_only = [a for a in anomalies_to_create if a.is_anomaly]
                for i, anom in enumerate(anomalies_only[:5]):  # Mostrar primeras 5 anomalías
                    logger.warning(
                        f"  🚨 Anomalía {i+1}: {anom.timestamp.isoformat()} - "
                        f"Likelihood={anom.anomaly_likelihood:.2f} > Threshold={anom.threshold:.2f}, "
                        f"Severidad={anom.severidad}"
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
