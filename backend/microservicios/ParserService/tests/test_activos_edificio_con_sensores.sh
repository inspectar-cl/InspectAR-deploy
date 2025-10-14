#!/bin/bash

# Test para verificar que /activo/edificio/:edificio_id retorna información de sensores

BASE_URL="http://localhost:8090"
EDIFICIO_ID=1

echo "=========================================="
echo "Test: Activos por Edificio con Sensores"
echo "=========================================="
echo ""

# Colores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Paso 1: Crear activos en el edificio 1
echo -e "${YELLOW}Paso 1: Crear activos en edificio 1${NC}"
echo "Creando activo 1 (Caldera Principal)..."
RESPONSE1=$(curl -s -X POST $BASE_URL/activo \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 101,
    "estado": "operativo",
    "edificio_id": 1,
    "sensores": [
      {"sensor_id": "TEMP_101", "tipo": "temperatura", "unidad": "°C"},
      {"sensor_id": "PRESS_101", "tipo": "presion", "unidad": "bar"}
    ]
  }')
echo "Respuesta: $RESPONSE1"
echo ""

echo "Creando activo 2 (Bomba Centrífuga)..."
RESPONSE2=$(curl -s -X POST $BASE_URL/activo \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 102,
    "estado": "operativo",
    "edificio_id": 1,
    "sensores": [
      {"sensor_id": "FLOW_102", "tipo": "flujo", "unidad": "L/min"},
      {"sensor_id": "VIBR_102", "tipo": "vibracion", "unidad": "Hz"}
    ]
  }')
echo "Respuesta: $RESPONSE2"
echo ""

# Paso 2: Crear activo en edificio 2 (para verificar filtrado)
echo -e "${YELLOW}Paso 2: Crear activo en edificio 2 (control)${NC}"
echo "Creando activo 3 en edificio 2..."
RESPONSE3=$(curl -s -X POST $BASE_URL/activo \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 103,
    "estado": "operativo",
    "edificio_id": 2,
    "sensores": [
      {"sensor_id": "TEMP_103", "tipo": "temperatura", "unidad": "°C"}
    ]
  }')
echo "Respuesta: $RESPONSE3"
echo ""

# Esperar un momento para que se creen los activos
sleep 2

# Paso 3: Enviar lecturas para algunos sensores (simular sensores conectados)
echo -e "${YELLOW}Paso 3: Enviar lecturas de sensores (solo algunos)${NC}"
echo "Enviando lectura para TEMP_101..."
curl -s -X POST $BASE_URL/lectura \
  -H "Content-Type: application/json" \
  -d '{
    "sensor_id": "TEMP_101",
    "valor": 75.5,
    "timestamp": "2025-10-11T12:00:00Z"
  }' > /dev/null
echo "✓ Lectura enviada"

echo "Enviando lectura para FLOW_102..."
curl -s -X POST $BASE_URL/lectura \
  -H "Content-Type: application/json" \
  -d '{
    "sensor_id": "FLOW_102",
    "valor": 120.0,
    "timestamp": "2025-10-11T12:00:00Z"
  }' > /dev/null
echo "✓ Lectura enviada"
echo ""

# Esperar para que se actualice el estado de los sensores
sleep 2

# Paso 4: Obtener activos del edificio 1 con información de sensores
echo -e "${YELLOW}Paso 4: Obtener activos del edificio 1 con información de sensores${NC}"
echo "GET $BASE_URL/activo/edificio/$EDIFICIO_ID"
echo ""

RESULTADO=$(curl -s -X GET "$BASE_URL/activo/edificio/$EDIFICIO_ID")

echo "Respuesta completa:"
echo "$RESULTADO" | jq '.'
echo ""

# Paso 5: Validaciones
echo -e "${YELLOW}Paso 5: Validaciones${NC}"

# Verificar que retorna activos del edificio correcto
TOTAL=$(echo "$RESULTADO" | jq '.total')
if [ "$TOTAL" -ge 2 ]; then
    echo -e "${GREEN}✓ Total de activos correcto: $TOTAL (esperado >= 2)${NC}"
else
    echo -e "${RED}✗ Total de activos incorrecto: $TOTAL (esperado >= 2)${NC}"
