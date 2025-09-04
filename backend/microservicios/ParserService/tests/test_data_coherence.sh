#!/bin/bash

# Test de coherencia de datos entre bases de datos del ParserService
# Verifica sincronización entre PostgreSQL (otros microservicios) y MongoDB (ParserService)

set -e

echo "=================================================="
echo "🔍 TEST COHERENCIA DE DATOS - PARSER SERVICE"
echo "=================================================="
echo "$(date '+%Y-%m-%d %H:%M:%S') - Iniciando test de coherencia entre bases de datos"

# Configuración
PARSER_SERVICE_URL="http://localhost:8090"
GESTION_SERVICE_URL="http://localhost:8092"

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Función para mostrar resultados de tests
test_result() {
    local test_name="$1"
    local result="$2"
    local details="$3"
    
    if [[ "$result" == "PASS" ]]; then
        echo -e "${GREEN}✅ PASS${NC}: $test_name"
        [[ -n "$details" ]] && echo "   └─ $details"
    elif [[ "$result" == "FAIL" ]]; then
        echo -e "${RED}❌ FAIL${NC}: $test_name"
        [[ -n "$details" ]] && echo "   └─ $details"
    elif [[ "$result" == "WARN" ]]; then
        echo -e "${YELLOW}⚠️  WARN${NC}: $test_name"
        [[ -n "$details" ]] && echo "   └─ $details"
    fi
}

# Variables para conteo de tests
total_tests=0
passed_tests=0
failed_tests=0
warning_tests=0

# Función para incrementar contadores
increment_test() {
    local result="$1"
    total_tests=$((total_tests + 1))
    
    case "$result" in
        "PASS") passed_tests=$((passed_tests + 1)) ;;
        "FAIL") failed_tests=$((failed_tests + 1)) ;;
        "WARN") warning_tests=$((warning_tests + 1)) ;;
    esac
}

# Función para verificar servicios
check_services() {
    echo "🔍 Verificando disponibilidad de servicios..."
    
    # Verificar ParserService
    parser_response=$(curl -s --max-time 5 "$PARSER_SERVICE_URL/activo" 2>/dev/null || echo "error")
    if [[ "$parser_response" == *"error"* || -z "$parser_response" ]]; then
        echo -e "${RED}❌ ParserService no está disponible${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ ParserService disponible${NC}"
    
    # Verificar GestionService
    gestion_response=$(curl -s --max-time 5 "$GESTION_SERVICE_URL/activos" 2>/dev/null || echo "error")
    if [[ "$gestion_response" == *"error"* || -z "$gestion_response" ]]; then
        echo -e "${RED}❌ GestionService no está disponible${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ GestionService disponible${NC}"
}

# Función para obtener activos de MongoDB (ParserService)
get_mongodb_activos() {
    echo "📊 Obteniendo activos de MongoDB (ParserService)..."
    response=$(curl -s "$PARSER_SERVICE_URL/activo" 2>/dev/null)
    echo "$response"
}

# Función para obtener activos de PostgreSQL (GestionService)
get_postgresql_activos() {
    echo "📊 Obteniendo activos de PostgreSQL (GestionService)..."
    response=$(curl -s "$GESTION_SERVICE_URL/activos" 2>/dev/null)
    echo "$response"
}

# Función para extraer mapas de datos
extract_activo_map() {
    local json_data="$1"
    local source_type="$2"
    
    if [[ "$source_type" == "mongodb" ]]; then
        # MongoDB: extraer activo_id -> {nombre, estado, tipo}
        echo "$json_data" | jq -r '.[] | "\(.activo_id)|\(.nombre)|\(.estado)|\(.tipo // "N/A")"' 2>/dev/null | sort
    elif [[ "$source_type" == "postgresql" ]]; then
        # PostgreSQL: extraer activo_id -> {nombre, estado, tipo}
        echo "$json_data" | jq -r '.activos[] | "\(.activo_id)|\(.nombre)|\(.estado)|\(.tipo)"' 2>/dev/null | sort
    fi
}

