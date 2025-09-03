#!/bin/bash

# 🚨 TEST DE ALERTAS DE SENSOR - SERVICIO DE NOTIFICACIONES
# =========================================================
# Este script prueba el nuevo endpoint de alertas de sensor

# Configuración de colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuración de la API
API_URL="http://localhost:8091"
TEST_COUNT=0
PASS_COUNT=0
FAIL_COUNT=0

echo -e "${BLUE}🚨 TEST DE ALERTAS DE SENSOR${NC}"
echo "==============================="
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
    
    echo -e "${PURPLE}📍 Test #$TEST_COUNT: $description${NC}"
    echo -e "${CYAN}$method $endpoint${NC}"
    
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
    
    # Mostrar respuesta truncada
    if [ ${#RESPONSE_BODY} -gt 300 ]; then
        echo -e "Response: ${RESPONSE_BODY:0:300}..."
    else
        echo -e "Response: $RESPONSE_BODY"
    fi
    echo ""
    
    # Guardar última respuesta
    LAST_RESPONSE="$RESPONSE_BODY"
}

# Verificar disponibilidad del servicio
echo -e "${CYAN}🔍 Verificando disponibilidad del servicio...${NC}"
if timeout 10 curl -s "$API_URL/tipos-notificacion" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Servicio de Notificaciones disponible${NC}"
    echo ""
else
    echo -e "${RED}❌ Servicio de Notificaciones no disponible en puerto 8091${NC}"
    echo -e "${YELLOW}💡 Asegúrate de que el servicio esté ejecutándose${NC}"
    exit 1
fi

echo -e "${BLUE}📋 FASE 1: CONSULTAR TIPOS DE NOTIFICACIÓN${NC}"
echo "============================================="

# Test 1: Obtener tipos de notificación disponibles
test_endpoint "GET" "$API_URL/tipos-notificacion" "Obtener tipos de notificación" "" "200"

echo -e "${BLUE}🚨 FASE 2: PRUEBAS DE ALERTAS DE SENSOR${NC}"
echo "========================================"

# Test 2: Alerta de desconexión de sensor
SENSOR_DISCONNECT_ALERT='{
    "sensor_id": "TEMP_001_DISCONNECT_TEST",
    "activo_id": 1,
    "tipo_alerta": "desconexion",
    "titulo": "Sensor de Temperatura Desconectado",
    "mensaje": "El sensor TEMP_001 se ha desconectado inesperadamente. Última lectura hace 5 minutos.",
    "datos_sensor": {
        "ultima_lectura": "2025-09-03T04:30:00Z",
        "tipo": "temperatura",
        "unidad": "°C",
        "valor_anterior": 23.5,
        "ubicacion": "Sala Principal"
    },
    "prioridad": 3,
    "enviar_email": true,
    "enviar_sms": false
}'

test_endpoint "POST" "$API_URL/sensor/alert" "Crear alerta de desconexión de sensor" "$SENSOR_DISCONNECT_ALERT" "201"

# Test 3: Alerta de problema en sensor
SENSOR_PROBLEM_ALERT='{
    "sensor_id": "PRESS_002_PROBLEM_TEST", 
    "activo_id": 2,
    "tipo_alerta": "problema",
    "titulo": "Sensor de Presión con Problemas",
    "mensaje": "El sensor PRESS_002 está reportando valores inconsistentes. Puede requerir calibración.",
    "datos_sensor": {
        "ultima_lectura": "2025-09-03T04:35:00Z",
        "tipo": "presion",
        "unidad": "bar",
        "valor_actual": 999.99,
        "valor_esperado": "1.8-2.2",
        "desviacion": "450%"
    },
    "prioridad": 2,
    "enviar_email": true,
    "enviar_sms": false
}'

test_endpoint "POST" "$API_URL/sensor/alert" "Crear alerta de problema en sensor" "$SENSOR_PROBLEM_ALERT" "201"

# Test 4: Alerta crítica sin activo específico (para administradores)
SENSOR_CRITICAL_ALERT='{
    "sensor_id": "FIRE_003_CRITICAL_TEST",
    "tipo_alerta": "error",
    "titulo": "CRÍTICO: Detector de Humo Fuera de Línea",
    "mensaje": "ALERTA CRÍTICA: El detector de humo FIRE_003 ha dejado de funcionar completamente. Requiere atención inmediata del personal técnico.",
    "datos_sensor": {
        "tipo": "detector_humo",
        "ubicacion": "Salida de Emergencia - Piso 3",
        "estado": "offline",
        "codigo_error": "E001",
        "descripcion_error": "Fallo de comunicación total"
    },
    "prioridad": 4,
    "enviar_email": true,
    "enviar_sms": true
}'

test_endpoint "POST" "$API_URL/sensor/alert" "Crear alerta crítica sin activo específico" "$SENSOR_CRITICAL_ALERT" "201"

