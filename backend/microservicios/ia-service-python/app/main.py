"""
IA Service Python - FastAPI Application
Microservicio de detección de anomalías con Machine Learning
"""

from fastapi import FastAPI, Query
from fastapi.middleware.cors import CORSMiddleware
from contextlib import asynccontextmanager
import logging
import sys

from app.config import config
from app.database import check_db_connection, init_db
from app.handlers.anomaly_handler import router as anomaly_router
from app.models.model_manager import model_manager
from app.services.training_service import training_service

# Configurar logging
logging.basicConfig(
    level=getattr(logging, config.log_level.upper()),
    format="[%(asctime)s] [%(levelname)s] %(name)s: %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
    handlers=[logging.StreamHandler(sys.stdout)]
)

logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """
    Lifecycle events para FastAPI
    Se ejecuta al iniciar y detener la aplicación
    """
    # Startup
    logger.info("🚀 Iniciando IA Service Python...")
    logger.info(f"   Versión: 1.0.0")
    logger.info(f"   Puerto: {config.port}")
    logger.info(f"   Base de datos: {config.database_host}:{config.database_port}")
    logger.info(f"   IOT Service: {config.iot_service_url}")
    logger.info(f"   Entrenamiento: {'Habilitado' if config.training_enabled else 'Deshabilitado'}")
    
    # Verificar conexión a BD
    if not check_db_connection():
        logger.error("❌ No se pudo conectar a la base de datos")
        sys.exit(1)
    
    # Inicializar BD (solo verifica, no crea tablas)
    init_db()
    
    # Cargar modelo de ML
    logger.info("🤖 Cargando modelo de Machine Learning...")
    success, message = model_manager.load_model()
    logger.info(f"   └─ {message}")
    
    # Iniciar scheduler de entrenamiento
    logger.info("⏰ Configurando scheduler de entrenamiento...")
    logger.info(f"   └─ Intervalo: {config.training_interval_minutes} minutos")
    training_service.start_scheduler()
    
    logger.info("✅ IA Service Python iniciado correctamente")
    
    yield
    
    # Shutdown
    logger.info("🛑 Deteniendo IA Service Python...")
    
    # Detener scheduler
    training_service.stop_scheduler()
    
    # Guardar modelo antes de cerrar
    logger.info("💾 Guardando modelo antes de cerrar...")
    model_manager.save_model()


# Crear aplicación FastAPI
app = FastAPI(
    title="IA Service Python",
    description="Microservicio de detección de anomalías con Machine Learning",
    version="1.0.0",
    lifespan=lifespan,
    docs_url="/docs",
    redoc_url="/redoc",
    openapi_url="/openapi.json"
)

# Configurar CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # En producción, especificar orígenes permitidos
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Registrar routers
app.include_router(anomaly_router)


# Root endpoint
@app.get("/")
async def root():
    """Endpoint raíz con información del servicio"""
    return {
        "service": "IA Service Python",
        "version": "1.0.0",
        "status": "running",
        "endpoints": {
            "health": "/health",
            "status": "/status",
            "docs": "/docs",
            "anomalies_by_activo": "/anomalies/activo/{activo_id}",
            "anomalies_by_sensor": "/anomalies/sensor/{sensor_id}",
            "detect": "/detect/{activo_id}",
            "training_status": "/training/status",
            "train_now": "/training/train",
            "model_info": "/training/model"
        }
    }


@app.get("/training/status")
async def get_training_status():
    """Obtiene el estado del entrenamiento y scheduler"""
    return training_service.get_training_stats()


@app.post("/training/train")
async def trigger_training(
    activo_id: int = Query(2, description="ID del activo a usar para entrenamiento"),
    page: int = Query(1, ge=1, description="Página a solicitar al IOT Service"),
    limit: int = Query(5000, ge=1, le=5000, description="Límite de registros a solicitar (max 5000)"),
):
    """Dispara un entrenamiento manual del modelo con parámetros configurables"""
    result = await training_service.train_model(activo_id=activo_id, page=page, limit=limit)
    return result


@app.get("/training/model")
async def get_model_info():
    """Obtiene información sobre el modelo actual"""
    return model_manager.get_model_info()


if __name__ == "__main__":
    import uvicorn
    
    uvicorn.run(
        "app.main:app",
        host=config.host,
        port=config.port,
        reload=config.debug,
        log_level=config.log_level.lower()
    )
