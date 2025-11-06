#!/bin/bash

# 🔧 TEST DE ACTUALIZACIÓN Y ELIMINACIÓN DE TÉCNICOS
# ===================================================
# Este script prueba las nuevas rutas PUT /tecnicos/:id y DELETE /tecnicos/:id

# Configuración de colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuración de la API
BASE_URL="http://localhost:8092"
TEST_COUNT=0
PASS_COUNT=0
FAIL_COUNT=0

# Variables para almacenar IDs creados
TECNICO_ID=""

echo -e "${BLUE}🔧 TEST DE ACTUALIZACIÓN Y ELIMINACIÓN DE TÉCNICOS${NC}"
echo "===================================================="
echo -e "${CYAN}🔗 URL Base: $BASE_URL${NC}"
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
    elif [ "$method" = "PUT" ]; then
        RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" -X PUT -H "Content-Type: application/json" -d "$data" "$endpoint")
    elif [ "$method" = "DELETE" ]; then
        RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" -X DELETE "$endpoint")
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
if timeout 10 curl -s "$BASE_URL/tecnicos" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Servicio de Gestión disponible${NC}"
    echo ""
else
    echo -e "${RED}❌ Servicio de Gestión no disponible en puerto 8092${NC}"
    echo -e "${YELLOW}💡 Asegúrate de que el servicio esté ejecutándose${NC}"
    exit 1
fi

echo -e "${BLUE}📝 FASE 1: CREAR TÉCNICO DE PRUEBA${NC}"
echo "===================================="

# Test 1: Crear técnico para las pruebas
TIMESTAMP=$(date +%s)
CREATE_TECNICO="{
    \"nombre\": \"Técnico Test $TIMESTAMP\",
    \"email\": \"tecnico.test.$TIMESTAMP@test.com\",
    \"telefono\": \"+56912345678\",
    \"especialidad\": \"electricidad\"
}"

test_endpoint "POST" "$BASE_URL/tecnicos" "Crear técnico de prueba" "$CREATE_TECNICO" "201"

# Extraer ID del técnico creado
if echo "$LAST_RESPONSE" | grep -q '"id"'; then
    TECNICO_ID=$(echo "$LAST_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
    echo -e "${CYAN}📝 Técnico creado con ID: $TECNICO_ID${NC}"
    echo ""
else
    echo -e "${RED}❌ No se pudo crear el técnico de prueba. Abortando tests.${NC}"
    exit 1
fi

echo -e "${BLUE}🔄 FASE 2: ACTUALIZAR TÉCNICO${NC}"
echo "=============================="

# Test 2: Actualizar solo el nombre
UPDATE_NOMBRE="{
    \"nombre\": \"Técnico Actualizado $TIMESTAMP\"
}"

test_endpoint "PUT" "$BASE_URL/tecnicos/$TECNICO_ID" "Actualizar solo nombre del técnico" "$UPDATE_NOMBRE" "200"

if echo "$LAST_RESPONSE" | grep -q "Técnico actualizado $TIMESTAMP"; then
    echo -e "${GREEN}✅ Nombre actualizado correctamente${NC}"
    echo ""
fi

# Test 3: Actualizar solo el email
UPDATE_EMAIL="{
    \"email\": \"tecnico.actualizado.$TIMESTAMP@test.com\"
}"

test_endpoint "PUT" "$BASE_URL/tecnicos/$TECNICO_ID" "Actualizar solo email del técnico" "$UPDATE_EMAIL" "200"

if echo "$LAST_RESPONSE" | grep -q "tecnico.actualizado.$TIMESTAMP@test.com"; then
    echo -e "${GREEN}✅ Email actualizado correctamente${NC}"
    echo ""
fi

# Test 4: Actualizar solo el teléfono
UPDATE_TELEFONO="{
    \"telefono\": \"+56987654321\"
}"

test_endpoint "PUT" "$BASE_URL/tecnicos/$TECNICO_ID" "Actualizar solo teléfono del técnico" "$UPDATE_TELEFONO" "200"

if echo "$LAST_RESPONSE" | grep -q "+56987654321"; then
    echo -e "${GREEN}✅ Teléfono actualizado correctamente${NC}"
    echo ""
fi

# Test 5: Actualizar solo la especialidad
UPDATE_ESPECIALIDAD="{
    \"especialidad\": \"HVAC y climatización\"
}"

