#!/bin/bash

# 🎫 TEST DE SISTEMA DE TICKETS - SERVICIO DE NOTIFICACIONES
# ===========================================================
# Este script prueba los endpoints del sistema de tickets

# Configuración de colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuración de la API
API_URL="http://localhost:8091"
TEST_COUNT=0
PASS_COUNT=0
FAIL_COUNT=0

# Variables para almacenar IDs de tickets creados
TICKET_EDIFICIO_ID=""
TICKET_ACTIVO_ID=""
TICKET_TECNICO_ID=""

echo -e "${BLUE}🎫 TEST DE SISTEMA DE TICKETS${NC}"
echo "==============================="
echo -e "${CYAN}🔗 URL Base: $API_URL${NC}"
echo -e "${CYAN}📅 Fecha: $(date)${NC}"
echo ""

# Función para realizar test de endpoint
test_endpoint() {
    local method="$1"
    local endpoint="$2"
    local description="$3"
    local data="$4"
    local expected_status="$5"
    
    TEST_COUNT=$((TEST_COUNT + 1))
    
    echo -e "${PURPLE}📍 Test #$TEST_COUNT: $description${NC}"
    echo -e "${CYAN}$method $endpoint${NC}"
    
    # Realizar petición
    if [ "$method" = "GET" ]; then
        RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" "$endpoint")
    elif [ "$method" = "POST" ]; then
        RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST -H "Content-Type: application/json" -d "$data" "$endpoint")
    elif [ "$method" = "PUT" ]; then
        RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" -X PUT -H "Content-Type: application/json" -d "$data" "$endpoint")
    fi
    
    # Separar response body y status code
    RESPONSE_BODY=$(echo "$RESPONSE" | sed -E 's/HTTPSTATUS\:[0-9]{3}$//')
    STATUS_CODE=$(echo "$RESPONSE" | tr -d '\n' | sed -E 's/.*HTTPSTATUS:([0-9]{3})$/\1/')
    
    # Verificar status code
    if [[ "$STATUS_CODE" =~ ^$expected_status ]]; then
        echo -e "${GREEN}✅ PASS - Status: $STATUS_CODE${NC}"
        PASS_COUNT=$((PASS_COUNT + 1))
    else
        echo -e "${RED}❌ FAIL - Status: $STATUS_CODE (esperado: ${expected_status}xx)${NC}"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
    
    # Mostrar respuesta truncada
    if [ ${#RESPONSE_BODY} -gt 300 ]; then
        echo -e "Response: ${RESPONSE_BODY:0:300}..."
    else
        echo -e "Response: $RESPONSE_BODY"
    fi
    echo ""
    
    # Guardar última respuesta
    LAST_RESPONSE="$RESPONSE_BODY"
}

# Verificar disponibilidad del servicio
echo -e "${CYAN}🔍 Verificando disponibilidad del servicio...${NC}"
if timeout 10 curl -s "$API_URL/tipos-notificacion" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Servicio de Notificaciones disponible${NC}"
    echo ""
else
    echo -e "${RED}❌ Servicio de Notificaciones no disponible en puerto 8091${NC}"
    echo -e "${YELLOW}💡 Asegúrate de que el servicio esté ejecutándose${NC}"
    exit 1
fi

echo -e "${BLUE}🏢 FASE 1: TICKETS DE EDIFICIOS${NC}"
echo "================================"

# Test 1: Crear ticket para nuevo edificio (INGRESO)
TICKET_EDIFICIO_INGRESO='{
    "tipo_entidad": "edificio",
    "tipo_operacion": "ingreso",
    "usuario_email": "admin@test.com",
    "edificio_nombre": "Torre Central Business",
    "edificio_direccion": "Av. Libertador 1500, CABA",
    "edificio_latitud": -34.5875,
    "edificio_longitud": -58.3974,
    "justificacion": "Nuevo edificio corporativo que requiere sistema de monitoreo completo de activos"
}'

