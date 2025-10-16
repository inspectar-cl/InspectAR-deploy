import logging
from . import config

# ============================================================
# INICIALIZACIÓN DEL MÓDULO APP
# ============================================================

# Configuración básica de logs
logging.basicConfig(
    level=getattr(logging, config.LOG_LEVEL.upper(), logging.INFO),
    format="[%(asctime)s] [%(levelname)s] %(name)s: %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
)

logger = logging.getLogger("ml_service")

logger.info(" Iniciando microservicio ML (HTM-like)")
logger.info(f"  Modo: {config.MODE.upper()} | Ventana W={config.WINDOW_SIZE}, α={config.ALPHA}, k={config.K_ADAPT}")
logger.info(f"  Resultados guardados: {config.SAVE_RESULTS} en {config.OUTPUT_PATH}")

# ============================================================
# EXPORTS (para que otros módulos puedan acceder fácilmente)
# ============================================================

__all__ = [
    "config",
    "logger",
]
