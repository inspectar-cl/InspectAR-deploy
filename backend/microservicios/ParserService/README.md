# Parser Service - Microservicio de Análisis de Datos IoT

Microservicio encargado de la gestión de activos industriales, sensores IoT y análisis de datos en tiempo real.

## 🎉 Estado Actual: COMPLETAMENTE IMPLEMENTADO

**✅ Sistema de Gestión IoT**: **IMPLEMENTADO AL 100%**
- Gestión completa de activos y sensores
- Monitoreo en tiempo real del estado de sensores
- Almacenamiento histórico de lecturas en InfluxDB
- Sistema de notificaciones automáticas
- Detección de desconexión de sensores

**📊 Estadísticas de Implementación:**
- **16 rutas totales** configuradas
- **16 rutas funcionando** (100% operativas)
- **0 rutas pendientes** de implementación
- **Sistema de monitoreo automático** activo (cada 2 minutos)
- **Detección de desconexión** (timeout de 5 minutos)

**🚀 Última Actualización:** 12 de Octubre 2025 - **Información de Sensores en Ruta de Activo Individual** 🆕

## Funcionalidades

### Sistema de Gestión de Activos
- ✅ Crear activos industriales
- ✅ Listar todos los activos del sistema
- ✅ **Obtener activo específico con información completa de sensores**
- ✅ **🆕 Filtrar activos por edificio con información de sensores**
- ✅ Actualizar estado de activos
- ✅ Agregar sensores a activos existentes
- ✅ Consultar activos con estado de sensores (resumen estadístico)

### Sistema de Gestión de Lecturas
- ✅ Registrar lecturas de sensores
- ✅ Obtener datos históricos por activo
- ✅ Obtener última lectura de sensores
- ✅ Almacenamiento persistente en InfluxDB

### Sistema de Monitoreo de Sensores
- ✅ Detección automática de desconexión (>5 min sin datos)
- ✅ Estado persistente en MongoDB
- ✅ Notificaciones automáticas de desconexión
- ✅ Monitoreo en tiempo real (cada 2 minutos)
- ✅ Estadísticas de sensores activos/inactivos
- ✅ Health check del sistema de monitoreo

## 📋 Tabla de Rutas - Vista Rápida

### 🔍 Rutas Básicas

| Método | Endpoint | Descripción | Estado |
|--------|----------|-------------|---------|
| `GET` | `/healthz` | Health check del servicio | ✅ Funcionando |

### 🏗️ Rutas de Activos

| Método | Endpoint | Descripción | Estado |
|--------|----------|-------------|---------|
| `POST` | `/activo` | Crear nuevo activo | ✅ Funcionando |
| `GET` | `/activo` | Listar todos los activos | ✅ Funcionando |
| `GET` | `/activo/:activo_id` | **Obtener activo con info de sensores** | ✅ **ACTUALIZADO** |
| `GET` | `/activo/:activo_id/sensores/estado` | Obtener activo con estado de sensores | ✅ Funcionando |
| `POST` | `/activo/:activo_id/sensores` | Agregar sensor a activo existente | ✅ Funcionando |
| `GET` | `/activo/edificio/:edificio_id` | 🆕 **Filtrar activos por edificio** | ✅ **NUEVO** |
| `PUT` | `/activo/:activo_id/estado` | Actualizar estado del activo | ✅ Funcionando |

### 📊 Rutas de Lecturas de Sensores

| Método | Endpoint | Descripción | Estado |
|--------|----------|-------------|---------|
| `POST` | `/lectura` | Registrar nueva lectura de sensor | ✅ Funcionando |
| `GET` | `/lectura/:activo_id/datos` | Obtener datos históricos de activo | ✅ Funcionando |
| `GET` | `/lectura/:activo_id/datos/ultimo` | Obtener última lectura de sensores | ✅ Funcionando |

### 🔔 Rutas de Monitoreo de Sensores

