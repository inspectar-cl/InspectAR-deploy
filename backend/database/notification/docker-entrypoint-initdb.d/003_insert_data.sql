-- Insertar datos de ejemplo en edificios
INSERT INTO edificios (direccion, numero_activos) VALUES
('Av. Libertador 1234, CABA', 0),
('Corrientes 5678, CABA', 0),
('Florida 9012, CABA', 0);

-- Insertar usuarios de ejemplo
INSERT INTO usuarios (scope, usuario, correo, numero, edificio_id) VALUES
('admin', 'admin_user', 'admin@inspectarar.com', '+541123456789', 1),
('manager', 'manager_edificio1', 'manager1@inspectarar.com', '+541123456790', 1),
('manager', 'manager_edificio2', 'manager2@inspectarar.com', '+541123456791', 2),
('tech', 'tech_user', 'tech@inspectarar.com', '+541123456792', NULL),
('operator', 'operator1', 'j.a.garcia13m@gmail.com', '+541123456793', 1),
('viewer', 'viewer1', 'example@gmail.com', '+541123456794', 2);

-- Insertar activos de ejemplo
INSERT INTO activos (nombre, edificio_id) VALUES
('HVAC Sistema Principal', 1),
('Caldera Edificio 1', 1),
('Bomba de Agua Principal', 1),
('Sistema Eléctrico', 2),
('Ascensor Principal', 2),
('Sistema Contra Incendios', 3),
('Sensor Temperatura Planta Baja', 1),
('Sensor Humedad Primer Piso', 1),
('Sensor Presión Caldera', 2),
('Monitor Calidad Aire', 2),
('Detector Humo Oficina', 3);

-- Actualizar numero_activos en edificios
UPDATE edificios SET numero_activos = (
    SELECT COUNT(*) FROM activos WHERE edificio_id = edificios.id
);
