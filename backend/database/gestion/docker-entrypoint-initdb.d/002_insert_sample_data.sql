-- Datos de ejemplo para testing
-- Archivo: 002_insert_sample_data.sql

-- Insertar empresas de mantención
INSERT INTO empresas (nombre, rut, telefono, email, direccion) VALUES
('Mantención Integral SpA', '12345678-9', '+56912345678', 'contacto@mantencion.cl', 'Av. Providencia 1234, Santiago'),
('Servicios Técnicos Ltda', '87654321-0', '+56987654321', 'info@servicios.cl', 'Av. Las Condes 567, Las Condes'),
('Especialistas Industriales', '11223344-5', '+56911223344', 'ventas@especialistas.cl', 'Av. Vitacura 890, Vitacura')
ON CONFLICT (rut) DO NOTHING;

-- Insertar edificios
INSERT INTO edificios (nombre, direccion) VALUES
('Edificio Central', 'Av. Providencia 123, Santiago'),
('Torre Norte', 'Av. Las Condes 456, Las Condes'),
('Complejo Industrial Sur', 'Av. Vicuña Mackenna 789, La Florida'),
('Centro de Distribución', 'Ruta 68 Km 15, Melipilla')
ON CONFLICT DO NOTHING;

-- Insertar usuarios del sistema (basados en los usuarios de OAuth2)
INSERT INTO usuarios (username, email) VALUES
('analista', 'analista@example.com'),
('tecnico', 'tecnico@example.com'),
('residente', 'residente@example.com'),
('admin', 'admin@example.com')
ON CONFLICT (username) DO NOTHING;

-- Crear relaciones usuarios-edificios
-- Los primeros 3 usuarios tienen acceso a todos los edificios
INSERT INTO usuarios_edificios (usuario_id, edificio_id) VALUES
-- admin tiene acceso a todos los edificios
(1, 1), (1, 2), (1, 3), (1, 4),
-- tecnico tiene acceso a todos los edificios
(2, 1), (2, 2), (2, 3), (2, 4),
-- usuario tiene acceso a todos los edificios
(3, 1), (3, 2), (3, 3), (3, 4),
-- residente_especial solo tiene acceso al edificio 1
(4, 1)
ON CONFLICT DO NOTHING;

-- Insertar técnicos especializados con empresas
INSERT INTO tecnicos (nombre, apellido, email, telefono, especialidad, empresa_id, autorizado, activo) VALUES
('Juan', 'Pérez', 'juan.perez@mantencion.cl', '+56912345678', 'Sistemas Hidráulicos', 1, true, true),
('María', 'González', 'maria.gonzalez@servicios.cl', '+56987654321', 'Electricidad Industrial', 2, true, true),
('Carlos', 'López', 'carlos.lopez@especialistas.cl', '+56955566677', 'Sistemas HVAC', 3, true, true),
('Ana', 'Martínez', 'ana.martinez@mantencion.cl', '+56944455566', 'Instrumentación', 1, true, true),
('Pedro', 'Silva', 'pedro.silva@servicios.cl', '+56933344455', 'Mecánica Industrial', 2, true, true),
('Laura', 'Fernández', 'laura.fernandez@especialistas.cl', '+56922233344', 'Calderas y Vapor', 3, true, true),
('Roberto', 'Morales', 'roberto.morales@mantencion.cl', '+56911122233', 'Sistemas de Control', 1, true, true),
('Carmen', 'Rojas', 'carmen.rojas@servicios.cl', '+56999988877', 'Seguridad Industrial', 2, true, true)
ON CONFLICT (email) DO NOTHING;