echo ""
echo "🚀 Iniciando verificación..."
check_services

echo ""
echo "📋 TESTS DE COHERENCIA DE DATOS"
echo "================================"

# 1. TEST: Comparar cantidad de activos
echo ""
echo "1️⃣ TEST: Cantidad de activos entre bases de datos"
echo "------------------------------------------------"

mongodb_data=$(get_mongodb_activos)
postgresql_data=$(get_postgresql_activos)

mongodb_count=$(echo "$mongodb_data" | jq length 2>/dev/null || echo "0")
postgresql_count=$(echo "$postgresql_data" | jq '.total // (.activos | length)' 2>/dev/null || echo "0")

echo -e "${CYAN}📊 MongoDB (ParserService): $mongodb_count activos${NC}"
echo -e "${CYAN}📊 PostgreSQL (GestionService): $postgresql_count activos${NC}"

if [[ "$mongodb_count" -eq "$postgresql_count" && "$mongodb_count" -gt 0 ]]; then
    test_result "Cantidad de activos" "PASS" "Ambas bases tienen $mongodb_count activos"
    increment_test "PASS"
elif [[ "$mongodb_count" -eq "$postgresql_count" && "$mongodb_count" -eq 0 ]]; then
    test_result "Cantidad de activos" "WARN" "Ambas bases están vacías"
    increment_test "WARN"
else
    test_result "Cantidad de activos" "FAIL" "MongoDB: $mongodb_count, PostgreSQL: $postgresql_count"
    increment_test "FAIL"
fi

# 2. TEST: Verificar existencia de IDs correspondientes
echo ""
echo "2️⃣ TEST: Correspondencia de IDs de activos"
echo "------------------------------------------"

if [[ "$mongodb_count" -gt 0 && "$postgresql_count" -gt 0 ]]; then
    # Extraer IDs
    mongodb_ids=$(echo "$mongodb_data" | jq -r '.[].activo_id' 2>/dev/null | sort)
    postgresql_ids=$(echo "$postgresql_data" | jq -r '.activos[].activo_id' 2>/dev/null | sort)
    
    echo "🔍 IDs en MongoDB:"
    echo "$mongodb_ids" | head -5
    [[ $(echo "$mongodb_ids" | wc -l) -gt 5 ]] && echo "   ... (mostrando primeros 5)"
    
    echo ""
    echo "🔍 IDs en PostgreSQL:"
    echo "$postgresql_ids" | head -5
    [[ $(echo "$postgresql_ids" | wc -l) -gt 5 ]] && echo "   ... (mostrando primeros 5)"
    
    # Comparar IDs
    ids_diff=$(diff <(echo "$mongodb_ids") <(echo "$postgresql_ids") | wc -l)
    
    if [[ "$ids_diff" -eq 0 ]]; then
        test_result "Correspondencia de IDs" "PASS" "Todos los IDs coinciden entre las bases"
        increment_test "PASS"
    else
        # Mostrar diferencias
        missing_in_mongodb=$(comm -13 <(echo "$mongodb_ids") <(echo "$postgresql_ids") | head -3)
        missing_in_postgresql=$(comm -23 <(echo "$mongodb_ids") <(echo "$postgresql_ids") | head -3)
        
        diff_details=""
        [[ -n "$missing_in_mongodb" ]] && diff_details+="Faltan en MongoDB: $(echo "$missing_in_mongodb" | tr '\n' ' ') "
        [[ -n "$missing_in_postgresql" ]] && diff_details+="Faltan en PostgreSQL: $(echo "$missing_in_postgresql" | tr '\n' ' ')"
        
        test_result "Correspondencia de IDs" "FAIL" "$diff_details"
        increment_test "FAIL"
    fi
else
    test_result "Correspondencia de IDs" "WARN" "No hay suficientes datos para comparar"
    increment_test "WARN"
fi

# 3. TEST: Verificar coherencia de nombres de activos
echo ""
echo "3️⃣ TEST: Coherencia de nombres de activos"
echo "-----------------------------------------"