| Método | Endpoint | Descripción | Estado |
|--------|----------|-------------|---------|
| `GET` | `/api/sensors/status` | Estado de todos los sensores | ✅ Funcionando |
| `GET` | `/api/sensors/status/:sensor_id` | Estado de un sensor específico | ✅ Funcionando |
| `GET` | `/api/sensors/stats` | Estadísticas generales (activos/inactivos) | ✅ Funcionando |
| `POST` | `/api/sensors/check-disconnected` | Verificación manual de desconexiones | ✅ Funcionando |
| `GET` | `/api/sensors/health` | Health check del sistema de monitoreo | ✅ Funcionando |

### 🎯 Resumen de Estado

- **✅ Funcionando**: 16 rutas operativas (100% IMPLEMENTADAS)
- **🔧 No implementado**: 0 rutas pendientes
- **Total**: 16 rutas configuradas

### ⚡ Tests Rápidos

```bash
# Verificar servicio
curl http://localhost:8090/healthz

# Crear activo
curl -X POST http://localhost:8090/activo \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 1,
    "nombre": "Caldera Principal",
    "estado": "operativo",
    "id_edificio": "1",
    "sensores": [
      {"sensor_id": "TEMP_001", "tipo": "temperatura", "unidad": "°C"}
    ]
  }'

# Listar todos los activos
curl http://localhost:8090/activo

# Obtener activo específico
curl http://localhost:8090/activo/1

# 🆕 NUEVA RUTA: Obtener activos por edificio (con info de sensores)
curl http://localhost:8090/activo/edificio/1

# Obtener activo con estado de sensores
curl http://localhost:8090/activo/1/sensores/estado

# Agregar sensor a activo
curl -X POST http://localhost:8090/activo/1/sensores \
  -H "Content-Type: application/json" \
  -d '{"sensor_id": "PRESS_001", "tipo": "presion", "unidad": "bar"}'

# Registrar lectura de sensor
curl -X POST http://localhost:8090/lectura \
  -H "Content-Type: application/json" \
  -d '{
    "sensor_id": "TEMP_001",
    "valor": 75.5,
    "timestamp": "2025-10-10T10:30:00Z"
  }'

# Obtener datos históricos
curl http://localhost:8090/lectura/1/datos

# Obtener última lectura
curl http://localhost:8090/lectura/1/datos/ultimo

# Actualizar estado de activo
curl -X PUT http://localhost:8090/activo/1/estado \
  -H "Content-Type: application/json" \
  -d '{"estado": "mantenimiento"}'

# Estadísticas de sensores
curl http://localhost:8090/api/sensors/stats

# Verificar desconexiones manualmente
curl -X POST http://localhost:8090/api/sensors/check-disconnected

# Health check de monitoreo
curl http://localhost:8090/api/sensors/health
```

## API Endpoints

### 🏥 Health Check

#### Verificar estado del servicio
```bash
curl -X GET http://localhost:8090/healthz
```
**Respuesta:**
```json
{
  "ok": true
}
```

---

### 🏗️ Activos

#### Crear activo
```bash
curl -X POST http://localhost:8090/activo \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 1,
    "nombre": "Caldera Principal",
    "estado": "operativo",
    "id_edificio": "1",
    "sensores": [
      {
        "sensor_id": "TEMP_001",
        "tipo": "temperatura",
        "unidad": "°C"
      },
      {
        "sensor_id": "PRESS_001",
        "tipo": "presion",
        "unidad": "bar"
      }
    ]
  }'
```
**Respuesta:**
```json
{
  "activo_id": "670586c9a3e45c001e8b4567"
}
```

#### Listar todos los activos
```bash
curl -X GET http://localhost:8090/activo
```
**Respuesta:**
```json
[
  {
    "id": "670586c9a3e45c001e8b4567",
    "activo_id": 1,
    "nombre": "Caldera Principal",
    "estado": "operativo",
    "id_edificio": "1",
    "sensores": [
      {
        "sensor_id": "TEMP_001",
        "tipo": "temperatura",
        "unidad": "°C"
      },
      {
        "sensor_id": "PRESS_001",
        "tipo": "presion",
        "unidad": "bar"
      }
    ]
  }
]
```

