#!/bin/bash

# ejecutarlo con "bash test_tipos_falla.sh"
# Script de prueba integral para tipos_falla y comentarios
# Crea 25 reportes de falla con 3-4 comentarios cada uno
# Luego verifica la paginación y estructura de respuesta

BASE_URL="http://localhost:8092"
EDIFICIO_ID=1

# Colores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Contadores
TOTAL_FALLAS_CREADAS=0
TOTAL_COMENTARIOS_CREADOS=0
ERRORES=0

# Arrays para emails y tipos de falla
EMAILS=("residente@example.com" "tecnico@example.com" "analista@example.com")
TIPOS_FALLA=("falla agua" "falla ascensor" "falla electricidad" "falla caldera")

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}  TEST INTEGRAL DE TIPOS_FALLA Y COMENTARIOS${NC}"
echo -e "${BLUE}================================================${NC}\n"

# ==============================================================================
# FASE 1: CREAR 25 TIPOS DE FALLA
# ==============================================================================
echo -e "${YELLOW}[FASE 1] Creando 25 tipos de falla...${NC}\n"

FALLA_IDS=()

for i in {1..25}; do
    # Seleccionar email alternando entre los 3 emails disponibles
    EMAIL_INDEX=$((($i - 1) % 3))
    EMAIL="${EMAILS[$EMAIL_INDEX]}"
    
    # Seleccionar tipo de falla alternando entre los 4 tipos
    TIPO_INDEX=$((($i - 1) % 4))
    TIPO="${TIPOS_FALLA[$TIPO_INDEX]}"
    
    # Crear el reporte de falla
    echo -e "${BLUE}Creando falla #$i: ${NC}tipo='$TIPO', email='$EMAIL'"
    
    RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/tipos-falla" \
        -H "Content-Type: application/json" \
        -d "{
            \"tipo\": \"$TIPO\",
            \"descripcion\": \"Descripción de prueba para falla #$i del tipo $TIPO\",
            \"email\": \"$EMAIL\",
            \"id_edificio\": $EDIFICIO_ID
        }")
    
    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | head -n-1)
    
    if [ "$HTTP_CODE" -eq 201 ]; then
        # Extraer id_falla del JSON response
        FALLA_ID=$(echo "$BODY" | grep -o '"id_falla":[0-9]*' | grep -o '[0-9]*')
        FALLA_IDS+=("$FALLA_ID")
        TOTAL_FALLAS_CREADAS=$((TOTAL_FALLAS_CREADAS + 1))
        echo -e "${GREEN}✓ Falla creada exitosamente (ID: $FALLA_ID)${NC}\n"
    else
        echo -e "${RED}✗ Error al crear falla: HTTP $HTTP_CODE${NC}"
        echo -e "${RED}Response: $BODY${NC}\n"
        ERRORES=$((ERRORES + 1))
    fi
    
    # Pequeña pausa para evitar sobrecarga
    sleep 0.1
done

echo -e "${GREEN}Total de fallas creadas: $TOTAL_FALLAS_CREADAS/25${NC}\n"

# ==============================================================================
# FASE 2: CREAR COMENTARIOS (3-4 POR CADA FALLA)
# ==============================================================================
echo -e "${YELLOW}[FASE 2] Creando comentarios para cada falla...${NC}\n"

