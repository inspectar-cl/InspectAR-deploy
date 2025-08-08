import requests

activo_id = "BombaDeAgua1"
# Endpoint correcto
url = f'http://localhost:8090/activo/{activo_id}/estado'

def put_estado(estado):
    payload = {
        "estado": estado
    }

    try:
        response = requests.put(url, json=payload)
        if response.status_code == 200:
            print(f"🌐 Actualizando estado a: {estado} | status: {response.status_code}")
        else:
            print(f"❌ API {response.status_code}: {response.text}")
    except Exception as e:
        print(f"🚫 Error de conexión a API: {e}")