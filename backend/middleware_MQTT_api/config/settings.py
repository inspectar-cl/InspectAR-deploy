import os
from dotenv import load_dotenv

# Cargar variables desde .env si existe
load_dotenv()

# Credenciales de Google Cloud
GOOGLE_APPLICATION_CREDENTIALS = os.getenv("GOOGLE_APPLICATION_CREDENTIALS", "ruta/por/defecto/credenciales.json")
GCP_PROJECT_ID = os.getenv("GCP_PROJECT_ID", "tu-proyecto-id")

# Configuración de MQTT
MQTT_BROKER = os.getenv("MQTT_BROKER", "localhost")
MQTT_PORT = int(os.getenv("MQTT_PORT", 1883))
MQTT_TOPIC = os.getenv("MQTT_TOPIC", "sensor/datos")

# Configuración de BigQuery
BQ_DATASET = os.getenv("BQ_DATASET", "nombre_dataset")
BQ_TABLE = os.getenv("BQ_TABLE", "nombre_tabla")

# Otras configuraciones globales si las necesitas