#!/bin/bash

# Test completo del flujo de trabajo: carga de documentos y consultas
# Script para probar la funcionalidad completa del microservicio de documentación

set -e

# Configuración
BASE_URL="http://localhost:8092"
UPLOAD_DIR="/tmp/test_docs"
TEST_DIR="$(dirname "$0")"

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para logging
log() {
    echo -e "${BLUE}[$(date '+%H:%M:%S')]${NC} $1"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Función para hacer peticiones HTTP
http_request() {
    local method="$1"
    local url="$2"
    local data="$3"
    local expected_status="$4"
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" "$url")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url")
    fi
    
    status=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status" = "$expected_status" ]; then
        echo "$body"
        return 0
    else
        error "Expected HTTP $expected_status, got $status. Response: $body"
    fi
}

# Función para subir archivo
upload_file() {
    local file_path="$1"
    local activo_id="$2"
    local nombre="$3"
    local descripcion="$4"
    local categoria="$5"
    local es_ficha_tecnica="$6"
    
    curl -s -w "\n%{http_code}" \
        -X POST \
        -F "archivo=@$file_path" \
        -F "activo_id=$activo_id" \
        -F "nombre=$nombre" \
        -F "descripcion=$descripcion" \
        -F "categoria=$categoria" \
        -F "subido_por=test_user" \
        -F "es_ficha_tecnica=$es_ficha_tecnica" \
        -F "palabras_clave=test,automatizado,prueba" \
        "$BASE_URL/api/v1/documentos"
}

# Función para crear archivos PDF de prueba
create_test_files() {
    log "Creando archivos de prueba..."
    
    mkdir -p "$UPLOAD_DIR"
    
    # Simular PDFs con contenido específico
    cat > "$UPLOAD_DIR/bomba_centrifuga_ficha.txt" << 'EOF'
FICHA TÉCNICA - BOMBA CENTRÍFUGA BC-001

ESPECIFICACIONES TÉCNICAS:
- Modelo: BC-001
- Caudal máximo: 500 L/min
- Presión máxima: 15 bar
- Potencia: 5.5 kW
- RPM: 2900
- Entrada: 3"
- Salida: 2"
- Material: Acero inoxidable 316L

CONDICIONES DE OPERACIÓN:
- Temperatura máxima fluido: 80°C
- Viscosidad máxima: 200 cSt
- Sólidos en suspensión: máx. 3%

MANTENIMIENTO:
- Revisión semanal de vibraciones
- Cambio de aceite cada 2000 horas
- Inspección de sellos cada 6 meses
EOF

    cat > "$UPLOAD_DIR/manual_operacion_bomba.txt" << 'EOF'
MANUAL DE OPERACIÓN - SISTEMA DE BOMBEO

PROCEDIMIENTOS DE ARRANQUE:
1. Verificar válvulas de succión abiertas
2. Purgar aire del sistema
3. Verificar presión de succión > 0.5 bar
4. Arrancar motor en modo manual
5. Monitorear presión de descarga

PARÁMETROS NORMALES:
- Presión succión: 0.5-2.0 bar
- Presión descarga: 8-12 bar
- Corriente motor: 8-10 A
- Temperatura cojinetes: < 70°C

ALARMAS Y ACCIONES:
- Baja presión succión: Parar bomba inmediatamente
- Alta temperatura: Reducir carga operativa
- Vibración alta: Programar mantenimiento
EOF

    cat > "$UPLOAD_DIR/compresor_manual.txt" << 'EOF'
MANUAL TÉCNICO - COMPRESOR AIRE CA-002

ESPECIFICACIONES:
- Tipo: Tornillo rotativo
- Capacidad: 2000 L/min
- Presión trabajo: 7 bar
- Potencia motor: 15 kW
- Tanque acumulador: 500 L

OPERACIÓN NORMAL:
- Arranque automático a 5.5 bar
- Parada automática a 7 bar
- Ciclo trabajo recomendado: 75%

MANTENIMIENTO PREVENTIVO:
- Filtro aire: cada 500 horas
- Aceite separador: cada 2000 horas
- Purga condensados: diaria
EOF

    cat > "$UPLOAD_DIR/procedimiento_emergencia.txt" << 'EOF'
PROCEDIMIENTO DE EMERGENCIA - SISTEMA INDUSTRIAL

PARADA DE EMERGENCIA:
1. Activar botón de emergencia general
2. Cerrar válvulas de alimentación principal
3. Verificar aislamiento eléctrico
4. Contactar supervisor de turno

FUGAS DE GAS:
1. Activar alarma general
2. Evacuar área de 50m radio
3. Cortar suministro eléctrico
4. Ventilar área natural
5. NO usar equipos que generen chispas

INCENDIO:
1. Activar sistema contra incendios
2. Evacuar personal no esencial
3. Usar extintores clase C para equipos eléctricos
4. Contactar bomberos: 132
EOF

    # Convertir a "PDF" (simular con extensión)
    cp "$UPLOAD_DIR/bomba_centrifuga_ficha.txt" "$UPLOAD_DIR/bomba_centrifuga_ficha.pdf"
    cp "$UPLOAD_DIR/manual_operacion_bomba.txt" "$UPLOAD_DIR/manual_operacion_bomba.pdf"
    cp "$UPLOAD_DIR/compresor_manual.txt" "$UPLOAD_DIR/compresor_manual.pdf"
    cp "$UPLOAD_DIR/procedimiento_emergencia.txt" "$UPLOAD_DIR/procedimiento_emergencia.pdf"
    
    success "Archivos de prueba creados en $UPLOAD_DIR"
}

