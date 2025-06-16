import threading
import time
import json
import random
from datetime import datetime
import paho.mqtt.client as mqtt

# Configuración del broker MQTT
BROKER_HOST = "localhost"  # Cambia a IP real si es necesario
BROKER_PORT = 1883

# Función general para cada sensor
def sensor_loop(sensor_id, topic, rango_min, rango_max, intervalo):
    client = mqtt.Client()
    client.connect(BROKER_HOST, BROKER_PORT)
    
    while True:
        valor = round(random.uniform(rango_min, rango_max), 2)
        timestamp = datetime.now().isoformat()
        payload = {
            "sensor_id": sensor_id,
            "valor": str(valor),
            "timestamp": timestamp
        }
        client.publish(topic, json.dumps(payload), retain=True)
        print(f"[{sensor_id}] Publicado en {topic}: {payload}")
        time.sleep(intervalo)

# Definir los sensores como hebras
def iniciar_simulador():
    sensores = [
        {
            "sensor_id": "temp-1",
            "topic": "sensors/temperatura",
            "min": 20.0,
            "max": 30.0,
            "intervalo": 3
        },
        {
            "sensor_id": "caudal-1",
            "topic": "sensors/caudal",
            "min": 50.0,
            "max": 120.0,
            "intervalo": 5
        },
        {
            "sensor_id": "presion-1",
            "topic": "sensors/presion",
            "min": 1.0,
            "max": 5.0,
            "intervalo": 4
        }
    ]

    for sensor in sensores:
        thread = threading.Thread(
            target=sensor_loop,
            args=(sensor["sensor_id"], sensor["topic"], sensor["min"], sensor["max"], sensor["intervalo"]),
            daemon=True
        )
        thread.start()

    # Mantener el programa vivo
    while True:
        time.sleep(1)

if __name__ == "__main__":
    iniciar_simulador()