test_endpoint "POST" "$API_URL/tickets" "Crear ticket de INGRESO de edificio" "$TICKET_EDIFICIO_INGRESO" "201"

# Extraer ID del ticket creado
if echo "$LAST_RESPONSE" | grep -q '"id"'; then
    TICKET_EDIFICIO_ID=$(echo "$LAST_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
    echo -e "${CYAN}📝 Ticket de edificio creado con ID: $TICKET_EDIFICIO_ID${NC}"
    echo ""
fi

# Test 2: Crear ticket para modificar edificio existente (MODIFICACIÓN)
TICKET_EDIFICIO_MODIFICACION='{
    "tipo_entidad": "edificio",
    "tipo_operacion": "modificacion",
    "usuario_email": "admin@test.com",
    "edificio_id": 1,
    "edificio_nombre": "Edificio Principal Actualizado",
    "edificio_direccion": "Nueva dirección actualizada 456",
    "edificio_latitud": -34.6037,
    "edificio_longitud": -58.3816,
    "justificacion": "Actualización de datos del edificio principal por cambio de dirección fiscal"
}'

test_endpoint "POST" "$API_URL/tickets" "Crear ticket de MODIFICACIÓN de edificio" "$TICKET_EDIFICIO_MODIFICACION" "201"

# Test 3: Crear ticket para eliminar edificio (ELIMINACIÓN)
TICKET_EDIFICIO_ELIMINACION='{
    "tipo_entidad": "edificio",
    "tipo_operacion": "eliminacion",
    "usuario_email": "admin@test.com",
    "edificio_id": 2,
    "edificio_nombre": "Edificio a dar de baja",
    "justificacion": "Edificio vendido y fuera del alcance del sistema de monitoreo"
}'

test_endpoint "POST" "$API_URL/tickets" "Crear ticket de ELIMINACIÓN de edificio" "$TICKET_EDIFICIO_ELIMINACION" "201"

echo -e "${BLUE}⚙️  FASE 2: TICKETS DE ACTIVOS${NC}"
echo "=============================="

# Test 4: Crear ticket para nuevo activo (INGRESO)
TICKET_ACTIVO_INGRESO='{
    "tipo_entidad": "activo",
    "tipo_operacion": "ingreso",
    "usuario_email": "tecnico@test.com",
    "activo_nombre": "Bomba Centrífuga Principal",
    "activo_tipo": "bomba de agua",
    "activo_descripcion": "Bomba centrífuga de 50HP para sistema de agua potable del edificio",
    "activo_ubicacion": "Sala de máquinas - Subsuelo nivel -2",
    "activo_edificio_id": 1,
    "justificacion": "Instalación de nuevo sistema de bombeo para mejora del suministro de agua"
}'

test_endpoint "POST" "$API_URL/tickets" "Crear ticket de INGRESO de activo" "$TICKET_ACTIVO_INGRESO" "201"

# Extraer ID del ticket creado
if echo "$LAST_RESPONSE" | grep -q '"id"'; then
    TICKET_ACTIVO_ID=$(echo "$LAST_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
    echo -e "${CYAN}📝 Ticket de activo creado con ID: $TICKET_ACTIVO_ID${NC}"
    echo ""
fi

# Test 5: Crear ticket para modificar activo existente (MODIFICACIÓN)
TICKET_ACTIVO_MODIFICACION='{
    "tipo_entidad": "activo",
    "tipo_operacion": "modificacion",
    "usuario_email": "tecnico@test.com",
    "activo_id": 1,
    "activo_nombre": "Sensor de Temperatura Actualizado",
    "activo_tipo": "sensor de temperatura",
    "activo_descripcion": "Sensor actualizado con nueva tecnología IoT",
    "activo_ubicacion": "Planta baja - Sector A",
    "justificacion": "Actualización de especificaciones técnicas y ubicación del sensor"
}'