#### Obtener activo específico con información de sensores
```bash
curl -X GET http://localhost:8090/activo/1
```
**Respuesta:**
```json
{
  "id": "670586c9a3e45c001e8b4567",
  "activo_id": 1,
  "estado": "operativo",
  "id_edificio": 1,
  "total_sensores": 2,
  "sensores": [
    {
      "sensor_id": "TEMP_001",
      "tipo": "temperatura",
      "unidad": "°C",
      "estado": "connected",
      "is_active": true,
      "last_seen": "2025-10-11T12:00:00Z",
      "first_seen": "2025-10-01T08:00:00Z",
      "total_reports": 1458,
      "created_at": "2025-10-01T08:00:00Z",
      "updated_at": "2025-10-11T12:00:00Z"
    },
    {
      "sensor_id": "PRESS_001",
      "tipo": "presion",
      "unidad": "bar",
      "estado": "disconnected",
      "is_active": false,
      "last_seen": "2025-10-10T09:45:00Z",
      "first_seen": "2025-10-01T08:00:00Z",
      "total_reports": 892,
      "created_at": "2025-10-01T08:00:00Z",
      "updated_at": "2025-10-10T09:45:00Z"
    }
  ]
}
```

**Características:**
- Incluye información completa del activo
- Array de sensores con estado detallado de cada uno
- Campos de estado: `estado`, `is_active`, `last_seen`, `first_seen`
- Métricas: `total_reports`, `created_at`, `updated_at`
- Campo `total_sensores` con cantidad de sensores del activo

**Estados de Sensores:**
- **`connected`**: Sensor activo enviando datos
- **`disconnected`**: Sensor que envió datos pero está inactivo (>5 min sin datos)
- **`never_connected`**: Sensor registrado pero nunca ha enviado datos

#### 🆕 **NUEVA RUTA: Obtener activos por edificio con información de sensores**
```bash
curl -X GET http://localhost:8090/activo/edificio/1
```
**Respuesta:**
```json
{
  "edificio_id": 1,
  "activos": [
    {
      "id": "670586c9a3e45c001e8b4567",
      "activo_id": 1,
      "estado": "operativo",
      "id_edificio": 1,
      "total_sensores": 2,
      "sensores": [
        {
          "sensor_id": "TEMP_001",
          "tipo": "temperatura",
          "unidad": "°C",
          "estado": "connected",
          "is_active": true,
          "last_seen": "2025-10-11T12:00:00Z",
          "first_seen": "2025-10-01T08:00:00Z",
          "total_reports": 1458,
          "created_at": "2025-10-01T08:00:00Z",
          "updated_at": "2025-10-11T12:00:00Z"
        },
        {
          "sensor_id": "PRESS_001",
          "tipo": "presion",
          "unidad": "bar",
          "estado": "disconnected",
          "is_active": false,
          "last_seen": "2025-10-10T09:45:00Z",
          "total_reports": 892
        }
      ]
    },
    {
      "id": "670586d2a3e45c001e8b4568",
      "activo_id": 2,
      "estado": "operativo",
      "id_edificio": 1,
      "total_sensores": 2,
      "sensores": [
        {
          "sensor_id": "FLOW_001",
          "tipo": "flujo",
          "unidad": "L/min",
          "estado": "connected",
          "is_active": true,
          "last_seen": "2025-10-11T12:00:00Z",
          "total_reports": 523
        },
        {
          "sensor_id": "VIBR_001",
          "tipo": "vibracion",
          "unidad": "Hz",
          "estado": "never_connected",
          "is_active": false,
          "total_reports": 0
        }
      ]
    }
  ],
  "total": 2
}
```

**Características:**
- Incluye información completa de cada activo del edificio
- Cada activo incluye el array de sensores con su estado
- Cada sensor muestra:
  - Datos básicos: `sensor_id`, `tipo`, `unidad`
  - Estado de conexión: `estado`, `is_active`
  - Métricas: `last_seen`, `first_seen`, `total_reports`
  - Timestamps: `created_at`, `updated_at`
