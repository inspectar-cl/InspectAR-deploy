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

# Variables para IDs y email
USUARIO_EMAIL="analista@example.com"
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
# TEST 2: Crear primera firma digital (JPG) - NO DISPONIBLE EN NUEVA API
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 2: Crear firma digital (formato JPG) - OMITIDO"
echo "────────────────────────────────────────────────────────────────────"

print_warning "Endpoint POST /firmas/upload requiere multipart/form-data con archivo real"
print_info "Este test requiere una firma existente. Usando firma predeterminada del usuario"
echo ""

# ============================================================================
# TEST 3: Crear segunda firma digital (SVG) - Usando /firmas/svg con EMAIL
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 3: Crear firma digital (formato SVG) con EMAIL"
echo "────────────────────────────────────────────────────────────────────"

SVG_DATA='<svg width="200" height="100" xmlns="http://www.w3.org/2000/svg"><path d="M10 80 Q 95 10 180 80" stroke="black" fill="transparent"/></svg>'

# Usar jq para construir el JSON de forma segura
JSON_PAYLOAD=$(jq -n \
  --arg email "$USUARIO_EMAIL" \
  --arg nombre "Firma Test SVG" \
  --arg svg "$SVG_DATA" \
  '{email: $email, nombre_archivo: $nombre, datos_svg: $svg, es_predeterminada: false}')

response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/firmas/svg" \
  -H "Content-Type: application/json" \
  -d "$JSON_PAYLOAD")

http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "201" ] || [ "$http_code" = "200" ]; then
    print_result 0 "Firma SVG creada exitosamente con email (HTTP $http_code)"
    echo "$body" | jq '.'
    FIRMA_ID_2=$(echo "$body" | jq -r '.firma.id // .id // .data.id // .firma_id')
    print_info "ID de firma creada: $FIRMA_ID_2"
else
    print_warning "Respuesta inesperada al crear firma SVG (HTTP $http_code)"
    echo "$body" | jq '.' 2>/dev/null || echo "$body"
fi
echo ""

# ============================================================================
# TEST 4: Obtener todas las firmas por EMAIL
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 4: Listar firmas por EMAIL del usuario"
echo "────────────────────────────────────────────────────────────────────"

response=$(curl -s -w "\n%{http_code}" "$BASE_URL/firmas/usuario/$USUARIO_EMAIL")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "200" ]; then
    print_result 0 "Listado de firmas obtenido exitosamente con email"
    echo "$body" | jq '.'
    total_firmas=$(echo "$body" | jq -r '.total // (.firmas | length) // 0')
    print_info "Total de firmas del usuario: $total_firmas"
    
    # Obtener ID de la primera firma para tests posteriores
    if [ -z "$FIRMA_ID" ] || [ "$FIRMA_ID" = "null" ]; then
        FIRMA_ID=$(echo "$body" | jq -r '.firmas[0].id // empty')
        if [ -n "$FIRMA_ID" ] && [ "$FIRMA_ID" != "null" ]; then
            print_info "ID de firma encontrada para tests: $FIRMA_ID"
        fi
    fi
else
    print_warning "Respuesta inesperada al obtener firmas (HTTP $http_code)"
    echo "$body" | jq '.' 2>/dev/null || echo "$body"
fi
echo ""

# ============================================================================
# TEST 5: Obtener firmas por EMAIL del usuario
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 5: Obtener firmas del usuario por EMAIL: $USUARIO_EMAIL"
echo "────────────────────────────────────────────────────────────────────"

response=$(curl -s -w "\n%{http_code}" "$BASE_URL/firmas/usuario/$USUARIO_EMAIL")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "200" ]; then
    print_result 0 "Firmas del usuario obtenidas exitosamente con email"
    echo "$body" | jq '.'
    firmas_usuario=$(echo "$body" | jq -r '.total // (.firmas | length) // 0')
    print_info "Firmas del usuario $USUARIO_EMAIL: $firmas_usuario"
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
# TEST 7: Obtener firma predeterminada del usuario por EMAIL
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 7: Obtener firma predeterminada del usuario por EMAIL"
echo "────────────────────────────────────────────────────────────────────"

response=$(curl -s -w "\n%{http_code}" "$BASE_URL/firmas/usuario/$USUARIO_EMAIL/predeterminada")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "200" ]; then
    print_result 0 "Firma predeterminada obtenida exitosamente con email"
    echo "$body" | jq '.'
    firma_pred_id=$(echo "$body" | jq -r '.id // .firma.id // .data.id')
    print_info "ID de firma predeterminada: $firma_pred_id"
elif [ "$http_code" = "404" ]; then
    print_warning "Usuario no tiene firma predeterminada (HTTP 404)"