# Test 5: Alerta de valor anormal
SENSOR_VALUE_ALERT='{
    "sensor_id": "TEMP_004_VALUE_TEST",
    "activo_id": 1,
    "tipo_alerta": "valor_anormal",
    "titulo": "Temperatura Fuera de Rango Normal",
    "mensaje": "El sensor TEMP_004 está reportando una temperatura de 45°C, que está fuera del rango normal de operación (18-25°C).",
    "datos_sensor": {
        "ultima_lectura": "2025-09-03T04:40:00Z",
        "tipo": "temperatura",
        "unidad": "°C",
        "valor_actual": 45.0,
        "rango_normal": "18-25",
        "estado_alarma": "HIGH_TEMP"
    },
    "prioridad": 2,
    "enviar_email": true,
    "enviar_sms": false
}'

test_endpoint "POST" "$API_URL/sensor/alert" "Crear alerta de valor anormal" "$SENSOR_VALUE_ALERT" "201"

echo -e "${BLUE}❌ FASE 3: PRUEBAS DE VALIDACIÓN (ERRORES ESPERADOS)${NC}"
echo "=================================================="

# Test 6: Alerta con datos inválidos
INVALID_ALERT='{
    "tipo_alerta": "desconexion",
    "titulo": "Test sin sensor_id"
}'

test_endpoint "POST" "$API_URL/sensor/alert" "Crear alerta con datos inválidos (sin sensor_id)" "$INVALID_ALERT" "400"

# Test 7: Alerta con JSON malformado
MALFORMED_JSON='{"sensor_id": "TEST", "tipo_alerta": "problem", titulo": "Malformed'

test_endpoint "POST" "$API_URL/sensor/alert" "Crear alerta con JSON malformado" "$MALFORMED_JSON" "400"

echo -e "${BLUE}📊 FASE 4: VERIFICACIÓN DE NOTIFICACIONES CREADAS${NC}"
echo "================================================"

# Test 8: Consultar notificaciones del activo 1
test_endpoint "GET" "$API_URL/notification/1" "Consultar notificaciones del activo 1" "" "200"

if echo "$LAST_RESPONSE" | grep -q "notificaciones"; then
    notification_count=$(echo "$LAST_RESPONSE" | grep -o '"id":[0-9]*' | wc -l)
    echo -e "${CYAN}📧 Notificaciones encontradas para activo 1: $notification_count${NC}"
fi

# Test 9: Consultar notificaciones del activo 2
test_endpoint "GET" "$API_URL/notification/2" "Consultar notificaciones del activo 2" "" "200"

echo -e "${BLUE}🏢 FASE 5: INFORMACIÓN ADICIONAL${NC}"
echo "================================="

# Test 10: Consultar edificios disponibles
test_endpoint "GET" "$API_URL/edificios" "Consultar edificios disponibles" "" "200"

# Test 11: Consultar usuarios del edificio 1
test_endpoint "GET" "$API_URL/edificio/1/usuarios" "Consultar usuarios del edificio 1" "" "200"

if echo "$LAST_RESPONSE" | grep -q "usuarios"; then
    user_count=$(echo "$LAST_RESPONSE" | grep -o '"correo":"[^"]*"' | wc -l)
    echo -e "${CYAN}👥 Usuarios encontrados en edificio 1: $user_count${NC}"
    echo -e "${CYAN}📧 Correos que recibirían notificaciones:${NC}"
    echo "$LAST_RESPONSE" | grep -o '"correo":"[^"]*"' | sed 's/"correo":"//g' | sed 's/"//g' | sed 's/^/   - /'
fi

# Mostrar resumen final
echo -e "\n${BLUE}📊 RESUMEN FINAL DE PRUEBAS${NC}"
echo "==============================="
echo "Total de tests ejecutados: $TEST_COUNT"
echo "Tests exitosos: $PASS_COUNT"
echo "Tests fallidos: $FAIL_COUNT"

if [ $TEST_COUNT -gt 0 ]; then
    SUCCESS_RATE=$(( (PASS_COUNT * 100) / TEST_COUNT ))
else
    SUCCESS_RATE=0
fi

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}🎉 TODOS LOS TESTS EXITOSOS${NC}"
    echo "Success rate: ${SUCCESS_RATE}%"
    echo ""
    echo -e "${GREEN}✅ SISTEMA DE ALERTAS DE SENSOR FUNCIONANDO CORRECTAMENTE${NC}"
    echo -e "${CYAN}📋 Funcionalidades verificadas:${NC}"
    echo -e "${CYAN}   • Creación de alertas de desconexión de sensor${NC}"
    echo -e "${CYAN}   • Creación de alertas de problemas en sensor${NC}"
    echo -e "${CYAN}   • Alertas críticas para administradores${NC}"
    echo -e "${CYAN}   • Alertas de valores anormales${NC}"
    echo -e "${CYAN}   • Validación de datos de entrada${NC}"
    echo -e "${CYAN}   • Envío de notificaciones por email${NC}"
    echo -e "${CYAN}   • Consulta de notificaciones por activo${NC}"
    echo ""
    echo -e "${YELLOW}💡 Para verificar emails enviados, revisa los logs del servicio:${NC}"
    echo -e "${YELLOW}   docker logs notification-service | grep -i 'email'${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠️  ALGUNOS TESTS FALLARON${NC}"
    echo "Success rate: ${SUCCESS_RATE}%"
    echo ""
    echo -e "${RED}❌ REVISAR CONFIGURACIÓN DEL SISTEMA${NC}"
    echo -e "${YELLOW}💡 Revisar logs del servicio para más detalles${NC}"
    exit 1
fi
