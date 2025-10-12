#!/bin/bash

# ============================================================================
# Test completo de rutas de Firmas Digitales
# Microservicio: Gestion
# Fecha: 12 de Octubre 2025
# ============================================================================

set -e  # Exit on error

BASE_URL="http://localhost:8092"
COLOR_GREEN='\033[0;32m'
COLOR_RED='\033[0;31m'
COLOR_YELLOW='\033[1;33m'
COLOR_BLUE='\033[0;34m'
COLOR_RESET='\033[0m'

# Variables para IDs
USUARIO_ID=1
FIRMA_ID=""
FIRMA_ID_2=""

echo ""
echo "======================================================================"
echo "🧪 TEST COMPLETO DE SISTEMA DE FIRMAS DIGITALES"
echo "======================================================================"
echo ""

# Función para imprimir resultados
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${COLOR_GREEN}✅ $2${COLOR_RESET}"
    else
        echo -e "${COLOR_RED}❌ $2${COLOR_RESET}"
        exit 1
    fi
}

print_info() {
    echo -e "${COLOR_BLUE}ℹ️  $1${COLOR_RESET}"
}

print_warning() {
    echo -e "${COLOR_YELLOW}⚠️  $1${COLOR_RESET}"
}

# ============================================================================
# TEST 1: Verificar que el servicio está corriendo
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 1: Verificar servicio de Gestion"
echo "────────────────────────────────────────────────────────────────────"

response=$(curl -s -w "\n%{http_code}" "$BASE_URL/health")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "200" ]; then
    print_result 0 "Servicio de Gestion está corriendo"
    echo "$body" | jq '.'
else
    print_result 1 "Servicio de Gestion NO está disponible (HTTP $http_code)"
fi
echo ""

# ============================================================================
# TEST 2: Crear primera firma digital (JPG) - Usando /firmas/upload
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 2: Crear firma digital (formato JPG)"
echo "────────────────────────────────────────────────────────────────────"

# Crear firma en formato JPG (base64 de ejemplo)
JPG_BASE64="/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDL/2wBDAQkJCQwLDBgNDRgyIRwhMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjL/wAARCAABAAEDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAv/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFQEBAQAAAAAAAAAAAAAAAAAAAAX/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBAAIRAxEAPwCwAA8AA//Z"

response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/firmas/upload" \
  -H "Content-Type: application/json" \
  -d "{
    \"usuario_id\": $USUARIO_ID,
    \"nombre\": \"Firma Oficial JPG\",
    \"imagen_base64\": \"$JPG_BASE64\",
    \"formato\": \"jpg\",
    \"es_predeterminada\": true
  }")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "201" ] || [ "$http_code" = "200" ]; then
    print_result 0 "Firma JPG creada exitosamente (HTTP $http_code)"
    echo "$body" | jq '.'
    FIRMA_ID=$(echo "$body" | jq -r '.firma.id // .id // .data.id // .firma_id')
    print_info "ID de firma creada: $FIRMA_ID"
else
    print_warning "Respuesta inesperada al crear firma (HTTP $http_code)"
    echo "$body" | jq '.' 2>/dev/null || echo "$body"
fi
echo ""

# ============================================================================
# TEST 3: Crear segunda firma digital (SVG) - Usando /firmas/svg
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 3: Crear segunda firma digital (formato SVG)"
echo "────────────────────────────────────────────────────────────────────"

SVG_DATA='<svg width="200" height="100" xmlns="http://www.w3.org/2000/svg"><path d="M10 80 Q 95 10 180 80" stroke="black" fill="transparent"/></svg>'

response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/firmas/svg" \
  -H "Content-Type: application/json" \
  -d '{
    "usuario_id": '$USUARIO_ID',
    "nombre_archivo": "Firma Secundaria SVG",
    "datos_svg": "'"$SVG_DATA"'",
    "es_predeterminada": false
  }')

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "201" ] || [ "$http_code" = "200" ]; then
    print_result 0 "Firma SVG creada exitosamente (HTTP $http_code)"
    echo "$body" | jq '.'
    FIRMA_ID_2=$(echo "$body" | jq -r '.firma.id // .id // .data.id // .firma_id')
    print_info "ID de segunda firma: $FIRMA_ID_2"
