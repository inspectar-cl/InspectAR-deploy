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
- ✅ Integración con ML Engine (HTM)
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
│   ├── models/              # Modelos de datos
│   │   ├── __init__.py
│   │   └── anomaly.py
│   │
│   ├── repositories/        # Capa de acceso a datos
│   │   ├── __init__.py
│   │   └── anomaly_repository.py
│   │
│   ├── services/            # Lógica de negocio
│   │   ├── __init__.py
│   │   └── anomaly_service.py
│   │
│   └── handlers/            # Controladores HTTP
│       ├── __init__.py
│       └── anomaly_handler.py
│
├── config/
│   ├── config.yaml          # Config desarrollo
│   └── config.docker.yaml   # Config Docker
│
├── requirements.txt
├── Dockerfile
└── README.md
```

## 🔧 Endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| **Anomalías** | | |
| `GET` | `/anomalies/activo/{activo_id}` | Listar anomalías por activo (con paginación) |
| `GET` | `/anomalies/sensor/{sensor_id}` | Listar anomalías por sensor |
| `GET` | `/anomalies/{anomaly_id}` | Obtener anomalía específica |
| `POST` | `/detect/{activo_id}` | Detectar anomalías manualmente (sync/async) |
| **Entrenamiento** | | |
| `GET` | `/training/status` | Estado del entrenamiento y scheduler |
| `POST` | `/training/train` | Disparar entrenamiento manual |
| `GET` | `/training/model` | Información del modelo actual |
| **Sistema** | | |
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

## 🧪 Testing

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
4. Service → transform_to_ml_format()
5. Service → MLEngine.predict_anomaly()
6. Service → Repository.bulk_create()
7. Repository → PostgreSQL INSERT
8. Service → Cliente (response)
```

## 🤝 Contribución

Este microservicio sigue los estándares del proyecto InspectAR.

## 📄 Licencia

Propiedad de InspectAR Team.
