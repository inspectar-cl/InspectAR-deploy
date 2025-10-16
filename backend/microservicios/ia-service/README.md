# IA-Service 🤖

Microservicio de Inteligencia Artificial para detección de anomalías en datos de sensores IoT del proyecto InspectAR.

## 📋 Descripción

Este servicio se encarga de:
- **Monitoreo automático**: Obtiene datos periódicamente del ParserService (cada 5 minutos por defecto)
- **Detección de anomalías con ML**: Procesa 1000 datos por consulta usando el motor de ML Python (HTM)
- **Transformación de datos**: Convierte datos agrupados por sensor a formato temporal para análisis ML
- **Almacenamiento inteligente**: Guarda solo las anomalías detectadas por el modelo en IA-db
- **API REST**: Expone endpoints para almacenar y consultar anomalías

## 🎯 Flujo de Datos

```
ParserService (1000 registros)
        ↓
IA-Service (transformación)
        ↓
ML Engine Python (análisis HTM)
        ↓
IA-Service (guardar resultados)
        ↓
IA-db (almacenamiento)
```

## 📊 Estadísticas de Implementación

- **5 rutas totales** configuradas
- **5 rutas funcionando** (100% operativas)
- **0 rutas pendientes** de implementación
- **Scheduler automático** ejecutándose en background
- **Integración ML Engine** completamente funcional

## 📋 Tabla de Rutas - Vista Rápida

### 🔍 Rutas Básicas

| Método | Endpoint | Descripción | Estado |
|--------|----------|-------------|---------|
| `GET` | `/health` | Health check del servicio | ✅ Funcionando |

### 🤖 Rutas de Anomalías

| Método | Endpoint | Descripción | Estado |
|--------|----------|-------------|---------|
| `POST` | `/anomalies/store` | Guardar nueva anomalía detectada | ✅ Funcionando |
| `GET` | `/status/` | Obtener última anomalía detectada | ✅ Funcionando |
| `GET` | `/anomalies/sensor/:sensor_id` | Obtener anomalías por sensor | ✅ Funcionando |
| `GET` | `/anomalies/activo/:activo_id` | **NUEVO:** Obtener anomalías por activo con paginación | ✅ Funcionando |

### ⚙️ Procesos en Background

| Proceso | Descripción | Frecuencia | Estado |
|---------|-------------|------------|---------|
| 🔄 Scheduler | Obtiene 1000 datos, los envía al ML Engine y guarda anomalías | Cada 5 min (configurable) | ✅ Activo |

## 🏗️ Arquitectura

```
IA-Service (Puerto 8095)
    ↓
    ├─→ PostgreSQL (ia-db:5432) - Almacenamiento de anomalías
    ├─→ IoT-Service (8090) - Obtención de datos de sensores
    └─→ ML Engine Python (8085) - Detección de anomalías con HTM
```

## 🔌 Configuración

### Variables de Entorno

| Variable | Descripción | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Puerto del servicio | `8095` |
| `DB_HOST` | Host de PostgreSQL | `ia-db` |
| `DB_PORT` | Puerto de PostgreSQL | `5432` |
| `DB_USER` | Usuario de BD | `ia_user` |
| `DB_PASSWORD` | Contraseña de BD | `ia_pass` |
| `DB_NAME` | Nombre de BD | `ia_db` |
| `PARSER_SERVICE_URL` | URL del ParserService | `http://iot-service:8090` |
| `ML_ENGINE_URL` | **NUEVO:** URL del ML Engine Python | `http://ml_engine_python:8085` |
| `FETCH_INTERVAL_MINUTES` | Intervalo de consulta (minutos) | `5` |
| `DEFAULT_ACTIVO_ID` | ID del activo a monitorear | `2` |

## 📡 Endpoints

### 1. Health Check

Verifica que el servicio esté funcionando.

```http
GET /health
```

**Respuesta exitosa:**
```json
{
  "status": "healthy",
  "service": "ia-service",
  "message": "Servicio de IA funcionando correctamente"
}
```

**Ejemplo:**
```bash
curl http://localhost:8095/health
```

---

### 2. Guardar Anomalía

Almacena los resultados de predicción de una anomalía.

```http
POST /anomalies/store
Content-Type: application/json
```

