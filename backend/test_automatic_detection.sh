#!/bin/bash

# 🚨 TEST DE DETECCIÓN AUTOMÁTICA DE SENSORES DESCONECTADOS
# Este test verifica que el IoT service detecte automáticamente sensores desconectados
# SIN forzar la verificación manual

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8090"
NOTIFICATION_URL="http://localhost:8091"

echo -e "${BLUE}🚨 TEST DE DETECCIÓN AUTOMÁTICA DE SENSORES DESCONECTADOS${NC}"
echo "=============================================================="
echo -e "🔗 IoT Service URL: $BASE_URL"
echo -e "📧 Notification Service URL: $NOTIFICATION_URL"
echo -e "📅 Fecha: $(date)"
echo ""

# Función para esperar que el servicio esté disponible
wait_for_service() {
    local url=$1
    local service_name=$2
    local max_attempts=30
    local attempt=1

    echo -e "${YELLOW}🔍 Esperando que $service_name esté disponible...${NC}"
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s "$url/api/sensors/health" >/dev/null 2>&1; then
            echo -e "${GREEN}✅ $service_name está disponible${NC}"
            return 0
        fi
        echo -e "   Intento $attempt/$max_attempts..."
        sleep 2
        ((attempt++))
    done
    
    echo -e "${RED}❌ $service_name no está disponible después de $max_attempts intentos${NC}"
    exit 1
}

# Esperar que los servicios estén disponibles
wait_for_service "$BASE_URL" "IoT Service"

echo ""
echo -e "${BLUE}📊 FASE 1: ESTADO INICIAL DEL SISTEMA${NC}"
echo "================================================"

# Verificar estado inicial
echo -e "${PURPLE}🔍 Estado inicial del sistema:${NC}"
curl -s "$BASE_URL/api/sensors/stats" | jq '.'

echo ""
echo -e "${BLUE}📡 FASE 2: CREAR SENSORES CON DATOS ANTIGUOS${NC}"
echo "================================================"

# Timestamp muy antiguo (6 minutos atrás) para asegurar que se detecte como desconectado
TIMESTAMP_OLD=$(date -d '6 minutes ago' -u +"%Y-%m-%dT%H:%M:%SZ")
echo -e "${PURPLE}📅 Timestamp antiguo para simular desconexión: $TIMESTAMP_OLD${NC}"

# Crear activo de prueba
ACTIVO_DATA='{
    "id": "AUTO_TEST_ACTIVO_001",
    "nombre": "Activo Test Detección Automática",
    "id_edificio": "test_building_auto"
}'

echo -e "${PURPLE}🏗️  Creando activo de prueba...${NC}"
ACTIVO_RESPONSE=$(curl -s -X POST "$BASE_URL/activo" \
    -H "Content-Type: application/json" \
    -d "$ACTIVO_DATA")
echo "Respuesta: $ACTIVO_RESPONSE"

# Enviar datos de sensores con timestamp antiguo
echo -e "${PURPLE}📊 Enviando datos de sensores con timestamp antiguo...${NC}"

# Sensor 1
SENSOR_DATA_1='{
    "sensor_id": "AUTO_TEMP_SENSOR_001",
    "activo_id": "AUTO_TEST_ACTIVO_001",
    "tipo": "temperatura",
    "valor": 22.5,
    "unidad": "°C",
    "timestamp": "'$TIMESTAMP_OLD'"
}'

echo -e "   📡 Enviando sensor temperatura..."
curl -s -X POST "$BASE_URL/lectura" \
    -H "Content-Type: application/json" \
    -d "$SENSOR_DATA_1"

# Sensor 2
SENSOR_DATA_2='{
    "sensor_id": "AUTO_PRESS_SENSOR_002", 
    "activo_id": "AUTO_TEST_ACTIVO_001",
    "tipo": "presion",
    "valor": 1.8,
    "unidad": "bar",
    "timestamp": "'$TIMESTAMP_OLD'"
}'

echo -e "   📡 Enviando sensor presión..."
curl -s -X POST "$BASE_URL/lectura" \
    -H "Content-Type: application/json" \
    -d "$SENSOR_DATA_2"

# Sensor 3
SENSOR_DATA_3='{
    "sensor_id": "AUTO_ELEC_SENSOR_003",
    "activo_id": "AUTO_TEST_ACTIVO_001", 
    "tipo": "corriente",
    "valor": 15.2,
    "unidad": "A",
    "timestamp": "'$TIMESTAMP_OLD'"
}'

echo -e "   📡 Enviando sensor eléctrico..."
curl -s -X POST "$BASE_URL/lectura" \
    -H "Content-Type: application/json" \
    -d "$SENSOR_DATA_3"

echo ""
echo -e "${GREEN}✅ Sensores creados con datos antiguos${NC}"

echo ""
echo -e "${BLUE}📊 FASE 3: VERIFICAR ESTADO DESPUÉS DE CREAR SENSORES${NC}"
echo "========================================================"

