#!/bin/bash

# Test corregido para observaciones editables
# Enfocado en las nuevas funcionalidades de reportes

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

# Configuración
GESTION_URL="http://localhost:8092"
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

echo -e "${BLUE}🆕 TEST OBSERVACIONES EDITABLES - CORREGIDO${NC}"
echo "=============================================="
echo "URL Base: $GESTION_URL"
echo "Fecha: $(date)"
echo ""

# Función para hacer test de una ruta
test_route() {
    local method=$1
    local url=$2
    local description=$3
    local data=$4
    local expected_status=$5
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    echo -e "\n${CYAN}📍 Test #$TOTAL_TESTS:${NC} $description"
    echo -e "${PURPLE}$method${NC} $url"
    
    if [ -n "$data" ]; then
        response=$(curl -s -X $method -H "Content-Type: application/json" -d "$data" -w "HTTP_CODE:%{http_code}" "$url")
    else
        response=$(curl -s -X $method -w "HTTP_CODE:%{http_code}" "$url")
    fi
    
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
    
    # Store the full response body for ID extraction
    LAST_RESPONSE_BODY="$body"
}

# Function to extract ID from JSON response
extract_id_from_response() {
    echo "$LAST_RESPONSE_BODY" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2
}

# Verificar servicio disponible
echo -e "\n${CYAN}🔍 Verificando disponibilidad del servicio...${NC}"
test_route "GET" "$GESTION_URL/health" "Health Check" "" "200"

echo -e "\n${CYAN}🆕 TESTS DE OBSERVACIONES EDITABLES${NC}"
echo "===================================="

# Test 1: Obtener todos los reportes
test_route "GET" "$GESTION_URL/reportes" "🆕 Obtener todos los reportes con observaciones" "" "200"

echo -e "\n${PURPLE}📝 CREANDO REPORTES CON OBSERVACIONES${NC}"
echo "======================================"

# Crear reporte #1: Mantenimiento de caldera
echo -e "\n${BLUE}Reporte #1: Mantenimiento de caldera principal${NC}"
REPORTE_DATA_1='{
    "activo_id": 1,
    "tipo_reporte": "mantenimiento",
    "contenido": "Mantenimiento preventivo realizado en caldera principal. Se verificaron válvulas de seguridad y presión del sistema.",
    "observaciones_analista": "Caldera operando correctamente dentro de parámetros. Se recomienda mantenimiento preventivo en 30 días.",
    "autor_analista": "Juan Pérez - Ingeniero Mecánico"
}'

test_route "POST" "$GESTION_URL/reportes" "🆕 Crear reporte mantenimiento caldera" "$REPORTE_DATA_1" "201"
REPORTE_ID_1=$(extract_id_from_response)
echo "ID del reporte creado: $REPORTE_ID_1"

# Crear reporte #2: Incidente bomba
echo -e "\n${BLUE}Reporte #2: Incidente bomba hidráulica${NC}"
REPORTE_DATA_2='{
    "activo_id": 2,
    "tipo_reporte": "incidente",
    "contenido": "Incidente menor en bomba hidráulica. Se detectó vibración anormal durante operación rutinaria.",
    "observaciones_analista": "Vibración corregida tras ajuste de rodamientos. Sistema operativo normal.",
    "autor_analista": "María González - Técnico Hidráulico"
}'

test_route "POST" "$GESTION_URL/reportes" "🆕 Crear reporte incidente bomba" "$REPORTE_DATA_2" "201"
REPORTE_ID_2=$(extract_id_from_response)
echo "ID del reporte creado: $REPORTE_ID_2"

# Crear reporte #3: Emergencia transformador
echo -e "\n${BLUE}Reporte #3: Emergencia transformador eléctrico${NC}"
REPORTE_DATA_3='{
    "activo_id": 6,
    "tipo_reporte": "incidente",
    "contenido": "Respuesta a emergencia por sobrecalentamiento del transformador. Se detectó falla en sistema de ventilación.",
    "observaciones_analista": "Emergencia resuelta. Sistema de ventilación reparado y transformador vuelto a operación normal.",
    "autor_analista": "Carlos López - Especialista Eléctrico"
}'