**Body:**
```json
{
  "activo_id": 2,
  "sensor_id": "sensor_001",
  "timestamp": "2025-10-14T15:30:00Z",
  "anomaly_score": 0.85,
  "anomaly_likelihood": 0.92,
  "severidad": "alta",
  "descripcion": "Valor anómalo detectado en temperatura",
  "threshold": 75.0,
  "is_anomaly": true
}
```

**Campos:**
- `activo_id` (int, requerido): ID del activo
- `sensor_id` (string, requerido): Identificador del sensor
- `timestamp` (datetime, requerido): Momento de la detección
- `anomaly_score` (float, requerido): Puntuación de anomalía (0-1)
- `anomaly_likelihood` (float, requerido): Probabilidad de anomalía real (0-1)
- `severidad` (string, requerido): `baja`, `media`, `alta`, `critica`
- `descripcion` (string, opcional): Descripción detallada
- `threshold` (float, requerido): Umbral usado para detección
- `is_anomaly` (bool): Confirmación de anomalía

**Respuesta exitosa (201):**
```json
{
  "message": "Anomalía guardada exitosamente",
  "data": {
    "id": 123,
    "activo_id": 2,
    "sensor_id": "sensor_001",
    "timestamp": "2025-10-14T15:30:00Z",
    "anomaly_score": 0.85,
    "anomaly_likelihood": 0.92,
    "severidad": "alta",
    "descripcion": "Valor anómalo detectado en temperatura",
    "threshold": 75.0,
    "is_anomaly": true,
    "created_at": "2025-10-14T15:35:00Z",
    "updated_at": "2025-10-14T15:35:00Z"
  }
}
```

**Ejemplo:**
```bash
curl -X POST http://localhost:8095/anomalies/store \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 2,
    "sensor_id": "sensor_001",
    "timestamp": "2025-10-14T15:30:00Z",
    "anomaly_score": 0.85,
    "anomaly_likelihood": 0.92,
    "severidad": "alta",
    "descripcion": "Temperatura fuera de rango normal",
    "threshold": 75.0,
    "is_anomaly": true
  }'
```

---

### 3. Obtener Última Anomalía

Retorna la anomalía más reciente detectada.

```http
GET /status/
```

**Respuesta exitosa (200):**
```json
{
  "message": "Última anomalía obtenida exitosamente",
  "data": {
    "id": 123,
    "activo_id": 2,
    "sensor_id": "sensor_001",
    "timestamp": "2025-10-14T15:30:00Z",
    "anomaly_score": 0.85,
    "anomaly_likelihood": 0.92,
    "severidad": "alta",
    "descripcion": "Valor anómalo detectado en temperatura",
    "threshold": 75.0,
    "is_anomaly": true,
    "created_at": "2025-10-14T15:35:00Z",
    "updated_at": "2025-10-14T15:35:00Z"
  }
}
```

**Respuesta sin anomalías (404):**
```json
{
  "message": "No se encontraron anomalías"
}
```

**Ejemplo:**
```bash
curl http://localhost:8095/status/
```

---

### 4. Obtener Anomalías por Sensor

Retorna las anomalías de un sensor específico.

```http
GET /anomalies/sensor/:sensor_id?limit=10
```

**Parámetros:**
- `sensor_id` (path, requerido): ID del sensor
- `limit` (query, opcional): Límite de resultados (default: 10)

**Respuesta exitosa (200):**
```json
{
  "message": "Anomalías obtenidas exitosamente",
  "count": 5,
  "sensor_id": "sensor_001",
  "data": [
    {
      "id": 123,
      "activo_id": 2,
      "sensor_id": "sensor_001",
      "timestamp": "2025-10-14T15:30:00Z",
      "anomaly_score": 0.85,
      "anomaly_likelihood": 0.92,
      "severidad": "alta",
      "descripcion": "Valor anómalo detectado",
      "threshold": 75.0,
      "is_anomaly": true,
      "created_at": "2025-10-14T15:35:00Z",
      "updated_at": "2025-10-14T15:35:00Z"
    }
  ]
}
```

**Ejemplo:**
```bash
curl http://localhost:8095/anomalies/sensor/sensor_001?limit=20
```

---

