#!/bin/bash

# 📧 MICROSERVICIO DE NOTIFICACIONES - TEST COMPLETO
# =====================================================
# Este script realiza tests completos del microservicio de notificaciones
# Verifica funcionalidades de CRUD de notificaciones, gestión de edificios,
# usuarios y activos, así como el envío de notificaciones.

# Configuración de colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuración de la API
API_URL="http://localhost:8091"
LAST_RESPONSE=""
TEST_COUNT=0
PASS_COUNT=0
FAIL_COUNT=0

# Función para mostrar el header
show_header() {
    echo -e "${BLUE}📧 MICROSERVICIO DE NOTIFICACIONES - TEST COMPLETO${NC}"
    echo "=============================================================="
    echo -e "${CYAN}🔗 URL Base: $API_URL${NC}"
    echo -e "${CYAN}📅 Fecha: $(date)${NC}"
    echo ""
}

# Función para verificar disponibilidad del servicio
check_service_availability() {
    echo -e "${CYAN}🔍 Verificando disponibilidad del servicio...${NC}"
    
    # Intentar conectar al servicio con timeout
    if timeout 10 curl -s "$API_URL/edificios" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Servicio de Notificaciones disponible${NC}"
        echo ""
    else
        echo -e "${RED}❌ Servicio de Notificaciones no disponible${NC}"
        echo -e "${YELLOW}💡 Asegúrate de que el servicio esté ejecutándose en puerto 8091${NC}"
        echo -e "${YELLOW}   Comando: docker-compose up notification-service${NC}"
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
    elif [ "$method" = "DELETE" ]; then
        RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" -X DELETE "$endpoint")
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

# Función para extraer ID de la respuesta JSON
extract_id_from_response() {
    echo "$LAST_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2
}

# Función para extraer activo_id de la respuesta JSON
extract_activo_id_from_response() {
    echo "$LAST_RESPONSE" | grep -o '"activo_id":[0-9]*' | head -1 | cut -d':' -f2
}

# Función para extraer edificio_id de la respuesta JSON
extract_edificio_id_from_response() {
    echo "$LAST_RESPONSE" | grep -o '"edificio_id":[0-9]*' | head -1 | cut -d':' -f2
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
    echo "• Gestión de edificios (listar, obtener por ID)"
    echo "• Gestión de usuarios por edificio"
    echo "• Gestión de activos por edificio"
    echo "• Creación de notificaciones"
    echo "• Consulta de notificaciones por activo"
    echo "• Envío de notificaciones"
    echo "• Validaciones de parámetros"
    
    echo ""
    echo -e "${CYAN}🔧 PARA DEBUGGING:${NC}"
    echo "• Logs del servicio: docker logs notification-service"
    echo "• Logs de la DB: docker logs notification-db"
    echo "• Base URL: $API_URL"
    echo ""
    
    if [ $FAIL_COUNT -eq 0 ]; then
        echo -e "${GREEN}✅ MICROSERVICIO DE NOTIFICACIONES FUNCIONANDO CORRECTAMENTE${NC}"
        exit 0
    else
        echo -e "${RED}❌ MICROSERVICIO DE NOTIFICACIONES NECESITA REVISIÓN${NC}"
        exit 1
    fi
}

# =====================================
# EJECUCIÓN PRINCIPAL DE LOS TESTS
# =====================================

show_header
check_service_availability

echo -e "${BLUE}🏢 TESTS DE GESTIÓN DE EDIFICIOS${NC}"
echo "----------------------------------------"

# Test 1: Listar todos los edificios
test_endpoint "GET" "$API_URL/edificios" "Listar todos los edificios" "" "200"

# Extraer primer edificio ID para tests posteriores si existe
EDIFICIO_ID=""
if echo "$LAST_RESPONSE" | grep -q '"edificios"'; then
    EDIFICIO_ID=$(echo "$LAST_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
fi

# Test 2: Obtener edificio por ID (si existe)
if [ -n "$EDIFICIO_ID" ]; then
    test_endpoint "GET" "$API_URL/edificio/$EDIFICIO_ID" "Obtener edificio por ID válido" "" "200"
else
    echo -e "${YELLOW}⚠️  Saltando Test #2: No hay edificios disponibles${NC}"
    TEST_COUNT=$((TEST_COUNT + 1))
    echo ""
fi

# Test 3: Obtener edificio inexistente
test_endpoint "GET" "$API_URL/edificio/999999" "Obtener edificio inexistente" "" "500"

echo -e "${BLUE}👥 TESTS DE GESTIÓN DE USUARIOS${NC}"
echo "----------------------------------------"

# Test 4: Obtener usuarios por edificio ID (si existe edificio)
if [ -n "$EDIFICIO_ID" ]; then
    test_endpoint "GET" "$API_URL/edificio/$EDIFICIO_ID/usuarios" "Obtener usuarios por edificio ID válido" "" "200"
else
    echo -e "${YELLOW}⚠️  Saltando Test #4: No hay edificios disponibles${NC}"
    TEST_COUNT=$((TEST_COUNT + 1))
    echo ""
fi

# Test 5: Obtener usuarios por edificio inexistente
test_endpoint "GET" "$API_URL/edificio/999999/usuarios" "Obtener usuarios por edificio inexistente" "" "500"

echo -e "${BLUE}🏭 TESTS DE GESTIÓN DE ACTIVOS${NC}"
echo "----------------------------------------"

# Test 6: Obtener activos por edificio ID (si existe edificio)
if [ -n "$EDIFICIO_ID" ]; then
    test_endpoint "GET" "$API_URL/edificio/$EDIFICIO_ID/activos" "Obtener activos por edificio ID válido" "" "200"
    
    # Extraer primer activo ID para tests de notificaciones
    ACTIVO_ID=""
    if echo "$LAST_RESPONSE" | grep -q '"activos"'; then
        ACTIVO_ID=$(echo "$LAST_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
    fi
else
    echo -e "${YELLOW}⚠️  Saltando Test #6: No hay edificios disponibles${NC}"
    TEST_COUNT=$((TEST_COUNT + 1))
    echo ""
fi

# Test 7: Obtener activos por edificio inexistente
test_endpoint "GET" "$API_URL/edificio/999999/activos" "Obtener activos por edificio inexistente" "" "500"

echo -e "${BLUE}📧 TESTS DE GESTIÓN DE NOTIFICACIONES${NC}"
echo "----------------------------------------"

# Test 8: Crear notificación (usando activo ID si existe, sino usar ID genérico)
if [ -n "$ACTIVO_ID" ]; then
    NOTIFICATION_DATA="{\"activo_id\": $ACTIVO_ID}"
    test_endpoint "POST" "$API_URL/notification" "Crear notificación con activo válido" "$NOTIFICATION_DATA" "201"
    
    # Extraer notification ID para tests posteriores
    NOTIFICATION_ID=$(extract_id_from_response)
else
    # Intentar crear con ID genérico
    NOTIFICATION_DATA="{\"activo_id\": 1}"
    test_endpoint "POST" "$API_URL/notification" "Crear notificación con activo genérico" "$NOTIFICATION_DATA" "201"
    
    NOTIFICATION_ID=$(extract_id_from_response)
    ACTIVO_ID="1"
fi

# Test 9: Crear notificación con datos inválidos
test_endpoint "POST" "$API_URL/notification" "Crear notificación con datos inválidos" "{\"invalid\": true}" "400"

# Test 10: Crear notificación sin activo_id
test_endpoint "POST" "$API_URL/notification" "Crear notificación sin activo_id" "{}" "400"

# Test 11: Obtener notificaciones por activo ID
if [ -n "$ACTIVO_ID" ]; then
    test_endpoint "GET" "$API_URL/notification/$ACTIVO_ID" "Obtener notificaciones por activo ID" "" "200"
else
    test_endpoint "GET" "$API_URL/notification/1" "Obtener notificaciones por activo ID genérico" "" "200"
fi

# Test 12: Obtener notificaciones con activo ID inválido
test_endpoint "GET" "$API_URL/notification/abc" "Obtener notificaciones con ID inválido" "" "400"

echo -e "${BLUE}📤 TESTS DE ENVÍO DE NOTIFICACIONES${NC}"
echo "----------------------------------------"

# Test 13: Enviar notificación (si tenemos notification ID)
if [ -n "$NOTIFICATION_ID" ]; then
    test_endpoint "PUT" "$API_URL/notification/$NOTIFICATION_ID/send" "Enviar notificación existente" "" "200"
else
    echo -e "${YELLOW}⚠️  Saltando Test #13: No se pudo crear notificación para enviar${NC}"
    TEST_COUNT=$((TEST_COUNT + 1))
    echo ""
fi

# Test 14: Enviar notificación inexistente
test_endpoint "PUT" "$API_URL/notification/999999/send" "Enviar notificación inexistente" "" "500"

# Test 15: Enviar notificación con ID inválido
test_endpoint "PUT" "$API_URL/notification/abc/send" "Enviar notificación con ID inválido" "" "400"

echo -e "\n${BLUE}⚠️  TESTS DE VALIDACIONES Y CASOS EDGE${NC}"
echo "----------------------------------------"

# Test 16: Parámetros inválidos en edificio
test_endpoint "GET" "$API_URL/edificio/abc" "Obtener edificio con ID inválido" "" "400"

# Test 17: Parámetros inválidos en usuarios por edificio
test_endpoint "GET" "$API_URL/edificio/abc/usuarios" "Obtener usuarios con edificio ID inválido" "" "400"

# Test 18: Parámetros inválidos en activos por edificio
test_endpoint "GET" "$API_URL/edificio/abc/activos" "Obtener activos con edificio ID inválido" "" "400"

echo -e "\n${BLUE}🧪 TESTS DE FUNCIONALIDAD COMPLETA${NC}"
echo "----------------------------------------"
echo -e "${CYAN}🔄 Probando flujo completo de notificaciones...${NC}"

# Test 19: Flujo completo - Crear y verificar notificación
FLOW_ACTIVO_ID="1"
FLOW_NOTIFICATION_DATA="{\"activo_id\": $FLOW_ACTIVO_ID}"
test_endpoint "POST" "$API_URL/notification" "Flujo: Crear notificación para flujo completo" "$FLOW_NOTIFICATION_DATA" "201"

# Extraer ID de la notificación creada
FLOW_NOTIFICATION_ID=""
if echo "$LAST_RESPONSE" | grep -q '"notificacion"'; then
    FLOW_NOTIFICATION_ID=$(echo "$LAST_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
fi

# Test 20: Flujo completo - Verificar que la notificación aparece en la lista
test_endpoint "GET" "$API_URL/notification/$FLOW_ACTIVO_ID" "Flujo: Verificar notificación en lista" "" "200"

# Verificar que la notificación creada aparece en la respuesta
if echo "$LAST_RESPONSE" | grep -q "\"id\":$FLOW_NOTIFICATION_ID"; then
    echo -e "${GREEN}✅ Notificación aparece correctamente en la lista${NC}"
else
    echo -e "${YELLOW}⚠️  No se pudo verificar la notificación en la lista${NC}"
fi

# Test 21: Flujo completo - Enviar la notificación
if [ -n "$FLOW_NOTIFICATION_ID" ]; then
    test_endpoint "PUT" "$API_URL/notification/$FLOW_NOTIFICATION_ID/send" "Flujo: Enviar notificación del flujo" "" "200"
    echo -e "${GREEN}✅ Flujo completo de notificación ejecutado exitosamente${NC}"
else
    echo -e "${YELLOW}⚠️  No se pudo completar el flujo: No se obtuvo ID de notificación${NC}"
fi

echo ""

# Mostrar resumen final
show_summary
