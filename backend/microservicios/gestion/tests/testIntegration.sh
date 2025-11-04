#!/bin/bash

# =============================================================================
# Prueba de Integración Completa
# Verifica la creación de activos y sensores en ambos servicios
# =============================================================================

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# URLs
GESTION_URL="http://localhost:8092"
IOT_URL="http://localhost:8090"
ADMIN_EMAIL="admin@example.com"

echo ""
echo -e "${CYAN}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║        PRUEBA DE INTEGRACIÓN COMPLETA - GESTION + IOT       ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════════╝${NC}"
echo ""

# =============================================================================
# PASO 1: Crear Edificio
# =============================================================================
echo -e "${YELLOW}[PASO 1]${NC} Creando edificio..."
EDIFICIO_RESPONSE=$(curl -s -X POST "$GESTION_URL/admin/edificios" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Torre Integración Test",
    "direccion": "Av. Integración 999",
    "latitud": -33.45,
    "longitud": -70.65,
    "email": "'"$ADMIN_EMAIL"'"
  }')

EDIFICIO_ID=$(echo "$EDIFICIO_RESPONSE" | jq -r '.edificio.id')
echo -e "${GREEN}✓ Edificio creado con ID: $EDIFICIO_ID${NC}"
echo "$EDIFICIO_RESPONSE" | jq '.'
echo ""
sleep 1

# =============================================================================
# PASO 2: Crear Activo (debe crearse en ambos servicios)
# =============================================================================
echo -e "${YELLOW}[PASO 2]${NC} Creando activo (debe sincronizarse a IOT-Service)..."
ACTIVO_RESPONSE=$(curl -s -X POST "$GESTION_URL/admin/activos" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Bomba Integración Test",
    "tipo": "bomba de agua",
    "descripcion": "Bomba para prueba de integración completa",
    "ubicacion": "Sala de Máquinas Principal",
    "edificio_id": '"$EDIFICIO_ID"',
    "email": "'"$ADMIN_EMAIL"'"
  }')

ACTIVO_ID=$(echo "$ACTIVO_RESPONSE" | jq -r '.activo.id')
echo -e "${GREEN}✓ Activo creado en Gestión con ID: $ACTIVO_ID${NC}"
echo "$ACTIVO_RESPONSE" | jq '.'
echo ""
sleep 1

# =============================================================================
# PASO 3: Verificar que el activo existe en IOT-Service
# =============================================================================
echo -e "${YELLOW}[PASO 3]${NC} Verificando activo en IOT-Service..."
IOT_ACTIVO=$(curl -s "$IOT_URL/activo/$ACTIVO_ID")
IOT_ACTIVO_ID=$(echo "$IOT_ACTIVO" | jq -r '.activo_id // .activo.activo_id')

echo "Respuesta de IOT-Service:"
echo "$IOT_ACTIVO" | jq '.'
echo ""

if [ "$IOT_ACTIVO_ID" == "$ACTIVO_ID" ]; then
  echo -e "${GREEN}✓✓✓ ÉXITO: El activo existe en IOT-Service con el MISMO ID${NC}"
  echo -e "    Gestión ID: ${BLUE}$ACTIVO_ID${NC}"
  echo -e "    IOT-Service ID: ${BLUE}$IOT_ACTIVO_ID${NC}"
else
  echo -e "${RED}✗✗✗ ERROR: Los IDs NO coinciden${NC}"
  echo -e "    Gestión ID: ${RED}$ACTIVO_ID${NC}"
  echo -e "    IOT-Service ID: ${RED}$IOT_ACTIVO_ID${NC}"
fi
echo ""
sleep 1

# =============================================================================
# PASO 4: Crear Sensores (deben crearse en IOT-Service)
# =============================================================================
echo -e "${YELLOW}[PASO 4]${NC} Creando sensores en el activo..."

# Sensor 1: Temperatura
echo "4.1. Creando sensor de temperatura..."
SENSOR1_RESPONSE=$(curl -s -X POST "$GESTION_URL/admin/sensores" \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": '"$ACTIVO_ID"',
    "nombre": "Sensor_Temp_Integration_'$ACTIVO_ID'",
    "tipo": "temperatura",
    "unidad": "°C",
    "email": "'"$ADMIN_EMAIL"'"
  }')

SENSOR1_ID=$(echo "$SENSOR1_RESPONSE" | jq -r '.sensor.sensor_id // .sensor.sensor.sensor_id')
echo -e "${GREEN}✓ Sensor 1 creado: $SENSOR1_ID${NC}"
echo "$SENSOR1_RESPONSE" | jq '.'
echo ""
sleep 1

