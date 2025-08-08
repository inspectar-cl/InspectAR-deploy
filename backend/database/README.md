# Database Directory

Esta carpeta contiene las configuraciones de bases de datos para todos los microservicios del proyecto.

## Estructura

```
database/
└── notification/
    ├── Dockerfile                      # Configuración PostgreSQL
    └── docker-entrypoint-initdb.d/
        ├── 001_create_tables.sql       # Creación de tablas
        ├── 002_create_triggers.sql     # Triggers y funciones
        └── 003_insert_data.sql         # Datos de ejemplo
```

## notification/

Base de datos PostgreSQL para el microservicio de notificaciones.

### Tablas:
- **edificios**: Información de edificios
- **usuarios**: Usuarios del sistema con scope y contacto
- **activos**: Activos asociados a edificios
- **notificaciones**: Notificaciones enviadas/pendientes

### Configuración:
- **Usuario**: notification_user
- **Password**: notification_pass
- **Base de datos**: notification_db
- **Puerto**: 5432

### Uso:
Esta base de datos se levanta automáticamente con `docker-compose up` y está disponible para el microservicio `notification-service`.