-- Insertar activos industriales (solo tipos permitidos: caldera, bomba de agua, ascensor, transformador)
INSERT INTO activos (nombre, tipo, estado, ubicacion, edificio_id) VALUES
('BombaDeAgua #1 ML', 'bomba de agua', 'operativo', 'Planta B', 1),
('Bomba Centrífuga A', 'bomba de agua', 'operativo', 'Sala de Bombas', 1),
('Bomba Hidráulica 1', 'bomba de agua', 'operativo', 'Sala de Bombas', 2),
('Ascensor Norte', 'ascensor', 'operativo', 'Torre Norte - Piso 1', 2),
('Transformador Secundario', 'transformador', 'mantenimiento', 'Subestación Secundaria', 3),
('Transformador Principal', 'transformador', 'operativo', 'Subestación Eléctrica', 3),
('Ascensor Central', 'ascensor', 'operativo', 'Edificio Central - Hall', 1),
('Bomba de Emergencia', 'bomba de agua', 'stand_by', 'Planta de Emergencia', 4)
ON CONFLICT (id) DO NOTHING;

-- Autorizar técnicos para activos específicos
INSERT INTO activos_tecnicos_autorizados (tecnico_id, activo_id, edificio_id) VALUES
-- Juan (Sistemas Hidráulicos) - Calderas y Bombas
(1, 1, 1), (1, 2, 1), (1, 3, 2), (1, 8, 4),
-- María (Electricidad) - Transformadores y Sistemas de Control
(2, 6, 3), (2, 7, 4),
-- Carlos (HVAC) - Sistemas HVAC
(3, 4, 2),
-- Ana (Instrumentación) - Sistemas de Control
(4, 7, 4),
-- Pedro (Mecánica) - Compresores y Bombas
(5, 5, 3), (5, 2, 1), (5, 3, 2),
-- Laura (Calderas) - Calderas y Vapor
(6, 1, 1),
-- Roberto (Control) - Sistemas de Control y PLC
(7, 7, 4),
-- Carmen (Seguridad) - Todos los activos para inspecciones de seguridad
(8, 1, 1), (8, 2, 1), (8, 3, 2), (8, 4, 2), (8, 5, 3), (8, 6, 3), (8, 7, 4), (8, 8, 4)
ON CONFLICT (tecnico_id, activo_id) DO NOTHING;

-- Insertar solicitudes de muestra
INSERT INTO solicitudes_tecnico (
    tecnico_id, residente_id, activo_id, edificio_id, tipo, asunto, descripcion, 
    prioridad, estado, medio_contacto, telefono_contacto, email_contacto
) VALUES
(1, 1001, 1, 1, 'mantenimiento', 'Revisión mensual de caldera', 'Solicito revisión programada de la caldera principal según protocolo de mantenimiento preventivo', 'media', 'pendiente', 'email', '+56912345678', 'residente1@edificio.cl'),
(2, 1002, 6, 3, 'emergencia', 'Falla eléctrica en transformador', 'Se detectó un ruido anormal en el transformador principal, requiere inspección urgente', 'critica', 'enviada', 'ambos', '+56987654321', 'residente2@edificio.cl'),
(3, 1003, 4, 2, 'reparacion', 'Sistema HVAC no enfría', 'El sistema de aire acondicionado de la Torre Norte no está enfriando correctamente', 'alta', 'en_proceso', 'telefono', '+56955566677', 'residente3@edificio.cl'),
(4, 1004, 7, 4, 'inspeccion', 'Calibración de instrumentos', 'Requiero calibración trimestral de los instrumentos del sistema de control', 'media', 'pendiente', 'email', '', 'residente4@edificio.cl'),
(8, 1005, 5, 3, 'consulta', 'Procedimiento de seguridad', 'Consulta sobre el procedimiento de seguridad para el mantenimiento del compresor industrial', 'baja', 'completada', 'email', '', 'residente5@edificio.cl')
ON CONFLICT DO NOTHING;