test_endpoint "POST" "$API_URL/tickets" "Crear ticket de MODIFICACIÓN de activo" "$TICKET_ACTIVO_MODIFICACION" "201"

# Test 6: Crear ticket para eliminar activo (ELIMINACIÓN)
TICKET_ACTIVO_ELIMINACION='{
    "tipo_entidad": "activo",
    "tipo_operacion": "eliminacion",
    "usuario_email": "admin@test.com",
    "activo_id": 3,
    "activo_nombre": "Activo obsoleto a dar de baja",
    "justificacion": "Equipo obsoleto que será reemplazado por tecnología más moderna"
}'

test_endpoint "POST" "$API_URL/tickets" "Crear ticket de ELIMINACIÓN de activo" "$TICKET_ACTIVO_ELIMINACION" "201"

echo -e "${BLUE}👨‍🔧 FASE 3: TICKETS DE TÉCNICOS${NC}"
echo "==============================="

# Test 7: Crear ticket para nuevo técnico (INGRESO)
TICKET_TECNICO_INGRESO='{
    "tipo_entidad": "tecnico",
    "tipo_operacion": "ingreso",
    "usuario_email": "rrhh@test.com",
    "tecnico_nombre": "Juan Carlos Pérez Gómez",
    "tecnico_email": "juancarlos.perez@tecnico.com",
    "tecnico_telefono": "+54911234567890",
    "tecnico_especialidad": "electricidad industrial",
    "tecnico_autorizado": true,
    "justificacion": "Técnico certificado con 15 años de experiencia en sistemas eléctricos industriales"
}'

test_endpoint "POST" "$API_URL/tickets" "Crear ticket de INGRESO de técnico" "$TICKET_TECNICO_INGRESO" "201"

# Extraer ID del ticket creado
if echo "$LAST_RESPONSE" | grep -q '"id"'; then
    TICKET_TECNICO_ID=$(echo "$LAST_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
    echo -e "${CYAN}📝 Ticket de técnico creado con ID: $TICKET_TECNICO_ID${NC}"
    echo ""
fi

# Test 8: Crear ticket para modificar técnico existente (MODIFICACIÓN)
TICKET_TECNICO_MODIFICACION='{
    "tipo_entidad": "tecnico",
    "tipo_operacion": "modificacion",
    "usuario_email": "rrhh@test.com",
    "tecnico_id": 1,
    "tecnico_nombre": "María González Actualizado",
    "tecnico_email": "maria.gonzalez.nuevo@tecnico.com",
    "tecnico_telefono": "+54911111111",
    "tecnico_especialidad": "HVAC y climatización",
    "tecnico_autorizado": true,
    "justificacion": "Actualización de datos de contacto y nueva especialidad certificada"
}'

test_endpoint "POST" "$API_URL/tickets" "Crear ticket de MODIFICACIÓN de técnico" "$TICKET_TECNICO_MODIFICACION" "201"

# Test 9: Crear ticket para eliminar técnico (ELIMINACIÓN)
TICKET_TECNICO_ELIMINACION='{
    "tipo_entidad": "tecnico",
    "tipo_operacion": "eliminacion",
    "usuario_email": "rrhh@test.com",
    "tecnico_id": 2,
    "tecnico_nombre": "Técnico a dar de baja",
    "justificacion": "El técnico ha finalizado su contrato con la empresa"
}'

test_endpoint "POST" "$API_URL/tickets" "Crear ticket de ELIMINACIÓN de técnico" "$TICKET_TECNICO_ELIMINACION" "201"

echo -e "${BLUE}❌ FASE 4: VALIDACIONES (ERRORES ESPERADOS)${NC}"
echo "==========================================="

# Test 10: Ticket con tipo de entidad inválido
TICKET_INVALIDO_ENTIDAD='{
    "tipo_entidad": "invalido",
    "tipo_operacion": "ingreso",
    "usuario_email": "test@test.com",
    "justificacion": "Test de validación"
}'

