# Parser Service - Microservicio de Análisis de Datos IoT

## 📊 Funcionalidades Principales

### 1. **Gestión de Activos y Sensores**
- Registro y gestión de activos industriales
- Almacenamiento de lecturas de sensores en InfluxDB
- Consulta de datos históricos y en tiempo real

### 2. **🆕 Sistema de Monitoreo de Estado de Sensores**
- **Detección automática de desconexión**: Si un sensor no envía datos por 5 minutos, se marca como desconectado
- **Estado persistente en MongoDB**: Cada sensor tiene un registro de estado (activo/inactivo)
- **Notificaciones automáticas**: Se envían notificaciones cuando un sensor se desconecta
- **Monitoreo en tiempo real**: Verificación cada 2 minutos de sensores desconectados

## � Respuesta de Consulta de Activo con Estado de Sensores

El endpoint `GET /activo/:activo_id/sensores/estado` retorna:

```json
{
  "id": "...",
  "activo_id": "CALDERA_001",
  "nombre": "Caldera Principal",
  "ubicacion": "Planta Baja",
  "estado": "activo",
  "id_edificio": "edificio_001",
  "total_sensores": 3,
  "sensores": [
    {
      "sensor_id": "TEMP_001",
      "tipo": "temperatura", 
      "unidad": "°C",
      "estado": "connected",           // connected | disconnected | never_connected
      "is_active": true,
      "last_seen": "2025-09-03T10:30:00Z",
      "first_seen": "2025-09-01T08:00:00Z",
      "total_reports": 1458,
      "created_at": "2025-09-01T08:00:00Z",
      "updated_at": "2025-09-03T10:30:00Z"
    },
    {
      "sensor_id": "PRESS_001",
      "tipo": "presion",
      "unidad": "bar", 
      "estado": "disconnected",
      "is_active": false,
      "last_seen": "2025-09-03T09:45:00Z",
      "total_reports": 892
    },
    {
      "sensor_id": "VIBR_001",
      "tipo": "vibracion",
      "unidad": "Hz",
      "estado": "never_connected",     // Nunca ha enviado datos
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

### Estados de Sensores:
- **`connected`**: Sensor activo enviando datos
- **`disconnected`**: Sensor que envió datos pero está inactivo (>5 min sin datos)
- **`never_connected`**: Sensor registrado pero nunca ha enviado datos

## �🚀 Endpoints Disponibles

### Gestión de Activos
```
POST /activo - Crear nuevo activo
GET /activo - Listar todos los activos
GET /activo/:activo_id - Obtener activo específico
GET /activo/:activo_id/sensores/estado - 🆕 Obtener activo con estado de sensores
PUT /activo/:activo_id/estado - Actualizar estado del activo
```

### Gestión de Lecturas
```
POST /lectura - Registrar nueva lectura de sensor
GET /lectura/:activo_id/datos - Obtener datos históricos
GET /lectura/:activo_id/datos/ultimo - Obtener última lectura
```

### 🆕 Monitoreo de Sensores
```
GET /api/sensors/status - Estado de todos los sensores
GET /api/sensors/status/:sensor_id - Estado de un sensor específico
GET /api/sensors/stats - Estadísticas generales (activos/inactivos)
POST /api/sensors/check-disconnected - Verificación manual de desconexiones
GET /api/sensors/health - Health check del sistema de monitoreo
```

## 📋 Modelos de Datos

### SensorStatus (MongoDB)
```go
type SensorStatus struct {
    SensorID     string    `json:"sensor_id"`
    IsActive     bool      `json:"is_active"`
    LastSeen     time.Time `json:"last_seen"`
    FirstSeen    time.Time `json:"first_seen"`
    TotalReports int64     `json:"total_reports"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

## ⚙️ Configuración

```yaml
# config/config.yaml
sensor:
  timeout_minutes: 5  # Tiempo para considerar desconectado

notification:
  url: "http://notification-service:8091/notification"  # Servicio de notificaciones
```

## 🔄 Flujo de Funcionamiento

1. **Llegada de Datos**: Cuando llega una lectura (`POST /lectura`):
   - Se guarda en InfluxDB
   - Se actualiza el estado del sensor en MongoDB (marca como activo)

2. **Monitoreo Automático**: Cada 2 minutos:
   - Se buscan sensores que no han enviado datos por 5+ minutos
   - Se marcan como inactivos
   - Se envía notificación (console.log + HTTP opcional)

3. **Consulta de Estado**: 
   - Los endpoints permiten consultar el estado en tiempo real
   - Se pueden obtener estadísticas de sensores activos/inactivos

## 🧪 Testing

### Test Completo de Monitoreo
```bash
cd tests
./test_monitoring.sh
```

### 🆕 Test de Consulta de Activo con Estado de Sensores
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

## 📦 Dependencias Principales

- **MongoDB**: Estado persistente de sensores
- **InfluxDB**: Datos de tiempo real
- **Gin**: Framework web
- **Viper**: Configuración

## 🚨 Notificaciones

Cuando un sensor se desconecta:
1. **Console Log**: Mensaje detallado en logs
2. **HTTP Request**: (Opcional) POST al servicio de notificaciones
3. **Estado Persistente**: Actualización en MongoDB

## 🔍 Ejemplos de Uso

### Consultar estado de sensores de un activo
```bash
# Obtener activo con estado detallado de sensores
curl http://localhost:8090/activo/CALDERA_001/sensores/estado

# Respuesta incluye:
# - Información completa del activo
# - Estado de cada sensor (connected/disconnected/never_connected) 
# - Última actividad y estadísticas por sensor
# - Resumen consolidado de estados
```

### Flujo típico de monitoreo
```bash
# 1. Consultar todos los activos
curl http://localhost:8090/activo

# 2. Ver estado de sensores de un activo específico  
curl http://localhost:8090/activo/ACTIVO_ID/sensores/estado

# 3. Ver estadísticas globales de sensores
curl http://localhost:8090/api/sensors/stats

# 4. Forzar verificación de desconexiones
curl -X POST http://localhost:8090/api/sensors/check-disconnected
```