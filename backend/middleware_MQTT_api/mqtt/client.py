import paho.mqtt.client as mqtt

# Configuración del broker EMQX
BROKER = 'localhost'  # Cambia esto por la IP o dominio de tu broker EMQX
PORT = 1883
TOPIC = 'sensor/datos'  # Cambia esto por el tópico que uses

def on_connect(client, userdata, flags, rc):
    if rc == 0:
        print("Conectado al broker MQTT!")
        client.subscribe(TOPIC)
    else:
        print(f"Fallo la conexión, código de error: {rc}")

def on_message(client, userdata, msg):
    print(f"Mensaje recibido en el tópico {msg.topic}: {msg.payload.decode()}")

client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect(BROKER, PORT, 60)
print("Esperando mensajes...")

client.loop_forever()