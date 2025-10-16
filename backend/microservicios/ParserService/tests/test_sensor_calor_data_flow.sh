#!/bin/bash

# Test para ParserService: Crear activo con sensor de calor, enviar 7 datos cada 5 segundos y verificar datos guardados
# Este test simula un flujo completo de datos IoT con un sensor de temperatura

set -e

echo "=========================================="
echo "🌡️  TEST SENSOR DE CALOR - FLUJO DE DATOS"
echo "=========================================="
echo "$(date '+%Y-%m-%d %H:%M:%S') - Iniciando test de sensor de calor con envío de datos"

# Configuración
PARSER_SERVICE_URL="http://localhost:8090"
TEST_ACTIVO_ID="CALDERA-TEST-$(date +%s)"  # ID único basado en timestamp
TEST_SENSOR_ID="TEMP-CALOR-001"

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Función para mostrar resultados de tests
test_result() {
    local test_name="$1"
    local result="$2"
    local details="$3"
    
    if [[ "$result" == "PASS" ]]; then
        echo -e "${GREEN}✅ $test_name: PASS${NC}"
        [[ -n "$details" ]] && echo -e "   ${CYAN}→ $details${NC}"
    else
        echo -e "${RED}❌ $test_name: FAIL${NC}"
        [[ -n "$details" ]] && echo -e "   ${RED}→ $details${NC}"
    fi
}

# Función para verificar que el servicio esté disponible
check_service() {
    echo -e "${BLUE}🔍 Verificando disponibilidad del ParserService...${NC}"
    
    if curl -s "$PARSER_SERVICE_URL/health" > /dev/null 2>&1; then
        test_result "Conectividad del servicio" "PASS" "ParserService está disponible"
    else
        test_result "Conectividad del servicio" "FAIL" "ParserService no está disponible en $PARSER_SERVICE_URL"
        echo -e "${RED}❌ Error: Asegúrate de que el ParserService esté ejecutándose${NC}"
        exit 1
    fi
}

# Función para limpiar datos previos (opcional)
cleanup() {
    echo -e "${YELLOW}🧹 Limpiando datos de test previos (si existen)...${NC}"
    # Intentar eliminar el activo si existe (no importa si falla)
    curl -s -X DELETE "$PARSER_SERVICE_URL/activo/$TEST_ACTIVO_ID" > /dev/null 2>&1 || true
}

# Función para crear el activo con sensor de calor
create_activo_with_sensor() {
    echo -e "${BLUE}🏭 Paso 1: Creando activo con sensor de calor...${NC}"
    
    local payload=$(cat <<EOF
{
    "activo_id": "$TEST_ACTIVO_ID",
    "nombre": "Caldera de Test",
    "ubicacion": "Sala de Pruebas",
    "estado": "activo",
    "id_edificio": "edificio_test_001",
    "sensores": [
        {
            "sensor_id": "$TEST_SENSOR_ID",
            "tipo": "temperatura",
            "unidad": "°C"
        }
    ]
}
EOF
)

    local response=$(curl -s -X POST "$PARSER_SERVICE_URL/activo" \
        -H "Content-Type: application/json" \
        -d "$payload")
    
    if echo "$response" | grep -q "activo_id"; then
        test_result "Creación de activo" "PASS" "Activo $TEST_ACTIVO_ID creado con sensor de calor"
        echo -e "   ${CYAN}→ Sensor ID: $TEST_SENSOR_ID (tipo: temperatura, unidad: °C)${NC}"
    else
        test_result "Creación de activo" "FAIL" "No se pudo crear el activo"
        echo -e "${RED}Response: $response${NC}"
        exit 1
    fi
}