-- Insertar acciones de mantenimiento colaborativas (compatibilidad con HdU13)
INSERT INTO acciones_mantenimiento (activo_id, tecnico_id, tipo, descripcion, estado, prioridad) VALUES
(1, 1, 'preventivo', 'Inspección mensual de calderas', 'pendiente', 'alta'),
(2, 1, 'correctivo', 'Reparación de bomba centrífuga', 'en_progreso', 'media'),
(3, 5, 'preventivo', 'Mantenimiento bomba hidráulica', 'pendiente', 'media'),
(4, 3, 'correctivo', 'Reparación sistema HVAC', 'completado', 'alta'),
(5, 5, 'emergencia', 'Reparación urgente compresor', 'pendiente', 'critica'),
(6, 2, 'preventivo', 'Revisión transformador', 'pendiente', 'alta'),
(7, 7, 'preventivo', 'Actualización sistema control', 'en_progreso', 'media'),
(8, 1, 'preventivo', 'Prueba bomba emergencia', 'pendiente', 'baja')
ON CONFLICT DO NOTHING;

-- Insertar reportes automáticos de muestra
INSERT INTO reportes (activo_id, tipo_reporte, contenido, estado) VALUES
(1, 'mantenimiento', 'Reporte mensual de caldera - Estado operativo normal', 'generado'),
(2, 'incidente', 'Reporte de falla en bomba centrífuga - Reparación completada', 'enviado'),
(3, 'semanal', 'Reporte semanal bomba hidráulica - Funcionamiento óptimo', 'generado'),
(4, 'mantenimiento', 'Reporte mantenimiento HVAC - Sistema reparado', 'archivado'),
(5, 'incidente', 'Reporte emergencia compresor - Pendiente reparación', 'generado')
ON CONFLICT DO NOTHING;

-- Insertar relaciones de compatibilidad (activos_tecnicos)
INSERT INTO activos_tecnicos (activo_id, tecnico_id) VALUES
(1, 1), (1, 6),  -- Caldera: Juan y Laura
(2, 1), (2, 5),  -- Bomba A: Juan y Pedro  
(3, 1), (3, 5),  -- Bomba 1: Juan y Pedro
(4, 3),          -- HVAC: Carlos
(5, 5),          -- Compresor: Pedro
(6, 2),          -- Transformador: María
(7, 4), (7, 7),  -- Control: Ana y Roberto
(8, 1)           -- Bomba Emergencia: Juan
ON CONFLICT DO NOTHING;

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

-- 🆕 INSERTAR REPORTES CON OBSERVACIONES EDITABLES (DATOS DE MUESTRA)
INSERT INTO reportes (
    activo_id, tipo_reporte, contenido, 
    observaciones_analista, autor_analista, 
    estructura_informe, metadata_informe, 
    version_reporte, estado_revision, 
    fecha_revision, revisor, 
    estado
) VALUES
(1, 'mantenimiento', 'Reporte detallado de mantenimiento preventivo de caldera principal',
 'Se observó un ligero incremento en la temperatura de funcionamiento. Se recomienda revisar el sistema de refrigeración en la próxima mantención.',
 'Juan Pérez',
 '{"resumen": "Mantenimiento preventivo caldera principal", "observaciones": "Ligero incremento de temperatura", "recomendaciones": ["Revisar sistema refrigeración", "Verificar sensores temperatura"], "conclusiones": "Equipo operativo con observaciones menores"}',
 '{"fecha_creacion": "2025-09-08T10:00:00Z", "tipo_activo": "caldera", "nombre_activo": "Caldera Principal", "version": 1}',
 1, 'aprobado', CURRENT_TIMESTAMP - INTERVAL '2 hours', 'María González', 'enviado'),

(3, 'incidente', 'Reporte de incidente menor en bomba hidráulica',
 'Durante la inspección rutinaria se detectó una pequeña fuga en el sello mecánico. Se programó reemplazo para el próximo mantenimiento programado.',
 'Carlos López',
 '{"resumen": "Incidente menor bomba hidráulica", "observaciones": "Fuga pequeña en sello mecánico", "recomendaciones": ["Reemplazar sello mecánico", "Verificar torque de conexiones"], "conclusiones": "Incidente menor, equipo puede continuar operando"}',
 '{"fecha_creacion": "2025-09-07T14:30:00Z", "tipo_activo": "bomba de agua", "nombre_activo": "Bomba Hidráulica 1", "version": 1}',
 1, 'en_revision', NULL, NULL, 'generado'),

