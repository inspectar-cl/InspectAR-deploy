import os
from dotenv import load_dotenv

# Cargar variables desde .env si existe
load_dotenv()

TARGET_MICROSERVICE_URL = os.getenv("TARGET_MICROSERVICE_URL", "http://localhost:8000/endpoint")

# Configuración de MQTT
MQTT_BROKER = os.getenv("MQTT_BROKER", "localhost")
MQTT_PORT = int(os.getenv("MQTT_PORT", 1883))
MQTT_TOPIC = os.getenv("MQTT_TOPIC", "sensor/datos")

# Otras configuraciones globales si las necesitas