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

### Técnicos

#### Crear técnico
```http
POST /tecnicos
Content-Type: application/json

{
  "nombre": "Juan Pérez",
  "email": "juan.perez@empresa.com",
  "telefono": "+56912345678",
  "especialidad": "Sistemas Hidráulicos"
}
```

#### Listar todos los técnicos
```http
GET /tecnicos
```

#### Listar técnicos por activo
```http
GET /tecnicos/activo/{activo_id}?solo_autorizados=true
```

#### Listar técnicos por edificio
```http
GET /tecnicos/edificio/{edificio_id}?solo_autorizados=true
```

#### Obtener técnico específico
```http
GET /tecnicos/{id}
```

#### Actualizar estado autorizado
```http
PUT /tecnicos/{id}/autorizado
Content-Type: application/json

{
  "autorizado": true
}
```

#### Asignar técnico a activo
```http
POST /activos/{activo_id}/tecnicos
Content-Type: application/json

{
  "tecnico_id": 1
}
```

### Acciones de Mantenimiento

#### Crear acción de mantenimiento
```http
POST /acciones
Content-Type: application/json

{
  "activo_id": 1,
  "tecnico_id": 1,
  "tipo": "preventivo",
  "descripcion": "Revisión mensual de válvulas",
  "prioridad": "alta"
}
```

#### Obtener acciones por técnico
```http
GET /acciones/tecnico/{tecnico_id}
```

#### Obtener acciones por activo
```http
GET /acciones/activo/{activo_id}
```

#### Actualizar estado de acción
```http
PUT /acciones/{id}/estado
Content-Type: application/json

{
  "estado": "en_progreso"
}
```

#### Listar acciones pendientes con prioridad
```http
GET /acciones/pendientes
```

### Reportes

#### Generar reporte PDF por activo
```http
POST /reportes/activo/{activo_id}
```
Retorna: Archivo PDF para descarga

#### Obtener reportes de un activo
```http
GET /reportes/activo/{activo_id}
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
