"""
Modelos de datos para anomalías
"""

from sqlalchemy import Column, Integer, String, Float, Boolean, DateTime, Text
from sqlalchemy.sql import func
from pydantic import BaseModel, Field, ConfigDict
from datetime import datetime
from typing import Optional, List, Dict, Any

from app.database import Base


# ============================================================================
# MODELO SQLALCHEMY (Base de datos)
# ============================================================================

class AnomalyDB(Base):
    """Modelo de base de datos para anomalías"""
    __tablename__ = "anomalias"
    
    id = Column(Integer, primary_key=True, index=True)
    activo_id = Column(Integer, nullable=False, index=True)
    sensor_id = Column(String(100), nullable=True, index=True)
    timestamp = Column(DateTime(timezone=True), nullable=False, index=True)
    anomaly_score = Column(Float, nullable=False, default=0.0)
    anomaly_likelihood = Column(Float, nullable=False, default=0.0)
    severidad = Column(String(20), nullable=False, default='baja', index=True)
    descripcion = Column(Text, nullable=True)
    threshold = Column(Float, nullable=True)
    is_anomaly = Column(Boolean, nullable=False, default=False, index=True)
    most_influential_variable = Column(String(100), nullable=True, default=None)
    contribution_magnitude = Column(Float, nullable=True, default=None)
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())


# ============================================================================
# MODELOS PYDANTIC (API)
# ============================================================================

class AnomalyBase(BaseModel):
    """Modelo base de anomalía"""
    activo_id: int
    sensor_id: Optional[str] = None
    timestamp: datetime
    anomaly_score: float = 0.0
    anomaly_likelihood: float = 0.0
    severidad: str = "baja"
    descripcion: Optional[str] = None
    threshold: Optional[float] = None
    is_anomaly: bool = False
    most_influential_variable: Optional[str] = None
    contribution_magnitude: Optional[float] = None



class AnomalyCreate(AnomalyBase):
    """Modelo para crear anomalía"""
    pass


class AnomalyResponse(AnomalyBase):
    """Modelo de respuesta de anomalía"""
    id: int
    created_at: datetime
    updated_at: Optional[datetime] = None
    
    model_config = ConfigDict(from_attributes=True)


class AnomalyListResponse(BaseModel):
    """Respuesta de lista de anomalías"""
    message: str
    count: int
    activo_id: Optional[int] = None
    sensor_id: Optional[str] = None
    limit: int
    offset: int
    data: List[AnomalyResponse]


# ============================================================================
# MODELOS PARA ML ENGINE
# ============================================================================

class MLSensorRecord(BaseModel):
    """Registro de sensor para ML Engine"""
    timestamp: str
    A_ACR_Mot_PV: float = Field(alias="A_ACR_Mot.PV", default=0.0)
    A_ACR_Mot_SV: float = Field(alias="A_ACR_Mot.SV", default=0.0)
    A_Desc_Bom_PV: float = Field(alias="A_Desc_Bom.PV", default=0.0)
    A_Desc_Bom_SV: float = Field(alias="A_Desc_Bom.SV", default=0.0)
    Flujo_Caudal: float = 0.0
    Nivel_Presion: float = 0.0
    Potencia_Activa: float = 0.0
    Temp_Agua: float = 0.0
    Temp_Ambiente: float = 0.0
    Vibracion: float = 0.0
    
    model_config = ConfigDict(populate_by_name=True)


class MLRequest(BaseModel):
    """Request al ML Engine"""
    pump: str
    records: List[MLSensorRecord]


class MLResult(BaseModel):
    """Resultado individual del ML Engine"""
    timestamp: str
    AnomalyScore: float
    AnomalyLikelihood: float
    Threshold: float
    is_anomaly: int
    Severity: str
    Description: str
    most_influential_variable: str
    contribution_magnitude: float


class MLResponse(BaseModel):
    """Respuesta del ML Engine"""
    pump: str
    n_rows: Optional[int] = None
    n_anomalies: Optional[int] = None
    results: List[MLResult]
    last: Optional[Dict[str, Any]] = None


# ============================================================================
# MODELOS PARA PARSER SERVICE
# ============================================================================

class DataPoint(BaseModel):
    """Punto de dato del ParserService"""
    tiempo: str
    valor: float


class SensorData(BaseModel):
    """Datos de un sensor del ParserService"""
    id_sensor: str = Field(alias="sensor_id")
    datos: List[DataPoint]
    
    model_config = ConfigDict(populate_by_name=True)


class ParserResponse(BaseModel):
    """Respuesta del ParserService"""
    activo_id: int
    page: int
    limit: int
    offset: int
    sensores: List[SensorData]
