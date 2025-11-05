#!/bin/bash

# Script de prueba para la ruta GET /edificios
# Este script valida que la ruta devuelva todos los edificios disponibles

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuración
BASE_URL="http://localhost:8092"

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}   PRUEBA DE RUTA GET /edificios${NC}"
echo -e "${CYAN}========================================${NC}"
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

echo -e "${BLUE}[2] Obteniendo todos los edificios...${NC}"
EDIFICIOS_RESPONSE=$(curl -s $BASE_URL/edificios)
echo -e "${YELLOW}Respuesta:${NC}"
echo "$EDIFICIOS_RESPONSE" | jq '.'

# Verificar que la respuesta contiene el campo 'edificios'
if echo "$EDIFICIOS_RESPONSE" | jq -e '.edificios' > /dev/null 2>&1; then
    print_result 0 "Campo 'edificios' presente en la respuesta"
else
    print_result 1 "Campo 'edificios' no presente en la respuesta" "$EDIFICIOS_RESPONSE"
fi
echo ""

echo -e "${BLUE}[3] Verificando campo 'total'...${NC}"
TOTAL=$(echo "$EDIFICIOS_RESPONSE" | jq -r '.total')
echo -e "${CYAN}Total de edificios:${NC} $TOTAL"

if [ "$TOTAL" -gt 0 ]; then
    print_result 0 "Total de edificios es mayor a 0 ($TOTAL edificios)"
else
    print_result 1 "Total de edificios es 0" "$EDIFICIOS_RESPONSE"
fi
echo ""

echo -e "${BLUE}[4] Verificando estructura de edificios...${NC}"
PRIMER_EDIFICIO=$(echo "$EDIFICIOS_RESPONSE" | jq -r '.edificios[0]')

# Verificar campos requeridos
TIENE_ID=$(echo "$PRIMER_EDIFICIO" | jq -e '.id' > /dev/null 2>&1 && echo "si" || echo "no")
TIENE_NOMBRE=$(echo "$PRIMER_EDIFICIO" | jq -e '.nombre' > /dev/null 2>&1 && echo "si" || echo "no")
TIENE_DIRECCION=$(echo "$PRIMER_EDIFICIO" | jq -e '.direccion' > /dev/null 2>&1 && echo "si" || echo "no")
TIENE_CREADO=$(echo "$PRIMER_EDIFICIO" | jq -e '.creado_en' > /dev/null 2>&1 && echo "si" || echo "no")

if [[ "$TIENE_ID" == "si" && "$TIENE_NOMBRE" == "si" && "$TIENE_DIRECCION" == "si" && "$TIENE_CREADO" == "si" ]]; then
    print_result 0 "Estructura de edificio es correcta (id, nombre, direccion, creado_en)"
    echo -e "${CYAN}Ejemplo de edificio:${NC}"
    echo "$PRIMER_EDIFICIO" | jq '.'
else
    print_result 1 "Estructura de edificio incorrecta" "$PRIMER_EDIFICIO"
fi
echo ""

echo -e "${BLUE}[5] Mostrando todos los edificios...${NC}"
echo -e "${CYAN}Lista completa de edificios:${NC}"
echo "$EDIFICIOS_RESPONSE" | jq -r '.edificios[] | "- \(.id): \(.nombre) (\(.direccion))"'
echo ""

echo -e "${BLUE}[6] Verificando orden alfabético...${NC}"
NOMBRES=$(echo "$EDIFICIOS_RESPONSE" | jq -r '.edificios[].nombre')
NOMBRES_ORDENADOS=$(echo "$NOMBRES" | sort)

if [ "$NOMBRES" == "$NOMBRES_ORDENADOS" ]; then
    print_result 0 "Edificios están ordenados alfabéticamente por nombre"
else
    echo -e "${YELLOW}⚠️  ADVERTENCIA: Edificios no están en orden alfabético${NC}"
    echo -e "${CYAN}Orden actual:${NC}"
    echo "$NOMBRES"
    echo -e "${CYAN}Orden esperado:${NC}"
    echo "$NOMBRES_ORDENADOS"
fi
echo ""

echo -e "${BLUE}[7] Verificando tiempo de respuesta...${NC}"
START_TIME=$(date +%s%N)
curl -s $BASE_URL/edificios > /dev/null
END_TIME=$(date +%s%N)
DURATION=$(( (END_TIME - START_TIME) / 1000000 ))

echo -e "${CYAN}Tiempo de respuesta:${NC} ${DURATION}ms"
if [ $DURATION -lt 1000 ]; then
    print_result 0 "Tiempo de respuesta aceptable (< 1000ms)"
else
    echo -e "${YELLOW}⚠️  ADVERTENCIA: Tiempo de respuesta alto (${DURATION}ms)${NC}"
fi
echo ""

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   ✅ TODAS LAS PRUEBAS PASARON${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${CYAN}Resumen:${NC}"
echo -e "- Servicio funcionando: ${GREEN}✓${NC}"
echo -e "- Endpoint /edificios: ${GREEN}✓${NC}"
echo -e "- Estructura correcta: ${GREEN}✓${NC}"
echo -e "- Total de edificios: ${CYAN}$TOTAL${NC}"
echo -e "- Tiempo de respuesta: ${CYAN}${DURATION}ms${NC}"
echo ""
