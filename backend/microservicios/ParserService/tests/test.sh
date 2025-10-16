#!/bin/bash

# ===================================================================
# SCRIPT DE PRUEBAS COMPLETO PARA PARSERSERVICE
# ===================================================================

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8090"
TESTS_PASSED=0
TESTS_FAILED=0

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     PRUEBAS COMPLETAS - PARSERSERVICE API                  ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}\n"

# Función para test exitoso
pass_test() {
    echo -e "${GREEN}✓ PASS${NC} - $1"
    ((TESTS_PASSED++))
}

# Función para test fallido
fail_test() {
    echo -e "${RED}✗ FAIL${NC} - $1"
    ((TESTS_FAILED++))
}

# ============================================================
# TEST 1: HEALTH CHECK
# ============================================================
echo -e "\n${CYAN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}📋 TEST 1: Health Check${NC}"
echo -e "${CYAN}═══════════════════════════════════════════════════════════${NC}"

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}/healthz")
if [ "$HTTP_CODE" -eq 200 ]; then
    pass_test "Servicio respondiendo correctamente"
else
    fail_test "Servicio no responde (HTTP $HTTP_CODE)"
fi

# ============================================================
# TEST 2: GET /activo - Obtener todos los activos
# ============================================================
echo -e "\n${CYAN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}📋 TEST 2: GET /activo - Obtener todos los activos${NC}"
echo -e "${CYAN}═══════════════════════════════════════════════════════════${NC}"

RESPONSE=$(curl -s "${BASE_URL}/activo")
TOTAL=$(echo "$RESPONSE" | jq 'length' 2>/dev/null)

if [ "$TOTAL" -eq 8 ]; then
    pass_test "Devuelve 8 activos correctamente"
    echo "$RESPONSE" | jq -c '.[] | {activo_id, estado, edificio_id}' | head -3
else
    fail_test "Esperaba 8 activos, obtuvo: $TOTAL"
fi

# Verificar estructura de un activo
FIRST_ACTIVO=$(echo "$RESPONSE" | jq '.[0]')
HAS_ACTIVO_ID=$(echo "$FIRST_ACTIVO" | jq 'has("activo_id")')
HAS_ESTADO=$(echo "$FIRST_ACTIVO" | jq 'has("estado")')
HAS_EDIFICIO_ID=$(echo "$FIRST_ACTIVO" | jq 'has("edificio_id")')

if [ "$HAS_ACTIVO_ID" == "true" ] && [ "$HAS_ESTADO" == "true" ] && [ "$HAS_EDIFICIO_ID" == "true" ]; then
    pass_test "Estructura de activo correcta (activo_id, estado, edificio_id)"
else
    fail_test "Estructura de activo incorrecta"
fi

# ============================================================
# TEST 3: GET /activo/:activo_id - Obtener activo específico
# ============================================================
echo -e "\n${CYAN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}📋 TEST 3: GET /activo/:activo_id - Obtener activo específico${NC}"
echo -e "${CYAN}═══════════════════════════════════════════════════════════${NC}"

RESPONSE=$(curl -s "${BASE_URL}/activo/1")
ACTIVO_ID=$(echo "$RESPONSE" | jq -r '.activo_id')
ESTADO=$(echo "$RESPONSE" | jq -r '.estado')

if [ "$ACTIVO_ID" -eq 1 ] && [ -n "$ESTADO" ]; then
    pass_test "Activo ID 1 obtenido correctamente"
    echo "$RESPONSE" | jq '{activo_id, estado, edificio_id}'
else
    fail_test "Error obteniendo activo ID 1"
fi

# Test con activo inexistente
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}/activo/999")
if [ "$HTTP_CODE" -eq 404 ]; then
    pass_test "Retorna 404 para activo inexistente"
else
    fail_test "Esperaba 404, obtuvo: $HTTP_CODE"
fi

# ============================================================
# TEST 4: GET /activo/edificio/:edificio_id - Filtrar por edificio
# ============================================================
echo -e "\n${CYAN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}📋 TEST 4: GET /activo/edificio/:edificio_id${NC}"
echo -e "${CYAN}═══════════════════════════════════════════════════════════${NC}"

