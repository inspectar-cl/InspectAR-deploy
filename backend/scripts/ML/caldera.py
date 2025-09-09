import time
import numpy as np
import requests
from datetime import datetime

API_URL = "http://localhost:8090/lectura/"
SENSOR_IDS = {
    "temperatura": "temp1",
    "presion": "pres1",
    "caudal": "caud1"
}

def enviar_a_api(sensor_id, valor):
    timestamp = datetime.utcnow().isoformat() + "Z"
    data = {
        "sensor_id": sensor_id,
        "valor": round(valor, 2),
        "timestamp": timestamp
    }
    try:
        response = requests.post(API_URL, json=data)
        if response.status_code == 200:
            print(f"🌐 [{timestamp}] API -> {sensor_id}: {valor:.2f}")
        else:
            print(f"❌ API {response.status_code}: {response.text}")
    except Exception as e:
        print(f"🚫 Error de conexión a API: {e}")

def generar_variable_sensor(base, amplitud, ruido, t, oscilacion=60):
    variacion = amplitud * np.sin(2 * np.pi * t / oscilacion)
    ruido_gauss = np.random.normal(0, ruido)
    return base + variacion + ruido_gauss

def run_simulador(data_queue, control_queue):
    print("🔧 Simulador activo (Ctrl+C para salir)")
    modo_anomalo = 0
    t = 0

    try:
        while True:
            # Leer comandos desde control_queue
            if not control_queue.empty():
                comando = control_queue.get()
                if comando in ["1", "2", "3"]:
                    modo_anomalo = int(comando)
                else:
                    modo_anomalo = 0

            timestamp = datetime.now().strftime('%Y-%m-%d %H:%M:%S')

            # Lógica de datos según modo
            if modo_anomalo == 1:
                temp = np.random.uniform(75, 80)
                presion_bar = np.random.uniform(3, 3.2)
                caudal = np.random.uniform(105, 110)
                etiqueta = "Anormal"

            elif modo_anomalo == 2:
                temp = np.random.uniform(85, 90)
                presion_bar = np.random.uniform(4, 5)
                caudal = np.random.uniform(115, 120)
                etiqueta = None

            elif modo_anomalo == 3:
                temp = np.random.uniform(170, 180)
                presion_bar = np.random.uniform(13, 14)
                caudal = np.random.uniform(200, 220)
                etiqueta = "Anormal"

            else:
                temp = generar_variable_sensor(70, 5, 0.5, t, 300)
                presion_bar = generar_variable_sensor(2.5, 0.3, 0.05, t, 150)
                caudal = generar_variable_sensor(100, 20, 5, t, 600)
                etiqueta = "Normal"

            presion_psi = presion_bar * 14.5038

            # 🌐 Enviar a API (sin etiqueta)
            enviar_a_api(SENSOR_IDS["temperatura"], temp)
            enviar_a_api(SENSOR_IDS["presion"], presion_psi)
            enviar_a_api(SENSOR_IDS["caudal"], caudal)

            # 📦 Enviar a modelo de ML (con etiqueta)
            data_queue.put({
                "temperatura": temp,
                "presion": presion_psi,
                "caudal": caudal,
                "etiqueta": etiqueta
            })

            print(f"📤 [{timestamp}] Temp: {temp:.2f} °C | Presión: {presion_psi:.2f} PSI | Caudal: {caudal:.2f} L/min | {'Sin etiqueta' if etiqueta is None else etiqueta}")
            time.sleep(5)
            t += 5

    except KeyboardInterrupt:
        print("\n🛑 Simulador detenido.")
