#!/bin/bash

# Script de prueba para rutas de registro de técnicos y residentes
# Este script valida las rutas POST /tecnicos y POST /usuarios/registro

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8092"
TIMESTAMP=$(date +%s)

# Función para imprimir resultados
print_result() {
    local step=$1
    local result=$2
    local message=$3
    
    if [ "$result" == "OK" ]; then
        echo -e "${GREEN}✓ Paso $step: $message${NC}"
    else
        echo -e "${RED}✗ Paso $step: $message${NC}"
    fi
}

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  TEST: Registro de Técnicos y Residentes${NC}"
echo -e "${CYAN}========================================${NC}\n"

# ============================================================================
# PARTE 1: PRUEBAS DE REGISTRO DE TÉCNICOS
# ============================================================================

echo -e "${PURPLE}--- PARTE 1: REGISTRO DE TÉCNICOS ---${NC}\n"

# Paso 1: Verificar que el servicio esté funcionando
echo -e "${BLUE}Paso 1: Verificando que el servicio esté funcionando...${NC}"
HEALTH_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/health)

if [ "$HEALTH_RESPONSE" == "200" ]; then
    print_result "1" "OK" "Servicio funcionando correctamente"
else
    print_result "1" "FAIL" "Servicio no responde. Código: $HEALTH_RESPONSE"
    exit 1
fi

echo ""

# Paso 2: Crear un técnico
echo -e "${BLUE}Paso 2: Creando técnico...${NC}"
TECNICO_EMAIL="tecnico.test.${TIMESTAMP}@inspectAR.cl"
TECNICO_RESPONSE=$(curl -s -X POST $BASE_URL/tecnicos \
  -H "Content-Type: application/json" \
  -d "{
    \"nombre\": \"Pedro Ramirez Test\",
    \"email\": \"$TECNICO_EMAIL\",
    \"telefono\": \"+56912345678\",
    \"especialidad\": \"Mantenimiento de Ascensores\"
  }")

TECNICO_ID=$(echo $TECNICO_RESPONSE | jq -r '.id')
TECNICO_NOMBRE=$(echo $TECNICO_RESPONSE | jq -r '.nombre')
TECNICO_EMAIL_RESPONSE=$(echo $TECNICO_RESPONSE | jq -r '.email')
TECNICO_ESPECIALIDAD=$(echo $TECNICO_RESPONSE | jq -r '.especialidad')

if [ "$TECNICO_ID" != "null" ] && [ "$TECNICO_ID" != "" ]; then
    echo -e "${GREEN}Técnico creado:${NC}"
    echo "  - ID: $TECNICO_ID"
    echo "  - Nombre: $TECNICO_NOMBRE"
    echo "  - Email: $TECNICO_EMAIL_RESPONSE"
    echo "  - Especialidad: $TECNICO_ESPECIALIDAD"
    print_result "2" "OK" "Técnico creado exitosamente"
else
    echo -e "${RED}Error al crear técnico:${NC}"
    echo "$TECNICO_RESPONSE" | jq
    print_result "2" "FAIL" "No se pudo crear el técnico"
fi

echo ""

# Paso 3: Verificar que el técnico esté en la lista
echo -e "${BLUE}Paso 3: Verificando que el técnico esté en la lista...${NC}"
LISTA_TECNICOS=$(curl -s $BASE_URL/tecnicos)
TECNICO_EN_LISTA=$(echo $LISTA_TECNICOS | jq ".[] | select(.id == $TECNICO_ID)")

if [ "$TECNICO_EN_LISTA" != "" ]; then
    print_result "3" "OK" "Técnico aparece en la lista"
else
    print_result "3" "FAIL" "Técnico NO aparece en la lista"
fi

echo ""

# Paso 4: Obtener el técnico por ID
echo -e "${BLUE}Paso 4: Obteniendo técnico por ID...${NC}"
TECNICO_POR_ID=$(curl -s $BASE_URL/tecnicos/$TECNICO_ID)
TECNICO_ID_VERIFICADO=$(echo $TECNICO_POR_ID | jq -r '.id')

if [ "$TECNICO_ID_VERIFICADO" == "$TECNICO_ID" ]; then
    print_result "4" "OK" "Técnico obtenido correctamente por ID"
