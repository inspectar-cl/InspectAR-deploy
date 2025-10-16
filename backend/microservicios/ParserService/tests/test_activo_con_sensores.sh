#!/bin/bash

# Test para verificar que /activo/:activo_id retorna información completa de sensores

BASE_URL="http://localhost:8090"

echo "=========================================="
echo "Test: Activo Individual con Sensores"
echo "=========================================="
echo ""

# Colores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Paso 1: Crear activo con sensores
echo -e "${YELLOW}Paso 1: Crear activo con sensores${NC}"
echo "Creando activo 201 con 3 sensores..."
RESPONSE=$(curl -s -X POST $BASE_URL/activo \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": 201,
    "estado": "operativo",
    "edificio_id": 5,
    "sensores": [
      {"sensor_id": "TEMP_201", "tipo": "temperatura", "unidad": "°C"},
      {"sensor_id": "PRESS_201", "tipo": "presion", "unidad": "bar"},
      {"sensor_id": "FLOW_201", "tipo": "flujo", "unidad": "L/min"}
    ]
  }')
echo "Respuesta: $RESPONSE"
echo ""

# Esperar para que se cree
sleep 2

# Paso 2: Enviar lecturas para algunos sensores
echo -e "${YELLOW}Paso 2: Enviar lecturas para sensores (solo 2 de 3)${NC}"
echo "Enviando lectura para TEMP_201..."
curl -s -X POST $BASE_URL/lectura \
  -H "Content-Type: application/json" \
  -d '{
    "sensor_id": "TEMP_201",
    "valor": 25.5,
    "timestamp": "2025-10-12T10:00:00Z"
  }' > /dev/null
echo "✓ Lectura enviada"

echo "Enviando lectura para PRESS_201..."
curl -s -X POST $BASE_URL/lectura \
  -H "Content-Type: application/json" \
  -d '{
    "sensor_id": "PRESS_201",
    "valor": 2.5,
    "timestamp": "2025-10-12T10:00:00Z"
  }' > /dev/null
echo "✓ Lectura enviada"
echo "(FLOW_201 sin datos - should be 'never_connected')"
echo ""

# Esperar para que se actualice el estado
sleep 2

# Paso 3: Obtener activo por ID
echo -e "${YELLOW}Paso 3: Obtener activo por ID con información de sensores${NC}"
echo "GET $BASE_URL/activo/201"
echo ""

RESULTADO=$(curl -s -X GET "$BASE_URL/activo/201")

echo "Respuesta completa:"
echo "$RESULTADO" | jq '.'
echo ""

# Paso 4: Validaciones
echo -e "${YELLOW}Paso 4: Validaciones${NC}"

# Verificar que retorna el activo correcto
ACTIVO_ID=$(echo "$RESULTADO" | jq -r '.activo_id')
if [ "$ACTIVO_ID" == "201" ]; then
    echo -e "${GREEN}✓ Activo ID correcto: $ACTIVO_ID${NC}"
else
    echo -e "${RED}✗ Activo ID incorrecto: $ACTIVO_ID (esperado: 201)${NC}"
fi

# Verificar que tiene campo total_sensores
TOTAL_SENSORES=$(echo "$RESULTADO" | jq -r '.total_sensores')
if [ "$TOTAL_SENSORES" == "3" ]; then
    echo -e "${GREEN}✓ Total de sensores correcto: $TOTAL_SENSORES${NC}"
else
    echo -e "${RED}✗ Total de sensores incorrecto: $TOTAL_SENSORES (esperado: 3)${NC}"
fi

# Verificar que tiene array de sensores
TIENE_SENSORES=$(echo "$RESULTADO" | jq 'has("sensores")')
if [ "$TIENE_SENSORES" == "true" ]; then
    echo -e "${GREEN}✓ Campo 'sensores' presente${NC}"
    
    # Verificar cantidad de sensores en el array
    CANT_SENSORES=$(echo "$RESULTADO" | jq '.sensores | length')
    if [ "$CANT_SENSORES" == "3" ]; then
        echo -e "${GREEN}✓ Array de sensores tiene 3 elementos${NC}"
    else
        echo -e "${RED}✗ Array de sensores tiene $CANT_SENSORES elementos (esperado: 3)${NC}"
    fi
else
    echo -e "${RED}✗ Campo 'sensores' no presente${NC}"
fi

# Verificar que los sensores tienen campos de estado
PRIMER_SENSOR=$(echo "$RESULTADO" | jq '.sensores[0]')
TIENE_ESTADO=$(echo "$PRIMER_SENSOR" | jq 'has("estado")')
TIENE_IS_ACTIVE=$(echo "$PRIMER_SENSOR" | jq 'has("is_active")')
TIENE_LAST_SEEN=$(echo "$PRIMER_SENSOR" | jq 'has("last_seen")')
TIENE_TOTAL_REPORTS=$(echo "$PRIMER_SENSOR" | jq 'has("total_reports")')

if [ "$TIENE_ESTADO" == "true" ] && [ "$TIENE_IS_ACTIVE" == "true" ] && \
   [ "$TIENE_LAST_SEEN" == "true" ] && [ "$TIENE_TOTAL_REPORTS" == "true" ]; then
    echo -e "${GREEN}✓ Sensores incluyen campos de estado completos${NC}"
    echo ""
    echo "Ejemplo de sensor con estado:"
    echo "$PRIMER_SENSOR" | jq '{sensor_id, tipo, estado, is_active, total_reports}'
else
    echo -e "${RED}✗ Sensores no incluyen todos los campos de estado${NC}"
fi

# Verificar estados de sensores
echo ""
echo -e "${YELLOW}Paso 5: Verificar estados de sensores${NC}"

SENSORES_CONNECTED=$(echo "$RESULTADO" | jq '[.sensores[] | select(.estado == "connected")] | length')
SENSORES_NEVER_CONNECTED=$(echo "$RESULTADO" | jq '[.sensores[] | select(.estado == "never_connected")] | length')

echo "Sensores conectados: $SENSORES_CONNECTED"
echo "Sensores nunca conectados: $SENSORES_NEVER_CONNECTED"

if [ "$SENSORES_CONNECTED" -ge 1 ]; then
    echo -e "${GREEN}✓ Al menos un sensor conectado detectado${NC}"
else
    echo -e "${YELLOW}⚠ No se detectaron sensores conectados${NC}"
fi

if [ "$SENSORES_NEVER_CONNECTED" -ge 1 ]; then
    echo -e "${GREEN}✓ Sensor sin datos detectado correctamente (never_connected)${NC}"
else
    echo -e "${YELLOW}⚠ No se detectaron sensores sin datos${NC}"
fi

echo ""
echo "=========================================="
echo -e "${GREEN}Test completado${NC}"
echo "=========================================="
echo ""

# Paso 6: Comparación con endpoint de estado
echo -e "${YELLOW}Paso 6: Comparar con endpoint /activo/:activo_id/sensores/estado${NC}"
echo "GET $BASE_URL/activo/201/sensores/estado"
echo ""

RESULTADO_ESTADO=$(curl -s -X GET "$BASE_URL/activo/201/sensores/estado")
echo "Respuesta (solo resumen):"
echo "$RESULTADO_ESTADO" | jq '{activo_id, total_sensores, resumen}'
echo ""

echo "=========================================="
echo "Resumen de campos en /activo/201:"
echo "=========================================="
echo "$RESULTADO" | jq '{
  activo_id,
  estado,
  id_edificio,
  total_sensores,
  sensores: [.sensores[] | {
    sensor_id,
    tipo,
    estado,
    is_active
  }]
}'
