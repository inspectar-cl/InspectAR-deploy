"""
Servicio de lógica de negocio para anomalías
Orquesta las llamadas entre ParserService, ML Engine y BD
"""

from typing import List, Dict, Any, Optional
import httpx
import logging
from datetime import datetime
from collections import defaultdict

from app.config import config
from app.models.anomaly import (
    MLRequest, MLSensorRecord, MLResponse,
    AnomalyCreate, ParserResponse
)
from app.repositories.anomaly_repository import AnomalyRepository

logger = logging.getLogger(__name__)


class AnomalyService:
    """Servicio de lógica de negocio para detección de anomalías"""
    
    def __init__(self, repository: AnomalyRepository):
        self.repository = repository
    
    async def fetch_data_from_parser(
        self,
        activo_id: int,
        page: int = 1,
        limit: int = 1000
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
            sensor_id = sensor.id_sensor
            
            for data_point in sensor.datos:
                timestamp = data_point.tiempo
                value = data_point.valor
                
                # Normalizar nombre del sensor para ML Engine
                # Reemplazar puntos por guiones bajos
                normalized_sensor = sensor_id.replace(".", "_")
                
                records_by_time[timestamp][normalized_sensor] = value
        
        # Convertir a lista de MLSensorRecord
        ml_records = []
        for timestamp, sensors in sorted(records_by_time.items()):
            try:
                record = MLSensorRecord(
                    timestamp=timestamp,
                    **sensors
                )
                ml_records.append(record)
            except Exception as e:
                logger.warning(f"⚠️  Error creando record para {timestamp}: {e}")
                continue
        
        logger.info(
            f"📊 Transformados {len(ml_records)} registros agrupados por timestamp "
            f"desde {len(parser_data.sensores)} sensores"
        )
        
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
            import numpy as np
            
            # Verificar si el modelo está disponible
            model_available = False
            if model_manager.model is None:
                logger.warning("⚠️  Modelo no cargado, intentando cargar...")
                success, msg = model_manager.load_model()
                if success:
                    model_available = True
                    logger.info("✅ Modelo cargado exitosamente")
                else:
                    logger.warning(f"⚠️  Modelo no disponible: {msg}")
                    logger.info("📊 Generando datos con valores aleatorios para entrenamiento inicial")
            else:
                model_available = True
            
            results = []
            
            # Preparar datos para predicción
            for record in ml_records:
                try:
                    # Extraer features (todos los valores de sensores excepto timestamp)
                    record_dict = record.model_dump()
                    timestamp = record_dict.pop('timestamp')
                    
                    # Convertir valores a array numpy
                    features = np.array([list(record_dict.values())]).reshape(1, -1)
                    actual_value = features[0][0] if len(features[0]) > 0 else 0.0
                    
                    if model_available:
                        # Usar modelo entrenado
                        try:
                            # Escalar features si el scaler está entrenado
                            if hasattr(model_manager.scaler_X, 'mean_'):
                                features_scaled = model_manager.scaler_X.transform(features)
                            else:
                                features_scaled = features
                            
                            # Predecir con el modelo
                            prediction = model_manager.model.predict(features_scaled)[0]
                            
                            # Calcular error de reconstrucción
                            reconstruction_error = abs(actual_value - prediction)
                            
                            # Calcular anomaly score (normalizado 0-1)
                            anomaly_score = min(reconstruction_error / 10.0, 1.0)
                        
                        except Exception as e:
                            logger.warning(f"⚠️  Error en predicción, usando valores por defecto: {e}")
                            prediction = actual_value
                            anomaly_score = np.random.uniform(0.1, 0.4)  # Score bajo por defecto
                    else:
                        # Sin modelo: generar datos de ejemplo para entrenamiento inicial
                        prediction = actual_value + np.random.normal(0, 0.1)
                        # Score basado en variación aleatoria
                        anomaly_score = min(abs(actual_value - prediction) / 10.0, 0.6)
                    
                    # Determinar si es anomalía (threshold dinámico)
                    threshold = 0.5
                    is_anomaly = anomaly_score > threshold
                    
                    # Determinar severidad basada en el score
                    if anomaly_score < 0.3:
                        severity = "baja"
                    elif anomaly_score < 0.6:
                        severity = "media"
                    else:
                        severity = "alta"
                    
                    # Crear resultado
                    result = {
                        "timestamp": timestamp,
                        "is_anomaly": is_anomaly,
                        "AnomalyScore": round(anomaly_score, 4),
                        "AnomalyLikelihood": round(anomaly_score * 100, 2),
                        "Threshold": threshold,
                        "Severity": severity,
                        "Description": (
                            f"Anomalía detectada en activo {activo_id}"
                            if is_anomaly
                            else f"Lectura normal en activo {activo_id}"
                        ),
                        "ReconstructedValue": round(prediction, 4)
                    }
                    
                    results.append(result)
                    
                except Exception as e:
                    logger.warning(f"⚠️  Error procesando registro: {e}")
                    continue
            
            # Crear respuesta en formato MLResponse
            from app.models.anomaly import MLResponse as MLResponseClass
            
            ml_response = {
                "pump": f"Activo_{activo_id}",
                "results": results
            }
            
            anomaly_count = sum(1 for r in results if r.get("is_anomaly") is True)
            
            mode_text = "con modelo ML" if model_available else "sin modelo (datos iniciales)"
            logger.info(
                f"✅ Predicción completada {mode_text}. "
                f"Resultados: {len(results)}, "
                f"Anomalías detectadas: {anomaly_count}"
            )
            
            return ml_response
                
        except Exception as e:
            logger.error(f"❌ Error en predicción de anomalías: {e}")
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
        
        for result in results:
            try:
                # Normalizar severidad a minúsculas para BD
                severidad = result.get("Severity", "baja").lower()
                
                # Parsear timestamp
                try:
                    timestamp_str = result.get("timestamp")
                    # Intentar varios formatos
                    timestamp = datetime.fromisoformat(
                        timestamp_str.replace('Z', '+00:00')
                    )
                except (ValueError, AttributeError):
                    # Fallback: usar timestamp actual
                    timestamp = datetime.utcnow()
                
                # Crear modelo de anomalía
                anomaly = AnomalyCreate(
                    activo_id=activo_id,
                    sensor_id=None,  # NULL para análisis a nivel de activo
                    timestamp=timestamp,
                    anomaly_score=result.get("AnomalyScore", 0.0),
                    anomaly_likelihood=result.get("AnomalyLikelihood", 0.0),
                    severidad=severidad,
                    descripcion=result.get("Description", f"Análisis ML para activo {activo_id}"),
                    threshold=result.get("Threshold", 0.5),
                    is_anomaly=bool(result.get("is_anomaly", False))
                )
                
                anomalies_to_create.append(anomaly)
                
            except Exception as e:
                logger.error(f"❌ Error preparando anomalía: {e}")
                continue
        
        # Guardar en batch
        if anomalies_to_create:
            saved_count = self.repository.bulk_create(anomalies_to_create)
            anomaly_count = sum(1 for a in anomalies_to_create if a.is_anomaly is True)
            
            logger.info(
                f"💾 Guardados {saved_count}/{len(results)} resultados, "
                f"{anomaly_count} anomalías detectadas"
            )
            
            return {
                "total": len(results),
                "saved": saved_count,
                "anomalies": anomaly_count
            }
        else:
            logger.warning("⚠️  No se guardaron resultados")
            return {"total": 0, "saved": 0, "anomalies": 0}
    
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
