from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from app.model_htm import detectar_htm_multivar
import pandas as pd

app = FastAPI(
    title="InspectAR ML Engine",
    description="Microservicio Python para predicción de anomalías en bombas centrífugas",
    version="1.0.0"
)

# ============================================================
# Entrada esperada
# ============================================================

class PumpWindow(BaseModel):
    pump: str
    window: list[dict]  # bloque W de datos del sensor (dicts con timestamp + features)


@app.get("/healthz")
def health_check():
    return {"status": "ok"}

@app.post("/predict_anomaly")
async def predict_anomaly(request: Request):
    """
    Recibe un bloque W de datos (JSON) y devuelve score, likelihood y severidad.
    """
    payload = await request.json()
    pump = payload.get("pump")
    records = payload.get("records", [])
    
    if not pump or not records:
        return JSONResponse({"error": "Datos insuficientes"}, status_code=400)

    df = pd.DataFrame(records)
    if 'timestamp' in df.columns:
        df['timestamp'] = pd.to_datetime(df['timestamp'], errors='coerce')

    # --- Ejecutar modelo HTM-like ---
    result = detectar_htm_multivar(df, pump)

    if result is None:
        return JSONResponse({"error": "No se pudo generar predicción"}, status_code=500)

    return JSONResponse(result)