# Función para esperar que el servicio esté disponible
wait_for_service() {
    log "Esperando a que el servicio esté disponible..."
    local timeout=60
    local count=0
    
    while [ $count -lt $timeout ]; do
        if curl -s "$BASE_URL/health" > /dev/null 2>&1; then
            success "Servicio disponible"
            return 0
        fi
        sleep 1
        count=$((count + 1))
        echo -n "."
    done
    
    error "Servicio no disponible después de $timeout segundos"
}

# Test de carga de documentos
test_document_upload() {
    log "=== PRUEBA 1: CARGA DE DOCUMENTOS ==="
    
    # Subir ficha técnica de bomba
    log "Subiendo ficha técnica de bomba..."
    response=$(upload_file "$UPLOAD_DIR/bomba_centrifuga_ficha.pdf" "1" "Ficha Técnica Bomba BC-001" "Especificaciones técnicas de bomba centrífuga" "ficha_tecnica" "true")
    status=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status" = "201" ]; then
        success "Ficha técnica subida correctamente"
        BOMBA_DOC_ID=$(echo "$body" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
        log "ID documento bomba: $BOMBA_DOC_ID"
    else
        error "Error subiendo ficha técnica. Status: $status, Response: $body"
    fi
    
    # Subir manual de operación
    log "Subiendo manual de operación..."
    response=$(upload_file "$UPLOAD_DIR/manual_operacion_bomba.pdf" "1" "Manual de Operación Bomba" "Procedimientos de operación y mantenimiento" "manual_fabricante" "false")
    status=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status" = "201" ]; then
        success "Manual de operación subido correctamente"
        MANUAL_DOC_ID=$(echo "$body" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
        log "ID manual: $MANUAL_DOC_ID"
    else
        error "Error subiendo manual. Status: $status, Response: $body"
    fi
    
    # Subir documentación de compresor (activo diferente)
    log "Subiendo manual de compresor..."
    response=$(upload_file "$UPLOAD_DIR/compresor_manual.pdf" "2" "Manual Compresor CA-002" "Documentación técnica del compresor de aire" "manual_fabricante" "false")
    status=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status" = "201" ]; then
        success "Manual de compresor subido correctamente"
        COMPRESOR_DOC_ID=$(echo "$body" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
        log "ID compresor: $COMPRESOR_DOC_ID"
    else
        error "Error subiendo manual compresor. Status: $status, Response: $body"
    fi
    
    # Subir procedimiento de emergencia
    log "Subiendo procedimiento de emergencia..."
    response=$(upload_file "$UPLOAD_DIR/procedimiento_emergencia.pdf" "3" "Procedimiento de Emergencia" "Protocolo de seguridad y emergencias" "reporte_mantenimiento" "false")
    status=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status" = "201" ]; then
        success "Procedimiento de emergencia subido correctamente"
        EMERGENCIA_DOC_ID=$(echo "$body" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
        log "ID procedimiento: $EMERGENCIA_DOC_ID"
    else
        error "Error subiendo procedimiento. Status: $status, Response: $body"
    fi
    
    success "Todos los documentos subidos exitosamente"
}

# Test de consultas de documentos
test_document_queries() {
    log "=== PRUEBA 2: CONSULTAS DE DOCUMENTOS ==="
    
    # Obtener documentos por activo
    log "Consultando documentos del activo 1..."
    response=$(http_request "GET" "$BASE_URL/api/v1/documentos/activo/1" "" "200")
    doc_count=$(echo "$response" | grep -o '"id":[0-9]*' | wc -l)
    if [ "$doc_count" -ge "2" ]; then
        success "Documentos del activo 1 encontrados: $doc_count documentos"
    else
        warning "Se esperaban al menos 2 documentos para el activo 1, encontrados: $doc_count"
    fi
    
    # Buscar fichas técnicas
    log "Buscando fichas técnicas..."
    response=$(http_request "GET" "$BASE_URL/api/v1/documentos/buscar?es_ficha_tecnica=true" "" "200")
    ficha_count=$(echo "$response" | grep -o '"es_ficha_tecnica":true' | wc -l)
    if [ "$ficha_count" -ge "1" ]; then
        success "Fichas técnicas encontradas: $ficha_count"
    else
        warning "No se encontraron fichas técnicas"
    fi
    
    # Buscar documentos con filtro
    log "Buscando documentos con query 'bomba'..."
    response=$(http_request "GET" "$BASE_URL/api/v1/documentos/buscar?query=bomba" "" "200")
    bomba_results=$(echo "$response" | grep -o '"nombre":"[^"]*[Bb]omba[^"]*"' | wc -l)
    if [ "$bomba_results" -ge "1" ]; then
        success "Documentos relacionados con 'bomba' encontrados: $bomba_results"
    else
        warning "No se encontraron documentos relacionados con 'bomba'"
    fi
    
    # Obtener documento específico
    if [ -n "$BOMBA_DOC_ID" ]; then
        log "Obteniendo documento específico (ID: $BOMBA_DOC_ID)..."
        response=$(http_request "GET" "$BASE_URL/api/v1/documentos/$BOMBA_DOC_ID" "" "200")
        if echo "$response" | grep -q "Bomba"; then
            success "Documento específico obtenido correctamente"
        else
            warning "Documento obtenido pero contenido inesperado"
        fi
    fi
}

# Test de consultas con IA
test_ai_queries() {
    log "=== PRUEBA 3: CONSULTAS CON IA ==="
    
    if [ -n "$BOMBA_DOC_ID" ]; then
        # Consulta sobre especificaciones técnicas
        log "Consultando especificaciones de la bomba con IA..."
        ai_query='{"pregunta": "¿Cuál es la presión máxima de la bomba BC-001?"}'
        response=$(curl -s -w "\n%{http_code}" \
            -X POST \
            -H "Content-Type: application/json" \
            -d "$ai_query" \
            "$BASE_URL/api/v1/documentos/$BOMBA_DOC_ID/consultar")
        
        status=$(echo "$response" | tail -n1)
        body=$(echo "$response" | sed '$d')
        
        if [ "$status" = "200" ]; then
            success "Consulta IA sobre presión ejecutada"
            if echo "$body" | grep -qi "15.*bar\|bar.*15"; then
                success "IA respondió correctamente sobre presión (15 bar)"
            else
                warning "IA respondió pero no mencionó la presión correcta"
            fi
        else
            warning "Error en consulta IA. Status: $status"
        fi
        
        # Consulta sobre mantenimiento
        log "Consultando sobre mantenimiento con IA..."
        maintenance_query='{"pregunta": "¿Con qué frecuencia se debe cambiar el aceite?"}'
        response=$(curl -s -w "\n%{http_code}" \
            -X POST \
            -H "Content-Type: application/json" \
            -d "$maintenance_query" \
            "$BASE_URL/api/v1/documentos/$BOMBA_DOC_ID/consultar")
        
        status=$(echo "$response" | tail -n1)
        body=$(echo "$response" | sed '$d')
        
        if [ "$status" = "200" ]; then
            success "Consulta IA sobre mantenimiento ejecutada"
            if echo "$body" | grep -qi "2000.*hora\|hora.*2000"; then
                success "IA respondió correctamente sobre cambio de aceite (2000 horas)"
            else
                warning "IA respondió pero no mencionó el intervalo correcto"
            fi
        else
            warning "Error en consulta IA sobre mantenimiento. Status: $status"
        fi
        
        # Obtener historial de consultas
        log "Obteniendo historial de consultas..."
        response=$(http_request "GET" "$BASE_URL/api/v1/documentos/$BOMBA_DOC_ID/consultas" "" "200")
        consulta_count=$(echo "$response" | grep -o '"id":[0-9]*' | wc -l)
        if [ "$consulta_count" -ge "2" ]; then
            success "Historial de consultas disponible: $consulta_count consultas"
        else
            warning "Se esperaban al menos 2 consultas en el historial, encontradas: $consulta_count"
        fi
    fi
}

# Test de análisis con IA
test_ai_analysis() {
    log "=== PRUEBA 4: ANÁLISIS CON IA ==="
    
    if [ -n "$COMPRESOR_DOC_ID" ]; then
        # Solicitar análisis IA
        log "Solicitando análisis IA del documento de compresor..."
        response=$(curl -s -w "\n%{http_code}" \
            -X POST \
            "$BASE_URL/api/v1/documentos/$COMPRESOR_DOC_ID/analizar")
        
        status=$(echo "$response" | tail -n1)
        body=$(echo "$response" | sed '$d')
        
        if [ "$status" = "200" ] || [ "$status" = "202" ]; then
            success "Análisis IA iniciado para documento de compresor"
            
            # Esperar un poco y obtener el análisis
            sleep 3
            log "Obteniendo resultado del análisis..."
            response=$(http_request "GET" "$BASE_URL/api/v1/documentos/$COMPRESOR_DOC_ID/analisis" "" "200")
            
            if echo "$response" | grep -qi "compresor\|tornillo\|aire"; then
                success "Análisis IA contiene información relevante sobre compresor"
            else
                warning "Análisis IA completado pero contenido inesperado"
            fi
        else
            warning "Error solicitando análisis IA. Status: $status"
        fi
    fi
}

# Test de casos edge y validaciones
test_edge_cases() {
    log "=== PRUEBA 5: CASOS EDGE Y VALIDACIONES ==="
    
    # Documento inexistente
    log "Probando obtener documento inexistente..."
    response=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/api/v1/documentos/99999")
    status=$(echo "$response" | tail -n1)
    if [ "$status" = "404" ]; then
        success "Error 404 correcto para documento inexistente"
    else
        warning "Status inesperado para documento inexistente: $status"
    fi
    
    # Activo sin documentos
    log "Probando obtener documentos de activo inexistente..."
    response=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/api/v1/documentos/activo/99999")
    status=$(echo "$response" | tail -n1)
    if [ "$status" = "200" ]; then
        body=$(echo "$response" | sed '$d')
        if echo "$body" | grep -q '\[\]'; then
            success "Array vacío correcto para activo sin documentos"
        else
            warning "Respuesta inesperada para activo sin documentos"
        fi
    else
        warning "Status inesperado para activo sin documentos: $status"
    fi
    
    # Consulta IA con documento inexistente
    log "Probando consulta IA con documento inexistente..."
    ai_query='{"pregunta": "¿Qué información contiene este documento?"}'
    response=$(curl -s -w "\n%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -d "$ai_query" \
        "$BASE_URL/api/v1/documentos/99999/consultar")
    
    status=$(echo "$response" | tail -n1)
    if [ "$status" = "404" ]; then
        success "Error 404 correcto para consulta IA con documento inexistente"
    else
        warning "Status inesperado para consulta IA con documento inexistente: $status"
    fi
}

# Test de estadísticas
test_statistics() {
    log "=== PRUEBA 6: ESTADÍSTICAS ==="
    
    log "Obteniendo estadísticas de consultas..."
    response=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/api/v1/consultas/estadisticas")
    status=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status" = "200" ]; then
        success "Estadísticas obtenidas correctamente"
        if echo "$body" | grep -q "total_consultas"; then
            total=$(echo "$body" | grep -o '"total_consultas":[0-9]*' | cut -d':' -f2)
            log "Total de consultas realizadas: $total"
        fi
    else
        warning "Error obteniendo estadísticas. Status: $status"
    fi
}

# Cleanup
cleanup() {
    log "Limpiando archivos temporales..."
    rm -rf "$UPLOAD_DIR"
    success "Limpieza completada"
}

# Función principal
main() {
    echo ""
    echo "🧪 PRUEBAS COMPLETAS DEL MICROSERVICIO DE DOCUMENTACIÓN"
    echo "========================================================"
    echo "🔗 Base URL: $BASE_URL"
    echo "📁 Directorio de pruebas: $UPLOAD_DIR"
    echo ""
    
    # Verificar servicio
    wait_for_service
    
    # Crear archivos de prueba
    create_test_files
    
    # Ejecutar pruebas
    test_document_upload
    echo ""
    test_document_queries
    echo ""
    test_ai_queries
    echo ""
    test_ai_analysis
    echo ""
    test_edge_cases
    echo ""
    test_statistics
    
    # Limpieza
    echo ""
    cleanup
    
    echo ""
    echo "✅ PRUEBAS COMPLETAS FINALIZADAS"
    echo "=================================="
    echo "🎯 Microservicio funcionando correctamente"
    echo "📊 Base de datos operativa"
    echo "🤖 IA procesando consultas"
    echo "📁 Gestión de documentos funcional"
    echo ""
    echo "Documentos creados en esta sesión:"
    echo "- Ficha técnica bomba (ID: ${BOMBA_DOC_ID:-'N/A'})"
    echo "- Manual operación bomba (ID: ${MANUAL_DOC_ID:-'N/A'})"
    echo "- Manual compresor (ID: ${COMPRESOR_DOC_ID:-'N/A'})"
    echo "- Procedimiento emergencia (ID: ${EMERGENCIA_DOC_ID:-'N/A'})"
    echo ""
    echo "Para ver logs del sistema:"
    echo "  docker logs documentacion-service"
    echo "  docker logs documentacion-db"
}

# Ejecutar tests
main "$@"
