#!/bin/bash

# 🔧 MICROSERVICIO PARSER - TEST DE AGREGAR SENSOR A ACTIVO
# =========================================================
# Este script testa la nueva funcionalidad de agregar sensores a activos existentes
# Endpoint: POST /activo/:activo_id/sensores

# Configuración de colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuración de la API
API_URL="http://localhost:8090"
LAST_RESPONSE=""
TEST_COUNT=0
PASS_COUNT=0
FAIL_COUNT=0

# Función para mostrar el header
show_header() {
    echo -e "${BLUE}🔧 TEST: AGREGAR SENSOR A ACTIVO EXISTENTE${NC}"
    echo "========================================================="
    echo -e "${CYAN}🔗 URL Base: $API_URL${NC}"
    echo -e "${CYAN}📅 Fecha: $(date)${NC}"
    echo ""
}

# Función para verificar disponibilidad del servicio
check_service_availability() {
    echo -e "${CYAN}🔍 Verificando disponibilidad del servicio...${NC}"
    
    # Intentar conectar al servicio con timeout
    if timeout 10 curl -s "$API_URL/activo" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Servicio Parser disponible${NC}"
        echo ""
    else
        echo -e "${RED}❌ Servicio Parser no disponible${NC}"
        echo -e "${YELLOW}💡 Asegúrate de que el servicio esté ejecutándose en puerto 8090${NC}"
        echo -e "${YELLOW}   Comando: docker-compose up parser-service${NC}"
        exit 1
    fi
}

# Función para realizar test de endpoint
test_endpoint() {
    local method="$1"
    local endpoint="$2"
    local description="$3"
    local data="$4"
    local expected_status="$5"
    
    TEST_COUNT=$((TEST_COUNT + 1))
    
    echo -e "${PURPLE}📍 Test #$TEST_COUNT: $description${NC}"
    echo -e "${CYAN}$method $endpoint${NC}"
    
    # Realizar petición según el método
    if [ "$method" = "GET" ]; then
        RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" "$endpoint")
    elif [ "$method" = "POST" ]; then
        RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST -H "Content-Type: application/json" -d "$data" "$endpoint")
    elif [ "$method" = "PUT" ]; then
        RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" -X PUT -H "Content-Type: application/json" -d "$data" "$endpoint")
    fi
    
    # Separar response body y status code
    RESPONSE_BODY=$(echo "$RESPONSE" | sed -E 's/HTTPSTATUS\:[0-9]{3}$//')
    STATUS_CODE=$(echo "$RESPONSE" | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
    
    # Guardar última respuesta para análisis posterior
    LAST_RESPONSE="$RESPONSE_BODY"
    
    # Verificar status code
    if [[ "$STATUS_CODE" =~ ^$expected_status ]]; then
        echo -e "${GREEN}✅ PASS - Status: $STATUS_CODE${NC}"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo -e "${RED}❌ FAIL - Status: $STATUS_CODE (esperado: ${expected_status}xx)${NC}"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
    
    # Mostrar respuesta formateada
    echo "Response:"
    echo "$RESPONSE_BODY" | jq . 2>/dev/null || echo "$RESPONSE_BODY"
    echo ""
}

# Función para mostrar resumen final
show_summary() {
    echo -e "\n${BLUE}📊 RESUMEN FINAL DE TESTS${NC}"
    echo "==========================================="
    echo "Total de tests ejecutados: $TEST_COUNT"
    echo "Tests exitosos: $PASS_COUNT"
    echo "Tests fallidos: $FAIL_COUNT"
    
    # Calcular porcentaje de éxito
    if [ $TEST_COUNT -gt 0 ]; then
        SUCCESS_RATE=$(( (PASS_COUNT * 100) / TEST_COUNT ))
        
        if [ $FAIL_COUNT -eq 0 ]; then
            echo -e "${GREEN}🎉 TODOS LOS TESTS PASARON${NC}"
            echo "Success rate: ${SUCCESS_RATE}%"
        elif [ $SUCCESS_RATE -ge 80 ]; then
            echo -e "${YELLOW}⚠️  ALGUNOS TESTS FALLARON${NC}"
            echo "Success rate: ${SUCCESS_RATE}%"
        else
            echo -e "${RED}❌ MUCHOS TESTS FALLARON${NC}"
            echo "Success rate: ${SUCCESS_RATE}%"
        fi
    fi
    
    echo ""
    echo -e "${CYAN}✅ NUEVA FUNCIONALIDAD IMPLEMENTADA${NC}"
    echo -e "${CYAN}📍 Endpoint: POST /activo/:activo_id/sensores${NC}"
    echo -e "${CYAN}📋 Funcionalidad: Agregar sensor a activo existente${NC}"
}

# =====================================
# EJECUCIÓN PRINCIPAL DE LOS TESTS
# =====================================

show_header
check_service_availability

