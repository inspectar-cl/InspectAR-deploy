import requests
import random
import time
from datetime import datetime

API_URL = "http://localhost:8000/lectura/"
SENSOR_ID = "temp1"  # Este sensor debe estar registrado en el activo previamente

def simular_sensor():
    while True:
        valor = round(random.uniform(60.0, 80.0), 2)  # temperatura simulada
        timestamp = datetime.utcnow().isoformat()

        data = {
            "sensor_id": SENSOR_ID,
            "valor": valor,
            "timestamp": timestamp
        }

        try:
            response = requests.post(API_URL, json=data)
            if response.status_code == 200:
                print(f"[{timestamp}] Enviado: {valor} °C")
            else:
                print(f"Error {response.status_code}: {response.text}")
        except Exception as e:
            print(f"Error de conexión: {e}")

        time.sleep(5)  # Espera 5 segundos antes de enviar el siguiente valor

if __name__ == "__main__":
    simular_sensor()
