from fastapi import FastAPI, Request, HTTPException
from fastapi.responses import JSONResponse
from pydantic import BaseModel, Field
from typing import List, Dict, Any, Optional, Union
import pandas as pd
import logging

from . import config, logger
from .model_htm import detectar_htm_multivar
from .utils import clean_timestamps, transform_sensores_to_records

# ============================================================
# Inicialización de FastAPI con configuración
# ============================================================

app = FastAPI(
    title=config.API_TITLE,
    description="Microservicio Python para predicción de anomalías en activos usando HTM-like",
    version=config.API_VERSION
)

# ============================================================
# Modelos de entrada
# ============================================================

class SensorDato(BaseModel):
    """Dato individual de un sensor"""
    tiempo: str
    valor: float

class SensorData(BaseModel):
    """Datos de un sensor específico"""
    sensor_id: str
    datos: List[SensorDato]


class UnifiedPayload(BaseModel):
    """
    Payload unificado que acepta ambos formatos:
    1. Formato flat (records): [{"timestamp": "...", "sensor1": val, ...}, ...]
    2. Formato agrupado (sensores): [{"sensor_id": "...", "datos": [...]}, ...]
    """
    pump: Union[int, str]
    edificio_id: Optional[int] = None
    estado: Optional[str] = None
    limit: Optional[int] = None
    offset: Optional[int] = None
    page: Optional[int] = None
    
    # Formato flat por timestamp
    records: Optional[List[Dict[str, Any]]] = None
    
    # Formato agrupado por sensor
    sensores: Optional[List[SensorData]] = None
    
    class Config:
        schema_extra = {
            "example_flat": {
                "pump": "A",
                "records": [
                    {
                        "timestamp": "2025-10-09T22:40:00Z",
                        "A_ACR_Mot.PV": 0.012,
                        "A_Temp.PV": 22.7
                    }
                ]
            },
            "example_grouped": {
                "pump": 2,
                "sensores": [
                    {
                        "sensor_id": "A_ACR_Mot.PV",
                        "datos": [
                            {"tiempo": "2024-06-11T15:59:59Z", "valor": 0.001735421}
                        ]
                    }
                ]
            }
        }


# ============================================================
# Endpoints
# ============================================================

@app.get("/")
def root():
    """Endpoint raíz con información del servicio"""
    return {
        "service": config.API_TITLE,
        "version": config.API_VERSION,
        "status": "running",
        "mode": config.MODE,
        "config": {
            "window_size": config.WINDOW_SIZE,
            "alpha": config.ALPHA,
            "k_adapt": config.K_ADAPT
        }
    }


@app.get("/healthz")
def health_check():
    """Health check endpoint"""
    return {
        "status": "ok",
        "service": "ml_engine",
        "mode": config.MODE
    }


@app.get("/config")
def get_config():
    """Retorna la configuración actual del modelo"""
    return {
        "window_size": config.WINDOW_SIZE,
        "alpha": config.ALPHA,
        "k_adapt": config.K_ADAPT,
        "score_scale": config.SCORE_SCALE,
        "min_periods": config.MIN_PERIODS,
        "rolling_window": config.ROLLING_WINDOW,
        "mode": config.MODE
    }