(6, 'semanal', 'Reporte semanal de monitoreo transformador',
 'Transformador funcionando dentro de parámetros normales. Se registró una ligera variación en los niveles de aceite dieléctrico que requiere seguimiento.',
 'Ana Martínez',
 '{"resumen": "Monitoreo semanal transformador", "observaciones": "Variación en niveles aceite dieléctrico", "recomendaciones": ["Monitorear niveles aceite", "Programar análisis de aceite"], "conclusiones": "Funcionamiento normal con observación preventiva"}',
 '{"fecha_creacion": "2025-09-06T09:15:00Z", "tipo_activo": "transformador", "nombre_activo": "Transformador Principal", "version": 1}',
 1, 'pendiente', NULL, NULL, 'generado');

-- 🆕 INSERTAR REPORTES ADICIONALES PARA TESTING COMPLETO
INSERT INTO reportes (
    activo_id, tipo_reporte, contenido,
    observaciones_analista, autor_analista,
    version_reporte, estado_revision,
    estado
) VALUES
(4, 'mantenimiento', 'Reporte de mantenimiento de ascensor',
 'Ascensor operando correctamente. Se realizó lubricación de guías y ajuste de puertas.',
 'Pedro Silva',
 1, 'aprobado', 'enviado'),

(2, 'semanal', 'Reporte semanal bomba centrífuga',
 'Bomba presentando vibración ligeramente superior a lo normal. Requiere balanceado.',
 'Laura Fernández',
 1, 'pendiente', 'generado');

-- -- Insertar relaciones entre activos y técnicos
-- INSERT INTO activos_tecnicos (activo_id, tecnico_id) VALUES
-- (1, 1), -- Juan Pérez asignado a Caldera Principal
-- (1, 6), -- Laura Fernández asignada a Caldera Principal
-- (2, 3), -- Carlos Rodríguez asignado a Compresor Auxiliar
-- (3, 1), -- Juan Pérez asignado a Bomba Hidráulica 1
-- (3, 4), -- Ana Silva asignada a Bomba Hidráulica 1
-- (4, 2), -- María González asignada a Motor Eléctrico Principal
-- (5, 2), -- María González asignada a Sistema de Ventilación
-- (6, 2), -- María González asignada a Transformador Eléctrico
-- (7, 7), -- Roberto Castro asignado a Chiller Industrial
-- (8, 5); -- Pedro Morales asignado a Generador de Emergencia

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

-- Insertar catálogo de fallos por tipo de activo
INSERT INTO fallos (tipo_activo, descripcion, prioridad, probabilidad_ocurrencia) VALUES
-- Fallos de CALDERA
('caldera', 'Fuga de agua en tubería principal', 'alta', 15.50),
('caldera', 'Presión irregular o fluctuante', 'media', 25.00),
('caldera', 'Fallo en válvula de seguridad', 'alta', 8.75),
('caldera', 'Corrosión en intercambiador de calor', 'media', 18.20),
('caldera', 'Obstrucción en sistema de combustión', 'alta', 12.30),
('caldera', 'Sensor de temperatura defectuoso', 'media', 22.40),
('caldera', 'Ruido anormal durante operación', 'media', 20.15),
('caldera', 'Pérdida de eficiencia energética', 'media', 28.60),

