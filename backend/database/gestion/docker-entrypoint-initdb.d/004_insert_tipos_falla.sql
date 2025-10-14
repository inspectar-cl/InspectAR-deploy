-- ============================================================================
-- INSERTAR REPORTES DE FALLAS DE USUARIOS (tipos_falla) Y COMENTARIOS - HU22
-- Archivo: 004_insert_tipos_falla.sql
-- ============================================================================

-- Tipos de falla: 'falla agua', 'falla ascensor', 'falla electricidad', 'falla caldera'
-- Estados: 'reportado', 'en_revision', 'resuelto', 'rechazado'

INSERT INTO tipos_falla (tipo, descripcion, id_usuario, id_edificio, estado, fecha_publicacion) VALUES
-- Fallas en Edificio Central (id: 1)
('falla agua', 'Fuga de agua en el baño del departamento 305. El agua sale por debajo del lavamanos.', 3, 1, 'reportado', CURRENT_TIMESTAMP - INTERVAL '2 hours'),
('falla ascensor', 'El ascensor central se detuvo entre el piso 4 y 5. Los residentes tuvieron que usar las escaleras.', 4, 1, 'en_revision', CURRENT_TIMESTAMP - INTERVAL '5 hours'),
('falla electricidad', 'Corte de luz en todo el tercer piso. Los enchufes no funcionan desde esta mañana.', 1, 1, 'reportado', CURRENT_TIMESTAMP - INTERVAL '1 hour'),
('falla caldera', 'No hay agua caliente en los departamentos del segundo piso desde ayer. La caldera parece no estar funcionando.', 2, 1, 'en_revision', CURRENT_TIMESTAMP - INTERVAL '1 day'),
('falla agua', 'Presión muy baja del agua en el departamento 801. Apenas sale un hilo de agua.', 4, 1, 'resuelto', CURRENT_TIMESTAMP - INTERVAL '3 days'),

-- Fallas en Torre Norte (id: 2)
('falla ascensor', 'El ascensor norte hace ruidos extraños y se mueve con sacudidas. Da miedo usarlo.', 3, 2, 'reportado', CURRENT_TIMESTAMP - INTERVAL '30 minutes'),
('falla electricidad', 'Las luces del pasillo del piso 10 parpadean constantemente. Puede causar problemas.', 1, 2, 'reportado', CURRENT_TIMESTAMP - INTERVAL '6 hours'),
('falla agua', 'Goteo constante en las tuberías del estacionamiento subterráneo. Se está formando un charco.', 2, 2, 'en_revision', CURRENT_TIMESTAMP - INTERVAL '12 hours'),

-- Fallas en Complejo Industrial Sur (id: 3)
('falla electricidad', 'Sobrecarga eléctrica en el sector de producción. Se han quemado 2 fusibles esta semana.', 1, 3, 'en_revision', CURRENT_TIMESTAMP - INTERVAL '8 hours'),
('falla caldera', 'La caldera industrial está generando mucho humo y tiene un olor extraño.', 2, 3, 'reportado', CURRENT_TIMESTAMP - INTERVAL '4 hours'),
('falla agua', 'Las bombas de agua no están funcionando correctamente. La presión es insuficiente.', 3, 3, 'resuelto', CURRENT_TIMESTAMP - INTERVAL '2 days'),

-- Fallas en Centro de Distribución (id: 4)
('falla ascensor', 'El montacargas se quedó atascado con mercancía dentro. Necesitamos ayuda urgente.', 4, 4, 'resuelto', CURRENT_TIMESTAMP - INTERVAL '1 day'),
('falla electricidad', 'El generador de emergencia no arranca. Esto es un problema de seguridad crítico.', 1, 4, 'reportado', CURRENT_TIMESTAMP - INTERVAL '3 hours'),
('falla agua', 'Fuga importante en la sala de baños del personal. El piso está inundado.', 2, 4, 'en_revision', CURRENT_TIMESTAMP - INTERVAL '7 hours'),
('falla caldera', 'Sistema de calefacción no funciona. Hace mucho frío en las oficinas.', 3, 4, 'reportado', CURRENT_TIMESTAMP - INTERVAL '10 hours')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- INSERTAR COMENTARIOS SOBRE REPORTES DE FALLAS - HU22
-- ============================================================================

