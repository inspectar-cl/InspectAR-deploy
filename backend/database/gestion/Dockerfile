FROM postgres:15

# Copiar scripts de inicialización
COPY docker-entrypoint-initdb.d/ /docker-entrypoint-initdb.d/

# Configurar timezone
ENV TZ=America/Santiago
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

EXPOSE 5432
