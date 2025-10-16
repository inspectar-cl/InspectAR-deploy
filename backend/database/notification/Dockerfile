FROM postgres:15

ENV POSTGRES_USER=notification_user
ENV POSTGRES_PASSWORD=notification_pass
ENV POSTGRES_DB=notification_db

# Copia los archivos de inicialización
COPY ./docker-entrypoint-initdb.d/ /docker-entrypoint-initdb.d/

EXPOSE 5432
