"""
Script de prueba standalone para el servicio ML Engine.
Genera datos sintéticos y prueba los endpoints.

Uso:
    python test_ml_service.py
"""

import requests
import json
import pandas as pd
import numpy as np
from datetime import datetime, timedelta
import time

# ============================================================
# Configuración
# ============================================================

BASE_URL = "http://localhost:8085"  # Ajusta el puerto según tu config
PUMP_ID = "bomba_test_1"
NUM_RECORDS = 200  # Más del mínimo requerido (180)

# ============================================================
# Generador de datos sintéticos
# ============================================================

def generate_synthetic_data(n_records=200, add_anomaly=True, invalid_timestamps=0):
    """
    Genera datos sintéticos de sensores de una bomba.
    
    Args:
        n_records: Número de registros a generar
        add_anomaly: Si True, inyecta una anomalía en los últimos datos
        invalid_timestamps: Número de timestamps inválidos a inyectar (para testing)
    
    Returns:
        Lista de diccionarios con datos de sensores
    """
    print(f"\n📊 Generando {n_records} registros sintéticos...")
    
    start_time = datetime.now() - timedelta(seconds=n_records)
    records = []
    
    for i in range(n_records):
        timestamp = start_time + timedelta(seconds=i)
        
        # Valores base con variación normal
        base_temp = 75.0
        base_pressure = 2.0
        base_vibration = 0.01
        
        # Agregar ruido normal
        noise_temp = np.random.normal(0, 2)
        noise_pressure = np.random.normal(0, 0.1)
        noise_vibration = np.random.normal(0, 0.005)
        
        # Inyectar anomalía en los últimos 20 registros
        if add_anomaly and i >= n_records - 20:
            anomaly_factor = (i - (n_records - 20)) / 20
            noise_temp += 15 * anomaly_factor  # Temperatura sube
            noise_vibration += 0.05 * anomaly_factor  # Vibración aumenta
        
        record = {
            "timestamp": timestamp.isoformat() + "Z",
            "A_ACR_Mot.PV": round(base_vibration + noise_vibration, 4),
            "A_ACR_Mot.SV": round(0.035 + np.random.normal(0, 0.002), 4),
            "A_ACR_Mot.TV": round(40.0 + np.random.normal(0, 1.5), 2),
            "A_ACR_Pmp.PV": round(base_vibration * 0.9 + noise_vibration * 0.8, 4),
            "A_ACR_Pmp.SV": round(0.031 + np.random.normal(0, 0.002), 4),
            "A_ACR_Pmp.TV": round(42.0 + np.random.normal(0, 1.5), 2),
            "A_Pres.PV": round(base_pressure + noise_pressure, 2),
            "A_Temp.PV": round(base_temp + noise_temp, 1),
            "Barometer": round(1013.0 + np.random.normal(0, 0.5), 1),
            "Temperature": round(22.0 + np.random.normal(0, 0.3), 1)
        }
        records.append(record)
    
    # Inyectar timestamps inválidos para testing
    if invalid_timestamps > 0:
        print(f"⚠️  Inyectando {invalid_timestamps} timestamps inválidos...")
        invalid_indices = np.random.choice(n_records, invalid_timestamps, replace=False)
        for idx in invalid_indices:
            records[idx]["timestamp"] = "INVALID_TIMESTAMP"
    
    print(f"✅ Datos generados: {len(records)} registros")
    if add_anomaly:
        print(f"⚠️  Anomalía inyectada en últimos 20 registros")
    
    return records


# ============================================================
# Funciones de prueba
# ============================================================

def test_health_check():
    """Prueba el endpoint de health check"""
    print("\n" + "="*60)
    print("🧪 Test 1: Health Check")
    print("="*60)
    
    try:
        response = requests.get(f"{BASE_URL}/healthz", timeout=5)
        print(f"Status Code: {response.status_code}")
        print(f"Response: {json.dumps(response.json(), indent=2)}")
        return response.status_code == 200
    except Exception as e:
        print(f"❌ Error: {e}")
        return False


def test_root():
    """Prueba el endpoint raíz"""
    print("\n" + "="*60)
    print("🧪 Test 2: Root Endpoint")
    print("="*60)
    
    try:
        response = requests.get(f"{BASE_URL}/", timeout=5)
        print(f"Status Code: {response.status_code}")
        print(f"Response: {json.dumps(response.json(), indent=2)}")
        return response.status_code == 200
    except Exception as e:
        print(f"❌ Error: {e}")
        return False


def test_get_config():
    """Prueba el endpoint de configuración"""
    print("\n" + "="*60)
    print("🧪 Test 3: Get Config")
    print("="*60)
    
    try:
        response = requests.get(f"{BASE_URL}/config", timeout=5)
        print(f"Status Code: {response.status_code}")
        print(f"Response: {json.dumps(response.json(), indent=2)}")
        return response.status_code == 200
    except Exception as e:
        print(f"❌ Error: {e}")
        return False