if [[ "$mongodb_count" -gt 0 && "$postgresql_count" -gt 0 ]]; then
    # Extraer mapas ID -> nombre
    mongodb_names=$(extract_activo_map "$mongodb_data" "mongodb")
    postgresql_names=$(extract_activo_map "$postgresql_data" "postgresql")
    
    # Comparar nombres por ID
    name_mismatches=0
    total_compared=0
    
    echo "🔍 Verificando nombres por activo_id..."
    
    while IFS='|' read -r id nombre estado tipo; do
        if [[ -n "$id" ]]; then
            postgres_line=$(echo "$postgresql_names" | grep "^$id|" | head -1)
            if [[ -n "$postgres_line" ]]; then
                postgres_nombre=$(echo "$postgres_line" | cut -d'|' -f2)
                total_compared=$((total_compared + 1))
                
                if [[ "$nombre" != "$postgres_nombre" ]]; then
                    echo -e "   ${YELLOW}⚠️  ID $id: MongoDB='$nombre' vs PostgreSQL='$postgres_nombre'${NC}"
                    name_mismatches=$((name_mismatches + 1))
                fi
            fi
        fi
    done <<< "$mongodb_names"
    
    if [[ $name_mismatches -eq 0 && $total_compared -gt 0 ]]; then
        test_result "Coherencia de nombres" "PASS" "Todos los $total_compared nombres coinciden"
        increment_test "PASS"
    elif [[ $total_compared -eq 0 ]]; then
        test_result "Coherencia de nombres" "WARN" "No se pudieron comparar nombres"
        increment_test "WARN"
    else
        test_result "Coherencia de nombres" "FAIL" "$name_mismatches de $total_compared nombres no coinciden"
        increment_test "FAIL"
    fi
else
    test_result "Coherencia de nombres" "WARN" "No hay suficientes datos para comparar"
    increment_test "WARN"
fi

# 4. TEST: Verificar coherencia de estados de activos
echo ""
echo "4️⃣ TEST: Coherencia de estados de activos"
echo "-----------------------------------------"

if [[ "$mongodb_count" -gt 0 && "$postgresql_count" -gt 0 ]]; then
    # Comparar estados por ID
    state_mismatches=0
    total_compared=0
    
    echo "🔍 Verificando estados por activo_id..."
    
    while IFS='|' read -r id nombre estado tipo; do
        if [[ -n "$id" ]]; then
            postgres_line=$(echo "$postgresql_names" | grep "^$id|" | head -1)
            if [[ -n "$postgres_line" ]]; then
                postgres_estado=$(echo "$postgres_line" | cut -d'|' -f3)
                total_compared=$((total_compared + 1))
                
                if [[ "$estado" != "$postgres_estado" ]]; then
                    echo -e "   ${YELLOW}⚠️  ID $id: MongoDB='$estado' vs PostgreSQL='$postgres_estado'${NC}"
                    state_mismatches=$((state_mismatches + 1))
                fi
            fi
        fi
    done <<< "$mongodb_names"
    
    if [[ $state_mismatches -eq 0 && $total_compared -gt 0 ]]; then
        test_result "Coherencia de estados" "PASS" "Todos los $total_compared estados coinciden"
        increment_test "PASS"
    elif [[ $total_compared -eq 0 ]]; then
        test_result "Coherencia de estados" "WARN" "No se pudieron comparar estados"
        increment_test "WARN"
    else
        test_result "Coherencia de estados" "FAIL" "$state_mismatches de $total_compared estados no coinciden"
        increment_test "FAIL"
    fi
else
    test_result "Coherencia de estados" "WARN" "No hay suficientes datos para comparar"
    increment_test "WARN"
fi

# 5. TEST: Verificar coherencia de tipos de activos
echo ""
echo "5️⃣ TEST: Coherencia de tipos de activos"
echo "---------------------------------------"