for i in "${!FALLA_IDS[@]}"; do
    FALLA_ID="${FALLA_IDS[$i]}"
    FALLA_NUM=$((i + 1))
    
    # Determinar número de comentarios: 13 primeras fallas = 4 comentarios, resto = 3
    if [ $FALLA_NUM -le 13 ]; then
        NUM_COMENTARIOS=4
    else
        NUM_COMENTARIOS=3
    fi
    
    echo -e "${BLUE}Creando $NUM_COMENTARIOS comentarios para falla #$FALLA_NUM (ID: $FALLA_ID)${NC}"
    
    for j in $(seq 1 $NUM_COMENTARIOS); do
        # Alternar email para los comentarios
        EMAIL_INDEX=$(((FALLA_NUM + j - 1) % 3))
        EMAIL="${EMAILS[$EMAIL_INDEX]}"
        
        RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/comentarios" \
            -H "Content-Type: application/json" \
            -d "{
                \"id_falla\": $FALLA_ID,
                \"email\": \"$EMAIL\",
                \"comentario\": \"Comentario #$j para la falla $FALLA_ID - Autor: $EMAIL\"
            }")
        
        HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
        
        if [ "$HTTP_CODE" -eq 201 ]; then
            TOTAL_COMENTARIOS_CREADOS=$((TOTAL_COMENTARIOS_CREADOS + 1))
            echo -e "${GREEN}  ✓ Comentario #$j creado${NC}"
        else
            echo -e "${RED}  ✗ Error en comentario #$j: HTTP $HTTP_CODE${NC}"
            ERRORES=$((ERRORES + 1))
        fi
    done
    
    echo ""
    sleep 0.1
done

echo -e "${GREEN}Total de comentarios creados: $TOTAL_COMENTARIOS_CREADOS/88${NC}\n"

# ==============================================================================
# FASE 3: VERIFICAR PAGINACIÓN - PÁGINA 1
# ==============================================================================
echo -e "${YELLOW}[FASE 3] Verificando paginación - Página 1 (primeros 10 items)${NC}\n"

RESPONSE=$(curl -s "$BASE_URL/tipos-falla/edificio/$EDIFICIO_ID?pagina=1")

# Contar items en la respuesta
ITEMS_COUNT=$(echo "$RESPONSE" | grep -o '"id_falla"' | wc -l)
echo -e "${BLUE}Items en página 1: $ITEMS_COUNT${NC}"

if [ "$ITEMS_COUNT" -eq 10 ]; then
    echo -e "${GREEN}✓ Página 1 contiene 10 items (correcto)${NC}\n"
else
    echo -e "${RED}✗ Página 1 debería contener 10 items, pero tiene $ITEMS_COUNT${NC}\n"
    ERRORES=$((ERRORES + 1))
fi

# Verificar que incluye username en lugar de id_usuario
HAS_USERNAME=$(echo "$RESPONSE" | grep -o '"username"' | head -n1)
if [ ! -z "$HAS_USERNAME" ]; then
    echo -e "${GREEN}✓ La respuesta incluye campo 'username'${NC}\n"
else
    echo -e "${RED}✗ La respuesta NO incluye campo 'username'${NC}\n"
    ERRORES=$((ERRORES + 1))
fi

# ==============================================================================
# FASE 4: VERIFICAR PAGINACIÓN - PÁGINA 2
# ==============================================================================
echo -e "${YELLOW}[FASE 4] Verificando paginación - Página 2 (items 11-20)${NC}\n"

RESPONSE=$(curl -s "$BASE_URL/tipos-falla/edificio/$EDIFICIO_ID?pagina=2")

ITEMS_COUNT=$(echo "$RESPONSE" | grep -o '"id_falla"' | wc -l)
echo -e "${BLUE}Items en página 2: $ITEMS_COUNT${NC}"

if [ "$ITEMS_COUNT" -eq 10 ]; then
    echo -e "${GREEN}✓ Página 2 contiene 10 items (correcto)${NC}\n"
else
    echo -e "${RED}✗ Página 2 debería contener 10 items, pero tiene $ITEMS_COUNT${NC}\n"
    ERRORES=$((ERRORES + 1))
fi

# ==============================================================================
# FASE 5: VERIFICAR PAGINACIÓN - PÁGINA 3
# ==============================================================================
echo -e "${YELLOW}[FASE 5] Verificando paginación - Página 3 (items 21-25)${NC}\n"

RESPONSE=$(curl -s "$BASE_URL/tipos-falla/edificio/$EDIFICIO_ID?pagina=3")

