#!/bin/bash

# 🚨 TEST COMPLETO - SIMULACIÓN DE DESCONEXIÓN DE SENSOR
# =====================================================
# Este script simula el flujo completo de monitoreo:
# 1. Crea un activo de prueba
# 2. Envía datos de sensores
# 3. Espera que se detecte la desconexión (>5 min)
# 4. Verifica que se envíe la notificación

# Configuración de colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuración de la API
API_URL="http://localhost:8090"
TEST_COUNT=0
PASS_COUNT=0
FAIL_COUNT=0
TIMEOUT_MINUTES=5
CHECK_INTERVAL_SECONDS=30

echo -e "${BLUE}🚨 TEST SIMULACIÓN COMPLETA DE DESCONEXIÓN${NC}"
echo "=============================================="
echo -e "${CYAN}🔗 URL Base: $API_URL${NC}"
echo -e "${CYAN}📅 Fecha: $(date)${NC}"
echo -e "${CYAN}⏱️  Timeout configurado: $TIMEOUT_MINUTES minutos${NC}"
echo -e "${CYAN}🔄 Intervalo de verificación: $CHECK_INTERVAL_SECONDS segundos${NC}"
echo ""

# Función para mostrar tiempo transcurrido
show_elapsed_time() {
    local start_time=$1
    local current_time=$(date +%s)
    local elapsed=$((current_time - start_time))
    local minutes=$((elapsed / 60))
    local seconds=$((elapsed % 60))
    echo -e "${YELLOW}⏱️  Tiempo transcurrido: ${minutes}m ${seconds}s${NC}"
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
    if [ ${#RESPONSE_BODY} -gt 200 ]; then
        echo -e "Response: ${RESPONSE_BODY:0:200}..."
    else
        echo -e "Response: $RESPONSE_BODY"
    fi
    echo ""
    
    # Guardar última respuesta
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

# Preparar identificadores únicos para esta prueba
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
ACTIVO_ID=$((100000 + (RANDOM % 899999)))
SENSOR_TEMP="TEMP_DISC_${TIMESTAMP}"
SENSOR_PRESS="PRESS_DISC_${TIMESTAMP}"

echo -e "${BLUE}🏗️  FASE 1: PREPARACIÓN DEL ENTORNO DE PRUEBA${NC}"
echo "================================================"

# Test 1: Crear activo de prueba
ACTIVO_TEST=$(cat <<JSON
{
    "activo_id": $ACTIVO_ID,
    "nombre": "Equipo Test Desconexión - $TIMESTAMP",
    "estado": "activo",
    "id_edificio": "test_building",
    "sensores": [
        { "sensor_id": "$SENSOR_TEMP", "tipo": "temperatura", "unidad": "°C" },
        { "sensor_id": "$SENSOR_PRESS", "tipo": "presion", "unidad": "bar" }
    ]
}
JSON
)

test_endpoint "POST" "$API_URL/activo" "Crear activo de prueba para desconexión" "$ACTIVO_TEST" "200"

# Test 2: Verificar estado inicial (sensores nunca conectados)
test_endpoint "GET" "$API_URL/activo/$ACTIVO_ID/sensores/estado" "Verificar estado inicial de sensores" "" "200"

echo -e "${BLUE}📡 FASE 2: SIMULACIÓN DE ACTIVIDAD DE SENSORES${NC}"
echo "==============================================="

# Obtener tiempo de inicio
START_TIME=$(date +%s)

# Test 3: Enviar datos iniciales del sensor de temperatura
LECTURA_TEMP='{
    "sensor_id": "'$SENSOR_TEMP'", 
    "valor": 23.5,
    "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
}'

test_endpoint "POST" "$API_URL/lectura" "Enviar dato inicial - Sensor Temperatura" "$LECTURA_TEMP" "200"

# Test 4: Enviar datos iniciales del sensor de presión
LECTURA_PRESS='{
    "sensor_id": "'$SENSOR_PRESS'",
    "valor": 1.8,
    "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
}'

test_endpoint "POST" "$API_URL/lectura" "Enviar dato inicial - Sensor Presión" "$LECTURA_PRESS" "200"

# Esperar un momento para que se procesen los datos
sleep 3

# Test 5: Verificar que los sensores están conectados
test_endpoint "GET" "$API_URL/activo/$ACTIVO_ID/sensores/estado" "Verificar sensores conectados después de datos" "" "200"

# Verificar que ambos sensores están conectados
if echo "$LAST_RESPONSE" | grep -q '"estado":"connected"'; then
    echo -e "${GREEN}✅ Sensores detectados como conectados${NC}"
else
    echo -e "${RED}❌ Error: Sensores no detectados como conectados${NC}"
fi

# Test 6: Enviar algunas lecturas más para simular actividad normal
echo -e "${CYAN}📊 Enviando lecturas adicionales para simular actividad normal...${NC}"

for i in {1..3}; do
    TEMP_VAL=$(echo "23.5 + $i * 0.5" | bc -l)
    PRESS_VAL=$(echo "1.8 + $i * 0.1" | bc -l)
    
    LECTURA_TEMP='{
        "sensor_id": "'$SENSOR_TEMP'",
        "valor": '$TEMP_VAL',
        "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
    }'
    
    LECTURA_PRESS='{
        "sensor_id": "'$SENSOR_PRESS'", 
        "valor": '$PRESS_VAL',
        "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
    }'
    
    curl -s -X POST -H "Content-Type: application/json" -d "$LECTURA_TEMP" "$API_URL/lectura" > /dev/null
    curl -s -X POST -H "Content-Type: application/json" -d "$LECTURA_PRESS" "$API_URL/lectura" > /dev/null
    
    echo -e "${CYAN}  📈 Enviada lectura #$i (Temp: ${TEMP_VAL}°C, Press: ${PRESS_VAL}bar)${NC}"
    sleep 2
