#!/bin/bash

# Script de pruebas para el microservicio de documentación
# Asegúrate de que el microservicio esté corriendo en el puerto 8092

BASE_URL="http://localhost:8092"
API_URL="$BASE_URL/api/v1"

echo "🧪 EJECUTANDO PRUEBAS DEL MICROSERVICIO DE DOCUMENTACIÓN"
echo "=================================================="
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

# Función para hacer petición HTTP y extraer código de estado
make_request() {
    local method="$1"
    local url="$2"
    local data="$3"
    local content_type="$4"
    
    if [ -n "$data" ] && [ -n "$content_type" ]; then
        curl -s -o /dev/null -w "%{http_code}" -X "$method" -H "Content-Type: $content_type" -d "$data" "$url"
    elif [ -n "$data" ]; then
        curl -s -o /dev/null -w "%{http_code}" -X "$method" -d "$data" "$url"
    else
        curl -s -o /dev/null -w "%{http_code}" -X "$method" "$url"
    fi
}

echo -e "${BLUE}1. HEALTH CHECK${NC}"
echo "=================="
status=$(make_request "GET" "$BASE_URL/health")
show_test_result "Health Check" "$status" "200"
echo ""

echo -e "${BLUE}2. GESTIÓN DE DOCUMENTOS${NC}"
echo "========================="

# GET /api/v1/documentos - Listar documentos
status=$(make_request "GET" "$API_URL/documentos")
show_test_result "Listar todos los documentos" "$status" "200"

# GET /api/v1/documentos?activo_id=1 - Filtrar por activo
status=$(make_request "GET" "$API_URL/documentos?activo_id=1")
show_test_result "Filtrar documentos por activo_id" "$status" "200"

# GET /api/v1/documentos?categoria=manual_fabricante - Filtrar por categoría
status=$(make_request "GET" "$API_URL/documentos?categoria=manual_fabricante")
show_test_result "Filtrar documentos por categoría" "$status" "200"

# GET /api/v1/documentos?solo_fichas_tecnicas=true - Solo fichas técnicas
status=$(make_request "GET" "$API_URL/documentos?solo_fichas_tecnicas=true")
show_test_result "Filtrar solo fichas técnicas" "$status" "200"

# GET /api/v1/documentos/1 - Obtener documento específico
status=$(make_request "GET" "$API_URL/documentos/1")
show_test_result "Obtener documento por ID (existente)" "$status" "200"

# GET /api/v1/documentos/999 - Documento inexistente
status=$(make_request "GET" "$API_URL/documentos/999")
show_test_result "Obtener documento inexistente" "$status" "404"

# GET /api/v1/documentos/abc - ID inválido
status=$(make_request "GET" "$API_URL/documentos/abc")
show_test_result "Obtener documento con ID inválido" "$status" "400"

# PUT /api/v1/documentos/1 - Actualizar documento
update_data='{"nombre":"Documento Actualizado","descripcion":"Nueva descripción"}'
status=$(make_request "PUT" "$API_URL/documentos/1" "$update_data" "application/json")
show_test_result "Actualizar documento existente" "$status" "200"

# PUT /api/v1/documentos/999 - Actualizar documento inexistente
status=$(make_request "PUT" "$API_URL/documentos/999" "$update_data" "application/json")
show_test_result "Actualizar documento inexistente" "$status" "404"

# DELETE /api/v1/documentos/1 - Eliminar documento
status=$(make_request "DELETE" "$API_URL/documentos/1")
show_test_result "Eliminar documento existente" "$status" "200"

# DELETE /api/v1/documentos/999 - Eliminar documento inexistente
status=$(make_request "DELETE" "$API_URL/documentos/999")
show_test_result "Eliminar documento inexistente" "$status" "404"

# GET /api/v1/documentos/1/descargar - Descargar documento
status=$(make_request "GET" "$API_URL/documentos/1/descargar")
show_test_result "Descargar documento existente" "$status" "200"

# GET /api/v1/documentos/999/descargar - Descargar documento inexistente
status=$(make_request "GET" "$API_URL/documentos/999/descargar")
show_test_result "Descargar documento inexistente" "$status" "404"

echo ""

echo -e "${BLUE}3. BÚSQUEDA DE DOCUMENTOS${NC}"
echo "============================"

# GET /api/v1/documentos/buscar?q=bomba - Búsqueda básica
status=$(make_request "GET" "$API_URL/documentos/buscar?q=bomba")
show_test_result "Búsqueda básica con query" "$status" "200"

# GET /api/v1/documentos/buscar?q=manual&categoria=manual_fabricante - Búsqueda con filtros
status=$(make_request "GET" "$API_URL/documentos/buscar?q=manual&categoria=manual_fabricante")
show_test_result "Búsqueda con filtros adicionales" "$status" "200"

# GET /api/v1/documentos/buscar - Búsqueda sin query
status=$(make_request "GET" "$API_URL/documentos/buscar")
show_test_result "Búsqueda sin parámetro query (debe fallar)" "$status" "400"

echo ""

echo -e "${BLUE}4. CONSULTAS INTERACTIVAS CON IA${NC}"
echo "===================================="

# POST /api/v1/documentos/1/consultar - Realizar consulta
consulta_data='{"pregunta":"¿Qué presión de agua tiene esta bomba?"}'
status=$(make_request "POST" "$API_URL/documentos/1/consultar" "$consulta_data" "application/json")
show_test_result "Realizar consulta sobre documento" "$status" "200"

