#!/bin/bash

# ===================================================================
# TEST SIMPLE DE SINCRONIZACIÓN DE BASES DE DATOS
# ===================================================================

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 TEST DE SINCRONIZACIÓN DE BASES DE DATOS${NC}"
echo "============================================="

# Verificar que los contenedores estén ejecutándose
echo -e "\n${YELLOW}🐳 Verificando contenedores...${NC}"
containers=("gestion-db" "notification-db" "documentacion-db" "mongu")
all_ok=true

for container in "${containers[@]}"; do
    if docker ps | grep -q "$container"; then
        echo -e "   ✅ $container está ejecutándose"
    else
        echo -e "   ❌ $container no está ejecutándose"
        all_ok=false
    fi
done

if [ "$all_ok" = false ]; then
    echo -e "\n${RED}❌ Algunos contenedores no están disponibles${NC}"
    echo "Ejecuta: docker-compose up -d"
    exit 1
fi

echo -e "\n${YELLOW}📊 Verificando conteos de datos...${NC}"

# Contar registros en cada DB
echo -e "\n${BLUE}📋 Conteos por base de datos:${NC}"

gestion_count=$(docker exec gestion-db psql -U gestion_user -d gestion_db -t -c "SELECT COUNT(*) FROM activos;" 2>/dev/null | tr -d ' ')
echo "   • Gestión DB: $gestion_count activos"

notification_count=$(docker exec notification-db psql -U notification_user -d notification_db -t -c "SELECT COUNT(*) FROM activos;" 2>/dev/null | tr -d ' ')
echo "   • Notification DB: $notification_count activos"

doc_count=$(docker exec documentacion-db psql -U documentacion_user -d documentacion_db -t -c "SELECT COUNT(DISTINCT activo_id) FROM documentos;" 2>/dev/null | tr -d ' ')
echo "   • Documentación DB: $doc_count activos únicos"

mongo_count=$(docker exec mongu mongosh iot_db --eval "db.activos.countDocuments()" --quiet 2>/dev/null | tail -1)
echo "   • MongoDB IoT: $mongo_count activos"

# Verificar sincronización básica
echo -e "\n${YELLOW}🔍 Análisis de sincronización:${NC}"

if [ "$gestion_count" = "8" ] && [ "$mongo_count" = "8" ]; then
    echo -e "   ✅ Gestión DB y MongoDB están sincronizados (8 activos cada uno)"
    sync_status="parcial"
else
    echo -e "   ❌ Gestión DB ($gestion_count) y MongoDB ($mongo_count) no están sincronizados"
    sync_status="fallido"
fi

if [ "$notification_count" = "8" ]; then
    echo -e "   ✅ Notification DB tiene el conteo correcto (8 activos)"
    if [ "$sync_status" = "parcial" ]; then
        sync_status="bueno"
    fi
else
    echo -e "   ❌ Notification DB tiene $notification_count activos (esperados: 8)"
fi

if [ "$doc_count" = "8" ]; then
    echo -e "   ✅ Documentación DB tiene activos para todos los IDs (8 únicos)"
    if [ "$sync_status" = "bueno" ]; then
        sync_status="perfecto"
    fi
else
    echo -e "   ❌ Documentación DB solo tiene $doc_count activos únicos (esperados: 8)"
fi

# Mostrar resumen
echo -e "\n${BLUE}📋 RESUMEN FINAL:${NC}"
case $sync_status in
    "perfecto")
        echo -e "   🎉 ${GREEN}SINCRONIZACIÓN PERFECTA${NC} - Todas las bases están sincronizadas"
        exit 0
        ;;
    "bueno")
        echo -e "   ✅ ${GREEN}SINCRONIZACIÓN BUENA${NC} - Core databases sincronizadas, Documentación necesita atención"
        echo -e "   💡 Recomendación: Agregar más documentos a la base de documentación"
        exit 0
        ;;
    "parcial")
        echo -e "   ⚠️ ${YELLOW}SINCRONIZACIÓN PARCIAL${NC} - MongoDB y Gestión DB sincronizados"
        echo -e "   🔧 Notification DB necesita corrección"
        exit 1
        ;;
    "fallido")
        echo -e "   ❌ ${RED}SINCRONIZACIÓN FALLIDA${NC} - Múltiples problemas detectados"
        echo -e "   🔧 Se requiere corrección completa"
        exit 1
        ;;
esac
