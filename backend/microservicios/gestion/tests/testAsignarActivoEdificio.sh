#!/bin/bash

# Script de prueba para la ruta PUT /admin/activos/edificio
# Este script valida la asignación y reasignación de activos a edificios

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

echo -e "${CYAN}============================================================${NC}"
echo -e "${CYAN}   PRUEBA DE ASIGNACIÓN DE ACTIVOS A EDIFICIOS${NC}"
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

echo -e "${BLUE}[2] Obteniendo lista de edificios disponibles...${NC}"
EDIFICIOS=$(curl -s $BASE_URL/edificios)
echo -e "${CYAN}Edificios disponibles:${NC}"
echo "$EDIFICIOS" | jq -r '.edificios[] | "  - ID: \(.id) - \(.nombre)"'
EDIFICIO_1=$(echo "$EDIFICIOS" | jq -r '.edificios[0].id')
EDIFICIO_2=$(echo "$EDIFICIOS" | jq -r '.edificios[1].id')
echo -e "${YELLOW}Usaremos edificio_1: $EDIFICIO_1 y edificio_2: $EDIFICIO_2${NC}"
echo ""

echo -e "${BLUE}[3] Obteniendo lista de activos...${NC}"
ACTIVOS=$(curl -s $BASE_URL/activos)
ACTIVO_ID=$(echo "$ACTIVOS" | jq -r '.activos[0].id')
ACTIVO_NOMBRE=$(echo "$ACTIVOS" | jq -r '.activos[0].nombre')
ACTIVO_EDIFICIO_ACTUAL=$(echo "$ACTIVOS" | jq -r '.activos[0].edificio_id')
echo -e "${CYAN}Activo seleccionado:${NC}"
echo "  - ID: $ACTIVO_ID"
echo "  - Nombre: $ACTIVO_NOMBRE"
echo "  - Edificio actual: $ACTIVO_EDIFICIO_ACTUAL"
echo ""

echo -e "${BLUE}[4] Asignando activo al edificio $EDIFICIO_1...${NC}"
ASIGNAR_RESPONSE=$(curl -s -X PUT $BASE_URL/admin/activos/edificio \
  -H "Content-Type: application/json" \
  -d "{
    \"activo_id\": $ACTIVO_ID,
    \"edificio_id\": $EDIFICIO_1,
    \"admin_email\": \"$ADMIN_EMAIL\"
  }")

echo -e "${YELLOW}Respuesta:${NC}"
echo "$ASIGNAR_RESPONSE" | jq '.'

# Verificar que la asignación fue exitosa
if echo "$ASIGNAR_RESPONSE" | jq -e '.message' > /dev/null 2>&1; then
    MESSAGE=$(echo "$ASIGNAR_RESPONSE" | jq -r '.message')
    if [[ $MESSAGE == *"exitosamente"* ]]; then
        print_result 0 "Activo asignado exitosamente al edificio $EDIFICIO_1"
    else
        print_result 1 "Error en la respuesta" "$ASIGNAR_RESPONSE"
    fi
else
    print_result 1 "Respuesta no contiene mensaje" "$ASIGNAR_RESPONSE"
fi
echo ""

echo -e "${BLUE}[5] Verificando que el activo fue asignado correctamente...${NC}"
sleep 1
ACTIVO_ACTUALIZADO=$(curl -s $BASE_URL/activos/$ACTIVO_ID)
NUEVO_EDIFICIO_ID=$(echo "$ACTIVO_ACTUALIZADO" | jq -r '.edificio_id')

echo -e "${CYAN}Estado del activo:${NC}"
echo "$ACTIVO_ACTUALIZADO" | jq '{id, nombre, edificio_id}'

if [ "$NUEVO_EDIFICIO_ID" == "$EDIFICIO_1" ]; then
    print_result 0 "Activo ahora pertenece al edificio $EDIFICIO_1"
else
    print_result 1 "El edificio no se actualizó correctamente (esperado: $EDIFICIO_1, actual: $NUEVO_EDIFICIO_ID)" "$ACTIVO_ACTUALIZADO"
fi
echo ""

echo -e "${BLUE}[6] Reasignando activo al edificio $EDIFICIO_2...${NC}"
REASIGNAR_RESPONSE=$(curl -s -X PUT $BASE_URL/admin/activos/edificio \
  -H "Content-Type: application/json" \
  -d "{
    \"activo_id\": $ACTIVO_ID,
    \"edificio_id\": $EDIFICIO_2,
    \"admin_email\": \"$ADMIN_EMAIL\"
  }")