ITEMS_COUNT=$(echo "$RESPONSE" | grep -o '"id_falla"' | wc -l)
echo -e "${BLUE}Items en página 3: $ITEMS_COUNT${NC}"

if [ "$ITEMS_COUNT" -ge 5 ]; then
    echo -e "${GREEN}✓ Página 3 contiene al menos 5 items (correcto)${NC}\n"
else
    echo -e "${RED}✗ Página 3 debería contener al menos 5 items, pero tiene $ITEMS_COUNT${NC}\n"
    ERRORES=$((ERRORES + 1))
fi

# ==============================================================================
# FASE 6: VERIFICAR ORDENAMIENTO (los más recientes primero)
# ==============================================================================
echo -e "${YELLOW}[FASE 6] Verificando ordenamiento DESC por fecha_publicacion${NC}\n"

# Obtener los primeros 2 IDs de la página 1
RESPONSE=$(curl -s "$BASE_URL/tipos-falla/edificio/$EDIFICIO_ID?pagina=1")

# Los IDs deberían estar en orden descendente (más reciente primero)
FIRST_ID=$(echo "$RESPONSE" | grep -o '"id_falla":[0-9]*' | head -n1 | grep -o '[0-9]*')
SECOND_ID=$(echo "$RESPONSE" | grep -o '"id_falla":[0-9]*' | head -n2 | tail -n1 | grep -o '[0-9]*')

echo -e "${BLUE}Primer ID: $FIRST_ID${NC}"
echo -e "${BLUE}Segundo ID: $SECOND_ID${NC}"

if [ "$FIRST_ID" -gt "$SECOND_ID" ]; then
    echo -e "${GREEN}✓ Ordenamiento correcto: IDs en orden descendente${NC}\n"
else
    echo -e "${RED}✗ Ordenamiento incorrecto: IDs deberían estar en orden descendente${NC}\n"
    ERRORES=$((ERRORES + 1))
fi

# ==============================================================================
# FASE 7: VERIFICAR QUE COMENTARIOS INCLUYEN USERNAME
# ==============================================================================
echo -e "${YELLOW}[FASE 7] Verificando que comentarios incluyen username${NC}\n"

RESPONSE=$(curl -s "$BASE_URL/tipos-falla/edificio/$EDIFICIO_ID?pagina=1")

# Buscar en la sección de comentarios
COMENTARIOS_USERNAME=$(echo "$RESPONSE" | grep -o '"comentarios":\[.*\]' | grep -o '"username"' | head -n1)

if [ ! -z "$COMENTARIOS_USERNAME" ]; then
    echo -e "${GREEN}✓ Los comentarios incluyen campo 'username'${NC}\n"
else
    echo -e "${RED}✗ Los comentarios NO incluyen campo 'username'${NC}\n"
    ERRORES=$((ERRORES + 1))
fi

# ==============================================================================
# FASE 8: OBTENER TODOS LOS TIPOS DE FALLA DEL EDIFICIO (SIN PAGINACIÓN)
# ==============================================================================
echo -e "${YELLOW}[FASE 8] Obteniendo todos los tipos de falla del edificio (sin paginación)${NC}\n"

# Llamar al mismo endpoint pero SIN el parámetro ?pagina para obtener todos los resultados
RESPONSE=$(curl -s "$BASE_URL/tipos-falla/edificio/$EDIFICIO_ID")

