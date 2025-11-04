"""
Modelos de datos del servicio IA
"""

from app.models.anomaly import (
    AnomalyDB,
    AnomalyBase,
    AnomalyCreate,
    AnomalyResponse,
    AnomalyListResponse,
    MLRequest,
    MLResponse,
    MLResult
)

__all__ = [
    "AnomalyDB",
    "AnomalyBase",
    "AnomalyCreate",
    "AnomalyResponse",
    "AnomalyListResponse",
    "MLRequest",
    "MLResponse",
    "MLResult"
]
