#!/bin/bash

# Test unificado para el microservicio de documentación
# Evalúa el estado completo del microservicio y sus funcionalidades principales

# Configuración
if [ -n "$API_BASE_URL" ]; then
    BASE_URL="$API_BASE_URL"
else
    BASE_URL="http://localhost:8093"
fi

API_URL="$BASE_URL/api/v1"

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Variables de conteo
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

echo -e "${CYAN}📄 MICROSERVICIO DE DOCUMENTACIÓN - TEST COMPLETO${NC}"
echo "=============================================================="
echo "🔗 URL Base: $BASE_URL"
echo "📅 Fecha: $(date)"
echo ""

# Directorio del script (para construir rutas absolutas a archivos de ejemplo)
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Función para hacer peticiones HTTP con timeout
make_request() {
    local method="$1"
    local url="$2"
    local data="$3"
    local headers="$4"
    
    if [ -n "$data" ] && [ -n "$headers" ]; then
        curl -s -X "$method" -w "HTTP_CODE:%{http_code}" --max-time 10 --connect-timeout 5 -H "$headers" -d "$data" "$url"
    elif [ -n "$headers" ]; then
        curl -s -X "$method" -w "HTTP_CODE:%{http_code}" --max-time 10 --connect-timeout 5 -H "$headers" "$url"
    elif [ -n "$data" ]; then
        curl -s -X "$method" -w "HTTP_CODE:%{http_code}" --max-time 10 --connect-timeout 5 -H "Content-Type: application/json" -d "$data" "$url"
    else
        curl -s -X "$method" -w "HTTP_CODE:%{http_code}" --max-time 10 --connect-timeout 5 "$url"
    fi
}