-- Comentarios de diferentes usuarios sobre las fallas reportadas
INSERT INTO comentarios (id_falla, id_usuario, comentario, fecha_comentario) VALUES
-- Comentarios sobre falla #1 (fuga agua edificio 1)
(1, 1, 'Ya envié al técnico de plomería. Debería llegar en 30 minutos.', CURRENT_TIMESTAMP - INTERVAL '1 hour 50 minutes'),
(1, 3, 'Gracias por la respuesta rápida. Esperaré al técnico.', CURRENT_TIMESTAMP - INTERVAL '1 hour 45 minutes'),
(1, 2, 'El técnico llegó y está trabajando en la reparación.', CURRENT_TIMESTAMP - INTERVAL '1 hour 30 minutes'),

-- Comentarios sobre falla #2 (ascensor atascado)
(2, 2, 'Técnico de ascensores notificado. Llegará en 1 hora. Mientras tanto usar escaleras.', CURRENT_TIMESTAMP - INTERVAL '4 hours 30 minutes'),
(2, 4, 'Entendido. ¿Tienen estimación de cuánto tardará la reparación?', CURRENT_TIMESTAMP - INTERVAL '4 hours'),
(2, 2, 'El técnico está revisando. Estima 2-3 horas para solucionar el problema.', CURRENT_TIMESTAMP - INTERVAL '3 hours'),

-- Comentarios sobre falla #3 (corte electricidad)
(3, 2, 'Electricista en camino. Por favor desconecten todos los electrodomésticos del tercer piso.', CURRENT_TIMESTAMP - INTERVAL '45 minutes'),
(3, 1, 'Ya avisé a los residentes. Todos desconectaron sus equipos.', CURRENT_TIMESTAMP - INTERVAL '30 minutes'),

-- Comentarios sobre falla #4 (caldera sin agua caliente)
(4, 1, 'Técnico de calderas revisando el sistema ahora mismo.', CURRENT_TIMESTAMP - INTERVAL '23 hours'),
(4, 2, 'Se detectó falla en el termostato. Reemplazando la pieza.', CURRENT_TIMESTAMP - INTERVAL '22 hours'),
(4, 1, 'Trabajo completado. El agua caliente debería volver en 1 hora.', CURRENT_TIMESTAMP - INTERVAL '21 hours'),
(4, 4, '¿Ya deberíamos tener agua caliente? En mi departamento aún sale fría.', CURRENT_TIMESTAMP - INTERVAL '20 hours'),
(4, 2, 'Denle tiempo. La caldera debe calentar todo el sistema. En 30 min más debería estar ok.', CURRENT_TIMESTAMP - INTERVAL '19 hours 30 minutes'),

-- Comentarios sobre falla #5 (presión baja - RESUELTA)
(5, 2, 'Se revisó y limpió el filtro de la tubería. Presión normalizada.', CURRENT_TIMESTAMP - INTERVAL '3 days'),
(5, 4, 'Confirmado. La presión ya está normal. Muchas gracias!', CURRENT_TIMESTAMP - INTERVAL '2 days 23 hours'),

-- Comentarios sobre falla #6 (ascensor con ruidos Torre Norte)
(6, 1, 'Urgente. Revisaremos hoy en la tarde. Por favor no usar ese ascensor.', CURRENT_TIMESTAMP - INTERVAL '20 minutes'),
(6, 3, 'Entendido. Usaré el otro ascensor. ¿Hay riesgo de accidente?', CURRENT_TIMESTAMP - INTERVAL '15 minutes'),
(6, 1, 'Por precaución lo deshabilitamos hasta que el técnico lo revise completamente.', CURRENT_TIMESTAMP - INTERVAL '10 minutes'),