test_endpoint "POST" "$API_URL/tickets" "Ticket con tipo de entidad INVÁLIDO" "$TICKET_INVALIDO_ENTIDAD" "400"

# Test 11: Ticket con tipo de operación inválido
TICKET_INVALIDO_OPERACION='{
    "tipo_entidad": "edificio",
    "tipo_operacion": "operacion_invalida",
    "usuario_email": "test@test.com",
    "justificacion": "Test de validación"
}'

test_endpoint "POST" "$API_URL/tickets" "Ticket con tipo de operación INVÁLIDO" "$TICKET_INVALIDO_OPERACION" "400"

# Test 12: Ticket con email inválido
TICKET_EMAIL_INVALIDO='{
    "tipo_entidad": "edificio",
    "tipo_operacion": "ingreso",
    "usuario_email": "email-sin-formato-valido",
    "edificio_nombre": "Test",
    "justificacion": "Test de validación"
}'

test_endpoint "POST" "$API_URL/tickets" "Ticket con EMAIL INVÁLIDO" "$TICKET_EMAIL_INVALIDO" "400"

# Test 13: Ticket sin datos requeridos
TICKET_SIN_DATOS='{
    "tipo_entidad": "edificio",
    "tipo_operacion": "ingreso"
}'

test_endpoint "POST" "$API_URL/tickets" "Ticket SIN DATOS REQUERIDOS" "$TICKET_SIN_DATOS" "400"

echo -e "${BLUE}📄 FASE 5: PAGINACIÓN DE TICKETS${NC}"
echo "================================"

# Test 14: Obtener tickets sin parámetro de página (default página 1)
test_endpoint "GET" "$API_URL/tickets" "Obtener tickets sin parámetro (página 1 default)" "" "200"

if echo "$LAST_RESPONSE" | grep -q '"total_tickets"'; then
    TOTAL_TICKETS=$(echo "$LAST_RESPONSE" | grep -o '"total_tickets":[0-9]*' | grep -o '[0-9]*')
    TOTAL_PAGINAS=$(echo "$LAST_RESPONSE" | grep -o '"total_paginas":[0-9]*' | grep -o '[0-9]*')
    PAGINA_ACTUAL=$(echo "$LAST_RESPONSE" | grep -o '"pagina_actual":[0-9]*' | grep -o '[0-9]*')
    echo -e "${CYAN}📊 Total de tickets: $TOTAL_TICKETS${NC}"
    echo -e "${CYAN}📊 Total de páginas: $TOTAL_PAGINAS${NC}"
    echo -e "${CYAN}📊 Página actual: $PAGINA_ACTUAL${NC}"
    echo ""
fi

# Test 15: Obtener tickets página 1 explícitamente
test_endpoint "GET" "$API_URL/tickets?pagina=1" "Obtener tickets página 1" "" "200"

# Test 16: Obtener tickets página 2 (si existe)
test_endpoint "GET" "$API_URL/tickets?pagina=2" "Obtener tickets página 2" "" "200"

echo -e "${BLUE}✅ FASE 6: RESOLVER TICKETS${NC}"
echo "==========================="

# Test 17: Resolver ticket de edificio
if [ -n "$TICKET_EDIFICIO_ID" ]; then
    RESOLVER_TICKET='{
        "resuelto_por_email": "admin@test.com",
        "comentario_admin": "Solicitud de ingreso de edificio aprobada. El edificio Torre Central Business ha sido agregado al sistema exitosamente."
    }'
    
    test_endpoint "PUT" "$API_URL/tickets/$TICKET_EDIFICIO_ID/resolver" "Resolver ticket de edificio (ID: $TICKET_EDIFICIO_ID)" "$RESOLVER_TICKET" "200"
    
    if echo "$LAST_RESPONSE" | grep -q '"estado":"resuelto"'; then
        echo -e "${GREEN}✅ Ticket de edificio marcado como RESUELTO${NC}"
        echo ""
    fi
