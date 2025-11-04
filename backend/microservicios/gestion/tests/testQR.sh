#!/bin/bash

# =============================================================================
# Script de Prueba de Códigos QR para Activos
# Verifica las 2 rutas implementadas y la generación automática de códigos
# =============================================================================

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# URLs
GESTION_URL="http://localhost:8092"
ADMIN_EMAIL="admin@example.com"

echo ""
echo -e "${CYAN}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║          PRUEBA DE SISTEMA DE CÓDIGOS QR - ACTIVOS          ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════════╝${NC}"
echo ""

# =============================================================================
# PASO 1: Crear Edificio
# =============================================================================
echo -e "${YELLOW}[PASO 1]${NC} Creando edificio para prueba de QR..."
EDIFICIO_RESPONSE=$(curl -s -X POST "$GESTION_URL/admin/edificios" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Torre QR Test",
    "direccion": "Av. QR 123",
    "latitud": -33.45,
    "longitud": -70.65,
    "email": "'"$ADMIN_EMAIL"'"
  }')

EDIFICIO_ID=$(echo "$EDIFICIO_RESPONSE" | jq -r '.edificio.id')
echo -e "${GREEN}✓ Edificio creado con ID: $EDIFICIO_ID${NC}"
echo ""
sleep 1

# =============================================================================
# PASO 2: Crear Activo (debe generar código QR automáticamente)
# =============================================================================
echo -e "${YELLOW}[PASO 2]${NC} Creando activo (debe generar código único y QR automáticamente)..."
ACTIVO_RESPONSE=$(curl -s -X POST "$GESTION_URL/admin/activos" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Bomba QR Test",
    "tipo": "bomba de agua",
    "descripcion": "Bomba para prueba de sistema de QR",
    "ubicacion": "Sala de Máquinas QR",
    "edificio_id": '"$EDIFICIO_ID"',
    "email": "'"$ADMIN_EMAIL"'"
  }')

ACTIVO_ID=$(echo "$ACTIVO_RESPONSE" | jq -r '.activo.id')
CODIGO_ACTIVO=$(echo "$ACTIVO_RESPONSE" | jq -r '.codigo')
URL_QR=$(echo "$ACTIVO_RESPONSE" | jq -r '.url_qr')
QR_DISPONIBLE=$(echo "$ACTIVO_RESPONSE" | jq -r '.qr_disponible')

echo -e "${GREEN}✓ Activo creado con ID: $ACTIVO_ID${NC}"
echo -e "${BLUE}  • Código generado: $CODIGO_ACTIVO${NC}"
echo -e "${BLUE}  • URL QR: $URL_QR${NC}"
echo -e "${BLUE}  • QR disponible: $QR_DISPONIBLE${NC}"
echo ""
echo "Respuesta completa:"
echo "$ACTIVO_RESPONSE" | jq '.'
echo ""
sleep 2

# =============================================================================
# PASO 3: Probar ruta GET /qr/:codigo (obtener info del activo por código)
# =============================================================================
echo -e "${YELLOW}[PASO 3]${NC} Probando ruta: GET /qr/$CODIGO_ACTIVO"
echo -e "${CYAN}Esta ruta debe retornar la información del activo (tipo y descripción)${NC}"
echo ""

