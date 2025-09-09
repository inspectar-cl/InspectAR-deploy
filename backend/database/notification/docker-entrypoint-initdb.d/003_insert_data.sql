-- Insertar edificios sincronizados con Gestión DB
INSERT INTO edificios (direccion, numero_activos) VALUES
('Av. Providencia 123, Santiago', 0),
('Av. Las Condes 456, Las Condes', 0),
('Av. Vicuña Mackenna 789, La Florida', 0),
('Ruta 68 Km 15, Melipilla', 0);

-- Insertar usuarios de ejemplo
INSERT INTO usuarios (scope, usuario, correo, numero, edificio_id) VALUES
('admin', 'admin_user', 'admin@inspectarar.com', '+56912345678', 1),
('manager', 'manager_edificio1', 'manager1@inspectarar.com', '+56912345679', 1),
('manager', 'manager_edificio2', 'manager2@inspectarar.com', '+56912345680', 2),
('manager', 'manager_edificio3', 'manager3@inspectarar.com', '+56912345681', 3),
('manager', 'manager_edificio4', 'manager4@inspectarar.com', '+56912345682', 4),
('tech', 'tech_user', 'tech@inspectarar.com', '+56912345683', NULL);

-- Insertar activos sincronizados con Gestión DB (mismo orden, nombres y edificios)
INSERT INTO activos (nombre, edificio_id) VALUES
('BombaDeAgua #1 ML', 1),              -- ID 1 -> AC-1001
('Bomba Centrífuga A', 1),             -- ID 2 -> AC-1002
('Bomba Hidráulica 1', 2),             -- ID 3 -> AC-1003
('Ascensor Norte', 2),                 -- ID 4 -> AC-1004
('Transformador Secundario', 3),       -- ID 5 -> AC-1005
('Transformador Principal', 3),        -- ID 6 -> AC-1006
('Ascensor Central', 1),               -- ID 7 -> AC-1007
('Bomba de Emergencia', 4);            -- ID 8 -> AC-1008

-- Actualizar numero_activos en edificios
UPDATE edificios SET numero_activos = (
    SELECT COUNT(*) FROM activos WHERE edificio_id = edificios.id
);