-- Fallos de BOMBA DE AGUA
('bomba_de_agua', 'Fuga en sello mecánico', 'alta', 32.50),
('bomba_de_agua', 'Cavitación en impulsores', 'media', 24.80),
('bomba_de_agua', 'Vibración excesiva', 'alta', 19.40),
('bomba_de_agua', 'Sobrecalentamiento del motor', 'alta', 16.70),
('bomba_de_agua', 'Desgaste de rodamientos', 'media', 35.20),
('bomba_de_agua', 'Bajo caudal o presión', 'media', 27.90),
('bomba_de_agua', 'Ruido anormal en operación', 'media', 21.50),
('bomba_de_agua', 'Fuga en conexiones de tuberías', 'alta', 18.30),
('bomba_de_agua', 'Fallo en sistema de control', 'media', 14.60),

-- Fallos de ASCENSOR
('ascensor', 'Fallo en puertas (no cierran correctamente)', 'alta', 38.70),
('ascensor', 'Ruido excesivo durante movimiento', 'media', 29.40),
('ascensor', 'Paradas bruscas o irregulares', 'alta', 15.80),
('ascensor', 'Fallo en botones o panel de control', 'media', 26.50),
('ascensor', 'Iluminación defectuosa en cabina', 'media', 22.10),
('ascensor', 'Desnivelación en paradas', 'alta', 17.30),
('ascensor', 'Fallo en sistema de emergencia', 'alta', 6.20),
('ascensor', 'Desgaste en cables de tracción', 'alta', 11.40),
('ascensor', 'Sobrecalentamiento del motor', 'alta', 9.80),

-- Fallos de TRANSFORMADOR
('transformador', 'Sobrecalentamiento del núcleo', 'alta', 14.20),
('transformador', 'Pérdida de aceite dieléctrico', 'alta', 18.50),
('transformador', 'Ruido anormal de operación', 'media', 24.30),
('transformador', 'Cortocircuito en bobinados', 'alta', 7.60),
('transformador', 'Degradación del aislamiento', 'alta', 12.90),
('transformador', 'Vibración excesiva', 'media', 19.70),
('transformador', 'Fallo en sistema de refrigeración', 'alta', 16.40),
('transformador', 'Desbalance de voltaje en fases', 'alta', 11.80),
('transformador', 'Corrosión en terminales', 'media', 22.60)
ON CONFLICT DO NOTHING;

-- Insertar fallos detectados en activos específicos (historial)
INSERT INTO activos_fallos (activo_id, fallo_id, estado, notas) VALUES
-- Activo 1: BombaDeAgua #1 ML
(1, 9, 'resuelto', 'Fuga detectada en sello mecánico durante inspección rutinaria. Reemplazado el 15/09/2025.'),
(1, 13, 'en_revision', 'Desgaste leve en rodamientos. Programado para mantenimiento preventivo.'),

-- Activo 2: Bomba Centrífuga A
(2, 11, 'resuelto', 'Vibración excesiva corregida mediante balanceo de impulsores.'),
(2, 15, 'detectado', 'Se detectó bajo caudal. Requiere inspección de filtros y válvulas.'),

-- Activo 3: Bomba Hidráulica 1
(3, 10, 'en_revision', 'Cavitación detectada. Verificando presión de succión.'),

-- Activo 4: Ascensor Norte
(4, 19, 'resuelto', 'Fallo en puertas corregido mediante ajuste de sensores.'),
(4, 20, 'detectado', 'Ruido leve durante movimiento. Programada lubricación de guías.'),

-- Activo 5: Transformador Secundario
(5, 30, 'en_revision', 'Sobrecalentamiento detectado. Verificando sistema de refrigeración.'),

-- Activo 6: Transformador Principal
(6, 31, 'resuelto', 'Pérdida menor de aceite dieléctrico. Nivel restaurado y sellado verificado.'),

-- Activo 7: Ascensor Central
(7, 22, 'detectado', 'Paradas con ligera irregularidad. Requiere calibración del sistema de control.'),

-- Activo 8: Bomba de Emergencia
(8, 14, 'resuelto', 'Verificación de funcionamiento completada. Sin fallos detectados.')
ON CONFLICT DO NOTHING;
