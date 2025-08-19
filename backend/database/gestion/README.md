# Base de Datos para Microservicio de Gestión

Esta carpeta contiene la configuración y scripts de inicialización para la base de datos PostgreSQL del microservicio de gestión.

## Estructura

```
database/gestion/
├── Dockerfile                           # Imagen personalizada de PostgreSQL
└── docker-entrypoint-initdb.d/
    ├── 001_create_tables.sql           # Creación de tablas y estructura
    └── 002_insert_sample_data.sql      # Datos de ejemplo para testing
```

## Tablas creadas

### tecnicos
- **Propósito**: Lista de contactos de técnicos especializados (HdU16)
- **Campos**: id, nombre, email, telefono, especialidad, disponible, creado_en
- **Datos ejemplo**: 8 técnicos con diferentes especialidades

### activos  
- **Propósito**: Activos industriales gestionados
- **Campos**: id, activo_id, nombre, tipo, estado, ubicacion, edificio_id, creado_en
- **Datos ejemplo**: 8 activos (calderas, compresores, bombas, motores, etc.)

### acciones_mantenimiento
- **Propósito**: Acciones de mantenimiento colaborativas (HdU13)
- **Campos**: id, activo_id, tecnico_id, tipo, descripcion, estado, prioridad, fecha_inicio, fecha_fin, creado_en
- **Datos ejemplo**: 10 acciones con diferentes estados y prioridades

### reportes
- **Propósito**: Reportes automáticos generados (HdU04)
- **Campos**: id, activo_id, tipo_reporte, contenido, generado_en, estado
- **Datos ejemplo**: 8 reportes de diferentes tipos

## Conexión

La base de datos está configurada con:
- **Host**: gestion-db (dentro de Docker) / localhost:5433 (desde fuera)
- **Usuario**: gestion_user
- **Contraseña**: gestion_pass
- **Base de datos**: gestion_db

## Datos de prueba incluidos

### Técnicos especializados:
- Juan Pérez (Sistemas Hidráulicos) ✅ Disponible
- María González (Electricidad Industrial) ✅ Disponible  
- Carlos Rodríguez (Mecánica Industrial) ❌ No disponible
- Ana Silva (Instrumentación) ✅ Disponible
- Pedro Morales (Sistemas de Control) ✅ Disponible
- Laura Fernández (Calderas y Vapor) ✅ Disponible
- Roberto Castro (Refrigeración) ❌ No disponible
- Isabel Herrera (Sistemas Neumáticos) ✅ Disponible

### Activos de ejemplo:
- AC-1001: Caldera Principal (operativo)
- AC-1002: Compresor Auxiliar (mantenimiento)
- AC-1003: Bomba Hidráulica 1 (operativo)
- AC-1004: Motor Eléctrico Principal (operativo)
- AC-1005: Sistema de Ventilación (alerta)
- AC-1006: Transformador Eléctrico (operativo)
- AC-1007: Chiller Industrial (operativo)
- AC-1008: Generador de Emergencia (standby)

### Acciones de mantenimiento:
- 3 acciones pendientes (diferentes prioridades)
- 2 acciones en progreso
- 5 acciones completadas

### Reportes generados:
- 4 reportes semanales
- 2 reportes mensuales
- 2 reportes de incidentes
- 1 reporte de mantenimiento

## Consultas de ejemplo

```sql
-- Listar técnicos disponibles por especialidad
SELECT * FROM tecnicos WHERE disponible = true ORDER BY especialidad;

-- Acciones pendientes de alta prioridad
SELECT am.*, a.nombre as activo_nombre, t.nombre as tecnico_nombre 
FROM acciones_mantenimiento am
JOIN activos a ON am.activo_id = a.id
JOIN tecnicos t ON am.tecnico_id = t.id
WHERE am.estado = 'pendiente' AND am.prioridad IN ('alta', 'critica');

-- Reportes de incidentes no archivados
SELECT r.*, a.nombre as activo_nombre
FROM reportes r
JOIN activos a ON r.activo_id = a.id
WHERE r.tipo_reporte = 'incidente' AND r.estado != 'archivado';
```
