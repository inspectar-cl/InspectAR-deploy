#!/bin/bash

# =============================================================================
# Script de Validación de Lógicas Críticas
# Valida que todas las lógicas implementadas funcionen correctamente
# =============================================================================

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# URLs
GESTION_URL="http://localhost:8092"
PARSER_URL="http://localhost:8090"
ADMIN_EMAIL="admin@example.com"

# Contadores
TESTS_PASSED=0
TESTS_FAILED=0
TOTAL_TESTS=8

echo ""
echo -e "${CYAN}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║     VALIDACIÓN DE LÓGICAS CRÍTICAS - SISTEMA DE ADMIN       ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Función para imprimir resultado de test
print_test_result() {
  local test_name="$1"
  local passed="$2"
  local details="$3"
  
  if [ "$passed" == "true" ]; then
    echo -e "  ${GREEN}✓ PASS${NC} - $test_name"
    [ -n "$details" ] && echo -e "         ${CYAN}→${NC} $details"
    ((TESTS_PASSED++))
  else
    echo -e "  ${RED}✗ FAIL${NC} - $test_name"
    [ -n "$details" ] && echo -e "         ${RED}→${NC} $details"
    ((TESTS_FAILED++))
  fi
  echo ""
}

# =============================================================================
# PREPARACIÓN: Crear datos de prueba
# =============================================================================
echo -e "${YELLOW}[PREPARACIÓN] Creando datos de prueba...${NC}"
echo ""

# Crear edificio
EDIFICIO_RESPONSE=$(curl -s -X POST "$GESTION_URL/admin/edificios" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Torre Validación",
    "direccion": "Calle Prueba 456",
    "ciudad": "Santiago",
    "pais": "Chile",
    "email": "'"$ADMIN_EMAIL"'"
  }')

EDIFICIO_ID=$(echo "$EDIFICIO_RESPONSE" | jq -r '.edificio.id')
echo "Edificio creado: ID = $EDIFICIO_ID"

# Crear activo
ACTIVO_RESPONSE=$(curl -s -X POST "$GESTION_URL/admin/activos" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Bomba Validación",
    "tipo": "bomba de agua",
    "descripcion": "Para validación de lógicas",
    "ubicacion": "Sala Test",
    "edificio_id": '"$EDIFICIO_ID"',
    "email": "'"$ADMIN_EMAIL"'"
  }')

ACTIVO_ID=$(echo "$ACTIVO_RESPONSE" | jq -r '.activo.id')
echo "Activo creado: ID = $ACTIVO_ID"

# Crear sensor
SENSOR_RESPONSE=$(curl -s -X POST "$GESTION_URL/admin/sensores" \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": '"$ACTIVO_ID"',
    "nombre": "Sensor_Validacion_'$ACTIVO_ID'",
    "tipo": "temperatura",
    "unidad": "°C",
    "email": "'"$ADMIN_EMAIL"'"
  }')

SENSOR_ID=$(echo "$SENSOR_RESPONSE" | jq -r '.sensor.sensor_id // .sensor.sensor.sensor_id')
echo "Sensor creado: ID = $SENSOR_ID"
echo ""
sleep 1

# =============================================================================
# TEST 1: Activo con mismo ID en ambos servicios
# =============================================================================
echo -e "${BLUE}[TEST 1]${NC} Verificando ID idéntico en Gestión y ParserService"

GESTION_ACTIVO=$(curl -s "$GESTION_URL/activos/$ACTIVO_ID")
GESTION_ACTIVO_ID=$(echo "$GESTION_ACTIVO" | jq -r '.id')

PARSER_ACTIVO=$(curl -s "$PARSER_URL/activo/$ACTIVO_ID")
PARSER_ACTIVO_ID=$(echo "$PARSER_ACTIVO" | jq -r '.activo.activo_id // .activo_id')

if [ "$GESTION_ACTIVO_ID" == "$PARSER_ACTIVO_ID" ] && [ "$GESTION_ACTIVO_ID" == "$ACTIVO_ID" ]; then
  print_test_result "ID idéntico en ambos servicios" "true" "Gestión: $GESTION_ACTIVO_ID = Parser: $PARSER_ACTIVO_ID"
else
  print_test_result "ID idéntico en ambos servicios" "false" "Gestión: $GESTION_ACTIVO_ID ≠ Parser: $PARSER_ACTIVO_ID"
fi