def test_predict_anomaly(records, pump_id=PUMP_ID):
    """Prueba el endpoint de predicción de anomalías"""
    print("\n" + "="*60)
    print("🧪 Test 4: Predict Anomaly")
    print("="*60)
    
    payload = {
        "pump": pump_id,
        "records": records
    }
    
    print(f"Enviando {len(records)} registros para bomba '{pump_id}'...")
    
    try:
        response = requests.post(
            f"{BASE_URL}/predict_anomaly",
            json=payload,
            timeout=30
        )
        
        print(f"Status Code: {response.status_code}")
        
        if response.status_code == 200:
            result = response.json()
            print(f"\n✅ Predicción exitosa!")
            print(f"  - Bomba: {result['pump']}")
            print(f"  - Registros analizados: {result['n_rows']}")
            print(f"  - Anomalías detectadas: {result['n_anomalies']}")
            
            if result['last']:
                print(f"\n📊 Último resultado:")
                last = result['last']
                print(f"  - Timestamp: {last['timestamp']}")
                print(f"  - AnomalyScore: {last['AnomalyScore']}")
                print(f"  - AnomalyLikelihood: {last['AnomalyLikelihood']}")
                print(f"  - Threshold: {last['Threshold']}")
                print(f"  - Severity: {last['Severity']}")
                print(f"  - Is Anomaly: {last['is_anomaly']}")
                print(f"  - Description: {last['Description']}")
            
            # Mostrar estadísticas de severidad
            if result['results']:
                severities = [r['Severity'] for r in result['results']]
                from collections import Counter
                severity_counts = Counter(severities)
                print(f"\n📈 Distribución de severidad:")
                for sev, count in severity_counts.items():
                    print(f"  - {sev}: {count}")
            
            return True
        else:
            print(f"❌ Error en predicción:")
            print(f"Response: {response.text}")
            return False
            
    except Exception as e:
        print(f"❌ Error: {e}")
        return False


def test_insufficient_data():
    """Prueba con datos insuficientes (debe fallar)"""
    print("\n" + "="*60)
    print("🧪 Test 5: Insufficient Data (debe fallar)")
    print("="*60)
    
    records = generate_synthetic_data(n_records=50, add_anomaly=False)
    
    payload = {
        "pump": "bomba_test_insuficiente",
        "records": records
    }
    
    try:
        response = requests.post(
            f"{BASE_URL}/predict_anomaly",
            json=payload,
            timeout=10
        )
        
        print(f"Status Code: {response.status_code}")
        
        if response.status_code == 400:
            print(f"✅ Validación correcta - Error esperado:")
            print(f"Response: {response.json()}")
            return True
        else:
            print(f"❌ Se esperaba error 400, pero se obtuvo: {response.status_code}")
            return False
            
    except Exception as e:
        print(f"❌ Error: {e}")
        return False


def test_invalid_timestamps():
    """Prueba con algunos timestamps inválidos"""
    print("\n" + "="*60)
    print("🧪 Test 6: Invalid Timestamps (5% inválidos)")
    print("="*60)
    
    records = generate_synthetic_data(n_records=200, add_anomaly=True, invalid_timestamps=10)
    
    payload = {
        "pump": "bomba_test_timestamps",
        "records": records
    }
    
    try:
        response = requests.post(
            f"{BASE_URL}/predict_anomaly",
            json=payload,
            timeout=30
        )
        
        print(f"Status Code: {response.status_code}")
        
        if response.status_code == 200:
            result = response.json()
            print(f"✅ Predicción exitosa a pesar de timestamps inválidos!")
            print(f"  - Registros analizados: {result['n_rows']}")
            print(f"  - Anomalías detectadas: {result['n_anomalies']}")
            return True
        else:
            print(f"Response: {response.json()}")
            return False
            
    except Exception as e:
        print(f"❌ Error: {e}")
        return False


# ============================================================
# Función principal
# ============================================================

def run_all_tests():
    """Ejecuta todos los tests"""
    print("\n" + "="*60)
    print("🚀 INICIANDO SUITE DE PRUEBAS - ML ENGINE")
    print("="*60)
    print(f"Base URL: {BASE_URL}")
    print(f"Timestamp: {datetime.now().isoformat()}")
    
    results = {}
    
    # Test 1: Health Check
    results['health_check'] = test_health_check()
    time.sleep(0.5)
    
    # Test 2: Root
    results['root'] = test_root()
    time.sleep(0.5)
    
    # Test 3: Config
    results['config'] = test_get_config()
    time.sleep(0.5)
    
    # Test 4: Predicción normal (con anomalía)
    records = generate_synthetic_data(n_records=NUM_RECORDS, add_anomaly=True)
    results['predict_with_anomaly'] = test_predict_anomaly(records)
    time.sleep(0.5)
    
    # Test 5: Datos insuficientes
    results['insufficient_data'] = test_insufficient_data()
    time.sleep(0.5)
    
    # Test 6: Timestamps inválidos
    results['invalid_timestamps'] = test_invalid_timestamps()
    
    # Resumen
    print("\n" + "="*60)
    print("📊 RESUMEN DE PRUEBAS")
    print("="*60)
    
    for test_name, passed in results.items():
        status = "✅ PASS" if passed else "❌ FAIL"
        print(f"{status} - {test_name}")
    
    total = len(results)
    passed = sum(results.values())
    print(f"\nTotal: {passed}/{total} pruebas pasadas")
    
    if passed == total:
        print("\n🎉 ¡Todos los tests pasaron exitosamente!")
    else:
        print(f"\n⚠️  {total - passed} test(s) fallaron")
    
    return passed == total


# ============================================================
# Entry point
# ============================================================

if __name__ == "__main__":
    print("🔬 Script de prueba para ML Engine Service")
    print("Asegúrate de que el servicio esté corriendo en", BASE_URL)
    
    input("\nPresiona ENTER para comenzar las pruebas...")
    
    success = run_all_tests()
    
    exit(0 if success else 1)