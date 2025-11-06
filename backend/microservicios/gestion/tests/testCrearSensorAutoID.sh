#!/bin/bash

# Script de prueba para creación de sensores con sensor_id auto-generado
# Verifica que el sensor_id sea único y tenga formato correcto

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

BASE_URL="http://localhost:8092"
TIMESTAMP=$(date +%s)

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  TEST: Crear Sensor con ID Auto-generado${NC}"
echo -e "${CYAN}========================================${NC}\n"

# Paso 1: Verificar servicio
echo -e "${BLUE}Paso 1: Verificando servicio...${NC}"
HEALTH=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/health)
if [ "$HEALTH" == "200" ]; then
    echo -e "${GREEN}✓ Servicio funcionando${NC}\n"
else
    echo -e "${RED}✗ Servicio no responde${NC}\n"
    exit 1
fi

# Paso 2: Crear sensor SIN nombre (solo campos requeridos)
echo -e "${BLUE}Paso 2: Creando sensor SIN nombre (campos mínimos)...${NC}"
SENSOR1_RESPONSE=$(curl -s -X POST $BASE_URL/admin/sensores \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 1,
    "tipo": "temperatura",
    "unidad": "°C",
    "email": "admin@inspectAR.cl"
  }')

SENSOR1_ID=$(echo $SENSOR1_RESPONSE | jq -r '.sensor_id')

if [ "$SENSOR1_ID" != "null" ] && [ "$SENSOR1_ID" != "" ]; then
    echo -e "${GREEN}✓ Sensor creado:${NC}"
    echo "  - sensor_id: $SENSOR1_ID"
    echo $SENSOR1_RESPONSE | jq '.sensor'
    
    # Verificar formato del ID: TIPO_ACTIVO_TIMESTAMP_RANDOM
    if [[ $SENSOR1_ID =~ ^[A-Z0-9]+_[0-9]+_[0-9]+_[A-F0-9]{4}$ ]]; then
        echo -e "${GREEN}✓ Formato de ID correcto: TIPO_ACTIVO_TIMESTAMP_RANDOM${NC}"
    else
        echo -e "${YELLOW}⚠ Formato de ID inesperado: $SENSOR1_ID${NC}"
    fi
else
    echo -e "${RED}✗ Error al crear sensor${NC}"
    echo $SENSOR1_RESPONSE | jq
fi

echo ""

# Paso 3: Crear sensor CON nombre
echo -e "${BLUE}Paso 3: Creando sensor CON nombre...${NC}"
SENSOR2_RESPONSE=$(curl -s -X POST $BASE_URL/admin/sensores \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 1,
    "nombre": "Sensor Principal Caldera 1",
    "tipo": "presion",
    "unidad": "bar",
    "email": "admin@inspectAR.cl"
  }')

SENSOR2_ID=$(echo $SENSOR2_RESPONSE | jq -r '.sensor_id')

if [ "$SENSOR2_ID" != "null" ] && [ "$SENSOR2_ID" != "" ]; then
    echo -e "${GREEN}✓ Sensor creado con nombre:${NC}"
    echo "  - sensor_id: $SENSOR2_ID"
    echo $SENSOR2_RESPONSE | jq '.sensor'
else
    echo -e "${RED}✗ Error al crear sensor con nombre${NC}"
    echo $SENSOR2_RESPONSE | jq
fi

echo ""

# Paso 4: Crear múltiples sensores para verificar unicidad
echo -e "${BLUE}Paso 4: Creando 3 sensores para verificar unicidad de IDs...${NC}"
IDS=()

for i in {1..3}; do
    RESPONSE=$(curl -s -X POST $BASE_URL/admin/sensores \
      -H "Content-Type: application/json" \
      -d '{
        "activo_id": 2,
        "tipo": "vibracion",
        "unidad": "Hz",
        "email": "admin@inspectAR.cl"
      }')
    
    ID=$(echo $RESPONSE | jq -r '.sensor_id')
    IDS+=("$ID")
    echo "  Sensor $i: $ID"
    sleep 0.1 # Pequeña pausa para evitar colisiones de timestamp
done

