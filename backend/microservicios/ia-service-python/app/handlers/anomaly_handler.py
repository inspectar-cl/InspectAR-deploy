"""
Controladores HTTP para anomalías
Capa de presentación - Endpoints REST
"""

from fastapi import APIRouter, Depends, HTTPException, Query, BackgroundTasks
from sqlalchemy.orm import Session
from typing import Optional
import logging

from app.database import get_db
from app.models.anomaly import AnomalyListResponse, AnomalyResponse
from app.repositories.anomaly_repository import AnomalyRepository
from app.services.anomaly_service import AnomalyService

logger = logging.getLogger(__name__)

router = APIRouter(tags=["anomalies"])


# ============================================================================
# DEPENDENCY INJECTION
# ============================================================================

def get_repository(db: Session = Depends(get_db)) -> AnomalyRepository:
    """Dependency injection para el repositorio"""
    return AnomalyRepository(db)


def get_service(repo: AnomalyRepository = Depends(get_repository)) -> AnomalyService:
    """Dependency injection para el servicio"""
    return AnomalyService(repo)


# ============================================================================
# ENDPOINTS
# ============================================================================

@router.get("/anomalies/activo/{activo_id}", response_model=AnomalyListResponse)
async def get_anomalies_by_activo(
    activo_id: int,
    limit: int = Query(50, ge=1, le=1000, description="Número de resultados"),
    offset: int = Query(0, ge=0, description="Offset para paginación"),
    only_anomalies: bool = Query(False, description="Solo anomalías confirmadas"),
    repo: AnomalyRepository = Depends(get_repository)
):
    """
    Obtiene anomalías de un activo con paginación
    
    - **activo_id**: ID del activo
    - **limit**: Cantidad máxima de resultados (default: 50, max: 1000)
    - **offset**: Desplazamiento para paginación (default: 0)
    - **only_anomalies**: Solo retornar anomalías confirmadas (is_anomaly=1)
    """
    try:
        anomalies, total = repo.get_by_activo(
            activo_id=activo_id,
            limit=limit,
            offset=offset,
            only_anomalies=only_anomalies
        )
        
        return AnomalyListResponse(
            message="Anomalías obtenidas exitosamente",
            count=len(anomalies),
            activo_id=activo_id,
            limit=limit,
            offset=offset,
            data=[AnomalyResponse.model_validate(a) for a in anomalies]
        )
        
    except Exception as e:
        logger.error(f"❌ Error obteniendo anomalías por activo: {e}")
        raise HTTPException(
            status_code=500,
            detail=f"Error al obtener anomalías: {str(e)}"
        )


@router.get("/anomalies/sensor/{sensor_id}", response_model=AnomalyListResponse)
async def get_anomalies_by_sensor(
    sensor_id: str,
    limit: int = Query(10, ge=1, le=1000, description="Número de resultados"),
    offset: int = Query(0, ge=0, description="Offset para paginación"),
    repo: AnomalyRepository = Depends(get_repository)
):
    """
    Obtiene anomalías de un sensor específico
    
    - **sensor_id**: Identificador del sensor
    - **limit**: Cantidad máxima de resultados
    - **offset**: Desplazamiento para paginación
    """
    try:
        anomalies, total = repo.get_by_sensor(
            sensor_id=sensor_id,
            limit=limit,
            offset=offset
        )
        
        return AnomalyListResponse(
            message="Anomalías obtenidas exitosamente",
            count=len(anomalies),
            sensor_id=sensor_id,
            limit=limit,
            offset=offset,
            data=[AnomalyResponse.model_validate(a) for a in anomalies]
        )
        
    except Exception as e:
        logger.error(f"❌ Error obteniendo anomalías por sensor: {e}")
        raise HTTPException(
            status_code=500,
            detail=f"Error al obtener anomalías: {str(e)}"
        )


@router.get("/anomalies/{anomaly_id}", response_model=AnomalyResponse)
async def get_anomaly_by_id(
    anomaly_id: int,
    repo: AnomalyRepository = Depends(get_repository)
):
    """
    Obtiene una anomalía específica por ID
    
    - **anomaly_id**: ID de la anomalía
    """
    try:
        anomaly = repo.get_by_id(anomaly_id)
        
        if not anomaly:
            raise HTTPException(
                status_code=404,
                detail=f"Anomalía {anomaly_id} no encontrada"
            )
        
        return AnomalyResponse.model_validate(anomaly)
        
    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"❌ Error obteniendo anomalía: {e}")
        raise HTTPException(
            status_code=500,
            detail=f"Error al obtener anomalía: {str(e)}"
        )


@router.post("/detect/{activo_id}")
async def trigger_anomaly_detection(
    activo_id: int,
    background_tasks: BackgroundTasks,
    page: int = Query(1, ge=1, description="Página de datos"),
    limit: int = Query(1000, ge=1, le=5000, description="Cantidad de registros"),
    async_mode: bool = Query(False, description="Ejecutar en segundo plano"),
    service: AnomalyService = Depends(get_service)
):
    """
    Dispara manualmente la detección de anomalías para un activo
    
    - **activo_id**: ID del activo a analizar
    - **page**: Página de datos a procesar
    - **limit**: Cantidad de registros a analizar
    - **async_mode**: Si es True, ejecuta en segundo plano
    """
    try:
        if async_mode:
            # Ejecutar en background
            background_tasks.add_task(
                service.detect_and_store_anomalies,
                activo_id=activo_id,
                page=page,
                limit=limit
            )
            
            return {
                "message": "Detección de anomalías iniciada en segundo plano",
                "activo_id": activo_id,
                "status": "processing"
            }
        else:
            # Ejecutar sincrónicamente
            result = await service.detect_and_store_anomalies(
                activo_id=activo_id,
                page=page,
                limit=limit
            )
            return result
        
    except Exception as e:
        logger.error(f"❌ Error detectando anomalías: {e}")
        raise HTTPException(
            status_code=500,
            detail=f"Error al detectar anomalías: {str(e)}"
        )


@router.get("/status")
async def get_status(
    activo_id: Optional[int] = Query(None, description="Filtrar por activo"),
    repo: AnomalyRepository = Depends(get_repository)
):
    """
    Obtiene el estado del servicio y estadísticas
    
    - **activo_id**: (Opcional) Filtrar estadísticas por activo
    """
    try:
        latest = repo.get_latest(activo_id=activo_id)
        total_anomalies = repo.count(activo_id=activo_id, only_anomalies=True)
        total_records = repo.count(activo_id=activo_id, only_anomalies=False)
        
        return {
            "status": "healthy",
            "service": "ia-service-python",
            "activo_id": activo_id,
            "total_records": total_records,
            "total_anomalies": total_anomalies,
            "anomaly_rate": (
                round(total_anomalies / total_records * 100, 2)
                if total_records > 0 else 0
            ),
            "latest_anomaly": (
                AnomalyResponse.model_validate(latest) if latest else None
            )
        }
        
    except Exception as e:
        logger.error(f"❌ Error obteniendo estado: {e}")
        raise HTTPException(
            status_code=500,
            detail=f"Error al obtener estado: {str(e)}"
        )


@router.get("/health")
async def health_check():
    """
    Health check endpoint simple
    Usado por Docker healthcheck
    """
    return {
        "status": "healthy",
        "service": "ia-service-python",
        "version": "1.0.0"
    }