### 5. **NUEVO:** Obtener Anomalías por Activo

Retorna todas las anomalías de un activo específico con paginación.

```http
GET /anomalies/activo/:activo_id?limit=50&offset=0
```

**Parámetros:**
- `activo_id` (path, requerido): ID del activo
- `limit` (query, opcional): Límite de resultados (default: 50)
- `offset` (query, opcional): Offset para paginación (default: 0)

**Respuesta exitosa (200):**
```json
{
  "message": "Anomalías obtenidas exitosamente",
  "count": 25,
  "activo_id": 2,
  "limit": 50,
  "offset": 0,
  "data": [
    {
      "id": 456,
      "activo_id": 2,
      "sensor_id": "combined",
      "timestamp": "2025-10-14T16:45:00Z",
      "anomaly_score": 0.92,
      "anomaly_likelihood": 0.88,
      "severidad": "critica",
      "descripcion": "Patrón anómalo detectado en múltiples sensores",
      "threshold": 0.75,
      "is_anomaly": 1,
      "created_at": "2025-10-14T16:50:00Z",
      "updated_at": "2025-10-14T16:50:00Z"
    },
    {
      "id": 455,
      "activo_id": 2,
      "sensor_id": "combined",
      "timestamp": "2025-10-14T16:40:00Z",
      "anomaly_score": 0.78,
      "anomaly_likelihood": 0.82,
      "severidad": "alta",
      "descripcion": "Desviación significativa en valores de sensores",
      "threshold": 0.75,
      "is_anomaly": 1,
      "created_at": "2025-10-14T16:45:00Z",
      "updated_at": "2025-10-14T16:45:00Z"
    }
  ]
}
```

**Ejemplo:**
```bash
# Obtener las primeras 50 anomalías
curl http://localhost:8095/anomalies/activo/2?limit=50&offset=0 | jq '.'

# Paginación: obtener anomalías 51-100
curl http://localhost:8095/anomalies/activo/2?limit=50&offset=50 | jq '.'

# Obtener solo las últimas 10 anomalías
curl http://localhost:8095/anomalies/activo/2?limit=10 | jq '.'
```

---

## 🤖 Integración con ML Engine

### Proceso de Detección de Anomalías

El servicio utiliza el **ML Engine Python** con algoritmo HTM (Hierarchical Temporal Memory) para detectar anomalías:

1. **Obtención de datos:** El scheduler consulta 1000 registros del ParserService
2. **Transformación:** Los datos se transforman de formato por-sensor a formato temporal
3. **Análisis ML:** Los datos se envían al ML Engine vía POST `/predict_anomaly`
4. **Procesamiento de resultados:** Solo las anomalías detectadas (`is_anomaly=1`) se guardan
5. **Almacenamiento:** Las anomalías se guardan con metadata completa (score, likelihood, severity)

### Formato de Datos para ML Engine

**Entrada (ParserService):**
```json
{
  "activo_id": 2,
  "sensores": [
    {
      "id_sensor": "A_ACR_Mot.PV",
      "datos": [
        {"timestamp": "2025-10-14T10:00:00Z", "valor": 45.2},
        {"timestamp": "2025-10-14T10:01:00Z", "valor": 45.8}
      ]
    },
    {
      "id_sensor": "Temp_Agua",
      "datos": [
        {"timestamp": "2025-10-14T10:00:00Z", "valor": 22.5},
        {"timestamp": "2025-10-14T10:01:00Z", "valor": 22.7}
      ]
    }
  ]
}
```

**Transformado (para ML Engine):**
```json
{
  "pump": "activo_2",
  "records": [
    {
      "timestamp": "2025-10-14T10:00:00Z",
      "A_ACR_Mot.PV": 45.2,
      "A_ACR_Mot.SV": 0.0,
      "A_Desc_Bom.PV": 0.0,
      "A_Desc_Bom.SV": 0.0,
      "Flujo_Caudal": 0.0,
      "Nivel_Presion": 0.0,
      "Potencia_Activa": 0.0,
      "Temp_Agua": 22.5,
      "Temp_Ambiente": 0.0,
      "Vibracion": 0.0
    },
    {
      "timestamp": "2025-10-14T10:01:00Z",
      "A_ACR_Mot.PV": 45.8,
      "Temp_Agua": 22.7,
      ...
    }
  ]
}
```

