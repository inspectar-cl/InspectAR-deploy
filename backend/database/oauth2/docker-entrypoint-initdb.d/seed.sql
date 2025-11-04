-- Insertar users de prueba con IDs específicos
INSERT INTO users (username, email, password, scope) 
VALUES ('carlos.mendoza', 'analista@example.com', '$2a$10$AVOlu9c1j6204IpKPjaqcelWG.E5VXPk2538Eu.0j/chBV1ZWBBFG', 'user-type:Analista');

INSERT INTO users (username, email, password, scope) 
VALUES ('maria.gonzalez', 'tecnico@example.com', '$2a$10$AVOlu9c1j6204IpKPjaqcelWG.E5VXPk2538Eu.0j/chBV1ZWBBFG', 'user-type:Tecnico');

INSERT INTO users (username, email, password, scope) 
VALUES ('pedro.silva', 'residente@example.com', '$2a$10$AVOlu9c1j6204IpKPjaqcelWG.E5VXPk2538Eu.0j/chBV1ZWBBFG', 'user-type:Residente');

INSERT INTO users (username, email, password, scope) 
VALUES ('ana.rodriguez', 'admin@example.com', '$2a$10$AVOlu9c1j6204IpKPjaqcelWG.E5VXPk2538Eu.0j/chBV1ZWBBFG', 'user-type:Residente');

-- Usuarios Administradores
INSERT INTO users (username, email, password, scope) 
VALUES ('admin.norte', 'admin.norte@inspectAR.com', '$2a$10$AVOlu9c1j6204IpKPjaqcelWG.E5VXPk2538Eu.0j/chBV1ZWBBFG', 'user-type:Administrador');

INSERT INTO users (username, email, password, scope) 
VALUES ('admin.centro', 'admin.centro@inspectAR.com', '$2a$10$AVOlu9c1j6204IpKPjaqcelWG.E5VXPk2538Eu.0j/chBV1ZWBBFG', 'user-type:Administrador');

-- Usuario Root
INSERT INTO users (username, email, password, scope) 
VALUES ('root.system', 'root@inspectAR.com', '$2a$10$AVOlu9c1j6204IpKPjaqcelWG.E5VXPk2538Eu.0j/chBV1ZWBBFG', 'user-type:Root');

-- Insertar configuración

INSERT INTO config (key_access, key_refresh, lifetime_access, lifetime_refresh) VALUES ('test_password', 'refresh_test_password', 15, 10080);

-- Reiniciar la secuencia de IDs para que PostgreSQL continúe desde el siguiente valor correcto
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));