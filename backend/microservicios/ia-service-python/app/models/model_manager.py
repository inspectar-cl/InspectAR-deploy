"""
Gestor del modelo de Machine Learning
Maneja carga, guardado y entrenamiento del modelo
"""

import os
import joblib
import logging
from typing import Optional, Tuple, Any
from sklearn.linear_model import SGDRegressor
from sklearn.preprocessing import StandardScaler

from app.config import config

logger = logging.getLogger(__name__)


class ModelManager:
    """Gestor del modelo de ML para detección de anomalías"""
    
    def __init__(self):
        self.model: Optional[SGDRegressor] = None
        self.scaler_X: Optional[StandardScaler] = None
        self.scaler_Y: Optional[StandardScaler] = None
        self.is_trained: bool = False  # ✅ NUEVO: Flag para saber si está entrenado
        self.threshold: float = 0.5  # ✅ NUEVO: Threshold adaptativo
        self.err_ref: float = 1.0    # ✅ NUEVO: Error de referencia (p95)
        self._ensure_model_dir()
    
    def _ensure_model_dir(self):
        """Crea el directorio de modelos si no existe"""
        os.makedirs(config.model_dir, exist_ok=True)
        logger.info(f"📁 Directorio de modelos: {config.model_dir}")
    
    def load_model(self) -> Tuple[bool, str]:
        """
        Carga el modelo y scalers desde disco si existen
        
        Returns:
            Tuple[bool, str]: (éxito, mensaje)
        """
        try:
            # Verificar que existan todos los archivos
            if not all(os.path.exists(p) for p in [
                config.model_path,
                config.scaler_x_path,
                config.scaler_y_path
            ]):
                logger.warning("⚠️  Archivos del modelo no encontrados, se creará uno nuevo")
                self._initialize_new_model()
                return False, "Modelo inicializado desde cero"
            
            # Cargar modelo y scalers
            self.model = joblib.load(config.model_path)
            self.scaler_X = joblib.load(config.scaler_x_path)
            self.scaler_Y = joblib.load(config.scaler_y_path)
            
            # ✅ NUEVO: Cargar threshold y err_ref
            threshold_path = os.path.join(config.model_dir, "threshold.pkl")
            err_ref_path = os.path.join(config.model_dir, "err_ref.pkl")
            
            if os.path.exists(threshold_path):
                self.threshold = joblib.load(threshold_path)
                logger.info(f"📊 Threshold cargado: {self.threshold:.4f}")
            else:
                self.threshold = 0.5
                logger.warning("⚠️  Threshold no encontrado, usando default: 0.5")
            
            if os.path.exists(err_ref_path):
                self.err_ref = joblib.load(err_ref_path)
                logger.info(f"📊 Error de referencia cargado: {self.err_ref:.4f}")
            else:
                self.err_ref = 1.0
                logger.warning("⚠️  Error de referencia no encontrado, usando default: 1.0")
            
            self.is_trained = True
            logger.info("✅ Modelo cargado exitosamente desde disco")
            return True, "Modelo cargado correctamente"
            
        except Exception as e:
            logger.error(f"❌ Error cargando modelo: {e}")
            logger.info("🔄 Inicializando modelo vacío")
            self._initialize_new_model()
            return False, f"Error al cargar, modelo reinicializado: {str(e)}"
    
    def _initialize_new_model(self):
        """Inicializa un nuevo modelo y scalers"""
        self.model = SGDRegressor(
            max_iter=1,
            learning_rate='constant',
            eta0=0.01,
            warm_start=True,
            random_state=42
        )
        self.scaler_X = StandardScaler()
        self.scaler_Y = StandardScaler()
        
        self.is_trained = False  # ✅ NUEVO: Modelo vacío NO está entrenado
        
        logger.info("🆕 Nuevo modelo SGDRegressor inicializado")
    
    def save_model(self) -> Tuple[bool, str]:
        """
        Guarda el modelo y scalers en disco
        
        Returns:
            Tuple[bool, str]: (éxito, mensaje)
        """
        try:
            if self.model is None:
                return False, "No hay modelo para guardar"
            
            # Guardar modelo y scalers
            joblib.dump(self.model, config.model_path)
            joblib.dump(self.scaler_X, config.scaler_x_path)
            joblib.dump(self.scaler_Y, config.scaler_y_path)
            
            # ✅ NUEVO: Guardar threshold y err_ref
            threshold_path = os.path.join(config.model_dir, "threshold.pkl")
            err_ref_path = os.path.join(config.model_dir, "err_ref.pkl")
            
            joblib.dump(self.threshold, threshold_path)
            joblib.dump(self.err_ref, err_ref_path)
            
            logger.info(f"💾 Modelo guardado en {config.model_path}")
            logger.info(f"💾 Threshold guardado: {self.threshold:.4f}")
            logger.info(f"💾 Error de referencia guardado: {self.err_ref:.4f}")
            
            self.is_trained = True
            return True, "Modelo guardado exitosamente"
            
        except Exception as e:
            logger.error(f"❌ Error guardando modelo: {e}")
            return False, f"Error al guardar: {str(e)}"
    
    def get_model_info(self) -> dict:
        """
        Retorna información sobre el modelo actual
        
        Returns:
            dict: Información del modelo
        """
        model_exists = os.path.exists(config.model_path)
        
        info = {
            "model_loaded": self.model is not None,
            "model_exists_on_disk": model_exists,
            "is_trained": self.is_trained,  # ✅ NUEVO: Incluir estado de entrenamiento
            "model_path": config.model_path,
            "model_type": type(self.model).__name__ if self.model else None,
            "threshold": round(self.threshold, 4),      # ✅ NUEVO
            "err_ref": round(self.err_ref, 4),          # ✅ NUEVO
        }
        
        if model_exists:
            try:
                # Información de archivos
                info["model_size_mb"] = round(
                    os.path.getsize(config.model_path) / (1024 * 1024), 2
                )
                info["last_modified"] = os.path.getmtime(config.model_path)
            except Exception as e:
                logger.warning(f"No se pudo obtener info del archivo: {e}")
        
        return info


# Instancia global del gestor de modelos
model_manager = ModelManager()
