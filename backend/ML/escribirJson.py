import json
from datetime import datetime
import os

ALERT_JSON_PATH = os.path.abspath(
    os.path.join(os.path.dirname(__file__), "..", "..", "frontend", "src", "mocks", "alerts.json")
)

def escribir_alerta(id_alerta: str, nombre_alerta: str):
    try:
        with open(ALERT_JSON_PATH, "r") as f:
            data = json.load(f)
    except (FileNotFoundError, json.JSONDecodeError):
        data = {"products": []}

    # Obtener la hora actual
    hora_actual = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    # Buscar si ya existe una alerta con ese ID
    existente = next((item for item in data["products"] if item["id"] == id_alerta), None)

    if existente:
        existente["updatedAt"] = hora_actual  # Actualizar hora
    else:
        # Agregar nueva alerta
        data["products"].append({
            "id": id_alerta,
            "name": nombre_alerta,
            "updatedAt": hora_actual
        })

    # Guardar de vuelta
    with open(ALERT_JSON_PATH, "w") as f:
        json.dump(data, f, indent=2)