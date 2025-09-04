import requests

API_URL = 'http://localhost:3500/api/notificacion/notification'

def enviar_correo(activo):
    # JSON que espera la API
    payload = {
        "activo_id": 1
    }

    headers = {
        "Content-Type": "application/json"
    }
    #Enviar correo

    try:
        response = requests.post(API_URL, json=payload, headers=headers)

        if response.status_code == 200 or response.status_code == 201:
            data = response.json()
            print("✅ Notificación creada exitosamente")
            print(data)
        else:
            print(f"❌ Error {response.status_code}: {response.text}")

    except requests.exceptions.RequestException as e:
        print("⚠️ Error en la conexión:", e)