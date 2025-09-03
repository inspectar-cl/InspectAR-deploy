-- Crear tabla edificios
CREATE TABLE edificios (
    id SERIAL PRIMARY KEY,
    direccion VARCHAR(255) NOT NULL,
    numero_activos INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Crear tabla usuarios
CREATE TABLE usuarios (
    id SERIAL PRIMARY KEY,
    scope VARCHAR(100) NOT NULL,
    usuario VARCHAR(100) UNIQUE NOT NULL,
    correo VARCHAR(255) UNIQUE NOT NULL,
    numero VARCHAR(20),
    edificio_id INTEGER REFERENCES edificios(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Crear tabla activos
CREATE TABLE activos (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    edificio_id INTEGER NOT NULL REFERENCES edificios(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Crear tabla notificaciones
CREATE TABLE notificaciones (
    id SERIAL PRIMARY KEY,
    building_id INTEGER REFERENCES edificios(id) ON DELETE CASCADE,
    asset_id INTEGER REFERENCES activos(id) ON DELETE CASCADE,
    sensor_id VARCHAR(255), -- ID del sensor que generó la alerta
    message TEXT NOT NULL,
    alert_type VARCHAR(100) NOT NULL, -- sensor_disconnected, sensor_problem, asset_offline, etc.
    tipo VARCHAR(50) NOT NULL, -- sensor, alerta, mantenimiento, sistema, etc.
    prioridad VARCHAR(20) NOT NULL DEFAULT 'medium', -- low, medium, high, critical
    notification_mail BOOLEAN DEFAULT FALSE,
    notification_sms BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'pending', -- pending, sent, failed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Crear índices para mejorar rendimiento
CREATE INDEX idx_usuarios_edificio_id ON usuarios(edificio_id);
CREATE INDEX idx_activos_edificio_id ON activos(edificio_id);
CREATE INDEX idx_notificaciones_building_id ON notificaciones(building_id);
CREATE INDEX idx_notificaciones_asset_id ON notificaciones(asset_id);
CREATE INDEX idx_notificaciones_sensor_id ON notificaciones(sensor_id);
CREATE INDEX idx_notificaciones_tipo ON notificaciones(tipo);
CREATE INDEX idx_notificaciones_prioridad ON notificaciones(prioridad);
CREATE INDEX idx_notificaciones_status ON notificaciones(status);
CREATE INDEX idx_notificaciones_created_at ON notificaciones(created_at);

-- Crear trigger para actualizar numero_activos en edificios
CREATE OR REPLACE FUNCTION update_numero_activos()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE edificios SET numero_activos = numero_activos + 1 WHERE id = NEW.edificio_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE edificios SET numero_activos = numero_activos - 1 WHERE id = OLD.edificio_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
