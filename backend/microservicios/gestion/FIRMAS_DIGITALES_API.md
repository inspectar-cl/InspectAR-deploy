# Sistema de Firmas Digitales - Documentación de API

## Descripción
Sistema completo para gestión de firmas digitales en reportes PDF del microservicio de gestión de InspectAR.

## Funcionalidades Implementadas

### ✅ 1. Tabla de Base de Datos
- **Tabla:** `firmas_digitales`
- **Ubicación:** `backend/database/gestion/docker-entrypoint-initdb.d/001_create_tables.sql`
- **Campos:**
  - `id`: ID único de la firma
  - `usuario_id`: ID del usuario propietario
  - `nombre_archivo`: Nombre descriptivo
  - `ruta_archivo`: Ruta física del archivo
  - `tipo_mime`: Tipo MIME (image/png, image/jpeg, image/svg+xml)
  - `formato`: Formato del archivo (png, jpeg, jpg, svg)
  - `datos_firma`: Datos binarios (BYTEA) para SVG
  - `tamano_bytes`: Tamaño del archivo
  - `es_predeterminada`: Booleano para firma por defecto
  - `creado_en`, `actualizado_en`: Timestamps

### ✅ 2. Endpoints Disponibles

#### 2.1 Subir Firma (Archivo de Imagen)
```bash
POST http://localhost:8092/firmas/upload

# Ejemplo con curl (form-data)
curl -X POST http://localhost:8092/firmas/upload \
  -F "email=usuario@example.com" \
  -F "nombre_archivo=Mi Firma Profesional" \
  -F "es_predeterminada=true" \
  -F "archivo=@/path/to/firma.png"
```

**Formatos aceptados:** PNG, JPEG, JPG, SVG

#### 2.2 Crear Firma desde SVG (Pizarra Digital)
```bash
POST http://localhost:8092/firmas/svg
Content-Type: application/json

{
  "email": "usuario@example.com",
  "nombre_archivo": "Firma Dibujada",
  "datos_svg": "<svg>...</svg>",
  "es_predeterminada": false
}
```

#### 2.3 Obtener Firma por ID
```bash
GET http://localhost:8092/firmas/:id

# Ejemplo
curl http://localhost:8092/firmas/1
```

#### 2.4 Obtener Imagen de la Firma
```bash
GET http://localhost:8092/firmas/:id/imagen

# Ejemplo - devuelve la imagen con el Content-Type correcto
curl http://localhost:8092/firmas/1/imagen --output firma.png
```

#### 2.5 Listar Firmas de un Usuario
```bash
GET http://localhost:8092/firmas/usuario/:email

# Ejemplo
curl http://localhost:8092/firmas/usuario/usuario@example.com | jq '.'
```

#### 2.6 Obtener Firma Predeterminada de un Usuario
```bash
GET http://localhost:8092/firmas/usuario/:email/predeterminada

# Ejemplo
curl http://localhost:8092/firmas/usuario/usuario@example.com/predeterminada | jq '.'
```

#### 2.7 Actualizar Firma
```bash
PUT http://localhost:8092/firmas/:id
Content-Type: application/json

{
  "nombre_archivo": "Nuevo Nombre",
  "es_predeterminada": true
}

# Ejemplo con curl
curl -X PUT http://localhost:8092/firmas/1 \
  -H "Content-Type: application/json" \
  -d '{"nombre_archivo":"Firma Principal","es_predeterminada":true}'
```

#### 2.8 Eliminar Firma
```bash
DELETE http://localhost:8092/firmas/:id

# Ejemplo
curl -X DELETE http://localhost:8092/firmas/1
```

#### 2.9 Establecer Firma como Predeterminada
```bash
POST http://localhost:8092/firmas/:id/predeterminada
Content-Type: application/json

{
  "email": "usuario@example.com"
}

# Ejemplo
curl -X POST http://localhost:8092/firmas/2/predeterminada \
  -H "Content-Type: application/json" \
  -d '{"email":"usuario@example.com"}'
```

### ✅ 3. Integración con Reportes PDF

