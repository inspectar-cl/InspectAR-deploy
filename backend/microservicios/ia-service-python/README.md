# IA Service Python

Microservicio de detección de anomalías utilizando Machine Learning.

## 📋 Descripción

Este servicio orquesta la detección de anomalías combinando datos de sensores IoT con algoritmos de Machine Learning. Reemplaza al servicio original escrito en Go.

## 🏗️ Arquitectura

```
┌─────────────┐
│  API REST   │ ← FastAPI (puerto 8095)
└──────┬──────┘
       │
┌──────▼──────────────────────────┐
│  Handlers (Controladores)       │
│  - anomaly_handler.py            │
└──────┬──────────────────────────┘
       │
┌──────▼──────────────────────────┐
│  Services (Lógica de negocio)   │
│  - anomaly_service.py            │
│    · Orquesta llamadas externas  │
│    · Transforma datos            │
│    · Predicción con modelo local │
│  - training_service.py           │
│    · Entrenamiento programado    │
└──────┬──────────────────────────┘
       │
┌──────▼──────────────────────────┐
│  Models (Machine Learning)      │
│  - model_manager.py              │
│    · Carga/guardado de modelo    │
│    · SGDRegressor + Scalers      │
└──────┬──────────────────────────┘
       │
┌──────▼──────────────────────────┐
│  Repositories (Acceso a datos)  │
│  - anomaly_repository.py         │
│    · CRUD en PostgreSQL          │
└─────────────┬────────────────────┘
              │
       ┌──────▼──────┐
       │  PostgreSQL │
       │   (ia_db)   │
       └─────────────┘
```

## 🚀 Características

- ✅ API REST con FastAPI
- ✅ **Modelo de ML integrado** (SGDRegressor con StandardScaler)
- ✅ Detección de anomalías sin dependencias externas
- ✅ Entrenamiento automático programable (APScheduler)
- ✅ Integración con IOT Service (ParserService)
- ✅ Almacenamiento en PostgreSQL
- ✅ Paginación y filtros
- ✅ Documentación automática (Swagger/ReDoc)
- ✅ Health checks
- ✅ Logging estructurado
- ✅ Multi-stage Docker build

## 📁 Estructura

```
ia-service-python/
├── app/
│   ├── __init__.py
│   ├── main.py              # Aplicación FastAPI
│   ├── config.py            # Configuración
│   ├── database.py          # Conexión PostgreSQL
│   │
│   ├── models/              # Modelos de datos y ML
│   │   ├── __init__.py
│   │   ├── anomaly.py       # Modelos Pydantic
│   │   └── model_manager.py # Gestor de modelo ML
│   │
│   ├── repositories/        # Capa de acceso a datos
│   │   ├── __init__.py
│   │   └── anomaly_repository.py
│   │
│   ├── services/            # Lógica de negocio
│   │   ├── __init__.py
│   │   ├── anomaly_service.py
│   │   └── training_service.py  # Entrenamiento automático
│   │
│   └── handlers/            # Controladores HTTP
│       ├── __init__.py
│       └── anomaly_handler.py
│
├── config/
│   ├── config.yaml          # Config desarrollo
│   └── config.docker.yaml   # Config Docker
│
├── models/                  # Modelos ML persistidos
│   ├── anomaly_model.pkl
│   ├── scaler_X.pkl
│   └── scaler_Y.pkl
│
├── requirements.txt
├── Dockerfile
└── README.md
```

## 🔧 Endpoints

### Anomalías

#### `GET /anomalies/activo/{activo_id}`

Obtiene anomalías de un activo específico con paginación.

**Parámetros de ruta:**
- `activo_id` (int): ID del activo

**Parámetros de query:**
- `limit` (int, opcional): Número de resultados (default: 50, max: 1000)
- `offset` (int, opcional): Offset para paginación (default: 0)
- `only_anomalies` (bool, opcional): Solo anomalías confirmadas (default: false)

**Ejemplo:**
```bash
curl "http://localhost:8095/anomalies/activo/1?limit=10&offset=0&only_anomalies=true"
```