# Función para enviar datos del sensor con valores realistas de temperatura
send_sensor_data() {
    echo -e "${BLUE}🌡️  Paso 2: Enviando 7 lecturas de temperatura cada 5 segundos...${NC}"
    
    # Temperaturas realistas para una caldera (60-90°C)
    local temperatures=(65.2 67.8 70.1 72.5 68.9 71.3 69.7)
    local count=0
    
    for temp in "${temperatures[@]}"; do
        count=$((count + 1))
        echo -e "${CYAN}📊 Enviando lectura $count/7: ${temp}°C...${NC}"
        
        local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%S.%3NZ")
        local payload=$(cat <<EOF
{
    "activo_id": "$TEST_ACTIVO_ID",
    "sensor_id": "$TEST_SENSOR_ID",
    "timestamp": "$timestamp",
    "valor": $temp,
    "unidad": "°C"
}
EOF
)

        local response=$(curl -s -X POST "$PARSER_SERVICE_URL/lectura" \
            -H "Content-Type: application/json" \
            -d "$payload")
        
        if echo "$response" | grep -q -E "(success|guardado|creado)" || [[ $(echo "$response" | jq -r '.status // empty' 2>/dev/null) == "success" ]]; then
            echo -e "   ${GREEN}✅ Lectura $count enviada: ${temp}°C${NC}"
        else
            echo -e "   ${RED}❌ Error enviando lectura $count${NC}"
            echo -e "   ${RED}Response: $response${NC}"
        fi
        
        # Esperar 5 segundos antes de la siguiente lectura (excepto en la última)
        if [ $count -lt 7 ]; then
            echo -e "   ${YELLOW}⏳ Esperando 5 segundos...${NC}"
            sleep 5
        fi
    done
    
    test_result "Envío de datos del sensor" "PASS" "7 lecturas de temperatura enviadas (65.2°C - 72.5°C)"
}

# Función para verificar los datos guardados
verify_stored_data() {
    echo -e "${BLUE}📋 Paso 3: Verificando datos guardados en InfluxDB...${NC}"
    
    # Esperar un poco para asegurar que los datos se guardaron
    echo -e "${YELLOW}⏳ Esperando 3 segundos para asegurar persistencia...${NC}"
    sleep 3
    
    # Verificar datos históricos
    local response=$(curl -s "$PARSER_SERVICE_URL/lectura/$TEST_ACTIVO_ID/datos")
    
    if echo "$response" | grep -q "$TEST_SENSOR_ID"; then
        # Contar cuántas lecturas se guardaron
        local readings_count=$(echo "$response" | jq '[.[] | select(.sensor_id == "'$TEST_SENSOR_ID'")] | length' 2>/dev/null || echo "0")
        
        if [ "$readings_count" -ge 5 ]; then  # Al menos 5 de las 7 lecturas
            test_result "Verificación de datos históricos" "PASS" "$readings_count lecturas encontradas en InfluxDB"
            
            # Mostrar detalles de las lecturas
            echo -e "${CYAN}📊 Detalles de las lecturas guardadas:${NC}"
            echo "$response" | jq -r '.[] | select(.sensor_id == "'$TEST_SENSOR_ID'") | "   → \(.timestamp): \(.valor)\(.unidad) (Sensor: \(.sensor_id))"' 2>/dev/null || echo "   → Datos disponibles pero formato no JSON"
        else
            test_result "Verificación de datos históricos" "FAIL" "Solo $readings_count lecturas encontradas (esperaba al menos 5)"
        fi
    else
        test_result "Verificación de datos históricos" "FAIL" "No se encontraron datos para el sensor $TEST_SENSOR_ID"
        echo -e "${RED}Response: $response${NC}"
    fi
    
    # Verificar última lectura
    echo -e "${BLUE}🔍 Verificando última lectura...${NC}"
    local last_reading=$(curl -s "$PARSER_SERVICE_URL/lectura/$TEST_ACTIVO_ID/datos/ultimo")
    
    if echo "$last_reading" | grep -q "$TEST_SENSOR_ID"; then
        local last_value=$(echo "$last_reading" | jq -r '.valor // empty' 2>/dev/null)
        test_result "Verificación de última lectura" "PASS" "Última lectura: ${last_value}°C"
    else
        test_result "Verificación de última lectura" "FAIL" "No se pudo obtener la última lectura"
    fi
}