-- Comentarios sobre falla #7 (luces parpadeando)
(7, 2, 'Electricista programado para mañana en la mañana. No es peligroso pero lo arreglaremos.', CURRENT_TIMESTAMP - INTERVAL '5 hours'),

-- Comentarios sobre falla #8 (goteo estacionamiento)
(8, 1, 'Técnico identificó la tubería. Necesita reemplazo de una sección.', CURRENT_TIMESTAMP - INTERVAL '11 hours'),
(8, 2, 'Materiales pedidos. Reparación se hará el miércoles.', CURRENT_TIMESTAMP - INTERVAL '10 hours'),

-- Comentarios sobre falla #9 (sobrecarga eléctrica industrial)
(9, 2, 'Electricista industrial revisando la carga. Posible necesidad de ampliar capacidad.', CURRENT_TIMESTAMP - INTERVAL '7 hours'),
(9, 1, 'Urgente resolver esto. Afecta la producción.', CURRENT_TIMESTAMP - INTERVAL '6 hours 30 minutes'),

-- Comentarios sobre falla #10 (caldera industrial)
(10, 1, 'Técnico especialista en camino. ETA 45 minutos.', CURRENT_TIMESTAMP - INTERVAL '3 hours 30 minutes'),
(10, 2, 'Por favor evacuen el área cerca de la caldera por seguridad.', CURRENT_TIMESTAMP - INTERVAL '3 hours 15 minutes'),

-- Comentarios sobre falla #11 (bombas agua - RESUELTA)
(11, 1, 'Se reemplazó el impulsor de la bomba principal. Todo funcionando.', CURRENT_TIMESTAMP - INTERVAL '2 days'),
(11, 3, 'Excelente. La presión ya está normal en toda la planta.', CURRENT_TIMESTAMP - INTERVAL '1 day 23 hours'),

-- Comentarios sobre falla #12 (montacargas atascado - RESUELTO)
(12, 2, 'Montacargas liberado. Mercancía recuperada sin daños.', CURRENT_TIMESTAMP - INTERVAL '23 hours'),
(12, 4, 'Perfecto. Gracias por la pronta respuesta.', CURRENT_TIMESTAMP - INTERVAL '22 hours'),

-- Comentarios sobre falla #13 (generador emergencia)
(13, 1, 'Técnico especializado en generadores notificado. Llegará en 2 horas.', CURRENT_TIMESTAMP - INTERVAL '2 hours 30 minutes'),
(13, 2, 'Esto es crítico para la seguridad. Necesitamos prioridad.', CURRENT_TIMESTAMP - INTERVAL '2 hours 15 minutes'),
(13, 1, 'Confirmado como prioridad alta. Técnico acelerando su llegada.', CURRENT_TIMESTAMP - INTERVAL '2 hours'),

-- Comentarios sobre falla #14 (inundación baños)
(14, 2, 'Plomero de emergencia en camino. Cierren la llave de paso principal.', CURRENT_TIMESTAMP - INTERVAL '6 hours 30 minutes'),
(14, 3, 'Llave cerrada. El agua dejó de salir pero el piso está muy mojado.', CURRENT_TIMESTAMP - INTERVAL '6 hours 15 minutes'),
(14, 2, 'Plomero llegó. También enviamos personal de limpieza.', CURRENT_TIMESTAMP - INTERVAL '6 hours'),

-- Comentarios sobre falla #15 (calefacción no funciona)
(15, 1, 'Revisando el sistema de calefacción ahora.', CURRENT_TIMESTAMP - INTERVAL '9 hours 30 minutes'),
(15, 3, 'Hace realmente frío. ¿Cuánto tiempo tomará?', CURRENT_TIMESTAMP - INTERVAL '9 hours'),
(15, 1, 'Encontramos el problema. Válvula principal cerrada por error. Solucionado en 10 minutos.', CURRENT_TIMESTAMP - INTERVAL '8 hours 45 minutes')
ON CONFLICT DO NOTHING;