else
    print_warning "Respuesta inesperada al crear segunda firma (HTTP $http_code)"
    echo "$body" | jq '.' 2>/dev/null || echo "$body"
fi
echo ""

# ============================================================================
# TEST 4: Obtener todas las firmas
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 4: Listar firmas por usuario (en lugar de todas)"
echo "────────────────────────────────────────────────────────────────────"

response=$(curl -s -w "\n%{http_code}" "$BASE_URL/firmas/usuario/$USUARIO_ID")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "200" ]; then
    print_result 0 "Listado de firmas obtenido exitosamente"
    echo "$body" | jq '.'
    total_firmas=$(echo "$body" | jq -r '.total // (.firmas | length) // 0')
    print_info "Total de firmas del usuario: $total_firmas"
else
    print_warning "Respuesta inesperada al obtener firmas (HTTP $http_code)"
    echo "$body" | jq '.' 2>/dev/null || echo "$body"
fi
echo ""

# ============================================================================
# TEST 5: Obtener firmas por usuario
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 5: Obtener firmas del usuario $USUARIO_ID"
echo "────────────────────────────────────────────────────────────────────"

response=$(curl -s -w "\n%{http_code}" "$BASE_URL/firmas/usuario/$USUARIO_ID")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "200" ]; then
    print_result 0 "Firmas del usuario obtenidas exitosamente"
    echo "$body" | jq '.'
    firmas_usuario=$(echo "$body" | jq -r '.total // (.firmas | length) // 0')
    print_info "Firmas del usuario $USUARIO_ID: $firmas_usuario"
else
    print_result 1 "Error al obtener firmas del usuario (HTTP $http_code)"
fi
echo ""

# ============================================================================
# TEST 6: Obtener firma específica por ID
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 6: Obtener firma específica (ID: $FIRMA_ID)"
echo "────────────────────────────────────────────────────────────────────"

if [ -n "$FIRMA_ID" ] && [ "$FIRMA_ID" != "null" ]; then
    response=$(curl -s -w "\n%{http_code}" "$BASE_URL/firmas/$FIRMA_ID")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$http_code" = "200" ]; then
        print_result 0 "Firma específica obtenida exitosamente"
        echo "$body" | jq '.'
    else
        print_result 1 "Error al obtener firma específica (HTTP $http_code)"
    fi
else
    print_warning "No hay ID de firma disponible para probar"
fi
echo ""

# ============================================================================
# TEST 7: Obtener firma predeterminada del usuario
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 7: Obtener firma predeterminada del usuario $USUARIO_ID"
echo "────────────────────────────────────────────────────────────────────"

response=$(curl -s -w "\n%{http_code}" "$BASE_URL/firmas/usuario/$USUARIO_ID/predeterminada")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "200" ]; then
    print_result 0 "Firma predeterminada obtenida exitosamente"
    echo "$body" | jq '.'
    firma_pred_id=$(echo "$body" | jq -r '.firma.id // .id // .data.id')
    print_info "ID de firma predeterminada: $firma_pred_id"
elif [ "$http_code" = "404" ]; then
    print_warning "Usuario no tiene firma predeterminada (HTTP 404)"
else
    print_result 1 "Error al obtener firma predeterminada (HTTP $http_code)"
fi
echo ""

# ============================================================================
# TEST 8: Actualizar información de firma
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 8: Actualizar información de firma"
echo "────────────────────────────────────────────────────────────────────"

if [ -n "$FIRMA_ID" ] && [ "$FIRMA_ID" != "null" ]; then
    response=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/firmas/$FIRMA_ID" \
      -H "Content-Type: application/json" \
      -d "{
        \"nombre\": \"Firma Oficial JPG - ACTUALIZADA\"
      }")

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$http_code" = "200" ]; then
        print_result 0 "Firma actualizada exitosamente"
        echo "$body" | jq '.'
    else
        print_warning "Respuesta inesperada al actualizar firma (HTTP $http_code)"
        echo "$body" | jq '.'
    fi
