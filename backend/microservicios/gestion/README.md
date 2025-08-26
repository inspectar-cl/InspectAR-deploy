# Microservicio de Gestión

Microservicio encargado de la gestión de técnicos especializados, acciones de mantenimiento colaborativas y generación de reportes automáticos.

## Funcionalidades

### HdU16 - Lista de contactos de técnicos especializados
- ✅ Crear técnicos con especialidades
- ✅ Listar todos los técnicos disponibles
- ✅ Listar técnicos relacionados con un activo específico
- ✅ Listar técnicos relacionados con un edificio (indirectamente a través de activos)
- ✅ Consultar información de un técnico específico
- ✅ Actualizar estado autorizado de técnicos
- ✅ Asignar técnicos a activos (relación muchos a muchos)
- 🎯 **✅ Obtener activos asociados a un técnico** (`GET /activos-de-tecnico/{tecnico_id}`)

### HdU13 - Acciones de mantención colaborativas
- ✅ Crear acciones de mantenimiento (preventivo, correctivo, emergencia)
- ✅ Asignar técnicos a acciones específicas
- ✅ Actualizar estado de acciones (pendiente, en_progreso, completado)
- ✅ Consultar acciones por técnico
- ✅ Consultar acciones por activo
- ✅ Listar acciones pendientes con prioridad

### HdU04 - Reportes automáticos
- ✅ Generación automática de reportes por activo (PDF)
- ✅ Reportes por activo con información completa

## API Endpoints

### 🏥 Health Check

#### Verificar estado del servicio
```bash
curl -X GET http://localhost:8092/health
```
**Respuesta:**
```json
{
  "status": "ok",
  "service": "gestion"
}
```

---

### 🧑‍🔧 Técnicos (HdU16)

#### Crear técnico
```bash
curl -X POST http://localhost:8092/tecnicos \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Ana Silva",
    "email": "ana.silva@empresa.com",
    "telefono": "+56955667788",
    "especialidad": "Instrumentación"
  }'
```
**Respuesta:**
```json
{
  "id": 9,
  "nombre": "Ana Silva",
  "email": "ana.silva@empresa.com",
  "telefono": "+56955667788",
  "especialidad": "Instrumentación",
  "autorizado": true,
  "creado_en": "2025-08-25T15:30:00Z"
}
```

#### Listar todos los técnicos autorizados
```bash
curl -X GET http://localhost:8092/tecnicos
```
**Respuesta:**
```json
[
  {
    "id": 1,
    "nombre": "Juan Pérez",
    "email": "juan.perez@inspectar.com",
    "telefono": "+56912345678",
    "especialidad": "Sistemas Hidráulicos",
    "autorizado": true,
    "creado_en": "2025-08-20T22:49:18.024565Z"
  },
  {
    "id": 2,
    "nombre": "María González",
    "email": "maria.gonzalez@inspectar.com",
    "telefono": "+56987654321",
    "especialidad": "Electricidad Industrial",
    "autorizado": true,
    "creado_en": "2025-08-20T22:49:18.024565Z"
  }
]
```

#### Listar técnicos por activo
```bash
curl -X GET http://localhost:8092/tecnicos/activo/1
```
**Respuesta:**
```json
[
  {
    "id": 1,
    "nombre": "Juan Pérez",
    "email": "juan.perez@inspectar.com",
    "telefono": "+56912345678",
    "especialidad": "Sistemas Hidráulicos",
    "autorizado": true,
    "creado_en": "2025-08-20T22:49:18.024565Z"
  },
  {
    "id": 6,
    "nombre": "Laura Fernández",
    "email": "laura.fernandez@inspectar.com",
    "telefono": "+56933445566",
    "especialidad": "Calderas y Vapor",
    "autorizado": true,
    "creado_en": "2025-08-20T22:49:18.024565Z"
  }
]
```

#### Listar técnicos por edificio
```bash
curl -X GET http://localhost:8092/tecnicos/edificio/1
```
**Respuesta:**
```json
[
  {
    "id": 1,
    "nombre": "Juan Pérez",
    "email": "juan.perez@inspectar.com",
    "telefono": "+56912345678",
    "especialidad": "Sistemas Hidráulicos",
    "autorizado": true,
    "creado_en": "2025-08-20T22:49:18.024565Z"
  }
]
```

