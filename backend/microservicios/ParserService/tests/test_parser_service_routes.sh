#!/bin/bash

# Test comprehensivo para todas las rutas del microservicio ParserService
# Excluye únicamente la ruta de estado de sensores (/activo/:activo_id/sensores/estado)

set -e

echo "=========================================="
echo "🔍 TEST COMPREHENSIVO PARSER SERVICE"
echo "=========================================="
echo "$(date '+%Y-%m-%d %H:%M:%S') - Iniciando test de todas las rutas del ParserService"

# Configuración
PARSER_SERVICE_URL="http://localhost:8090"
TEST_ACTIVO_ID="AC-TEST-001"
TEST_SENSOR_ID="TEMP-001"

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para mostrar resultados de tests
test_result() {
    local test_name="$1"
    local result="$2"
    local details="$3"
    
    if [[ "$result" == "PASS" ]]; then
        echo -e "${GREEN}✅ PASS${NC}: $test_name"
        [[ -n "$details" ]] && echo "   └─ $details"
    elif [[ "$result" == "FAIL" ]]; then
        echo -e "${RED}❌ FAIL${NC}: $test_name"
        [[ -n "$details" ]] && echo "   └─ $details"
    elif [[ "$result" == "SKIP" ]]; then
        echo -e "${YELLOW}⚠️  SKIP${NC}: $test_name"
        [[ -n "$details" ]] && echo "   └─ $details"
    fi
}

# Función para verificar si el servicio está disponible
check_service() {
    echo "🔍 Verificando disponibilidad del servicio ParserService..."
    
    response=$(curl -s --max-time 5 "$PARSER_SERVICE_URL/activo" 2>/dev/null || echo "error")
    
    if [[ "$response" == *"error"* || -z "$response" ]]; then
        echo -e "${RED}❌ ParserService no está disponible en $PARSER_SERVICE_URL${NC}"
        echo "🔧 Asegúrate de que Docker Compose esté ejecutándose: docker-compose up -d"
        exit 1
    fi
    
    echo -e "${GREEN}✅ ParserService está disponible${NC}"
}

# Variables para conteo de tests
total_tests=0
passed_tests=0
failed_tests=0
skipped_tests=0

# Función para incrementar contadores
increment_test() {
    local result="$1"
    total_tests=$((total_tests + 1))
    
    case "$result" in
        "PASS") passed_tests=$((passed_tests + 1)) ;;
        "FAIL") failed_tests=$((failed_tests + 1)) ;;
        "SKIP") skipped_tests=$((skipped_tests + 1)) ;;
    esac
}

echo ""
echo "🚀 Iniciando verificación del servicio..."
check_service

echo ""
echo "📋 INICIANDO TESTS DE RUTAS"
echo "============================"

# 1. TEST: GET /activo - Obtener todos los activos
echo ""
echo "1️⃣ TEST: GET /activo (Obtener todos los activos)"
echo "------------------------------------------------"

response=$(curl -s -w "HTTP_CODE:%{http_code}" "$PARSER_SERVICE_URL/activo" 2>/dev/null)
http_code=$(echo "$response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')

if [[ "$http_code" == "200" ]]; then
    count=$(echo "$body" | jq length 2>/dev/null || echo "0")
    test_result "GET /activo" "PASS" "Respuesta HTTP 200, $count activos encontrados"
    increment_test "PASS"
else
    test_result "GET /activo" "FAIL" "HTTP $http_code - $body"
    increment_test "FAIL"
fi

# 2. TEST: POST /activo - Crear un activo de prueba
echo ""
echo "2️⃣ TEST: POST /activo (Crear activo)"
echo "-----------------------------------"

test_activo_data='{
    "activo_id": "'$TEST_ACTIVO_ID'",
    "nombre": "Activo de Prueba",
    "estado": "operativo",
    "ubicacion": "Laboratorio de Testing",
    "edificio_id": 999,
    "tipo": "prueba"
}'

