# InspectAR Backend

Este proyecto corresponde al backend de **InspectAR**, encargado de gestionar la lógica del servidor y la conexión con la base de datos.

## 🐳 Inicio rápido con Docker

### Requisitos

- Tener [Docker Desktop](https://www.docker.com/products/docker-desktop/) instalado y en ejecución.
- Verificar que el archivo `apigateway/.env` esté correctamente configurado con las variables de entorno necesarias.

### Iniciar los servicios

Para levantar los servicios con Docker, ejecuta:

```bash
docker compose up --build
```