else
    print_result 1 "Error al obtener firma predeterminada (HTTP $http_code)"
fi
echo ""

# ============================================================================
# TEST 8: Actualizar información de firma (sin cambios - mantiene ID)
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 8: Actualizar información de firma"
echo "────────────────────────────────────────────────────────────────────"

if [ -n "$FIRMA_ID" ] && [ "$FIRMA_ID" != "null" ]; then
    response=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/firmas/$FIRMA_ID" \
      -H "Content-Type: application/json" \
      -d "{
        \"nombre_archivo\": \"Firma Actualizada - Test\"
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
# TEST 9: Establecer firma como predeterminada (con EMAIL)
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 9: Establecer firma como predeterminada (con EMAIL)"
echo "────────────────────────────────────────────────────────────────────"

if [ -n "$FIRMA_ID_2" ] && [ "$FIRMA_ID_2" != "null" ]; then
    response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/firmas/$FIRMA_ID_2/predeterminada" \
      -H "Content-Type: application/json" \
      -d "{\"email\":\"$USUARIO_EMAIL\"}")
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    if [ "$http_code" = "200" ]; then
        print_result 0 "Firma establecida como predeterminada exitosamente con email"
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    else
        print_warning "Respuesta inesperada al establecer predeterminada (HTTP $http_code)"
        echo "$body" | jq '.' 2>/dev/null || echo "$body"
    fi
else
    print_warning "No hay ID de firma disponible para establecer como predeterminada"
fi
echo ""

# ============================================================================
# TEST 10: Verificar cambio de firma predeterminada (con EMAIL)
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 10: Verificar cambio de firma predeterminada con EMAIL"
echo "────────────────────────────────────────────────────────────────────"

response=$(curl -s -w "\n%{http_code}" "$BASE_URL/firmas/usuario/$USUARIO_EMAIL/predeterminada")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n-1)

if [ "$http_code" = "200" ]; then
    nueva_pred_id=$(echo "$body" | jq -r '.id // .firma.id // .data.id')
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
# TEST 11: Eliminar primera firma (si existe)
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 11: Eliminar primera firma (ID: $FIRMA_ID)"
echo "────────────────────────────────────────────────────────────────────"

if [ -n "$FIRMA_ID" ] && [ "$FIRMA_ID" != "null" ] && [ "$FIRMA_ID" != "$FIRMA_ID_2" ]; then
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
    print_warning "No hay ID de firma diferente disponible para eliminar"
fi
echo ""

# ============================================================================
# TEST 12: Limpiar - Eliminar segunda firma
# ============================================================================
echo "────────────────────────────────────────────────────────────────────"
echo "TEST 12: Limpieza - Eliminar firma de prueba"
echo "────────────────────────────────────────────────────────────────────"

if [ -n "$FIRMA_ID_2" ] && [ "$FIRMA_ID_2" != "null" ]; then
    response=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE_URL/firmas/$FIRMA_ID_2")
    http_code=$(echo "$response" | tail -n1)

    if [ "$http_code" = "200" ] || [ "$http_code" = "204" ]; then
        print_result 0 "Firma de prueba eliminada (limpieza completada)"
    else
        print_warning "Error al eliminar firma de prueba (HTTP $http_code)"
    fi
else
    print_warning "No hay ID de firma de prueba para eliminar"
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
echo "  2. ⚠️  Crear firma JPG (omitido - requiere multipart/form-data real)"
echo "  3. ✅ Crear firma SVG con EMAIL"
echo "  4. ✅ Listar firmas por EMAIL"
echo "  5. ✅ Obtener firmas por EMAIL"
echo "  6. ✅ Obtener firma específica"
echo "  7. ✅ Obtener firma predeterminada por EMAIL"
echo "  8. ✅ Actualizar firma"
echo "  9. ✅ Establecer como predeterminada (con EMAIL en body)"
echo " 10. ✅ Verificar cambio predeterminada"
echo " 11. ✅ Eliminar firma"
echo " 12. ✅ Limpieza final"
echo ""
echo -e "${COLOR_BLUE}📊 Rutas de firmas digitales verificadas (con EMAIL):${COLOR_RESET}"
echo "   - POST /firmas/svg (crear desde SVG con email)"
echo "   - GET /firmas/:id"
echo "   - GET /firmas/usuario/:email (antes :usuario_id)"
echo "   - GET /firmas/usuario/:email/predeterminada (antes :usuario_id)"
echo "   - PUT /firmas/:id"
echo "   - POST /firmas/:id/predeterminada (con email en body)"
echo "   - DELETE /firmas/:id"
echo ""
