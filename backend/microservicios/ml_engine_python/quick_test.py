"""
Prueba rápida del servicio ML - Genera datos y obtiene predicción.
"""

import requests
import json
from datetime import datetime, timedelta
import numpy as np

BASE_URL = "http://localhost:8086"

# Generar datos rápidos
print("Generando datos sintéticos...")
start = datetime.now() - timedelta(seconds=200)
records = []

for i in range(200):
    ts = start + timedelta(seconds=i)
    # Inyectar anomalía en últimos registros
    anomaly = 10 if i > 180 else 0
    
    records.append({
        "timestamp": ts.isoformat() + "Z",
        "A_ACR_Mot.PV": round(0.012 + np.random.normal(0, 0.005) + anomaly * 0.01, 4),
        "A_ACR_Mot.TV": round(41.0 + np.random.normal(0, 2) + anomaly, 2),
        "A_Pres.PV": round(1.85 + np.random.normal(0, 0.1), 2),
        "A_Temp.PV": round(75.0 + np.random.normal(0, 3) + anomaly, 1),
        "Temperature": round(22.0 + np.random.normal(0, 0.5), 1)
    })

print(f"✅ {len(records)} registros generados")
print(f"⚠️  Anomalía inyectada en últimos 20 registros\n")

# Enviar request
print("Enviando request a ML Engine...")
response = requests.post(
    f"{BASE_URL}/predict_anomaly",
    json={"pump": "test_pump", "records": records}
)

if response.status_code == 200:
    result = response.json()
    print(f"\n✅ Predicción exitosa!")
    print(f"Anomalías detectadas: {result['n_anomalies']} de {result['n_rows']}")
    print(f"\nÚltimo resultado:")
    print(json.dumps(result['last'], indent=2))
else:
    print(f"\n❌ Error: {response.status_code}")
    print(response.text)