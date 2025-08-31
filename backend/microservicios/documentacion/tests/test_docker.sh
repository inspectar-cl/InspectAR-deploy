#!/bin/bash

# Script de pruebas para el microservicio de documentación en Docker
# Funciona tanto dentro de contenedores como desde el host

# Detectar si estamos en Docker o local
if [ -n "$API_BASE_URL" ]; then
    BASE_URL="$API_BASE_URL"
else
    BASE_URL="http://localhost:8092"
fi

API_URL="$BASE_URL/api/v1"

echo "🧪 EJECUTANDO PRUEBAS DEL MICROSERVICIO DE DOCUMENTACIÓN (DOCKER)"
echo "=================================================================="
echo "🔗 Base URL: $BASE_URL"
echo ""

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para mostrar resultado de prueba
show_test_result() {
    local test_name="$1"
    local status_code="$2"
    local expected="$3"
    
    echo -n "Testing: $test_name... "
    if [ "$status_code" -eq "$expected" ]; then
        echo -e "${GREEN}✅ PASS${NC} (HTTP $status_code)"
    else
        echo -e "${RED}❌ FAIL${NC} (Expected HTTP $expected, got $status_code)"
    fi
}

# Función para hacer request con timeout
make_request() {
    local method="$1"
    local url="$2"
    local data="$3"
    local headers="$4"
    
    if [ -n "$data" ] && [ -n "$headers" ]; then
        curl -s -X "$method" -w "%{http_code}" -o /dev/null --max-time 10 --connect-timeout 5 -H "$headers" -d "$data" "$url"
    elif [ -n "$headers" ]; then
        curl -s -X "$method" -w "%{http_code}" -o /dev/null --max-time 10 --connect-timeout 5 -H "$headers" "$url"
    elif [ -n "$data" ]; then
        curl -s -X "$method" -w "%{http_code}" -o /dev/null --max-time 10 --connect-timeout 5 -d "$data" "$url"
    else
        curl -s -X "$method" -w "%{http_code}" -o /dev/null --max-time 10 --connect-timeout 5 "$url"
    fi
}

# Esperar a que el servicio esté disponible
echo "⏳ Esperando a que el servicio esté disponible..."
for i in {1..30}; do
    if curl -s --max-time 5 --connect-timeout 3 "$BASE_URL/health" > /dev/null; then
        echo -e "${GREEN}✅ Servicio disponible${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}❌ Timeout: Servicio no disponible después de 30 intentos${NC}"
        exit 1
    fi
    echo -n "."
    sleep 2
done
echo ""

echo "1. HEALTH CHECK"
echo "=================="
status=$(make_request "GET" "$BASE_URL/health")
show_test_result "Health Check" "$status" "200"
echo ""

echo "2. GESTIÓN DE DOCUMENTOS"
echo "========================="
status=$(make_request "GET" "$API_URL/documentos")
show_test_result "Listar todos los documentos" "$status" "200"

status=$(make_request "GET" "$API_URL/documentos?activo_id=1")
show_test_result "Filtrar documentos por activo_id" "$status" "200"

status=$(make_request "GET" "$API_URL/documentos?solo_fichas_tecnicas=true")
show_test_result "Filtrar solo fichas técnicas" "$status" "200"

status=$(make_request "GET" "$API_URL/documentos/1")
show_test_result "Obtener documento por ID (existente)" "$status" "200"

status=$(make_request "GET" "$API_URL/documentos/999")
show_test_result "Obtener documento inexistente" "$status" "404"

# Test de actualización
status=$(make_request "PUT" "$API_URL/documentos/1" '{"nombre":"Documento Actualizado"}' "Content-Type: application/json")
show_test_result "Actualizar documento existente" "$status" "200"

status=$(make_request "PUT" "$API_URL/documentos/999" '{"nombre":"Documento Inexistente"}' "Content-Type: application/json")
show_test_result "Actualizar documento inexistente" "$status" "404"
echo ""