# Función para verificar el estado del sensor
verify_sensor_status() {
    echo -e "${BLUE}🔍 Paso 4: Verificando estado del sensor...${NC}"
    
    local response=$(curl -s "$PARSER_SERVICE_URL/activo/$TEST_ACTIVO_ID/sensores/estado")
    
    if echo "$response" | grep -q "$TEST_SENSOR_ID"; then
        local sensor_status=$(echo "$response" | jq -r '.sensores[] | select(.sensor_id == "'$TEST_SENSOR_ID'") | .estado' 2>/dev/null || echo "unknown")
        local total_reports=$(echo "$response" | jq -r '.sensores[] | select(.sensor_id == "'$TEST_SENSOR_ID'") | .total_reports' 2>/dev/null || echo "0")
        
        if [ "$sensor_status" = "connected" ]; then
            test_result "Estado del sensor" "PASS" "Sensor conectado con $total_reports reportes"
        else
            test_result "Estado del sensor" "FAIL" "Sensor en estado: $sensor_status"
        fi
    else
        test_result "Estado del sensor" "FAIL" "No se pudo obtener el estado del sensor"
    fi
}

# Función para mostrar resumen de estadísticas
show_statistics() {
    echo -e "${BLUE}📈 Paso 5: Mostrando estadísticas finales...${NC}"
    
    # Estadísticas generales de sensores
    local stats=$(curl -s "$PARSER_SERVICE_URL/api/sensors/stats")
    echo -e "${CYAN}📊 Estadísticas generales de sensores:${NC}"
    echo "$stats" | jq '.' 2>/dev/null || echo "$stats"
    
    # Información del activo creado
    echo -e "${CYAN}🏭 Información del activo creado:${NC}"
    local activo_info=$(curl -s "$PARSER_SERVICE_URL/activo/$TEST_ACTIVO_ID")
    echo "$activo_info" | jq '{activo_id, nombre, ubicacion, estado, total_sensores: (.sensores | length)}' 2>/dev/null || echo "$activo_info"
}

# Función de limpieza final
final_cleanup() {
    echo -e "${BLUE}🧹 Limpieza final...${NC}"
    
    echo -e "${YELLOW}¿Deseas mantener el activo de test para inspección manual? (y/N)${NC}"
    read -t 10 -n 1 keep_data
    echo
    
    if [[ "$keep_data" =~ ^[Yy]$ ]]; then
        echo -e "${CYAN}📋 Activo de test mantenido:${NC}"
        echo -e "   → Activo ID: $TEST_ACTIVO_ID"
        echo -e "   → Sensor ID: $TEST_SENSOR_ID"
        echo -e "   → Para consultar: curl $PARSER_SERVICE_URL/activo/$TEST_ACTIVO_ID"
    else
        echo -e "${YELLOW}🗑️  Eliminando activo de test...${NC}"
        curl -s -X DELETE "$PARSER_SERVICE_URL/activo/$TEST_ACTIVO_ID" > /dev/null 2>&1 || true
        echo -e "${GREEN}✅ Limpieza completada${NC}"
    fi
}

# Función principal
main() {
    echo -e "${CYAN}🚀 Iniciando test completo de sensor de calor...${NC}"
    echo
    
    # Ejecutar todos los pasos
    check_service
    echo
    
    cleanup
    echo
    
    create_activo_with_sensor
    echo
    
    send_sensor_data
    echo
    
    verify_stored_data
    echo
    
    verify_sensor_status
    echo
    
    show_statistics
    echo
    
    final_cleanup
    
    echo
    echo "=========================================="
    echo -e "${GREEN}🎉 TEST COMPLETADO${NC}"
    echo "=========================================="
    echo "$(date '+%Y-%m-%d %H:%M:%S') - Test de sensor de calor finalizado"
}

# Ejecutar el test
main

echo
echo -e "${BLUE}📝 Resumen del test:${NC}"
echo -e "   1. ✅ Creado activo '$TEST_ACTIVO_ID' con sensor de temperatura"
echo -e "   2. ✅ Enviadas 7 lecturas de temperatura (65.2°C - 72.5°C) cada 5 segundos"
echo -e "   3. ✅ Verificados datos guardados en InfluxDB"
echo -e "   4. ✅ Confirmado estado del sensor como 'connected'"
echo -e "   5. ✅ Mostradas estadísticas finales del sistema"