# Edificio 1
RESPONSE=$(curl -s "${BASE_URL}/activo/edificio/1")
TOTAL=$(echo "$RESPONSE" | jq -r '.total')
EDIFICIO_ID=$(echo "$RESPONSE" | jq -r '.edificio_id')

if [ "$TOTAL" -eq 3 ] && [ "$EDIFICIO_ID" -eq 1 ]; then
    pass_test "Edificio 1 devuelve 3 activos"
    echo "$RESPONSE" | jq -c '.activos[] | {activo_id, estado}'
else
    fail_test "Edificio 1 esperaba 3 activos, obtuvo: $TOTAL"
fi

# Edificio 2
RESPONSE=$(curl -s "${BASE_URL}/activo/edificio/2")
TOTAL=$(echo "$RESPONSE" | jq -r '.total')

if [ "$TOTAL" -eq 2 ]; then
    pass_test "Edificio 2 devuelve 2 activos"
else
    fail_test "Edificio 2 esperaba 2 activos, obtuvo: $TOTAL"
fi

# Edificio 3
RESPONSE=$(curl -s "${BASE_URL}/activo/edificio/3")
TOTAL=$(echo "$RESPONSE" | jq -r '.total')

if [ "$TOTAL" -eq 2 ]; then
    pass_test "Edificio 3 devuelve 2 activos"
else
    fail_test "Edificio 3 esperaba 2 activos, obtuvo: $TOTAL"
fi

# Edificio inexistente
RESPONSE=$(curl -s "${BASE_URL}/activo/edificio/999")
TOTAL=$(echo "$RESPONSE" | jq -r '.total')

if [ "$TOTAL" -eq 0 ]; then
    pass_test "Edificio inexistente devuelve 0 activos"
else
    fail_test "Edificio inexistente esperaba 0, obtuvo: $TOTAL"
fi

# Parámetro inválido
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}/activo/edificio/abc")
if [ "$HTTP_CODE" -eq 400 ]; then
    pass_test "Retorna 400 para parámetro inválido"
else
    fail_test "Esperaba 400, obtuvo: $HTTP_CODE"
fi

# ============================================================
# TEST 5: POST /activo - Crear activo
# ============================================================
echo -e "\n${CYAN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}📋 TEST 5: POST /activo - Crear activo${NC}"
echo -e "${CYAN}═══════════════════════════════════════════════════════════${NC}"

RESPONSE=$(curl -s -X POST "${BASE_URL}/activo" \
    -H "Content-Type: application/json" \
    -d '{
        "activo_id": 100,
        "estado": "OK",
        "edificio_id": 1
    }')

ACTIVO_ID=$(echo "$RESPONSE" | jq -r '.activo_id')

if [ "$ACTIVO_ID" == "100" ]; then
    pass_test "Activo creado correctamente (ID 100)"
else
    fail_test "Error creando activo"
fi

# Intentar crear duplicado
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${BASE_URL}/activo" \
    -H "Content-Type: application/json" \
    -d '{
        "activo_id": 100,
        "estado": "OK",
        "edificio_id": 1
    }')

if [ "$HTTP_CODE" -eq 409 ]; then
    pass_test "Retorna 409 al intentar crear duplicado"
else
    fail_test "Esperaba 409, obtuvo: $HTTP_CODE"
fi

# ============================================================
# TEST 6: PUT /activo/:activo_id/estado - Actualizar estado
# ============================================================
echo -e "\n${CYAN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}📋 TEST 6: PUT /activo/:activo_id/estado${NC}"
echo -e "${CYAN}═══════════════════════════════════════════════════════════${NC}"

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X PUT "${BASE_URL}/activo/100/estado" \
    -H "Content-Type: application/json" \
    -d '{"estado": "Crítico"}')

if [ "$HTTP_CODE" -eq 200 ]; then
    pass_test "Estado actualizado correctamente"
    
    # Verificar que el estado cambió
    RESPONSE=$(curl -s "${BASE_URL}/activo/100")
    ESTADO=$(echo "$RESPONSE" | jq -r '.estado')
    
    if [ "$ESTADO" == "Crítico" ]; then
        pass_test "Estado verificado: Crítico"
    else
        fail_test "Estado no cambió correctamente"
    fi
else
    fail_test "Error actualizando estado (HTTP $HTTP_CODE)"
fi

