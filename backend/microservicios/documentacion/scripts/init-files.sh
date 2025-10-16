#!/bin/bash
# Script de inicialización para copiar archivos PDF iniciales

echo "Inicializando archivos PDF del microservicio de documentación..."

# Crear directorio de uploads
mkdir -p /app/uploads

# Verificar si ya existen archivos (para evitar sobrescribir en reinicios)
if [ ! -f "/app/uploads/1_ficha_tecnica_caldera_bosch.pdf" ]; then
    echo "Copiando archivos PDF iniciales..."
    
    # Copiar archivos desde el directorio de archivos iniciales
    if [ -d "/app/storage/initial-files" ]; then
        cp /app/storage/initial-files/*.pdf /app/uploads/
        echo "✅ Archivos PDF copiados exitosamente:"
        ls -la /app/uploads/
    else
        echo "⚠️  Directorio de archivos iniciales no encontrado"
    fi
else
    echo "✅ Archivos PDF ya existen, omitiendo copia inicial"
fi

echo "🎯 Inicialización de archivos PDF completada"
