import requests
from config import settings

def insert_row_to_bigquery(data: dict):
    """
    Envía los datos a otro microservicio vía HTTP POST.
    :param data: Diccionario con los datos a enviar.
    """
    url = settings.TARGET_MICROSERVICE_URL  # Debes agregar esta variable en settings.py
    try:
        response = requests.post(url, json=data, timeout=5)
        response.raise_for_status()
        print("Datos enviados correctamente al microservicio.")
    except Exception as e:
        print(f"Error al enviar datos al microservicio: {e}")