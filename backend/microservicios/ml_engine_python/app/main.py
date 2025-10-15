from fastapi import FastAPI, Request, HTTPException
from fastapi.responses import JSONResponse
from pydantic import BaseModel, Field
from typing import List, Dict, Any
import pandas as pd
import logging

from . import config, logger
from .model_htm import detectar_htm_multivar
from .utils import clean_timestamps

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

class SensorRecord(BaseModel):
    """Registro individual de sensor con timestamp y features"""
    timestamp: str
    # Los demás campos son dinámicos (temperatura, presion, vibration, etc.)
    
    class Config:
        extra = "allow"  # Permite campos adicionales


class PumpWindow(BaseModel):
    """Ventana de datos para un activo/bomba"""
    pump: str = Field(..., description="Identificador del activo/bomba")
    records: List[Dict[str, Any]] = Field(..., description="Lista de registros de sensores")
    
    class Config:
        schema_extra = {
            "example": {
                "pump": "A",
                "records": [
                    {
                        "timestamp": "2025-10-09T22:40:00Z",
                        "A_ACR_Mot.PV": 0.012,
                        "A_ACR_Mot.SV": 0.035,
                        "A_ACR_Mot.TV": 41.2,
                        "A_ACR_Pmp.PV": 0.009,
                        "A_ACR_Pmp.SV": 0.031,
                        "A_ACR_Pmp.TV": 42.0,
                        "A_Pres.PV": 1.84,
                        "A_Temp.PV": 22.7,
                        "Barometer": 1013.3,
                        "Temperature": 22.5
                    },
                    {
                        "timestamp": "2025-10-09T22:40:01Z",
                        "A_ACR_Mot.PV": 0.011,
                        "A_ACR_Mot.SV": 0.034,
                        "A_ACR_Mot.TV": 41.3,
                        "A_ACR_Pmp.PV": 0.010,
                        "A_ACR_Pmp.SV": 0.032,
                        "A_ACR_Pmp.TV": 42.1,
                        "A_Pres.PV": 1.86,
                        "A_Temp.PV": 22.8,
                        "Barometer": 1013.4,
                        "Temperature": 22.6
                    }
                    # Se requieren al menos 180 registros para análisis completo
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
async def predict_anomaly(payload: PumpWindow):
    """
    Recibe un bloque de datos históricos y devuelve análisis de anomalías.
    
    - **pump**: Identificador del activo/bomba
    - **records**: Lista de registros con timestamp y features de sensores
    
    **Retorna**:
    - pump: Identificador del activo
    - n_rows: Número de filas analizadas
    - n_anomalies: Número de anomalías detectadas
    - results: Array con cada punto analizado (AnomalyScore, AnomalyLikelihood, Threshold, Severity, etc.)
    - last: Último resultado (más reciente)
    """
    try:
        logger.info(f"Procesando predicción para activo: {payload.pump}")
        logger.debug(f"Número de registros recibidos: {len(payload.records)}")
        
        # Validación básica
        if not payload.pump:
            raise HTTPException(
                status_code=400, 
                detail="El campo 'pump' es requerido"
            )
        
        if not payload.records or len(payload.records) == 0:
            raise HTTPException(
                status_code=400,
                detail="Se requiere al menos un registro en 'records'"
            )
        
        # Convertir a DataFrame
        df = pd.DataFrame(payload.records)
        
        # Validar que tenga timestamp
        if 'timestamp' not in df.columns:
            raise HTTPException(
                status_code=400,
                detail="Cada registro debe incluir un campo 'timestamp'"
            )
        
        # Limpiar timestamps con estrategia adaptativa
        try:
            df_clean, cleaning_stats = clean_timestamps(df, max_invalid_percent=5.0)
        except ValueError as ve:
            # Error de calidad de datos
            raise HTTPException(
                status_code=400,
                detail=str(ve)
            )
        
        # Log de calidad de datos (solo interno)
        if cleaning_stats['invalid_timestamps'] > 0:
            logger.info(
                f"Limpieza de datos aplicada - Estrategia: {cleaning_stats['strategy']}, "
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
        feature_cols = [c for c in df_clean.columns if c != 'timestamp']
        if len(feature_cols) == 0:
            raise HTTPException(
                status_code=400,
                detail="Se requiere al menos una columna de features (A_ACR_Mot.PV, A_Temp.PV, etc.)"
            )
        
        logger.debug(f"Features detectadas: {len(feature_cols)} columnas")
        
        # --- Ejecutar modelo HTM-like ---
        result = detectar_htm_multivar(df_clean, payload.pump)
        
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
    logger.info("🚀 Servicio ML Engine detenido")