echo -e "${PURPLE}🔍 Estado del sistema después de crear sensores:${NC}"
curl -s "$BASE_URL/api/sensors/stats" | jq '.'

echo ""
echo -e "${PURPLE}🔍 Estado detallado de todos los sensores:${NC}"
curl -s "$BASE_URL/api/sensors/status" | jq '.'

echo ""
echo -e "${BLUE}⏳ FASE 4: ESPERAR DETECCIÓN AUTOMÁTICA${NC}"
echo "=============================================="

echo -e "${YELLOW}⏱️  El IoT service verifica sensores cada 2 minutos${NC}"
echo -e "${YELLOW}🔍 Esperando detección automática de sensores desconectados...${NC}"
echo -e "${YELLOW}📊 Monitoreando por 3 ciclos de verificación (6 minutos máximo)${NC}"

# Monitorear por 6 minutos máximo, verificando cada 30 segundos
MAX_WAIT_TIME=360  # 6 minutos
CHECK_INTERVAL=30  # 30 segundos
elapsed_time=0
notifications_detected=0

while [ $elapsed_time -lt $MAX_WAIT_TIME ]; do
    sleep $CHECK_INTERVAL
    elapsed_time=$((elapsed_time + CHECK_INTERVAL))
    
    echo ""
    echo -e "${PURPLE}⏱️  Tiempo transcurrido: ${elapsed_time}s / ${MAX_WAIT_TIME}s${NC}"
    
    # Verificar estadísticas
    stats=$(curl -s "$BASE_URL/api/sensors/stats")
    inactive_count=$(echo "$stats" | jq -r '.inactive_sensors')
    
    echo -e "${PURPLE}📊 Sensores inactivos detectados: $inactive_count${NC}"
    
    # Si encontramos sensores inactivos, verificar notificaciones
    if [ "$inactive_count" -gt 0 ]; then
        echo -e "${GREEN}🎉 ¡Sensores detectados como inactivos!${NC}"
        
        # Verificar logs del IoT service para notificaciones
        echo -e "${PURPLE}🔍 Verificando logs de notificaciones...${NC}"
        recent_notifications=$(docker logs iot-service 2>&1 | grep -i "alerta.*enviada exitosamente" | tail -5)
        
        if [ ! -z "$recent_notifications" ]; then
            notifications_detected=$((notifications_detected + 1))
            echo -e "${GREEN}📧 Notificaciones encontradas en logs:${NC}"
            echo "$recent_notifications"
        fi
        
        # Si detectamos notificaciones, salir del loop
        if [ $notifications_detected -gt 0 ]; then
            echo -e "${GREEN}✅ ¡Detección automática y notificaciones funcionando!${NC}"
            break
        fi
    fi
    
    echo -e "${YELLOW}🔄 Continuando monitoreo...${NC}"
done

echo ""
echo -e "${BLUE}📊 FASE 5: RESULTADOS FINALES${NC}"
echo "=================================="

# Estado final del sistema
echo -e "${PURPLE}🔍 Estado final del sistema:${NC}"
final_stats=$(curl -s "$BASE_URL/api/sensors/stats")
echo "$final_stats" | jq '.'

final_inactive=$(echo "$final_stats" | jq -r '.inactive_sensors')

# Verificar logs finales de notificaciones
echo ""
echo -e "${PURPLE}📧 Logs finales de notificaciones (últimas 10):${NC}"
docker logs iot-service 2>&1 | grep -i "alerta.*enviada exitosamente\|sensor.*desconectó\|notificación" | tail -10

# Verificar logs del servicio de notificaciones
echo ""
echo -e "${PURPLE}📧 Logs del servicio de notificaciones (últimos 10):${NC}"
docker logs notification-service 2>&1 | grep -i "correo.*enviado\|email.*enviado\|notificación" | tail -10

echo ""
echo -e "${BLUE}📋 RESUMEN DEL TEST${NC}"
echo "======================"

if [ "$final_inactive" -gt 0 ] && [ $notifications_detected -gt 0 ]; then
    echo -e "${GREEN}🎉 ¡TEST EXITOSO!${NC}"
    echo -e "${GREEN}✅ Sensores detectados como inactivos: $final_inactive${NC}"
    echo -e "${GREEN}✅ Notificaciones automáticas enviadas: $notifications_detected${NC}"
    echo -e "${GREEN}✅ Sistema de detección automática funcionando correctamente${NC}"
    exit 0
elif [ "$final_inactive" -gt 0 ]; then
    echo -e "${YELLOW}⚠️  TEST PARCIAL${NC}"
    echo -e "${GREEN}✅ Sensores detectados como inactivos: $final_inactive${NC}"
    echo -e "${RED}❌ No se detectaron notificaciones automáticas${NC}"
    exit 1
else
    echo -e "${RED}❌ TEST FALLIDO${NC}"
    echo -e "${RED}❌ No se detectaron sensores inactivos${NC}"
    echo -e "${RED}❌ Sistema de detección automática no funcionó${NC}"
    exit 1
fi
