-- ============================================================================
-- BASE DE DATOS DE NOTIFICACIONES - INSPECCIONAR
-- ============================================================================
-- Este script inicializa el esquema de la base de datos de notificaciones
-- Incluye tablas para edificios, usuarios, activos, notificaciones y tickets
-- ============================================================================

-- ============================================================================
-- TABLAS PRINCIPALES
-- ============================================================================

-- Tabla de edificios
CREATE TABLE edificios (
    id SERIAL PRIMARY KEY,
    direccion VARCHAR(255) NOT NULL,
    numero_activos INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de usuarios
CREATE TABLE usuarios (
    id SERIAL PRIMARY KEY,
    scope VARCHAR(100) NOT NULL,
    usuario VARCHAR(100) UNIQUE NOT NULL,
    correo VARCHAR(255) NOT NULL,
    numero VARCHAR(20),
    edificio_id INTEGER REFERENCES edificios(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de activos
CREATE TABLE activos (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    edificio_id INTEGER NOT NULL REFERENCES edificios(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de notificaciones
CREATE TABLE notificaciones (
    id SERIAL PRIMARY KEY,
    
    -- Referencias
    building_id INTEGER REFERENCES edificios(id) ON DELETE CASCADE,
    asset_id INTEGER REFERENCES activos(id) ON DELETE CASCADE,
    sensor_id VARCHAR(255),
    
    -- Contenido de la notificación
    message TEXT NOT NULL,
    alert_type VARCHAR(100) NOT NULL,
    tipo VARCHAR(50) NOT NULL,
    prioridad VARCHAR(20) NOT NULL DEFAULT 'medium',
    
    -- Canales de notificación
    notification_mail BOOLEAN DEFAULT FALSE,
    notification_sms BOOLEAN DEFAULT FALSE,
    
    -- Estado
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- SISTEMA DE TICKETS
-- ============================================================================

-- Tabla de tickets para solicitudes de ingreso, modificación y eliminación
CREATE TABLE tickets (
    id SERIAL PRIMARY KEY,
    
    -- ========================================================================
    -- CONFIGURACIÓN DEL TICKET
    -- ========================================================================
    
    tipo_entidad VARCHAR(50) NOT NULL 
        CHECK (tipo_entidad IN ('edificio', 'activo', 'tecnico')),
    tipo_operacion VARCHAR(50) NOT NULL 
        CHECK (tipo_operacion IN ('ingreso', 'modificacion', 'eliminacion')),
    estado VARCHAR(50) NOT NULL DEFAULT 'no_resuelto' 
        CHECK (estado IN ('resuelto', 'no_resuelto')),
    
    -- ========================================================================
    -- USUARIO SOLICITANTE
    -- ========================================================================
    
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE SET NULL,
    usuario_email VARCHAR(255) NOT NULL,
    
    -- ========================================================================
    -- DATOS PARA EDIFICIO
    -- ========================================================================
    
    edificio_id INTEGER REFERENCES edificios(id) ON DELETE SET NULL,
    edificio_nombre VARCHAR(255),
    edificio_direccion VARCHAR(255),
    edificio_latitud DECIMAL(10, 8),
    edificio_longitud DECIMAL(11, 8),
    
    -- ========================================================================
    -- DATOS PARA ACTIVO
    -- ========================================================================
    
    activo_id INTEGER REFERENCES activos(id) ON DELETE SET NULL,
    activo_nombre VARCHAR(255),
    activo_tipo VARCHAR(100),
    activo_descripcion TEXT,
    activo_ubicacion VARCHAR(255),
    activo_edificio_id INTEGER REFERENCES edificios(id) ON DELETE SET NULL,
    
    -- ========================================================================
    -- DATOS PARA TÉCNICO
    -- ========================================================================
    
    tecnico_id INTEGER,
    tecnico_nombre VARCHAR(255),
    tecnico_email VARCHAR(255),
    tecnico_telefono VARCHAR(20),
    tecnico_especialidad VARCHAR(100),
    tecnico_autorizado BOOLEAN,
    
    -- ========================================================================
    -- INFORMACIÓN ADICIONAL Y RESOLUCIÓN
    -- ========================================================================
    
    justificacion TEXT,
    comentario_admin TEXT,
    
    -- ========================================================================
    -- METADATA Y TIMESTAMPS
    -- ========================================================================
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    fecha_resolucion TIMESTAMP,
    resuelto_por INTEGER REFERENCES usuarios(id) ON DELETE SET NULL,
    resuelto_por_email VARCHAR(255),
    
    -- ========================================================================
    -- CONSTRAINTS DE VALIDACIÓN
    -- ========================================================================
    
    -- Validación para datos de edificio
    CONSTRAINT chk_edificio_data CHECK (
        (tipo_entidad = 'edificio' AND tipo_operacion = 'ingreso' AND 
         edificio_nombre IS NOT NULL AND 
         edificio_direccion IS NOT NULL AND 
         edificio_latitud IS NOT NULL AND 
         edificio_longitud IS NOT NULL) OR
        (tipo_entidad = 'edificio' AND tipo_operacion IN ('modificacion', 'eliminacion') AND 
         edificio_id IS NOT NULL) OR
        (tipo_entidad IN ('activo', 'tecnico'))
    ),
    
    -- Validación para datos de activo
    CONSTRAINT chk_activo_data CHECK (
        (tipo_entidad = 'activo' AND tipo_operacion = 'ingreso' AND 
         activo_nombre IS NOT NULL AND 
         activo_tipo IS NOT NULL) OR
        (tipo_entidad = 'activo' AND tipo_operacion IN ('modificacion', 'eliminacion') AND 
         activo_id IS NOT NULL) OR
        (tipo_entidad IN ('edificio', 'tecnico'))
    ),
    
    -- Validación para datos de técnico
    CONSTRAINT chk_tecnico_data CHECK (
        (tipo_entidad = 'tecnico' AND tipo_operacion = 'ingreso' AND 
         tecnico_nombre IS NOT NULL AND 
         tecnico_email IS NOT NULL AND 
         tecnico_telefono IS NOT NULL AND 
         tecnico_especialidad IS NOT NULL AND 
         tecnico_autorizado IS NOT NULL) OR
        (tipo_entidad = 'tecnico' AND tipo_operacion IN ('modificacion', 'eliminacion') AND 
         tecnico_id IS NOT NULL) OR
        (tipo_entidad IN ('edificio', 'activo'))
    )
);

-- ============================================================================
-- ÍNDICES PARA OPTIMIZACIÓN
-- ============================================================================

-- Índices para tablas principales
CREATE INDEX idx_usuarios_edificio_id ON usuarios(edificio_id);
CREATE INDEX idx_activos_edificio_id ON activos(edificio_id);

-- Índices para notificaciones
CREATE INDEX idx_notificaciones_building_id ON notificaciones(building_id);
CREATE INDEX idx_notificaciones_asset_id ON notificaciones(asset_id);
CREATE INDEX idx_notificaciones_sensor_id ON notificaciones(sensor_id);
CREATE INDEX idx_notificaciones_tipo ON notificaciones(tipo);
CREATE INDEX idx_notificaciones_prioridad ON notificaciones(prioridad);
CREATE INDEX idx_notificaciones_status ON notificaciones(status);
CREATE INDEX idx_notificaciones_created_at ON notificaciones(created_at);

-- Índices para tickets
CREATE INDEX idx_tickets_tipo_entidad ON tickets(tipo_entidad);
CREATE INDEX idx_tickets_tipo_operacion ON tickets(tipo_operacion);
CREATE INDEX idx_tickets_estado ON tickets(estado);
CREATE INDEX idx_tickets_usuario_id ON tickets(usuario_id);
CREATE INDEX idx_tickets_created_at ON tickets(created_at);
CREATE INDEX idx_tickets_edificio_id ON tickets(edificio_id);
CREATE INDEX idx_tickets_activo_id ON tickets(activo_id);
CREATE INDEX idx_tickets_tecnico_id ON tickets(tecnico_id);
CREATE INDEX idx_tickets_estado_created ON tickets(estado, created_at);

-- ============================================================================
-- FUNCIONES Y TRIGGERS
-- ============================================================================

-- Función para actualizar el contador de activos en edificios
CREATE OR REPLACE FUNCTION update_numero_activos()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE edificios 
        SET numero_activos = numero_activos + 1 
        WHERE id = NEW.edificio_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE edificios 
        SET numero_activos = numero_activos - 1 
        WHERE id = OLD.edificio_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Trigger que actualiza el contador de activos automáticamente
CREATE TRIGGER trigger_update_numero_activos
AFTER INSERT OR DELETE ON activos
FOR EACH ROW
EXECUTE FUNCTION update_numero_activos();

-- ============================================================================
-- COMENTARIOS DE DOCUMENTACIÓN
-- ============================================================================

-- Tablas principales
COMMENT ON TABLE edificios IS 'Edificios monitoreados en el sistema';
COMMENT ON TABLE usuarios IS 'Usuarios del sistema de notificaciones';
COMMENT ON TABLE activos IS 'Activos asociados a edificios';
COMMENT ON TABLE notificaciones IS 'Notificaciones y alertas del sistema';

-- Sistema de tickets
COMMENT ON TABLE tickets IS 'Sistema de tickets para solicitudes de ingreso, modificación y eliminación de edificios, activos y técnicos';

COMMENT ON COLUMN tickets.tipo_entidad IS 'Tipo de entidad del ticket: edificio, activo o tecnico';

COMMENT ON COLUMN tickets.tipo_operacion IS 'Tipo de operación solicitada: ingreso, modificacion o eliminacion';

COMMENT ON COLUMN tickets.estado IS 'Estado actual del ticket: resuelto o no_resuelto';

COMMENT ON COLUMN tickets.justificacion IS 'Razón o justificación detallada de la solicitud';

COMMENT ON COLUMN tickets.comentario_admin IS 'Comentario del administrador al momento de resolver el ticket';

-- Campos de técnico
COMMENT ON COLUMN tickets.tecnico_nombre IS 'Nombre completo del técnico';
COMMENT ON COLUMN tickets.tecnico_email IS 'Correo electrónico de contacto del técnico';
COMMENT ON COLUMN tickets.tecnico_telefono IS 'Número de teléfono de contacto del técnico';
COMMENT ON COLUMN tickets.tecnico_especialidad IS 'Especialidad técnica o área de expertise del técnico';
COMMENT ON COLUMN tickets.tecnico_autorizado IS 'Indica si el técnico está autorizado o no';
