"""
Handlers HTTP del servicio IA
"""

from app.handlers.anomaly_handler import router as anomaly_router

__all__ = ["anomaly_router"]