# Función para ejecutar tests
test_endpoint() {
    local method="$1"
    local url="$2"
    local description="$3"
    local data="$4"
    local expected_status="$5"
    local headers="$6"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    echo -e "\n${CYAN}📍 Test #$TOTAL_TESTS:${NC} $description"
    echo -e "${PURPLE}$method${NC} $url"
    
    response=$(make_request "$method" "$url" "$data" "$headers")
    http_code=$(echo "$response" | grep -o 'HTTP_CODE:[0-9]*' | cut -d: -f2)
    body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')
    
    # Verificar código de estado esperado (por defecto 200-299)
    if [ -z "$expected_status" ]; then
        expected_status="2"
    fi
    
    if [[ $http_code =~ ^${expected_status}[0-9]*$ ]]; then
        echo -e "${GREEN}✅ PASS - Status: $http_code${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ FAIL - Status: $http_code (esperado: ${expected_status}xx)${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    
    # Mostrar respuesta truncada
    if [ ${#body} -gt 200 ]; then
        echo "Response: ${body:0:200}..."
    else
        echo "Response: $body"
    fi
    
    # Store the response for further analysis
    LAST_RESPONSE="$body"
    return $http_code
}

# Verificar disponibilidad del servicio
echo -e "${CYAN}🔍 Verificando disponibilidad del servicio...${NC}"
for i in {1..30}; do
    if curl -s --max-time 5 --connect-timeout 3 "$BASE_URL/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Servicio de Documentación disponible${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}❌ TIMEOUT: Servicio no disponible después de 30 intentos${NC}"
        exit 1
    fi
    echo -n "."
    sleep 2
done

echo -e "\n${BLUE}🏥 TESTS BÁSICOS DE CONECTIVIDAD${NC}"
echo "----------------------------------------"

# Test 1: Health Check
test_endpoint "GET" "$BASE_URL/health" "Health Check" "" "200"

# Test 2: Endpoint de test
test_endpoint "GET" "$BASE_URL/test" "Test Route" "" "200"

echo -e "\n${BLUE}📄 TESTS DE GESTIÓN DE DOCUMENTOS${NC}"
echo "----------------------------------------"

# Test 3: Listar todos los documentos
test_endpoint "GET" "$API_URL/documentos" "Listar todos los documentos" "" "200"

# Test 4: Filtrar documentos por activo
test_endpoint "GET" "$API_URL/documentos?activo_id=1" "Filtrar documentos por activo_id" "" "200"

# Test 5: Filtrar solo fichas técnicas
test_endpoint "GET" "$API_URL/documentos?solo_fichas_tecnicas=true" "Filtrar solo fichas técnicas" "" "200"

# Test 5.1: Obtener ficha técnica por activo (nuevo endpoint)
test_endpoint "GET" "$API_URL/documentos/activo/1/ficha-tecnica" "Obtener ficha técnica por activo (existente)" "" "200"

# Test 5.2: Ficha técnica para activo inexistente
test_endpoint "GET" "$API_URL/documentos/activo/999/ficha-tecnica" "Obtener ficha técnica por activo inexistente" "" "404"

echo -e "\n${BLUE}⬆️  TESTS DE SUBIDA Y DESCARGA DE ARCHIVOS${NC}"
echo "----------------------------------------"

# Helper: test de subida multipart/form-data
test_upload_file() {
    local url="$1"
    local file_path="$2"
    local nombre="$3"
    local categoria="$4"
    local esperado="$5" # 201 para éxito, 4xx para fallo

    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -e "\n${CYAN}📍 Test #$TOTAL_TESTS:${NC} Subir archivo '$file_path' (categoria=$categoria)"
    echo -e "${PURPLE}POST${NC} $url"

    if [ ! -f "$file_path" ]; then
        echo -e "${YELLOW}⚠️  Archivo no encontrado: $file_path${NC}"
    fi

    response=$(curl -s -X POST -w "HTTP_CODE:%{http_code}" \
        -F "archivo=@${file_path}" \
        -F "activo_id=1" \
        -F "nombre=${nombre}" \
        -F "categoria=${categoria}" \
        -F "descripcion=Documento de prueba subido por test.sh" \
        -F "subido_por=tests" \
        "$url")

    http_code=$(echo "$response" | grep -o 'HTTP_CODE:[0-9]*' | cut -d: -f2)
    body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')

    if [[ $http_code == $esperado* ]]; then
        echo -e "${GREEN}✅ PASS - Status: $http_code${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ FAIL - Status: $http_code (esperado: ${esperado}xx)${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi

    # Guardar body y posible ID
    LAST_RESPONSE="$body"
    LAST_UPLOAD_ID=$(echo "$body" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
    if [ -n "$LAST_UPLOAD_ID" ]; then
        echo "ID subido: $LAST_UPLOAD_ID"
    fi
}

# Test 6: Subir PDF válido (usa archivo de ejemplo del repo)
# Construir ruta absoluta al archivo de ejemplo
PDF_SAMPLE="$SCRIPT_DIR/../storage/initial-files/1_manual_operacion_caldera.pdf"
test_upload_file "$API_URL/documentos" "$PDF_SAMPLE" "Documento prueba PDF" "manual_fabricante" "201"

# Si se obtuvo ID, probar descarga
if [ -n "$LAST_UPLOAD_ID" ]; then
  test_endpoint "GET" "$API_URL/documentos/$LAST_UPLOAD_ID/download" "Descargar archivo recién subido" "" "200"
fi

# Test 7: Subir archivo con extensión no permitida
TMP_TXT="/tmp/test_documento_invalido.txt"
echo "archivo de prueba" > "$TMP_TXT"
test_upload_file "$API_URL/documentos" "$TMP_TXT" "Documento no permitido" "manual_fabricante" "4"

# Test 6: Obtener documento específico
test_endpoint "GET" "$API_URL/documentos/1" "Obtener documento por ID (existente)" "" "200"

# Test 7: Obtener documento inexistente
test_endpoint "GET" "$API_URL/documentos/999" "Obtener documento inexistente" "" "404"

# Test 8: Actualizar documento existente
test_endpoint "PUT" "$API_URL/documentos/1" "Actualizar documento existente" '{"nombre":"Documento Test Actualizado"}' "200"

# Test 9: Actualizar documento inexistente
test_endpoint "PUT" "$API_URL/documentos/999" "Actualizar documento inexistente" '{"nombre":"Documento Inexistente"}' "404"

echo -e "\n${BLUE}🔍 TESTS DE BÚSQUEDA${NC}"
echo "----------------------------------------"

# Test 10: Búsqueda básica con query
test_endpoint "GET" "$API_URL/documentos/buscar?q=bomba" "Búsqueda básica con query" "" "200"

# Test 11: Búsqueda sin parámetro query
test_endpoint "GET" "$API_URL/documentos/buscar" "Búsqueda sin parámetro query (debe fallar)" "" "400"

# Test 12: Búsqueda con múltiples términos
test_endpoint "GET" "$API_URL/documentos/buscar?q=bomba%20agua" "Búsqueda con múltiples términos" "" "200"

echo -e "\n${BLUE}🤖 TESTS DE IA Y CONSULTAS INTERACTIVAS${NC}"
echo "----------------------------------------"

# Test 13: Consulta sobre cantidad de bombas
test_endpoint "POST" "$API_URL/documentos/1/consultar" "Consulta sobre cantidad de bombas" '{"pregunta":"¿Cuántas bombas de agua tiene el sistema?"}' "200"

# Test 14: Consulta sobre presión
test_endpoint "POST" "$API_URL/documentos/1/consultar" "Consulta sobre presión" '{"pregunta":"¿Qué presión tiene esta bomba?"}' "200"

# Test 15: Consulta sobre documento inexistente
test_endpoint "POST" "$API_URL/documentos/999/consultar" "Consulta sobre documento inexistente" '{"pregunta":"Test"}' "404"

# Test 16: Obtener historial de consultas
test_endpoint "GET" "$API_URL/documentos/1/consultas" "Obtener historial de consultas" "" "200"

# Test 17: Obtener estadísticas de consultas
test_endpoint "GET" "$API_URL/consultas/estadisticas" "Obtener estadísticas de consultas" "" "200"

echo -e "\n${BLUE}⚠️  TESTS DE VALIDACIONES Y CASOS EDGE${NC}"
echo "----------------------------------------"

# Test 18: ID de documento inválido (no numérico)
test_endpoint "POST" "$API_URL/documentos/abc/consultar" "Consulta con ID de documento inválido" '{"pregunta":"Test"}' "400"

# Test 19: JSON malformado
test_endpoint "POST" "$API_URL/documentos/1/consultar" "Consulta con JSON malformado" '{"invalid_json":}' "400"

# Test 20: Consulta vacía
test_endpoint "POST" "$API_URL/documentos/1/consultar" "Consulta con pregunta vacía" '{"pregunta":""}' "400"

# Test 21: ID negativo
test_endpoint "GET" "$API_URL/documentos/-1" "Obtener documento con ID negativo" "" "400"

echo -e "\n${BLUE}🧪 TESTS DE FUNCIONALIDAD IA AVANZADA${NC}"
echo "----------------------------------------"

echo -e "${CYAN}🤖 Probando respuestas inteligentes de IA...${NC}"

# Test funcionalidad específica de IA
test_endpoint "POST" "$API_URL/documentos/1/consultar" "Test respuesta IA sobre bombas" '{"pregunta":"¿Cuántas bombas de agua tiene el sistema?"}' "200"

# Analizar respuesta de IA
if echo "$LAST_RESPONSE" | grep -q "bomba" && echo "$LAST_RESPONSE" | grep -q "agua"; then
    echo -e "${GREEN}✅ IA reconoce y responde preguntas sobre bombas de agua${NC}"
else
    echo -e "${YELLOW}⚠️  IA respuesta sobre bombas necesita revisión${NC}"
fi

test_endpoint "POST" "$API_URL/documentos/1/consultar" "Test respuesta IA sobre presión" '{"pregunta":"¿Qué presión tiene esta bomba?"}' "200"

# Verificar confianza de respuesta
if echo "$LAST_RESPONSE" | grep -q '"confianza"'; then
    echo -e "${GREEN}✅ IA proporciona nivel de confianza en respuestas${NC}"
else
    echo -e "${YELLOW}⚠️  Revisar nivel de confianza en respuestas IA${NC}"
fi

echo -e "\n${CYAN}📊 RESUMEN FINAL DE TESTS${NC}"
echo "============================================================="
echo "Total de tests ejecutados: $TOTAL_TESTS"
echo "Tests exitosos: $PASSED_TESTS"
echo "Tests fallidos: $FAILED_TESTS"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}🎉 TODOS LOS TESTS PASARON${NC}"
    echo "Success rate: 100%"
    exit_code=0
else
    echo -e "${RED}⚠️  ALGUNOS TESTS FALLARON${NC}"
    success_rate=$((PASSED_TESTS * 100 / TOTAL_TESTS))
    echo "Success rate: ${success_rate}%"
    exit_code=1
fi

echo ""
echo -e "${CYAN}📋 FUNCIONALIDADES EVALUADAS:${NC}"
echo "• Health check y conectividad básica"
echo "• Gestión completa de documentos (CRUD)"
echo "• Sistema de búsqueda y filtrado"
echo "• Consultas interactivas con IA"
echo "• Validaciones y casos edge"
echo "• Funcionalidad avanzada de IA"

echo ""
echo -e "${CYAN}🔧 PARA DEBUGGING:${NC}"
echo "• Logs del servicio: docker logs documentacion-service"
echo "• Logs de la DB: docker logs documentacion-db"
echo "• Base URL: $BASE_URL"

echo ""
if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}✅ MICROSERVICIO DE DOCUMENTACIÓN FUNCIONANDO CORRECTAMENTE${NC}"
else
    echo -e "${RED}❌ MICROSERVICIO DE DOCUMENTACIÓN NECESITA REVISIÓN${NC}"
fi

exit $exit_code
