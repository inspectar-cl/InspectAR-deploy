# 🧠 ML Engine Python - Detección de Anomalías HTM-like

Microservicio Python para detección de anomalías en tiempo real en activos industriales utilizando un algoritmo inspirado en Hierarchical Temporal Memory (HTM).

## 📋 Descripción

Este microservicio analiza datos de sensores en tiempo real y detecta anomalías utilizando:
- **Ventanas temporales deslizantes** para capturar patrones históricos
- **Suavizado exponencial (EMA)** para reducir ruido
- **Umbrales adaptativos** que se ajustan automáticamente
- **Clasificación de severidad** (Baja, Media, Alta)
- **Transformación automática** de formatos de entrada
- **Modelo incremental (SGDRegressor)** con persistencia

### Características principales

- ✅ Análisis multivariado (múltiples sensores simultáneamente)
- ✅ Detección en tiempo real
- ✅ Clasificación automática de severidad
- ✅ API REST con FastAPI
- ✅ Manejo robusto de datos incompletos/inválidos
- ✅ Configuración flexible mediante variables de entorno
- ✅ **Soporta múltiples formatos de entrada** (flat y sensor-agrupado)
- ✅ **Transformación automática** de datos del backend
- ✅ **Persistencia de modelo entrenado** entre sesiones
- ✅ **Limpieza adaptativa de timestamps** con múltiples estrategias
- ✅ **Interpolación automática** de datos faltantes

## 🏗️ Arquitectura

```
┌─────────────────┐
│  API Gateway    │
│  (Go Service)   │
└────────┬────────┘
         │ HTTP POST
         ▼
┌─────────────────┐
│   ML Engine     │◄─── config.py (parámetros)
│   (FastAPI)     │
└────────┬────────┘
         │
    ┌────┴─────┬──────────────┬──────────────┐
    ▼          ▼              ▼              ▼
┌────────┐  ┌──────────┐  ┌─────────────┐  ┌──────────┐
│ HTM    │  │ Cleaning │  │ Transform   │  │ Model    │
│ Model  │  │ Utils    │  │ Sensores    │  │ Persist  │
└────────┘  └──────────┘  └─────────────┘  └──────────┘
```

## 🐳 Despliegue con Docker (Producción)

### Prerrequisitos
- Docker instalado en WSL2
- Docker Compose (incluido en Docker Desktop)

### 1. Construir y ejecutar con Docker Compose

Desde el directorio raíz del backend:

```bash
cd backend
docker-compose up ml_engine_python -d
```

O construir manualmente:

```bash
cd backend/microservicios/ml_engine_python

# Construir imagen
docker build -t ml_engine_python:latest .

# Ejecutar contenedor
docker run -d \
  --name ml_engine \
  -p 8086:8086 \
  -v $(pwd)/app/models:/app/app/models \
  -e HTM_WINDOW_SIZE=180 \
  -e HTM_ALPHA=0.01 \
  -e HTM_K_ADAPT=2.0 \
  -e LOG_LEVEL=INFO \
  ml_engine_python:latest
```

**Nota:** El volumen `-v $(pwd)/app/models:/app/app/models` permite persistir el modelo entrenado entre reinicios del contenedor.

### 2. Verificar el servicio

```bash
# Ver logs
docker logs -f ml_engine

# Health check
curl http://localhost:8086/healthz

# Ver documentación
# Navega a http://localhost:8086/docs
```

### 3. Detener el servicio

```bash
docker stop ml_engine
docker rm ml_engine
```

## 🧪 Desarrollo y Pruebas Unitarias

### Prerrequisitos
- Python 3.9 o superior
- pip (gestor de paquetes)

### 1. Configurar entorno de desarrollo (Windows)

#### Opción A: PowerShell
```powershell
# Navegar al directorio
cd backend\microservicios\ml_engine_python

# Crear entorno virtual
python -m venv venv

# Activar entorno virtual
.\venv\Scripts\Activate.ps1

# Si hay error de permisos:
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Instalar dependencias
pip install -r requirements.txt
```

#### Opción B: CMD
```cmd
cd backend\microservicios\ml_engine_python
python -m venv venv
venv\Scripts\activate.bat
pip install -r requirements.txt
```

### 2. Ejecutar el servicio localmente