- Campo `total_sensores` por activo
- Campo `total` con cantidad de activos en el edificio

**Estados de Sensores:**
- **`connected`**: Sensor activo enviando datos
- **`disconnected`**: Sensor que envió datos pero está inactivo (>5 min sin datos)
- **`never_connected`**: Sensor registrado pero nunca ha enviado datos

#### Obtener activo con estado de sensores
```bash
curl -X GET http://localhost:8090/activo/1/sensores/estado
```
**Respuesta:**
```json
{
  "id": "670586c9a3e45c001e8b4567",
  "activo_id": 1,
  "nombre": "Caldera Principal",
  "estado": "operativo",
  "id_edificio": "1",
  "total_sensores": 3,
  "sensores": [
    {
      "sensor_id": "TEMP_001",
      "tipo": "temperatura",
      "unidad": "°C",
      "estado": "connected",
      "is_active": true,
      "last_seen": "2025-10-10T10:30:00Z",
      "first_seen": "2025-10-01T08:00:00Z",
      "total_reports": 1458,
      "created_at": "2025-10-01T08:00:00Z",
      "updated_at": "2025-10-10T10:30:00Z"
    },
    {
      "sensor_id": "PRESS_001",
      "tipo": "presion",
      "unidad": "bar",
      "estado": "disconnected",
      "is_active": false,
      "last_seen": "2025-10-10T09:45:00Z",
      "total_reports": 892
    },
    {
      "sensor_id": "VIBR_001",
      "tipo": "vibracion",
      "unidad": "Hz",
      "estado": "never_connected",
      "is_active": false,
      "total_reports": 0
    }
  ],
  "resumen": {
    "sensores_activos": 1,
    "sensores_desconectados": 1,
    "sensores_nunca_conectados": 1
  }
}
```

**Estados de Sensores:**
- **`connected`**: Sensor activo enviando datos
- **`disconnected`**: Sensor que envió datos pero está inactivo (>5 min sin datos)
- **`never_connected`**: Sensor registrado pero nunca ha enviado datos

#### Agregar sensor a activo existente
```bash
curl -X POST http://localhost:8090/activo/1/sensores \
  -H "Content-Type: application/json" \
  -d '{
    "sensor_id": "TEMP_002",
    "tipo": "temperatura",
    "unidad": "°C"
  }'
```
**Respuesta:**
```json
{
  "mensaje": "Sensor agregado correctamente al activo",
  "activo_id": 1,
  "sensor": {
    "sensor_id": "TEMP_002",
    "tipo": "temperatura",
    "unidad": "°C"
  }
}
```

**Validaciones:**
- ✅ El activo debe existir (404 si no existe)
- ✅ Campos requeridos: `sensor_id`, `tipo`, `unidad` (400 si faltan)
- ✅ Formato JSON válido (400 si es inválido)

#### Actualizar estado de activo
```bash
curl -X PUT http://localhost:8090/activo/1/estado \
  -H "Content-Type: application/json" \
  -d '{
    "estado": "mantenimiento"
  }'
```
**Respuesta:**
```json
{
  "message": "Estado actualizado correctamente"
}
```

---

### 📊 Lecturas de Sensores

#### Registrar nueva lectura
```bash
curl -X POST http://localhost:8090/lectura \
  -H "Content-Type: application/json" \
  -d '{
    "sensor_id": "TEMP_001",
    "valor": 75.5,
    "timestamp": "2025-10-10T10:30:00Z"
  }'
```
**Respuesta:**
```json
{
  "mensaje": "Lectura registrada en InfluxDB"
}
```

**Comportamiento:**
- Guarda la lectura en InfluxDB
- Actualiza automáticamente el estado del sensor en MongoDB
- Marca el sensor como activo (`connected`)

#### Obtener datos históricos de activo
```bash
curl -X GET http://localhost:8090/lectura/1/datos
```
**Respuesta:**
```json
{
  "id": "670586c9a3e45c001e8b4567",
  "activo_id": 1,
  "nombre": "Caldera Principal",
  "estado": "operativo",
  "sensores": [
    {
      "sensor_id": "TEMP_001",
      "datos": [
        {
          "timestamp": "2025-10-10T10:30:00Z",
          "valor": 75.5
        },
        {
          "timestamp": "2025-10-10T10:25:00Z",
          "valor": 74.8
        }
      ]
    }
  ]
}
```