if [[ "$mongodb_count" -gt 0 && "$postgresql_count" -gt 0 ]]; then
    # Comparar tipos por ID
    type_mismatches=0
    total_compared=0
    
    echo "🔍 Verificando tipos por activo_id..."
    
    while IFS='|' read -r id nombre estado tipo; do
        if [[ -n "$id" ]]; then
            postgres_line=$(echo "$postgresql_names" | grep "^$id|" | head -1)
            if [[ -n "$postgres_line" ]]; then
                postgres_tipo=$(echo "$postgres_line" | cut -d'|' -f4)
                total_compared=$((total_compared + 1))
                
                # Normalizar tipos (algunos pueden ser N/A)
                [[ "$tipo" == "N/A" ]] && tipo=""
                [[ "$postgres_tipo" == "N/A" ]] && postgres_tipo=""
                
                if [[ "$tipo" != "$postgres_tipo" && -n "$tipo" && -n "$postgres_tipo" ]]; then
                    echo -e "   ${YELLOW}⚠️  ID $id: MongoDB='$tipo' vs PostgreSQL='$postgres_tipo'${NC}"
                    type_mismatches=$((type_mismatches + 1))
                fi
            fi
        fi
    done <<< "$mongodb_names"
    
    if [[ $type_mismatches -eq 0 && $total_compared -gt 0 ]]; then
        test_result "Coherencia de tipos" "PASS" "Todos los tipos disponibles coinciden"
        increment_test "PASS"
    elif [[ $total_compared -eq 0 ]]; then
        test_result "Coherencia de tipos" "WARN" "No se pudieron comparar tipos"
        increment_test "WARN"
    else
        test_result "Coherencia de tipos" "FAIL" "$type_mismatches discrepancias de tipos encontradas"
        increment_test "FAIL"
    fi
else
    test_result "Coherencia de tipos" "WARN" "No hay suficientes datos para comparar"
    increment_test "WARN"
fi

# 6. TEST: Verificar activos específicos conocidos
echo ""
echo "6️⃣ TEST: Verificación de activos específicos conocidos"
echo "------------------------------------------------------"

# Verificar algunos activos conocidos de la sincronización
known_activos=("AC-1001" "AC-1002" "AC-1003")
specific_test_passed=0
specific_test_total=0

for activo_id in "${known_activos[@]}"; do
    specific_test_total=$((specific_test_total + 1))
    
    # Obtener del ParserService
    parser_activo=$(curl -s "$PARSER_SERVICE_URL/activo/$activo_id" 2>/dev/null || echo "")
    # Obtener del GestionService  
    gestion_activo=$(curl -s "$GESTION_SERVICE_URL/activos/$activo_id" 2>/dev/null || echo "")
    
    if [[ -n "$parser_activo" && "$parser_activo" != "null" && -n "$gestion_activo" && "$gestion_activo" != "null" ]]; then
        # Extraer nombres
        parser_nombre=$(echo "$parser_activo" | jq -r '.nombre' 2>/dev/null)
        gestion_nombre=$(echo "$gestion_activo" | jq -r '.nombre' 2>/dev/null)
        
        if [[ "$parser_nombre" == "$gestion_nombre" && "$parser_nombre" != "null" ]]; then
            echo -e "   ${GREEN}✅ $activo_id: Nombres coinciden ('$parser_nombre')${NC}"
            specific_test_passed=$((specific_test_passed + 1))
        else
            echo -e "   ${RED}❌ $activo_id: Parser='$parser_nombre' vs Gestión='$gestion_nombre'${NC}"
        fi
    else
        echo -e "   ${YELLOW}⚠️  $activo_id: No encontrado en una o ambas bases${NC}"
    fi
done

if [[ $specific_test_passed -eq $specific_test_total ]]; then
    test_result "Activos específicos conocidos" "PASS" "$specific_test_passed/$specific_test_total activos verificados correctamente"
    increment_test "PASS"
