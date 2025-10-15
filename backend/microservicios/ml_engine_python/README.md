# 🧠 ML Engine Python - Detección de Anomalías HTM-like

Microservicio Python para detección de anomalías en tiempo real en activos industriales utilizando un algoritmo inspirado en Hierarchical Temporal Memory (HTM).

## 📋 Descripción

Este microservicio analiza datos de sensores en tiempo real y detecta anomalías utilizando:
- **Ventanas temporales deslizantes** para capturar patrones históricos
- **Suavizado exponencial (EMA)** para reducir ruido
- **Umbrales adaptativos** que se ajustan automáticamente
- **Clasificación de severidad** (Baja, Media, Alta)

### Características principales

- ✅ Análisis multivariado (múltiples sensores simultáneamente)
- ✅ Detección en tiempo real
- ✅ Clasificación automática de severidad
- ✅ API REST con FastAPI
- ✅ Manejo robusto de datos incompletos/inválidos
- ✅ Configuración flexible mediante variables de entorno

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
    ┌────┴─────┐
    ▼          ▼
┌────────┐  ┌──────────┐
│ HTM    │  │ Cleaning │
│ Model  │  │ Utils    │
└────────┘  └──────────┘
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
  -p 8085:8085 \
  -e HTM_WINDOW_SIZE=180 \
  -e HTM_ALPHA=0.01 \
  -e HTM_K_ADAPT=2.0 \
  -e LOG_LEVEL=INFO \
  ml_engine_python:latest
```

### 2. Verificar el servicio

```bash
# Ver logs
docker logs -f ml_engine

# Health check
curl http://localhost:8085/healthz

# Ver documentación
# Navega a http://localhost:8085/docs
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
python -m uvicorn app.main:app --host 0.0.0.0 --port 8085 --reload
```

El servicio estará disponible en:
- API: `http://localhost:8085`
- Swagger UI: `http://localhost:8085/docs`
- ReDoc: `http://localhost:8085/redoc`

### 3. Ejecutar pruebas unitarias

#### Prueba rápida (genera datos sintéticos y prueba predicción)

```powershell
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
python test_ml_service.py
```

**Pruebas incluidas:**
1. ✅ Health Check - Verifica disponibilidad del servicio
2. ✅ Root Endpoint - Información del servicio
3. ✅ Get Config - Configuración actual
4. ✅ Predict Anomaly - Predicción con anomalía inyectada
5. ✅ Insufficient Data - Valida rechazo de datos insuficientes
6. ✅ Invalid Timestamps - Manejo de timestamps inválidos

### 4. Pruebas manuales con curl

```bash
# Health check
curl http://localhost:8085/healthz

# Ver configuración
curl http://localhost:8085/config

# Predicción (usar archivo JSON)
curl -X POST http://localhost:8085/predict_anomaly \
  -H "Content-Type: application/json" \
  -d @test_payload.json
```

Ejemplo `test_payload.json`:
```json
{
  "pump": "bomba_1",
  "records": [
    {
      "timestamp": "2025-10-15T10:00:00Z",
      "A_ACR_Mot.PV": 0.012,
      "A_Temp.PV": 75.2,
      "A_Pres.PV": 1.85
    }
    // ... más registros (mínimo 180)
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
| `API_PORT` | Puerto del servicio | `8085` |
| `LOG_LEVEL` | Nivel de logging | `INFO` |
| `RANDOM_SEED` | Semilla para reproducibilidad | `42` |

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

### `POST /predict_anomaly`
Predicción de anomalías.

**Request Body:**
```json
{
  "pump": "bomba_centrifuga_1",
  "records": [
    {
      "timestamp": "2025-10-15T10:00:00Z",
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

**Response:**
```json
{
  "pump": "bomba_centrifuga_1",
  "n_rows": 20,
  "n_anomalies": 5,
  "results": [
    {
      "timestamp": "2025-10-15T10:00:00Z",
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
    "timestamp": "2025-10-15T10:03:20Z",
    "AnomalyScore": 85.32,
    "AnomalyLikelihood": 92.15,
    "Threshold": 65.40,
    "is_anomaly": 1,
    "Severity": "Alta",
    "Description": "Anomalía crítica detectada. Atención prioritaria requerida."
  }
}
```

## 🔧 Estructura del Proyecto

```
ml_engine_python/
├── app/
│   ├── __init__.py          # Inicialización del módulo + logger
│   ├── config.py            # Configuración con variables de entorno
│   ├── main.py              # FastAPI app + endpoints
│   ├── model_htm.py         # Modelo HTM-like
│   └── utils.py             # Utilidades (limpieza, cálculos)
├── Dockerfile               # Imagen Docker
├── requirements.txt         # Dependencias Python
├── quick_test.py           # Prueba rápida
├── test_ml_service.py      # Suite de tests
└── README.md               # Este archivo
```

## 🧮 Algoritmo HTM-like

### Flujo de procesamiento

1. **Entrada:** Ventana de W registros históricos
2. **Preprocesamiento:** Limpieza de timestamps inválidos
3. **Codificación temporal:** Uso de SGDRegressor incremental
4. **Cálculo de error normalizado:** Comparación con media/std móvil
5. **Suavizado EMA:** `likelihood = (1-α) * likelihood + α * score`
6. **Umbral adaptativo:** `threshold = mean + k * std`
7. **Clasificación:** Comparación likelihood vs threshold
8. **Salida:** Score, likelihood, severidad, descripción

### Parámetros clave

- **W (window_size):** Mayor = más contexto, menor = más reactivo
- **α (alpha):** Mayor = más reactivo, menor = más estable
- **k (k_adapt):** Mayor = menos sensible, menor = más sensible

## 🐛 Troubleshooting

### El servicio no inicia

```bash
# Verificar logs
docker logs ml_engine

# Verificar puerto
netstat -ano | findstr :8085  # Windows
lsof -i :8085                 # Linux/WSL
```

### Error: "Insufficient data"

```
Error: Se requieren al menos 180 registros para el análisis
```

**Solución:** Envía al menos 180 registros en el campo `records`.

### Error: "Invalid timestamps"

```
Error: Calidad de datos insuficiente: 25% de timestamps inválidos
```

**Solución:** 
- Verifica formato ISO 8601: `YYYY-MM-DDTHH:MM:SSZ`
- Máximo 20% de timestamps inválidos permitido
- Entre 5-20% se intentará interpolación automática

### No detecta anomalías

**Posibles causas:**
- Datos muy estables (sin variación real)
- Parámetro `k_adapt` muy alto (reduce sensibilidad)
- Ventana muy corta

**Solución:** Ajusta parámetros:
```bash
# Más sensible
HTM_K_ADAPT=1.5
HTM_ALPHA=0.02
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

## 🔐 Seguridad

- ✅ Validación de entrada con Pydantic
- ✅ Límites de tamaño en requests
- ✅ Sanitización de timestamps
- ✅ Manejo seguro de errores
- ⚠️ **Nota:** Este servicio debe estar detrás del API Gateway, no expuesto directamente

## 🚀 Próximas mejoras

- [ ] Caché de modelos entrenados
- [ ] Soporte para múltiples activos simultáneos
- [ ] Métricas de Prometheus
- [ ] Reentrenamiento automático
- [ ] Soporte para datos streaming (MQTT)
- [ ] Tests de integración con InfluxDB

## 📝 Licencia

Proyecto InspectAR - Uso interno

## 👥 Contacto

Para preguntas o soporte, contactar al equipo de desarrollo de InspectAR.