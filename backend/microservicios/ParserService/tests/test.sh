#!/bin/bash

# 🧪 Test runner compacto para ParserService (activo_id entero)
set -euo pipefail

API_URL="http://localhost:8090"
VERBOSE=${VERBOSE:-1}
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

req() {
  local method="$1"; shift
  local url="$1"; shift
  local data="${1:-}"
  local response
  LAST_METHOD="$method"; LAST_URL="$url"; LAST_DATA="$data"
  if [[ -n "$data" ]]; then
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X "$method" -H 'Content-Type: application/json' -d "$data" "$url")
  else
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X "$method" "$url")
  fi
  CODE=$(printf '%s' "$response" | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
  BODY=$(printf '%s' "$response" | sed -E 's/HTTPSTATUS:[0-9]{3}$//')
  if [[ "$VERBOSE" == "1" ]]; then
    echo ">> $method $url"
    if [[ -n "$data" ]]; then echo "payload: $data"; fi
    echo "<< status: $CODE"
    if [[ -n "$BODY" ]]; then
      if (( ${#BODY} > 600 )); then
        echo "body: ${BODY:0:600}..."
      else
        echo "body: $BODY"
      fi
    fi
  fi
}

assert_2xx() {
  local code="$1"; shift
  if [[ ! "$code" =~ ^2 ]]; then
    echo -e "${RED}❌ HTTP $code${NC}"; exit 1
  fi
}

echo -e "${YELLOW}🔎 Healthcheck${NC}"
req GET "$API_URL/healthz"
assert_2xx "$CODE"; echo -e "${GREEN}OK${NC}"

TS=$(date +%s)
ACTIVO_ID=$((200000 + (RANDOM % 700000)))
SENS_A="TEMP_$TS"
SENS_B="PRES_$TS"

echo -e "${YELLOW}➕ Crear Activo${NC}"
payload_create=$(cat <<JSON
{
  "activo_id": $ACTIVO_ID,
  "nombre": "Equipo $TS",
  "estado": "operativo",
  "id_edificio": "ED-01",
  "sensores": [
    {"sensor_id": "$SENS_A", "tipo": "temperatura", "unidad": "°C"}
  ]
}
JSON
)
req POST "$API_URL/activo" "$payload_create"
assert_2xx "$CODE"; echo "$BODY" | grep -q 'activo_id'

echo -e "${YELLOW}📚 Listar Activos${NC}"
req GET "$API_URL/activo"
assert_2xx "$CODE"

echo -e "${YELLOW}🔎 Obtener Activo${NC}"
req GET "$API_URL/activo/$ACTIVO_ID"
assert_2xx "$CODE"

echo -e "${YELLOW}➕ Agregar Sensor${NC}"
payload_sensor=$(cat <<JSON
{"sensor_id":"$SENS_B","tipo":"presion","unidad":"bar"}
JSON
)
req POST "$API_URL/activo/$ACTIVO_ID/sensores" "$payload_sensor"
assert_2xx "$CODE"

echo -e "${YELLOW}📡 Enviar Lecturas${NC}"
now=$(date -u +%Y-%m-%dT%H:%M:%SZ)
req POST "$API_URL/lectura" "{\"sensor_id\":\"$SENS_A\",\"valor\":22.5,\"timestamp\":\"$now\"}"
assert_2xx "$CODE"
req POST "$API_URL/lectura" "{\"sensor_id\":\"$SENS_B\",\"valor\":1.2,\"timestamp\":\"$now\"}"
assert_2xx "$CODE"

echo -e "${YELLOW}📈 Datos por Activo${NC}"
req GET "$API_URL/lectura/$ACTIVO_ID/datos"
assert_2xx "$CODE"

echo -e "${YELLOW}🕒 Últimos Datos por Activo${NC}"
req GET "$API_URL/lectura/$ACTIVO_ID/datos/ultimo"
assert_2xx "$CODE"

echo -e "${YELLOW}⚙️  Cambiar Estado del Activo${NC}"
req PUT "$API_URL/activo/$ACTIVO_ID/estado" '{"estado":"mantenimiento"}'
assert_2xx "$CODE"

echo -e "${YELLOW}🧭 Estado de Sensores del Activo${NC}"
req GET "$API_URL/activo/$ACTIVO_ID/sensores/estado"
assert_2xx "$CODE"

echo -e "${YELLOW}📊 Stats Monitoreo${NC}"
req GET "$API_URL/api/sensors/stats"
assert_2xx "$CODE"

echo -e "${YELLOW}🩺 Health Monitoreo${NC}"
req GET "$API_URL/api/sensors/health"
assert_2xx "$CODE"

echo -e "${GREEN}✅ Todos los endpoints probados con éxito (activo_id=$ACTIVO_ID)${NC}"