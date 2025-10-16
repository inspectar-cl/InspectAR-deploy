import requests
from datetime import datetime, timezone
import random
import time

# Endpoint correcto
url = 'http://localhost:8090/lectura'

sensor_ids = ['sensortemp']

def generar_valor(sensor_id):
    if sensor_id == 'sensortemp':
        return round(random.uniform(60, 80), 2)
    # elif sensor_id == 'caud1':
    #     return round(random.uniform(50, 90), 2)
    # elif sensor_id == 'pres1':
    #     return round(random.uniform(60, 80), 2)
    else:
        return 0

while True:
    # Timestamp en formato ISO con Z al final (UTC)
    timestamp = datetime.now(timezone.utc).isoformat().replace("+00:00", "Z")

    for sensor_id in sensor_ids:
        valor = generar_valor(sensor_id)

        payload = {
            "sensor_id": sensor_id,
            "valor": valor,
            "timestamp": timestamp
        }

        try:
            response = requests.post(url, json=payload)
            print(f"Enviado a {sensor_id}: {valor} @ {timestamp} | status: {response.status_code}")
        except Exception as e:
            print(f"Error al enviar a {sensor_id}: {e}")

    time.sleep(2)  # Espera 5 segundos antes del siguiente grupo
