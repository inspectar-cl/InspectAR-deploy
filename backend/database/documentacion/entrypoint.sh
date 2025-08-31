#!/bin/bash
set -e

# Script de entrada personalizado para base de datos de documentación
echo "🚀 Iniciando PostgreSQL para microservicio de documentación..."

# Ejecutar el entrypoint original de PostgreSQL
exec docker-entrypoint.sh "$@"