**Respuesta:**
```json
{
  "message": "Anomalías obtenidas exitosamente",
  "count": 10,
  "activo_id": 1,
  "limit": 10,
  "offset": 0,
  "data": [
    {
      "id": 123,
      "activo_id": 1,
      "sensor_id": null,
      "timestamp": "2025-10-28T10:30:00",
      "anomaly_score": 0.85,
      "anomaly_likelihood": 85.0,
      "severidad": "alta",
      "descripcion": "Anomalía detectada en activo 1",
      "threshold": 0.5,
      "is_anomaly": 1
    }
  ]
}
```

#### `GET /anomalies/sensor/{sensor_id}`

Obtiene anomalías de un sensor específico.

**Parámetros:**
- `sensor_id` (str): Identificador del sensor
- `limit` (int, opcional): Número de resultados (default: 10)
- `offset` (int, opcional): Offset para paginación (default: 0)

#### `GET /anomalies/{anomaly_id}`

Obtiene una anomalía específica por ID.

#### `POST /detect/{activo_id}`

Dispara manualmente la detección de anomalías para un activo.

**Parámetros de query:**
- `page` (int, opcional): Página de datos (default: 1)
- `limit` (int, opcional): Cantidad de registros (default: 1000, max: 5000)
- `async_mode` (bool, opcional): Ejecutar en segundo plano (default: false)

**Ejemplo:**
```bash
curl -X POST "http://localhost:8095/detect/1?limit=1000"
```

**Respuesta:**
```json
{
  "activo_id": 1,
  "status": "success",
  "data_points": 1000,
  "results_saved": 1000,
  "anomalies_detected": 15
}
```

### Entrenamiento

#### `GET /training/status`

Obtiene el estado del entrenamiento y scheduler.

#### `POST /training/train`

Dispara un entrenamiento manual del modelo.

#### `GET /training/model`

Obtiene información sobre el modelo actual.

### Sistema

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/health` | Health check del servicio |
| `GET` | `/status` | Estado y estadísticas del servicio |
| `GET` | `/` | Información general del servicio |
| `GET` | `/docs` | Documentación Swagger UI |
| `GET` | `/redoc` | Documentación ReDoc |

## 🐳 Docker

### Build

```bash
docker build -t ia-service-python:latest .
```

### Run

```bash
docker run -d \
  -p 8095:8095 \
  -e DATABASE_HOST=ia-db \
  -e DATABASE_PORT=5432 \
  -e DATABASE_NAME=ia_db \
  -e DATABASE_USER=ia_user \
  -e DATABASE_PASSWORD=ia_pass \
  -e IOT_SERVICE_URL=http://iot-service:8090 \
  -e CONFIG_FILE=/app/config/config.yaml \
  -v $(pwd)/config:/app/config:ro \
  -v ia_service_models:/app/models \
  --name ia-service-python \
  ia-service-python:latest
```

## 🔨 Desarrollo Local

### Requisitos

- Python 3.12+
- PostgreSQL 16
- pip

### Instalación

```bash
# Crear entorno virtual
python -m venv venv
source venv/bin/activate  # Linux/Mac
# venv\Scripts\activate  # Windows

# Instalar dependencias
pip install -r requirements.txt

# Configurar variables de entorno
export DATABASE_HOST=localhost
export DATABASE_PORT=5436
export CONFIG_FILE=config/config.yaml

# Ejecutar
python -m app.main
```

## 📊 Variables de Entorno

| Variable | Descripción | Default |
|----------|-------------|---------|
| `DATABASE_HOST` | Host de PostgreSQL | `ia-db` |
| `DATABASE_PORT` | Puerto de PostgreSQL | `5432` |
| `DATABASE_NAME` | Nombre de la BD | `ia_db` |
| `DATABASE_USER` | Usuario de la BD | `ia_user` |
| `DATABASE_PASSWORD` | Contraseña de la BD | `ia_pass` |
| `IOT_SERVICE_URL` | URL del IOT Service | `http://iot-service:8090` |
| `LOG_LEVEL` | Nivel de logging | `INFO` |
| `CONFIG_FILE` | Archivo de configuración YAML | `None` |

