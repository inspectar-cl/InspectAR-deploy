import time
import paho.mqtt.client as mqtt
from database.inserter import insert_row_to_bigquery 

BROKER = 'localhost'  # Cambia por la IP o dominio de tu broker EMQX
PORT = 1883
TOPIC = 'sensor/datos'  # Cambia por el tópico que uses

def on_connect(client, userdata, flags, rc):
    if rc == 0:
        print("Conectado a EMQX correctamente.")
        client.subscribe(TOPIC)
    else:
        print(f"Error de conexión: {rc}")

def on_message(client, userdata, msg):
    print(f"Mensaje recibido en {msg.topic}: {msg.payload.decode()}")
    try:
        # Suponiendo que el payload es JSON
        import json
        data = json.loads(msg.payload.decode())
        insert_row_to_bigquery(data)
    except Exception as e:
        print(f"Error procesando el mensaje: {e}")

time.sleep(5)  # Espera 5 segundos antes de iniciar el cliente

client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect(BROKER, PORT, 60)
print("Esperando mensajes de EMQX...")

client.loop_forever()