**Nota:** Por defecto retorna datos de los últimos 30 minutos.

#### Obtener última lectura de sensores
```bash
curl -X GET http://localhost:8090/lectura/1/datos/ultimo
```
**Respuesta:**
```json
{
  "TEMP_001": {
    "timestamp": "2025-10-10T10:30:00Z",
    "valor": 75.5
  },
  "PRESS_001": {
    "timestamp": "2025-10-10T10:30:00Z",
    "valor": 8.5
  }
}
```

---

### 🔔 Monitoreo de Sensores

#### Obtener estado de todos los sensores
```bash
curl -X GET http://localhost:8090/api/sensors/status
```
**Respuesta:**
```json
[
  {
    "sensor_id": "TEMP_001",
    "is_active": true,
    "last_seen": "2025-10-10T10:30:00Z",
    "first_seen": "2025-10-01T08:00:00Z",
    "total_reports": 1458,
    "created_at": "2025-10-01T08:00:00Z",
    "updated_at": "2025-10-10T10:30:00Z"
  },
  {
    "sensor_id": "PRESS_001",
    "is_active": false,
    "last_seen": "2025-10-10T09:45:00Z",
    "first_seen": "2025-10-01T08:00:00Z",
    "total_reports": 892,
    "created_at": "2025-10-01T08:00:00Z",
    "updated_at": "2025-10-10T10:30:00Z"
  }
]
```

#### Obtener estado de un sensor específico
```bash
curl -X GET http://localhost:8090/api/sensors/status/TEMP_001
```
**Respuesta:**
```json
{
  "sensor_id": "TEMP_001",
  "is_active": true,
  "last_seen": "2025-10-10T10:30:00Z",
  "first_seen": "2025-10-01T08:00:00Z",
  "total_reports": 1458,
  "created_at": "2025-10-01T08:00:00Z",
  "updated_at": "2025-10-10T10:30:00Z"
}
```

#### Obtener estadísticas generales
```bash
curl -X GET http://localhost:8090/api/sensors/stats
```
**Respuesta:**
```json
{
  "total_sensors": 15,
  "active_sensors": 12,
  "inactive_sensors": 3,
  "timestamp": "2025-10-10T10:30:00Z"
}
```

#### Verificar desconexiones manualmente
```bash
curl -X POST http://localhost:8090/api/sensors/check-disconnected
```
**Respuesta:**
```json
{
  "message": "Verificación de sensores completada",
  "disconnected_count": 2,
  "disconnected_sensors": ["PRESS_001", "VIBR_002"]
}
```

**Comportamiento:**
- Busca sensores sin datos por más de 5 minutos
- Los marca como inactivos en MongoDB
- Envía notificaciones automáticas (opcional)

#### Health check del sistema de monitoreo
```bash
curl -X GET http://localhost:8090/api/sensors/health
```
**Respuesta:**
```json
{
  "status": "healthy",
  "monitoring_active": true,
  "last_check": "2025-10-10T10:28:00Z",
  "check_interval": "2 minutes"
}
```

---

## 📋 Modelos de Datos

### Activo (MongoDB)
```go
type Activo struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    ActivoID    int               `json:"activo_id" bson:"activo_id"`
    Nombre      string            `json:"nombre" bson:"nombre"`
    Estado      string            `json:"estado" bson:"estado"`
    Id_edificio string            `json:"id_edificio" bson:"id_edificio"`
    Sensores    []Sensor          `json:"sensores" bson:"sensores"`
}
```

### Sensor
```go
type Sensor struct {
    SensorID string `json:"sensor_id" bson:"sensor_id"`
    Tipo     string `json:"tipo" bson:"tipo"`
    Unidad   string `json:"unidad" bson:"unidad"`
}
```