**Salida (ML Engine):**
```json
{
  "pump": "activo_2",
  "results": [
    {
      "timestamp": "2025-10-14T10:05:00Z",
      "anomaly_score": 0.92,
      "anomaly_likelihood": 0.88,
      "threshold": 0.75,
      "is_anomaly": 1,
      "severity": "critica",
      "description": "Patrón anómalo detectado en múltiples sensores"
    }
  ]
}
```

### Requisitos del ML Engine

- **Mínimo 180 registros**: El modelo HTM requiere al menos 180 registros para análisis
- **10 campos de sensores**: El formato espera valores para 10 tipos de sensores
- **Timeout**: 60 segundos para el análisis completo
- **Endpoint**: `POST http://ml_engine_python:8085/predict_anomaly`

---

## 🔄 Funcionamiento del Scheduler

El servicio incluye un **scheduler automático con integración ML** que:

1. Se ejecuta cada `FETCH_INTERVAL_MINUTES` minutos (default: 5 minutos)
2. Consulta al ParserService: `GET /lectura/{activo_id}/window?page=1&limit=1000`
3. Transforma los datos de formato por-sensor a formato temporal
4. Valida que haya al menos 180 registros (requerido por el modelo HTM)
5. Envía los datos al ML Engine: `POST /predict_anomaly`
6. Procesa las predicciones del modelo
7. Guarda automáticamente solo las anomalías detectadas (`is_anomaly=1`) en IA-db

### Logs del Scheduler

```
🔄 Scheduler iniciado: Obteniendo datos cada 5 minutos para activo_id=2
📡 Solicitando datos del ParserService: http://iot-service:8090/lectura/2/window?page=1&limit=1000
✅ Datos recibidos: 10 sensores, 1000 lecturas totales
📊 Transformados 856 registros agrupados por timestamp
🤖 Enviando 856 registros al ML Engine: http://ml_engine_python:8085/predict_anomaly
✅ ML Engine procesó 856 registros, detectó 12 anomalías
✅ Anomalía guardada: ID=45, Sensor=combined, Score=0.9234, Severidad=critica
✅ Anomalía guardada: ID=46, Sensor=combined, Score=0.8567, Severidad=alta
```

---

## 🚀 Despliegue

### Con Docker Compose

```bash
# Construir y levantar el servicio
cd backend
docker-compose up -d ia-service

# Ver logs
docker logs -f ia-service

# Verificar estado
curl http://localhost:8095/health
```

### Desarrollo Local

```bash
# Instalar dependencias
cd microservicios/ia-service
go mod download

# Ejecutar en desarrollo
go run cmd/main.go

# Compilar
go build -o ia-service cmd/main.go
./ia-service
```

---

## 🧪 Testing

### Test de Health Check
```bash
curl http://localhost:8095/health
```

### Test de Guardar Anomalía
```bash
curl -X POST http://localhost:8095/anomalies/store \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 2,
    "sensor_id": "test_sensor",
    "timestamp": "2025-10-14T10:00:00Z",
    "anomaly_score": 0.95,
    "anomaly_likelihood": 0.88,
    "severidad": "critica",
    "descripcion": "Test de anomalía crítica",
    "threshold": 100.0,
    "is_anomaly": 1
  }'
```

### Test de Consultar Última Anomalía
```bash
curl http://localhost:8095/status/ | jq '.'
```

### Test de Consultar por Sensor
```bash
curl http://localhost:8095/anomalies/sensor/test_sensor?limit=5 | jq '.'
```

### **NUEVO:** Test de Consultar por Activo
```bash
# Obtener anomalías del activo 2 con paginación
curl http://localhost:8095/anomalies/activo/2?limit=10&offset=0 | jq '.'

# Obtener todas las anomalías (máximo 100)
curl http://localhost:8095/anomalies/activo/2?limit=100 | jq '.'
```

---

## 📊 Base de Datos

### Tabla: `anomalias`

Ver esquema completo en `/backend/database/ia/README.md`

**Conexión:**
```bash
# Desde host
psql -h localhost -p 5436 -U ia_user -d ia_db

# Desde Docker
docker exec -it ia-db psql -U ia_user -d ia_db
```