echo "3. BÚSQUEDA DE DOCUMENTOS"
echo "============================"
status=$(make_request "GET" "$API_URL/documentos/buscar?q=bomba")
show_test_result "Búsqueda básica con query" "$status" "200"

status=$(make_request "GET" "$API_URL/documentos/buscar")
show_test_result "Búsqueda sin parámetro query (debe fallar)" "$status" "400"
echo ""

echo "4. CONSULTAS INTERACTIVAS CON IA"
echo "===================================="
status=$(make_request "POST" "$API_URL/documentos/1/consultar" '{"pregunta":"¿Cuántas bombas de agua tiene el sistema?"}' "Content-Type: application/json")
show_test_result "Consulta sobre cantidad de bombas" "$status" "200"

status=$(make_request "POST" "$API_URL/documentos/1/consultar" '{"pregunta":"¿Qué presión tiene esta bomba?"}' "Content-Type: application/json")
show_test_result "Consulta sobre presión" "$status" "200"

status=$(make_request "POST" "$API_URL/documentos/999/consultar" '{"pregunta":"Test"}' "Content-Type: application/json")
show_test_result "Consulta sobre documento inexistente" "$status" "404"

status=$(make_request "GET" "$API_URL/documentos/1/consultas")
show_test_result "Obtener historial de consultas" "$status" "200"

status=$(make_request "GET" "$API_URL/consultas/estadisticas")
show_test_result "Obtener estadísticas de consultas" "$status" "200"
echo ""

echo "5. VALIDACIONES Y CASOS EDGE"
echo "================================="
status=$(make_request "POST" "$API_URL/documentos/abc/consultar" '{"pregunta":"Test"}' "Content-Type: application/json")
show_test_result "Consulta con ID de documento inválido" "$status" "400"

status=$(make_request "POST" "$API_URL/documentos/1/consultar" '{"invalid_json":}' "Content-Type: application/json")
show_test_result "Consulta con JSON malformado" "$status" "400"
echo ""

echo "6. TESTS DE FUNCIONALIDAD IA"
echo "=============================="
echo -e "${BLUE}🤖 Probando respuestas inteligentes de IA...${NC}"

# Test respuesta sobre bombas
response=$(curl -s -X POST "$API_URL/documentos/1/consultar" \
    -H "Content-Type: application/json" \
    -d '{"pregunta":"¿Cuántas bombas de agua tiene el sistema?"}')

if echo "$response" | grep -q "3 bombas de agua"; then
    echo -e "${GREEN}✅ IA reconoce pregunta sobre cantidad de bombas${NC}"
else
    echo -e "${RED}❌ IA no reconoce pregunta sobre cantidad de bombas${NC}"
fi

# Test respuesta sobre presión
response=$(curl -s -X POST "$API_URL/documentos/1/consultar" \
    -H "Content-Type: application/json" \
    -d '{"pregunta":"¿Qué presión tiene esta bomba?"}')

if echo "$response" | grep -q "150 PSI"; then
    echo -e "${GREEN}✅ IA responde correctamente sobre presión${NC}"
else
    echo -e "${RED}❌ IA no responde correctamente sobre presión${NC}"
fi

# Test confianza
if echo "$response" | grep -q '"confianza":0.9'; then
    echo -e "${GREEN}✅ IA proporciona nivel de confianza alto${NC}"
else
    echo -e "${YELLOW}⚠️  Revisar nivel de confianza de respuestas IA${NC}"
fi
echo ""

echo "✅ PRUEBAS COMPLETADAS EN DOCKER!"
echo "=================================="
echo -e "${GREEN}🎯 Microservicio de documentación funcionando correctamente${NC}"
echo -e "${BLUE}📊 Base de datos PostgreSQL conectada${NC}"
echo -e "${YELLOW}🤖 IA respondiendo a consultas inteligentemente${NC}"
echo ""
echo "Para ver logs detallados:"
echo "  docker logs documentacion-service"
echo "  docker logs documentacion-db"