### SensorStatus (MongoDB)
```go
type SensorStatus struct {
    SensorID     string    `json:"sensor_id" bson:"sensor_id"`
    IsActive     bool      `json:"is_active" bson:"is_active"`
    LastSeen     time.Time `json:"last_seen" bson:"last_seen"`
    FirstSeen    time.Time `json:"first_seen" bson:"first_seen"`
    TotalReports int64     `json:"total_reports" bson:"total_reports"`
    CreatedAt    time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}
```

### LecturaSensorRequest
```go
type LecturaSensorRequest struct {
    SensorID  string  `json:"sensor_id"`
    Valor     float64 `json:"valor"`
    Timestamp string  `json:"timestamp"`
}
```

---

## ⚙️ Configuración

### Variables de Entorno
```yaml
# MongoDB
MONGO_URI: mongodb://parser_user:parser_pass@parser-db:27017/parserdb?authSource=admin

# InfluxDB
INFLUX_URL: http://influxdb:8086
INFLUX_TOKEN: your-influx-token
INFLUX_ORG: inspectAR
INFLUX_BUCKET: sensor_data

# Servicio de Notificaciones (opcional)
NOTIFICATION_URL: http://notification-service:8091/notification

# Configuración de Monitoreo
SENSOR_TIMEOUT_MINUTES: 5  # Tiempo para considerar desconectado
CHECK_INTERVAL_MINUTES: 2   # Intervalo de verificación automática
```

### Archivo config.yaml
```yaml
sensor:
  timeout_minutes: 5  # Tiempo para considerar desconectado

notification:
  url: "http://notification-service:8091/notification"
  enabled: true

activos:
  primary_collection: "activos"
  sensors_collection: "sensores"
  source_names:
    - "activos"
```

---

## 🔄 Flujo de Funcionamiento

### 1. Registro de Lecturas
```
POST /lectura
    ↓
Guardar en InfluxDB
    ↓
Actualizar SensorStatus en MongoDB (marca como activo)
    ↓
Respuesta 200 OK
```

### 2. Monitoreo Automático (cada 2 minutos)
```
Sistema de Monitoreo
    ↓
Buscar sensores sin datos por >5 minutos
    ↓
Marcar como inactivos en MongoDB
    ↓
Enviar notificaciones (opcional)
    ↓
Log en consola
```

### 3. Consulta de Estado
```
GET /activo/:activo_id/sensores/estado
    ↓
Obtener activo de MongoDB
    ↓
Obtener estado de cada sensor
    ↓
Calcular resumen (activos/desconectados/nunca_conectados)
    ↓
Respuesta con datos completos
```

### 4. Filtrado por Edificio 🆕
```
GET /activo/edificio/:edificio_id
    ↓
Buscar activos con id_edificio en MongoDB
    ↓
Enriquecer con sensores
    ↓
Retornar lista completa con total
```

---

## Estructura de Base de Datos

### MongoDB Collections

#### `activos`
- **Propósito**: Almacenar activos industriales y sus sensores
- **Campos**: id, activo_id, nombre, estado, id_edificio, sensores[]
- **Índices**: activo_id (unique), id_edificio

#### `sensores` (opcional)
- **Propósito**: Almacenamiento alternativo de sensores (normalizado)
- **Campos**: sensor_id, activo_id, tipo, unidad
- **Índices**: sensor_id (unique), activo_id

#### `sensor_status`
- **Propósito**: Estado de monitoreo de sensores
- **Campos**: sensor_id, is_active, last_seen, first_seen, total_reports, created_at, updated_at
- **Índices**: sensor_id (unique), is_active, last_seen

### InfluxDB Measurements

#### `sensor_data`
- **Tags**: sensor_id, activo_id, tipo
- **Fields**: valor (float)
- **Timestamp**: Marca de tiempo de la lectura
- **Retention**: Configurable (default: 30 días)

---

## 🧪 Testing

### Test Completo de Monitoreo
```bash
cd tests
./test_monitoring.sh
```

Este test verifica:
- Creación de activos con sensores
- Registro de lecturas
- Detección de sensores conectados
- Detección de desconexión (después de timeout)
- Notificaciones de estado

