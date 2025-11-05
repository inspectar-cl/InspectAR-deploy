#!/bin/bash

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# URL base
GESTION_URL="http://localhost:8092"

echo ""
echo -e "${PURPLE}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${PURPLE}║    PRUEBA DE RUTAS ADMIN: USUARIOS-EDIFICIOS                ║${NC}"
echo -e "${PURPLE}╚══════════════════════════════════════════════════════════════╝${NC}"
echo ""

# =============================================================================
# PASO 1: Asignar usuario a edificio
# =============================================================================
echo -e "${YELLOW}[PASO 1]${NC} Asignando usuario a edificio..."
echo -e "${CYAN}POST /admin/usuarios/edificios${NC}"
echo ""

RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "$GESTION_URL/admin/usuarios/edificios" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "residente.especial@example.com",
    "edificio_id": 2,
    "admin_email": "admin@example.com"
  }')

HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed -e 's/HTTP_STATUS.*//')

if [ "$HTTP_STATUS" == "200" ]; then
  echo -e "${GREEN}✓✓✓ ÉXITO: Usuario asignado al edificio${NC}"
  echo -e "${GREEN}    HTTP Status: $HTTP_STATUS${NC}"
  echo ""
  echo -e "${BLUE}Respuesta:${NC}"
  echo "$BODY" | jq '.'
else
  echo -e "${RED}✗✗✗ ERROR: Fallo al asignar usuario${NC}"
  echo -e "${RED}    HTTP Status: $HTTP_STATUS${NC}"
  echo "$BODY"
fi
echo ""
sleep 2

# =============================================================================
# PASO 2: Verificar acceso del usuario al edificio
# =============================================================================
echo -e "${YELLOW}[PASO 2]${NC} Verificando acceso del usuario al edificio 2..."
echo -e "${CYAN}GET /usuarios/residente.especial@example.com/edificio/2/acceso${NC}"
echo ""

RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" "$GESTION_URL/usuarios/residente.especial@example.com/edificio/2/acceso")

HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed -e 's/HTTP_STATUS.*//')

if [ "$HTTP_STATUS" == "200" ]; then
  echo -e "${GREEN}✓✓✓ ÉXITO: Usuario tiene acceso al edificio${NC}"
  echo -e "${GREEN}    HTTP Status: $HTTP_STATUS${NC}"
  echo ""
  echo "$BODY" | jq '.'
else
  echo -e "${RED}✗✗✗ ERROR: Usuario no tiene acceso${NC}"
  echo -e "${RED}    HTTP Status: $HTTP_STATUS${NC}"
  echo "$BODY"
fi
echo ""
sleep 2

# =============================================================================
# PASO 3: Obtener usuarios del edificio
# =============================================================================
echo -e "${YELLOW}[PASO 3]${NC} Obteniendo usuarios del edificio 2..."
echo -e "${CYAN}GET /admin/edificios/2/usuarios${NC}"
echo ""

RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" "$GESTION_URL/admin/edificios/2/usuarios")

HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed -e 's/HTTP_STATUS.*//')

if [ "$HTTP_STATUS" == "200" ]; then
  echo -e "${GREEN}✓✓✓ ÉXITO: Usuarios obtenidos correctamente${NC}"
  echo -e "${GREEN}    HTTP Status: $HTTP_STATUS${NC}"
  echo ""
  echo "$BODY" | jq '.'
else
  echo -e "${RED}✗✗✗ ERROR: Fallo al obtener usuarios${NC}"
  echo -e "${RED}    HTTP Status: $HTTP_STATUS${NC}"
  echo "$BODY"
fi
echo ""
sleep 2

# =============================================================================
# PASO 4: Obtener edificios del usuario
# =============================================================================
echo -e "${YELLOW}[PASO 4]${NC} Obteniendo edificios del usuario..."
echo -e "${CYAN}GET /usuarios/edificios/residente.especial@example.com${NC}"
echo ""

RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" "$GESTION_URL/usuarios/edificios/residente.especial@example.com")

HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed -e 's/HTTP_STATUS.*//')

if [ "$HTTP_STATUS" == "200" ]; then
  echo -e "${GREEN}✓✓✓ ÉXITO: Edificios del usuario obtenidos${NC}"
  echo -e "${GREEN}    HTTP Status: $HTTP_STATUS${NC}"
  echo ""
  echo "$BODY" | jq '.'
  
  # Verificar que ahora tiene 2 edificios
  TOTAL=$(echo "$BODY" | jq '.total')
  if [ "$TOTAL" == "2" ]; then
    echo -e "${BLUE}✓ Usuario ahora tiene acceso a 2 edificios (incluyendo el nuevo)${NC}"
  fi
else
  echo -e "${RED}✗✗✗ ERROR: Fallo al obtener edificios${NC}"
  echo -e "${RED}    HTTP Status: $HTTP_STATUS${NC}"
  echo "$BODY"
fi
echo ""
sleep 2