```powershell
# Asegúrate de que el entorno virtual esté activo
python -m uvicorn app.main:app --host 0.0.0.0 --port 8086 --reload
```

El servicio estará disponible en:
- API: `http://localhost:8086`
- Swagger UI: `http://localhost:8086/docs`
- ReDoc: `http://localhost:8086/redoc`

### 3. Ejecutar pruebas unitarias

#### Prueba rápida (genera datos sintéticos y prueba predicción)

```powershell
cd test
python quick_test.py
```

**Salida esperada:**
```
Generando datos sintéticos...
✅ 200 registros generados
⚠️  Anomalía inyectada en últimos 20 registros

Enviando request a ML Engine...

✅ Predicción exitosa!
Anomalías detectadas: 18 de 20

Último resultado:
{
  "timestamp": "2025-10-15T...",
  "AnomalyScore": 85.32,
  "AnomalyLikelihood": 92.15,
  "Threshold": 65.40,
  "is_anomaly": 1,
  "Severity": "Alta",
  "Description": "Anomalía crítica detectada..."
}
```

#### Suite completa de pruebas

```powershell
cd test
python test_ml_service.py
```

**Pruebas incluidas:**
1. ✅ Health Check - Verifica disponibilidad del servicio
2. ✅ Root Endpoint - Información del servicio
3. ✅ Get Config - Configuración actual
4. ✅ Predict Anomaly - Predicción con anomalía inyectada
5. ✅ Insufficient Data - Valida rechazo de datos insuficientes
6. ✅ Invalid Timestamps - Manejo de timestamps inválidos

#### Test con formato de sensores agrupados

```powershell
cd test
python test_ml_sensors.py
```

**Valida:**
- ✅ Conversión de formato sensor-agrupado a flat
- ✅ Manejo de datos faltantes (sparse data)
- ✅ Ordenamiento correcto de timestamps
- ✅ Persistencia del modelo entre llamadas
- ✅ Casos extremos (valores cero, negativos, etc.)

### 4. Pruebas manuales con curl

```bash
# Health check
curl http://localhost:8086/healthz

# Ver configuración
curl http://localhost:8086/config

# Predicción (formato flat - usar archivo JSON)
curl -X POST http://localhost:8086/predict_anomaly \
  -H "Content-Type: application/json" \
  -d @test_payload.json

# Predicción (formato backend - sensores agrupados)
curl -X POST http://localhost:8086/predict_anomaly_from_sensors \
  -H "Content-Type: application/json" \
  -d @test_payload_backend.json
```

Ejemplo `test_payload.json` (formato flat):
```json
{
  "pump": "bomba_1",
  "records": [
    {
      "Timestamp": "2025-10-15T10:00:00Z",
      "A_ACR_Mot.PV": 0.012,
      "A_Temp.PV": 75.2,
      "A_Pres.PV": 1.85
    }
    // ... más registros (mínimo 180)
  ]
}
```

Ejemplo `test_payload_backend.json` (formato backend):
```json
{
  "pump": 2,
  "sensores": [
    {
      "sensor_id": "A_ACR_Mot.PV",
      "datos": [
        {"tiempo": "2024-06-11T15:59:59Z", "valor": 0.001735421},
        {"tiempo": "2024-06-11T15:59:58Z", "valor": 0.001735421}
      ]
    },
    {
      "sensor_id": "A_Temp.PV",
      "datos": [
        {"tiempo": "2024-06-11T15:59:59Z", "valor": 31.99942017},
        {"tiempo": "2024-06-11T15:59:58Z", "valor": 31.99942017}
      ]
    }
    // ... más sensores (mínimo 180 timestamps únicos)
  ]
}
```

## ⚙️ Configuración

### Variables de entorno

| Variable | Descripción | Default |
|----------|-------------|---------|
| `HTM_WINDOW_SIZE` | Tamaño de ventana temporal | `180` |
| `HTM_ALPHA` | Factor de suavizado (0-1) | `0.01` |
| `HTM_K_ADAPT` | Sensibilidad del umbral | `2.0` |
| `HTM_SCORE_SCALE` | Escala de normalización | `50.0` |
| `HTM_MIN_PERIODS` | Períodos mínimos para threshold | `50` |
| `HTM_ROLLING_WINDOW` | Ventana rolling para threshold | `300` |
| `HTM_MODE` | Modo de operación | `offline` |
| `API_PORT` | Puerto del servicio | `8086` |
| `LOG_LEVEL` | Nivel de logging | `INFO` |
| `RANDOM_SEED` | Semilla para reproducibilidad | `42` |
| `SAVE_RESULTS` | Guardar resultados en archivos | `false` |
| `OUTPUT_PATH` | Carpeta para guardar resultados | `/app/results` |