**Consultas útiles:**
```sql
-- Total de anomalías por severidad
SELECT severidad, COUNT(*) 
FROM anomalias 
WHERE is_anomaly = true 
GROUP BY severidad;

-- Últimas 10 anomalías críticas
SELECT * FROM anomalias 
WHERE severidad = 'critica' AND is_anomaly = true 
ORDER BY timestamp DESC 
LIMIT 10;

-- Sensores con más anomalías
SELECT sensor_id, COUNT(*) as total
FROM anomalias
WHERE is_anomaly = true
GROUP BY sensor_id
ORDER BY total DESC;
```

---

## 🛠️ Estructura del Proyecto

```
ia-service/
├── api/
│   └── router.go          # Configuración de rutas
├── cmd/
│   └── main.go            # Punto de entrada
├── config/
│   └── config.go          # Configuración y variables de entorno
├── internal/
│   ├── handler/
│   │   └── anomaly_handler.go    # Handlers HTTP
│   ├── models/
│   │   └── anomaly.go     # Modelos de datos
│   ├── repository/
│   │   └── anomaly_repo.go        # Acceso a BD
│   └── services/
│       └── anomaly_service.go     # Lógica de negocio y scheduler
├── Dockerfile
├── go.mod
└── README.md
```

---

## 📝 Notas Importantes

### Detección de Anomalías con ML Engine

El servicio está **completamente integrado con el ML Engine Python** que utiliza:

- **Algoritmo HTM** (Hierarchical Temporal Memory): Modelo avanzado para detección de anomalías en series temporales
- **Análisis multivariado**: Procesa 10 sensores simultáneamente para detectar patrones anómalos
- **Requisito mínimo**: 180 registros temporales para análisis efectivo
- **Clasificación automática**: El modelo asigna severidad (baja, media, alta, crítica) basado en la desviación

**Características del modelo HTM:**
- Aprende patrones temporales automáticamente
- No requiere entrenamiento previo con datos etiquetados
- Detecta anomalías contextuales (valores normales en contexto anómalo)
- Proporciona scores de confianza (likelihood) para cada predicción

### Transformación de Datos

El servicio realiza una **transformación crítica** de datos:

**Antes (ParserService):** Datos agrupados por sensor
```
sensor1: [timestamp1: valor1, timestamp2: valor2, ...]
sensor2: [timestamp1: valor3, timestamp2: valor4, ...]
```

**Después (ML Engine):** Datos agrupados por timestamp
```
timestamp1: {sensor1: valor1, sensor2: valor3, ...}
timestamp2: {sensor1: valor2, sensor2: valor4, ...}
```

Esta transformación es necesaria porque el modelo HTM analiza **todos los sensores en cada instante** para detectar correlaciones anómalas entre sensores.

### Configuración del Scheduler

Para cambiar la frecuencia de consultas, modifica la variable de entorno:

```yaml
# docker-compose.yml
environment:
  - FETCH_INTERVAL_MINUTES=10  # Cada 10 minutos
```

### Integración con Frontend

El frontend puede consumir estos endpoints para:
- Mostrar dashboard de anomalías en tiempo real
- Alertas push cuando se detectan anomalías críticas
- Gráficos de tendencias de anomalías por sensor
- Historial de eventos anómalos

---

## 🐛 Troubleshooting

### El servicio no se conecta a IA-db

```bash
# Verificar que ia-db esté running
docker ps | grep ia-db

# Verificar logs de ia-db
docker logs ia-db

# Verificar conectividad
docker exec ia-service ping ia-db
```

### El scheduler no consulta datos

```bash
# Verificar logs del scheduler
docker logs -f ia-service | grep "Scheduler"

# Verificar que iot-service esté disponible
curl http://localhost:8090/lectura/2/window?page=1&limit=10
```

### Errores de compilación Go

```bash
# Limpiar cache y recompilar
cd microservicios/ia-service
go clean -cache
go mod tidy
go build -o ia-service cmd/main.go
```

---

## 📞 Soporte

Para reportar problemas o solicitar features, crear un issue en el repositorio del proyecto InspectAR.

---

## 📜 Licencia

Este microservicio es parte del proyecto InspectAR.
