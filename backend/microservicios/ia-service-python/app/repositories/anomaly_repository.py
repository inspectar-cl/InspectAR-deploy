"""
Repositorio para operaciones de base de datos de anomalías
Capa de acceso a datos (DAL)
"""

from sqlalchemy.orm import Session
from sqlalchemy import desc, and_, func
from typing import List, Optional, Tuple
import logging

from app.models.anomaly import AnomalyDB, AnomalyCreate

logger = logging.getLogger(__name__)


class AnomalyRepository:
    """Repositorio para operaciones CRUD de anomalías"""
    
    def __init__(self, db: Session):
        self.db = db
    
    def create(self, anomaly: AnomalyCreate) -> AnomalyDB:
        """Crea una nueva anomalía en la base de datos"""
        try:
            db_anomaly = AnomalyDB(**anomaly.model_dump())
            self.db.add(db_anomaly)
            self.db.commit()
            self.db.refresh(db_anomaly)
            
            logger.debug(
                f"💾 Anomalía guardada: ID={db_anomaly.id}, "
                f"Activo={db_anomaly.activo_id}, "
                f"IsAnomaly={db_anomaly.is_anomaly}, "
                f"Score={db_anomaly.anomaly_score:.2f}"
            )
            
            return db_anomaly
            
        except Exception as e:
            self.db.rollback()
            logger.error(f"❌ Error guardando anomalía: {e}")
            raise
    
    def bulk_create(self, anomalies: List[AnomalyCreate]) -> int:
        """Crea múltiples anomalías en una sola transacción"""
        try:
            db_anomalies = [AnomalyDB(**a.model_dump()) for a in anomalies]
            self.db.bulk_save_objects(db_anomalies)
            self.db.commit()
            
            logger.info(f"💾 Guardadas {len(db_anomalies)} anomalías en batch")
            
            return len(db_anomalies)
            
        except Exception as e:
            self.db.rollback()
            logger.error(f"❌ Error en bulk create: {e}")
            raise
    
    def get_by_id(self, anomaly_id: int) -> Optional[AnomalyDB]:
        """Obtiene una anomalía por ID"""
        try:
            return self.db.query(AnomalyDB).filter(
                AnomalyDB.id == anomaly_id
            ).first()
        except Exception as e:
            logger.error(f"❌ Error obteniendo anomalía por ID: {e}")
            raise
    
    def get_by_activo(
        self,
        activo_id: int,
        limit: int = 50,
        offset: int = 0,
        only_anomalies: bool = False
    ) -> Tuple[List[AnomalyDB], int]:
        """
        Obtiene anomalías por activo con paginación
        
        Returns:
            Tuple[List[AnomalyDB], int]: (lista de anomalías, total de registros)
        """
        try:
            # Query base
            query = self.db.query(AnomalyDB).filter(
                AnomalyDB.activo_id == activo_id
            )
            
            # Filtrar solo anomalías si se solicita
            if only_anomalies:
                query = query.filter(AnomalyDB.is_anomaly == 1)
            
            # Total de registros
            total = query.count()
            
            # Aplicar paginación y ordenamiento
            anomalies = query.order_by(
                desc(AnomalyDB.timestamp)
            ).offset(offset).limit(limit).all()
            
            logger.info(
                f"📊 Obtenidas {len(anomalies)} anomalías del activo {activo_id} "
                f"(total: {total}, limit: {limit}, offset: {offset})"
            )
            
            return anomalies, total
            
        except Exception as e:
            logger.error(f"❌ Error obteniendo anomalías por activo: {e}")
            raise
    
    def get_by_sensor(
        self,
        sensor_id: str,
        limit: int = 10,
        offset: int = 0
    ) -> Tuple[List[AnomalyDB], int]:
        """Obtiene anomalías por sensor específico"""
        try:
            query = self.db.query(AnomalyDB).filter(
                AnomalyDB.sensor_id == sensor_id
            )
            
            total = query.count()
            
            anomalies = query.order_by(
                desc(AnomalyDB.timestamp)
            ).offset(offset).limit(limit).all()
            
            logger.info(
                f"📊 Obtenidas {len(anomalies)} anomalías del sensor '{sensor_id}'"
            )
            
            return anomalies, total
            
        except Exception as e:
            logger.error(f"❌ Error obteniendo anomalías por sensor: {e}")
            raise
    
    def get_latest(
        self,
        activo_id: Optional[int] = None,
        only_anomalies: bool = True
    ) -> Optional[AnomalyDB]:
        """
        Obtiene la última anomalía detectada
        
        Args:
            activo_id: Filtrar por activo específico
            only_anomalies: Solo anomalías confirmadas (is_anomaly=1)
        """
        try:
            query = self.db.query(AnomalyDB)
            
            if only_anomalies:
                query = query.filter(AnomalyDB.is_anomaly == 1)
            
            if activo_id:
                query = query.filter(AnomalyDB.activo_id == activo_id)
            
            latest = query.order_by(desc(AnomalyDB.timestamp)).first()
            
            return latest
            
        except Exception as e:
            logger.error(f"❌ Error obteniendo última anomalía: {e}")
            raise
    
    def count(
        self,
        activo_id: Optional[int] = None,
        only_anomalies: bool = False
    ) -> int:
        """Cuenta el total de registros"""
        try:
            query = self.db.query(func.count(AnomalyDB.id))
            
            if activo_id:
                query = query.filter(AnomalyDB.activo_id == activo_id)
            
            if only_anomalies:
                query = query.filter(AnomalyDB.is_anomaly == 1)
            
            return query.scalar()
            
        except Exception as e:
            logger.error(f"❌ Error contando anomalías: {e}")
            raise
    
    def delete_by_activo(self, activo_id: int) -> int:
        """Elimina todas las anomalías de un activo"""
        try:
            deleted = self.db.query(AnomalyDB).filter(
                AnomalyDB.activo_id == activo_id
            ).delete()
            
            self.db.commit()
            
            logger.info(f"🗑️  Eliminadas {deleted} anomalías del activo {activo_id}")
            
            return deleted
            
        except Exception as e:
            self.db.rollback()
            logger.error(f"❌ Error eliminando anomalías: {e}")
            raise
