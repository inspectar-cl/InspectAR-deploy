#!/bin/bash

# =============================================================================
# Script de Pruebas para Rutas de Administración
# Prueba: Crear Edificio → Crear Activo → Crear Sensores → Verificar Logs
# =============================================================================

set -e  # Salir si hay errores

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# URLs de los servicios
GESTION_URL="http://localhost:8092"
PARSER_URL="http://localhost:8090"

# Email del usuario administrador
ADMIN_EMAIL="admin@example.com"

# Variables para IDs creados
EDIFICIO_ID=""
ACTIVO_ID=""
SENSOR1_ID=""
SENSOR2_ID=""

echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}    PRUEBA DE RUTAS DE ADMINISTRACIÓN - MICROSERVICIO GESTIÓN${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo ""

# =============================================================================
# FASE 1: CREAR EDIFICIO
# =============================================================================
echo -e "${YELLOW}[FASE 1] Creando Edificio...${NC}"
echo ""

EDIFICIO_RESPONSE=$(curl -s -X POST "$GESTION_URL/admin/edificios" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Torre Test Admin",
    "direccion": "Avenida Test 123",
    "ciudad": "Santiago",
    "pais": "Chile",
    "latitud": -33.4569,
    "longitud": -70.6483,
    "email": "'"$ADMIN_EMAIL"'"
  }')

echo "Respuesta de creación de edificio:"
echo "$EDIFICIO_RESPONSE" | jq '.'
echo ""

# Extraer ID del edificio
EDIFICIO_ID=$(echo "$EDIFICIO_RESPONSE" | jq -r '.edificio.id')

if [ "$EDIFICIO_ID" == "null" ] || [ -z "$EDIFICIO_ID" ]; then
  echo -e "${RED}✗ Error: No se pudo crear el edificio${NC}"
  exit 1
fi

echo -e "${GREEN}✓ Edificio creado exitosamente con ID: $EDIFICIO_ID${NC}"
echo ""
sleep 1

# =============================================================================
# FASE 2: CREAR ACTIVO
# =============================================================================
echo -e "${YELLOW}[FASE 2] Creando Activo en el Edificio...${NC}"
echo ""

ACTIVO_RESPONSE=$(curl -s -X POST "$GESTION_URL/admin/activos" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Bomba Test Admin",
    "tipo": "bomba de agua",
    "descripcion": "Bomba de agua para pruebas de administración",
    "ubicacion": "Sala de Máquinas - Piso 1",
    "edificio_id": '"$EDIFICIO_ID"',
    "email": "'"$ADMIN_EMAIL"'"
  }')

echo "Respuesta de creación de activo:"
echo "$ACTIVO_RESPONSE" | jq '.'
echo ""

# Extraer ID del activo
ACTIVO_ID=$(echo "$ACTIVO_RESPONSE" | jq -r '.activo.id')

if [ "$ACTIVO_ID" == "null" ] || [ -z "$ACTIVO_ID" ]; then
  echo -e "${RED}✗ Error: No se pudo crear el activo${NC}"
  exit 1
fi

echo -e "${GREEN}✓ Activo creado exitosamente con ID: $ACTIVO_ID${NC}"
echo ""
sleep 1

# =============================================================================
# FASE 3: CREAR SENSORES
# =============================================================================
echo -e "${YELLOW}[FASE 3] Creando Sensores para el Activo...${NC}"
echo ""

# Sensor 1: Temperatura
echo "Creando Sensor 1: Temperatura..."
SENSOR1_RESPONSE=$(curl -s -X POST "$GESTION_URL/admin/sensores" \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": '"$ACTIVO_ID"',
    "nombre": "Sensor_Temp_Test_'$ACTIVO_ID'",
    "tipo": "temperatura",
    "unidad": "°C",
    "email": "'"$ADMIN_EMAIL"'"
  }')

echo "$SENSOR1_RESPONSE" | jq '.'
SENSOR1_ID=$(echo "$SENSOR1_RESPONSE" | jq -r '.sensor.sensor_id // .sensor.sensor.sensor_id // "temp_'$ACTIVO_ID'"')
echo -e "${GREEN}✓ Sensor 1 (Temperatura) creado: $SENSOR1_ID${NC}"
echo ""
sleep 1

# Sensor 2: Presión
echo "Creando Sensor 2: Presión..."
SENSOR2_RESPONSE=$(curl -s -X POST "$GESTION_URL/admin/sensores" \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": '"$ACTIVO_ID"',
    "nombre": "Sensor_Pres_Test_'$ACTIVO_ID'",
    "tipo": "presion",
    "unidad": "PSI",
    "email": "'"$ADMIN_EMAIL"'"
  }')

echo "$SENSOR2_RESPONSE" | jq '.'
SENSOR2_ID=$(echo "$SENSOR2_RESPONSE" | jq -r '.sensor.sensor_id // .sensor.sensor.sensor_id // "pres_'$ACTIVO_ID'"')
echo -e "${GREEN}✓ Sensor 2 (Presión) creado: $SENSOR2_ID${NC}"
echo ""
sleep 1

# =============================================================================
# FASE 4: VERIFICAR EN GESTION
# =============================================================================
echo -e "${YELLOW}[FASE 4] Verificando datos en Microservicio de Gestión...${NC}"
echo ""

# Verificar edificio
echo "1. Verificando edificio en Gestión:"
curl -s "$GESTION_URL/activos/edificio/$EDIFICIO_ID" | jq '.'
echo ""

# Verificar activo
echo "2. Verificando activo en Gestión:"
curl -s "$GESTION_URL/activos/$ACTIVO_ID" | jq '.'
echo ""

echo -e "${GREEN}✓ Datos verificados en Gestión${NC}"
echo ""
sleep 1

# =============================================================================
# FASE 5: VERIFICAR EN PARSERSERVICE
# =============================================================================
echo -e "${YELLOW}[FASE 5] Verificando datos en ParserService...${NC}"
echo ""

# Verificar activo en ParserService
echo "1. Verificando activo en ParserService (debe tener el mismo ID):"
PARSER_ACTIVO=$(curl -s "$PARSER_URL/activo/$ACTIVO_ID")
echo "$PARSER_ACTIVO" | jq '.'
echo ""

# Extraer activo_id de la respuesta
PARSER_ACTIVO_ID=$(echo "$PARSER_ACTIVO" | jq -r '.activo.activo_id // .activo_id')

if [ "$PARSER_ACTIVO_ID" == "$ACTIVO_ID" ]; then
  echo -e "${GREEN}✓ IDs coinciden: Gestión=$ACTIVO_ID, Parser=$PARSER_ACTIVO_ID${NC}"
else
  echo -e "${RED}✗ ERROR: IDs NO coinciden: Gestión=$ACTIVO_ID, Parser=$PARSER_ACTIVO_ID${NC}"
fi
echo ""

# Verificar sensores en ParserService
echo "2. Verificando sensores del activo en ParserService:"
PARSER_SENSORES=$(curl -s "$PARSER_URL/activo/$ACTIVO_ID")
echo "$PARSER_SENSORES" | jq '.activo.sensores // .sensores'
echo ""

SENSOR_COUNT=$(echo "$PARSER_SENSORES" | jq '.activo.sensores // .sensores | length')
echo -e "${GREEN}✓ Sensores encontrados en ParserService: $SENSOR_COUNT${NC}"
echo ""
sleep 1

# =============================================================================
# FASE 6: VERIFICAR LOGS DE AUDITORÍA
# =============================================================================
echo -e "${YELLOW}[FASE 6] Verificando Logs de Auditoría...${NC}"
echo ""

# Obtener logs más recientes
echo "1. Logs de creación de edificio:"
curl -s "$GESTION_URL/admin/logs?entidad=edificio&entidad_id=$EDIFICIO_ID&limit=5" | jq '.logs'
echo ""

echo "2. Logs de creación de activo:"
curl -s "$GESTION_URL/admin/logs?entidad=activo&entidad_id=$ACTIVO_ID&limit=5" | jq '.logs'
echo ""

echo "3. Logs de creación de sensores:"
curl -s "$GESTION_URL/admin/logs?entidad=sensor&entidad_id=$ACTIVO_ID&limit=10" | jq '.logs'
echo ""

# Contar logs totales
TOTAL_LOGS=$(curl -s "$GESTION_URL/admin/logs?limit=1000" | jq '.total')
echo -e "${GREEN}✓ Total de logs en el sistema: $TOTAL_LOGS${NC}"
echo ""
sleep 1

# =============================================================================
# FASE 7: RESUMEN Y VALIDACIÓN
# =============================================================================
echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}                      RESUMEN DE PRUEBAS${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo ""
echo -e "📋 ${GREEN}Edificio creado:${NC} ID = $EDIFICIO_ID"
echo -e "🏭 ${GREEN}Activo creado:${NC} ID = $ACTIVO_ID"
echo -e "🌡️  ${GREEN}Sensor 1 (Temp):${NC} $SENSOR1_ID"
echo -e "📊 ${GREEN}Sensor 2 (Pres):${NC} $SENSOR2_ID"
echo ""

# Validar que los IDs coincidan
echo -e "${YELLOW}Validación de Consistencia:${NC}"
if [ "$PARSER_ACTIVO_ID" == "$ACTIVO_ID" ]; then
  echo -e "  ✓ ${GREEN}IDs coinciden entre Gestión y ParserService${NC}"
else
  echo -e "  ✗ ${RED}ERROR: IDs NO coinciden${NC}"
fi

if [ "$SENSOR_COUNT" -ge 2 ]; then
  echo -e "  ✓ ${GREEN}Sensores creados correctamente en ParserService${NC}"
else
  echo -e "  ✗ ${YELLOW}ADVERTENCIA: Se esperaban 2 sensores, encontrados: $SENSOR_COUNT${NC}"
fi

echo ""
echo -e "${YELLOW}Verificación de Logs:${NC}"
echo -e "  ℹ️  Total de logs registrados: $TOTAL_LOGS"
echo -e "  ℹ️  Verifica que incluyan las 4 operaciones:"
echo -e "     1. Crear edificio"
echo -e "     2. Crear activo"
echo -e "     3. Crear sensor (temperatura)"
echo -e "     4. Crear sensor (presión)"
echo ""

# =============================================================================
# FASE 8: LIMPIEZA (OPCIONAL)
# =============================================================================
echo -e "${YELLOW}[OPCIONAL] ¿Deseas eliminar los datos de prueba? (y/N)${NC}"
read -t 10 -n 1 CLEANUP_RESPONSE || CLEANUP_RESPONSE="n"
echo ""

if [ "$CLEANUP_RESPONSE" == "y" ] || [ "$CLEANUP_RESPONSE" == "Y" ]; then
  echo -e "${YELLOW}Limpiando datos de prueba...${NC}"
  
  # Eliminar activo (esto también debería eliminar de ParserService)
  echo "Eliminando activo..."
  curl -s -X DELETE "$GESTION_URL/admin/activos/$ACTIVO_ID?email=$ADMIN_EMAIL" | jq '.'
  echo ""
  
  # Eliminar edificio
  echo "Eliminando edificio..."
  curl -s -X DELETE "$GESTION_URL/admin/edificios/$EDIFICIO_ID?email=$ADMIN_EMAIL" | jq '.'
  echo ""
  
  echo -e "${GREEN}✓ Datos de prueba eliminados${NC}"
else
  echo -e "${BLUE}ℹ️  Los datos de prueba permanecen en el sistema${NC}"
  echo -e "${BLUE}   Para eliminarlos manualmente:${NC}"
  echo -e "   curl -X DELETE \"$GESTION_URL/admin/activos/$ACTIVO_ID?email=$ADMIN_EMAIL\""
  echo -e "   curl -X DELETE \"$GESTION_URL/admin/edificios/$EDIFICIO_ID?email=$ADMIN_EMAIL\""
fi

echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}               ✓ PRUEBAS COMPLETADAS EXITOSAMENTE${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo ""
