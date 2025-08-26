-- Crear tablas para el microservicio de gestión
-- Archivo: 001_create_tables.sql

-- Tabla de técnicos especializados
CREATE TABLE IF NOT EXISTS tecnicos (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    telefono VARCHAR(50),
    especialidad VARCHAR(100) NOT NULL,
    autorizado BOOLEAN DEFAULT true,
    creado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de edificios
CREATE TABLE IF NOT EXISTS edificios (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    direccion VARCHAR(255),
    creado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de activos (espejo del ParserService para gestión)
CREATE TABLE IF NOT EXISTS activos (
    id SERIAL PRIMARY KEY,
    activo_id VARCHAR(100) UNIQUE NOT NULL,
    nombre VARCHAR(255) NOT NULL,
    tipo VARCHAR(100) NOT NULL,
    estado VARCHAR(50) DEFAULT 'operativo',
    ubicacion VARCHAR(255),
    edificio_id INTEGER REFERENCES edificios(id) ON DELETE SET NULL,
    creado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de acciones de mantenimiento colaborativas
CREATE TABLE IF NOT EXISTS acciones_mantenimiento (
    id SERIAL PRIMARY KEY,
    activo_id INTEGER REFERENCES activos(id) ON DELETE CASCADE,
    tecnico_id INTEGER REFERENCES tecnicos(id) ON DELETE SET NULL,
    tipo VARCHAR(50) NOT NULL CHECK (tipo IN ('preventivo', 'correctivo', 'emergencia')),
    descripcion TEXT,
    estado VARCHAR(50) DEFAULT 'pendiente' CHECK (estado IN ('pendiente', 'en_progreso', 'completado', 'cancelado')),
    prioridad VARCHAR(20) DEFAULT 'media' CHECK (prioridad IN ('baja', 'media', 'alta', 'critica')),
    fecha_inicio TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    fecha_fin TIMESTAMP,
    creado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de reportes automáticos
CREATE TABLE IF NOT EXISTS reportes (
    id SERIAL PRIMARY KEY,
    activo_id INTEGER REFERENCES activos(id) ON DELETE CASCADE,
    tipo_reporte VARCHAR(50) NOT NULL CHECK (tipo_reporte IN ('semanal', 'mensual', 'incidente', 'mantenimiento')),
    contenido TEXT,
    generado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    estado VARCHAR(50) DEFAULT 'generado' CHECK (estado IN ('generado', 'enviado', 'archivado'))
);

-- Tabla intermedia para la relación muchos a muchos entre activos y técnicos
CREATE TABLE IF NOT EXISTS activos_tecnicos (
    activo_id INTEGER NOT NULL REFERENCES activos(id) ON DELETE CASCADE,
    tecnico_id INTEGER NOT NULL REFERENCES tecnicos(id) ON DELETE CASCADE,
    asignado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (activo_id, tecnico_id)
);

-- Índices para mejorar rendimiento
CREATE INDEX IF NOT EXISTS idx_tecnicos_especialidad ON tecnicos(especialidad);
CREATE INDEX IF NOT EXISTS idx_tecnicos_autorizado ON tecnicos(autorizado);
CREATE INDEX IF NOT EXISTS idx_activos_estado ON activos(estado);
CREATE INDEX IF NOT EXISTS idx_activos_tipo ON activos(tipo);
CREATE INDEX IF NOT EXISTS idx_activos_edificio ON activos(edificio_id);
CREATE INDEX IF NOT EXISTS idx_acciones_estado ON acciones_mantenimiento(estado);
CREATE INDEX IF NOT EXISTS idx_acciones_prioridad ON acciones_mantenimiento(prioridad);
CREATE INDEX IF NOT EXISTS idx_acciones_activo ON acciones_mantenimiento(activo_id);
CREATE INDEX IF NOT EXISTS idx_reportes_tipo ON reportes(tipo_reporte);
CREATE INDEX IF NOT EXISTS idx_reportes_fecha ON reportes(generado_en);
CREATE INDEX IF NOT EXISTS idx_activos_tecnicos_activo ON activos_tecnicos(activo_id);
CREATE INDEX IF NOT EXISTS idx_activos_tecnicos_tecnico ON activos_tecnicos(tecnico_id);

COMMENT ON TABLE tecnicos IS 'Técnicos especializados para mantenimiento (HdU16)';
COMMENT ON TABLE edificios IS 'Edificios donde se ubican los activos';
COMMENT ON TABLE activos IS 'Activos industriales gestionados';
COMMENT ON TABLE acciones_mantenimiento IS 'Acciones de mantenimiento colaborativas (HdU13)';
COMMENT ON TABLE reportes IS 'Reportes automáticos generados (HdU04)';