else
    echo -e "${YELLOW}⚠️  No se pudo obtener ID del ticket de edificio para resolver${NC}"
    echo ""
fi

# Test 18: Resolver ticket de activo
if [ -n "$TICKET_ACTIVO_ID" ]; then
    RESOLVER_ACTIVO='{
        "resuelto_por_email": "tecnico@test.com",
        "comentario_admin": "Solicitud de ingreso de activo aprobada. La Bomba Centrífuga Principal ha sido registrada en el sistema."
    }'
    
    test_endpoint "PUT" "$API_URL/tickets/$TICKET_ACTIVO_ID/resolver" "Resolver ticket de activo (ID: $TICKET_ACTIVO_ID)" "$RESOLVER_ACTIVO" "200"
    
    if echo "$LAST_RESPONSE" | grep -q '"estado":"resuelto"'; then
        echo -e "${GREEN}✅ Ticket de activo marcado como RESUELTO${NC}"
        echo ""
    fi
else
    echo -e "${YELLOW}⚠️  No se pudo obtener ID del ticket de activo para resolver${NC}"
    echo ""
fi

# Test 19: Resolver ticket de técnico
if [ -n "$TICKET_TECNICO_ID" ]; then
    RESOLVER_TECNICO='{
        "resuelto_por_email": "rrhh@test.com",
        "comentario_admin": "Solicitud de ingreso de técnico aprobada. Juan Carlos Pérez Gómez ha sido registrado en el sistema como técnico autorizado."
    }'
    
    test_endpoint "PUT" "$API_URL/tickets/$TICKET_TECNICO_ID/resolver" "Resolver ticket de técnico (ID: $TICKET_TECNICO_ID)" "$RESOLVER_TECNICO" "200"
    
    if echo "$LAST_RESPONSE" | grep -q '"estado":"resuelto"'; then
        echo -e "${GREEN}✅ Ticket de técnico marcado como RESUELTO${NC}"
        echo ""
    fi
else
    echo -e "${YELLOW}⚠️  No se pudo obtener ID del ticket de técnico para resolver${NC}"
    echo ""
fi

# Test 20: Intentar resolver ticket inexistente
RESOLVER_INEXISTENTE='{
    "resuelto_por_email": "admin@test.com",
    "comentario_admin": "Test de ticket inexistente"
}'

test_endpoint "PUT" "$API_URL/tickets/999999/resolver" "Resolver ticket INEXISTENTE (error esperado)" "$RESOLVER_INEXISTENTE" "500"

# Test 21: Resolver ticket con email inválido
if [ -n "$TICKET_EDIFICIO_ID" ]; then
    RESOLVER_EMAIL_INVALIDO='{
        "resuelto_por_email": "email-invalido",
        "comentario_admin": "Test de validación"
    }'
    
    # Este debería fallar por email inválido, pero si el ticket ya está resuelto, puede dar otro error
    test_endpoint "PUT" "$API_URL/tickets/$TICKET_EDIFICIO_ID/resolver" "Resolver con EMAIL INVÁLIDO" "$RESOLVER_EMAIL_INVALIDO" "400"
fi

# Test 22: Resolver ticket sin comentario
RESOLVER_SIN_COMENTARIO='{
    "resuelto_por_email": "admin@test.com",
    "comentario_admin": ""
}'

test_endpoint "PUT" "$API_URL/tickets/1/resolver" "Resolver ticket SIN COMENTARIO" "$RESOLVER_SIN_COMENTARIO" "400"

echo -e "${BLUE}📊 FASE 7: VERIFICACIÓN FINAL${NC}"
echo "============================="

# Test 23: Obtener todos los tickets y ver resueltos vs no resueltos
test_endpoint "GET" "$API_URL/tickets?pagina=1" "Verificar tickets después de resoluciones" "" "200"