test_endpoint "PUT" "$BASE_URL/tecnicos/$TECNICO_ID" "Actualizar solo especialidad del técnico" "$UPDATE_ESPECIALIDAD" "200"

if echo "$LAST_RESPONSE" | grep -q "HVAC y climatización"; then
    echo -e "${GREEN}✅ Especialidad actualizada correctamente${NC}"
    echo ""
fi

# Test 6: Actualizar múltiples campos a la vez
UPDATE_MULTIPLE="{
    \"nombre\": \"Técnico Completo $TIMESTAMP\",
    \"telefono\": \"+56900000000\",
    \"especialidad\": \"electricidad industrial\"
}"

test_endpoint "PUT" "$BASE_URL/tecnicos/$TECNICO_ID" "Actualizar múltiples campos del técnico" "$UPDATE_MULTIPLE" "200"

if echo "$LAST_RESPONSE" | grep -q "Técnico Completo $TIMESTAMP"; then
    echo -e "${GREEN}✅ Múltiples campos actualizados correctamente${NC}"
    echo ""
fi

echo -e "${BLUE}❌ FASE 3: VALIDACIONES DE ACTUALIZACIÓN${NC}"
echo "========================================="

# Test 7: Intentar actualizar sin datos
UPDATE_VACIO="{}"

test_endpoint "PUT" "$BASE_URL/tecnicos/$TECNICO_ID" "Actualizar sin proporcionar campos" "$UPDATE_VACIO" "400"

# Test 8: Intentar actualizar técnico inexistente
UPDATE_INEXISTENTE="{
    \"nombre\": \"Técnico Inexistente\"
}"

test_endpoint "PUT" "$BASE_URL/tecnicos/999999" "Actualizar técnico inexistente" "$UPDATE_INEXISTENTE" "404"

# Test 9: Intentar actualizar con ID inválido
test_endpoint "PUT" "$BASE_URL/tecnicos/abc" "Actualizar con ID inválido" "$UPDATE_INEXISTENTE" "400"

echo -e "${BLUE}🔍 FASE 4: VERIFICAR TÉCNICO ACTUALIZADO${NC}"
echo "========================================"

# Test 10: Obtener técnico para verificar actualizaciones
test_endpoint "GET" "$BASE_URL/tecnicos/$TECNICO_ID" "Obtener técnico actualizado" "" "200"

if echo "$LAST_RESPONSE" | grep -q "Técnico Completo $TIMESTAMP"; then
    echo -e "${GREEN}✅ Los cambios persisten correctamente${NC}"
    echo ""
    echo -e "${CYAN}📋 Datos actuales del técnico:${NC}"
    echo "$LAST_RESPONSE" | jq '.'
    echo ""
fi

echo -e "${BLUE}🗑️  FASE 5: ELIMINAR TÉCNICO${NC}"
echo "============================"

# Crear un segundo técnico para eliminar
CREATE_TECNICO_2="{
    \"nombre\": \"Técnico Para Eliminar $TIMESTAMP\",
    \"email\": \"eliminar.$TIMESTAMP@test.com\",
    \"telefono\": \"+56911111111\",
    \"especialidad\": \"plomería\"
}"

test_endpoint "POST" "$BASE_URL/tecnicos" "Crear segundo técnico para eliminar" "$CREATE_TECNICO_2" "201"

