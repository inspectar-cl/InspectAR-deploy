"""
Servicio de entrenamiento del modelo ML
Gestiona el entrenamiento periódico y bajo demanda
"""

import logging
import numpy as np
import httpx
from datetime import datetime
from typing import Dict, Any
from apscheduler.schedulers.asyncio import AsyncIOScheduler
from apscheduler.triggers.interval import IntervalTrigger
from collections import defaultdict

from app.config import config
from app.models.model_manager import model_manager
from app.models.anomaly import ParserResponse

logger = logging.getLogger(__name__)


class TrainingService:
    """Servicio para entrenar el modelo ML periódicamente"""
    
    def __init__(self):
        self.scheduler: AsyncIOScheduler = None
        self.training_count = 0
        self.last_training_time = None
        self.last_training_status = None
        self.is_training = False
    
    async def train_model(self) -> Dict[str, Any]:
        """
        Entrena el modelo con datos reales del IOT Service
        
        Returns:
            Dict con el resultado del entrenamiento
        """
        if self.is_training:
            logger.warning("⚠️  Ya hay un entrenamiento en progreso")
            return {
                "status": "skipped",
                "message": "Entrenamiento ya en progreso"
            }
        
        self.is_training = True
        self.training_count += 1
        start_time = datetime.now()
        
        try:
            logger.info("=" * 70)
            logger.info(f"🚀 INICIANDO ENTRENAMIENTO #{self.training_count}")
            logger.info("=" * 70)
            
            # Fase 1: Obtener datos del IOT Service
            logger.info("📊 Fase 1: Obtención de datos")
            activo_id = 2
            url = f"{config.iot_service_url}/lectura/{activo_id}/window"
            
            async with httpx.AsyncClient(timeout=30.0) as client:
                response = await client.get(url, params={"page": 1, "limit": 500})
                response.raise_for_status()
                parser_response = ParserResponse(**response.json())
            
            # Transformar a formato ML
            records_by_time = defaultdict(dict)
            for sensor in parser_response.sensores:
                if sensor.datos:
                    for dp in sensor.datos:
                        records_by_time[dp.tiempo][sensor.id_sensor] = dp.valor
            
            ml_records = [
                {"timestamp": t, **vals}
                for t, vals in records_by_time.items()
                if len(vals) >= 3
            ]
            
            if len(ml_records) < 50:
                logger.warning(f"⚠️  Datos insuficientes: {len(ml_records)}")
                return {
                    "status": "insufficient_data",
                    "records_found": len(ml_records)
                }
            
            logger.info(f"   ✓ {len(ml_records)} registros obtenidos")
            
            # Fase 2: Preprocesamiento
            logger.info("🔧 Fase 2: Preprocesamiento")
            X_list = []
            y_list = []
            
            for rec in ml_records:
                vals = [v for k, v in rec.items() if k != 'timestamp']
                if vals:
                    X_list.append(vals)
                    y_list.append(np.mean(vals))
            
            max_feat = max(len(x) for x in X_list)
            X = np.array([
                x + [0.0] * (max_feat - len(x)) if len(x) < max_feat else x[:max_feat]
                for x in X_list
            ])
            y = np.array(y_list)
            
            logger.info(f"   ✓ X: {X.shape}, y: {y.shape}")
            
            # Fase 3: Entrenamiento
            logger.info("🤖 Fase 3: Entrenamiento")
            if model_manager.model is None:
                model_manager._initialize_new_model()
            
            model_manager.scaler_X.fit(X)
            model_manager.scaler_Y.fit(y.reshape(-1, 1))
            
            X_scaled = model_manager.scaler_X.transform(X)
            y_scaled = model_manager.scaler_Y.transform(y.reshape(-1, 1)).ravel()
            
            model_manager.model.partial_fit(X_scaled, y_scaled)
            logger.info("   ✓ Modelo entrenado")
            
            # Fase 4: Evaluación
            logger.info("📈 Fase 4: Evaluación")
            pred_scaled = model_manager.model.predict(X_scaled)
            pred = model_manager.scaler_Y.inverse_transform(pred_scaled.reshape(-1, 1)).ravel()
            
            mse = float(np.mean((pred - y) ** 2))
            mae = float(np.mean(np.abs(pred - y)))
            r2 = float(1 - np.sum((y - pred) ** 2) / np.sum((y - np.mean(y)) ** 2))
            
            logger.info(f"   ✓ MSE: {mse:.4f}, MAE: {mae:.4f}, R²: {r2:.4f}")
            
            # Fase 5: Guardar modelo
            logger.info("💾 Fase 5: Guardando modelo")
            success, _ = model_manager.save_model()
            
            end_time = datetime.now()
            duration = (end_time - start_time).total_seconds()
            
            logger.info("=" * 70)
            logger.info(f"✅ COMPLETADO EN {duration:.2f}s")
            logger.info("=" * 70)
            
            self.last_training_time = end_time
            self.last_training_status = "success"
            
            return {
                "status": "success",
                "training_number": self.training_count,
                "duration_seconds": round(duration, 2),
                "records_used": len(X),
                "features_count": X.shape[1],
                "mse": round(mse, 4),
                "mae": round(mae, 4),
                "r2_score": round(r2, 4),
                "model_saved": success,
                "message": "Entrenamiento completado con datos del IOT Service"
            }
            
        except Exception as e:
            logger.error(f"❌ Error: {e}", exc_info=True)
            self.last_training_status = "error"
            return {
                "status": "error",
                "error": str(e)
            }
        finally:
            self.is_training = False
    
    def start_scheduler(self):
        """Inicia el scheduler"""
        if not config.training_enabled:
            return
        
        self.scheduler = AsyncIOScheduler()
        self.scheduler.add_job(
            self.train_model,
            trigger=IntervalTrigger(minutes=config.training_interval_minutes),
            id='train_model_job',
            max_instances=1
        )
        self.scheduler.start()
        logger.info("🔄 Scheduler iniciado")
    
    def stop_scheduler(self):
        """Detiene el scheduler"""
        if self.scheduler:
            self.scheduler.shutdown()
    
    def get_training_stats(self) -> Dict[str, Any]:
        """Retorna estadísticas"""
        return {
            "training_enabled": config.training_enabled,
            "total_trainings": self.training_count,
            "last_training_status": self.last_training_status
        }


training_service = TrainingService()
