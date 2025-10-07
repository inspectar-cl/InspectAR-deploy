-- Insertar users de prueba con IDs específicos
INSERT INTO users (username, email, password, scope) 
VALUES ('admin', 'admin@example.com', '$2a$10$AVOlu9c1j6204IpKPjaqcelWG.E5VXPk2538Eu.0j/chBV1ZWBBFG', 'user-type:Analista');

INSERT INTO users (username, email, password, scope) 
VALUES ('tecnico', 'tecnico@example.com', '$2a$10$AVOlu9c1j6204IpKPjaqcelWG.E5VXPk2538Eu.0j/chBV1ZWBBFG', 'user-type:Tecnico');

INSERT INTO users (username, email, password, scope) 
VALUES ('usuario', 'usuario@example.com', '$2a$10$AVOlu9c1j6204IpKPjaqcelWG.E5VXPk2538Eu.0j/chBV1ZWBBFG', 'user-type:Residente');

-- Insertar configuración

INSERT INTO config (key_access, key_refresh, lifetime_access, lifetime_refresh) VALUES ('test_password', 'refresh_test_password', 15, 10080);

-- Reiniciar la secuencia de IDs para que PostgreSQL continúe desde el siguiente valor correcto
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));