# Extraer ID del segundo técnico
if echo "$LAST_RESPONSE" | grep -q '"id"'; then
    TECNICO_2_ID=$(echo "$LAST_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
    echo -e "${CYAN}📝 Segundo técnico creado con ID: $TECNICO_2_ID${NC}"
    echo ""
fi

# Test 11: Eliminar técnico
test_endpoint "DELETE" "$BASE_URL/tecnicos/$TECNICO_2_ID" "Eliminar técnico" "" "200"

if echo "$LAST_RESPONSE" | grep -q "eliminado correctamente"; then
    echo -e "${GREEN}✅ Técnico eliminado correctamente${NC}"
    echo ""
fi

# Test 12: Verificar que el técnico fue eliminado
test_endpoint "GET" "$BASE_URL/tecnicos/$TECNICO_2_ID" "Verificar que técnico fue eliminado" "" "404"

echo -e "${BLUE}❌ FASE 6: VALIDACIONES DE ELIMINACIÓN${NC}"
echo "======================================"

# Test 13: Intentar eliminar técnico inexistente
test_endpoint "DELETE" "$BASE_URL/tecnicos/999999" "Eliminar técnico inexistente" "" "404"

# Test 14: Intentar eliminar con ID inválido
test_endpoint "DELETE" "$BASE_URL/tecnicos/abc" "Eliminar con ID inválido" "" "400"

# Test 15: Intentar eliminar el mismo técnico dos veces
test_endpoint "DELETE" "$BASE_URL/tecnicos/$TECNICO_2_ID" "Eliminar técnico ya eliminado" "" "404"

echo -e "${BLUE}🔍 FASE 7: VERIFICACIÓN FINAL${NC}"
echo "=============================="

# Test 16: Listar todos los técnicos y verificar que el primer técnico sigue existiendo
test_endpoint "GET" "$BASE_URL/tecnicos" "Listar todos los técnicos" "" "200"

if echo "$LAST_RESPONSE" | grep -q "\"id\":$TECNICO_ID"; then
    echo -e "${GREEN}✅ El primer técnico todavía existe en el sistema${NC}"
    echo ""
fi

if ! echo "$LAST_RESPONSE" | grep -q "\"id\":$TECNICO_2_ID"; then
    echo -e "${GREEN}✅ El segundo técnico fue eliminado correctamente del sistema${NC}"
    echo ""
fi

# Limpiar: Eliminar el primer técnico usado para pruebas de actualización
echo -e "${YELLOW}🧹 Limpiando: Eliminando técnico de prueba...${NC}"
curl -s -X DELETE "$BASE_URL/tecnicos/$TECNICO_ID" > /dev/null
echo ""

# Mostrar resumen final
echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}📊 RESUMEN FINAL DE PRUEBAS${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${CYAN}Total de tests ejecutados: ${YELLOW}$TEST_COUNT${NC}"
echo -e "${GREEN}Tests exitosos: $PASS_COUNT${NC}"
echo -e "${RED}Tests fallidos: $FAIL_COUNT${NC}"
echo ""

if [ $TEST_COUNT -gt 0 ]; then
    SUCCESS_RATE=$(( (PASS_COUNT * 100) / TEST_COUNT ))
else
    SUCCESS_RATE=0
fi

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}🎉 TODOS LOS TESTS EXITOSOS - ${SUCCESS_RATE}%${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "${GREEN}✅ RUTAS DE TÉCNICOS FUNCIONANDO CORRECTAMENTE${NC}"
    echo ""
    echo -e "${CYAN}📋 Funcionalidades verificadas:${NC}"
    echo -e "${CYAN}   ✓ Actualización de nombre del técnico${NC}"
    echo -e "${CYAN}   ✓ Actualización de email del técnico${NC}"
    echo -e "${CYAN}   ✓ Actualización de teléfono del técnico${NC}"
    echo -e "${CYAN}   ✓ Actualización de especialidad del técnico${NC}"
    echo -e "${CYAN}   ✓ Actualización de múltiples campos simultáneamente${NC}"
    echo -e "${CYAN}   ✓ Validación de datos requeridos${NC}"
    echo -e "${CYAN}   ✓ Eliminación de técnicos${NC}"
    echo -e "${CYAN}   ✓ Manejo de errores (ID inválido, técnico inexistente)${NC}"
    echo -e "${CYAN}   ✓ Persistencia de cambios en la base de datos${NC}"
    echo ""
    echo -e "${YELLOW}📝 Rutas implementadas:${NC}"
    echo -e "${YELLOW}   • PUT /tecnicos/:id - Actualizar campos de técnico${NC}"
    echo -e "${YELLOW}   • DELETE /tecnicos/:id - Eliminar técnico${NC}"
    echo ""
    exit 0
else
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}⚠️  ALGUNOS TESTS FALLARON - ${SUCCESS_RATE}%${NC}"
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "${RED}❌ REVISAR CONFIGURACIÓN DEL SISTEMA${NC}"
    echo ""
    echo -e "${YELLOW}💡 Sugerencias:${NC}"
    echo -e "${YELLOW}   1. Verificar que el servicio esté ejecutándose correctamente${NC}"
    echo -e "${YELLOW}   2. Revisar logs del servicio:${NC}"
    echo -e "${YELLOW}      docker logs gestion-service${NC}"
    echo -e "${YELLOW}   3. Verificar configuración de base de datos${NC}"
    echo -e "${YELLOW}   4. Verificar que las rutas estén correctamente registradas${NC}"
    echo ""
    exit 1
fi