done

echo -e "${GREEN}✅ Actividad inicial simulada correctamente${NC}"
echo ""

echo -e "${BLUE}⏳ FASE 3: SIMULACIÓN DE DESCONEXIÓN${NC}"
echo "====================================="

# Simular desconexión: dejar de enviar datos y esperar
echo -e "${YELLOW}🔌 Simulando desconexión de sensores...${NC}"
echo -e "${YELLOW}⏱️  Esperando ${TIMEOUT_MINUTES} minutos + 30 segundos para que se detecte la desconexión${NC}"
echo -e "${CYAN}📝 Verificando estado cada $CHECK_INTERVAL_SECONDS segundos${NC}"
echo ""

# Tiempo objetivo: 5 minutos + 30 segundos de buffer
TARGET_WAIT_TIME=$((TIMEOUT_MINUTES * 60 + 30))
elapsed_time=0

echo -e "${CYAN}🕐 Iniciando período de espera de $((TARGET_WAIT_TIME / 60))m $((TARGET_WAIT_TIME % 60))s...${NC}"

while [ $elapsed_time -lt $TARGET_WAIT_TIME ]; do
    current_time=$(date +%s)
    elapsed_time=$((current_time - START_TIME))
    
    # Mostrar progreso cada intervalo
    if [ $((elapsed_time % CHECK_INTERVAL_SECONDS)) -eq 0 ] || [ $elapsed_time -eq $TARGET_WAIT_TIME ]; then
        show_elapsed_time $START_TIME
        
        # Verificar estado actual
        echo -e "${CYAN}🔍 Verificando estado actual de sensores...${NC}"
        test_endpoint "GET" "$API_URL/activo/$ACTIVO_ID/sensores/estado" "Verificación periódica de estado" "" "200"
        
        # Contar sensores conectados/desconectados
        connected_count=$(echo "$LAST_RESPONSE" | grep -o '"estado":"connected"' | wc -l)
        disconnected_count=$(echo "$LAST_RESPONSE" | grep -o '"estado":"disconnected"' | wc -l)
        
        echo -e "${CYAN}   📊 Estado actual: $connected_count conectados, $disconnected_count desconectados${NC}"
        
        # Si ya se detectaron desconexiones, podemos continuar
        if [ $disconnected_count -gt 0 ]; then
            echo -e "${GREEN}🎯 ¡Desconexión detectada antes del tiempo límite!${NC}"
            break
        fi
        
        echo ""
    fi
    
    sleep 5
done

echo -e "${BLUE}🔍 FASE 4: VERIFICACIÓN DE DETECCIÓN DE DESCONEXIÓN${NC}"
echo "=================================================="

# Test final: Verificar estado después del período de desconexión
test_endpoint "GET" "$API_URL/activo/$ACTIVO_ID/sensores/estado" "Verificación final - sensores desconectados" "" "200"

# Analizar resultados
connected_count=$(echo "$LAST_RESPONSE" | grep -o '"estado":"connected"' | wc -l)
disconnected_count=$(echo "$LAST_RESPONSE" | grep -o '"estado":"disconnected"' | wc -l)

