-- Regla para enrutar mensajes publicados en 'sensor/datos' a 'procesados/datos'
CREATE RULE route_sensor_to_procesados
WHEN message.publish
AND topic = 'sensor/datos'
THEN
    mqtt:publish("procesados/datos", payload);

-- Puedes agregar más reglas según tus necesidades
-- Por ejemplo, enrutar de 'sensor/alerta' a 'procesados/alerta'
CREATE RULE route_alerta_to_procesados
WHEN message.publish
AND topic = 'sensor/alerta'
THEN
    mqtt:publish("procesados/alerta", payload);