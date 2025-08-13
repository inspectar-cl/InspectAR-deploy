# ParserService - Microservicio IoT para InspectAR

El ParserService (iot-service) es un microservicio encargado de la gestión de activos industriales y sus lecturas de sensores en tiempo real. Proporciona una API REST para crear, consultar y actualizar activos, así como para registrar y consultar lecturas de sensores.

## Características principales

- Gestión de activos industriales con múltiples sensores
- Registro de lecturas de sensores en tiempo real
- Consultas históricas y en tiempo real de datos de sensores
- Actualización de estados de activos
- Persistencia en MongoDB (metadatos) e InfluxDB (series temporales)

## Arquitectura

- **API REST**: Implementada con Gin Framework
- **Base de datos**: 
  - MongoDB: Almacenamiento de activos y metadatos
  - InfluxDB: Almacenamiento de series temporales (lecturas de sensores)
- **Implementado en**: Go 1.23

## Endpoints disponibles

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/activo` | Listar todos los activos |
| POST | `/activo` | Crear un nuevo activo |
| GET | `/activo/:activo_id` | Obtener un activo específico |
| PUT | `/activo/:activo_id/estado` | Actualizar el estado de un activo |
| POST | `/lectura` | Registrar una lectura de sensor |
| GET | `/lectura/:activo_id/datos` | Obtener lecturas de todos los sensores de un activo |
| GET | `/lectura/:activo_id/datos/ultimo` | Obtener la última lectura de cada sensor de un activo |

## Ejemplos de uso

### 1. Listar todos los activos

**Petición:**
```bash
curl http://34.39.160.195:8090/activo
```

**Respuesta:**
```json
[
  {
    "id": "64f8a9c87c1b3e001e0b6d23",
    "activo_id": "AC-1001",
    "nombre": "Caldera Principal",
    "estado": "operativo",
    "ubicacion": "Sala 1",
    "sensores": [
      {
        "sensor_id": "TEMP-01",
        "tipo": "temperatura",
        "unidad": "C"
      },
      {
        "sensor_id": "PRES-01",
        "tipo": "presion",
        "unidad": "bar"
      }
    ],
    "id_edificio": "ED-01"
  },
  {
    "id": "64f8b1a87c1b3e001e0b6d24",
    "activo_id": "AC-1002",
    "nombre": "Compresor Auxiliar",
    "estado": "mantenimiento",
    "ubicacion": "Sala 2",
    "sensores": [
      {
        "sensor_id": "TEMP-02",
        "tipo": "temperatura",
        "unidad": "C"
      }
    ],
    "id_edificio": "ED-01"
  }
]
```

### 2. Crear un activo

**Petición:**
```bash
curl -X POST http://34.39.160.195:8090/activo \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": "AC-1003",
    "nombre": "Bomba Hidráulica", 
    "estado": "operativo",
    "ubicacion": "Sala 3",
    "sensores": [
      {"sensor_id": "FLOW-01", "tipo": "flujo", "unidad": "l/min"},
      {"sensor_id": "PRES-03", "tipo": "presion", "unidad": "bar"}
    ],
    "id_edificio": "ED-01"
  }'
```

**Respuesta:**
```json
{
  "activo_id": "AC-1003"
}
```

### 3. Obtener un activo específico

**Petición:**
```bash
curl http://34.39.160.195:8090/activo/AC-1001
```

**Respuesta:**
```json
{
  "id": "64f8a9c87c1b3e001e0b6d23",
  "activo_id": "AC-1001",
  "nombre": "Caldera Principal",
  "estado": "operativo",
  "ubicacion": "Sala 1",
  "sensores": [
    {
      "sensor_id": "TEMP-01",
      "tipo": "temperatura",
      "unidad": "C"
    },
    {
      "sensor_id": "PRES-01",
      "tipo": "presion",
      "unidad": "bar"
    }
  ],
  "id_edificio": "ED-01"
}
```

### 4. Registrar una lectura de sensor

**Petición:**
```bash
curl -X POST http://34.39.160.195:8090/lectura \
  -H "Content-Type: application/json" \
  -d '{
    "sensor_id": "TEMP-01",
    "valor": 85.4,
    "timestamp": "2025-08-13T03:45:10Z"
  }'
```

**Respuesta:**
```json
{
  "mensaje": "Lectura registrada en InfluxDB"
}
```

### 5. Obtener lecturas de un activo

**Petición:**
```bash
curl http://34.39.160.195:8090/lectura/AC-1001/datos
```

**Respuesta:**
```json
{
  "id": "64f8a9c87c1b3e001e0b6d23",
  "activo_id": "AC-1001",
  "nombre": "Caldera Principal",
  "ubicacion": "Sala 1",
  "estado": "operativo",
  "sensores": [
    {
      "sensor_id": "TEMP-01",
      "datos": [
        {
          "timestamp": "2025-08-13T03:45:10Z",
          "valor": 85.4
        },
        {
          "timestamp": "2025-08-13T03:46:15Z",
          "valor": 86.2
        }
      ]
    },
    {
      "sensor_id": "PRES-01",
      "datos": [
        {
          "timestamp": "2025-08-13T03:45:11Z",
          "valor": 3.2
        },
        {
          "timestamp": "2025-08-13T03:46:16Z",
          "valor": 3.3
        }
      ]
    }
  ]
}
```

### 6. Obtener último dato de cada sensor

**Petición:**
```bash
curl http://34.39.160.195:8090/lectura/AC-1001/datos/ultimo
```

**Respuesta:**
```json
{
  "TEMP-01": {
    "timestamp": "2025-08-13T03:46:15Z",
    "valor": 86.2
  },
  "PRES-01": {
    "timestamp": "2025-08-13T03:46:16Z",
    "valor": 3.3
  }
}
```

### 7. Actualizar estado de un activo

**Petición:**
```bash
curl -X PUT http://34.39.160.195:8090/activo/AC-1001/estado \
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

## Configuración

El servicio utiliza un archivo de configuración `config.yaml` que define las conexiones a las bases de datos:

```yaml
mongu:
  uri: "mongodb://34.39.160.195:27017"
  database: "iot_db"

influxdb:
  url: "http://34.39.160.195:8086"
  token: "my-super-token"
  org: "my-org"
  bucket: "sensores"

server:
  port: "8090"
```

## Ejecución

### Con Docker Compose
```bash
# Reconstruir la imagen (si se modificó la configuración)
docker-compose build iot-service

# Ejecutar el servicio
docker-compose up iot-service
```

### Directamente (para desarrollo)
```bash
cd microservicios/ParserService
go run cmd/main.go
```

## Requisitos para desarrollo

- Go 1.23 o superior
- MongoDB 4.4 o superior
- InfluxDB 2.7 o superior
