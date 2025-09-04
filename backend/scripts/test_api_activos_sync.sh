#!/bin/bash

# Test complejo para verificar sincronización de activos usando APIs de microservicios
# Verifica que los datos de id_activo sean consistentes entre todos los servicios

set -e

echo "==================================="
echo "🔍 TEST API ACTIVOS SINCRONIZACIÓN"
echo "==================================="
echo "$(date '+%Y-%m-%d %H:%M:%S') - Iniciando test complejo de APIs de microservicios"

# Configuración de URLs de API Gateway
API_GATEWAY="http://localhost:3500"
GESTION_ENDPOINT="$API_GATEWAY/api/gestion"
PARSER_ENDPOINT="$API_GATEWAY/api/parser"
DOCUMENTACION_ENDPOINT="$API_GATEWAY/api/documentacion"

# URLs directas de microservicios (fallback)
GESTION_DIRECT="http://localhost:8092"
PARSER_DIRECT="http://localhost:8090"
DOCUMENTACION_DIRECT="http://localhost:8093"

# Función para verificar si un servicio está disponible
check_service() {
    local service_name=$1
    local url=$2
    local timeout=5
    
    echo "🔍 Verificando $service_name en $url..."
    if curl -s --max-time $timeout "$url/health" > /dev/null 2>&1; then
        echo "✅ $service_name está disponible"
        return 0
    else
        echo "❌ $service_name no está disponible en $url"
        return 1
    fi
}

# Función para obtener activos de un endpoint
get_activos() {
    local service_name=$1
    local endpoint=$2
    local fallback_endpoint=$3
    
    # Determinar la ruta correcta según el servicio
    local route="/activos"
    if [[ "$service_name" == *"Parser"* || "$service_name" == *"IoT"* ]]; then
        route="/activo"
    fi
    
    # Intentar primero con API Gateway
    response=$(curl -s --max-time 10 "$endpoint$route" 2>/dev/null || echo "")
    
    # Si falla, intentar directo
    if [[ -z "$response" || "$response" == *"error"* ]]; then
        response=$(curl -s --max-time 10 "$fallback_endpoint$route" 2>/dev/null || echo "")
    fi
    
    if [[ -z "$response" || "$response" == *"error"* ]]; then
        return 1
    fi
    
    echo "$response"
}

# Función para extraer IDs de la respuesta JSON
extract_ids() {
    local json_response=$1
    local service_name=$2
    
    if [[ "$service_name" == "ParserService" ]]; then
        # Para IoT/Parser service, extraer activo_id de la estructura MongoDB
        if [[ "$json_response" == "null" || -z "$json_response" ]]; then
            echo ""
        else
            echo "$json_response" | jq -r '.[] | .activo_id' 2>/dev/null | sort
        fi
    else
        # Para servicios PostgreSQL, extraer id desde .activos array
        if echo "$json_response" | jq -e '.activos' > /dev/null 2>&1; then
            echo "$json_response" | jq -r '.activos[] | .id' 2>/dev/null | sort -n
        else
            echo "$json_response" | jq -r '.[] | .id // .activo_id' 2>/dev/null | sort -n
        fi
    fi
}

# Función para contar activos
count_activos() {
    local json_response=$1
    
    if [[ "$json_response" == "null" || -z "$json_response" ]]; then
        echo "0"
    elif [[ "$json_response" == *"404"* ]]; then
        echo "0"
    elif echo "$json_response" | jq -e '.total' > /dev/null 2>&1; then
        # Si tiene campo total, usarlo
        echo "$json_response" | jq '.total' 2>/dev/null || echo "0"
    elif echo "$json_response" | jq -e '.activos' > /dev/null 2>&1; then
        # Si tiene array activos, contar elementos
        echo "$json_response" | jq '.activos | length' 2>/dev/null || echo "0"
    else
        # Array directo
        echo "$json_response" | jq length 2>/dev/null || echo "0"
    fi
}

echo ""
echo "🚀 Iniciando verificación de servicios..."

# Verificar que los servicios estén disponibles
services_available=true

if ! check_service "API Gateway" "http://localhost:3500"; then
    echo "⚠️  API Gateway no disponible, usando conexiones directas"
fi

if ! check_service "Gestión Service" "$GESTION_DIRECT"; then
    services_available=false
fi

if ! check_service "Parser/IoT Service" "$PARSER_DIRECT"; then
    services_available=false
fi

# Documentación service es opcional para este test
check_service "Documentación Service" "$DOCUMENTACION_DIRECT" || echo "⚠️  Servicio de documentación no disponible (opcional)"

if [ "$services_available" = false ]; then
    echo ""
    echo "❌ ERROR: Servicios críticos no están disponibles"
    echo "🔧 Asegúrate de que Docker Compose esté ejecutándose: docker-compose up -d"
    exit 1
fi

echo ""
echo "📊 Obteniendo datos de activos de cada microservicio..."

# Variables para almacenar resultados
gestion_response=""
parser_response=""
documentacion_response=""

