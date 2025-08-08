import json
import time
import paho.mqtt.client as mqtt

# Callback al conectar al broker
def on_connect(client, userdata, flags, rc):
    if rc == 0:
        print("✅ Conectado al broker MQTT")
        client.subscribe("sensors/#")  # Suscribirse a todos los sensores
    else:
        print("❌ Error de conexión:", rc)

# Callback al recibir un mensaje
def on_message(client, userdata, msg):
    try:
        payload = msg.payload.decode('utf-8')
        data = json.loads(payload)

        sensor_id = data.get("sensor_id")
        valor = data.get("valor")
        timestamp = data.get("timestamp")

        print(f"📨 Mensaje recibido en topic {msg.topic}")
        print(f"  Sensor ID: {sensor_id}")
        print(f"  Valor: {valor}")
        print(f"  Timestamp: {timestamp}")
        print("-" * 40)

    except Exception as e:
        print("⚠️ Error procesando mensaje:", e)

# Crear cliente MQTT
client = mqtt.Client(client_id="python-subscriber", protocol=mqtt.MQTTv311)
client.on_connect = on_connect
client.on_message = on_message

# Conectar al broker EMQX
client.connect("localhost", 1883, 60)  # Cambia host si usas Docker

# Bucle principal
client.loop_start()

try:
    while True:
        time.sleep(1)
except KeyboardInterrupt:
    print("🛑 Cerrando conexión...")
    client.loop_stop()
    client.disconnect()