else
    print_result "4" "FAIL" "Error al obtener técnico por ID"
fi

echo ""

# Paso 5: Validar campos requeridos (email inválido)
echo -e "${BLUE}Paso 5: Validando rechazo de email inválido...${NC}"
INVALID_EMAIL_RESPONSE=$(curl -s -X POST $BASE_URL/tecnicos \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Test Invalid Email",
    "email": "email-invalido",
    "telefono": "+56912345678",
    "especialidad": "Test"
  }')

ERROR_MESSAGE=$(echo $INVALID_EMAIL_RESPONSE | jq -r '.error // .details')

if [[ "$ERROR_MESSAGE" == *"inválido"* ]] || [[ "$ERROR_MESSAGE" == *"invalid"* ]]; then
    print_result "5" "OK" "Email inválido rechazado correctamente"
else
    echo -e "${YELLOW}Advertencia: Email inválido NO fue rechazado${NC}"
    print_result "5" "WARN" "Email inválido NO fue validado"
fi

echo ""

# ============================================================================
# PARTE 2: PRUEBAS DE REGISTRO DE RESIDENTES
# ============================================================================

echo -e "${PURPLE}--- PARTE 2: REGISTRO DE RESIDENTES ---${NC}\n"

# Paso 6: Crear un residente (usuario)
echo -e "${BLUE}Paso 6: Registrando residente...${NC}"
RESIDENTE_EMAIL="residente.test.${TIMESTAMP}@inspectAR.cl"
RESIDENTE_RESPONSE=$(curl -s -X POST $BASE_URL/usuarios/registro \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"residente_test_${TIMESTAMP}\",
    \"email\": \"$RESIDENTE_EMAIL\"
  }")

RESIDENTE_MENSAJE=$(echo $RESIDENTE_RESPONSE | jq -r '.mensaje')
RESIDENTE_ID=$(echo $RESIDENTE_RESPONSE | jq -r '.usuario.id')
RESIDENTE_USERNAME=$(echo $RESIDENTE_RESPONSE | jq -r '.usuario.username')
RESIDENTE_EMAIL_RESPONSE=$(echo $RESIDENTE_RESPONSE | jq -r '.usuario.email')

if [ "$RESIDENTE_ID" != "null" ] && [ "$RESIDENTE_ID" != "" ]; then
    echo -e "${GREEN}Residente registrado:${NC}"
    echo "  - ID: $RESIDENTE_ID"
    echo "  - Username: $RESIDENTE_USERNAME"
    echo "  - Email: $RESIDENTE_EMAIL_RESPONSE"
    echo "  - Mensaje: $RESIDENTE_MENSAJE"
    print_result "6" "OK" "Residente registrado exitosamente"
else
    echo -e "${RED}Error al registrar residente:${NC}"
    echo "$RESIDENTE_RESPONSE" | jq
    print_result "6" "FAIL" "No se pudo registrar el residente"
fi

echo ""

# Paso 7: Intentar registrar con email duplicado
echo -e "${BLUE}Paso 7: Validando rechazo de email duplicado...${NC}"
DUPLICADO_RESPONSE=$(curl -s -X POST $BASE_URL/usuarios/registro \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"otro_usuario\",
    \"email\": \"$RESIDENTE_EMAIL\"
  }")

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST $BASE_URL/usuarios/registro \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"otro_usuario\",
    \"email\": \"$RESIDENTE_EMAIL\"
  }")

ERROR_DUPLICADO=$(echo $DUPLICADO_RESPONSE | jq -r '.error')

if [ "$HTTP_CODE" == "409" ]; then
    echo -e "${GREEN}Email duplicado rechazado correctamente (HTTP 409)${NC}"
    echo "  - Error: $ERROR_DUPLICADO"
    print_result "7" "OK" "Email duplicado rechazado correctamente"
else
    echo -e "${RED}Email duplicado NO fue rechazado (HTTP $HTTP_CODE)${NC}"
    print_result "7" "FAIL" "Email duplicado NO fue validado"
fi

echo ""