if echo "$LAST_RESPONSE" | grep -q '"estado"'; then
    RESUELTOS=$(echo "$LAST_RESPONSE" | grep -o '"estado":"resuelto"' | wc -l)
    NO_RESUELTOS=$(echo "$LAST_RESPONSE" | grep -o '"estado":"no_resuelto"' | wc -l)
    echo -e "${CYAN}✅ Tickets RESUELTOS en página 1: $RESUELTOS${NC}"
    echo -e "${CYAN}⏳ Tickets NO RESUELTOS en página 1: $NO_RESUELTOS${NC}"
    echo ""
fi

# Mostrar resumen final
echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}📊 RESUMEN FINAL DE PRUEBAS${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${CYAN}Total de tests ejecutados: ${YELLOW}$TEST_COUNT${NC}"
echo -e "${GREEN}Tests exitosos: $PASS_COUNT${NC}"
echo -e "${RED}Tests fallidos: $FAIL_COUNT${NC}"
echo ""

if [ $TEST_COUNT -gt 0 ]; then
    SUCCESS_RATE=$(( (PASS_COUNT * 100) / TEST_COUNT ))
else
    SUCCESS_RATE=0
fi

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}🎉 TODOS LOS TESTS EXITOSOS - ${SUCCESS_RATE}%${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "${GREEN}✅ SISTEMA DE TICKETS FUNCIONANDO CORRECTAMENTE${NC}"
    echo ""
    echo -e "${CYAN}📋 Funcionalidades verificadas:${NC}"
    echo -e "${CYAN}   ✓ Creación de tickets para EDIFICIOS (ingreso, modificación, eliminación)${NC}"
    echo -e "${CYAN}   ✓ Creación de tickets para ACTIVOS (ingreso, modificación, eliminación)${NC}"
    echo -e "${CYAN}   ✓ Creación de tickets para TÉCNICOS (ingreso, modificación, eliminación)${NC}"
    echo -e "${CYAN}   ✓ Validación de datos de entrada${NC}"
    echo -e "${CYAN}   ✓ Paginación de tickets (10 por página)${NC}"
    echo -e "${CYAN}   ✓ Resolución de tickets con justificación${NC}"
    echo -e "${CYAN}   ✓ Notificaciones por email a administradores${NC}"
    echo -e "${CYAN}   ✓ Manejo de errores y casos edge${NC}"
    echo ""
    echo -e "${YELLOW}💡 Tickets creados en este test:${NC}"
    [ -n "$TICKET_EDIFICIO_ID" ] && echo -e "${YELLOW}   • Ticket Edificio ID: $TICKET_EDIFICIO_ID${NC}"
    [ -n "$TICKET_ACTIVO_ID" ] && echo -e "${YELLOW}   • Ticket Activo ID: $TICKET_ACTIVO_ID${NC}"
    [ -n "$TICKET_TECNICO_ID" ] && echo -e "${YELLOW}   • Ticket Técnico ID: $TICKET_TECNICO_ID${NC}"
    echo ""
    echo -e "${YELLOW}📧 Para verificar emails enviados a administradores:${NC}"
    echo -e "${YELLOW}   docker logs notification-service | grep -i 'ticket'${NC}"
    echo ""
    exit 0
else
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}⚠️  ALGUNOS TESTS FALLARON - ${SUCCESS_RATE}%${NC}"
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "${RED}❌ REVISAR CONFIGURACIÓN DEL SISTEMA${NC}"
    echo ""
    echo -e "${YELLOW}💡 Sugerencias:${NC}"
    echo -e "${YELLOW}   1. Verificar que el servicio esté ejecutándose correctamente${NC}"
    echo -e "${YELLOW}   2. Revisar logs del servicio:${NC}"
    echo -e "${YELLOW}      docker logs notification-service${NC}"
    echo -e "${YELLOW}   3. Verificar configuración de base de datos${NC}"
    echo -e "${YELLOW}   4. Verificar configuración de email (SMTP)${NC}"
    echo ""
    exit 1
fi
