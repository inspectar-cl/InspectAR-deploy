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
    activo_id INTEGER NOT NULL REFERENCES activos(id) ON DELETE CASCADE,
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE SET NULL,
    mensaje TEXT NOT NULL,
    enviado BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    sent_at TIMESTAMP NULL
);

-- Crear índices para mejorar rendimiento
CREATE INDEX idx_usuarios_edificio_id ON usuarios(edificio_id);
CREATE INDEX idx_activos_edificio_id ON activos(edificio_id);
CREATE INDEX idx_notificaciones_activo_id ON notificaciones(activo_id);
CREATE INDEX idx_notificaciones_usuario_id ON notificaciones(usuario_id);

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