# =============================================================================
# TEST 2: Sensores NO están en base de datos de gestión
# =============================================================================
echo -e "${BLUE}[TEST 2]${NC} Verificando que sensores NO están en gestion_db"

# Verificar que el activo en gestión no tenga tabla de sensores
GESTION_ACTIVO_DETAIL=$(curl -s "$GESTION_URL/activos/$ACTIVO_ID")
HAS_SENSORS_IN_GESTION=$(echo "$GESTION_ACTIVO_DETAIL" | jq 'has("sensores")')

if [ "$HAS_SENSORS_IN_GESTION" == "false" ]; then
  print_test_result "Sensores NO en gestion_db" "true" "El activo en gestión no contiene array de sensores"
else
  print_test_result "Sensores NO en gestion_db" "false" "ADVERTENCIA: El activo contiene sensores en la respuesta"
fi

# =============================================================================
# TEST 3: Sensores SÍ están en ParserService
# =============================================================================
echo -e "${BLUE}[TEST 3]${NC} Verificando que sensores SÍ están en ParserService"

PARSER_ACTIVO_SENSORES=$(curl -s "$PARSER_URL/activo/$ACTIVO_ID")
SENSOR_COUNT=$(echo "$PARSER_ACTIVO_SENSORES" | jq '.activo.sensores // .sensores | length')

if [ "$SENSOR_COUNT" -gt 0 ]; then
  print_test_result "Sensores en ParserService" "true" "Encontrados $SENSOR_COUNT sensor(es) en MongoDB"
else
  print_test_result "Sensores en ParserService" "false" "No se encontraron sensores en ParserService"
fi

# =============================================================================
# TEST 4: Logs de auditoría para creación de edificio
# =============================================================================
echo -e "${BLUE}[TEST 4]${NC} Verificando log de auditoría para creación de edificio"

EDIFICIO_LOGS=$(curl -s "$GESTION_URL/admin/logs?entidad=edificio&entidad_id=$EDIFICIO_ID")
EDIFICIO_LOG_COUNT=$(echo "$EDIFICIO_LOGS" | jq '.logs | length')

if [ "$EDIFICIO_LOG_COUNT" -gt 0 ]; then
  LOG_ACCION=$(echo "$EDIFICIO_LOGS" | jq -r '.logs[0].accion')
  LOG_IP=$(echo "$EDIFICIO_LOGS" | jq -r '.logs[0].ip_origen')
  print_test_result "Log de creación de edificio" "true" "Acción: $LOG_ACCION, IP: $LOG_IP"
else
  print_test_result "Log de creación de edificio" "false" "No se encontró log para edificio ID $EDIFICIO_ID"
fi

# =============================================================================
# TEST 5: Logs de auditoría para creación de activo
# =============================================================================
echo -e "${BLUE}[TEST 5]${NC} Verificando log de auditoría para creación de activo"

ACTIVO_LOGS=$(curl -s "$GESTION_URL/admin/logs?entidad=activo&entidad_id=$ACTIVO_ID")
ACTIVO_LOG_COUNT=$(echo "$ACTIVO_LOGS" | jq '.logs | length')

if [ "$ACTIVO_LOG_COUNT" -gt 0 ]; then
  LOG_ACCION=$(echo "$ACTIVO_LOGS" | jq -r '.logs[0].accion')
  LOG_DATOS_NUEVOS=$(echo "$ACTIVO_LOGS" | jq -r '.logs[0].datos_nuevos.nombre')
  print_test_result "Log de creación de activo" "true" "Acción: $LOG_ACCION, Nombre: $LOG_DATOS_NUEVOS"
else
  print_test_result "Log de creación de activo" "false" "No se encontró log para activo ID $ACTIVO_ID"
fi

# =============================================================================
# TEST 6: Logs de auditoría para creación de sensor
# =============================================================================
echo -e "${BLUE}[TEST 6]${NC} Verificando log de auditoría para creación de sensor"

SENSOR_LOGS=$(curl -s "$GESTION_URL/admin/logs?entidad=sensor&entidad_id=$ACTIVO_ID")
SENSOR_LOG_COUNT=$(echo "$SENSOR_LOGS" | jq '.logs | length')

if [ "$SENSOR_LOG_COUNT" -gt 0 ]; then
  LOG_ACCION=$(echo "$SENSOR_LOGS" | jq -r '.logs[0].accion')
  LOG_SENSOR_NOMBRE=$(echo "$SENSOR_LOGS" | jq -r '.logs[0].datos_nuevos.nombre')
  print_test_result "Log de creación de sensor" "true" "Acción: $LOG_ACCION, Sensor: $LOG_SENSOR_NOMBRE"