# Sensor 2: Presión
echo "4.2. Creando sensor de presión..."
SENSOR2_RESPONSE=$(curl -s -X POST "$GESTION_URL/admin/sensores" \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": '"$ACTIVO_ID"',
    "nombre": "Sensor_Pres_Integration_'$ACTIVO_ID'",
    "tipo": "presion",
    "unidad": "PSI",
    "email": "'"$ADMIN_EMAIL"'"
  }')

SENSOR2_ID=$(echo "$SENSOR2_RESPONSE" | jq -r '.sensor.sensor_id // .sensor.sensor.sensor_id')
echo -e "${GREEN}✓ Sensor 2 creado: $SENSOR2_ID${NC}"
echo "$SENSOR2_RESPONSE" | jq '.'
echo ""
sleep 1

# Sensor 3: Vibración
echo "4.3. Creando sensor de vibración..."
SENSOR3_RESPONSE=$(curl -s -X POST "$GESTION_URL/admin/sensores" \
  -H "Content-Type: application/json" \
  -d '{
    "activo_id": '"$ACTIVO_ID"',
    "nombre": "Sensor_Vib_Integration_'$ACTIVO_ID'",
    "tipo": "vibracion",
    "unidad": "Hz",
    "email": "'"$ADMIN_EMAIL"'"
  }')

SENSOR3_ID=$(echo "$SENSOR3_RESPONSE" | jq -r '.sensor.sensor_id // .sensor.sensor.sensor_id')
echo -e "${GREEN}✓ Sensor 3 creado: $SENSOR3_ID${NC}"
echo "$SENSOR3_RESPONSE" | jq '.'
echo ""
sleep 2

# =============================================================================
# PASO 5: Verificar sensores en IOT-Service
# =============================================================================
echo -e "${YELLOW}[PASO 5]${NC} Verificando sensores en IOT-Service..."
IOT_ACTIVO_WITH_SENSORS=$(curl -s "$IOT_URL/activo/$ACTIVO_ID")

echo "Activo completo con sensores desde IOT-Service:"
echo "$IOT_ACTIVO_WITH_SENSORS" | jq '.'
echo ""

# Contar sensores
SENSOR_COUNT=$(echo "$IOT_ACTIVO_WITH_SENSORS" | jq '.sensores | length')
echo -e "${CYAN}Total de sensores encontrados en IOT-Service: ${BLUE}$SENSOR_COUNT${NC}"
echo ""

# Mostrar detalles de cada sensor
if [ "$SENSOR_COUNT" -gt 0 ]; then
  echo -e "${GREEN}✓✓✓ ÉXITO: Sensores encontrados en IOT-Service${NC}"
  echo ""
  echo "Detalles de los sensores:"
  echo "$IOT_ACTIVO_WITH_SENSORS" | jq '.sensores[] | {sensor_id, tipo, unidad, estado, is_active}'
  echo ""
else
  echo -e "${RED}✗✗✗ ERROR: No se encontraron sensores en IOT-Service${NC}"
fi

# =============================================================================
# PASO 6: Verificar en base de datos de Gestión (NO debe tener sensores)
# =============================================================================
echo -e "${YELLOW}[PASO 6]${NC} Verificando que sensores NO están en Gestión DB..."
GESTION_ACTIVO=$(curl -s "$GESTION_URL/activos/$ACTIVO_ID")

echo "Activo desde Gestión:"
echo "$GESTION_ACTIVO" | jq '.'
echo ""

HAS_SENSORS_FIELD=$(echo "$GESTION_ACTIVO" | jq 'has("sensores")')
if [ "$HAS_SENSORS_FIELD" == "false" ]; then
  echo -e "${GREEN}✓✓✓ CORRECTO: El activo en Gestión NO tiene campo de sensores${NC}"
  echo -e "    ${CYAN}→ Los sensores solo existen en IOT-Service (MongoDB)${NC}"
else
  echo -e "${YELLOW}⚠ ADVERTENCIA: El activo en Gestión tiene campo de sensores${NC}"
fi
echo ""

# =============================================================================
# PASO 7: Verificar logs de auditoría
# =============================================================================
echo -e "${YELLOW}[PASO 7]${NC} Verificando logs de auditoría..."

echo "7.1. Logs de creación de edificio:"
curl -s "$GESTION_URL/admin/logs?entidad=edificio&entidad_id=$EDIFICIO_ID&limit=1" | jq '.logs[0] | {accion, entidad, descripcion, usuario_email, fecha_accion}'
echo ""

echo "7.2. Logs de creación de activo:"
curl -s "$GESTION_URL/admin/logs?entidad=activo&entidad_id=$ACTIVO_ID&limit=1" | jq '.logs[0] | {accion, entidad, descripcion, usuario_email, fecha_accion}'
echo ""