### Test de Consulta de Activo con Estado de Sensores
```bash
cd tests
./test_activo_sensores_estado.sh
```

Este test verifica:
- Creación de activo con sensores
- Consulta de sensores sin datos (never_connected)
- Envío de datos y detección de sensores conectados
- Resumen de estadísticas por activo
- Validación de estructura de respuesta

### Test de Agregar Sensores a Activos
```bash
cd tests
./test_add_sensor_to_activo.sh
```

Este test verifica:
- Agregar sensores a activos existentes
- Validaciones de campos requeridos
- Manejo de errores (activo inexistente, datos inválidos)
- Verificación de la integridad de los datos agregados

### 🆕 Test de Filtrado por Edificio
```bash
# Crear activos en diferentes edificios
curl -X POST http://localhost:8090/activo -H "Content-Type: application/json" \
  -d '{"activo_id": 1, "nombre": "Caldera A", "estado": "operativo", "id_edificio": "1", "sensores": []}'

curl -X POST http://localhost:8090/activo -H "Content-Type: application/json" \
  -d '{"activo_id": 2, "nombre": "Bomba A", "estado": "operativo", "id_edificio": "1", "sensores": []}'

curl -X POST http://localhost:8090/activo -H "Content-Type: application/json" \
  -d '{"activo_id": 3, "nombre": "Caldera B", "estado": "operativo", "id_edificio": "2", "sensores": []}'

# Verificar filtrado por edificio 1 (debe retornar 2 activos)
curl http://localhost:8090/activo/edificio/1

# Verificar filtrado por edificio 2 (debe retornar 1 activo)
curl http://localhost:8090/activo/edificio/2
```

---

## 📦 Dependencias Principales

- **MongoDB**: Almacenamiento de activos y estado de sensores
- **InfluxDB**: Series temporales de lecturas
- **Gin**: Framework web HTTP
- **Viper**: Gestión de configuración
- **InfluxDB Client**: Interacción con base de datos de series temporales
- **MongoDB Driver**: Driver oficial de MongoDB para Go

---

## 🚨 Sistema de Notificaciones

### Comportamiento al Detectar Desconexión

Cuando un sensor se desconecta (>5 minutos sin datos):

1. **Console Log**: Mensaje detallado en logs del servidor
   ```
   SENSOR DESCONECTADO: TEMP_001
   Última vez visto: 2025-10-10T10:25:00Z
   Total reportes: 1458
   ```

2. **HTTP Request** (opcional): POST al servicio de notificaciones
   ```json
   {
     "tipo": "sensor_desconectado",
     "sensor_id": "TEMP_001",
     "last_seen": "2025-10-10T10:25:00Z",
     "mensaje": "El sensor TEMP_001 no ha enviado datos en 5 minutos"
   }
   ```

3. **Estado Persistente**: Actualización en MongoDB
   - `is_active` → `false`
   - `updated_at` → timestamp actual

---

## 🔍 Ejemplos de Uso

### Flujo típico de operación

```bash
# 1. Crear un activo con sensores
curl -X POST http://localhost:8090/activo \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 1,
    "nombre": "Caldera Principal",
    "estado": "operativo",
    "id_edificio": "1",
    "sensores": [
      {"sensor_id": "TEMP_001", "tipo": "temperatura", "unidad": "°C"},
      {"sensor_id": "PRESS_001", "tipo": "presion", "unidad": "bar"}
    ]
  }'

# 2. Registrar lecturas de sensores
curl -X POST http://localhost:8090/lectura \
  -H "Content-Type: application/json" \
  -d '{"sensor_id": "TEMP_001", "valor": 75.5, "timestamp": "2025-10-10T10:30:00Z"}'

curl -X POST http://localhost:8090/lectura \
  -H "Content-Type: application/json" \
  -d '{"sensor_id": "PRESS_001", "valor": 8.5, "timestamp": "2025-10-10T10:30:00Z"}'

# 3. Consultar estado del activo con sus sensores
curl http://localhost:8090/activo/1/sensores/estado

# 4. Agregar un nuevo sensor al activo existente
curl -X POST http://localhost:8090/activo/1/sensores \
  -H "Content-Type: application/json" \
  -d '{"sensor_id": "VIBR_001", "tipo": "vibracion", "unidad": "Hz"}'

# 5. 🆕 Obtener todos los activos del edificio
curl http://localhost:8090/activo/edificio/1

# 6. Obtener datos históricos
curl http://localhost:8090/lectura/1/datos

# 7. Obtener última lectura
curl http://localhost:8090/lectura/1/datos/ultimo

# 8. Ver estadísticas de sensores
curl http://localhost:8090/api/sensors/stats

# 9. Verificar desconexiones manualmente
curl -X POST http://localhost:8090/api/sensors/check-disconnected

# 10. Actualizar estado del activo
curl -X PUT http://localhost:8090/activo/1/estado \
  -H "Content-Type: application/json" \
  -d '{"estado": "mantenimiento"}'
```

