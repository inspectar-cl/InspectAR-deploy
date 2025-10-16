#!/bin/bash
# Script de entrada que inicializa archivos PDF y luego inicia el servicio

echo "🚀 Iniciando microservicio de documentación..."

# Ejecutar inicialización de archivos PDF
if [ -f "./scripts/init-files.sh" ]; then
    echo "📁 Ejecutando inicialización de archivos PDF..."
    bash ./scripts/init-files.sh
fi

# Iniciar el servicio principal
echo "🎯 Iniciando servicio de documentación..."
exec "./documentacion-service"