response=$(curl -s -w "HTTP_CODE:%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -d "$test_activo_data" \
    "$PARSER_SERVICE_URL/activo" 2>/dev/null)

http_code=$(echo "$response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')

if [[ "$http_code" == "200" ]]; then
    activo_id=$(echo "$body" | jq -r '.activo_id' 2>/dev/null)
    test_result "POST /activo" "PASS" "Activo creado: $activo_id"
    increment_test "PASS"
    ACTIVO_CREATED=true
elif [[ "$http_code" == "409" ]]; then
    test_result "POST /activo" "PASS" "Activo ya existe (conflicto esperado)"
    increment_test "PASS"
    ACTIVO_CREATED=true
else
    test_result "POST /activo" "FAIL" "HTTP $http_code - $body"
    increment_test "FAIL"
    ACTIVO_CREATED=false
fi

# 3. TEST: GET /activo/:activo_id - Obtener activo específico
echo ""
echo "3️⃣ TEST: GET /activo/:activo_id (Obtener activo específico)"
echo "--------------------------------------------------------"

if [[ "$ACTIVO_CREATED" == true ]]; then
    response=$(curl -s -w "HTTP_CODE:%{http_code}" "$PARSER_SERVICE_URL/activo/$TEST_ACTIVO_ID" 2>/dev/null)
    http_code=$(echo "$response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
    body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')
    
    if [[ "$http_code" == "200" ]]; then
        activo_name=$(echo "$body" | jq -r '.nombre' 2>/dev/null)
        test_result "GET /activo/:activo_id" "PASS" "Activo encontrado: $activo_name"
        increment_test "PASS"
    else
        test_result "GET /activo/:activo_id" "FAIL" "HTTP $http_code - $body"
        increment_test "FAIL"
    fi
else
    test_result "GET /activo/:activo_id" "SKIP" "No se pudo crear activo de prueba"
    increment_test "SKIP"
fi

# 4. TEST: PUT /activo/:activo_id/estado - Actualizar estado del activo
echo ""
echo "4️⃣ TEST: PUT /activo/:activo_id/estado (Actualizar estado)"
echo "--------------------------------------------------------"

if [[ "$ACTIVO_CREATED" == true ]]; then
    update_data='{"estado": "mantenimiento"}'
    
    response=$(curl -s -w "HTTP_CODE:%{http_code}" -X PUT \
        -H "Content-Type: application/json" \
        -d "$update_data" \
        "$PARSER_SERVICE_URL/activo/$TEST_ACTIVO_ID/estado" 2>/dev/null)
    
    http_code=$(echo "$response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
    body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')
    
    if [[ "$http_code" == "200" ]]; then
        test_result "PUT /activo/:activo_id/estado" "PASS" "Estado actualizado correctamente"
        increment_test "PASS"
    else
        test_result "PUT /activo/:activo_id/estado" "FAIL" "HTTP $http_code - $body"
        increment_test "FAIL"
    fi
else
    test_result "PUT /activo/:activo_id/estado" "SKIP" "No se pudo crear activo de prueba"
    increment_test "SKIP"
fi

# 5. TEST: POST /lectura - Crear lectura de sensor
echo ""
echo "5️⃣ TEST: POST /lectura (Crear lectura de sensor)"
echo "------------------------------------------------"

if [[ "$ACTIVO_CREATED" == true ]]; then
    lectura_data='{
        "activo_id": "'$TEST_ACTIVO_ID'",
        "sensor_id": "'$TEST_SENSOR_ID'",
        "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
        "valor": 25.5,
        "unidad": "°C",
        "tipo_sensor": "temperatura"
    }'
    
    response=$(curl -s -w "HTTP_CODE:%{http_code}" -X POST \
        -H "Content-Type: application/json" \
        -d "$lectura_data" \
        "$PARSER_SERVICE_URL/lectura" 2>/dev/null)
    
    http_code=$(echo "$response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
    body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')
    
    if [[ "$http_code" == "200" || "$http_code" == "201" ]]; then
        test_result "POST /lectura" "PASS" "Lectura creada correctamente"
        increment_test "PASS"
        LECTURA_CREATED=true
    else
        test_result "POST /lectura" "FAIL" "HTTP $http_code - $body"
        increment_test "FAIL"
        LECTURA_CREATED=false
    fi
else
    test_result "POST /lectura" "SKIP" "No se pudo crear activo de prueba"
    increment_test "SKIP"
    LECTURA_CREATED=false
fi

# 6. TEST: GET /lectura/:activo_id/datos - Obtener datos de sensores por activo
echo ""
echo "6️⃣ TEST: GET /lectura/:activo_id/datos (Datos de sensores)"
echo "--------------------------------------------------------"

if [[ "$LECTURA_CREATED" == true ]]; then
    response=$(curl -s -w "HTTP_CODE:%{http_code}" "$PARSER_SERVICE_URL/lectura/$TEST_ACTIVO_ID/datos" 2>/dev/null)
    http_code=$(echo "$response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
    body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')
    
    if [[ "$http_code" == "200" ]]; then
        test_result "GET /lectura/:activo_id/datos" "PASS" "Datos de sensores obtenidos"
        increment_test "PASS"
    else
        test_result "GET /lectura/:activo_id/datos" "FAIL" "HTTP $http_code - $body"
        increment_test "FAIL"
    fi
else
    test_result "GET /lectura/:activo_id/datos" "SKIP" "No se pudo crear lectura de prueba"
    increment_test "SKIP"
fi

# 7. TEST: GET /lectura/:activo_id/datos/ultimo - Obtener última lectura
echo ""
echo "7️⃣ TEST: GET /lectura/:activo_id/datos/ultimo (Última lectura)"
echo "------------------------------------------------------------"

if [[ "$LECTURA_CREATED" == true ]]; then
    response=$(curl -s -w "HTTP_CODE:%{http_code}" "$PARSER_SERVICE_URL/lectura/$TEST_ACTIVO_ID/datos/ultimo" 2>/dev/null)
    http_code=$(echo "$response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
    body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')
    
    if [[ "$http_code" == "200" ]]; then
        test_result "GET /lectura/:activo_id/datos/ultimo" "PASS" "Última lectura obtenida"
        increment_test "PASS"
    else
        test_result "GET /lectura/:activo_id/datos/ultimo" "FAIL" "HTTP $http_code - $body"
        increment_test "FAIL"
    fi
else
    test_result "GET /lectura/:activo_id/datos/ultimo" "SKIP" "No se pudo crear lectura de prueba"
    increment_test "SKIP"
fi

echo ""
echo "🔧 TESTS DE RUTAS DE MONITOREO DE SENSORES"
echo "==========================================="

# 8. TEST: GET /api/sensors/status - Estado de todos los sensores
echo ""
echo "8️⃣ TEST: GET /api/sensors/status (Estado de sensores)"
echo "----------------------------------------------------"

response=$(curl -s -w "HTTP_CODE:%{http_code}" "$PARSER_SERVICE_URL/api/sensors/status" 2>/dev/null)
http_code=$(echo "$response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')

if [[ "$http_code" == "200" ]]; then
    test_result "GET /api/sensors/status" "PASS" "Estado de sensores obtenido"
    increment_test "PASS"
else
    test_result "GET /api/sensors/status" "FAIL" "HTTP $http_code - $body"
    increment_test "FAIL"
fi

# 9. TEST: GET /api/sensors/status/:sensor_id - Estado de sensor específico
echo ""
echo "9️⃣ TEST: GET /api/sensors/status/:sensor_id (Sensor específico)"
echo "--------------------------------------------------------------"

if [[ "$LECTURA_CREATED" == true ]]; then
    response=$(curl -s -w "HTTP_CODE:%{http_code}" "$PARSER_SERVICE_URL/api/sensors/status/$TEST_SENSOR_ID" 2>/dev/null)
    http_code=$(echo "$response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
    body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')
    
    if [[ "$http_code" == "200" ]]; then
        test_result "GET /api/sensors/status/:sensor_id" "PASS" "Estado de sensor específico obtenido"
        increment_test "PASS"
    elif [[ "$http_code" == "404" ]]; then
        test_result "GET /api/sensors/status/:sensor_id" "PASS" "Sensor no encontrado (esperado para sensor de prueba)"
        increment_test "PASS"
    else
        test_result "GET /api/sensors/status/:sensor_id" "FAIL" "HTTP $http_code - $body"
        increment_test "FAIL"
    fi
else
    test_result "GET /api/sensors/status/:sensor_id" "SKIP" "No se pudo crear sensor de prueba"
    increment_test "SKIP"
fi

# 10. TEST: POST /api/sensors/check-disconnected - Verificar sensores desconectados
echo ""
echo "🔟 TEST: POST /api/sensors/check-disconnected (Sensores desconectados)"
echo "--------------------------------------------------------------------"

response=$(curl -s -w "HTTP_CODE:%{http_code}" -X POST "$PARSER_SERVICE_URL/api/sensors/check-disconnected" 2>/dev/null)
http_code=$(echo "$response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')

if [[ "$http_code" == "200" ]]; then
    test_result "POST /api/sensors/check-disconnected" "PASS" "Verificación de sensores desconectados completada"
    increment_test "PASS"
else
    test_result "POST /api/sensors/check-disconnected" "FAIL" "HTTP $http_code - $body"
    increment_test "FAIL"
fi

# 11. TEST: GET /api/sensors/stats - Estadísticas de sensores
echo ""
echo "1️⃣1️⃣ TEST: GET /api/sensors/stats (Estadísticas de sensores)"
echo "-----------------------------------------------------------"

response=$(curl -s -w "HTTP_CODE:%{http_code}" "$PARSER_SERVICE_URL/api/sensors/stats" 2>/dev/null)
http_code=$(echo "$response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')

if [[ "$http_code" == "200" ]]; then
    test_result "GET /api/sensors/stats" "PASS" "Estadísticas de sensores obtenidas"
    increment_test "PASS"
else
    test_result "GET /api/sensors/stats" "FAIL" "HTTP $http_code - $body"
    increment_test "FAIL"
fi

# 12. TEST: GET /api/sensors/health - Health check de sensores
echo ""
echo "1️⃣2️⃣ TEST: GET /api/sensors/health (Health check)"
echo "------------------------------------------------"

response=$(curl -s -w "HTTP_CODE:%{http_code}" "$PARSER_SERVICE_URL/api/sensors/health" 2>/dev/null)
http_code=$(echo "$response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')

if [[ "$http_code" == "200" ]]; then
    test_result "GET /api/sensors/health" "PASS" "Health check exitoso"
    increment_test "PASS"
else
    test_result "GET /api/sensors/health" "FAIL" "HTTP $http_code - $body"
    increment_test "FAIL"
fi

echo ""
echo "🧹 LIMPIEZA DE DATOS DE PRUEBA"
echo "==============================="

# Cleanup: Intentar eliminar el activo de prueba (si existe endpoint DELETE)
# Nota: El ParserService no parece tener endpoint DELETE, pero verificamos
echo "🗑️  Verificando endpoints de limpieza..."

cleanup_response=$(curl -s -w "HTTP_CODE:%{http_code}" -X DELETE "$PARSER_SERVICE_URL/activo/$TEST_ACTIVO_ID" 2>/dev/null)
cleanup_http_code=$(echo "$cleanup_response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)

if [[ "$cleanup_http_code" == "200" || "$cleanup_http_code" == "204" ]]; then
    echo "✅ Activo de prueba eliminado correctamente"
elif [[ "$cleanup_http_code" == "404" ]]; then
    echo "ℹ️  Endpoint DELETE no disponible o activo ya eliminado"
else
    echo "⚠️  No se pudo eliminar activo de prueba (no crítico)"
fi

echo ""
echo "📊 RESUMEN FINAL DE TESTS"
echo "========================="
echo ""
echo "📈 Estadísticas:"
echo "   Total de tests ejecutados: $total_tests"
echo -e "   ${GREEN}✅ Tests pasados: $passed_tests${NC}"
echo -e "   ${RED}❌ Tests fallidos: $failed_tests${NC}"
echo -e "   ${YELLOW}⚠️  Tests omitidos: $skipped_tests${NC}"

echo ""
echo "📋 Rutas probadas (excluye /activo/:activo_id/sensores/estado):"
echo "   1. GET /activo - Listar todos los activos"
echo "   2. POST /activo - Crear nuevo activo"
echo "   3. GET /activo/:activo_id - Obtener activo específico"
echo "   4. PUT /activo/:activo_id/estado - Actualizar estado de activo"
echo "   5. POST /lectura - Crear lectura de sensor"
echo "   6. GET /lectura/:activo_id/datos - Obtener datos de sensores"
echo "   7. GET /lectura/:activo_id/datos/ultimo - Obtener última lectura"
echo "   8. GET /api/sensors/status - Estado de todos los sensores"
echo "   9. GET /api/sensors/status/:sensor_id - Estado de sensor específico"
echo "   10. POST /api/sensors/check-disconnected - Verificar sensores desconectados"
echo "   11. GET /api/sensors/stats - Estadísticas de sensores"
echo "   12. GET /api/sensors/health - Health check de sensores"

echo ""
success_rate=$((passed_tests * 100 / total_tests))

if [[ $failed_tests -eq 0 ]]; then
    echo -e "${GREEN}🎉 ¡ÉXITO TOTAL! Todos los tests críticos pasaron ($success_rate% éxito)${NC}"
    echo "✅ El microservicio ParserService está funcionando correctamente"
    exit 0
elif [[ $success_rate -ge 80 ]]; then
    echo -e "${YELLOW}⚠️  ÉXITO PARCIAL: $success_rate% de tests pasaron${NC}"
    echo "🔧 Algunos tests fallaron, pero el servicio está mayormente funcional"
    exit 1
else
    echo -e "${RED}❌ FALLO: Solo $success_rate% de tests pasaron${NC}"
    echo "🚨 Problemas graves detectados en el microservicio ParserService"
    exit 2
fi