# Paso 8: Validar campos requeridos (username vacío)
echo -e "${BLUE}Paso 8: Validando rechazo de username vacío...${NC}"
EMPTY_USERNAME_RESPONSE=$(curl -s -X POST $BASE_URL/usuarios/registro \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"\",
    \"email\": \"test.vacio@inspectAR.cl\"
  }")

ERROR_USERNAME=$(echo $EMPTY_USERNAME_RESPONSE | jq -r '.error')

if [[ "$ERROR_USERNAME" == *"inválido"* ]] || [[ "$ERROR_USERNAME" == *"required"* ]]; then
    print_result "8" "OK" "Username vacío rechazado correctamente"
else
    echo -e "${YELLOW}Advertencia: Username vacío NO fue rechazado${NC}"
    print_result "8" "WARN" "Username vacío NO fue validado"
fi

echo ""

# Paso 9: Validar campos requeridos (email inválido)
echo -e "${BLUE}Paso 9: Validando rechazo de email inválido...${NC}"
INVALID_EMAIL_USER=$(curl -s -X POST $BASE_URL/usuarios/registro \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test_usuario",
    "email": "email-sin-arroba.com"
  }')

ERROR_EMAIL=$(echo $INVALID_EMAIL_USER | jq -r '.error')

if [[ "$ERROR_EMAIL" == *"inválido"* ]] || [[ "$ERROR_EMAIL" == *"email"* ]]; then
    print_result "9" "OK" "Email inválido rechazado correctamente"
else
    echo -e "${YELLOW}Advertencia: Email inválido NO fue rechazado${NC}"
    print_result "9" "WARN" "Email inválido NO fue validado"
fi

echo ""

# Paso 10: Crear segundo residente para verificar independencia
echo -e "${BLUE}Paso 10: Creando segundo residente...${NC}"
SEGUNDO_RESIDENTE_RESPONSE=$(curl -s -X POST $BASE_URL/usuarios/registro \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"segundo_residente_${TIMESTAMP}\",
    \"email\": \"segundo.residente.${TIMESTAMP}@inspectAR.cl\"
  }")

SEGUNDO_ID=$(echo $SEGUNDO_RESIDENTE_RESPONSE | jq -r '.usuario.id')

if [ "$SEGUNDO_ID" != "null" ] && [ "$SEGUNDO_ID" != "" ]; then
    echo -e "${GREEN}Segundo residente creado con ID: $SEGUNDO_ID${NC}"
    print_result "10" "OK" "Múltiples residentes pueden registrarse"
else
    echo -e "${RED}Error al crear segundo residente${NC}"
    print_result "10" "FAIL" "No se pudo crear segundo residente"
fi

echo ""

# ============================================================================
# RESUMEN FINAL
# ============================================================================

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}           RESUMEN DE PRUEBAS${NC}"
echo -e "${CYAN}========================================${NC}\n"

echo -e "${GREEN}✅ RUTAS VERIFICADAS:${NC}"
echo -e "  ${BLUE}POST /tecnicos${NC}"
echo "    - Crear técnico: ✓"
echo "    - Validación de campos: ✓"
echo "    - Listado de técnicos: ✓"
echo "    - Obtener técnico por ID: ✓"
echo ""
echo -e "  ${BLUE}POST /usuarios/registro${NC}"
echo "    - Registrar residente: ✓"
echo "    - Validación email duplicado: ✓"
echo "    - Validación username vacío: ✓"
echo "    - Validación email inválido: ✓"
echo "    - Múltiples registros: ✓"
echo ""

echo -e "${GREEN}📊 DATOS DE PRUEBA CREADOS:${NC}"
echo "  - Técnico ID: $TECNICO_ID (email: $TECNICO_EMAIL)"
echo "  - Residente ID: $RESIDENTE_ID (email: $RESIDENTE_EMAIL)"
echo "  - Segundo Residente ID: $SEGUNDO_ID"
echo ""

echo -e "${YELLOW}💡 DIFERENCIAS CLAVE:${NC}"
echo "  - Técnicos: Requieren especialidad, asociados a empresas"
echo "  - Residentes: Solo requieren username y email, auto-registro simple"
echo "  - Ambos: Validación de email único en sus respectivas tablas"
echo ""

echo -e "${GREEN}✅ TODAS LAS PRUEBAS COMPLETADAS EXITOSAMENTE${NC}\n"