#### 3.1 Generar Reporte con Firma Específica
```bash
POST http://localhost:8092/reportes/activo/:activo_id
Content-Type: application/json

{
  "campos": ["ubicacion", "historial_mantenimientos", "ultima_acciones"],
  "firma_id": 1
}

# Ejemplo
curl -X POST http://localhost:8092/reportes/activo/1 \
  -H "Content-Type: application/json" \
  -d '{
    "campos": ["ubicacion", "historial_mantenimientos"],
    "firma_id": 1
  }' \
  --output reporte_firmado.pdf
```

#### 3.2 Generar Reporte con Firma Predeterminada del Usuario
```bash
POST http://localhost:8092/reportes/activo/:activo_id
Content-Type: application/json

{
  "campos": ["ubicacion", "historial_mantenimientos", "ultima_acciones"],
  "usar_firma_predeterminada": true,
  "email": "usuario@example.com"
}

# Ejemplo
curl -X POST http://localhost:8092/reportes/activo/1 \
  -H "Content-Type: application/json" \
  -d '{
    "campos": ["ubicacion", "historial_mantenimientos"],
    "usar_firma_predeterminada": true,
    "email": "usuario@example.com"
  }' \
  --output reporte_firmado.pdf
```

### ✅ 4. Flujo de Uso Completo

#### Escenario 1: Usuario sube firma por primera vez
```bash
# 1. Subir firma y marcarla como predeterminada
curl -X POST http://localhost:8092/firmas/upload \
  -F "email=usuario@example.com" \
  -F "nombre_archivo=Firma Oficial" \
  -F "es_predeterminada=true" \
  -F "archivo=@mi_firma.png"

# 2. Generar reporte usando firma predeterminada
curl -X POST http://localhost:8092/reportes/activo/1 \
  -H "Content-Type: application/json" \
  -d '{
    "campos": ["ubicacion", "historial_mantenimientos"],
    "usar_firma_predeterminada": true,
    "email": "usuario@example.com"
  }' \
  --output reporte_firmado.pdf
```

#### Escenario 2: Usuario dibuja firma en pizarra digital
```bash
# 1. Crear firma desde SVG capturado de la pizarra
curl -X POST http://localhost:8092/firmas/svg \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@example.com",
    "nombre_archivo": "Firma Digital Dibujada",
    "datos_svg": "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"300\" height=\"150\"><path d=\"M10 100 Q 50 50 100 100\" stroke=\"black\" fill=\"none\"/></svg>",
    "es_predeterminada": true
  }'

# 2. Generar reporte
curl -X POST http://localhost:8092/reportes/activo/1 \
  -H "Content-Type: application/json" \
  -d '{
    "campos": ["ubicacion"],
    "usar_firma_predeterminada": true,
    "email": "usuario@example.com"
  }' \
  --output reporte_firmado.pdf
```

#### Escenario 3: Usuario gestiona múltiples firmas
```bash
# 1. Listar firmas actuales
curl http://localhost:8092/firmas/usuario/usuario@example.com | jq '.'

# 2. Subir segunda firma
curl -X POST http://localhost:8092/firmas/upload \
  -F "email=usuario@example.com" \
  -F "nombre_archivo=Firma Secundaria" \
  -F "es_predeterminada=false" \
  -F "archivo=@firma2.png"

# 3. Cambiar firma predeterminada
curl -X POST http://localhost:8092/firmas/2/predeterminada \
  -H "Content-Type: application/json" \
  -d '{"email":"usuario@example.com"}'

# 4. Eliminar firma antigua
curl -X DELETE http://localhost:8092/firmas/1
```

## Características Técnicas

### Almacenamiento de Archivos
- **Ubicación:** `storage/firmas/`
- **Formato de nombre:** `{usuario_id}_{timestamp}.{extension}`
- **Ejemplo:** `1_1728648000.png`

### Validación de Formatos
- PNG: `image/png`
- JPEG: `image/jpeg`, `image/jpg`
- SVG: `image/svg+xml`