# Obtener activos del servicio de gestión
echo ""
echo "1️⃣ === SERVICIO DE GESTIÓN ==="
echo "📡 Obteniendo activos de Gestión Service..."
gestion_response=$(get_activos "Gestión Service" "$GESTION_ENDPOINT" "$GESTION_DIRECT")
if [[ $? -ne 0 ]]; then
    echo "❌ Error crítico: No se pudo obtener datos del servicio de gestión"
    exit 1
fi

gestion_count=$(count_activos "$gestion_response")
echo "📋 Activos encontrados en Gestión: $gestion_count"

# Obtener activos del servicio parser/IoT
echo ""
echo "2️⃣ === SERVICIO PARSER/IOT ==="
echo "📡 Obteniendo activos de Parser/IoT Service..."
parser_response=$(get_activos "Parser/IoT Service" "$PARSER_ENDPOINT" "$PARSER_DIRECT")
if [[ $? -ne 0 ]]; then
    echo "❌ Error crítico: No se pudo obtener datos del servicio Parser/IoT"
    exit 1
fi

parser_count=$(count_activos "$parser_response")
echo "📋 Activos encontrados en Parser/IoT: $parser_count"

# Intentar obtener activos del servicio de documentación (opcional)
echo ""
echo "3️⃣ === SERVICIO DOCUMENTACIÓN ==="
echo "📡 Obteniendo activos de Documentación Service..."
documentacion_response=$(get_activos "Documentación Service" "$DOCUMENTACION_ENDPOINT" "$DOCUMENTACION_DIRECT" 2>/dev/null || echo "")
if [[ -n "$documentacion_response" && "$documentacion_response" != *"error"* ]]; then
    documentacion_count=$(count_activos "$documentacion_response")
    echo "📋 Activos encontrados en Documentación: $documentacion_count"
else
    echo "⚠️  Servicio de documentación no disponible o sin endpoint /activos"
    documentacion_count=0
fi

echo ""
echo "🔍 ANÁLISIS DE CONSISTENCIA DE DATOS"
echo "===================================="

# Extraer IDs de cada servicio
echo "📝 Extrayendo IDs de activos..."

gestion_ids=$(extract_ids "$gestion_response" "GestionService")
parser_ids=$(extract_ids "$parser_response" "ParserService")

echo ""
echo "🔢 IDs de Gestión Service:"
echo "$gestion_ids" | head -10
if [[ $(echo "$gestion_ids" | wc -l) -gt 10 ]]; then
    echo "... (mostrando primeros 10)"
fi

echo ""
echo "🔢 IDs de Parser/IoT Service:"
echo "$parser_ids" | head -10
if [[ $(echo "$parser_ids" | wc -l) -gt 10 ]]; then
    echo "... (mostrando primeros 10)"
fi

# Verificar correspondencia de IDs
echo ""
echo "🔄 VERIFICACIÓN DE CORRESPONDENCIA"
echo "================================="

# Mapeo esperado: PostgreSQL IDs 1-8 ↔ MongoDB códigos AC-1001 a AC-1008
expected_postgresql_ids="1 2 3 4 5 6 7 8"
expected_mongodb_codes="AC-1001 AC-1002 AC-1003 AC-1004 AC-1005 AC-1006 AC-1007 AC-1008"

# Verificar IDs de PostgreSQL (Gestión)
echo "🔍 Verificando IDs de PostgreSQL..."
gestion_ids_clean=$(echo "$gestion_ids" | tr '\n' ' ' | xargs)
if [[ "$gestion_ids_clean" == "$expected_postgresql_ids" ]]; then
    echo "✅ IDs de Gestión Service correctos: $gestion_ids_clean"
    gestion_sync=true
else
    echo "❌ IDs de Gestión Service incorrectos:"
    echo "   Esperado: $expected_postgresql_ids"
    echo "   Actual:   $gestion_ids_clean"
    gestion_sync=false
fi

# Verificar códigos de MongoDB (Parser/IoT)
echo ""
echo "🔍 Verificando códigos de MongoDB..."
parser_ids_clean=$(echo "$parser_ids" | tr '\n' ' ' | xargs)
if [[ "$parser_ids_clean" == "$expected_mongodb_codes" ]]; then
    echo "✅ Códigos de Parser/IoT Service correctos: $parser_ids_clean"
    parser_sync=true
else
    echo "❌ Códigos de Parser/IoT Service incorrectos:"
    echo "   Esperado: $expected_mongodb_codes"
    echo "   Actual:   $parser_ids_clean"
    parser_sync=false
fi

# Verificar correspondencia 1:1
echo ""
echo "🔍 Verificando correspondencia 1:1..."
gestion_count_clean=$(echo "$gestion_count" | tr -d '\n\r ')
parser_count_clean=$(echo "$parser_count" | tr -d '\n\r ')

if [[ "$gestion_count_clean" -eq 8 && "$parser_count_clean" -eq 8 ]]; then
    echo "✅ Ambos servicios tienen 8 activos (correspondencia numérica correcta)"
    count_sync=true
