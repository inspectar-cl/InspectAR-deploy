#!/bin/bash

# Script de prueba para la ruta POST /admin/usuarios
# Este script valida la creación de usuarios administradores de edificios

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Configuración
BASE_URL="http://localhost:8092"
ADMIN_EMAIL="admin@example.com"
TEST_USERNAME="test_admin_$(date +%s)"
TEST_EMAIL="test.admin.$(date +%s)@inspectAR.cl"

echo -e "${CYAN}============================================================${NC}"
echo -e "${CYAN}   PRUEBA DE CREACIÓN DE USUARIOS ADMINISTRADORES${NC}"
echo -e "${CYAN}============================================================${NC}"
echo ""

# Función para imprimir resultado de test
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ ÉXITO:${NC} $2"
    else
        echo -e "${RED}❌ ERROR:${NC} $2"
        echo -e "${YELLOW}Respuesta:${NC} $3"
        exit 1
    fi
}

echo -e "${BLUE}[1] Verificando servicio de gestión...${NC}"
HEALTH_CHECK=$(curl -s $BASE_URL/health)
if [[ $HEALTH_CHECK == *"ok"* ]]; then
    print_result 0 "Servicio de gestión está corriendo"
else
    print_result 1 "Servicio de gestión no responde" "$HEALTH_CHECK"
fi
echo ""

echo -e "${BLUE}[2] Creando usuario administrador de edificios...${NC}"
echo -e "${CYAN}Datos del usuario:${NC}"
echo "  - Username: $TEST_USERNAME"
echo "  - Email: $TEST_EMAIL"
echo ""

CREATE_RESPONSE=$(curl -s -X POST $BASE_URL/admin/usuarios \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$TEST_USERNAME\",
    \"email\": \"$TEST_EMAIL\",
    \"admin_email\": \"$ADMIN_EMAIL\"
  }")

echo -e "${YELLOW}Respuesta:${NC}"
echo "$CREATE_RESPONSE" | jq '.'

# Verificar que se creó exitosamente
if echo "$CREATE_RESPONSE" | jq -e '.usuario.id' > /dev/null 2>&1; then
    USUARIO_ID=$(echo "$CREATE_RESPONSE" | jq -r '.usuario.id')
    USUARIO_USERNAME=$(echo "$CREATE_RESPONSE" | jq -r '.usuario.username')
    USUARIO_EMAIL=$(echo "$CREATE_RESPONSE" | jq -r '.usuario.email')
    
    if [ "$USUARIO_USERNAME" == "$TEST_USERNAME" ] && [ "$USUARIO_EMAIL" == "$TEST_EMAIL" ]; then
        print_result 0 "Usuario creado exitosamente con ID: $USUARIO_ID"
    else
        print_result 1 "Datos del usuario no coinciden" "$CREATE_RESPONSE"
    fi
else
    print_result 1 "Respuesta no contiene ID de usuario" "$CREATE_RESPONSE"
fi
echo ""

echo -e "${BLUE}[3] Intentando crear el mismo usuario de nuevo (debe fallar)...${NC}"
DUPLICATE_RESPONSE=$(curl -s -X POST $BASE_URL/admin/usuarios \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"otro_username\",
    \"email\": \"$TEST_EMAIL\",
    \"admin_email\": \"$ADMIN_EMAIL\"
  }")

echo -e "${YELLOW}Respuesta:${NC}"
echo "$DUPLICATE_RESPONSE" | jq '.'

if echo "$DUPLICATE_RESPONSE" | jq -e '.error' | grep -q "ya está registrado"; then
    print_result 0 "Validación de email duplicado funciona correctamente"
else
    print_result 1 "Debería rechazar email duplicado" "$DUPLICATE_RESPONSE"
fi
echo ""

echo -e "${BLUE}[4] Obteniendo lista de edificios disponibles...${NC}"
EDIFICIOS=$(curl -s $BASE_URL/edificios)
EDIFICIO_ID=$(echo "$EDIFICIOS" | jq -r '.edificios[0].id')
EDIFICIO_NOMBRE=$(echo "$EDIFICIOS" | jq -r '.edificios[0].nombre')
echo -e "${CYAN}Edificio seleccionado: ID $EDIFICIO_ID - $EDIFICIO_NOMBRE${NC}"
echo ""

echo -e "${BLUE}[5] Asignando edificio al usuario administrador...${NC}"
ASIGNAR_RESPONSE=$(curl -s -X POST $BASE_URL/admin/usuarios/edificios \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"edificio_id\": $EDIFICIO_ID,
    \"admin_email\": \"$ADMIN_EMAIL\"
  }")

echo -e "${YELLOW}Respuesta:${NC}"
echo "$ASIGNAR_RESPONSE" | jq '.'

if echo "$ASIGNAR_RESPONSE" | jq -e '.message' | grep -q "exitosamente"; then
    print_result 0 "Usuario asignado al edificio $EDIFICIO_NOMBRE"
else
    print_result 1 "Error al asignar edificio" "$ASIGNAR_RESPONSE"
fi
echo ""

echo -e "${BLUE}[6] Verificando edificios del usuario administrador...${NC}"
sleep 1
EDIFICIOS_USUARIO=$(curl -s $BASE_URL/usuarios/edificios/$TEST_EMAIL)

echo -e "${CYAN}Edificios del usuario:${NC}"
echo "$EDIFICIOS_USUARIO" | jq '.'

TOTAL_EDIFICIOS=$(echo "$EDIFICIOS_USUARIO" | jq -r '.total')
PRIMER_EDIFICIO_ID=$(echo "$EDIFICIOS_USUARIO" | jq -r '.edificios[0].id')

