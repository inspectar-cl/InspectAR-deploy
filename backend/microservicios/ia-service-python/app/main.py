"""
IA Service Python - FastAPI Application
Microservicio de detección de anomalías con Machine Learning
"""

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from contextlib import asynccontextmanager
import logging
import sys

from app.config import config
from app.database import check_db_connection, init_db
from app.handlers.anomaly_handler import router as anomaly_router

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
    logger.info(f"   ML Engine: {config.ml_engine_url}")
    logger.info(f"   IOT Service: {config.iot_service_url}")
    
    # Verificar conexión a BD
    if not check_db_connection():
        logger.error("❌ No se pudo conectar a la base de datos")
        sys.exit(1)
    
    # Inicializar BD (solo verifica, no crea tablas)
    init_db()
    
    logger.info("✅ IA Service Python iniciado correctamente")
    
    yield
    
    # Shutdown
    logger.info("🛑 Deteniendo IA Service Python...")


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
            "detect": "/detect/{activo_id}"
        }
    }


if __name__ == "__main__":
    import uvicorn
    
    uvicorn.run(
        "app.main:app",
        host=config.host,
        port=config.port,
        reload=config.debug,
        log_level=config.log_level.lower()
    )