### Ajustar parámetros

```bash
# En Docker
docker run -d \
  -e HTM_WINDOW_SIZE=240 \
  -e HTM_ALPHA=0.02 \
  -e HTM_K_ADAPT=2.5 \
  ml_engine_python:latest

# En desarrollo (crear .env)
echo "HTM_WINDOW_SIZE=240" > .env
echo "HTM_ALPHA=0.02" >> .env
echo "HTM_K_ADAPT=2.5" >> .env
```

## 📊 API Endpoints

### `GET /`
Información general del servicio.

**Response:**
```json
{
  "service": "Microservicio ML - Detección de Anomalías HTM-like",
  "version": "1.0.0",
  "status": "running",
  "mode": "offline",
  "config": {
    "window_size": 180,
    "alpha": 0.01,
    "k_adapt": 2.0
  }
}
```

### `GET /healthz`
Health check del servicio.

### `GET /config`
Configuración actual del modelo.

### `POST /predict_anomaly` ⭐ **[RECOMENDADO]**
Predicción de anomalías con **detección automática de formato**.

Este endpoint **unificado** acepta tanto formato **flat** como **sensor-agrupado** y detecta automáticamente cuál usar.

#### Formato 1: Flat (timestamp como índice)

**Request Body:**
```json
{
  "pump": "bomba_centrifuga_1",
  "records": [
    {
      "Timestamp": "2025-10-15T10:00:00Z",
      "A_ACR_Mot.PV": 0.012,
      "A_ACR_Mot.SV": 0.035,
      "A_ACR_Mot.TV": 41.2,
      "A_Pres.PV": 1.84,
      "A_Temp.PV": 22.7,
      "Barometer": 1013.3,
      "Temperature": 22.5
    }
    // ... mínimo 180 registros
  ]
}
```

#### Formato 2: Sensor-agrupado (formato backend)

**Request Body:**
```json
{
  "pump": 2,
  "sensores": [
    {
      "sensor_id": "A_ACR_Mot.PV",
      "datos": [
        {"tiempo": "2024-06-11T15:59:59Z", "valor": 0.001735421},
        {"tiempo": "2024-06-11T15:59:58Z", "valor": 0.001735421}
      ]
    },
    {
      "sensor_id": "A_Temp.PV",
      "datos": [
        {"tiempo": "2024-06-11T15:59:59Z", "valor": 31.99942017},
        {"tiempo": "2024-06-11T15:59:58Z", "valor": 31.99942017}
      ]
    }
    // ... más sensores (mínimo 180 timestamps únicos)
  ]
}
```

**Transformación automática:**
```
Entrada (sensor-agrupado):       →       Salida (flat por timestamp):
┌──────────────────────┐                ┌──────────────────────────┐
│ A_Temp.PV:           │                │ Timestamp: 15:59:58      │
│  - 15:59:58 → 31.98  │                │   A_Temp.PV: 31.98       │
│  - 15:59:59 → 31.99  │    ────────→   │   A_Pres.PV: 0.488       │
│ A_Pres.PV:           │                ├──────────────────────────┤
│  - 15:59:58 → 0.488  │                │ Timestamp: 15:59:59      │
│  - 15:59:59 → 0.490  │                │   A_Temp.PV: 31.99       │
└──────────────────────┘                │   A_Pres.PV: 0.490       │
                                        └──────────────────────────┘
```

**Response:**
```json
{
  "pump": "bomba_centrifuga_1",
  "n_rows": 20,
  "n_anomalies": 5,
  "results": [
    {
      "Timestamp": "2025-10-15T10:00:00Z",
      "AnomalyScore": 45.32,
      "AnomalyLikelihood": 48.15,
      "Threshold": 65.40,
      "is_anomaly": 0,
      "Severity": "Baja",
      "Description": "Funcionamiento dentro del rango esperado."
    }
    // ... más resultados
  ],
  "last": {
    "Timestamp": "2025-10-15T10:03:20Z",
    "AnomalyScore": 85.32,
    "AnomalyLikelihood": 92.15,
    "Threshold": 65.40,
    "is_anomaly": 1,
    "Severity": "Alta",
    "Description": "Anomalía crítica detectada. Atención prioritaria requerida."
  }
}
```