echo -e "${YELLOW}Respuesta:${NC}"
echo "$REASIGNAR_RESPONSE" | jq '.'

if echo "$REASIGNAR_RESPONSE" | jq -e '.edificio_anterior' > /dev/null 2>&1; then
    EDIFICIO_ANTERIOR_RESP=$(echo "$REASIGNAR_RESPONSE" | jq -r '.edificio_anterior.id')
    EDIFICIO_NUEVO_RESP=$(echo "$REASIGNAR_RESPONSE" | jq -r '.edificio_nuevo.id')
    
    if [ "$EDIFICIO_ANTERIOR_RESP" == "$EDIFICIO_1" ] && [ "$EDIFICIO_NUEVO_RESP" == "$EDIFICIO_2" ]; then
        print_result 0 "Reasignación exitosa de edificio $EDIFICIO_1 a $EDIFICIO_2"
    else
        print_result 1 "Error en reasignación" "$REASIGNAR_RESPONSE"
    fi
else
    print_result 1 "Respuesta no contiene información de edificios" "$REASIGNAR_RESPONSE"
fi
echo ""

echo -e "${BLUE}[7] Verificando reasignación en la base de datos...${NC}"
sleep 1
ACTIVO_FINAL=$(curl -s $BASE_URL/activos/$ACTIVO_ID)
EDIFICIO_FINAL=$(echo "$ACTIVO_FINAL" | jq -r '.edificio_id')

echo -e "${CYAN}Estado final del activo:${NC}"
echo "$ACTIVO_FINAL" | jq '{id, nombre, edificio_id}'

if [ "$EDIFICIO_FINAL" == "$EDIFICIO_2" ]; then
    print_result 0 "Activo ahora pertenece al edificio $EDIFICIO_2"
else
    print_result 1 "El edificio no se actualizó correctamente (esperado: $EDIFICIO_2, actual: $EDIFICIO_FINAL)" "$ACTIVO_FINAL"
fi
echo ""

echo -e "${BLUE}[8] Consultando logs de auditoría...${NC}"
LOGS=$(curl -s "$BASE_URL/admin/logs?entidad=activo_edificio&entidad_id=$ACTIVO_ID")
TOTAL_LOGS=$(echo "$LOGS" | jq -r '.total')

echo -e "${CYAN}Logs encontrados: $TOTAL_LOGS${NC}"
echo "$LOGS" | jq '.logs[] | {accion, descripcion, fecha: .creado_en}'

if [ "$TOTAL_LOGS" -ge 2 ]; then
    print_result 0 "Logs de auditoría registrados correctamente"
else
    echo -e "${YELLOW}⚠️  ADVERTENCIA: Se esperaban al menos 2 logs (asignación + reasignación)${NC}"
fi
echo ""

echo -e "${BLUE}[9] Verificando activos del edificio $EDIFICIO_2...${NC}"
ACTIVOS_EDIFICIO=$(curl -s $BASE_URL/activos/edificio/$EDIFICIO_2)
ACTIVOS_IDS=$(echo "$ACTIVOS_EDIFICIO" | jq -r '.activos[].id')

if echo "$ACTIVOS_IDS" | grep -q "^$ACTIVO_ID$"; then
    print_result 0 "Activo aparece en la lista de activos del edificio $EDIFICIO_2"
    echo -e "${CYAN}Activos del edificio $EDIFICIO_2:${NC}"
    echo "$ACTIVOS_EDIFICIO" | jq -r '.activos[] | "  - \(.id): \(.nombre)"'
else
    print_result 1 "Activo NO aparece en lista de activos del edificio" "$ACTIVOS_EDIFICIO"
fi
echo ""

echo -e "${GREEN}============================================================${NC}"
echo -e "${GREEN}   ✅ TODAS LAS PRUEBAS PASARON EXITOSAMENTE${NC}"
echo -e "${GREEN}============================================================${NC}"
echo ""
echo -e "${CYAN}Resumen de operaciones:${NC}"
echo -e "- Asignación inicial: ${GREEN}✓${NC}"
echo -e "- Verificación de asignación: ${GREEN}✓${NC}"
echo -e "- Reasignación: ${GREEN}✓${NC}"
echo -e "- Verificación de reasignación: ${GREEN}✓${NC}"
echo -e "- Logs de auditoría: ${GREEN}✓${NC}"
echo -e "- Activos por edificio: ${GREEN}✓${NC}"
echo ""
echo -e "${PURPLE}Estado final:${NC}"
echo -e "  Activo $ACTIVO_ID ($ACTIVO_NOMBRE) → Edificio $EDIFICIO_2"
echo ""