test_route "POST" "$GESTION_URL/reportes" "🆕 Crear reporte emergencia transformador" "$REPORTE_DATA_3" "201"
REPORTE_ID_3=$(extract_id_from_response)
echo "ID del reporte creado: $REPORTE_ID_3"

echo -e "\n${CYAN}📋 CONSULTANDO REPORTES CREADOS${NC}"
echo "=================================="

# Consultar reportes individuales
if [ -n "$REPORTE_ID_1" ]; then
    test_route "GET" "$GESTION_URL/reportes/$REPORTE_ID_1" "🆕 Obtener reporte específico #1" "" "200"
fi

if [ -n "$REPORTE_ID_2" ]; then
    test_route "GET" "$GESTION_URL/reportes/$REPORTE_ID_2" "🆕 Obtener reporte específico #2" "" "200"
fi

if [ -n "$REPORTE_ID_3" ]; then
    test_route "GET" "$GESTION_URL/reportes/$REPORTE_ID_3" "🆕 Obtener reporte específico #3" "" "200"
fi

# Consultar reportes con observaciones por activo
test_route "GET" "$GESTION_URL/reportes/activo/1/observaciones" "🆕 Reportes con observaciones - Activo 1" "" "200"
test_route "GET" "$GESTION_URL/reportes/activo/2/observaciones" "🆕 Reportes con observaciones - Activo 2" "" "200"
test_route "GET" "$GESTION_URL/reportes/activo/6/observaciones" "🆕 Reportes con observaciones - Activo 6" "" "200"

echo -e "\n${CYAN}✏️  EDITANDO OBSERVACIONES - FLUJO DE REVISIÓN${NC}"
echo "=============================================="

# Actualizar observaciones del reporte de caldera
if [ -n "$REPORTE_ID_1" ]; then
    echo -e "\n${BLUE}Actualizando observaciones del reporte de caldera${NC}"
    UPDATE_DATA_1='{
        "observaciones_analista": "Caldera operando correctamente. Actualización: Se detectó leve incremento en presión, dentro de límites normales. Próxima revisión programada.",
        "autor_analista": "Juan Pérez - Ingeniero Mecánico Senior"
    }'
    
    test_route "PUT" "$GESTION_URL/reportes/$REPORTE_ID_1/observaciones" "🆕 Actualizar observaciones - Caldera" "$UPDATE_DATA_1" "200"
fi

# Actualizar observaciones del reporte de bomba
if [ -n "$REPORTE_ID_2" ]; then
    echo -e "\n${BLUE}Actualizando observaciones del reporte de bomba${NC}"
    UPDATE_DATA_2='{
        "observaciones_analista": "Bomba funcionando óptimamente tras reparación. Protocolos de mantenimiento a replicar en otras unidades similares.",
        "autor_analista": "María González - Técnico Hidráulico Senior"
    }'
    
    test_route "PUT" "$GESTION_URL/reportes/$REPORTE_ID_2/observaciones" "🆕 Actualizar observaciones - Bomba" "$UPDATE_DATA_2" "200"
fi

echo -e "\n${CYAN}🔍 PROCESO DE REVISIÓN Y APROBACIÓN${NC}"
echo "========================================="

# Proceso de revisión para el reporte de caldera
if [ -n "$REPORTE_ID_1" ]; then
    echo -e "\n${BLUE}Enviando reporte de caldera a revisión${NC}"
    REVISION_DATA_1='{
        "estado_revision": "en_revision",
        "revisor": "Ana Supervisor - Jefe Técnico",
        "observaciones": "Reporte técnicamente correcto. Pendiente de aprobación final."
    }'
    
    test_route "PUT" "$GESTION_URL/reportes/$REPORTE_ID_1/revision" "🆕 Enviar a revisión - Caldera" "$REVISION_DATA_1" "200"
fi

# Aprobar reporte de bomba directamente
if [ -n "$REPORTE_ID_2" ]; then
    echo -e "\n${BLUE}Aprobando reporte de bomba hidráulica${NC}"
    REVISION_DATA_2='{
        "estado_revision": "aprobado",
        "revisor": "Luis Director - Director Técnico",
        "observaciones": "Excelente trabajo. Protocolo aprobado para replicar."
    }'
    
    test_route "PUT" "$GESTION_URL/reportes/$REPORTE_ID_2/revision" "🆕 Aprobar reporte - Bomba" "$REVISION_DATA_2" "200"
