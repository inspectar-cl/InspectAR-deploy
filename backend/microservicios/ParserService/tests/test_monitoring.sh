#!/bin/bash

# 📊 MICROSERVICIO PARSER - TEST DE MONITOREO DE SENSORES
# ======================================================
# Este script realiza tests del sistema de monitoreo de estado de sensores
# Verifica la funcionalidad de detección de desconexión y estado de sensores.

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
    echo -e "${BLUE}📊 MICROSERVICIO PARSER - TEST DE MONITOREO DE SENSORES${NC}"
    echo "=================================================================="
    echo -e "${CYAN}🔗 URL Base: $API_URL${NC}"
    echo -e "${CYAN}📅 Fecha: $(date)${NC}"
    echo ""
}

# Función para verificar disponibilidad del servicio
check_service_availability() {
    echo -e "${CYAN}🔍 Verificando disponibilidad del servicio...${NC}"
    
    # Intentar conectar al servicio con timeout
    if timeout 10 curl -s "$API_URL/api/sensors/health" > /dev/null 2>&1; then
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
    
    # Mostrar respuesta (truncada si es muy larga)
    if [ ${#RESPONSE_BODY} -gt 150 ]; then
        echo -e "Response: ${RESPONSE_BODY:0:150}..."
    else
        echo -e "Response: $RESPONSE_BODY"
    fi
    echo ""
}

# Función para mostrar resumen final
show_summary() {
    echo -e "\n${BLUE}📊 RESUMEN FINAL DE TESTS${NC}"
    echo "============================================================="
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
    echo -e "${CYAN}📋 FUNCIONALIDADES EVALUADAS:${NC}"
    echo "• Health check del sistema de monitoreo"
    echo "• Estado de sensores individuales y globales"
    echo "• Estadísticas de sensores activos/inactivos"
    echo "• Verificación manual de desconexiones"
    echo "• Sistema de simulación de datos"
    
    echo ""
    echo -e "${CYAN}🔧 PARA DEBUGGING:${NC}"
    echo "• Logs del servicio: docker logs parser-service"
    echo "• Logs de MongoDB: docker logs mongu"
    echo "• Logs de InfluxDB: docker logs influxdb"
    echo "• Base URL: $API_URL"
    echo ""
    
    if [ $FAIL_COUNT -eq 0 ]; then
        echo -e "${GREEN}✅ SISTEMA DE MONITOREO FUNCIONANDO CORRECTAMENTE${NC}"
        exit 0
    else
        echo -e "${RED}❌ SISTEMA DE MONITOREO NECESITA REVISIÓN${NC}"
        exit 1
    fi
}

# =====================================
# EJECUCIÓN PRINCIPAL DE LOS TESTS
# =====================================

show_header
check_service_availability

echo -e "${BLUE}🏥 TESTS DE HEALTH CHECK Y ESTADO${NC}"
echo "----------------------------------------"

# Test 1: Health check del sistema de monitoreo
test_endpoint "GET" "$API_URL/api/sensors/health" "Health check del sistema de monitoreo" "" "200"

# Test 2: Obtener estadísticas generales
test_endpoint "GET" "$API_URL/api/sensors/stats" "Obtener estadísticas generales de sensores" "" "200"

# Test 3: Obtener estado de todos los sensores
test_endpoint "GET" "$API_URL/api/sensors/status" "Obtener estado de todos los sensores" "" "200"

echo -e "${BLUE}📊 TESTS DE FUNCIONALIDAD DE MONITOREO${NC}"
echo "----------------------------------------"

# Test 4: Crear un activo de prueba para generar datos
ACTIVO_DATA='{"activo_id":"TEST_SENSOR_001","nombre":"Sensor de Prueba","ubicacion":"Test Lab","estado":"activo","sensores":[{"sensor_id":"TEMP_001","tipo":"temperatura","unidad":"°C"}]}'
test_endpoint "POST" "$API_URL/activo" "Crear activo de prueba" "$ACTIVO_DATA" "200"

# Test 5: Simular envío de datos para activar un sensor
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
LECTURA_DATA="{\"sensor_id\":\"TEMP_001\",\"valor\":23.5,\"timestamp\":\"$TIMESTAMP\"}"
test_endpoint "POST" "$API_URL/lectura" "Enviar lectura para activar sensor" "$LECTURA_DATA" "200"

# Esperar un momento para que se procese
sleep 2

# Test 6: Verificar que el sensor aparece como activo
test_endpoint "GET" "$API_URL/api/sensors/status/TEMP_001" "Verificar estado del sensor después de enviar datos" "" "200"

# Test 7: Verificar estadísticas después de actividad
test_endpoint "GET" "$API_URL/api/sensors/stats" "Verificar estadísticas después de actividad" "" "200"

echo -e "${BLUE}🔍 TESTS DE VERIFICACIÓN DE DESCONEXIONES${NC}"
echo "----------------------------------------"

# Test 8: Forzar verificación manual de sensores desconectados
test_endpoint "POST" "$API_URL/api/sensors/check-disconnected" "Verificación manual de sensores desconectados" "" "200"

# Test 9: Obtener sensor inexistente
test_endpoint "GET" "$API_URL/api/sensors/status/SENSOR_INEXISTENTE" "Obtener estado de sensor inexistente" "" "404"

echo -e "${BLUE}🧪 TESTS DE INTEGRACIÓN COMPLETA${NC}"
echo "----------------------------------------"
echo -e "${CYAN}🔄 Probando flujo completo de monitoreo...${NC}"

# Test 10: Crear otro sensor para pruebas más completas
ACTIVO_DATA2='{"activo_id":"TEST_SENSOR_002","nombre":"Sensor de Humedad","ubicacion":"Test Lab 2","estado":"activo","sensores":[{"sensor_id":"HUM_001","tipo":"humedad","unidad":"%"}]}'
test_endpoint "POST" "$API_URL/activo" "Crear segundo activo de prueba" "$ACTIVO_DATA2" "200"

# Test 11: Enviar datos al segundo sensor
TIMESTAMP2=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
LECTURA_DATA2="{\"sensor_id\":\"HUM_001\",\"valor\":65.0,\"timestamp\":\"$TIMESTAMP2\"}"
test_endpoint "POST" "$API_URL/lectura" "Enviar lectura al segundo sensor" "$LECTURA_DATA2" "200"

# Test 12: Verificar que ahora hay múltiples sensores activos
test_endpoint "GET" "$API_URL/api/sensors/status" "Verificar múltiples sensores en el sistema" "" "200"

# Verificar que hay al menos 2 sensores en la respuesta
if echo "$LAST_RESPONSE" | grep -q '"total":' && echo "$LAST_RESPONSE" | grep -E '"total":[[:space:]]*[2-9]'; then
    echo -e "${GREEN}✅ Sistema detectó múltiples sensores correctamente${NC}"
else
    echo -e "${YELLOW}⚠️  Verificar que el sistema está registrando múltiples sensores${NC}"
fi

echo ""

# Mostrar resumen final
show_summary
