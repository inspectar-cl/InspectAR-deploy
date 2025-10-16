# IA-Service 🤖

Microservicio de Inteligencia Artificial para detección de anomalías en datos de sensores IoT del proyecto InspectAR.

## 📋 Descripción

Este servicio se encarga de:
- **Monitoreo automático**: Obtiene datos periódicamente del ParserService (cada 5 minutos por defecto)
- **Detección de anomalías**: Procesa 1000 datos por consulta buscando patrones anómalos
- **Almacenamiento**: Guarda resultados de predicción en la base de datos IA-db
- **API REST**: Expone endpoints para almacenar y consultar anomalías

## 📊 Estadísticas de Implementación

- **4 rutas totales** configuradas
- **4 rutas funcionando** (100% operativas)
- **0 rutas pendientes** de implementación
- **Scheduler automático** ejecutándose en background

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

### ⚙️ Procesos en Background

| Proceso | Descripción | Frecuencia | Estado |
|---------|-------------|------------|---------|
| 🔄 Scheduler | Obtiene 1000 datos del ParserService y procesa anomalías | Cada 5 min (configurable) | ✅ Activo |

## 🏗️ Arquitectura

```
IA-Service (Puerto 8095)
    ↓
    ├─→ PostgreSQL (ia-db:5432) - Almacenamiento de anomalías
    └─→ IoT-Service (8090) - Obtención de datos de sensores
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

## 🔄 Funcionamiento del Scheduler

El servicio incluye un **scheduler automático** que:

1. Se ejecuta cada `FETCH_INTERVAL_MINUTES` minutos (default: 5 minutos)
2. Consulta al ParserService: `GET /lectura/{activo_id}/window?page=1&limit=1000`
3. Procesa los 1000 datos obtenidos
4. Detecta anomalías usando el algoritmo configurado
5. Guarda automáticamente las anomalías detectadas en IA-db

### Logs del Scheduler

```
🔄 Scheduler iniciado: Obteniendo datos cada 5 minutos para activo_id=2
📡 Solicitando datos del ParserService: http://iot-service:8090/lectura/2/window?page=1&limit=1000
✅ Datos recibidos: 10 sensores, 1000 lecturas totales
✅ Anomalía guardada: ID=45, Sensor=sensor_003, Score=0.7854, Severidad=alta
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
    "is_anomaly": true
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

### Detección de Anomalías

La implementación actual incluye un **algoritmo simple de detección** basado en umbrales. Para producción, se recomienda:

1. Integrar modelos de ML entrenados (TensorFlow, PyTorch, scikit-learn)
2. Usar algoritmos avanzados como:
   - Isolation Forest
   - LSTM para series temporales
   - Autoencoders
   - Prophet de Facebook

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