#### Obtener técnico específico
```bash
curl -X GET http://localhost:8092/tecnicos/1
```
**Respuesta:**
```json
{
  "id": 1,
  "nombre": "Juan Pérez",
  "email": "juan.perez@inspectar.com",
  "telefono": "+56912345678",
  "especialidad": "Sistemas Hidráulicos",
  "autorizado": true,
  "creado_en": "2025-08-20T22:49:18.024565Z"
}
```

#### 🎯 **RUTA PRINCIPAL: Obtener activos asociados a un técnico**
```bash
curl -X GET http://localhost:8092/activos-de-tecnico/1
```
**Respuesta:**
```json
[
  {
    "id": 3,
    "activo_id": "AC-1003",
    "nombre": "Bomba Hidráulica 1",
    "tipo": "Bomba",
    "estado": "operativo",
    "ubicacion": "Sala de Bombas",
    "edificio_id": 2,
    "creado_en": "2025-08-20T22:49:18.026008Z"
  },
  {
    "id": 1,
    "activo_id": "AC-1001",
    "nombre": "Caldera Principal",
    "tipo": "Caldera",
    "estado": "operativo",
    "ubicacion": "Sala de Calderas 1",
    "edificio_id": 1,
    "creado_en": "2025-08-20T22:49:18.026008Z"
  }
]
```

#### Actualizar estado autorizado
```bash
curl -X PUT http://localhost:8092/tecnicos/1/autorizado \
  -H "Content-Type: application/json" \
  -d '{"autorizado": false}'
```
**Respuesta:**
```json
{
  "mensaje": "Estado autorizado actualizado correctamente"
}
```

#### Asignar técnico a activo
```bash
curl -X POST http://localhost:8092/activos/1/tecnicos \
  -H "Content-Type: application/json" \
  -d '{"tecnico_id": 2}'
```
**Respuesta:**
```json
{
  "mensaje": "Técnico asignado al activo correctamente"
}
```

---

### 🔧 Acciones de Mantenimiento (HdU13)

#### Crear acción de mantenimiento
```bash
curl -X POST http://localhost:8092/acciones \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 1,
    "tecnico_id": 1,
    "tipo": "preventivo",
    "descripcion": "Revisión mensual de válvulas",
    "prioridad": "alta"
  }'
```
**Respuesta:**
```json
{
  "id": 11,
  "activo_id": 1,
  "tecnico_id": 1,
  "tipo": "preventivo",
  "descripcion": "Revisión mensual de válvulas",
  "estado": "pendiente",
  "prioridad": "alta",
  "fecha_inicio": "2025-08-25T15:30:00Z",
  "creado_en": "2025-08-25T15:30:00Z"
}
```

#### Obtener acciones por técnico
```bash
curl -X GET http://localhost:8092/acciones/tecnico/1
```
**Respuesta:**
```json
[
  {
    "id": 1,
    "activo_id": 1,
    "tecnico_id": 1,
    "tipo": "preventivo",
    "descripcion": "Inspección mensual de calderas",
    "estado": "pendiente",
    "prioridad": "alta",
    "fecha_inicio": "2025-08-20T22:49:18.027847Z",
    "creado_en": "2025-08-20T22:49:18.027847Z",
    "activo": "Caldera Principal",
    "tecnico": "Juan Pérez"
  }
]
```

#### Obtener acciones por activo
```bash
curl -X GET http://localhost:8092/acciones/activo/1
```
**Respuesta:**
```json
[
  {
    "id": 1,
    "activo_id": 1,
    "tecnico_id": 1,
    "tipo": "preventivo",
    "descripcion": "Inspección mensual de calderas",
    "estado": "pendiente",
    "prioridad": "alta",
    "fecha_inicio": "2025-08-20T22:49:18.027847Z",
    "creado_en": "2025-08-20T22:49:18.027847Z",
    "activo": "Caldera Principal",
    "tecnico": "Juan Pérez"
  }
]
```

