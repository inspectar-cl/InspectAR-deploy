#!/bin/bash

# Script de test para la funcionalidad de filtrado de activos por edificio
# ParserService - HU Filtrado por Edificio

# Colores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8090"
ERRORES=0

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   TEST: Activos por Edificio${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Función para verificar respuesta
check_response() {
    local response=$1
    local expected_count=$2
    local edificio_id=$3
    
    # Verificar que la respuesta contiene el edificio_id
    if echo "$response" | grep -q "\"edificio_id\":$edificio_id"; then
        echo -e "${GREEN}✓${NC} Edificio ID correcto: $edificio_id"
    else
        echo -e "${RED}✗${NC} Edificio ID incorrecto"
        ((ERRORES++))
    fi
    
    # Verificar el total de activos
    actual_count=$(echo "$response" | grep -o '"total":[0-9]*' | grep -o '[0-9]*')
    if [ "$actual_count" == "$expected_count" ]; then
        echo -e "${GREEN}✓${NC} Total de activos correcto: $actual_count"
    else
        echo -e "${RED}✗${NC} Total esperado: $expected_count, obtenido: $actual_count"
        ((ERRORES++))
    fi
}

# FASE 1: Health Check
echo -e "${YELLOW}[FASE 1]${NC} Verificando servicio..."
HEALTH=$(curl -s $BASE_URL/healthz)
if echo "$HEALTH" | grep -q '"ok":true'; then
    echo -e "${GREEN}✓${NC} Servicio funcionando correctamente"
else
    echo -e "${RED}✗${NC} Servicio no disponible"
    exit 1
fi
echo ""

# FASE 2: Crear activos de prueba en diferentes edificios
echo -e "${YELLOW}[FASE 2]${NC} Creando activos de prueba..."

# Edificio 1 - 3 activos (se suman a los 3 existentes: 1, 2, 7)
echo "Creando activos para Edificio 1..."
curl -s -X POST $BASE_URL/activo \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 101,
    "estado": "operativo",
    "edificio_id": 1
  }' > /dev/null

curl -s -X POST $BASE_URL/activo \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 102,
    "estado": "operativo",
    "edificio_id": 1
  }' > /dev/null

curl -s -X POST $BASE_URL/activo \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 103,
    "estado": "operativo",
    "edificio_id": 1
  }' > /dev/null

echo -e "${GREEN}✓${NC} 3 activos creados para Edificio 1"

# Edificio 2 - 2 activos (se suman a los 2 existentes: 3, 4)
echo "Creando activos para Edificio 2..."
curl -s -X POST $BASE_URL/activo \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 201,
    "estado": "operativo",
    "edificio_id": 2
  }' > /dev/null

curl -s -X POST $BASE_URL/activo \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 202,
    "estado": "mantenimiento",
    "edificio_id": 2
  }' > /dev/null

echo -e "${GREEN}✓${NC} 2 activos creados para Edificio 2"

# Edificio 3 - 1 activo (se suma a los 2 existentes: 5, 6)
echo "Creando activos para Edificio 3..."
curl -s -X POST $BASE_URL/activo \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 301,
    "estado": "stand_by",
    "edificio_id": 3
  }' > /dev/null

echo -e "${GREEN}✓${NC} 1 activo creado para Edificio 3"
echo ""

# FASE 3: Probar filtrado por edificio
echo -e "${YELLOW}[FASE 3]${NC} Probando filtrado por edificio..."

# Test Edificio 1 (debe retornar 6 activos: 3 iniciales + 3 nuevos)
echo ""
echo "Test 1: Obtener activos del Edificio 1 (esperados: 6)"
RESPONSE_1=$(curl -s $BASE_URL/activo/edificio/1)
check_response "$RESPONSE_1" 6 1

# Test Edificio 2 (debe retornar 4 activos: 2 iniciales + 2 nuevos)
echo ""
echo "Test 2: Obtener activos del Edificio 2 (esperados: 4)"
RESPONSE_2=$(curl -s $BASE_URL/activo/edificio/2)
check_response "$RESPONSE_2" 4 2

