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
    latitud DECIMAL(10, 8),
    longitud DECIMAL(11, 8),
    creado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de activos (espejo del ParserService para gestión)
CREATE TABLE IF NOT EXISTS activos (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    tipo VARCHAR(100) NOT NULL CHECK (tipo IN ('caldera', 'bomba de agua', 'ascensor', 'transformador')),
    descripcion TEXT,
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
    -- NUEVOS CAMPOS PARA OBSERVACIONES Y ESTRUCTURA
    observaciones_analista TEXT DEFAULT '', -- Comentarios editables del analista
    autor_analista VARCHAR(255) DEFAULT '', -- Nombre del analista que hizo el reporte
    estructura_informe JSONB DEFAULT '{}', -- Estructura completa del informe en JSON
    metadata_informe JSONB DEFAULT '{}', -- Metadata adicional (versión, secciones, etc.)
    version_reporte INTEGER DEFAULT 1, -- Versión del reporte (para tracking de cambios)
    estado_revision VARCHAR(50) DEFAULT 'pendiente' CHECK (estado_revision IN ('pendiente', 'en_revision', 'aprobado', 'rechazado')),
    fecha_revision TIMESTAMP NULL, -- Fecha de la última revisión
    revisor VARCHAR(255) DEFAULT '', -- Nombre del revisor
    --
    generado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    estado VARCHAR(50) DEFAULT 'generado' CHECK (estado IN ('generado', 'enviado', 'archivado'))
);