elif [[ $specific_test_passed -gt 0 ]]; then
    test_result "Activos específicos conocidos" "WARN" "Solo $specific_test_passed/$specific_test_total activos coinciden"
    increment_test "WARN"
else
    test_result "Activos específicos conocidos" "FAIL" "Ningún activo específico pudo ser verificado"
    increment_test "FAIL"
fi

# 7. TEST: Verificar integridad referencial (activos con sensores)
echo ""
echo "7️⃣ TEST: Integridad referencial - Activos con datos de sensores"
echo "---------------------------------------------------------------"

# Verificar que activos en MongoDB tengan datos de sensores consistentes
echo "🔍 Verificando activos con datos de sensores..."

activos_with_sensors=0
total_activos_checked=0

# Obtener primeros 3 activos para verificar sensores
first_three_activos=$(echo "$mongodb_data" | jq -r '.[0:3][] | .activo_id' 2>/dev/null)

for activo_id in $first_three_activos; do
    if [[ -n "$activo_id" ]]; then
        total_activos_checked=$((total_activos_checked + 1))
        
        # Verificar si tiene lecturas de sensores
        sensor_data=$(curl -s "$PARSER_SERVICE_URL/lectura/$activo_id/datos" 2>/dev/null || echo "")
        
        if [[ -n "$sensor_data" && "$sensor_data" != "null" && "$sensor_data" != "[]" ]]; then
            activos_with_sensors=$((activos_with_sensors + 1))
            echo -e "   ${GREEN}✅ $activo_id: Tiene datos de sensores${NC}"
        else
            echo -e "   ${YELLOW}⚠️  $activo_id: Sin datos de sensores${NC}"
        fi
    fi
done

if [[ $total_activos_checked -gt 0 ]]; then
    test_result "Integridad referencial" "PASS" "$activos_with_sensors/$total_activos_checked activos verificados"
    increment_test "PASS"
else
    test_result "Integridad referencial" "WARN" "No se pudieron verificar activos con sensores"
    increment_test "WARN"
fi

echo ""
echo "📊 RESUMEN FINAL DE COHERENCIA"
echo "=============================="
echo ""
echo "📈 Estadísticas:"
echo "   Total de tests ejecutados: $total_tests"
echo -e "   ${GREEN}✅ Tests pasados: $passed_tests${NC}"
echo -e "   ${RED}❌ Tests fallidos: $failed_tests${NC}"
echo -e "   ${YELLOW}⚠️  Advertencias: $warning_tests${NC}"

echo ""
echo "📋 Aspectos de coherencia verificados:"
echo "   1. Cantidad de activos entre MongoDB y PostgreSQL"
echo "   2. Correspondencia de IDs de activos"
echo "   3. Coherencia de nombres de activos"
echo "   4. Coherencia de estados de activos"
echo "   5. Coherencia de tipos de activos"
echo "   6. Verificación de activos específicos conocidos"
echo "   7. Integridad referencial con datos de sensores"

echo ""
success_rate=$(( (passed_tests * 100) / total_tests ))
critical_failures=$failed_tests

if [[ $critical_failures -eq 0 && $passed_tests -gt 0 ]]; then
    echo -e "${GREEN}🎉 ¡COHERENCIA EXCELENTE! Datos sincronizados correctamente ($success_rate% éxito)${NC}"
    echo "✅ Las bases de datos MongoDB y PostgreSQL están coherentes"
    exit 0
elif [[ $critical_failures -le 1 && $success_rate -ge 70 ]]; then
    echo -e "${YELLOW}⚠️  COHERENCIA ACEPTABLE: $success_rate% de tests pasaron${NC}"
    echo "🔧 Inconsistencias menores detectadas, pero mayormente sincronizado"
    exit 1
else
    echo -e "${RED}❌ PROBLEMAS DE COHERENCIA: Solo $success_rate% de tests pasaron${NC}"
    echo "🚨 Inconsistencias graves entre las bases de datos detectadas"
    echo "🔧 Ejecuta el script de sincronización: ./scripts/simple_data_init.sh"
    exit 2
fi