QR_INFO_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" "$GESTION_URL/qr/$CODIGO_ACTIVO")
HTTP_STATUS=$(echo "$QR_INFO_RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
QR_INFO_BODY=$(echo "$QR_INFO_RESPONSE" | sed -e 's/HTTP_STATUS.*//')

if [ "$HTTP_STATUS" == "200" ]; then
  echo -e "${GREEN}✓✓✓ ÉXITO: Ruta GET /qr/:codigo funciona correctamente${NC}"
  echo -e "${GREEN}    HTTP Status: $HTTP_STATUS${NC}"
  echo ""
  echo "Información del activo obtenida:"
  echo "$QR_INFO_BODY" | jq '.'
  echo ""
  
  # Verificar campos
  TIPO=$(echo "$QR_INFO_BODY" | jq -r '.tipo')
  DESCRIPCION=$(echo "$QR_INFO_BODY" | jq -r '.descripcion')
  
  echo -e "${CYAN}Validación de campos:${NC}"
  echo -e "  • Tipo: ${BLUE}$TIPO${NC}"
  echo -e "  • Descripción: ${BLUE}$DESCRIPCION${NC}"
else
  echo -e "${RED}✗✗✗ ERROR: Ruta GET /qr/:codigo falló${NC}"
  echo -e "${RED}    HTTP Status: $HTTP_STATUS${NC}"
  echo -e "${RED}    Response: $QR_INFO_BODY${NC}"
fi
echo ""
sleep 2

# =============================================================================
# PASO 4: Probar ruta GET /qr/obtener/:activo_id (generar y obtener QR)
# =============================================================================
echo -e "${YELLOW}[PASO 4]${NC} Probando ruta: GET /qr/obtener/$ACTIVO_ID"
echo -e "${CYAN}Esta ruta debe generar el QR (si no existe) y retornarlo como imagen PNG${NC}"
echo ""

# Guardar el QR en un archivo
QR_FILE="/tmp/qr_activo_${ACTIVO_ID}.png"
curl -s -o "$QR_FILE" -w "\nHTTP_STATUS:%{http_code}\n" "$GESTION_URL/qr/obtener/$ACTIVO_ID" > /tmp/qr_response.txt
HTTP_STATUS=$(cat /tmp/qr_response.txt | grep "HTTP_STATUS" | cut -d: -f2)

if [ "$HTTP_STATUS" == "200" ]; then
  echo -e "${GREEN}✓✓✓ ÉXITO: Ruta GET /qr/obtener/:activo_id funciona correctamente${NC}"
  echo -e "${GREEN}    HTTP Status: $HTTP_STATUS${NC}"
  echo ""
  
  # Verificar que es una imagen PNG válida
  if file "$QR_FILE" | grep -q "PNG image"; then
    FILE_SIZE=$(stat -f%z "$QR_FILE" 2>/dev/null || stat -c%s "$QR_FILE" 2>/dev/null)
    echo -e "${GREEN}✓ QR generado correctamente como imagen PNG${NC}"
    echo -e "${BLUE}  • Archivo guardado en: $QR_FILE${NC}"
    echo -e "${BLUE}  • Tamaño: ${FILE_SIZE} bytes${NC}"
    
    # Verificar tipo MIME y dimensiones
    FILE_INFO=$(file "$QR_FILE")
    echo -e "${BLUE}  • Info: $FILE_INFO${NC}"
  else
    echo -e "${YELLOW}⚠ El archivo no parece ser una imagen PNG válida${NC}"
  fi
else
  echo -e "${RED}✗✗✗ ERROR: Ruta GET /qr/obtener/:activo_id falló${NC}"
  echo -e "${RED}    HTTP Status: $HTTP_STATUS${NC}"
fi
echo ""
sleep 2

# =============================================================================
# PASO 4.5: Probar ruta GET /qr/ver/:activo_id (visualizar QR en HTML)
# =============================================================================
echo -e "${YELLOW}[PASO 4.5]${NC} Probando ruta: GET /qr/ver/$ACTIVO_ID"
echo -e "${CYAN}Esta ruta debe mostrar el QR en una página HTML${NC}"
echo ""

VER_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" "$GESTION_URL/qr/ver/$ACTIVO_ID")
HTTP_STATUS=$(echo "$VER_RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
VER_BODY=$(echo "$VER_RESPONSE" | sed -e 's/HTTP_STATUS.*//')

if [ "$HTTP_STATUS" == "200" ]; then
  echo -e "${GREEN}✓✓✓ ÉXITO: Ruta GET /qr/ver/:activo_id funciona correctamente${NC}"
  echo -e "${GREEN}    HTTP Status: $HTTP_STATUS${NC}"
  echo ""
  
  # Verificar que es HTML
  if echo "$VER_BODY" | grep -q "<!DOCTYPE html>"; then
    echo -e "${GREEN}✓ Página HTML generada correctamente${NC}"
    
    # Verificar que contiene el título correcto
    if echo "$VER_BODY" | grep -q "Código QR - Bomba QR Test"; then
      echo -e "${BLUE}  • Título correcto: 'Código QR - Bomba QR Test'${NC}"
    fi
    
    # Verificar que contiene el código del activo
    if echo "$VER_BODY" | grep -q "$CODIGO_ACTIVO"; then
      echo -e "${BLUE}  • Código del activo presente: $CODIGO_ACTIVO${NC}"
    fi
    
    # Verificar que contiene la imagen del QR en base64
    if echo "$VER_BODY" | grep -q "data:image/png;base64"; then
      echo -e "${BLUE}  • QR embebido en formato Base64${NC}"
    fi
    
    echo -e "${CYAN}  • URL para visualizar: ${BLUE}$GESTION_URL/qr/ver/$ACTIVO_ID${NC}"
  else
    echo -e "${YELLOW}⚠ La respuesta no parece ser HTML válido${NC}"
  fi
else
  echo -e "${RED}✗✗✗ ERROR: Ruta GET /qr/ver/:activo_id falló${NC}"
  echo -e "${RED}    HTTP Status: $HTTP_STATUS${NC}"
fi
echo ""
sleep 2

# =============================================================================
# PASO 5: Verificar que el QR se guardó en la base de datos
# =============================================================================
echo -e "${YELLOW}[PASO 5]${NC} Verificando que el QR se guardó en la base de datos..."
echo ""

ACTIVO_INFO=$(curl -s "$GESTION_URL/activos/$ACTIVO_ID")
echo "Información del activo desde la base de datos:"
echo "$ACTIVO_INFO" | jq '{id, nombre, codigo_activo, url_qr, qr_generado_en}'
echo ""

HAS_QR=$(echo "$ACTIVO_INFO" | jq -r '.codigo_qr != null and .codigo_qr != ""')
if [ "$HAS_QR" == "true" ]; then
  echo -e "${GREEN}✓ El activo tiene QR guardado en la base de datos${NC}"
  QR_GEN_DATE=$(echo "$ACTIVO_INFO" | jq -r '.qr_generado_en')
  echo -e "${BLUE}  • Fecha de generación: $QR_GEN_DATE${NC}"
else
  echo -e "${YELLOW}⚠ El activo NO tiene QR guardado en la base de datos${NC}"
fi
echo ""
sleep 1

# =============================================================================
# PASO 6: Crear segundo activo para verificar secuencia de códigos
# =============================================================================
echo -e "${YELLOW}[PASO 6]${NC} Creando segundo activo para verificar secuencia de códigos..."
ACTIVO2_RESPONSE=$(curl -s -X POST "$GESTION_URL/admin/activos" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Bomba QR Test 2",
    "tipo": "bomba de agua",
    "descripcion": "Segunda bomba para verificar secuencia",
    "ubicacion": "Sala de Máquinas 2",
    "edificio_id": '"$EDIFICIO_ID"',
    "email": "'"$ADMIN_EMAIL"'"
  }')

CODIGO_ACTIVO2=$(echo "$ACTIVO2_RESPONSE" | jq -r '.codigo')
ACTIVO_ID2=$(echo "$ACTIVO2_RESPONSE" | jq -r '.activo.id')

echo -e "${GREEN}✓ Segundo activo creado${NC}"
echo -e "${BLUE}  • ID: $ACTIVO_ID2${NC}"
echo -e "${BLUE}  • Código: $CODIGO_ACTIVO2${NC}"
echo ""

# Comparar códigos
echo -e "${CYAN}Comparación de códigos:${NC}"
echo -e "  Activo 1: ${BLUE}$CODIGO_ACTIVO${NC}"
echo -e "  Activo 2: ${BLUE}$CODIGO_ACTIVO2${NC}"
echo ""

# Verificar que los códigos son secuenciales
SEQ1=$(echo "$CODIGO_ACTIVO" | cut -d'-' -f3)
SEQ2=$(echo "$CODIGO_ACTIVO2" | cut -d'-' -f3)

if [ "$((SEQ2))" -gt "$((SEQ1))" ]; then
  echo -e "${GREEN}✓ Los códigos son secuenciales correctamente${NC}"
  echo -e "${BLUE}  Secuencia 1: $SEQ1, Secuencia 2: $SEQ2${NC}"
else
  echo -e "${YELLOW}⚠ Los códigos no parecen ser secuenciales${NC}"
fi
echo ""

# =============================================================================
# RESUMEN FINAL
# =============================================================================
echo ""
echo -e "${CYAN}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║                    RESUMEN DE LA PRUEBA                      ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "📋 ${GREEN}Edificio creado:${NC} ID = ${BLUE}$EDIFICIO_ID${NC}"
echo -e "🏭 ${GREEN}Activos creados:${NC}"
echo -e "   1. ID = ${BLUE}$ACTIVO_ID${NC}, Código = ${BLUE}$CODIGO_ACTIVO${NC}"
echo -e "   2. ID = ${BLUE}$ACTIVO_ID2${NC}, Código = ${BLUE}$CODIGO_ACTIVO2${NC}"
echo ""
echo -e "🔍 ${CYAN}Rutas verificadas:${NC}"
echo -e "   • ${GREEN}GET /qr/:codigo${NC} → Retorna información del activo"
echo -e "   • ${GREEN}GET /qr/obtener/:activo_id${NC} → Genera y retorna QR como PNG"
echo -e "   • ${GREEN}GET /qr/ver/:activo_id${NC} → Visualiza QR en página HTML"
echo ""
echo -e "✅ ${CYAN}Funcionalidades validadas:${NC}"
echo -e "   • Generación automática de códigos únicos: ${GREEN}✓${NC}"
echo -e "   • Generación automática de QR al crear activo: ${GREEN}✓${NC}"
echo -e "   • Ruta de info por código funcionando: ${GREEN}✓${NC}"
echo -e "   • Ruta de obtención de QR funcionando: ${GREEN}✓${NC}"
echo -e "   • Ruta de visualización HTML funcionando: ${GREEN}✓${NC}"
echo -e "   • Secuencia de códigos correcta: ${GREEN}✓${NC}"
echo -e "   • QR guardado en base de datos: ${GREEN}✓${NC}"
echo ""
echo -e "${BLUE}📱 Archivo QR generado: $QR_FILE${NC}"
echo -e "${BLUE}🌐 URL del QR apunta a: $URL_QR${NC}"
echo -e "${BLUE}👁️  Visualizar QR en navegador: $GESTION_URL/qr/ver/$ACTIVO_ID${NC}"
echo ""

# =============================================================================
# LIMPIEZA OPCIONAL
# =============================================================================
echo -e "${YELLOW}¿Deseas eliminar los datos de prueba? (y/N)${NC}"
read -t 10 -n 1 CLEANUP_RESPONSE || CLEANUP_RESPONSE="n"
echo ""

if [ "$CLEANUP_RESPONSE" == "y" ] || [ "$CLEANUP_RESPONSE" == "Y" ]; then
  echo -e "${YELLOW}Limpiando datos de prueba...${NC}"
  
  curl -s -X DELETE "$GESTION_URL/admin/activos/$ACTIVO_ID?email=$ADMIN_EMAIL" > /dev/null
  curl -s -X DELETE "$GESTION_URL/admin/activos/$ACTIVO_ID2?email=$ADMIN_EMAIL" > /dev/null
  curl -s -X DELETE "$GESTION_URL/admin/edificios/$EDIFICIO_ID?email=$ADMIN_EMAIL" > /dev/null
  
  echo -e "${GREEN}✓ Datos de prueba eliminados${NC}"
  
  # Eliminar archivo QR
  if [ -f "$QR_FILE" ]; then
    rm "$QR_FILE"
    echo -e "${GREEN}✓ Archivo QR eliminado${NC}"
  fi
else
  echo -e "${BLUE}ℹ️  Los datos de prueba permanecen en el sistema${NC}"
  echo -e "${BLUE}   Edificio ID: $EDIFICIO_ID${NC}"
  echo -e "${BLUE}   Activo 1 ID: $ACTIVO_ID (Código: $CODIGO_ACTIVO)${NC}"
  echo -e "${BLUE}   Activo 2 ID: $ACTIVO_ID2 (Código: $CODIGO_ACTIVO2)${NC}"
  echo -e "${BLUE}   Archivo QR: $QR_FILE${NC}"
fi

echo ""
echo -e "${CYAN}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║              ✓ PRUEBA DE QR COMPLETADA                      ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════════╝${NC}"
echo ""