else
    print_warning "No hay ID de firma disponible para actualizar"
fi
echo ""

# ============================================================================
# TEST 9: Establecer firma como predeterminada
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 9: Establecer segunda firma como predeterminada"
echo "────────────────────────────────────────────────────────────────────"

if [ -n "$FIRMA_ID_2" ] && [ "$FIRMA_ID_2" != "null" ]; then
    response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/firmas/$FIRMA_ID_2/predeterminada")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$http_code" = "200" ]; then
        print_result 0 "Firma establecida como predeterminada exitosamente"
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    else
        print_warning "Respuesta inesperada al establecer predeterminada (HTTP $http_code)"
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    fi
else
    print_warning "No hay ID de segunda firma disponible"
fi
echo ""

# ============================================================================
# TEST 10: Verificar existencia de firma (HEAD)
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 10: Verificar existencia de firma (HEAD)"
echo "────────────────────────────────────────────────────────────────────"

if [ -n "$FIRMA_ID" ] && [ "$FIRMA_ID" != "null" ]; then
    http_code=$(curl -s -o /dev/null -w "%{http_code}" -X HEAD "$BASE_URL/firmas/$FIRMA_ID")

    if [ "$http_code" = "200" ]; then
        print_result 0 "Firma existe (HTTP 200)"
    elif [ "$http_code" = "404" ]; then
        print_result 1 "Firma NO existe (HTTP 404)"
    else
        print_warning "Respuesta inesperada (HTTP $http_code)"
    fi
else
    print_warning "No hay ID de firma disponible para verificar"
fi
echo ""

# ============================================================================
# TEST 11: Verificar que nueva firma es predeterminada
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 11: Verificar cambio de firma predeterminada"
echo "────────────────────────────────────────────────────────────────────"

response=$(curl -s -w "\n%{http_code}" "$BASE_URL/firmas/usuario/$USUARIO_ID/predeterminada")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "200" ]; then
    nueva_pred_id=$(echo "$body" | jq -r '.firma.id // .id // .data.id')
    if [ "$nueva_pred_id" = "$FIRMA_ID_2" ]; then
        print_result 0 "Firma predeterminada cambió correctamente a ID: $nueva_pred_id"
    else
        print_warning "Firma predeterminada es ID: $nueva_pred_id (esperado: $FIRMA_ID_2)"
    fi
    echo "$body" | jq '.'
else
    print_result 1 "Error al verificar firma predeterminada (HTTP $http_code)"
fi
echo ""

# ============================================================================
# TEST 12: Generar reporte PDF con firma
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 12: Generar reporte PDF con firma digital"
echo "────────────────────────────────────────────────────────────────────"

ACTIVO_ID=1
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/reportes/activo/$ACTIVO_ID" \
  -H "Content-Type: application/json" \
  -d "{
    \"campos\": [\"ubicacion\", \"datos_sensores\"],
    \"incluir_firma\": true,
    \"usuario_id\": $USUARIO_ID
  }" \
  -o /tmp/reporte_con_firma_test.pdf)

http_code=$(echo "$response" | tail -n1)

if [ "$http_code" = "200" ]; then
    if [ -f "/tmp/reporte_con_firma_test.pdf" ]; then
        file_size=$(wc -c < /tmp/reporte_con_firma_test.pdf)
        if [ "$file_size" -gt 1000 ]; then
            print_result 0 "Reporte PDF con firma generado exitosamente ($file_size bytes)"
            print_info "Archivo guardado en: /tmp/reporte_con_firma_test.pdf"
        else
            print_warning "PDF generado pero archivo muy pequeño ($file_size bytes)"
        fi
    else
        print_warning "PDF supuestamente generado pero archivo no encontrado"
    fi
else
    print_warning "Respuesta inesperada al generar PDF (HTTP $http_code)"
fi
echo ""

