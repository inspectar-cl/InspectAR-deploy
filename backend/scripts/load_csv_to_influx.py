#!/usr/bin/env python3
"""
Script para cargar datos CSV en InfluxDB
"""
import csv
import sys
from datetime import datetime
import requests
from pathlib import Path

# Configuración de InfluxDB
INFLUX_URL = "http://influxdb:8086"
INFLUX_ORG = "my-org"
INFLUX_BUCKET = "sensores"
INFLUX_TOKEN = "my-super-token"

def parse_timestamp(ts_str):
    """Convierte timestamp ISO8601 a nanosegundos"""
    try:
        dt = datetime.fromisoformat(ts_str.replace('Z', '+00:00'))
        return int(dt.timestamp() * 1_000_000_000)
    except Exception as e:
        print(f"Error parsing timestamp {ts_str}: {e}")
        return None

def load_csv_to_influx(csv_file, activo_id, batch_size=500):
    """Carga un archivo CSV en InfluxDB"""
    print(f"   Procesando archivo: {Path(csv_file).name}")
    
    batch = []
    line_count = 0
    batch_count = 0
    error_count = 0
    
    with open(csv_file, 'r') as f:
        reader = csv.DictReader(f)
        
        for row in reader:
            timestamp_ns = parse_timestamp(row['Timestamp'])
            if not timestamp_ns:
                continue
            
            # Crear line protocol para cada sensor
            sensors = [
                ('A_ACR_Mot.PV', row['A_ACR_Mot.PV']),
                ('A_ACR_Mot.SV', row['A_ACR_Mot.SV']),
                ('A_ACR_Mot.TV', row['A_ACR_Mot.TV']),
                ('A_ACR_Pmp.PV', row['A_ACR_Pmp.PV']),
                ('A_ACR_Pmp.SV', row['A_ACR_Pmp.SV']),
                ('A_ACR_Pmp.TV', row['A_ACR_Pmp.TV']),
                ('A_Pres.PV', row['A_Pres.PV']),
                ('A_Temp.PV', row['A_Temp.PV']),
                ('Barometer', row['Barometer']),
                ('Temperature', row['Temperature'])
            ]
            
            for sensor_id, value in sensors:
                try:
                    # Formato: measurement,tag1=value1,tag2=value2 field1=value1 timestamp
                    line = f"sensor_reading,sensor_id={sensor_id},activo_id={activo_id} valor={float(value)} {timestamp_ns}"
                    batch.append(line)
                except ValueError:
                    continue
            
            line_count += 1
            
            # Enviar batch cuando alcance el tamaño
            if len(batch) >= batch_size * 10:  # 10 sensores por línea
                success = send_batch(batch)
                if success:
                    batch_count += 1
                    print(f"      ✓ Batch {batch_count} enviado ({line_count} líneas, {len(batch)} puntos)")
                else:
                    error_count += 1
                    print(f"      ✗ Error en batch {batch_count + 1}")
                batch = []
    
    # Enviar el último batch
    if batch:
        success = send_batch(batch)
        if success:
            batch_count += 1
            print(f"      ✓ Último batch enviado ({len(batch)} puntos)")
        else:
            error_count += 1
    
    total_points = line_count * 10
    print(f"   ✅ Archivo completado: {line_count} líneas = {total_points} puntos en {batch_count} batches")
    if error_count > 0:
        print(f"   ⚠️  {error_count} batches con errores")
    
    return line_count, error_count

def send_batch(batch):
    """Envía un batch de datos a InfluxDB"""
    url = f"{INFLUX_URL}/api/v2/write?org={INFLUX_ORG}&bucket={INFLUX_BUCKET}&precision=ns"
    headers = {
        "Authorization": f"Token {INFLUX_TOKEN}",
        "Content-Type": "text/plain; charset=utf-8"
    }
    
    data = "\n".join(batch)
    
    try:
        response = requests.post(url, headers=headers, data=data, timeout=30)
        return response.status_code == 204
    except Exception as e:
        print(f"      Error enviando batch: {e}")
        return False

def main():
    """Función principal"""
    print("📊 Cargando datos históricos de sensores en InfluxDB...")
    
    # Esperar a que InfluxDB esté listo
    print("⏳ Verificando conexión con InfluxDB...")
    max_attempts = 30
    for attempt in range(1, max_attempts + 1):
        try:
            response = requests.get(f"{INFLUX_URL}/health", timeout=5)
            if response.status_code == 200:
                print("✅ InfluxDB está listo")
                break
        except:
            pass
        print(f"   Intento {attempt}/{max_attempts}...")
        import time
        time.sleep(2)
    else:
        print("❌ InfluxDB no está disponible")
        sys.exit(1)
    
    # Cargar archivos CSV
    csv_files = [
        "/scripts/dataset/A_2024-04-10.csv",
        "/scripts/dataset/A_2024-06-11.csv",
        "/scripts/dataset/A_2024-10-30.csv"
    ]
    
    total_lines = 0
    total_errors = 0
    files_loaded = 0
    
    for csv_file in csv_files:
        if Path(csv_file).exists():
            lines, errors = load_csv_to_influx(csv_file, activo_id=2)
            total_lines += lines
            total_errors += errors
            files_loaded += 1
        else:
            print(f"   ⚠️  Archivo no encontrado: {csv_file}")
    
    print(f"\n✅ {files_loaded} archivos CSV cargados en InfluxDB")
    print(f"📊 Total: {total_lines} líneas = {total_lines * 10} puntos de datos")
    if total_errors > 0:
        print(f"⚠️  {total_errors} batches con errores")

if __name__ == "__main__":
    main()