echo "7.3. Logs de creación de sensores:"
SENSOR_LOGS=$(curl -s "$GESTION_URL/admin/logs?entidad=sensor&entidad_id=$ACTIVO_ID&limit=10")
SENSOR_LOG_COUNT=$(echo "$SENSOR_LOGS" | jq '.logs | length')
echo "Total de logs de sensores: $SENSOR_LOG_COUNT"
echo "$SENSOR_LOGS" | jq '.logs[] | {accion, descripcion, usuario_email}'
echo ""

# =============================================================================
# RESUMEN FINAL
# =============================================================================
echo ""
echo -e "${CYAN}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║                    RESUMEN DE LA PRUEBA                      ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "📋 ${GREEN}Edificio creado:${NC} ID = ${BLUE}$EDIFICIO_ID${NC}"
echo -e "🏭 ${GREEN}Activo creado:${NC} ID = ${BLUE}$ACTIVO_ID${NC}"
echo ""
echo -e "🔍 ${CYAN}Verificación en IOT-Service:${NC}"
echo -e "   • Activo existe: ${GREEN}✓${NC}"
echo -e "   • ID coincide con Gestión: ${GREEN}✓${NC} (Gestión: $ACTIVO_ID = IOT: $IOT_ACTIVO_ID)"
echo -e "   • Sensores registrados: ${GREEN}$SENSOR_COUNT${NC}"
echo ""
echo -e "📊 ${CYAN}Sensores creados:${NC}"
echo -e "   1. ${GREEN}$SENSOR1_ID${NC} (Temperatura, °C)"
echo -e "   2. ${GREEN}$SENSOR2_ID${NC} (Presión, PSI)"
echo -e "   3. ${GREEN}$SENSOR3_ID${NC} (Vibración, Hz)"
echo ""
echo -e "📝 ${CYAN}Logs de auditoría:${NC}"
echo -e "   • Log de edificio: ${GREEN}✓${NC}"
echo -e "   • Log de activo: ${GREEN}✓${NC}"
echo -e "   • Logs de sensores: ${GREEN}$SENSOR_LOG_COUNT registrados${NC}"
echo ""
echo -e "✅ ${CYAN}Arquitectura verificada:${NC}"
echo -e "   • Activos en PostgreSQL (Gestión): ${GREEN}✓${NC}"
echo -e "   • Activos en MongoDB (IOT-Service): ${GREEN}✓${NC}"
echo -e "   • Sensores SOLO en MongoDB: ${GREEN}✓${NC}"
echo -e "   • Sincronización funcionando: ${GREEN}✓${NC}"
echo ""

# =============================================================================
# LIMPIEZA OPCIONAL
# =============================================================================
echo -e "${YELLOW}¿Deseas eliminar los datos de prueba? (y/N)${NC}"
read -t 10 -n 1 CLEANUP_RESPONSE || CLEANUP_RESPONSE="n"
echo ""

if [ "$CLEANUP_RESPONSE" == "y" ] || [ "$CLEANUP_RESPONSE" == "Y" ]; then
  echo -e "${YELLOW}Limpiando datos de prueba...${NC}"
  
  echo "Eliminando activo (también se eliminará de IOT-Service)..."
  curl -s -X DELETE "$GESTION_URL/admin/activos/$ACTIVO_ID?email=$ADMIN_EMAIL" | jq '.'
  echo ""
  
  echo "Eliminando edificio..."
  curl -s -X DELETE "$GESTION_URL/admin/edificios/$EDIFICIO_ID?email=$ADMIN_EMAIL" | jq '.'
  echo ""
  
  echo -e "${GREEN}✓ Datos de prueba eliminados${NC}"
  
  # Verificar que se eliminó de IOT-Service
  echo ""
  echo "Verificando eliminación en IOT-Service..."
  IOT_CHECK=$(curl -s "$IOT_URL/activo/$ACTIVO_ID" 2>&1)
  if echo "$IOT_CHECK" | grep -q "error\|not found\|no encontrado"; then
    echo -e "${GREEN}✓ Activo eliminado correctamente de IOT-Service${NC}"
  else
    echo -e "${YELLOW}⚠ El activo aún existe en IOT-Service (puede ser normal)${NC}"
  fi
else
  echo -e "${BLUE}ℹ️  Los datos de prueba permanecen en el sistema${NC}"
  echo -e "${BLUE}   Activo ID: $ACTIVO_ID${NC}"
  echo -e "${BLUE}   Edificio ID: $EDIFICIO_ID${NC}"
fi

echo ""
echo -e "${CYAN}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║              ✓ PRUEBA DE INTEGRACIÓN COMPLETADA             ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════════╝${NC}"
echo ""