**Características:**
- ✅ **Detección automática** de formato (flat vs sensor-agrupado)
- ✅ Maneja datos faltantes (algunos sensores pueden no tener todos los timestamps)
- ✅ Ordena automáticamente por timestamp (más antiguos primero)
- ✅ Valida calidad de datos antes del análisis
- ✅ Limpieza adaptativa de timestamps inválidos
- ✅ Usa el mismo motor de detección HTM-like

---

### `POST /predict_anomaly_from_sensors` ⚠️ **[DEPRECATED]**

> **⚠️ DEPRECATED:** Este endpoint se mantiene solo por compatibilidad retroactiva.  
> **Usa `/predict_anomaly`** que ahora soporta ambos formatos automáticamente.

Este endpoint redirige internamente a `/predict_anomaly` y generará un warning en los logs:

```
WARNING: ⚠️ Endpoint deprecated: /predict_anomaly_from_sensors. Usa /predict_anomaly
```

**Migración recomendada:**

```diff
# Antes (deprecated)
- POST /predict_anomaly_from_sensors

# Ahora (recomendado)
+ POST /predict_anomaly
```

El payload es exactamente el mismo, solo cambia la URL.

## 🔧 Estructura del Proyecto

```
ml_engine_python/
├── app/
│   ├── __init__.py          # Inicialización del módulo + logger
│   ├── config.py            # Configuración con variables de entorno
│   ├── main.py              # FastAPI app + endpoints
│   ├── model_htm.py         # Modelo HTM-like
│   ├── utils.py             # Utilidades (limpieza, transformación, cálculos)
│   └── models/              # Directorio para modelos persistentes
│       ├── sgd_model.pkl    # Modelo SGDRegressor entrenado
│       ├── scaler_X.pkl     # Scaler para features
│       └── scaler_Y.pkl     # Scaler para targets
├── test/
│   ├── quick_test.py        # Prueba rápida
│   ├── test_ml_service.py   # Suite de tests
│   └── test_ml_sensors.py   # Tests de transformación
├── Dockerfile               # Imagen Docker
├── requirements.txt         # Dependencias Python
└── README.md               # Este archivo
```

## 🧮 Algoritmo HTM-like

### Flujo de procesamiento

1. **Entrada:** Ventana de W registros históricos (formato flat o sensor-agrupado)
2. **Detección de formato:** Automática (flat si tiene `records`, agrupado si tiene `sensores`)
3. **Transformación:** Si es formato backend, convierte a formato flat
4. **Limpieza adaptativa:** Estrategias para manejar timestamps inválidos
5. **Preprocesamiento:** Interpolación de valores faltantes
6. **Codificación temporal:** Uso de SGDRegressor incremental con persistencia
7. **Cálculo de error normalizado:** Comparación con media/std móvil
8. **Suavizado EMA:** `likelihood = (1-α) * likelihood + α * score`
9. **Umbral adaptativo:** `threshold = mean + k * std`
10. **Clasificación:** Comparación likelihood vs threshold
11. **Persistencia:** Guardado automático del modelo entrenado
12. **Salida:** Score, likelihood, severidad, descripción

### Parámetros clave

- **W (window_size):** Mayor = más contexto, menor = más reactivo
- **α (alpha):** Mayor = más reactivo, menor = más estable
- **k (k_adapt):** Mayor = menos sensible, menor = más sensible

## 🔄 Transformación y Limpieza de Datos

### Función: `transform_sensores_to_records`

Ubicada en `app/utils.py`, esta función convierte el formato sensor-agrupado del backend al formato flat requerido por el modelo HTM.

**Características:**
- Pivotea datos de sensor-centric a timestamp-centric
- Maneja datos faltantes (sparse data)
- Ordena cronológicamente (más antiguos primero)
- Preserva valores especiales (cero, negativos)
- Valida integridad de datos

### Función: `clean_timestamps`

Limpieza adaptativa de timestamps con múltiples estrategias según calidad de datos:

**Estrategias automáticas:**

1. **< 5% inválidos:** Eliminar registros inválidos
   ```
   Calidad: 98% → Acción: Eliminar 2% de registros
   ```

2. **5-20% inválidos:** Interpolar timestamps
   ```
   Calidad: 85% → Acción: Forward/backward fill
   ```

3. **> 20% inválidos:** Rechazar datos
   ```
   Calidad: 75% → Error: Calidad insuficiente
   ```

**Características:**
- ✅ Parseo flexible de formatos ISO8601/RFC3339
- ✅ Detección automática de frecuencia de muestreo
- ✅ Interpolación inteligente con forward/backward fill
- ✅ Logging detallado de estrategia aplicada
- ✅ Estadísticas de limpieza en respuesta

**Uso programático:**
```python
from app.utils import clean_timestamps

df, stats = clean_timestamps(df, max_invalid_percent=5.0)

# stats = {
#     "total_records": 200,
#     "invalid_timestamps": 10,
#     "dropped_records": 0,
#     "interpolated_records": 10,
#     "strategy": "interpolate",
#     "invalid_percent": 5.0
# }
```

### Función: `preprocess_timeseries`

Preprocesamiento final de series temporales:
- Ordenamiento por timestamp
- Interpolación de valores faltantes (NaN)
- Forward/backward fill en extremos

## 💾 Persistencia del Modelo

El servicio guarda automáticamente el modelo entrenado y los scalers después de cada análisis:

**Archivos guardados:**
- `app/models/sgd_model.pkl` - Modelo SGDRegressor entrenado
- `app/models/scaler_X.pkl` - StandardScaler para features
- `app/models/scaler_Y.pkl` - StandardScaler para targets

**Ventajas:**
- ✅ Continuidad entre reinicios del servicio
- ✅ Mejora progresiva con más datos
- ✅ Reducción de tiempo de warm-up
- ✅ Mejor detección de patrones históricos

**Nota:** En Docker, usa un volumen para persistir entre contenedores:
```bash
docker run -v $(pwd)/app/models:/app/app/models ml_engine_python:latest
```

## 🐛 Troubleshooting

### El servicio no inicia

```bash
# Verificar logs
docker logs ml_engine

# Verificar puerto
netstat -ano | findstr :8086  # Windows
lsof -i :8086                 # Linux/WSL
```

### Error: "Insufficient data"

```
Error: Se requieren al menos 180 registros para el análisis
```

**Solución:** 
- Para formato flat (`records`): Envía al menos 180 registros
- Para formato agrupado (`sensores`): Asegúrate de que haya al menos 180 timestamps únicos entre todos los sensores

### Error: "Calidad de datos insuficiente"

```
Error: Calidad de datos insuficiente: 25% de timestamps inválidos
```

**Solución:** 
- Verifica formato ISO 8601: `YYYY-MM-DDTHH:MM:SSZ`
- Máximo 20% de timestamps inválidos permitido
- Entre 5-20% se intentará interpolación automática

### Error: "No se pudieron generar registros"

```
Error: No se pudieron generar registros válidos a partir de los datos proporcionados
```

**Solución:**
- Verifica que estés usando `records` (flat) o `sensores` (agrupado)
- Para `sensores`: Cada sensor debe tener `sensor_id` y `datos`
- Para `sensores`: Cada dato debe tener `tiempo` y `valor`
- Asegúrate de que al menos un sensor tenga datos válidos

### Warning: Endpoint deprecated

```
WARNING: ⚠️ Endpoint deprecated: /predict_anomaly_from_sensors. Usa /predict_anomaly
```

**Solución:**
- Actualiza tu código para usar `/predict_anomaly` en lugar de `/predict_anomaly_from_sensors`
- El formato del payload es idéntico, solo cambia la URL
- El endpoint deprecated se eliminará en futuras versiones

### No detecta anomalías

**Posibles causas:**
- Datos muy estables (sin variación real)
- Parámetro `k_adapt` muy alto (reduce sensibilidad)
- Ventana muy corta
- Modelo aún en fase de aprendizaje

**Solución:** Ajusta parámetros:
```bash
# Más sensible
HTM_K_ADAPT=1.5
HTM_ALPHA=0.02
```

### Datos faltantes en algunos timestamps