# Verificar si la respuesta tiene datos
if [ ! -z "$RESPONSE" ] && [ "$RESPONSE" != "null" ]; then
    TOTAL_ITEMS=$(echo "$RESPONSE" | grep -o '"id_falla"' | wc -l)
    echo -e "${BLUE}Total de items retornados: $TOTAL_ITEMS${NC}"
    
    if [ "$TOTAL_ITEMS" -ge 25 ]; then
        echo -e "${GREEN}✓ Endpoint sin parámetro pagina retorna todos los items creados (al menos 25)${NC}"
    else
        echo -e "${RED}✗ Endpoint sin parámetro pagina retorna solo $TOTAL_ITEMS items (esperados: 25+)${NC}"
        ERRORES=$((ERRORES + 1))
    fi
    
    # Verificar que incluye comentarios
    HAS_COMENTARIOS=$(echo "$RESPONSE" | grep -o '"comentarios"' | head -n1)
    if [ ! -z "$HAS_COMENTARIOS" ]; then
        echo -e "${GREEN}✓ La respuesta incluye comentarios anidados${NC}"
    else
        echo -e "${RED}✗ La respuesta NO incluye comentarios${NC}"
        ERRORES=$((ERRORES + 1))
    fi
    
    # Verificar estructura de respuesta (sin paginación no debe tener "pagina")
    HAS_PAGINA=$(echo "$RESPONSE" | grep -o '"pagina"' | head -n1)
    HAS_TOTAL_ITEMS=$(echo "$RESPONSE" | grep -o '"total_items"' | head -n1)
    HAS_DATA=$(echo "$RESPONSE" | grep -o '"data"' | head -n1)
    
    if [ -z "$HAS_PAGINA" ] && [ ! -z "$HAS_TOTAL_ITEMS" ] && [ ! -z "$HAS_DATA" ]; then
        echo -e "${GREEN}✓ Estructura de respuesta correcta (NO tiene pagina, tiene total_items y data)${NC}"
    else
        echo -e "${RED}✗ Estructura de respuesta incorrecta${NC}"
        ERRORES=$((ERRORES + 1))
    fi
    
    # Contar total de comentarios en la respuesta
    TOTAL_COMENTARIOS_RESPUESTA=$(echo "$RESPONSE" | grep -o '"id_comentario"' | wc -l)
    echo -e "${BLUE}Total de comentarios en respuesta: $TOTAL_COMENTARIOS_RESPUESTA${NC}"
    
    if [ "$TOTAL_COMENTARIOS_RESPUESTA" -ge 80 ]; then
        echo -e "${GREEN}✓ Todos los comentarios están incluidos en la respuesta${NC}\n"
    else
        echo -e "${YELLOW}⚠ Se esperaban al menos 80 comentarios, pero hay $TOTAL_COMENTARIOS_RESPUESTA${NC}\n"
    fi
else
    echo -e "${RED}✗ No se pudo obtener respuesta del endpoint sin paginación${NC}\n"
    ERRORES=$((ERRORES + 1))
fi

# ==============================================================================
# FASE 9: RESUMEN Y ESTADÍSTICAS
# ==============================================================================
echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}  RESUMEN DE RESULTADOS${NC}"
echo -e "${BLUE}================================================${NC}\n"

echo -e "${BLUE}Tipos de falla creados:${NC} $TOTAL_FALLAS_CREADAS/25"
echo -e "${BLUE}Comentarios creados:${NC} $TOTAL_COMENTARIOS_CREADOS/88"
echo -e "${BLUE}Promedio de comentarios por falla:${NC} $(awk "BEGIN {printf \"%.2f\", $TOTAL_COMENTARIOS_CREADOS/$TOTAL_FALLAS_CREADAS}")"
echo -e "${BLUE}Total de errores:${NC} $ERRORES\n"

if [ $ERRORES -eq 0 ]; then
    echo -e "${GREEN}================================================${NC}"
    echo -e "${GREEN}  ✓ TODAS LAS PRUEBAS PASARON EXITOSAMENTE${NC}"
    echo -e "${GREEN}================================================${NC}\n"
    exit 0
else
    echo -e "${RED}================================================${NC}"
    echo -e "${RED}  ✗ ALGUNAS PRUEBAS FALLARON${NC}"
    echo -e "${RED}================================================${NC}\n"
    exit 1
fi