**Nota:** La configuración de entrenamiento (`training_interval_minutes`) se maneja en `config/config.yaml`

## � Modelo de Machine Learning

### Arquitectura del Modelo

Este servicio utiliza un modelo de **regresión incremental** para detectar anomalías:

- **Algoritmo:** `SGDRegressor` (Stochastic Gradient Descent Regressor)
- **Preprocesamiento:** `StandardScaler` para features (X) y target (Y)
- **Método de detección:** Error de reconstrucción

### Funcionamiento

1. **Entrenamiento:**
   - El modelo aprende de datos históricos de sensores
   - Se entrena de forma incremental con `warm_start=True`
   - Entrenamiento programado cada `training_interval_minutes` (configurable)

2. **Predicción:**
   - Recibe valores de sensores como features
   - Predice el valor esperado
   - Calcula el error de reconstrucción: `|valor_real - valor_predicho|`
   - Si el error supera un threshold → **Anomalía detectada**

3. **Scoring:**
   - `anomaly_score`: Error normalizado (0-1)
   - `anomaly_likelihood`: Probabilidad porcentual (0-100)
   - `severidad`: Clasificación basada en score
     - `< 0.3` → baja
     - `0.3 - 0.6` → media
     - `> 0.6` → alta

### Persistencia

Los modelos se guardan en el volumen Docker:

```
/app/models/
├── anomaly_model.pkl    # Modelo SGDRegressor
├── scaler_X.pkl         # Scaler para features
└── scaler_Y.pkl         # Scaler para target
```

### Configuración de Entrenamiento

En `config/config.yaml`:

```yaml
training_enabled: true
training_interval_minutes: 60  # Cada 1 hora
model_dir: "/app/models"
model_filename: "anomaly_model.pkl"
```

## �🧪 Testing

```bash
# Tests unitarios
pytest tests/

# Tests con coverage
pytest --cov=app tests/
```

## 📝 Logging

El servicio utiliza logging estructurado:

```
[2025-10-16 08:30:00] [INFO] app.main: 🚀 Iniciando IA Service Python...
[2025-10-16 08:30:01] [INFO] app.database: ✅ Conexión a PostgreSQL exitosa
[2025-10-16 08:30:05] [INFO] app.services.anomaly_service: 📡 Solicitando datos del ParserService
[2025-10-16 08:30:10] [INFO] app.services.anomaly_service: ✅ Datos recibidos: 10 sensores
```

## 🔄 Flujo de Detección

```
1. Cliente → POST /detect/{activo_id}
2. Handler → Service.detect_and_store_anomalies()
3. Service → ParserService.get_data()
   └─ GET /lectura/{activo_id}/window
4. Service → transform_to_ml_format()
   └─ Agrupa por timestamp
5. Service → predict_anomalies() [MODELO LOCAL]
   ├─ model_manager.load_model()
   ├─ scaler_X.transform(features)
   ├─ model.predict(features)
   ├─ calcular error de reconstrucción
   └─ determinar is_anomaly (threshold)
6. Service → Repository.bulk_create()
7. Repository → PostgreSQL INSERT
8. Service → Cliente (response)
```

### Diferencia con ia-service (Go)

**ia-service-python** tiene el modelo **integrado** en el mismo microservicio:

- ✅ **Sin dependencias externas** de ml_engine_python
- ✅ **Predicción local** con SGDRegressor
- ✅ **Entrenamiento automático** programado
- ✅ **Persistencia** en volumen Docker

**ia-service (Go)** solo **consulta** anomalías ya almacenadas:

- 📊 GET /anomalies/activo/:activo_id (lectura desde BD)
- 📊 No hace predicciones ni entrenamiento

## 🤝 Contribución

Este microservicio sigue los estándares del proyecto InspectAR.

## 📄 Licencia

Propiedad de InspectAR Team.
