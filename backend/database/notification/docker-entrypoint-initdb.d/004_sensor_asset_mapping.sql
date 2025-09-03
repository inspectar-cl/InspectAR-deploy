-- Crear tabla para mapeo de sensores a activos
CREATE TABLE sensor_activo_mapping (
    id SERIAL PRIMARY KEY,
    sensor_id VARCHAR(255) NOT NULL,
    activo_id INTEGER NOT NULL REFERENCES activos(id) ON DELETE CASCADE,
    descripcion TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Crear índices para la nueva tabla
CREATE INDEX idx_sensor_activo_mapping_sensor_id ON sensor_activo_mapping(sensor_id);
CREATE INDEX idx_sensor_activo_mapping_activo_id ON sensor_activo_mapping(activo_id);
CREATE INDEX idx_sensor_activo_mapping_active ON sensor_activo_mapping(is_active);

-- Insertar mapeos de sensores a activos de ejemplo
INSERT INTO sensor_activo_mapping (sensor_id, activo_id, descripcion) VALUES
-- Sensores del HVAC Sistema Principal (activo_id: 1)
('TEMP_HVAC_001', 1, 'Sensor de temperatura del sistema HVAC principal'),
('PRESS_HVAC_001', 1, 'Sensor de presión del sistema HVAC principal'),
('FLOW_HVAC_001', 1, 'Sensor de flujo de aire del sistema HVAC'),

-- Sensores de la Caldera Edificio 1 (activo_id: 2)
('TEMP_CALDERA_001', 2, 'Sensor de temperatura de la caldera'),
('PRESS_CALDERA_001', 2, 'Sensor de presión de la caldera'),
('GAS_CALDERA_001', 2, 'Sensor de detección de gas de la caldera'),

-- Sensores de la Bomba de Agua Principal (activo_id: 3)
('FLOW_AGUA_001', 3, 'Sensor de flujo de la bomba de agua'),
('PRESS_AGUA_001', 3, 'Sensor de presión de la bomba de agua'),
('VIBR_AGUA_001', 3, 'Sensor de vibración de la bomba de agua'),

-- Sensores del Sistema Eléctrico (activo_id: 4)
('VOLT_ELEC_001', 4, 'Sensor de voltaje del sistema eléctrico'),
('AMP_ELEC_001', 4, 'Sensor de amperaje del sistema eléctrico'),
('FREQ_ELEC_001', 4, 'Sensor de frecuencia del sistema eléctrico'),

-- Sensores del Ascensor Principal (activo_id: 5)
('PESO_ASC_001', 5, 'Sensor de peso del ascensor'),
('VEL_ASC_001', 5, 'Sensor de velocidad del ascensor'),
('POSIC_ASC_001', 5, 'Sensor de posición del ascensor'),

-- Sensores del Sistema Contra Incendios (activo_id: 6)
('HUMO_INC_001', 6, 'Sensor de humo del sistema contra incendios'),
('TEMP_INC_001', 6, 'Sensor de temperatura del sistema contra incendios'),
('AGUA_INC_001', 6, 'Sensor de presión de agua del sistema contra incendios'),

-- Sensores adicionales ya existentes en el sistema
('sensor_test_001', 1, 'Sensor de prueba asociado al HVAC'),
('sensor_integration_test', 2, 'Sensor de integración asociado a la caldera'),
('sensor_test_medium', 3, 'Sensor de prueba asociado a la bomba de agua'),
('sensor_test_critical', 4, 'Sensor crítico asociado al sistema eléctrico');

-- Crear vista para consultas rápidas de sensores con información completa
CREATE VIEW v_sensores_completos AS
SELECT 
    sam.sensor_id,
    sam.descripcion as sensor_descripcion,
    a.id as activo_id,
    a.nombre as activo_nombre,
    e.id as edificio_id,
    e.direccion as edificio_direccion,
    sam.is_active as sensor_activo,
    sam.created_at as sensor_mapping_created
FROM sensor_activo_mapping sam
JOIN activos a ON sam.activo_id = a.id
JOIN edificios e ON a.edificio_id = e.id
WHERE sam.is_active = TRUE;

-- Crear función para obtener activo_id de un sensor
CREATE OR REPLACE FUNCTION get_activo_by_sensor(sensor_id_param VARCHAR(255))
RETURNS TABLE(
    activo_id INTEGER,
    activo_nombre VARCHAR(255),
    edificio_id INTEGER,
    edificio_direccion VARCHAR(255)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        a.id as activo_id,
        a.nombre as activo_nombre,
        e.id as edificio_id,
        e.direccion as edificio_direccion
    FROM sensor_activo_mapping sam
    JOIN activos a ON sam.activo_id = a.id
    JOIN edificios e ON a.edificio_id = e.id
    WHERE sam.sensor_id = sensor_id_param 
    AND sam.is_active = TRUE
    LIMIT 1;
END;
$$ LANGUAGE plpgsql;

-- Crear función para obtener usuarios relacionados con un sensor
CREATE OR REPLACE FUNCTION get_usuarios_by_sensor(sensor_id_param VARCHAR(255))
RETURNS TABLE(
    usuario_id INTEGER,
    usuario_nombre VARCHAR(100),
    usuario_correo VARCHAR(255),
    usuario_numero VARCHAR(20),
    usuario_scope VARCHAR(100),
    activo_nombre VARCHAR(255),
    edificio_direccion VARCHAR(255)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.id as usuario_id,
        u.usuario as usuario_nombre,
        u.correo as usuario_correo,
        u.numero as usuario_numero,
        u.scope as usuario_scope,
        a.nombre as activo_nombre,
        e.direccion as edificio_direccion
    FROM sensor_activo_mapping sam
    JOIN activos a ON sam.activo_id = a.id
    JOIN edificios e ON a.edificio_id = e.id
    JOIN usuarios u ON u.edificio_id = e.id
    WHERE sam.sensor_id = sensor_id_param 
    AND sam.is_active = TRUE;
END;
$$ LANGUAGE plpgsql;
