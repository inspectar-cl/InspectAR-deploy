-- Datos de ejemplo para testing
-- Archivo: 002_insert_sample_data.sql

-- Insertar edificios de ejemplo
INSERT INTO edificios (nombre, direccion) VALUES
('Edificio Central', 'Av. Principal 123'),
('Edificio Norte', 'Calle Secundaria 456'),
('Edificio Sur', 'Av. Industrial 789');

-- Insertar técnicos especializados
INSERT INTO tecnicos (nombre, email, telefono, especialidad, autorizado) VALUES
('Juan Pérez', 'juan.perez@inspectar.com', '+56912345678', 'Sistemas Hidráulicos', true),
('María González', 'maria.gonzalez@inspectar.com', '+56987654321', 'Electricidad Industrial', true),
('Carlos Rodríguez', 'carlos.rodriguez@inspectar.com', '+56911223344', 'Mecánica Industrial', false),
('Ana Silva', 'ana.silva@inspectar.com', '+56955667788', 'Instrumentación', true),
('Pedro Morales', 'pedro.morales@inspectar.com', '+56944556677', 'Sistemas de Control', true),
('Laura Fernández', 'laura.fernandez@inspectar.com', '+56933445566', 'Calderas y Vapor', true),
('Roberto Castro', 'roberto.castro@inspectar.com', '+56922334455', 'Refrigeración', false),
('Isabel Herrera', 'isabel.herrera@inspectar.com', '+56911445577', 'Sistemas Neumáticos', true);

-- Insertar activos de ejemplo
INSERT INTO activos (activo_id, nombre, tipo, estado, ubicacion, edificio_id) VALUES
('AC-1001', 'Caldera Principal', 'Caldera', 'operativo', 'Sala de Calderas 1', 1),
('AC-1002', 'Compresor Auxiliar', 'Compresor', 'mantenimiento', 'Sala de Compresores', 1),
('AC-1003', 'Bomba Hidráulica 1', 'Bomba', 'operativo', 'Sala de Bombas', 2),
('AC-1004', 'Motor Eléctrico Principal', 'Motor', 'operativo', 'Sala de Motores', 1),
('AC-1005', 'Sistema de Ventilación', 'Ventilación', 'alerta', 'Techo Edificio 1', 1),
('AC-1006', 'Transformador Eléctrico', 'Transformador', 'operativo', 'Subestación', 2),
('AC-1007', 'Chiller Industrial', 'Refrigeración', 'operativo', 'Sala de Refrigeración', 2),
('AC-1008', 'Generador de Emergencia', 'Generador', 'standby', 'Sala de Generadores', 1);

-- Insertar acciones de mantenimiento
INSERT INTO acciones_mantenimiento (activo_id, tecnico_id, tipo, descripcion, estado, prioridad) VALUES
(1, 1, 'preventivo', 'Revisión mensual de válvulas de seguridad en caldera principal', 'pendiente', 'alta'),
(2, 3, 'correctivo', 'Reparación de compresor auxiliar - falla en motor', 'en_progreso', 'critica'),
(3, 1, 'preventivo', 'Cambio de aceite y filtros en bomba hidráulica', 'completado', 'media'),
(4, 2, 'preventivo', 'Inspección eléctrica de motor principal', 'pendiente', 'media'),
(5, 2, 'correctivo', 'Sistema de ventilación con ruido anormal', 'pendiente', 'alta'),
(6, 2, 'preventivo', 'Medición de aislamiento en transformador', 'completado', 'baja'),
(7, 7, 'preventivo', 'Limpieza de serpentines en chiller industrial', 'pendiente', 'media'),
(8, 5, 'preventivo', 'Prueba semanal de arranque de generador', 'en_progreso', 'alta'),
(1, 6, 'correctivo', 'Caldera con presión irregular - revisar sensores', 'pendiente', 'critica'),
(3, 4, 'preventivo', 'Calibración de instrumentos de bomba hidráulica', 'pendiente', 'baja');

-- Insertar reportes de ejemplo
INSERT INTO reportes (activo_id, tipo_reporte, contenido, estado) VALUES
(1, 'semanal', 'Reporte semanal: Caldera operando normalmente. Presión estable en 8.5 bar. Temperatura promedio: 85°C.', 'enviado'),
(2, 'incidente', 'Incidente: Compresor auxiliar fuera de servicio por falla en motor. Técnico asignado para reparación.', 'generado'),
(3, 'mensual', 'Reporte mensual: Bomba hidráulica con rendimiento óptimo. Última mantención realizada exitosamente.', 'enviado'),
(4, 'semanal', 'Reporte semanal: Motor eléctrico principal operando en parámetros normales. Sin anomalías detectadas.', 'archivado'),
(5, 'incidente', 'Incidente: Sistema de ventilación generando ruido anormal. Requiere inspección urgente.', 'generado'),
(6, 'mensual', 'Reporte mensual: Transformador eléctrico con mediciones de aislamiento dentro de especificaciones.', 'enviado'),
(7, 'mantenimiento', 'Mantenimiento programado: Chiller industrial requiere limpieza de serpentines la próxima semana.', 'generado'),
(8, 'semanal', 'Reporte semanal: Generador de emergencia probado exitosamente. Tiempo de arranque: 12 segundos.', 'enviado');

-- Insertar relaciones entre activos y técnicos
INSERT INTO activos_tecnicos (activo_id, tecnico_id) VALUES
(1, 1), -- Juan Pérez asignado a Caldera Principal
(1, 6), -- Laura Fernández asignada a Caldera Principal
(2, 3), -- Carlos Rodríguez asignado a Compresor Auxiliar
(3, 1), -- Juan Pérez asignado a Bomba Hidráulica 1
(3, 4), -- Ana Silva asignada a Bomba Hidráulica 1
(4, 2), -- María González asignada a Motor Eléctrico Principal
(5, 2), -- María González asignada a Sistema de Ventilación
(6, 2), -- María González asignada a Transformador Eléctrico
(7, 7), -- Roberto Castro asignado a Chiller Industrial
(8, 5); -- Pedro Morales asignado a Generador de Emergencia

-- Actualizar algunas fechas para tener datos más realistas
UPDATE acciones_mantenimiento 
SET fecha_fin = fecha_inicio + INTERVAL '2 hours'
WHERE estado = 'completado';

UPDATE acciones_mantenimiento 
SET fecha_inicio = CURRENT_TIMESTAMP - INTERVAL '1 day'
WHERE estado = 'en_progreso';

UPDATE reportes 
SET generado_en = CURRENT_TIMESTAMP - INTERVAL '7 days'
WHERE tipo_reporte = 'semanal';

UPDATE reportes 
SET generado_en = CURRENT_TIMESTAMP - INTERVAL '30 days'
WHERE tipo_reporte = 'mensual';
