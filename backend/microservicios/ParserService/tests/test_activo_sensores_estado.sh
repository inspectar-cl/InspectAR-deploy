#!/bin/bash

# 🔍 TEST ESPECÍFICO - CONSULTA DE ACTIVO CON ESTADO DE SENSORES
# ============================================================
# Este script prueba la nueva funcionalidad de consultar un activo
# y obtener el estado de conexión de todos sus sensores

# Configuración de colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuración de la API
API_URL="http://localhost:8090"
TEST_COUNT=0
PASS_COUNT=0
FAIL_COUNT=0

echo -e "${BLUE}🔍 TEST: CONSULTA DE ACTIVO CON ESTADO DE SENSORES${NC}"
echo "========================================================"
echo -e "${CYAN}🔗 URL Base: $API_URL${NC}"
echo -e "${CYAN}📅 Fecha: $(date)${NC}"
echo ""

# Función para realizar test de endpoint
test_endpoint() {
    local method="$1"
    local endpoint="$2"
    local description="$3"
    local data="$4"
    local expected_status="$5"
    
    TEST_COUNT=$((TEST_COUNT + 1))
    
    echo -e "${CYAN}📍 Test #$TEST_COUNT: $description${NC}"
    echo -e "${YELLOW}$method $endpoint${NC}"
    
    # Realizar petición
    if [ "$method" = "GET" ]; then
        RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" "$endpoint")
    elif [ "$method" = "POST" ]; then
        RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST -H "Content-Type: application/json" -d "$data" "$endpoint")
    fi
    
    # Separar response body y status code
    RESPONSE_BODY=$(echo "$RESPONSE" | sed -E 's/HTTPSTATUS\:[0-9]{3}$//')
    STATUS_CODE=$(echo "$RESPONSE" | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
    
    # Verificar status code
    if [[ "$STATUS_CODE" =~ ^$expected_status ]]; then
        echo -e "${GREEN}✅ PASS - Status: $STATUS_CODE${NC}"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo -e "${RED}❌ FAIL - Status: $STATUS_CODE (esperado: ${expected_status}xx)${NC}"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
    
    # Mostrar respuesta formateada si es JSON
    if echo "$RESPONSE_BODY" | grep -q '"'; then
        echo -e "${CYAN}Response:${NC}"
        echo "$RESPONSE_BODY" | python3 -m json.tool 2>/dev/null || echo "$RESPONSE_BODY"
    else
        echo -e "Response: $RESPONSE_BODY"
    fi
    echo ""
    
    # Guardar última respuesta para análisis
    LAST_RESPONSE="$RESPONSE_BODY"
}

# Verificar disponibilidad del servicio
echo -e "${CYAN}🔍 Verificando disponibilidad del servicio...${NC}"
if timeout 10 curl -s "$API_URL/activo" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Servicio Parser disponible${NC}"
    echo ""
else
    echo -e "${RED}❌ Servicio Parser no disponible en puerto 8090${NC}"
    echo -e "${YELLOW}💡 Asegúrate de que el servicio esté ejecutándose${NC}"
    exit 1
fi

echo -e "${BLUE}🏭 TESTS DE CONSULTA DE ACTIVO CON ESTADO DE SENSORES${NC}"
echo "--------------------------------------------------------"

# Test 1: Crear un activo de prueba
ACTIVO_TEST='{
    "activo_id": "TEST_SENSOR_STATUS_001",
    "nombre": "Caldera de Prueba - Estado Sensores",
    "ubicacion": "Planta Test",
    "estado": "activo",
    "id_edificio": "edificio_001",
    "sensores": [
        {
            "sensor_id": "TEMP_TEST_001",
            "tipo": "temperatura",
            "unidad": "°C"
        },
        {
            "sensor_id": "PRESS_TEST_001", 
            "tipo": "presion",
            "unidad": "bar"
        },
        {
            "sensor_id": "VIBR_TEST_001",
            "tipo": "vibracion", 
            "unidad": "Hz"
        }
    ]
}'

test_endpoint "POST" "$API_URL/activo" "Crear activo de prueba con sensores" "$ACTIVO_TEST" "200"

# Test 2: Consultar el activo con estado de sensores (sensores nunca conectados)
test_endpoint "GET" "$API_URL/activo/TEST_SENSOR_STATUS_001/sensores/estado" "Consultar activo con estado de sensores (sin datos)" "" "200"

# Analizar respuesta del test anterior
if echo "$LAST_RESPONSE" | grep -q '"never_connected"'; then
    echo -e "${GREEN}✅ Correctamente identifica sensores nunca conectados${NC}"
else
    echo -e "${YELLOW}⚠️  No detectó sensores nunca conectados${NC}"
fi

# Test 3: Simular datos de un sensor para cambiar su estado
LECTURA_1='{
    "sensor_id": "TEMP_TEST_001",
    "valor": 25.5,
    "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
}'