### Escenario: Monitoreo de múltiples edificios con información de sensores 🆕

```bash
# Obtener todos los activos del sistema
curl http://localhost:8090/activo

# Filtrar activos por edificio específico (incluye estado de sensores)
curl http://localhost:8090/activo/edificio/1
curl http://localhost:8090/activo/edificio/2
curl http://localhost:8090/activo/edificio/3

# Obtener estado detallado de sensores de cada edificio
for edificio in 1 2 3; do
  echo "=== Edificio $edificio ==="
  curl http://localhost:8090/activo/edificio/$edificio | jq '.activos[] | {activo_id, estado, total_sensores, sensores: .sensores | map({sensor_id, tipo, estado})}'
done
```

---

## Ejecución

### Con Docker Compose
```bash
# Iniciar todos los servicios
cd backend
docker-compose up parser-service parser-db influxdb

# Ver logs
docker-compose logs -f parser-service

# Detener servicios
docker-compose down
```

### Desarrollo Local
```bash
# Instalar dependencias
go mod download

# Ejecutar el servicio
go run cmd/main.go

# El servicio estará disponible en http://localhost:8090
```

---

## Puerto por Defecto

El servicio Parser Service está disponible en el puerto **8090**.

---

## Integración con Otros Servicios

### Microservicio de Gestión
El Parser Service se integra con el microservicio de Gestión para:
- Sincronizar información de activos
- Compartir datos de sensores para reportes PDF
- Proporcionar métricas calculadas para reportes

**Endpoint de integración:**
```
GET http://localhost:8090/activo/edificio/:edificio_id
```

### Servicio de Notificaciones
Envía notificaciones automáticas cuando:
- Un sensor se desconecta (>5 minutos sin datos)
- Un activo cambia de estado
- Se detectan anomalías en lecturas

---

## 🎯 Próximas Funcionalidades

- [ ] Alertas configurables por umbrales de sensor
- [ ] Dashboard en tiempo real con WebSockets
- [ ] Exportación de datos históricos (CSV, Excel)
- [ ] Análisis predictivo con Machine Learning
- [ ] Integración con servicios de BI externos
- [ ] API GraphQL alternativa

---

## 📝 Notas de Desarrollo

### Arquitectura
- **Patrón Repository**: Separación de lógica de acceso a datos
- **Clean Architecture**: Capas bien definidas (handler → service → repository)
- **Múltiples fuentes de datos**: Soporte para merge de colecciones MongoDB
- **Monitoreo automático**: Goroutine en background para verificación periódica

### Consideraciones de Rendimiento
- Uso de índices en MongoDB para consultas rápidas
- Paginación implícita en InfluxDB (últimos 30 minutos por defecto)
- Cache de estados de sensores (opcional, configurable)
- Queries optimizadas con proyecciones específicas

### Seguridad
- Validación de entrada en todos los endpoints
- Sanitización de parámetros de query
- Autenticación y autorización (pendiente integración con API Gateway)
- Logging de todas las operaciones críticas

---

**Desarrollado por el equipo de InspectAR** 🚀