### Firma en PDF
- La firma se incrusta en el PDF usando base64
- Se muestra en una sección dedicada antes del footer
- Incluye fecha de la firma
- Opcionalmente muestra el nombre del firmante

### Seguridad
- Cada firma está asociada a un usuario específico
- Solo se puede tener una firma predeterminada por usuario
- Los archivos se almacenan con permisos restringidos
- Se validan tipos MIME en la subida

## Estructura de Archivos Modificados/Creados

```
backend/
├── database/gestion/docker-entrypoint-initdb.d/
│   └── 001_create_tables.sql (MODIFICADO - agregada tabla firmas_digitales)
├── microservicios/gestion/
│   ├── Dockerfile (MODIFICADO - copia directorio storage)
│   ├── internal/
│   │   ├── models/
│   │   │   └── models.go (MODIFICADO - agregados modelos de firma)
│   │   ├── repository/
│   │   │   └── firma_repository.go (NUEVO)
│   │   ├── services/
│   │   │   ├── firma_service.go (NUEVO)
│   │   │   └── reporte_service.go (MODIFICADO - integración con firmas)
│   │   └── handlers/
│   │       └── firma_handler.go (NUEVO)
│   ├── api/router/
│   │   └── router.go (MODIFICADO - agregadas rutas de firmas)
│   ├── cmd/
│   │   └── main.go (MODIFICADO - inicialización de firma handler)
│   ├── templates/
│   │   └── reporte_activo.html (MODIFICADO - sección de firma)
│   └── storage/firmas/
│       └── .gitkeep (NUEVO)
```

## Testing

### Ejemplo de Prueba Completa
```bash
#!/bin/bash

# 1. Verificar salud del servicio
curl http://localhost:8092/health

# 2. Crear usuario de prueba (si no existe)
# Nota: Esto depende de tu sistema de autenticación

# 3. Subir firma
FIRMA_RESPONSE=$(curl -X POST http://localhost:8092/firmas/upload \
  -F "email=usuario@example.com" \
  -F "nombre_archivo=Test Firma" \
  -F "es_predeterminada=true" \
  -F "archivo=@test_firma.png")

FIRMA_ID=$(echo $FIRMA_RESPONSE | jq -r '.firma.id')
echo "Firma creada con ID: $FIRMA_ID"

# 4. Verificar firma
curl http://localhost:8092/firmas/$FIRMA_ID | jq '.'

# 5. Generar reporte con firma
curl -X POST http://localhost:8092/reportes/activo/1 \
  -H "Content-Type: application/json" \
  -d "{
    \"campos\": [\"ubicacion\", \"historial_mantenimientos\"],
    \"firma_id\": $FIRMA_ID
  }" \
  --output test_reporte_firmado.pdf

echo "Reporte generado: test_reporte_firmado.pdf"

# 6. Verificar PDF
file test_reporte_firmado.pdf
```

## Notas Importantes

1. **Prerequisitos:** Asegúrate de tener usuarios creados en la tabla `usuarios` antes de subir firmas.

2. **Tamaño de Archivos:** No hay límite explícito, pero se recomienda mantener firmas bajo 2MB.

3. **Formato SVG:** Ideal para firmas dibujadas en canvas/pizarra digital por su tamaño reducido.

4. **Persistencia:** Los archivos se guardan en el sistema de archivos. Considera usar un volumen Docker para persistencia.

5. **Docker Volume:** Para persistir firmas entre reinicios de contenedores, agrega en `docker-compose.yml`:
   ```yaml
   services:
     gestion-service:
       volumes:
         - ./microservicios/gestion/storage/firmas:/app/storage/firmas
   ```

## Próximos Pasos Sugeridos

1. **Autenticación:** Integrar con sistema OAuth2 para obtener `usuario_id` automáticamente
2. **Frontend:** Implementar componente React para pizarra de dibujo de firmas
3. **Validación de Tamaño:** Agregar límite de tamaño de archivo
4. **Compresión:** Comprimir imágenes PNG/JPEG antes de guardar
5. **Thumbnails:** Generar miniaturas para vista previa rápida
6. **Auditoría:** Registrar quién y cuándo usó cada firma en reportes
