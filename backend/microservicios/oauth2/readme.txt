Para probar el microservicio es necesario tener una tabla users con esto (por ahora):

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    password TEXT NOT NULL
);
Luego insertar usuario admin para probar:

INSERT INTO users (username, password) VALUES ('admin', '$2a$10$AVOlu9c1j6204IpKPjaqcelWG.E5VXPk2538Eu.0j/chBV1ZWBBFG');

Aun no tiene el tema de hash de contraseñas al crear usuario nuevo por ejemplo (Aunque no deberia ser parte del servicio de utenticacion).

Luego para testear, en postman colocar lo siguiente:
{
    "username": "admin",
    "password": "123456"
}

deberia retornar el token