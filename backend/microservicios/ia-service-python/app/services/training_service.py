"""
Servicio de entrenamiento del modelo ML
Gestiona el entrenamiento periódico y bajo demanda
"""

import logging
from datetime import datetime
from typing import Dict, Any
from apscheduler.schedulers.asyncio import AsyncIOScheduler
from apscheduler.triggers.interval import IntervalTrigger

from app.config import config
from app.models.model_manager import model_manager

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
        Entrena el modelo con los datos disponibles
        
        Por ahora solo hace log, en el futuro implementará
        el entrenamiento real con datos de la base de datos
        
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
            logger.info(f"🚀 INICIANDO ENTRENAMIENTO DEL MODELO #{self.training_count}")
            logger.info(f"⏰ Timestamp: {start_time.isoformat()}")
            logger.info("=" * 70)
            
            # TODO: Implementar lógica real de entrenamiento
            # 1. Obtener datos de entrenamiento desde la BD
            # 2. Preparar features y targets
            # 3. Entrenar modelo con partial_fit()
            # 4. Evaluar performance
            # 5. Guardar modelo si mejora
            
            logger.info("📊 Fase 1: Recolección de datos")
            logger.info("   └─ TODO: Consultar base de datos ia_db")
            
            logger.info("🔧 Fase 2: Preprocesamiento")
            logger.info("   └─ TODO: Normalizar y transformar datos")
            
            logger.info("🤖 Fase 3: Entrenamiento incremental")
            logger.info("   └─ TODO: Ejecutar model.partial_fit()")
            
            logger.info("📈 Fase 4: Evaluación")
            logger.info("   └─ TODO: Calcular métricas de performance")
            
            logger.info("💾 Fase 5: Persistencia")
            logger.info("   └─ TODO: Guardar modelo si mejora")
            
            # Simulación de guardado del modelo
            success, message = model_manager.save_model()
            
            end_time = datetime.now()
            duration = (end_time - start_time).total_seconds()
            
            result = {
                "status": "success",
                "training_number": self.training_count,
                "start_time": start_time.isoformat(),
                "end_time": end_time.isoformat(),
                "duration_seconds": round(duration, 2),
                "model_saved": success,
                "message": "Entrenamiento completado (simulación)"
            }
            
            logger.info("=" * 70)
            logger.info(f"✅ ENTRENAMIENTO COMPLETADO EN {duration:.2f}s")
            logger.info(f"💾 Modelo guardado: {success}")
            logger.info("=" * 70)
            
            self.last_training_time = end_time
            self.last_training_status = "success"
            
            return result
            
        except Exception as e:
            logger.error(f"❌ Error durante el entrenamiento: {e}", exc_info=True)
            self.last_training_status = "error"
            return {
                "status": "error",
                "training_number": self.training_count,
                "error": str(e),
                "message": "Error durante el entrenamiento"
            }
        finally:
            self.is_training = False
    
    def start_scheduler(self):
        """Inicia el scheduler para entrenamiento periódico"""
        if not config.training_enabled:
            logger.info("⏸️  Entrenamiento automático deshabilitado")
            return
        
        self.scheduler = AsyncIOScheduler()
        
        # Agregar job de entrenamiento periódico
        self.scheduler.add_job(
            self.train_model,
            trigger=IntervalTrigger(minutes=config.training_interval_minutes),
            id='train_model_job',
            name='Entrenamiento Periódico del Modelo',
            replace_existing=True,
            max_instances=1  # Solo una instancia a la vez
        )
        
        self.scheduler.start()
        
        logger.info("=" * 70)
        logger.info("🔄 SCHEDULER DE ENTRENAMIENTO INICIADO")
        logger.info(f"⏱️  Intervalo: {config.training_interval_minutes} minutos")
        logger.info(f"📅 Próximo entrenamiento en: {config.training_interval_minutes} minutos")
        logger.info("=" * 70)
    
    def stop_scheduler(self):
        """Detiene el scheduler"""
        if self.scheduler:
            self.scheduler.shutdown()
            logger.info("🛑 Scheduler de entrenamiento detenido")
    
    def get_training_stats(self) -> Dict[str, Any]:
        """Retorna estadísticas del entrenamiento"""
        next_run = None
        if self.scheduler:
            job = self.scheduler.get_job('train_model_job')
            if job and job.next_run_time:
                next_run = job.next_run_time.isoformat()
        
        return {
            "training_enabled": config.training_enabled,
            "interval_minutes": config.training_interval_minutes,
            "total_trainings": self.training_count,
            "last_training_time": self.last_training_time.isoformat() if self.last_training_time else None,
            "last_training_status": self.last_training_status,
            "is_training": self.is_training,
            "next_training": next_run,
            "scheduler_running": self.scheduler is not None and self.scheduler.running
        }


# Instancia global del servicio de entrenamiento
training_service = TrainingService()