# ============================================================
# TEST 7: POST /lectura - Insertar lectura de sensor
# ============================================================
echo -e "\n${CYAN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}📋 TEST 7: POST /lectura - Insertar lectura${NC}"
echo -e "${CYAN}═══════════════════════════════════════════════════════════${NC}"

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${BASE_URL}/lectura" \
    -H "Content-Type: application/json" \
    -d '{
        "sensor_id": "temp1",
        "valor": 75.5,
        "timestamp": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'"
    }')

if [ "$HTTP_CODE" -eq 200 ]; then
    pass_test "Lectura insertada correctamente"
else
    fail_test "Error insertando lectura (HTTP $HTTP_CODE)"
fi

# ============================================================
# TEST 8: GET /sensor/:sensor_id - Obtener datos del sensor
# ============================================================
echo -e "\n${CYAN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}📋 TEST 8: GET /sensor/:sensor_id - Datos del sensor${NC}"
echo -e "${CYAN}═══════════════════════════════════════════════════════════${NC}"

sleep 2  # Esperar a que InfluxDB procese

RESPONSE=$(curl -s "${BASE_URL}/sensor/temp1")
DATOS_COUNT=$(echo "$RESPONSE" | jq 'length' 2>/dev/null)

if [ "$DATOS_COUNT" -ge 1 ]; then
    pass_test "Sensor temp1 devuelve datos (${DATOS_COUNT} lecturas)"
else
    fail_test "Sensor temp1 no devuelve datos"
fi

# ============================================================
# TEST 9: GET /sensor/:sensor_id/last - Última lectura
# ============================================================
echo -e "\n${CYAN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}📋 TEST 9: GET /sensor/:sensor_id/last - Última lectura${NC}"
echo -e "${CYAN}═══════════════════════════════════════════════════════════${NC}"

RESPONSE=$(curl -s "${BASE_URL}/sensor/temp1/last")
HAS_TIEMPO=$(echo "$RESPONSE" | jq 'has("tiempo")')
HAS_VALOR=$(echo "$RESPONSE" | jq 'has("valor")')

if [ "$HAS_TIEMPO" == "true" ] && [ "$HAS_VALOR" == "true" ]; then
    pass_test "Última lectura obtenida correctamente"
    echo "$RESPONSE" | jq '{tiempo, valor}'
else
    fail_test "Estructura de última lectura incorrecta"
fi

# ============================================================
# TEST 10: GET /activo/:activo_id/sensores - Sensores con datos
# ============================================================
echo -e "\n${CYAN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}📋 TEST 10: GET /activo/:activo_id/sensores${NC}"
echo -e "${CYAN}═══════════════════════════════════════════════════════════${NC}"

RESPONSE=$(curl -s "${BASE_URL}/activo/1/sensores")
HAS_ACTIVO_ID=$(echo "$RESPONSE" | jq 'has("activo_id")')
HAS_SENSORES=$(echo "$RESPONSE" | jq 'has("sensores")')

if [ "$HAS_ACTIVO_ID" == "true" ] && [ "$HAS_SENSORES" == "true" ]; then
    pass_test "Estructura correcta con activo_id y sensores"
    SENSORES_COUNT=$(echo "$RESPONSE" | jq '.sensores | length')
    echo "  • Sensores encontrados: $SENSORES_COUNT"
else
    fail_test "Estructura incorrecta"
fi

# ============================================================
# RESUMEN
# ============================================================
echo -e "\n${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                    RESUMEN DE PRUEBAS                      ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"

TOTAL_TESTS=$((TESTS_PASSED + TESTS_FAILED))
PERCENTAGE=$((TESTS_PASSED * 100 / TOTAL_TESTS))

echo -e "\n${GREEN}✓ Tests Exitosos:${NC} $TESTS_PASSED"
echo -e "${RED}✗ Tests Fallidos:${NC} $TESTS_FAILED"
echo -e "${CYAN}📊 Total Tests:${NC} $TOTAL_TESTS"
echo -e "${CYAN}📈 Porcentaje Éxito:${NC} ${PERCENTAGE}%\n"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 ¡TODOS LOS TESTS PASARON EXITOSAMENTE!${NC}\n"
    exit 0
else
    echo -e "${RED}❌ Algunos tests fallaron. Revisar los errores arriba.${NC}\n"
    exit 1
fi