# Test Edificio 3 (debe retornar 3 activos: 2 iniciales + 1 nuevo)
echo ""
echo "Test 3: Obtener activos del Edificio 3 (esperados: 3)"
RESPONSE_3=$(curl -s $BASE_URL/activo/edificio/3)
check_response "$RESPONSE_3" 3 3

# Test Edificio sin activos (debe retornar 0)
echo ""
echo "Test 4: Obtener activos del Edificio 99 (esperados: 0)"
RESPONSE_99=$(curl -s $BASE_URL/activo/edificio/99)
check_response "$RESPONSE_99" 0 99

echo ""

# FASE 4: Verificar estructura de respuesta
echo -e "${YELLOW}[FASE 4]${NC} Verificando estructura de respuesta..."

# Verificar que contiene los campos esperados
if echo "$RESPONSE_1" | grep -q '"edificio_id"' && \
   echo "$RESPONSE_1" | grep -q '"activos"' && \
   echo "$RESPONSE_1" | grep -q '"total"'; then
    echo -e "${GREEN}✓${NC} Estructura de respuesta correcta"
else
    echo -e "${RED}✗${NC} Estructura de respuesta incorrecta"
    ((ERRORES++))
fi

# Verificar que los activos tienen los campos necesarios (nueva estructura simplificada)
if echo "$RESPONSE_1" | grep -q '"activo_id"' && \
   echo "$RESPONSE_1" | grep -q '"estado"' && \
   echo "$RESPONSE_1" | grep -q '"edificio_id"'; then
    echo -e "${GREEN}✓${NC} Campos de activo correctos (activo_id, estado, edificio_id)"
else
    echo -e "${RED}✗${NC} Campos de activo faltantes"
    ((ERRORES++))
fi

echo ""

# FASE 5: Verificar que los activos pertenecen al edificio correcto
echo -e "${YELLOW}[FASE 5]${NC} Verificando correspondencia edificio-activos..."

# Extraer activos del Edificio 1 y verificar que todos tienen edificio_id = 1
# Contar solo los activo_id dentro del array "activos"
ACTIVOS_ED1=$(echo "$RESPONSE_1" | jq '.activos | length')
# Contar cuántos tienen edificio_id = 1
CORRECT_ED1=$(echo "$RESPONSE_1" | jq '[.activos[] | select(.edificio_id == 1)] | length')

if [ "$ACTIVOS_ED1" -eq "$CORRECT_ED1" ] && [ "$CORRECT_ED1" -ge 6 ]; then
    echo -e "${GREEN}✓${NC} Todos los $ACTIVOS_ED1 activos del Edificio 1 tienen edificio_id correcto"
else
    echo -e "${RED}✗${NC} Problema con edificio_id. Activos totales: $ACTIVOS_ED1, Con edificio_id=1: $CORRECT_ED1"
    ((ERRORES++))
fi

echo ""

# RESUMEN FINAL
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}         RESUMEN DE TESTS${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "Total de activos creados en este test: 6"
echo "Total de activos en la BD (incluye iniciales):"
echo "  - Edificio 1: 6 activos (3 iniciales + 3 nuevos)"
echo "  - Edificio 2: 4 activos (2 iniciales + 2 nuevos)"
echo "  - Edificio 3: 3 activos (2 iniciales + 1 nuevo)"
echo "  - Edificio 4: 1 activo (inicial)"
echo ""

if [ $ERRORES -eq 0 ]; then
    echo -e "${GREEN}✓ TODOS LOS TESTS PASARON${NC}"
    echo -e "${GREEN}✓ La funcionalidad de filtrado por edificio funciona correctamente${NC}"
    exit 0
else
    echo -e "${RED}✗ $ERRORES ERROR(ES) DETECTADO(S)${NC}"
    echo -e "${RED}✗ Revisa los logs para más detalles${NC}"
    exit 1
fi