else
  print_test_result "Log de creación de sensor" "false" "No se encontró log para sensor"
fi

# =============================================================================
# TEST 7: IP y User-Agent capturados en logs
# =============================================================================
echo -e "${BLUE}[TEST 7]${NC} Verificando captura de IP y User-Agent"

RECENT_LOG=$(curl -s "$GESTION_URL/admin/logs?limit=1")
LOG_IP=$(echo "$RECENT_LOG" | jq -r '.logs[0].ip_origen')
LOG_USER_AGENT=$(echo "$RECENT_LOG" | jq -r '.logs[0].user_agent')

if [ "$LOG_IP" != "null" ] && [ "$LOG_USER_AGENT" != "null" ]; then
  print_test_result "Captura de IP y User-Agent" "true" "IP: $LOG_IP, User-Agent: ${LOG_USER_AGENT:0:30}..."
else
  print_test_result "Captura de IP y User-Agent" "false" "IP o User-Agent faltante en logs"
fi

# =============================================================================
# TEST 8: Test de Rollback (simulación)
# =============================================================================
echo -e "${BLUE}[TEST 8]${NC} Verificando mecanismo de rollback"

# Intentar crear activo con edificio inexistente (debería fallar)
ROLLBACK_TEST=$(curl -s -X POST "$GESTION_URL/admin/activos" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Activo Rollback Test",
    "tipo": "caldera",
    "descripcion": "Test",
    "ubicacion": "Test",
    "edificio_id": 99999,
    "email": "'"$ADMIN_EMAIL"'"
  }')

ROLLBACK_ERROR=$(echo "$ROLLBACK_TEST" | jq -r '.error // "sin_error"')

if [ "$ROLLBACK_ERROR" != "sin_error" ]; then
  print_test_result "Mecanismo de rollback activo" "true" "Sistema rechaza operación inválida correctamente"
else
  print_test_result "Mecanismo de rollback activo" "false" "Sistema no rechazó operación inválida"
fi

# =============================================================================
# LIMPIEZA
# =============================================================================
echo -e "${YELLOW}[LIMPIEZA] Eliminando datos de prueba...${NC}"
echo ""

curl -s -X DELETE "$GESTION_URL/admin/activos/$ACTIVO_ID?email=$ADMIN_EMAIL" > /dev/null
echo "✓ Activo eliminado"

curl -s -X DELETE "$GESTION_URL/admin/edificios/$EDIFICIO_ID?email=$ADMIN_EMAIL" > /dev/null
echo "✓ Edificio eliminado"
echo ""

# =============================================================================
# RESUMEN FINAL
# =============================================================================
echo ""
echo -e "${CYAN}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║                     RESUMEN DE RESULTADOS                    ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "  Total de tests: ${BLUE}$TOTAL_TESTS${NC}"
echo -e "  ${GREEN}Pasados: $TESTS_PASSED${NC}"
echo -e "  ${RED}Fallidos: $TESTS_FAILED${NC}"
echo ""

PASS_PERCENTAGE=$((TESTS_PASSED * 100 / TOTAL_TESTS))

if [ $TESTS_FAILED -eq 0 ]; then
  echo -e "${GREEN}╔══════════════════════════════════════════════════════════════╗${NC}"
  echo -e "${GREEN}║         ✓ TODAS LAS LÓGICAS VALIDADAS EXITOSAMENTE          ║${NC}"
  echo -e "${GREEN}╚══════════════════════════════════════════════════════════════╝${NC}"
  exit 0
elif [ $PASS_PERCENTAGE -ge 75 ]; then
  echo -e "${YELLOW}╔══════════════════════════════════════════════════════════════╗${NC}"
  echo -e "${YELLOW}║      ⚠ ALGUNAS LÓGICAS NECESITAN REVISIÓN ($PASS_PERCENTAGE% OK)        ║${NC}"
  echo -e "${YELLOW}╚══════════════════════════════════════════════════════════════╝${NC}"
  exit 1
else
  echo -e "${RED}╔══════════════════════════════════════════════════════════════╗${NC}"
  echo -e "${RED}║         ✗ MÚLTIPLES LÓGICAS FALLARON ($PASS_PERCENTAGE% OK)             ║${NC}"
  echo -e "${RED}╚══════════════════════════════════════════════════════════════╝${NC}"
  exit 2
fi
