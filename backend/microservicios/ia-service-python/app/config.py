"""
Configuración del servicio IA
Gestión de variables de entorno y configuración
"""

from pydantic_settings import BaseSettings
from typing import Optional
import yaml
import os


class Settings(BaseSettings):
    """Configuración del servicio de IA"""
    
    # Configuración del servidor
    app_name: str = "IA Service Python"
    host: str = "0.0.0.0"
    port: int = 8095
    debug: bool = False
    log_level: str = "INFO"
    
    # Base de datos PostgreSQL
    database_host: str = "ia-db"
    database_port: int = 5432
    database_name: str = "ia_db"
    database_user: str = "ia_user"
    database_password: str = "ia_pass"
    
    # ML Engine
    ml_engine_url: str = "http://ml_engine_python:9999"
    ml_engine_timeout: int = 300
    
    # IOT Service (ParserService)
    iot_service_url: str = "http://iot-service:8090"
    iot_service_timeout: int = 60
    
    # Scheduler configuration
    scheduler_enabled: bool = True
    scheduler_interval_minutes: int = 2
    
    # Configuración de ventana de datos
    window_size: int = 1000
    window_page: int = 1
    
    @property
    def database_url(self) -> str:
        """Construye la URL de conexión a PostgreSQL"""
        return (
            f"postgresql://{self.database_user}:{self.database_password}"
            f"@{self.database_host}:{self.database_port}/{self.database_name}"
        )
    
    class Config:
        env_file = ".env"
        case_sensitive = False
        env_prefix = ""


def load_config(config_file: Optional[str] = None) -> Settings:
    """
    Carga la configuración desde archivo YAML si existe,
    sino usa variables de entorno
    """
    settings = Settings()
    
    if config_file and os.path.exists(config_file):
        try:
            with open(config_file, 'r') as f:
                yaml_config = yaml.safe_load(f)
                
            # Actualizar settings con valores del YAML
            for key, value in yaml_config.items():
                if hasattr(settings, key):
                    setattr(settings, key, value)
        except Exception as e:
            print(f"⚠️  Error cargando config YAML: {e}")
    
    return settings


# Instancia global de configuración
config = load_config(os.getenv("CONFIG_FILE"))