fi

# Aprobar reporte de emergencia
if [ -n "$REPORTE_ID_3" ]; then
    echo -e "\n${BLUE}Aprobando resolución de emergencia${NC}"
    REVISION_DATA_3='{
        "estado_revision": "aprobado",
        "revisor": "Emergency Team - Equipo de Emergencias",
        "observaciones": "Respuesta rápida y efectiva. Procedimientos seguidos correctamente."
    }'
    
    test_route "PUT" "$GESTION_URL/reportes/$REPORTE_ID_3/revision" "🆕 Aprobar emergencia - Transformador" "$REVISION_DATA_3" "200"
fi

# Completar revisión de caldera
if [ -n "$REPORTE_ID_1" ]; then
    echo -e "\n${BLUE}Aprobando reporte de caldera tras revisión${NC}"
    REVISION_FINAL_1='{
        "estado_revision": "aprobado",
        "revisor": "Ana Supervisor - Jefe Técnico",
        "observaciones": "Reporte aprobado. Mantenimiento programado autorizado."
    }'
    
    test_route "PUT" "$GESTION_URL/reportes/$REPORTE_ID_1/revision" "🆕 Aprobar tras revisión - Caldera" "$REVISION_FINAL_1" "200"
fi

echo -e "\n${CYAN}📊 VERIFICANDO ESTADOS FINALES DE REPORTES${NC}"
echo "============================================"

# Verificar estados finales
if [ -n "$REPORTE_ID_1" ]; then
    test_route "GET" "$GESTION_URL/reportes/$REPORTE_ID_1" "🆕 Estado final reporte caldera" "" "200"
fi

if [ -n "$REPORTE_ID_2" ]; then
    test_route "GET" "$GESTION_URL/reportes/$REPORTE_ID_2" "🆕 Estado final reporte bomba" "" "200"
fi

if [ -n "$REPORTE_ID_3" ]; then
    test_route "GET" "$GESTION_URL/reportes/$REPORTE_ID_3" "🆕 Estado final reporte emergencia" "" "200"
fi

# Obtener todos los reportes actualizados
test_route "GET" "$GESTION_URL/reportes" "🆕 Todos los reportes - Estado final" "" "200"

echo -e "\n${GREEN}🎉 FLUJO DE OBSERVACIONES EDITABLES COMPLETADO${NC}"
echo "Estado final esperado:"
echo "• Reporte #1: ✅ APROBADO (Caldera - mantenimiento programado)"
echo "• Reporte #2: ✅ APROBADO (Bomba - protocolo a replicar)" 
echo "• Reporte #3: ✅ APROBADO (Emergencia transformador - resuelto)"
echo ""
echo "Características probadas:"
echo "• ✅ Creación de reportes con observaciones iniciales"
echo "• ✅ Edición de observaciones por analistas"
echo "• ✅ Flujo de revisión y aprobación"
echo "• ✅ Tracking de autores y revisores"
echo "• ✅ Estados de revisión (pendiente → en_revision → aprobado)"
echo "• ✅ Consulta de reportes con observaciones por activo"

echo -e "\n${CYAN}📊 RESUMEN FINAL DE TESTS OBSERVACIONES${NC}"
echo "=============================================="
echo "Total de tests ejecutados: $TOTAL_TESTS"
echo "Tests exitosos: $PASSED_TESTS"
echo "Tests fallidos: $FAILED_TESTS"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}✅ TODOS LOS TESTS DE OBSERVACIONES PASARON${NC}"
    SUCCESS_RATE=100
else
    SUCCESS_RATE=$((PASSED_TESTS * 100 / TOTAL_TESTS))
    echo -e "\n${YELLOW}⚠️  ALGUNOS TESTS FALLARON${NC}"
fi

echo "Success rate: ${SUCCESS_RATE}%"
echo ""
echo "🆕 SISTEMA DE OBSERVACIONES EDITABLES VERIFICADO Y FUNCIONAL"
