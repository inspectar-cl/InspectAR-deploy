# Microservicio de Gestión de Activos y Técnicos

Este microservicio maneja la gestión de activos industriales, técnicos especializados y acciones de mantenimiento colaborativas.

## Funcionalidades implementadas

### HdU16 - Lista de contactos de técnicos especializados
- Crear técnicos con especialidades
- Listar todos los técnicos disponibles
- Consultar información de un técnico específico
- Actualizar disponibilidad de técnicos

### HdU13 - Acciones de mantención colaborativas
- Crear acciones de mantenimiento (preventivo, correctivo, emergencia)
- Asignar técnicos a acciones específicas
- Actualizar estado de acciones (pendiente, en_progreso, completado)
- Consultar acciones por activo
- Listar acciones pendientes con prioridad

### HdU04 - Reportes automáticos (TODO)
- Generación automática de reportes
- Reportes por activo
- Reportes periódicos

## Estructura del proyecto

```
gestion/
├── cmd/
│   └── main.go                    # Punto de entrada
├── config/
│   └── config.yaml               # Configuración
├── api/
│   └── router/
│       └── router.go             # Rutas HTTP
├── internal/
│   ├── models/
│   │   └── models.go             # Estructuras de datos
│   ├── handlers/
│   │   ├── tecnico_handler.go    # Handlers de técnicos
│   │   └── accion_handler.go     # Handlers de acciones
│   ├── services/
│   │   ├── tecnico_service.go    # Lógica de negocio técnicos
│   │   └── accion_service.go     # Lógica de negocio acciones
│   ├── repository/
│   │   ├── tecnico_repository.go # Acceso a datos técnicos
│   │   └── accion_repository.go  # Acceso a datos acciones
│   └── database/
│       └── postgres.go           # Conexión PostgreSQL
├── Dockerfile
└── README.md
```

## Endpoints disponibles

### Técnicos

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/tecnicos` | Listar todos los técnicos |
| POST | `/tecnicos` | Crear un nuevo técnico |
| GET | `/tecnicos/:id` | Obtener técnico por ID |
| PUT | `/tecnicos/:id/disponibilidad` | Actualizar disponibilidad |

### Acciones de Mantenimiento

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/acciones` | Crear nueva acción de mantenimiento |
| GET | `/activos/:id/acciones` | Obtener acciones de un activo |
| PUT | `/acciones/:id/estado` | Actualizar estado de acción |
| GET | `/acciones/pendientes` | Listar acciones pendientes |

### Health Check

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/health` | Estado del servicio |

## Ejemplos de uso

### Crear un técnico
```bash
curl -X POST http://localhost:8092/tecnicos \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Juan Pérez",
    "email": "juan.perez@empresa.com",
    "telefono": "+56912345678",
    "especialidad": "Sistemas Hidráulicos"
  }'
```

### Crear acción de mantenimiento
```bash
curl -X POST http://localhost:8092/acciones \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 1,
    "tecnico_id": 1,
    "tipo": "preventivo",
    "descripcion": "Revisión mensual de válvulas",
    "prioridad": "media"
  }'
```

### Listar técnicos disponibles
```bash
curl http://localhost:8092/tecnicos
```

### Actualizar estado de acción
```bash
curl -X PUT http://localhost:8092/acciones/1/estado \
  -H "Content-Type: application/json" \
  -d '{"estado": "en_progreso"}'
```

## Configuración

Edita `config/config.yaml` para configurar la conexión a PostgreSQL:

```yaml
postgres:
  host: "localhost"
  port: 5432
  user: "gestion_user"
  password: "gestion_pass"
  dbname: "gestion_db"
  sslmode: "disable"

server:
  port: "8092"
```

## Ejecución

### Con Go directo
```bash
cd microservicios/gestion
go run cmd/main.go
```

### Con Docker
```bash
docker build -t gestion-service .
docker run -p 8092:8092 gestion-service
```

## Base de datos

El servicio requiere una base de datos PostgreSQL con las siguientes tablas:
- `tecnicos`: Información de técnicos especializados
- `activos`: Catálogo de activos (sincronizado con ParserService)
- `acciones_mantenimiento`: Acciones de mantenimiento colaborativas
- `reportes`: Reportes automáticos generados

**Nota**: La estructura de la base de datos debe ser creada previamente. Ver `database/gestion/` para scripts de inicialización.

## Próximos pasos

1. Implementar handlers de reportes (HdU04)
2. Sincronización con ParserService para activos
3. Sistema de notificaciones para acciones críticas
4. Dashboard de gestión en tiempo real