else
    echo "❌ Discrepancia en cantidad de activos:"
    echo "   Gestión: $gestion_count_clean activos"
    echo "   Parser/IoT: $parser_count_clean activos"
    count_sync=false
fi

# Test de endpoints específicos
echo ""
echo "🎯 TEST DE ENDPOINTS ESPECÍFICOS"
echo "================================"

# Test de activo específico en Gestión Service
echo "📡 Testando GET /activos/1 en Gestión Service..."
activo1_gestion=$(curl -s --max-time 5 "$GESTION_DIRECT/activos/1" 2>/dev/null || echo "")
if [[ -n "$activo1_gestion" && "$activo1_gestion" != *"error"* ]]; then
    activo1_name=$(echo "$activo1_gestion" | jq -r '.nombre // .activo_nombre // "N/A"' 2>/dev/null)
    echo "✅ Activo 1 encontrado: $activo1_name"
    specific_test_gestion=true
else
    echo "❌ Error al obtener activo específico ID=1 de Gestión"
    specific_test_gestion=false
fi

# Test de activo específico en Parser/IoT Service  
echo "📡 Testando GET /activo/AC-1001 en Parser/IoT Service..."
activo_ac1001=$(curl -s --max-time 5 "$PARSER_DIRECT/activo/AC-1001" 2>/dev/null || echo "")
if [[ -n "$activo_ac1001" && "$activo_ac1001" != *"error"* ]]; then
    activo_ac1001_name=$(echo "$activo_ac1001" | jq -r '.nombre // .activo_nombre // "N/A"' 2>/dev/null)
    echo "✅ Activo AC-1001 encontrado: $activo_ac1001_name"
    specific_test_parser=true
else
    echo "❌ Error al obtener activo específico AC-1001 de Parser/IoT"
    specific_test_parser=false
fi

# Verificar si los nombres coinciden (si ambos tests pasaron)
if [[ "$specific_test_gestion" = true && "$specific_test_parser" = true ]]; then
    if [[ "$activo1_name" == "$activo_ac1001_name" ]]; then
        echo "✅ Nombres de activos coinciden: ID=1 ↔ AC-1001 = '$activo1_name'"
        name_consistency=true
    else
        echo "⚠️  Nombres de activos difieren:"
        echo "   ID=1 (Gestión): '$activo1_name'"
        echo "   AC-1001 (Parser): '$activo_ac1001_name'"
        name_consistency=false
    fi
else
    name_consistency=false
fi

# RESUMEN FINAL
echo ""
echo "🏁 RESUMEN DE RESULTADOS"
echo "========================"

total_tests=6
passed_tests=0

echo "📊 Resultados por categoría:"

if [[ "$gestion_sync" = true ]]; then
    echo "✅ 1. IDs PostgreSQL (Gestión): CORRECTO"
    passed_tests=$((passed_tests + 1))
else
    echo "❌ 1. IDs PostgreSQL (Gestión): INCORRECTO"
fi

if [[ "$parser_sync" = true ]]; then
    echo "✅ 2. Códigos MongoDB (Parser/IoT): CORRECTO"
    passed_tests=$((passed_tests + 1))
else
    echo "❌ 2. Códigos MongoDB (Parser/IoT): INCORRECTO"
fi

if [[ "$count_sync" = true ]]; then
    echo "✅ 3. Cantidad de activos: CORRECTO"
    passed_tests=$((passed_tests + 1))
else
    echo "❌ 3. Cantidad de activos: INCORRECTO"
fi

if [[ "$specific_test_gestion" = true ]]; then
    echo "✅ 4. Endpoint específico Gestión: CORRECTO"
    passed_tests=$((passed_tests + 1))
else
    echo "❌ 4. Endpoint específico Gestión: INCORRECTO"
fi

if [[ "$specific_test_parser" = true ]]; then
    echo "✅ 5. Endpoint específico Parser/IoT: CORRECTO"
    passed_tests=$((passed_tests + 1))
else
    echo "❌ 5. Endpoint específico Parser/IoT: INCORRECTO"
fi

if [[ "$name_consistency" = true ]]; then
    echo "✅ 6. Consistencia de nombres: CORRECTO"
    passed_tests=$((passed_tests + 1))
else
    echo "❌ 6. Consistencia de nombres: INCORRECTO"
fi

# Resultado final
echo ""
echo "🎯 RESULTADO FINAL: $passed_tests/$total_tests tests pasaron"

if [[ $passed_tests -eq $total_tests ]]; then
    echo "🎉 ¡ÉXITO! Todos los tests de sincronización de APIs pasaron"
    echo "✅ Los activos están perfectamente sincronizados entre microservicios"
    exit 0
elif [[ $passed_tests -ge 4 ]]; then
    echo "⚠️  PARCIAL: La mayoría de tests pasaron, pero hay problemas menores"
    exit 1
else
    echo "❌ FALLO: Problemas graves de sincronización detectados"
    echo "🔧 Ejecuta el script de inicialización: ./scripts/simple_data_init.sh"
    exit 2
fi
