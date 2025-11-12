import threading
import time
import json
import random
from datetime import datetime, timedelta
import paho.mqtt.client as mqtt

# Configuración del broker MQTT
BROKER_HOST = "localhost"  # Cambia a IP real si es necesario
BROKER_PORT = 1883

# Función general para cada sensor
def sensor_loop(sensor_id, topic, rango_min, rango_max, intervalo):
    client = mqtt.Client()
    client.connect(BROKER_HOST, BROKER_PORT)
    
    contador = 0
    while True:
        contador += 1
        
        # Simular desconexiones ocasionales para sensores críticos
        if sensor_id.startswith("SENSOR_AC1005") or sensor_id.startswith("SENSOR_AC1008"):
            # 30% de probabilidad de desconexión para sensores críticos
            if random.random() < 0.3:
                valor = 0.0
                print(f"[{sensor_id}] DESCONECTADO - Valor: {valor}")
            else:
                valor = round(random.uniform(rango_min, rango_max), 2)
        else:
            # Generar valor normal con ocasionales picos anómalos
            if random.random() < 0.05:  # 5% probabilidad de anomalía
                # Generar valor anómalo (fuera del rango normal)
                if random.random() < 0.5:
                    valor = round(rango_max * random.uniform(1.2, 1.8), 2)  # Pico alto
                else:
                    valor = round(rango_min * random.uniform(0.2, 0.8), 2)  # Pico bajo
                print(f"[{sensor_id}] ANOMALÍA DETECTADA - Valor: {valor}")
            else:
                valor = round(random.uniform(rango_min, rango_max), 2)
        
        # Timestamp con 3 horas menos
        timestamp = (datetime.now() - timedelta(hours=3)).isoformat()
        payload = {
            "sensor_id": sensor_id,
            "valor": str(valor),
            "timestamp": timestamp
        }
        client.publish(topic, json.dumps(payload), retain=True)
        
        # Log cada 10 envíos para reducir spam
        if contador % 10 == 0 or "ANOMALÍA" in str(valor) or "DESCONECTADO" in str(valor):
            print(f"[{sensor_id}] Publicado en {topic}: {payload}")
        
        time.sleep(intervalo)