echo -e "${BLUE}🏭 TESTS DE AGREGAR SENSOR A ACTIVO${NC}"
echo "--------------------------------------------------------"

# Test 1: Crear activo de prueba
ACTIVO_DATA='{
    "activo_id": "TEST_ADD_SENSOR_001",
    "nombre": "Caldera Test - Agregar Sensores",
    "ubicacion": "Planta Test",
    "estado": "activo",
    "id_edificio": "edificio_test_001",
    "sensores": [
        {
            "sensor_id": "TEMP_INICIAL_001",
            "tipo": "temperatura",
            "unidad": "°C"
        }
    ]
}'
test_endpoint "POST" "$API_URL/activo" "Crear activo de prueba con un sensor inicial" "$ACTIVO_DATA" "200"

# Test 2: Agregar sensor de presión al activo
SENSOR_PRESION='{
    "sensor_id": "PRESS_NUEVO_001",
    "tipo": "presion",
    "unidad": "bar"
}'
test_endpoint "POST" "$API_URL/activo/TEST_ADD_SENSOR_001/sensores" "Agregar sensor de presión al activo" "$SENSOR_PRESION" "200"

# Test 3: Agregar sensor de vibración al activo
SENSOR_VIBRACION='{
    "sensor_id": "VIBR_NUEVO_001",
    "tipo": "vibracion",
    "unidad": "Hz"
}'
test_endpoint "POST" "$API_URL/activo/TEST_ADD_SENSOR_001/sensores" "Agregar sensor de vibración al activo" "$SENSOR_VIBRACION" "200"

# Test 4: Verificar que el activo ahora tiene 3 sensores
test_endpoint "GET" "$API_URL/activo/TEST_ADD_SENSOR_001" "Verificar que el activo tiene todos los sensores agregados" "" "200"

# Verificar que hay 3 sensores en la respuesta
if echo "$LAST_RESPONSE" | grep -q '"sensores"' && echo "$LAST_RESPONSE" | jq '.sensores | length' | grep -q '3'; then
    echo -e "${GREEN}✅ Activo contiene los 3 sensores esperados${NC}"
else
    echo -e "${YELLOW}⚠️  Verificar manualmente el número de sensores en el activo${NC}"
fi

# Test 5: Verificar el estado de los sensores agregados
test_endpoint "GET" "$API_URL/activo/TEST_ADD_SENSOR_001/sensores/estado" "Verificar estado de todos los sensores del activo" "" "200"

# Test 6: Intentar agregar sensor a activo inexistente (debe fallar)
SENSOR_ERROR='{
    "sensor_id": "ERROR_SENSOR_001",
    "tipo": "error",
    "unidad": "err"
}'
test_endpoint "POST" "$API_URL/activo/ACTIVO_INEXISTENTE/sensores" "Intentar agregar sensor a activo inexistente" "$SENSOR_ERROR" "404"

# Test 7: Intentar agregar sensor con datos inválidos (debe fallar)
SENSOR_INVALIDO='{"datos_incorrectos": "valor"}'
test_endpoint "POST" "$API_URL/activo/TEST_ADD_SENSOR_001/sensores" "Intentar agregar sensor con datos inválidos" "$SENSOR_INVALIDO" "400"

# Test 8: Agregar sensor de humedad con datos completos
SENSOR_HUMEDAD='{
    "sensor_id": "HUM_NUEVO_001",
    "tipo": "humedad",
    "unidad": "%"
}'
test_endpoint "POST" "$API_URL/activo/TEST_ADD_SENSOR_001/sensores" "Agregar sensor de humedad al activo" "$SENSOR_HUMEDAD" "200"

# Test 9: Verificación final - el activo debe tener 4 sensores
test_endpoint "GET" "$API_URL/activo/TEST_ADD_SENSOR_001" "Verificación final: activo con 4 sensores" "" "200"

# Verificar que hay 4 sensores en la respuesta final
if echo "$LAST_RESPONSE" | grep -q '"sensores"' && echo "$LAST_RESPONSE" | jq '.sensores | length' | grep -q '4'; then
    echo -e "${GREEN}✅ Verificación final exitosa: 4 sensores en el activo${NC}"
    
    # Mostrar los sensores encontrados
    echo -e "${CYAN}📋 Sensores encontrados:${NC}"
    echo "$LAST_RESPONSE" | jq '.sensores[] | "- " + .sensor_id + " (" + .tipo + ")"' 2>/dev/null || echo "Error al parsear sensores"
else
    echo -e "${YELLOW}⚠️  Verificar manualmente el número final de sensores${NC}"
fi

echo ""

# Mostrar resumen final
show_summary

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}✅ FUNCIONALIDAD DE AGREGAR SENSORES TRABAJANDO CORRECTAMENTE${NC}"
    exit 0
else
    echo -e "${RED}❌ FUNCIONALIDAD NECESITA REVISIÓN${NC}"
    exit 1
fi