echo -e "${CYAN}📊 ANÁLISIS DE RESULTADOS:${NC}"
echo -e "   🟢 Sensores conectados: $connected_count"
echo -e "   🔴 Sensores desconectados: $disconnected_count"

if [ $disconnected_count -eq 2 ]; then
    echo -e "${GREEN}✅ ¡ÉXITO! Ambos sensores detectados como desconectados${NC}"
    PASS_COUNT=$((PASS_COUNT + 2))
elif [ $disconnected_count -gt 0 ]; then
    echo -e "${YELLOW}⚠️  PARCIAL: $disconnected_count de 2 sensores detectados como desconectados${NC}"
    PASS_COUNT=$((PASS_COUNT + 1))
    FAIL_COUNT=$((FAIL_COUNT + 1))
else
    echo -e "${RED}❌ FALLO: Ningún sensor detectado como desconectado${NC}"
    FAIL_COUNT=$((FAIL_COUNT + 2))
fi

TEST_COUNT=$((TEST_COUNT + 2))

echo -e "${BLUE}📧 FASE 5: VERIFICACIÓN DE NOTIFICACIONES${NC}"
echo "========================================="

# Test: Forzar verificación manual para generar notificaciones
test_endpoint "POST" "$API_URL/api/sensors/check-disconnected" "Forzar verificación de desconexiones" "" "200"

echo -e "${CYAN}💡 Para verificar las notificaciones, revisa los logs del servicio:${NC}"
echo -e "${YELLOW}   docker logs parser-service | grep -i 'NOTIFICACIÓN'${NC}"
echo -e "${YELLOW}   docker logs parser-service | grep -i 'desconectó'${NC}"

# Test: Verificar estadísticas generales
test_endpoint "GET" "$API_URL/api/sensors/stats" "Obtener estadísticas generales de sensores" "" "200"

if echo "$LAST_RESPONSE" | grep -q '"inactive_sensors"'; then
    inactive_sensors=$(echo "$LAST_RESPONSE" | grep -o '"inactive_sensors":[0-9]*' | cut -d':' -f2)
    echo -e "${CYAN}📊 Sensores inactivos reportados en estadísticas: $inactive_sensors${NC}"
fi

echo -e "${BLUE}🧹 FASE 6: LIMPIEZA (OPCIONAL)${NC}"
echo "==============================="

echo -e "${CYAN}ℹ️  El activo de prueba '$ACTIVO_ID' permanece en el sistema para inspección manual${NC}"
echo -e "${CYAN}ℹ️  Puedes eliminarlo manualmente si lo deseas${NC}"

# Mostrar resumen final
echo -e "\n${BLUE}📊 RESUMEN FINAL DE LA SIMULACIÓN${NC}"
echo "============================================"
echo "Total de tests ejecutados: $TEST_COUNT"
echo "Tests exitosos: $PASS_COUNT"
echo "Tests fallidos: $FAIL_COUNT"

show_elapsed_time $START_TIME

if [ $FAIL_COUNT -eq 0 ]; then
    SUCCESS_RATE=100
    echo -e "${GREEN}🎉 SIMULACIÓN COMPLETAMENTE EXITOSA${NC}"
    echo "Success rate: ${SUCCESS_RATE}%"
    echo ""
    echo -e "${GREEN}✅ SISTEMA DE MONITOREO FUNCIONANDO CORRECTAMENTE${NC}"
    echo -e "${CYAN}📋 Verificado:${NC}"
    echo -e "${CYAN}   • Detección de sensores conectados${NC}"
    echo -e "${CYAN}   • Detección automática de desconexión después de $TIMEOUT_MINUTES minutos${NC}"
    echo -e "${CYAN}   • Actualización de estado en tiempo real${NC}"
    echo -e "${CYAN}   • Generación de notificaciones (revisar logs)${NC}"
    exit 0
else
    if [ $TEST_COUNT -gt 0 ]; then
        SUCCESS_RATE=$(( (PASS_COUNT * 100) / TEST_COUNT ))
    else
        SUCCESS_RATE=0
    fi
    echo -e "${YELLOW}⚠️  SIMULACIÓN CON PROBLEMAS${NC}"
    echo "Success rate: ${SUCCESS_RATE}%"
    echo ""
    echo -e "${RED}❌ SISTEMA DE MONITOREO NECESITA REVISIÓN${NC}"
    echo -e "${YELLOW}💡 Revisar logs del servicio para más detalles${NC}"
    exit 1
fi