#### Actualizar estado de acción
```bash
curl -X PUT http://localhost:8092/acciones/1/estado \
  -H "Content-Type: application/json" \
  -d '{"estado": "en_progreso"}'
```
**Respuesta:**
```json
{
  "mensaje": "Estado de acción actualizado correctamente"
}
```

#### Listar acciones pendientes con prioridad
```bash
curl -X GET http://localhost:8092/acciones/pendientes
```
**Respuesta:**
```json
[
  {
    "id": 1,
    "activo_id": 1,
    "tecnico_id": 1,
    "tipo": "preventivo",
    "descripcion": "Inspección mensual de calderas",
    "estado": "pendiente",
    "prioridad": "alta",
    "fecha_inicio": "2025-08-20T22:49:18.027847Z",
    "creado_en": "2025-08-20T22:49:18.027847Z",
    "activo": "Caldera Principal",
    "tecnico": "Juan Pérez"
  }
]
```

---

### 📊 Reportes (HdU04)

#### Generar reporte PDF por activo
```bash
curl -X POST http://localhost:8092/reportes/activo/1
```
**Respuesta:** Archivo PDF para descarga

#### Obtener reportes de un activo
```bash
curl -X GET http://localhost:8092/reportes/activo/1
```
**Respuesta:**
```json
[
  {
    "id": 1,
    "activo_id": 1,
    "tipo": "mantenimiento",
    "formato": "PDF",
    "contenido": "Reporte detallado de mantenimiento...",
    "generado_por": "Sistema",
    "creado_en": "2025-08-20T22:49:18.028743Z",
    "activo": "Caldera Principal"
  }
]
```

## Estructura de Base de Datos

### Tablas principales:
- **`tecnicos`** - Técnicos especializados con campo `autorizado`
- **`edificios`** - Edificios donde se ubican los activos
- **`activos`** - Activos industriales (relacionados con edificios)
- **`acciones_mantenimiento`** - Acciones de mantenimiento colaborativas
- **`reportes`** - Reportes automáticos generados
- **`activos_tecnicos`** - Tabla intermedia para relación muchos a muchos

### Relaciones:
- `edificios` ↔ `activos` (uno a muchos)
- `activos` ↔ `tecnicos` (muchos a muchos vía `activos_tecnicos`)
- `activos` ↔ `acciones_mantenimiento` (uno a muchos)
- `tecnicos` ↔ `acciones_mantenimiento` (uno a muchos)

## Ejemplos de Uso

### 🎯 Obtener activos de un técnico específico (RUTA PRINCIPAL)
```bash
curl -X GET http://localhost:8092/activos-de-tecnico/1
```

### Listar técnicos autorizados
```bash
curl -X GET http://localhost:8092/tecnicos
```

### Obtener técnicos asociados a un activo
```bash
curl -X GET http://localhost:8092/tecnicos/activo/1
```

### Crear nueva acción de mantenimiento
```bash
curl -X POST http://localhost:8092/acciones \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 1,
    "tecnico_id": 1,
    "tipo": "preventivo",
    "descripcion": "Revisión mensual de válvulas",
    "prioridad": "alta"
  }'
```

### Generar reporte automático
```bash
curl -X POST http://localhost:8092/reportes/activo/1
```

## Configuración

El microservicio utiliza PostgreSQL como base de datos. La configuración se encuentra en:
- `config/config.yaml` - Configuración general
- Variables de entorno para conexión a base de datos

## Ejecución

```bash
# Ejecutar la base de datos
docker-compose up gestion-db

# Ejecutar el microservicio
go run cmd/main.go
```

El servicio estará disponible en el puerto `8092`.

## Estados y Tipos

### Estados de Acciones:
- `pendiente` - Acción creada pero no iniciada
- `en_progreso` - Acción siendo ejecutada
- `completado` - Acción finalizada
- `cancelado` - Acción cancelada

### Tipos de Mantenimiento:
- `preventivo` - Mantenimiento programado
- `correctivo` - Reparación de fallas
- `emergencia` - Mantenimiento urgente

### Prioridades:
- `baja` - Prioridad baja
- `media` - Prioridad media
- `alta` - Prioridad alta
- `critica` - Prioridad crítica