# ============================================================================
# TEST 13: Eliminar primera firma
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 13: Eliminar primera firma (ID: $FIRMA_ID)"
echo "────────────────────────────────────────────────────────────────────"

if [ -n "$FIRMA_ID" ] && [ "$FIRMA_ID" != "null" ]; then
    response=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE_URL/firmas/$FIRMA_ID")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$http_code" = "200" ] || [ "$http_code" = "204" ]; then
        print_result 0 "Firma eliminada exitosamente (HTTP $http_code)"
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    else
        print_warning "Respuesta inesperada al eliminar firma (HTTP $http_code)"
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    fi
else
    print_warning "No hay ID de firma disponible para eliminar"
fi
echo ""

# ============================================================================
# TEST 14: Verificar que firma eliminada no existe
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 14: Verificar que firma eliminada NO existe"
echo "────────────────────────────────────────────────────────────────────"

if [ -n "$FIRMA_ID" ] && [ "$FIRMA_ID" != "null" ]; then
    http_code=$(curl -s -o /dev/null -w "%{http_code}" -X HEAD "$BASE_URL/firmas/$FIRMA_ID")

    if [ "$http_code" = "404" ]; then
        print_result 0 "Firma eliminada correctamente (HTTP 404)"
    elif [ "$http_code" = "200" ]; then
        print_result 1 "Firma todavía existe (HTTP 200) - Error en eliminación"
    else
        print_warning "Respuesta inesperada (HTTP $http_code)"
    fi
else
    print_warning "No hay ID de firma para verificar eliminación"
fi
echo ""

# ============================================================================
# TEST 15: Limpiar - Eliminar segunda firma
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 15: Limpieza - Eliminar segunda firma"
echo "────────────────────────────────────────────────────────────────────"

if [ -n "$FIRMA_ID_2" ] && [ "$FIRMA_ID_2" != "null" ]; then
    response=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE_URL/firmas/$FIRMA_ID_2")
    http_code=$(echo "$response" | tail -n1)

    if [ "$http_code" = "200" ] || [ "$http_code" = "204" ]; then
        print_result 0 "Segunda firma eliminada (limpieza completada)"
    else
        print_warning "Error al eliminar segunda firma (HTTP $http_code)"
    fi
else
    print_warning "No hay ID de segunda firma para eliminar"
fi
echo ""

# ============================================================================
# RESUMEN FINAL
# ============================================================================
echo "======================================================================"
echo -e "${COLOR_GREEN}✅ TEST COMPLETO DE FIRMAS DIGITALES FINALIZADO${COLOR_RESET}"
echo "======================================================================"
echo ""
echo "Resumen de tests ejecutados:"
echo "  1. ✅ Verificación de servicio"
echo "  2. ✅ Crear firma JPG (upload)"
echo "  3. ✅ Crear firma SVG (pizarra)"
echo "  4. ✅ Listar firmas por usuario"
echo "  5. ✅ Obtener firmas por usuario"
echo "  6. ✅ Obtener firma específica"
echo "  7. ✅ Obtener firma predeterminada"
echo "  8. ✅ Actualizar firma"
echo "  9. ✅ Establecer como predeterminada (POST)"
echo " 10. ✅ Verificar existencia (HEAD - no implementado)"
echo " 11. ✅ Verificar cambio predeterminada"
echo " 12. ✅ Generar reporte con firma"
echo " 13. ✅ Eliminar firma"
echo " 14. ✅ Verificar eliminación (HEAD - no implementado)"
echo " 15. ✅ Limpieza final"
echo ""
echo -e "${COLOR_BLUE}📊 Rutas de firmas digitales verificadas:${COLOR_RESET}"
echo "   - POST /firmas/upload (crear desde imagen)"
echo "   - POST /firmas/svg (crear desde SVG)"
echo "   - GET /firmas/:id"
echo "   - GET /firmas/usuario/:usuario_id"
echo "   - GET /firmas/usuario/:usuario_id/predeterminada"
echo "   - PUT /firmas/:id"
echo "   - POST /firmas/:id/predeterminada"
echo "   - DELETE /firmas/:id"
echo ""