**Comportamiento esperado:** El servicio maneja automáticamente datos faltantes (sparse data). Si un sensor no tiene valor para cierto timestamp, simplemente no aparecerá en ese registro.

**Ejemplo:**
```json
// Sensor A_Temp tiene datos en 3 timestamps
// Sensor A_Pres solo tiene datos en 2 timestamps
// Resultado: 3 registros, el segundo solo tendrá A_Temp
[
  {"Timestamp": "...", "A_Temp": 31.99, "A_Pres": 0.49},
  {"Timestamp": "...", "A_Temp": 32.00},  // ← A_Pres faltante
  {"Timestamp": "...", "A_Temp": 32.01, "A_Pres": 0.51}
]
```

### Modelo no mejora con más datos

**Verificar:**
1. ¿El modelo se está guardando? Revisa logs: `✓ Modelo guardado correctamente`
2. ¿Existe el volumen en Docker? Usa `-v` para persistir
3. ¿Los datos son consistentes? Variaciones grandes resetean el aprendizaje

**Solución:**
```bash
# En Docker, verificar volumen
docker inspect ml_engine | grep Mounts

# Manualmente guardar modelo
curl -X POST http://localhost:8086/predict_anomaly -d @data.json
# Revisar logs: [INFO] Modelo guardado en /app/app/models/...
```

## 📈 Monitoreo y Logs

### Niveles de log

```bash
# Desarrollo (verbose)
LOG_LEVEL=DEBUG

# Producción (estándar)
LOG_LEVEL=INFO

# Producción (errores solamente)
LOG_LEVEL=ERROR
```

### Ver logs en Docker

```bash
# Logs en tiempo real
docker logs -f ml_engine

# Últimas 100 líneas
docker logs --tail 100 ml_engine

# Logs con timestamps
docker logs -t ml_engine
```

### Logs relevantes

```
INFO: Procesando predicción para activo: 2
INFO: 📦 Formato detectado: AGRUPADO POR SENSORES (10 sensores)
DEBUG: ✓ Transformación completada: 180 registros generados
INFO: 🧹 Limpieza de datos - Estrategia: interpolate, Calidad: 98.5%
DEBUG: Features detectadas: 10 columnas
INFO: ✓ Modelo guardado correctamente en /app/app/models/
INFO: ✅ Predicción completada para 2: 5 anomalías detectadas de 180 puntos
```

### Logs de limpieza de timestamps

```
WARNING: Timestamps inválidos: 10 de 200 (5.00%)
INFO: ✓ Estrategia: Interpolar todos (10 timestamps interpolados)
```

```
WARNING: Timestamps inválidos: 3 de 200 (1.50%)
INFO: ✓ Estrategia: Eliminar registros inválidos (3 registros, 1.50%)
```

```
ERROR: Calidad de datos insuficiente: 22.0% de timestamps inválidos
```

### Logs de endpoint deprecated

```
WARNING: ⚠️ Endpoint deprecated: /predict_anomaly_from_sensors. Usa /predict_anomaly
```

## 🔐 Seguridad

- ✅ Validación de entrada con Pydantic
- ✅ Límites de tamaño en requests
- ✅ Sanitización de timestamps
- ✅ Manejo seguro de errores
- ✅ Validación de estructura de datos
- ✅ Validación de calidad de timestamps
- ⚠️ **Nota:** Este servicio debe estar detrás del API Gateway, no expuesto directamente

## 🚀 Próximas mejoras

- [x] Persistencia de modelos entrenados
- [x] Limpieza adaptativa de timestamps
- [x] Interpolación automática de datos
- [x] Endpoint unificado con detección automática de formato
- [ ] Caché de modelos por activo
- [ ] Soporte para múltiples activos simultáneos en batch
- [ ] Métricas de Prometheus
- [ ] Reentrenamiento automático periódico
- [ ] Soporte para datos streaming (MQTT)
- [ ] Tests de integración con InfluxDB
- [ ] Endpoint para análisis histórico de tendencias
- [ ] Dashboard de métricas del modelo
- [ ] Eliminación completa de endpoint deprecated (v2.0.0)

## 📝 Licencia

Proyecto InspectAR - Uso interno

## 👥 Contacto

Para preguntas o soporte, contactar al equipo de desarrollo de InspectAR.