# POST /api/v1/documentos/999/consultar - Consulta sobre documento inexistente
status=$(make_request "POST" "$API_URL/documentos/999/consultar" "$consulta_data" "application/json")
show_test_result "Consulta sobre documento inexistente" "$status" "404"

# POST /api/v1/documentos/1/consultar - Consulta sin pregunta
consulta_vacia='{"pregunta":""}'
status=$(make_request "POST" "$API_URL/documentos/1/consultar" "$consulta_vacia" "application/json")
show_test_result "Consulta con pregunta vacía (debe fallar)" "$status" "400"

# GET /api/v1/documentos/1/consultas - Obtener historial de consultas
status=$(make_request "GET" "$API_URL/documentos/1/consultas")
show_test_result "Obtener historial de consultas" "$status" "200"

# GET /api/v1/documentos/999/consultas - Historial de documento inexistente
status=$(make_request "GET" "$API_URL/documentos/999/consultas")
show_test_result "Historial de documento inexistente" "$status" "404"

# GET /api/v1/consultas/estadisticas - Estadísticas generales
status=$(make_request "GET" "$API_URL/consultas/estadisticas")
show_test_result "Obtener estadísticas de consultas" "$status" "200"

echo ""

echo -e "${BLUE}5. VALIDACIONES Y CASOS EDGE${NC}"
echo "================================="

# POST /api/v1/documentos/abc/consultar - ID inválido
status=$(make_request "POST" "$API_URL/documentos/abc/consultar" "$consulta_data" "application/json")
show_test_result "Consulta con ID de documento inválido" "$status" "400"

# POST /api/v1/documentos/1/consultar - JSON malformado
json_malformado='{"pregunta":"¿Qué presión"'
status=$(make_request "POST" "$API_URL/documentos/1/consultar" "$json_malformado" "application/json")
show_test_result "Consulta con JSON malformado" "$status" "400"

# GET /api/v1/documentos?limit=abc - Parámetro limit inválido
status=$(make_request "GET" "$API_URL/documentos?limit=abc")
show_test_result "Listar con parámetro limit inválido" "$status" "400"

# GET /api/v1/documentos?categoria=invalida - Categoría inválida
status=$(make_request "GET" "$API_URL/documentos?categoria=categoria_que_no_existe")
show_test_result "Filtrar con categoría inválida" "$status" "400"

echo ""

echo -e "${YELLOW}6. PRUEBAS DE CREACIÓN DE DOCUMENTOS${NC}"
echo "========================================"
echo "⚠️  Las siguientes pruebas requieren archivos reales:"
echo ""
echo "Para probar la creación de documentos, ejecuta manualmente:"
echo ""
echo -e "${GREEN}# Crear un archivo PDF de prueba${NC}"
echo "echo 'Contenido de prueba' > test_document.pdf"
echo ""
echo -e "${GREEN}# Subir documento${NC}"
echo "curl -X POST $API_URL/documentos \\"
echo "  -F 'archivo=@test_document.pdf' \\"
echo "  -F 'activo_id=1' \\"
echo "  -F 'nombre=Documento de Prueba' \\"
echo "  -F 'descripcion=Descripción de prueba' \\"
echo "  -F 'categoria=manual_fabricante' \\"
echo "  -F 'subido_por=test@empresa.com' \\"
echo "  -F 'palabras_clave=prueba,test'"
echo ""

echo -e "${BLUE}7. RESUMEN DE ENDPOINTS DISPONIBLES${NC}"
echo "====================================="
echo ""
echo -e "${GREEN}Gestión de Documentos:${NC}"
echo "  GET    $API_URL/documentos                     - Listar documentos"
echo "  GET    $API_URL/documentos/{id}                - Obtener documento por ID"
echo "  POST   $API_URL/documentos                     - Crear nuevo documento"
echo "  PUT    $API_URL/documentos/{id}                - Actualizar documento"
echo "  DELETE $API_URL/documentos/{id}                - Eliminar documento"
echo "  GET    $API_URL/documentos/{id}/descargar      - Descargar archivo"
echo ""
echo -e "${GREEN}Búsqueda:${NC}"
echo "  GET    $API_URL/documentos/buscar              - Búsqueda de texto completo"
echo ""
echo -e "${GREEN}Consultas IA:${NC}"
echo "  POST   $API_URL/documentos/{id}/consultar      - Realizar consulta"
echo "  GET    $API_URL/documentos/{id}/consultas      - Historial de consultas"
echo "  GET    $API_URL/consultas/estadisticas         - Estadísticas generales"
echo ""
echo -e "${GREEN}Otros:${NC}"
echo "  GET    $BASE_URL/health                        - Health check"
echo ""

echo -e "${YELLOW}Parámetros de consulta disponibles para GET /documentos:${NC}"
echo "  ?activo_id=<int>           - Filtrar por ID de activo"
echo "  ?tecnico_id=<int>          - Filtrar por ID de técnico"
echo "  ?categoria=<string>        - Filtrar por categoría"
echo "  ?solo_fichas_tecnicas=true - Solo fichas técnicas"
echo "  ?limit=<int>               - Limitar número de resultados"
echo "  ?offset=<int>              - Offset para paginación"
echo ""

echo -e "${YELLOW}Categorías válidas:${NC}"
echo "  - ficha_tecnica"
echo "  - manual_fabricante"
echo "  - reporte_mantenimiento"
echo "  - diagnostico"
echo "  - certificacion"
echo ""

echo "✅ Pruebas completadas!"
echo ""
echo "Para ejecutar este script:"
echo "  chmod +x test_routes.sh"
echo "  ./test_routes.sh"