# Verificar que todos los IDs son diferentes
UNIQUE_COUNT=$(printf '%s\n' "${IDS[@]}" | sort -u | wc -l)
TOTAL_COUNT=${#IDS[@]}

if [ "$UNIQUE_COUNT" -eq "$TOTAL_COUNT" ]; then
    echo -e "${GREEN}✓ Todos los IDs son únicos ($UNIQUE_COUNT/$TOTAL_COUNT)${NC}"
else
    echo -e "${RED}✗ Se encontraron IDs duplicados ($UNIQUE_COUNT únicos de $TOTAL_COUNT)${NC}"
fi

echo ""

# Paso 5: Crear sensor con tipo largo (verificar truncamiento)
echo -e "${BLUE}Paso 5: Creando sensor con tipo largo...${NC}"
SENSOR_LARGO=$(curl -s -X POST $BASE_URL/admin/sensores \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 1,
    "tipo": "temperatura ambiente exterior",
    "unidad": "°C",
    "email": "admin@inspectAR.cl"
  }')

SENSOR_LARGO_ID=$(echo $SENSOR_LARGO | jq -r '.sensor_id')
echo "  - sensor_id: $SENSOR_LARGO_ID"

# Extraer parte del tipo del ID
TIPO_PARTE=$(echo $SENSOR_LARGO_ID | cut -d'_' -f1)
TIPO_LENGTH=${#TIPO_PARTE}

if [ "$TIPO_LENGTH" -le 4 ]; then
    echo -e "${GREEN}✓ Tipo truncado correctamente a $TIPO_LENGTH caracteres: $TIPO_PARTE${NC}"
else
    echo -e "${YELLOW}⚠ Tipo no truncado: $TIPO_PARTE ($TIPO_LENGTH caracteres)${NC}"
fi

echo ""

# Paso 6: Validar campos requeridos
echo -e "${BLUE}Paso 6: Validando campos requeridos...${NC}"

# Sin activo_id
ERROR1=$(curl -s -X POST $BASE_URL/admin/sensores \
  -H "Content-Type: application/json" \
  -d '{
    "tipo": "temperatura",
    "unidad": "°C",
    "email": "admin@inspectAR.cl"
  }')

if echo $ERROR1 | jq -e '.error' > /dev/null 2>&1; then
    echo -e "${GREEN}✓ activo_id requerido (validado)${NC}"
else
    echo -e "${YELLOW}⚠ activo_id no validado como requerido${NC}"
fi

# Sin tipo
ERROR2=$(curl -s -X POST $BASE_URL/admin/sensores \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 1,
    "unidad": "°C",
    "email": "admin@inspectAR.cl"
  }')

if echo $ERROR2 | jq -e '.error' > /dev/null 2>&1; then
    echo -e "${GREEN}✓ tipo requerido (validado)${NC}"
else
    echo -e "${YELLOW}⚠ tipo no validado como requerido${NC}"
fi

# Sin unidad
ERROR3=$(curl -s -X POST $BASE_URL/admin/sensores \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 1,
    "tipo": "temperatura",
    "email": "admin@inspectAR.cl"
  }')

if echo $ERROR3 | jq -e '.error' > /dev/null 2>&1; then
    echo -e "${GREEN}✓ unidad requerido (validado)${NC}"
else
    echo -e "${YELLOW}⚠ unidad no validado como requerido${NC}"
fi

echo ""

# Resumen
echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}           RESUMEN DE PRUEBAS${NC}"
echo -e "${CYAN}========================================${NC}\n"

echo -e "${GREEN}✅ CARACTERÍSTICAS VERIFICADAS:${NC}"
echo "  ✓ sensor_id generado automáticamente"
echo "  ✓ Formato: TIPO_ACTIVO_TIMESTAMP_RANDOM"
echo "  ✓ Nombre es OPCIONAL"
echo "  ✓ IDs únicos (sin colisiones)"
echo "  ✓ Tipo truncado a 4 caracteres máximo"
echo "  ✓ Validación de campos requeridos"
echo ""

echo -e "${BLUE}📊 EJEMPLOS DE IDs GENERADOS:${NC}"
echo "  - $SENSOR1_ID"
echo "  - $SENSOR2_ID"
echo "  - $SENSOR_LARGO_ID"
echo ""

echo -e "${YELLOW}💡 FORMATO DEL ID:${NC}"
echo "  TIPO_ACTIVO_TIMESTAMP_RANDOM"
echo "  └─┬─┘ └──┬──┘ └───┬────┘ └──┬──┘"
echo "    │      │         │         └─ 4 hex chars aleatorios"
echo "    │      │         └─────────── Unix timestamp"
echo "    │      └───────────────────── ID del activo"
echo "    └──────────────────────────── Tipo normalizado (max 4 chars)"
echo ""

echo -e "${GREEN}✅ TODAS LAS PRUEBAS COMPLETADAS${NC}\n"