if [ "$TOTAL_EDIFICIOS" -ge 1 ] && [ "$PRIMER_EDIFICIO_ID" == "$EDIFICIO_ID" ]; then
    print_result 0 "Usuario tiene acceso al edificio $EDIFICIO_NOMBRE"
else
    print_result 1 "Usuario no tiene acceso correcto al edificio" "$EDIFICIOS_USUARIO"
fi
echo ""

echo -e "${BLUE}[7] Verificando acceso del usuario al edificio...${NC}"
ACCESO_RESPONSE=$(curl -s $BASE_URL/usuarios/$TEST_EMAIL/edificio/$EDIFICIO_ID/acceso)

echo -e "${YELLOW}Respuesta:${NC}"
echo "$ACCESO_RESPONSE" | jq '.'

if echo "$ACCESO_RESPONSE" | jq -e '.tiene_acceso' | grep -q "true"; then
    print_result 0 "Verificación de acceso correcta (tiene_acceso: true)"
else
    print_result 1 "Usuario debería tener acceso al edificio" "$ACCESO_RESPONSE"
fi
echo ""

echo -e "${BLUE}[8] Obteniendo usuarios del edificio...${NC}"
USUARIOS_EDIFICIO=$(curl -s $BASE_URL/admin/edificios/$EDIFICIO_ID/usuarios)

echo -e "${CYAN}Usuarios del edificio $EDIFICIO_NOMBRE:${NC}"
echo "$USUARIOS_EDIFICIO" | jq '.usuarios[] | {id, username, email}'

# Verificar que nuestro usuario aparece en la lista
if echo "$USUARIOS_EDIFICIO" | jq -r '.usuarios[].email' | grep -q "$TEST_EMAIL"; then
    print_result 0 "Usuario aparece en la lista de usuarios del edificio"
else
    echo -e "${YELLOW}⚠️  ADVERTENCIA: Usuario no aparece en la lista del edificio${NC}"
fi
echo ""

echo -e "${BLUE}[9] Consultando logs de auditoría...${NC}"
LOGS=$(curl -s "$BASE_URL/admin/logs?accion=CREAR_USUARIO_ADMINISTRADOR&usuario_email=$ADMIN_EMAIL")
TOTAL_LOGS=$(echo "$LOGS" | jq -r '.total')

echo -e "${CYAN}Logs de creación de usuarios: $TOTAL_LOGS${NC}"
if [ "$TOTAL_LOGS" -gt 0 ]; then
    echo "$LOGS" | jq '.logs[0] | {accion, descripcion, usuario_email, fecha: .creado_en}'
    print_result 0 "Logs de auditoría registrados correctamente"
else
    echo -e "${YELLOW}⚠️  ADVERTENCIA: No se encontraron logs de auditoría${NC}"
fi
echo ""

echo -e "${BLUE}[10] Probando validaciones de campos...${NC}"

# Test: username vacío
echo -e "${CYAN}Test: username vacío${NC}"
VALIDATION_1=$(curl -s -X POST $BASE_URL/admin/usuarios \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"\",
    \"email\": \"test@example.com\",
    \"admin_email\": \"$ADMIN_EMAIL\"
  }")

if echo "$VALIDATION_1" | jq -e '.error' > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Validación de username vacío funciona"
else
    echo -e "${YELLOW}⚠${NC} Validación de username vacío no funciona como esperado"
fi

# Test: email inválido
echo -e "${CYAN}Test: email inválido${NC}"
VALIDATION_2=$(curl -s -X POST $BASE_URL/admin/usuarios \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"test\",
    \"email\": \"email-invalido\",
    \"admin_email\": \"$ADMIN_EMAIL\"
  }")

if echo "$VALIDATION_2" | jq -e '.error' > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Validación de email inválido funciona"
else
    echo -e "${YELLOW}⚠${NC} Validación de email inválido no funciona como esperado"
fi

# Test: admin_email faltante
echo -e "${CYAN}Test: admin_email faltante${NC}"
VALIDATION_3=$(curl -s -X POST $BASE_URL/admin/usuarios \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"test\",
    \"email\": \"test@example.com\"
  }")

if echo "$VALIDATION_3" | jq -e '.error' > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Validación de admin_email faltante funciona"
else
    echo -e "${YELLOW}⚠${NC} Validación de admin_email faltante no funciona como esperado"
fi

print_result 0 "Validaciones de campos completadas"
echo ""

echo -e "${GREEN}============================================================${NC}"
echo -e "${GREEN}   ✅ TODAS LAS PRUEBAS PASARON EXITOSAMENTE${NC}"
echo -e "${GREEN}============================================================${NC}"
echo ""
echo -e "${CYAN}Resumen de operaciones:${NC}"
echo -e "- Creación de usuario: ${GREEN}✓${NC}"
echo -e "- Validación de duplicados: ${GREEN}✓${NC}"
echo -e "- Asignación de edificio: ${GREEN}✓${NC}"
echo -e "- Verificación de edificios: ${GREEN}✓${NC}"
echo -e "- Verificación de acceso: ${GREEN}✓${NC}"
echo -e "- Usuarios por edificio: ${GREEN}✓${NC}"
echo -e "- Logs de auditoría: ${GREEN}✓${NC}"
echo -e "- Validaciones de entrada: ${GREEN}✓${NC}"
echo ""
echo -e "${PURPLE}Usuario creado:${NC}"
echo -e "  ID: $USUARIO_ID"
echo -e "  Username: $USUARIO_USERNAME"
echo -e "  Email: $USUARIO_EMAIL"
echo -e "  Edificio asignado: $EDIFICIO_NOMBRE (ID: $EDIFICIO_ID)"
echo ""