-- Tabla de empresas de mantención
CREATE TABLE IF NOT EXISTS empresas (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL UNIQUE,
    rut VARCHAR(12) NOT NULL UNIQUE,
    telefono VARCHAR(20),
    email VARCHAR(150),
    direccion VARCHAR(200),
    activo BOOLEAN DEFAULT true,
    fecha_registro TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    fecha_actualizacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Actualizar tabla de técnicos para incluir empresa
ALTER TABLE tecnicos 
ADD COLUMN IF NOT EXISTS apellido VARCHAR(255),
ADD COLUMN IF NOT EXISTS empresa_id INTEGER REFERENCES empresas(id),
ADD COLUMN IF NOT EXISTS activo BOOLEAN DEFAULT true,
ADD COLUMN IF NOT EXISTS fecha_registro TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS fecha_actualizacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

-- Renombrar columna creado_en a fecha_registro para consistencia
DO $$ 
BEGIN
    IF EXISTS(SELECT * FROM information_schema.columns WHERE table_name='tecnicos' AND column_name='creado_en') THEN
        UPDATE tecnicos SET fecha_registro = creado_en WHERE fecha_registro IS NULL;
        ALTER TABLE tecnicos DROP COLUMN creado_en;
    END IF;
END $$;

-- Tabla de solicitudes técnicas (HdU16)
CREATE TABLE IF NOT EXISTS solicitudes_tecnico (
    id SERIAL PRIMARY KEY,
    tecnico_id INTEGER NOT NULL REFERENCES tecnicos(id) ON DELETE CASCADE,
    residente_id INTEGER NOT NULL,
    activo_id INTEGER NOT NULL REFERENCES activos(id) ON DELETE CASCADE,
    edificio_id INTEGER NOT NULL REFERENCES edificios(id) ON DELETE CASCADE,
    
    -- Contenido de la solicitud
    tipo VARCHAR(50) NOT NULL CHECK (tipo IN ('mantenimiento', 'reparacion', 'inspeccion', 'emergencia', 'consulta')),
    asunto VARCHAR(200) NOT NULL,
    descripcion TEXT NOT NULL,
    prioridad VARCHAR(20) NOT NULL DEFAULT 'media' CHECK (prioridad IN ('baja', 'media', 'alta', 'critica', 'emergencia')),
    
    -- Estado y seguimiento
    estado VARCHAR(50) NOT NULL DEFAULT 'pendiente' CHECK (estado IN ('pendiente', 'enviada', 'recibida', 'en_proceso', 'completada', 'cancelada', 'rechazada')),
    fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    fecha_envio TIMESTAMP,
    fecha_recepcion TIMESTAMP,
    fecha_completado TIMESTAMP,
    fecha_actualizacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Información de contacto
    medio_contacto VARCHAR(50) NOT NULL CHECK (medio_contacto IN ('email', 'telefono', 'sms', 'ambos')),
    telefono_contacto VARCHAR(20),
    email_contacto VARCHAR(150),
    
    -- Respuesta del técnico
    respuesta_tecnico TEXT,
    notas_internas TEXT
);

-- Tabla de archivos adjuntos para solicitudes
CREATE TABLE IF NOT EXISTS archivos_solicitud (
    id SERIAL PRIMARY KEY,
    solicitud_id INTEGER NOT NULL REFERENCES solicitudes_tecnico(id) ON DELETE CASCADE,
    nombre_archivo VARCHAR(255) NOT NULL,
    ruta_archivo VARCHAR(500) NOT NULL,
    tipo_archivo VARCHAR(50),
    tamano_bytes BIGINT,
    fecha_subida TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla intermedia para activos autorizados por técnico
CREATE TABLE IF NOT EXISTS activos_tecnicos_autorizados (
    id SERIAL PRIMARY KEY,
    tecnico_id INTEGER NOT NULL REFERENCES tecnicos(id) ON DELETE CASCADE,
    activo_id INTEGER NOT NULL REFERENCES activos(id) ON DELETE CASCADE,
    edificio_id INTEGER NOT NULL REFERENCES edificios(id) ON DELETE CASCADE,
    fecha_autorizacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    activo BOOLEAN DEFAULT true,
    UNIQUE(tecnico_id, activo_id)
);

-- Tabla intermedia para la relación muchos a muchos entre activos y técnicos (mantener compatibilidad)
CREATE TABLE IF NOT EXISTS activos_tecnicos (
    activo_id INTEGER NOT NULL REFERENCES activos(id) ON DELETE CASCADE,
    tecnico_id INTEGER NOT NULL REFERENCES tecnicos(id) ON DELETE CASCADE,
    asignado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (activo_id, tecnico_id)
);

-- Tabla de usuarios del sistema
CREATE TABLE IF NOT EXISTS usuarios (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    creado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla intermedia para la relación muchos a muchos entre usuarios y edificios
CREATE TABLE IF NOT EXISTS usuarios_edificios (
    usuario_id INTEGER NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    edificio_id INTEGER NOT NULL REFERENCES edificios(id) ON DELETE CASCADE,
    asignado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (usuario_id, edificio_id)
);

-- Tabla de fallos por tipo de activo
CREATE TABLE IF NOT EXISTS fallos (
    id SERIAL PRIMARY KEY,
    tipo_activo VARCHAR(100) NOT NULL CHECK (tipo_activo IN ('caldera', 'bomba_de_agua', 'ascensor', 'transformador')),
    descripcion TEXT NOT NULL,
    prioridad VARCHAR(20) NOT NULL DEFAULT 'media' CHECK (prioridad IN ('media', 'alta')),
    probabilidad_ocurrencia DECIMAL(5,2) NOT NULL CHECK (probabilidad_ocurrencia >= 0 AND probabilidad_ocurrencia <= 100),
    creado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    actualizado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla intermedia para relacionar activos con fallos específicos
CREATE TABLE IF NOT EXISTS activos_fallos (
    id SERIAL PRIMARY KEY,
    activo_id INTEGER NOT NULL REFERENCES activos(id) ON DELETE CASCADE,
    fallo_id INTEGER NOT NULL REFERENCES fallos(id) ON DELETE CASCADE,
    fecha_deteccion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    estado VARCHAR(50) DEFAULT 'detectado' CHECK (estado IN ('detectado', 'en_revision', 'resuelto', 'pendiente')),
    notas TEXT,
    UNIQUE(activo_id, fallo_id, fecha_deteccion)
);

-- Tabla de reportes de fallas hechos por usuarios (tipos_falla)
CREATE TABLE IF NOT EXISTS tipos_falla (
    id_falla SERIAL PRIMARY KEY,
    tipo VARCHAR(50) NOT NULL CHECK (tipo IN ('falla agua', 'falla ascensor', 'falla electricidad', 'falla caldera')),
    descripcion TEXT,
    fecha_publicacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    id_usuario INTEGER NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    id_edificio INTEGER NOT NULL REFERENCES edificios(id) ON DELETE CASCADE,
    estado VARCHAR(50) DEFAULT 'reportado' CHECK (estado IN ('reportado', 'en_revision', 'resuelto', 'rechazado'))
);

-- Tabla de comentarios sobre reportes de fallas
CREATE TABLE IF NOT EXISTS comentarios (
    id_comentario SERIAL PRIMARY KEY,
    id_falla INTEGER NOT NULL REFERENCES tipos_falla(id_falla) ON DELETE CASCADE,
    id_usuario INTEGER NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    comentario TEXT NOT NULL,
    fecha_comentario TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de firmas digitales para usuarios
CREATE TABLE IF NOT EXISTS firmas_digitales (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    nombre_archivo VARCHAR(255) NOT NULL,
    ruta_archivo VARCHAR(500) NOT NULL,
    tipo_mime VARCHAR(100) NOT NULL CHECK (tipo_mime IN ('image/png', 'image/jpeg', 'image/jpg', 'image/svg+xml')),
    formato VARCHAR(20) NOT NULL CHECK (formato IN ('png', 'jpeg', 'jpg', 'svg')),
    datos_firma BYTEA, -- Para almacenar SVG o datos binarios de la firma dibujada
    tamano_bytes BIGINT,
    es_predeterminada BOOLEAN DEFAULT false,
    creado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    actualizado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices para mejorar rendimiento
CREATE INDEX IF NOT EXISTS idx_tecnicos_especialidad ON tecnicos(especialidad);
CREATE INDEX IF NOT EXISTS idx_tecnicos_autorizado ON tecnicos(autorizado);
CREATE INDEX IF NOT EXISTS idx_tecnicos_activo ON tecnicos(activo);
CREATE INDEX IF NOT EXISTS idx_tecnicos_empresa ON tecnicos(empresa_id);
CREATE INDEX IF NOT EXISTS idx_activos_tipo ON activos(tipo);
CREATE INDEX IF NOT EXISTS idx_activos_edificio ON activos(edificio_id);
CREATE INDEX IF NOT EXISTS idx_acciones_estado ON acciones_mantenimiento(estado);
CREATE INDEX IF NOT EXISTS idx_acciones_prioridad ON acciones_mantenimiento(prioridad);
CREATE INDEX IF NOT EXISTS idx_acciones_activo ON acciones_mantenimiento(activo_id);
CREATE INDEX IF NOT EXISTS idx_reportes_tipo ON reportes(tipo_reporte);
CREATE INDEX IF NOT EXISTS idx_reportes_fecha ON reportes(generado_en);
-- NUEVOS ÍNDICES PARA FUNCIONALIDADES EXTENDIDAS
CREATE INDEX IF NOT EXISTS idx_reportes_estado_revision ON reportes(estado_revision);
CREATE INDEX IF NOT EXISTS idx_reportes_autor ON reportes(autor_analista);
CREATE INDEX IF NOT EXISTS idx_reportes_version ON reportes(version_reporte);
CREATE INDEX IF NOT EXISTS idx_activos_tecnicos_activo ON activos_tecnicos(activo_id);
CREATE INDEX IF NOT EXISTS idx_activos_tecnicos_tecnico ON activos_tecnicos(tecnico_id);

-- Índices para solicitudes técnicas
CREATE INDEX IF NOT EXISTS idx_solicitudes_tecnico ON solicitudes_tecnico(tecnico_id);
CREATE INDEX IF NOT EXISTS idx_solicitudes_residente ON solicitudes_tecnico(residente_id);
CREATE INDEX IF NOT EXISTS idx_solicitudes_activo ON solicitudes_tecnico(activo_id);
CREATE INDEX IF NOT EXISTS idx_solicitudes_edificio ON solicitudes_tecnico(edificio_id);
CREATE INDEX IF NOT EXISTS idx_solicitudes_estado ON solicitudes_tecnico(estado);
CREATE INDEX IF NOT EXISTS idx_solicitudes_tipo ON solicitudes_tecnico(tipo);
CREATE INDEX IF NOT EXISTS idx_solicitudes_prioridad ON solicitudes_tecnico(prioridad);
CREATE INDEX IF NOT EXISTS idx_solicitudes_fecha_creacion ON solicitudes_tecnico(fecha_creacion);
CREATE INDEX IF NOT EXISTS idx_activos_tecnicos_autorizados_tecnico ON activos_tecnicos_autorizados(tecnico_id);
CREATE INDEX IF NOT EXISTS idx_activos_tecnicos_autorizados_activo ON activos_tecnicos_autorizados(activo_id);
CREATE INDEX IF NOT EXISTS idx_activos_tecnicos_autorizados_edificio ON activos_tecnicos_autorizados(edificio_id);

-- Índices para usuarios y relación con edificios
CREATE INDEX IF NOT EXISTS idx_usuarios_edificios_usuario ON usuarios_edificios(usuario_id);
CREATE INDEX IF NOT EXISTS idx_usuarios_edificios_edificio ON usuarios_edificios(edificio_id);

-- Índices para fallos
CREATE INDEX IF NOT EXISTS idx_fallos_tipo_activo ON fallos(tipo_activo);
CREATE INDEX IF NOT EXISTS idx_fallos_prioridad ON fallos(prioridad);
CREATE INDEX IF NOT EXISTS idx_fallos_probabilidad ON fallos(probabilidad_ocurrencia);
CREATE INDEX IF NOT EXISTS idx_activos_fallos_activo ON activos_fallos(activo_id);
CREATE INDEX IF NOT EXISTS idx_activos_fallos_fallo ON activos_fallos(fallo_id);
CREATE INDEX IF NOT EXISTS idx_activos_fallos_estado ON activos_fallos(estado);

-- Índices para reportes de fallas de usuarios
CREATE INDEX IF NOT EXISTS idx_tipos_falla_tipo ON tipos_falla(tipo);
CREATE INDEX IF NOT EXISTS idx_tipos_falla_usuario ON tipos_falla(id_usuario);
CREATE INDEX IF NOT EXISTS idx_tipos_falla_edificio ON tipos_falla(id_edificio);
CREATE INDEX IF NOT EXISTS idx_tipos_falla_estado ON tipos_falla(estado);
CREATE INDEX IF NOT EXISTS idx_tipos_falla_fecha ON tipos_falla(fecha_publicacion);
CREATE INDEX IF NOT EXISTS idx_comentarios_falla ON comentarios(id_falla);
CREATE INDEX IF NOT EXISTS idx_comentarios_usuario ON comentarios(id_usuario);
CREATE INDEX IF NOT EXISTS idx_comentarios_fecha ON comentarios(fecha_comentario);

-- Índices para firmas digitales
CREATE INDEX IF NOT EXISTS idx_firmas_usuario ON firmas_digitales(usuario_id);
CREATE INDEX IF NOT EXISTS idx_firmas_predeterminada ON firmas_digitales(es_predeterminada);

-- Comentarios sobre las tablas creadas
COMMENT ON TABLE tecnicos IS 'Técnicos especializados para mantenimiento (HdU16)';
COMMENT ON TABLE edificios IS 'Edificios donde se ubican los activos';
COMMENT ON TABLE activos IS 'Activos industriales gestionados';
COMMENT ON TABLE acciones_mantenimiento IS 'Acciones de mantenimiento colaborativas (HdU13)';
COMMENT ON TABLE reportes IS 'Reportes automáticos generados (HdU04)';
COMMENT ON TABLE empresas IS 'Empresas de mantención que emplean técnicos';
COMMENT ON TABLE solicitudes_tecnico IS 'Solicitudes de trabajo enviadas a técnicos especializados (HdU16)';
COMMENT ON TABLE archivos_solicitud IS 'Archivos adjuntos a solicitudes técnicas';
COMMENT ON TABLE activos_tecnicos_autorizados IS 'Técnicos autorizados para trabajar en activos específicos';
COMMENT ON TABLE usuarios IS 'Usuarios del sistema con acceso a edificios';
COMMENT ON TABLE usuarios_edificios IS 'Relación muchos a muchos entre usuarios y edificios';
COMMENT ON TABLE fallos IS 'Catálogo de fallos comunes por tipo de activo';
COMMENT ON TABLE activos_fallos IS 'Relación entre activos específicos y fallos detectados';
COMMENT ON TABLE tipos_falla IS 'Reportes de fallas hechos por usuarios residentes en edificios';
COMMENT ON TABLE comentarios IS 'Comentarios de usuarios sobre reportes de fallas';
COMMENT ON TABLE firmas_digitales IS 'Firmas digitales de usuarios para firma de reportes (HdU Firmas Digitales)';

-- Comentarios sobre las tablas creadas
