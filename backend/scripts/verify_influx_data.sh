#!/bin/bash

# Script para verificar datos cargados en InfluxDB

echo "📊 Verificando datos en InfluxDB para activo_id 2..."
echo "=================================================="
echo ""

echo "🔍 Conteo de registros por sensor:"
docker exec influxdb influx query '
from(bucket: "sensores")
  |> range(start: 2024-01-01T00:00:00Z)
  |> filter(fn: (r) => r["activo_id"] == "2")
  |> group(columns: ["sensor_id"])
  |> count()
  |> yield(name: "count")
' --org my-org --token my-super-token

echo ""
echo "📅 Rango de fechas de los datos:"
docker exec influxdb influx query '
from(bucket: "sensores")
  |> range(start: 2024-01-01T00:00:00Z)
  |> filter(fn: (r) => r["activo_id"] == "2")
  |> group(columns: ["sensor_id"])
  |> keep(columns: ["_time", "sensor_id"])
  |> min(column: "_time")
' --org my-org --token my-super-token | head -20

echo ""
echo "✅ Verificación completada"