# Definir los sensores como hebras
def iniciar_simulador():
    sensores = [
        # Sensores para Activo 2 (Bomba Centrífuga A) - 10 sensores
        {
            "sensor_id": "A_ACR_Mot.PV",
            "topic": "sensors/motor_pv",
            "min": 40.0,
            "max": 50.0,
            "intervalo": 1
        },
        {
            "sensor_id": "A_ACR_Mot.SV",
            "topic": "sensors/motor_sv",
            "min": 45.0,
            "max": 55.0,
            "intervalo": 1
        },
        {
            "sensor_id": "A_ACR_Mot.TV",
            "topic": "sensors/motor_tv",
            "min": 45.0,
            "max": 52.0,
            "intervalo": 1
        },
        {
            "sensor_id": "A_ACR_Pmp.PV",
            "topic": "sensors/bomba_pv",
            "min": 70.0,
            "max": 75.0,
            "intervalo": 1
        },
        {
            "sensor_id": "A_ACR_Pmp.SV", 
            "topic": "sensors/bomba_sv",
            "min": 72.0,
            "max": 78.0,
            "intervalo": 1
        },
        {
            "sensor_id": "A_ACR_Pmp.TV",
            "topic": "sensors/bomba_tv",
            "min": 72.0,
            "max": 76.0,
            "intervalo": 1
        },
        {
            "sensor_id": "A_Pres.PV",
            "topic": "sensors/presion_pv",
            "min": 3.0,
            "max": 4.0,
            "intervalo": 2
        },
        {
            "sensor_id": "A_Temp.PV",
            "topic": "sensors/temperatura_pv",
            "min": 62.0,
            "max": 68.0,
            "intervalo": 1
        },
        {
            "sensor_id": "Barometer",
            "topic": "sensors/barometer",
            "min": 1010.0,
            "max": 1020.0,
            "intervalo": 2
        },
        {
            "sensor_id": "Temperature",
            "topic": "sensors/temperature_ambient",
            "min": 20.0,
            "max": 25.0,
            "intervalo": 2
        },
        # Sensores para Activo 3
        {
            "sensor_id": "SENSOR_AC1003_01",
            "topic": "sensors/caudal_ac3",
            "min": 270.0,
            "max": 290.0,
            "intervalo": 1
        },
        {
            "sensor_id": "SENSOR_AC1003_02",
            "topic": "sensors/presion_ac3",
            "min": 7.5,
            "max": 8.2,
            "intervalo": 2
        },
        # Sensores para Activo 4
        {
            "sensor_id": "SENSOR_AC1004_01",
            "topic": "sensors/vibracion_ac4",
            "min": 8.5,
            "max": 9.8,
            "intervalo": 1
        },
        {
            "sensor_id": "SENSOR_AC1004_02", 
            "topic": "sensors/voltaje_ac4",
            "min": 395.0,
            "max": 405.0,
            "intervalo": 1
        },
        # Sensores para Activo 5 (simulando desconexiones ocasionales)
        {
            "sensor_id": "SENSOR_AC1005_01",
            "topic": "sensors/voltaje_ac5",
            "min": 0.0,
            "max": 220.0,
            "intervalo": 3
        },
        {
            "sensor_id": "SENSOR_AC1005_02",
            "topic": "sensors/temperatura_ac5", 
            "min": 0.0,
            "max": 45.0,
            "intervalo": 3
        },
        # Sensores para Activo 6
        {
            "sensor_id": "SENSOR_AC1006_01",
            "topic": "sensors/voltaje_ac6",
            "min": 218.0,
            "max": 224.0,
            "intervalo": 1
        },
        {
            "sensor_id": "SENSOR_AC1006_02",
            "topic": "sensors/temperatura_ac6",
            "min": 40.0,
            "max": 45.0,
            "intervalo": 2
        },
        # Sensores para Activo 7
        {
            "sensor_id": "SENSOR_AC1007_01",
            "topic": "sensors/vibracion_ac7",
            "min": 7.8,
            "max": 8.6,
            "intervalo": 1
        },
        {
            "sensor_id": "SENSOR_AC1007_02",
            "topic": "sensors/voltaje_ac7",
            "min": 390.0,
            "max": 400.0,
            "intervalo": 1
        },
        # Sensores para Activo 8 (simulando fallas críticas)
        {
            "sensor_id": "SENSOR_AC1008_01", 
            "topic": "sensors/presion_ac8",
            "min": 0.0,
            "max": 2.0,
            "intervalo": 3
        },
        {
            "sensor_id": "SENSOR_AC1008_02",
            "topic": "sensors/caudal_ac8",
            "min": 0.0,
            "max": 50.0,
            "intervalo": 3
        }
    ]

    print("\n" + "="*60)
    print("🚀 INICIANDO SIMULADOR MQTT - TODOS LOS SENSORES")
    print("="*60)
    print(f"📡 Broker: {BROKER_HOST}:{BROKER_PORT}")
    print(f"🔧 Total de sensores: {len(sensores)}")
    print("📊 Distribución por activo:")
    
    # Contar sensores por activo
    activos = {}
    for sensor in sensores:
        if sensor["sensor_id"].startswith("A_"):
            activo = "Activo 2 (Bomba Centrífuga A)"
        elif "1003" in sensor["sensor_id"]:
            activo = "Activo 3"
        elif "1004" in sensor["sensor_id"]:
            activo = "Activo 4"
        elif "1005" in sensor["sensor_id"]:
            activo = "Activo 5 (Crítico)"
        elif "1006" in sensor["sensor_id"]:
            activo = "Activo 6"
        elif "1007" in sensor["sensor_id"]:
            activo = "Activo 7"
        elif "1008" in sensor["sensor_id"]:
            activo = "Activo 8 (Crítico)"
        else:
            activo = "Otros"
        
        activos[activo] = activos.get(activo, 0) + 1
    
    for activo, count in activos.items():
        print(f"   • {activo}: {count} sensores")
    
    print("\n🔄 Características especiales:")
    print("   • Activos 5 y 8: Simulan desconexiones (30% prob.)")
    print("   • Todos los sensores: 5% probabilidad de anomalías")
    print("   • Logs reducidos: Cada 10 envíos + anomalías")
    print("="*60)
    
    for sensor in sensores:
        thread = threading.Thread(
            target=sensor_loop,
            args=(sensor["sensor_id"], sensor["topic"], sensor["min"], sensor["max"], sensor["intervalo"]),
            daemon=True
        )
        thread.start()
        print(f"✅ Iniciado: {sensor['sensor_id']} -> {sensor['topic']} (cada {sensor['intervalo']}s)")

    print(f"\n🎯 {len(sensores)} sensores activos. Presiona Ctrl+C para detener.")
    
    # Mantener el programa vivo
    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print("\n\n🛑 Simulador detenido por el usuario")
        print("="*60)

if __name__ == "__main__":
    iniciar_simulador()
