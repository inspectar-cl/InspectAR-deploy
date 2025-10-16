#!/bin/bash
# Script para generar firmas de prueba usando Python PIL

echo "=== Generador de Firmas de Prueba ==="
echo ""

# Verificar si Python está instalado
if ! command -v python3 &> /dev/null; then
    echo "❌ Python3 no está instalado. Por favor instálalo primero."
    exit 1
fi

# Instalar Pillow si no está instalado
echo "📦 Verificando dependencias..."
python3 -c "import PIL" 2>/dev/null
if [ $? -ne 0 ]; then
    echo "Instalando Pillow..."
    pip3 install Pillow --quiet
fi

# Directorio de destino
DEST_DIR="$(dirname "$0")"

# Crear firmas con Python
python3 << 'EOF'
from PIL import Image, ImageDraw, ImageFont
import os

def crear_firma(nombre, texto, filename):
    # Crear imagen
    width, height = 400, 150
    img = Image.new('RGB', (width, height), color='white')
    draw = ImageDraw.Draw(img)
    
    # Intentar usar una fuente, si no existe usar la por defecto
    try:
        font = ImageFont.truetype("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 32)
        font_small = ImageFont.truetype("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 20)
    except:
        font = ImageFont.load_default()
        font_small = ImageFont.load_default()
    
    # Dibujar línea decorativa superior
    draw.line([(20, 30), (380, 30)], fill='#007bff', width=2)
    
    # Dibujar texto principal (firma)
    text_width = draw.textlength(texto, font=font)
    x = (width - text_width) / 2
    draw.text((x, 55), texto, fill='#000000', font=font)
    
    # Dibujar texto pequeño abajo
    small_text = nombre
    small_text_width = draw.textlength(small_text, font=font_small)
    x_small = (width - small_text_width) / 2
    draw.text((x_small, 105), small_text, fill='#666666', font=font_small)
    
    # Dibujar línea decorativa inferior
    draw.line([(20, 120), (380, 120)], fill='#007bff', width=2)
    
    # Guardar
    dest = os.path.join(os.path.dirname(os.path.abspath(__file__)), filename)
    img.save(dest, 'JPEG', quality=95)
    print(f"✓ {filename} creada")

# Crear las 4 firmas
crear_firma("Analista del Sistema", "Juan Analista", "firma_usuario_1.jpg")
crear_firma("Técnico Especializado", "María Técnica", "firma_usuario_2.jpg")
crear_firma("Residente", "Pedro Residente", "firma_usuario_3.jpg")
crear_firma("Administrador", "Ana Administradora", "firma_usuario_4.jpg")

print("\n✅ Todas las firmas han sido generadas exitosamente")
EOF

echo ""
echo "📁 Las firmas han sido guardadas en: $DEST_DIR"
echo ""
echo "Archivos creados:"
ls -lh "$DEST_DIR"/firma_usuario_*.jpg 2>/dev/null || echo "No se encontraron archivos generados"
echo ""
echo "🚀 Ahora puedes reconstruir el contenedor con: docker-compose build gestion-service"
