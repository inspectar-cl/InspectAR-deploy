FROM postgres:16-alpine

# Copiar scripts de inicialización
COPY docker-entrypoint-initdb.d/*.sql /docker-entrypoint-initdb.d/

# Los scripts se ejecutarán automáticamente en orden alfabético
