-- Insertar datos de ejemplo en edificios
INSERT INTO edificios (direccion, numero_activos) VALUES
('Av. Libertador 1234, CABA', 0),
('Corrientes 5678, CABA', 0),
('Florida 9012, CABA', 0);

-- Insertar datos de ejemplo en usuarios
INSERT INTO usuarios (scope, usuario, correo, numero, edificio_id) VALUES
('admin', 'admin_user', 'joytan33334@gmail.com', '+54911234567', 1),
('operator', 'operator1', 'j.a.garcia13m@gmail.com', '+54911234568', 1),
('viewer', 'viewer1', 'example@gmail.com', '+54911234569', 2),
('operator', 'operator2', 'example2@gmail.com', '+54911234570', 2),
('viewer', 'viewer2', 'example3@gmail.com', '+54911234571', 3);

-- Insertar datos de ejemplo en activos
INSERT INTO activos (nombre, edificio_id) VALUES
('Sensor Temperatura Planta Baja', 1),
('Sensor Humedad Primer Piso', 1),
('Sistema HVAC Principal', 1),
('Sensor Presión Caldera', 2),
('Monitor Calidad Aire', 2),
('Detector Humo Oficina', 3);
