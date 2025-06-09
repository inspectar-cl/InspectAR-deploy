from fastapi import FastAPI, HTTPException, Response
from fastapi.responses import JSONResponse, StreamingResponse
from pydantic import BaseModel
from typing import List
from pymongo import MongoClient
from influxdb_client import InfluxDBClient, Point, WriteOptions
from influxdb_client.client.query_api import QueryApi
import os
from datetime import datetime
import pandas as pd
import numpy as np
from io import StringIO, BytesIO
import csv

# Configuración desde variables de entorno (recomendado en producción)
MONGO_URI = os.getenv("MONGO_URI", "mongudb://localhost:27017")
INFLUX_URL = os.getenv("INFLUX_URL", "http://localhost:8086")
INFLUX_TOKEN = os.getenv("INFLUX_TOKEN", "my-super-token")
INFLUX_ORG = os.getenv("INFLUX_ORG", "my-org")
INFLUX_BUCKET = os.getenv("INFLUX_BUCKET", "sensores")

# Conexiones
mongo_client = MongoClient(MONGO_URI)
mongo_db = mongo_client["iot"]
activos_col = mongo_db["activos"]

influx_client = InfluxDBClient(url=INFLUX_URL, token=INFLUX_TOKEN, org=INFLUX_ORG)
query_api: QueryApi = influx_client.query_api()
write_api = influx_client.write_api(write_options=WriteOptions(batch_size=1))

app = FastAPI()

# Esquemas de entrada para la API
class Sensor(BaseModel):
    sensor_id: str
    tipo: str
    unidad: str

class Activo(BaseModel):
    activo_id: str
    nombre: str
    ubicacion: str
    sensores: List[Sensor]

class Lectura(BaseModel):
    sensor_id: str
    valor: float
    timestamp: datetime = datetime.utcnow()

# ========================= RUTAS ==========================

@app.post("/activo/")
def crear_activo(activo: Activo):
    if activos_col.find_one({"activo_id": activo.activo_id}):
        raise HTTPException(status_code=400, detail="Activo ya existe")
    activos_col.insert_one(activo.dict())
    return {"mensaje": "Activo registrado exitosamente"}

@app.post("/lectura/")
def registrar_lectura(data: Lectura):
    point = Point("mediciones")\
        .tag("sensor_id", data.sensor_id)\
        .field("valor", data.valor)\
        .time(data.timestamp)
    write_api.write(bucket=INFLUX_BUCKET, org=INFLUX_ORG, record=point)
    return {"mensaje": "Lectura registrada"}

@app.get("/activo/{activo_id}")
def get_activo_con_datos(activo_id: str):
    activo = activos_col.find_one({"activo_id": activo_id})
    if not activo:
        raise HTTPException(status_code=404, detail="Activo no encontrado")

    sensores = activo.get("sensores", [])
    resultado = {
        "activo": activo.get("nombre"),
        "ubicacion": activo.get("ubicacion"),
        "sensores": []
    }

    for sensor in sensores:
        sensor_id = sensor.get("sensor_id")
        tipo = sensor.get("tipo")
        unidad = sensor.get("unidad")

        query = f'''
        from(bucket: "{INFLUX_BUCKET}")
          |> range(start: -1h)
          |> filter(fn: (r) => r["sensor_id"] == "{sensor_id}")
          |> last()
        '''
        tables = query_api.query(query=query)
        valor = None
        timestamp = None
        for table in tables:
            for record in table.records:
                valor = record.get_value()
                timestamp = record.get_time()

        resultado["sensores"].append({
            "sensor_id": sensor_id,
            "tipo": tipo,
            "valor_actual": valor,
            "unidad": unidad,
            "timestamp": timestamp
        })

    return resultado

@app.get("/activo/{activo_id}/datos")
async def get_datos_historicos(activo_id: str, format: str = "csv"):
    """
    Obtiene los datos históricos de un activo.
    Formatos soportados:
    - parquet: Archivo binario optimizado para datos tabulares
    - csv: Archivo de texto con valores separados por comas
    - excel: Archivo Excel
    - hdf: Archivo HDF5 (hierarchical data format)
    """
    activo = activos_col.find_one({"activo_id": activo_id})
    if not activo:
        raise HTTPException(status_code=404, detail="Activo no encontrado")

    # Preparar listas para el DataFrame
    datos = []
    
    for sensor in activo.get("sensores", []):
        sensor_id = sensor.get("sensor_id")
        tipo = sensor.get("tipo")
        unidad = sensor.get("unidad")

        query = f'''
        from(bucket: "{INFLUX_BUCKET}")
          |> range(start: 0)
          |> filter(fn: (r) => r["sensor_id"] == "{sensor_id}")
          |> sort(columns: ["_time"])
        '''
        
        try:
            tables = query_api.query(query=query)
            for table in tables:
                for record in table.records:
                    datos.append({
                        "timestamp": record.get_time(),
                        "sensor_id": sensor_id,
                        "tipo_sensor": tipo,
                        "valor": record.get_value(),
                        "unidad": unidad,
                        "activo": activo.get("nombre"),
                        "ubicacion": activo.get("ubicacion")
                    })
        except Exception as e:
            raise HTTPException(status_code=500, detail=f"Error al consultar InfluxDB: {str(e)}")    # Si no hay datos, retornar error
    if not datos:
        raise HTTPException(status_code=404, detail="No se encontraron datos para este activo")

    # Crear DataFrame y procesar datos
    df = pd.DataFrame(datos)
    df = df.sort_values("timestamp")

    # Convertir tipos de datos
    df["valor"] = pd.to_numeric(df["valor"])
    df["timestamp"] = pd.to_datetime(df["timestamp"])

    # Crear un buffer en memoria para guardar el archivo
    buffer = StringIO() if format.lower() == "csv" else BytesIO()

    # Generar el archivo según el formato solicitado
    filename = f"{activo_id}_datos"
    content_type = ""    
    try:
        if format.lower() == "csv":
            df.to_csv(buffer, index=False)
            filename += ".csv"
            content_type = "text/csv"
            
        elif format.lower() == "parquet":
            df.to_parquet(buffer, engine="pyarrow", index=False)
            filename += ".parquet"
            content_type = "application/octet-stream"
        
        elif format.lower() == "excel":
            df.to_excel(buffer, index=False, engine="openpyxl")
            filename += ".xlsx"
            content_type = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
        
        elif format.lower() == "hdf":
            df.to_hdf(buffer, key="data", mode="w", format="table")
            filename += ".h5"
            content_type = "application/x-hdf5"
        
        else:
            raise HTTPException(status_code=400, detail=f"Formato {format} no soportado")

        # Preparar la respuesta
        headers = {
            "Content-Disposition": f'attachment; filename="{filename}"'
        }

        # Si es CSV, necesitamos obtener el valor como string
        if format.lower() == "csv":
            content = buffer.getvalue()
        else:
            content = buffer.getvalue()
            buffer.seek(0)

        return StreamingResponse(
            iter([content]),
            media_type=content_type,
            headers=headers
        )

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error al generar el archivo: {str(e)}")
    
    finally:
        buffer.close()