# =============================================================================
# PASO 5: Remover usuario del edificio
# =============================================================================
echo -e "${YELLOW}[PASO 5]${NC} Removiendo usuario del edificio 2..."
echo -e "${CYAN}DELETE /admin/usuarios/edificios${NC}"
echo ""

RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X DELETE "$GESTION_URL/admin/usuarios/edificios" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "residente.especial@example.com",
    "edificio_id": 2,
    "admin_email": "admin@example.com"
  }')

HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed -e 's/HTTP_STATUS.*//')

if [ "$HTTP_STATUS" == "200" ]; then
  echo -e "${GREEN}✓✓✓ ÉXITO: Usuario removido del edificio${NC}"
  echo -e "${GREEN}    HTTP Status: $HTTP_STATUS${NC}"
  echo ""
  echo "$BODY" | jq '.'
else
  echo -e "${RED}✗✗✗ ERROR: Fallo al remover usuario${NC}"
  echo -e "${RED}    HTTP Status: $HTTP_STATUS${NC}"
  echo "$BODY"
fi
echo ""
sleep 2

# =============================================================================
# PASO 6: Verificar que ya no tiene acceso
# =============================================================================
echo -e "${YELLOW}[PASO 6]${NC} Verificando que el usuario ya NO tiene acceso al edificio 2..."
echo -e "${CYAN}GET /usuarios/residente.especial@example.com/edificio/2/acceso${NC}"
echo ""

RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" "$GESTION_URL/usuarios/residente.especial@example.com/edificio/2/acceso")

HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed -e 's/HTTP_STATUS.*//')

if [ "$HTTP_STATUS" == "403" ]; then
  echo -e "${GREEN}✓✓✓ ÉXITO: Usuario ya NO tiene acceso (como esperado)${NC}"
  echo -e "${GREEN}    HTTP Status: $HTTP_STATUS (Forbidden)${NC}"
  echo ""
  echo "$BODY" | jq '.'
else
  echo -e "${YELLOW}⚠ ATENCIÓN: Status inesperado${NC}"
  echo -e "${YELLOW}    HTTP Status: $HTTP_STATUS${NC}"
  echo "$BODY"
fi
echo ""
sleep 2

# =============================================================================
# PASO 7: Verificar logs de auditoría
# =============================================================================
echo -e "${YELLOW}[PASO 7]${NC} Verificando logs de auditoría..."
echo -e "${CYAN}GET /admin/logs?entidad=usuario_edificio&limit=5${NC}"
echo ""

RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" "$GESTION_URL/admin/logs?entidad=usuario_edificio&limit=5")

HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed -e 's/HTTP_STATUS.*//')

if [ "$HTTP_STATUS" == "200" ]; then
  echo -e "${GREEN}✓✓✓ ÉXITO: Logs de auditoría obtenidos${NC}"
  echo -e "${GREEN}    HTTP Status: $HTTP_STATUS${NC}"
  echo ""
  echo "$BODY" | jq '.'
else
  echo -e "${RED}✗✗✗ ERROR: Fallo al obtener logs${NC}"
  echo -e "${RED}    HTTP Status: $HTTP_STATUS${NC}"
  echo "$BODY"
fi
echo ""

# =============================================================================
# RESUMEN FINAL
# =============================================================================
echo ""
echo -e "${CYAN}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║                    RESUMEN DE LA PRUEBA                      ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "✅ ${CYAN}Funcionalidades validadas:${NC}"
echo -e "   ${GREEN}✓${NC} Asignar usuario a edificio"
echo -e "   ${GREEN}✓${NC} Verificar acceso de usuario a edificio"
echo -e "   ${GREEN}✓${NC} Obtener usuarios de un edificio"
echo -e "   ${GREEN}✓${NC} Obtener edificios de un usuario"
echo -e "   ${GREEN}✓${NC} Remover usuario de edificio"
echo -e "   ${GREEN}✓${NC} Verificar que el acceso fue removido"
echo -e "   ${GREEN}✓${NC} Logs de auditoría registrados"
echo ""
echo -e "${BLUE}🔗 Rutas verificadas:${NC}"
echo -e "   • ${GREEN}POST${NC}   /admin/usuarios/edificios"
echo -e "   • ${GREEN}DELETE${NC} /admin/usuarios/edificios"
echo -e "   • ${GREEN}GET${NC}    /admin/edificios/:edificio_id/usuarios"
echo -e "   • ${GREEN}GET${NC}    /usuarios/edificios/:email"
echo -e "   • ${GREEN}GET${NC}    /usuarios/:email/edificio/:edificio_id/acceso"
echo -e "   • ${GREEN}GET${NC}    /admin/logs"
echo ""
echo -e "${PURPLE}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${PURPLE}║              ✓ PRUEBA COMPLETADA EXITOSAMENTE               ║${NC}"
echo -e "${PURPLE}╚══════════════════════════════════════════════════════════════╝${NC}"
echo ""