fi

# Verificar que cada activo tiene información de sensores
ACTIVOS_CON_SENSORES=$(echo "$RESULTADO" | jq '[.activos[] | select(.sensores != null and (.sensores | length) > 0)] | length')
echo -e "${GREEN}✓ Activos con información de sensores: $ACTIVOS_CON_SENSORES${NC}"

# Verificar que los sensores tienen estado
PRIMER_ACTIVO_SENSORES=$(echo "$RESULTADO" | jq '.activos[0].sensores')
if [ "$PRIMER_ACTIVO_SENSORES" != "null" ] && [ "$PRIMER_ACTIVO_SENSORES" != "[]" ]; then
    echo -e "${GREEN}✓ Los activos incluyen información detallada de sensores${NC}"
    
    # Mostrar ejemplo de sensor con estado
    echo ""
    echo "Ejemplo de sensor con estado:"
    echo "$RESULTADO" | jq '.activos[0].sensores[0]'
    echo ""
    
    # Verificar campos del sensor
    SENSOR_TIENE_ESTADO=$(echo "$RESULTADO" | jq '.activos[0].sensores[0] | has("estado")')
    SENSOR_TIENE_IS_ACTIVE=$(echo "$RESULTADO" | jq '.activos[0].sensores[0] | has("is_active")')
    
    if [ "$SENSOR_TIENE_ESTADO" == "true" ] && [ "$SENSOR_TIENE_IS_ACTIVE" == "true" ]; then
        echo -e "${GREEN}✓ Sensores incluyen campos de estado (estado, is_active)${NC}"
    else
        echo -e "${RED}✗ Sensores no incluyen campos de estado completos${NC}"
    fi
else
    echo -e "${RED}✗ Los activos NO incluyen información de sensores${NC}"
fi

# Verificar que algunos sensores están conectados
SENSORES_CONECTADOS=$(echo "$RESULTADO" | jq '[.activos[].sensores[] | select(.estado == "connected")] | length')
if [ "$SENSORES_CONECTADOS" -gt 0 ]; then
    echo -e "${GREEN}✓ Sensores conectados detectados: $SENSORES_CONECTADOS${NC}"
else
    echo -e "${YELLOW}⚠ No se detectaron sensores conectados (es posible que no hayan enviado datos)${NC}"
fi

# Verificar que solo retorna activos del edificio 1
EDIFICIOS_DIFERENTES=$(echo "$RESULTADO" | jq '[.activos[].id_edificio] | unique | length')
if [ "$EDIFICIOS_DIFERENTES" -eq 1 ]; then
    echo -e "${GREEN}✓ Filtrado por edificio correcto (solo edificio $EDIFICIO_ID)${NC}"
else
    echo -e "${RED}✗ Filtrado por edificio incorrecto (múltiples edificios)${NC}"
fi

echo ""
echo "=========================================="
echo -e "${GREEN}Test completado${NC}"
echo "=========================================="
echo ""

# Paso 6: Prueba adicional - Edificio 2
echo -e "${YELLOW}Prueba adicional: Edificio 2${NC}"
echo "GET $BASE_URL/activo/edificio/2"
echo ""

RESULTADO_EDIFICIO_2=$(curl -s -X GET "$BASE_URL/activo/edificio/2")
echo "Respuesta edificio 2:"
echo "$RESULTADO_EDIFICIO_2" | jq '.'
echo ""

TOTAL_EDIFICIO_2=$(echo "$RESULTADO_EDIFICIO_2" | jq '.total')
if [ "$TOTAL_EDIFICIO_2" -ge 1 ]; then
    echo -e "${GREEN}✓ Edificio 2 tiene activos: $TOTAL_EDIFICIO_2${NC}"
else
    echo -e "${YELLOW}⚠ Edificio 2 sin activos${NC}"
fi

echo ""
echo "=========================================="
echo "Resumen de estructura de respuesta:"
echo "=========================================="
echo "$RESULTADO" | jq '{
  edificio_id,
  total,
  primer_activo: .activos[0] | {
    id,
    activo_id,
    estado,
    id_edificio,
    total_sensores,
    ejemplo_sensor: .sensores[0]
  }
}'
