FROM postgres:15

# Variables de entorno para PostgreSQL
ENV POSTGRES_DB=documentacion_db
ENV POSTGRES_USER=documentacion_user
ENV POSTGRES_PASSWORD=documentacion_pass
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8

# Copiar script de inicialización
COPY init.sql /docker-entrypoint-initdb.d/

# Exponer puerto estándar
EXPOSE 5434

# Configurar entrypoint personalizado
# COPY entrypoint.sh /usr/local/bin/
# RUN chmod +x /usr/local/bin/entrypoint.sh

# ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["postgres"]