test_endpoint "POST" "$API_URL/lectura" "Enviar datos de sensor TEMP_TEST_001" "$LECTURA_1" "200"

# Esperar un momento para que se procese
sleep 2

# Test 4: Consultar nuevamente el activo (ahora debería tener un sensor conectado)
test_endpoint "GET" "$API_URL/activo/TEST_SENSOR_STATUS_001/sensores/estado" "Consultar activo después de enviar datos a un sensor" "" "200"

# Analizar si detectó el sensor conectado
if echo "$LAST_RESPONSE" | grep -q '"connected"'; then
    echo -e "${GREEN}✅ Correctamente detecta sensor conectado después de recibir datos${NC}"
else
    echo -e "${YELLOW}⚠️  No detectó el sensor como conectado${NC}"
fi

# Verificar estadísticas en la respuesta
if echo "$LAST_RESPONSE" | grep -q '"resumen"'; then
    echo -e "${GREEN}✅ Incluye resumen de estadísticas de sensores${NC}"
else
    echo -e "${YELLOW}⚠️  No incluye resumen de estadísticas${NC}"
fi

# Test 5: Enviar datos a otro sensor
LECTURA_2='{
    "sensor_id": "PRESS_TEST_001",
    "valor": 1.2,
    "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
}'

test_endpoint "POST" "$API_URL/lectura" "Enviar datos de sensor PRESS_TEST_001" "$LECTURA_2" "200"

sleep 2

# Test 6: Consultar por tercera vez (ahora con 2 sensores conectados)
test_endpoint "GET" "$API_URL/activo/TEST_SENSOR_STATUS_001/sensores/estado" "Consultar activo con 2 sensores conectados" "" "200"

# Test 7: Consultar activo inexistente
test_endpoint "GET" "$API_URL/activo/ACTIVO_INEXISTENTE/sensores/estado" "Consultar activo inexistente" "" "404"

# Test 8: Verificar estructura de la respuesta
echo -e "${CYAN}🔍 Analizando estructura de la respuesta...${NC}"

# Obtener la respuesta del último test exitoso
test_endpoint "GET" "$API_URL/activo/TEST_SENSOR_STATUS_001/sensores/estado" "Verificación final de estructura" "" "200"

# Verificar campos obligatorios
campos_requeridos=("id" "activo_id" "nombre" "ubicacion" "estado" "sensores" "total_sensores" "resumen")
for campo in "${campos_requeridos[@]}"; do
    if echo "$LAST_RESPONSE" | grep -q "\"$campo\""; then
        echo -e "${GREEN}✅ Campo '$campo' presente${NC}"
    else
        echo -e "${RED}❌ Campo '$campo' faltante${NC}"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
done

# Verificar campos de sensor
campos_sensor=("sensor_id" "tipo" "unidad" "estado" "is_active")
echo -e "\n${CYAN}🔍 Verificando campos de sensores...${NC}"
for campo in "${campos_sensor[@]}"; do
    if echo "$LAST_RESPONSE" | grep -q "\"$campo\""; then
        echo -e "${GREEN}✅ Campo de sensor '$campo' presente${NC}"
    else
        echo -e "${RED}❌ Campo de sensor '$campo' faltante${NC}"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
done

# Mostrar resumen final
echo -e "\n${BLUE}📊 RESUMEN FINAL DE TESTS${NC}"
echo "==========================================="
echo "Total de tests ejecutados: $TEST_COUNT"
echo "Tests exitosos: $PASS_COUNT"
echo "Tests fallidos: $FAIL_COUNT"

if [ $FAIL_COUNT -eq 0 ]; then
    SUCCESS_RATE=100
    echo -e "${GREEN}🎉 TODOS LOS TESTS PASARON${NC}"
    echo "Success rate: ${SUCCESS_RATE}%"
    echo ""
    echo -e "${GREEN}✅ NUEVA FUNCIONALIDAD TRABAJANDO CORRECTAMENTE${NC}"
    echo -e "${CYAN}📍 Endpoint: GET /activo/:activo_id/sensores/estado${NC}"
    echo -e "${CYAN}📋 Funcionalidad: Consulta activo con estado de sensores${NC}"
    exit 0
else
    if [ $TEST_COUNT -gt 0 ]; then
        SUCCESS_RATE=$(( (PASS_COUNT * 100) / TEST_COUNT ))
    else
        SUCCESS_RATE=0
    fi
    echo -e "${YELLOW}⚠️  ALGUNOS TESTS FALLARON${NC}"
    echo "Success rate: ${SUCCESS_RATE}%"
    echo ""
    echo -e "${RED}❌ NUEVA FUNCIONALIDAD NECESITA REVISIÓN${NC}"
    exit 1
fi
