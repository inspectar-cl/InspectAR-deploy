FROM postgres:15

ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=joytan
ENV POSTGRES_DB=db-oauth

# Copia solo los archivos necesarios para oauth2
COPY ./docker-entrypoint-initdb.d/ /docker-entrypoint-initdb.d/