@app.post("/predict_anomaly")
async def predict_anomaly(payload: UnifiedPayload):
    """
    Recibe datos en cualquier formato (flat o agrupado) y devuelve análisis de anomalías.
    
    **Formatos soportados**:
    
    1. **Formato flat** (records):
    ```json
    {
        "pump": "A",
        "records": [
            {"timestamp": "2025-10-09T22:40:00Z", "A_ACR_Mot.PV": 0.012, "A_Temp.PV": 22.7},
            ...
        ]
    }
    ```
    
    2. **Formato agrupado por sensor** (sensores):
    ```json
    {
        "pump": 2,
        "sensores": [
            {
                "sensor_id": "A_ACR_Mot.PV",
                "datos": [{"tiempo": "2024-06-11T15:59:59Z", "valor": 0.001735421}]
            }
        ]
    }
    ```
    
    **Retorna**:
    - pump: Identificador del activo
    - n_rows: Número de filas analizadas
    - n_anomalies: Número de anomalías detectadas
    - results: Array con cada punto analizado (AnomalyScore, AnomalyLikelihood, Threshold, Severity, etc.)
    - last: Último resultado (más reciente)
    """
    try:
        logger.info(f"Procesando predicción para activo: {payload.pump}")
        
        # ============================================================
        # 1. DETECTAR Y TRANSFORMAR FORMATO DE ENTRADA
        # ============================================================
        
        records = None
        
        # Caso 1: Formato agrupado por sensores (tiene 'sensores')
        if payload.sensores is not None and len(payload.sensores) > 0:
            logger.info(f"📦 Formato detectado: AGRUPADO POR SENSORES ({len(payload.sensores)} sensores)")
            
            # Convertir a diccionario para transform_sensores_to_records
            data_dict = {
                "pump": payload.pump,
                "sensores": [sensor.dict() for sensor in payload.sensores]
            }
            
            # Transformar a formato flat
            transformed = transform_sensores_to_records(data_dict)
            records = transformed.get("records", [])
            
            logger.debug(f"✓ Transformación completada: {len(records)} registros generados")
        
        # Caso 2: Formato flat (tiene 'records')
        elif payload.records is not None and len(payload.records) > 0:
            logger.info(f"📋 Formato detectado: FLAT ({len(payload.records)} registros)")
            records = payload.records
            
            # Normalizar el nombre de la columna timestamp si es necesario
            if len(records) > 0:
                first_record = records[0]
                # Buscar timestamp en minúsculas o variantes
                timestamp_key = None
                for key in first_record.keys():
                    if key.lower() in ['timestamp', 'tiempo']:
                        timestamp_key = key
                        break
                
                if timestamp_key and timestamp_key != 'Timestamp':
                    logger.debug(f"Normalizando clave de timestamp: '{timestamp_key}' → 'Timestamp'")
                    records = [
                        {('Timestamp' if k == timestamp_key else k): v for k, v in record.items()}
                        for record in records
                    ]
        
        # Caso 3: No se proporcionó ningún formato válido
        else:
            raise HTTPException(
                status_code=400,
                detail="Se requiere 'records' (formato flat) o 'sensores' (formato agrupado)"
            )
        
        # ============================================================
        # 2. VALIDACIONES
        # ============================================================
        
        if not records or len(records) == 0:
            raise HTTPException(
                status_code=400,
                detail="No se pudieron generar registros válidos a partir de los datos proporcionados"
            )
        
        logger.debug(f"Número de registros a procesar: {len(records)}")
        
        # Convertir a DataFrame
        df = pd.DataFrame(records)
        
        # Debug: mostrar los primeros Timestamps recibidos
        if len(df) > 0 and 'Timestamp' in df.columns:
            logger.info(f"🔍 Primeros Timestamps recibidos (tipo: {type(df['Timestamp'].iloc[0])}): {df['Timestamp'].head(3).tolist()}")
        
        # Validar que tenga Timestamp
        if 'Timestamp' not in df.columns:
            raise HTTPException(
                status_code=400,
                detail="Cada registro debe incluir un campo 'Timestamp' o 'timestamp'"
            )
        
        # ============================================================
        # 3. LIMPIEZA DE TIMESTAMPS
        # ============================================================
        
        try:
            df_clean, cleaning_stats = clean_timestamps(df, max_invalid_percent=5.0)
        except ValueError as ve:
            # Error de calidad de datos
            raise HTTPException(
                status_code=400,
                detail=str(ve)
            )
        
        # Log de calidad de datos
        if cleaning_stats['invalid_timestamps'] > 0:
            logger.info(
                f"🧹 Limpieza de datos aplicada - Estrategia: {cleaning_stats['strategy']}, "
                f"Eliminados: {cleaning_stats['dropped_records']}, "
                f"Interpolados: {cleaning_stats['interpolated_records']}, "
                f"Calidad: {100 - cleaning_stats['invalid_percent']:.1f}%"
            )
        
        # Verificar si quedaron suficientes datos después de la limpieza
        if len(df_clean) < config.WINDOW_SIZE:
            raise HTTPException(
                status_code=400,
                detail=f"Después de limpiar datos, quedan {len(df_clean)} registros. "
                       f"Se requieren al menos {config.WINDOW_SIZE}."
            )
        
        # Verificar que haya al menos una columna de features
        feature_cols = [c for c in df_clean.columns if c != 'Timestamp']
        if len(feature_cols) == 0:
            raise HTTPException(
                status_code=400,
                detail="Se requiere al menos una columna de features (A_ACR_Mot.PV, A_Temp.PV, etc.)"
            )
        
        logger.debug(f"Features detectadas: {len(feature_cols)} columnas - {feature_cols[:5]}...")
        
        # ============================================================
        # 4. EJECUTAR MODELO HTM-like
        # ============================================================
        
        result = detectar_htm_multivar(df_clean, str(payload.pump))
        
        if isinstance(result, ValueError):
            raise HTTPException(
                status_code=400,
                detail=str(result)
            )
        
        logger.info(
            f"✅ Predicción completada para {payload.pump}: "
            f"{result['n_anomalies']} anomalías detectadas de {result['n_rows']} puntos"
        )
        
        return result
        
    except HTTPException:
        raise
    except ValueError as ve:
        logger.error(f"Error de validación: {str(ve)}")
        raise HTTPException(status_code=400, detail=str(ve))
    except Exception as e:
        logger.error(f"Error inesperado en predict_anomaly: {str(e)}", exc_info=True)
        raise HTTPException(
            status_code=500,
            detail=f"Error interno del servidor: {str(e)}"
        )


@app.post("/predict_anomaly_from_sensors")
async def predict_anomaly_from_sensors(payload: UnifiedPayload):
    """
    [DEPRECATED] Usa /predict_anomaly que ahora soporta ambos formatos.
    
    Este endpoint se mantiene por compatibilidad retroactiva.
    """
    logger.warning("⚠️ Endpoint deprecated: /predict_anomaly_from_sensors. Usa /predict_anomaly")
    return await predict_anomaly(payload)


# ============================================================
# Middleware de logging
# ============================================================

@app.middleware("http")
async def log_requests(request: Request, call_next):
    """Middleware para logging de requests"""
    logger.debug(f"Request: {request.method} {request.url.path}")
    response = await call_next(request)
    logger.debug(f"Response status: {response.status_code}")
    return response


# ============================================================
# Startup/Shutdown events
# ============================================================

@app.on_event("startup")
async def startup_event():
    logger.info("🚀 Servicio ML Engine iniciado correctamente")
    logger.info(f"   Puerto: {config.API_PORT}")
    logger.info(f"   Modo: {config.MODE}")
    logger.info(f"   Window Size: {config.WINDOW_SIZE}")
    logger.info(f"   Alpha: {config.ALPHA}")
    logger.info(f"   K_adapt: {config.K_ADAPT}")


@app.on_event("shutdown")
async def shutdown_event():
    logger.info("🛑 Servicio ML Engine detenido")