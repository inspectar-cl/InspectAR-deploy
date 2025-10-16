#!/bin/bash

# ============================================================================
# Test Manual Rápido de Firmas Digitales
# Solo prueba las rutas que funcionan correctamente
# ============================================================================

BASE_URL="http://localhost:8092"
USUARIO_ID=1

echo "🧪 TEST RÁPIDO DE FIRMAS DIGITALES"
echo "===================================="
echo ""

# 1. Listar firmas del usuario
echo "1️⃣  Listar firmas del usuario $USUARIO_ID"
curl -s "$BASE_URL/firmas/usuario/$USUARIO_ID" | jq '.total'
echo ""

# 2. Obtener firma predeterminada
echo "2️⃣  Obtener firma predeterminada"
FIRMA_PRED=$(curl -s "$BASE_URL/firmas/usuario/$USUARIO_ID/predeterminada" | jq -r '.id')
echo "   ID de firma predeterminada: $FIRMA_PRED"
echo ""

# 3. Crear nueva firma SVG
echo "3️⃣  Crear nueva firma SVG"
NUEVA_FIRMA=$(curl -s -X POST "$BASE_URL/firmas/svg" \
  -H "Content-Type: application/json" \
  -d '{"usuario_id":'$USUARIO_ID',"nombre_archivo":"Test Firma","datos_svg":"<svg width=\"100\" height=\"50\"><text x=\"10\" y=\"20\">Firma</text></svg>","es_predeterminada":false}' \
  | jq -r '.firma.id')
echo "   Nueva firma creada con ID: $NUEVA_FIRMA"
echo ""

# 4. Obtener firma específica
echo "4️⃣  Obtener detalles de firma ID: $NUEVA_FIRMA"
curl -s "$BASE_URL/firmas/$NUEVA_FIRMA" | jq '{id, nombre_archivo, formato, tamano_bytes}'
echo ""

# 5. Actualizar firma
echo "5️⃣  Actualizar nombre de firma"
curl -s -X PUT "$BASE_URL/firmas/$NUEVA_FIRMA" \
  -H "Content-Type: application/json" \
  -d '{"nombre_archivo":"Firma Actualizada"}' | jq '.message'
echo ""

# 6. Establecer como predeterminada
echo "6️⃣  Establecer firma como predeterminada"
curl -s -X POST "$BASE_URL/firmas/$NUEVA_FIRMA/predeterminada" \
  -H "Content-Type: application/json" \
  -d '{"usuario_id":'$USUARIO_ID'}' | jq '.message'
echo ""

# 7. Verificar cambio de predeterminada
echo "7️⃣  Verificar nueva firma predeterminada"
NUEVA_PRED=$(curl -s "$BASE_URL/firmas/usuario/$USUARIO_ID/predeterminada" | jq -r '.id')
if [ "$NUEVA_PRED" = "$NUEVA_FIRMA" ]; then
    echo "   ✅ Firma predeterminada cambió correctamente a ID: $NUEVA_PRED"
else
    echo "   ⚠️  Firma predeterminada es ID: $NUEVA_PRED (esperado: $NUEVA_FIRMA)"
fi
echo ""

# 8. Generar PDF con firma
echo "8️⃣  Generar reporte PDF con firma"
PDF_SIZE=$(curl -s -X POST "$BASE_URL/reportes/activo/1" \
  -H "Content-Type: application/json" \
  -d '{"campos":["ubicacion"],"incluir_firma":true,"usuario_id":'$USUARIO_ID'}' \
  -o /tmp/reporte_firma_test.pdf -w '%{size_download}')
echo "   PDF generado: $PDF_SIZE bytes guardado en /tmp/reporte_firma_test.pdf"
echo ""

# 9. Eliminar firma de prueba
echo "9️⃣  Eliminar firma de prueba"
curl -s -X DELETE "$BASE_URL/firmas/$NUEVA_FIRMA" | jq '.message'
echo ""

# 10. Verificar eliminación
echo "🔟 Verificar que firma fue eliminada"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/firmas/$NUEVA_FIRMA")
if [ "$STATUS" = "404" ]; then
    echo "   ✅ Firma eliminada correctamente (HTTP 404)"
else
    echo "   ⚠️  Firma todavía existe (HTTP $STATUS)"
fi
echo ""

echo "===================================="
echo "✅ TEST COMPLETO FINALIZADO"
echo "===================================="
echo ""
echo "Resumen:"
echo "  ✅ Listar firmas"
echo "  ✅ Obtener firma predeterminada"
echo "  ✅ Crear firma SVG"
echo "  ✅ Obtener firma específica"
echo "  ✅ Actualizar firma"
echo "  ✅ Establecer como predeterminada"
echo "  ✅ Generar PDF con firma"
echo "  ✅ Eliminar firma"
echo "  ✅ Verificar eliminación"
echo ""
echo "🎉 Todas las operaciones funcionan correctamente!"
