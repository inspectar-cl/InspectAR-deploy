"""
Servicio de entrenamiento del modelo ML
Gestiona el entrenamiento periódico y bajo demanda
"""

import logging
import numpy as np
import httpx
from datetime import datetime, timedelta
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
            
            errors = pred - y
            abs_errors = np.abs(errors)
            
            # Métricas básicas
            mse = float(np.mean(errors ** 2))
            mae = float(np.mean(abs_errors))
            rmse = float(np.sqrt(mse))
            r2 = float(1 - np.sum(errors ** 2) / np.sum((y - np.mean(y)) ** 2))
            
            '''# ✅ NUEVO: Threshold adaptativo (μ + 3σ estilo ml_engine)
            threshold_mean = float(np.mean(abs_errors))
            threshold_std = float(np.std(abs_errors))
            threshold_adaptive = threshold_mean + 3.0 * threshold_std
            
            # ✅ NUEVO: Error de referencia para normalización (percentil 95)
            err_ref = float(np.percentile(abs_errors, 95))
            
            # Coeficientes de variación
            cv_errors = float(np.std(errors) / threshold_mean * 100) if threshold_mean > 0 else 0.0
            cv_pred = float(np.std(pred) / np.mean(pred) * 100) if np.mean(pred) > 0 else 0.0
            cv_actual = float(np.std(y) / np.mean(y) * 100) if np.mean(y) > 0 else 0.0
            '''
            logger.info(f"   ✓ MSE: {mse:.4f}, MAE: {mae:.4f}, RMSE: {rmse:.4f}, R²: {r2:.4f}")
            '''logger.info(f"   ✓ CV Errores: {cv_errors:.2f}%, CV Pred: {cv_pred:.2f}%, CV Real: {cv_actual:.2f}%")
            logger.info(f"   ✓ Threshold Adaptativo: {threshold_adaptive:.4f} (μ={threshold_mean:.4f}, σ={threshold_std:.4f})")
            logger.info(f"   ✓ Error de Referencia (p95): {err_ref:.4f}")
            '''
            '''# ✅ NUEVO: Guardar threshold y err_ref en el modelo
            model_manager.threshold = threshold_adaptive
            model_manager.err_ref = err_ref'''
            
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
                "rmse": round(rmse, 4),
                "r2_score": round(r2, 4),
                '''"cv_errors_percent": round(cv_errors, 2),
                "cv_predictions_percent": round(cv_pred, 2),
                "cv_actual_percent": round(cv_actual, 2),
                "threshold_adaptive": round(threshold_adaptive, 4),  # ✅ NUEVO
                "threshold_mean": round(threshold_mean, 4),          # ✅ NUEVO
                "threshold_std": round(threshold_std, 4),            # ✅ NUEVO
                "err_ref_p95": round(err_ref, 4),                    # ✅ NUEVO
            '''    "model_saved": success,
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
            logger.info("⏸️  Entrenamiento deshabilitado en configuración")
            return
        
        # Configuración robusta del scheduler
        job_defaults = {
            'coalesce': True,
            'max_instances': 1,
            'misfire_grace_time': 300
        }
        
        self.scheduler = AsyncIOScheduler(job_defaults=job_defaults)
        
        # Calcular cuándo debe ser el primer entrenamiento
        if not model_manager.is_trained:
            # Si no hay modelo, entrenar en 2 segundos
            first_run = datetime.now() + timedelta(seconds=2)
            logger.info("🎓 Modelo no entrenado. Primer entrenamiento en 2 segundos")
        else:
            # Si hay modelo, esperar el intervalo completo
            first_run = datetime.now() + timedelta(minutes=config.training_interval_minutes)
            logger.info(f"✅ Modelo entrenado. Próximo entrenamiento en {config.training_interval_minutes} min")
        
        # UN SOLO JOB con tiempo de inicio calculado
        self.scheduler.add_job(
            self.train_model,
            trigger=IntervalTrigger(
                minutes=config.training_interval_minutes,
                start_date=first_run
            ),
            id='train_model_job',
            max_instances=1,
            replace_existing=True
        )
        
        self.scheduler.start()
        logger.info(f"🔄 Scheduler iniciado (cada {config.training_interval_minutes} minutos)")
        logger.info(f"📅 Primera ejecución: {first_run.strftime('%H:%M:%S')}")
        logger.info(f"📅 Próximo entrenamiento periódico: ~{config.training_interval_minutes} minutos")
    
    def stop_scheduler(self):
        """Detiene el scheduler"""
        if self.scheduler:
            self.scheduler.shutdown()
            logger.info("🛑 Scheduler detenido")
    
    def get_training_stats(self) -> Dict[str, Any]:
        """Retorna estadísticas"""
        return {
            "training_enabled": config.training_enabled,
            "interval_minutes": config.training_interval_minutes,
            "total_trainings": self.training_count,
            "last_training_status": self.last_training_status
        }


training_service = TrainingService()
