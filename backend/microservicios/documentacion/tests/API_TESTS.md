# Colección de Pruebas - Microservicio de Documentación
# =====================================================

Este archivo contiene ejemplos de todas las rutas disponibles para el microservicio de documentación.
Puedes usar estos ejemplos con Postman, curl, o cualquier cliente HTTP.

BASE_URL: http://localhost:8092
API_BASE: http://localhost:8092/api/v1

## 1. HEALTH CHECK

### GET /health
```http
GET http://localhost:8092/health
```

**Respuesta esperada:**
```json
{
  "status": "ok",
  "service": "documentacion",
  "timestamp": "2025-08-29T10:00:00Z"
}
```

---

## 2. GESTIÓN DE DOCUMENTOS

### GET /api/v1/documentos - Listar todos los documentos
```http
GET http://localhost:8092/api/v1/documentos
```

### GET /api/v1/documentos - Con filtros
```http
GET http://localhost:8092/api/v1/documentos?activo_id=1&categoria=manual_fabricante&limit=10&offset=0
```

**Parámetros disponibles:**
- `activo_id`: ID del activo
- `tecnico_id`: ID del técnico
- `categoria`: Categoría del documento
- `solo_fichas_tecnicas`: true/false
- `limit`: Número máximo de resultados
- `offset`: Número de resultados a saltar

**Categorías válidas:**
- `ficha_tecnica`
- `manual_fabricante`
- `reporte_mantenimiento`
- `diagnostico`
- `certificacion`

### GET /api/v1/documentos/{id} - Obtener documento específico
```http
GET http://localhost:8092/api/v1/documentos/1
```

### POST /api/v1/documentos - Crear nuevo documento
```http
POST http://localhost:8092/api/v1/documentos
Content-Type: multipart/form-data

archivo: [archivo PDF/DOCX]
activo_id: 1
nombre: "Manual de Bomba Hidráulica XY-2000"
descripcion: "Manual completo de operación y mantenimiento"
categoria: "manual_fabricante"
tecnico_id: 1
fecha_emision: "2025-01-15"
subido_por: "tecnico@empresa.com"
palabras_clave: "bomba,hidráulica,manual,mantenimiento"
es_ficha_tecnica: false
```

**Campos requeridos:**
- `archivo`: Archivo PDF o DOCX
- `activo_id`: ID del activo (entero)
- `nombre`: Nombre del documento
- `categoria`: Categoría válida
- `subido_por`: Email o nombre del usuario

**Campos opcionales:**
- `descripcion`: Descripción del documento
- `tecnico_id`: ID del técnico
- `fecha_emision`: Fecha en formato YYYY-MM-DD
- `palabras_clave`: Palabras clave separadas por comas
- `es_ficha_tecnica`: true/false

### PUT /api/v1/documentos/{id} - Actualizar documento
```http
PUT http://localhost:8092/api/v1/documentos/1
Content-Type: application/json

{
  "nombre": "Manual Actualizado",
  "descripcion": "Nueva descripción del manual",
  "categoria": "manual_fabricante",
  "palabras_clave": "bomba,hidráulica,actualizado"
}
```

**Campos actualizables:**
- `nombre`
- `descripcion`
- `categoria`
- `palabras_clave`
- `es_ficha_tecnica`

### DELETE /api/v1/documentos/{id} - Eliminar documento
```http
DELETE http://localhost:8092/api/v1/documentos/1
```

### GET /api/v1/documentos/{id}/descargar - Descargar archivo
```http
GET http://localhost:8092/api/v1/documentos/1/descargar
```

---

## 3. BÚSQUEDA DE DOCUMENTOS

### GET /api/v1/documentos/buscar - Búsqueda de texto completo
```http
GET http://localhost:8092/api/v1/documentos/buscar?q=bomba hidráulica
```

### Con filtros adicionales
```http
GET http://localhost:8092/api/v1/documentos/buscar?q=manual&categoria=manual_fabricante&activo_id=1
```

**Parámetros:**
- `q`: Texto a buscar (requerido)
- Todos los filtros de la ruta `/documentos` también aplicables

---

## 4. CONSULTAS INTERACTIVAS CON IA

### POST /api/v1/documentos/{id}/consultar - Realizar consulta
```http
POST http://localhost:8092/api/v1/documentos/1/consultar
Content-Type: application/json

{
  "pregunta": "¿Qué presión de agua tiene esta bomba?"
}
```

**Respuesta esperada:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "documento_id": 1,
    "pregunta": "¿Qué presión de agua tiene esta bomba?",
    "respuesta": "Según las especificaciones técnicas, la bomba opera con una presión máxima de 150 PSI y una presión nominal de 120 PSI.",
    "confianza": 0.92,
    "fuentes": ["Tabla de especificaciones técnicas", "Sección 2.1"],
    "tiempo_respuesta_ms": 1250,
    "creado_en": "2025-08-29T10:30:00Z"
  }
}
```

### Ejemplos de preguntas frecuentes:
```http
POST http://localhost:8092/api/v1/documentos/1/consultar
Content-Type: application/json

{
  "pregunta": "¿Cuándo fue el último mantenimiento?"
}
```

```http
POST http://localhost:8092/api/v1/documentos/1/consultar
Content-Type: application/json

{
  "pregunta": "¿Cuáles son las especificaciones técnicas principales?"
}
```

```http
POST http://localhost:8092/api/v1/documentos/1/consultar
Content-Type: application/json

{
  "pregunta": "¿Qué procedimientos de seguridad debo seguir?"
}
```

```http
POST http://localhost:8092/api/v1/documentos/1/consultar
Content-Type: application/json

{
  "pregunta": "¿Cuál es la frecuencia de mantenimiento recomendada?"
}
```

### GET /api/v1/documentos/{id}/consultas - Historial de consultas
```http
GET http://localhost:8092/api/v1/documentos/1/consultas
```

**Respuesta esperada:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "documento_id": 1,
      "pregunta": "¿Qué presión de agua tiene esta bomba?",
      "respuesta": "La bomba opera con una presión máxima de 150 PSI...",
      "confianza": 0.92,
      "fuentes": ["Tabla de especificaciones técnicas"],
      "tiempo_respuesta_ms": 1200,
      "creado_en": "2025-08-29T09:15:00Z"
    }
  ]
}
```

### GET /api/v1/consultas/estadisticas - Estadísticas generales
```http
GET http://localhost:8092/api/v1/consultas/estadisticas
```

**Respuesta esperada:**
```json
{
  "success": true,
  "data": {
    "total_consultas": 47,
    "consultas_hoy": 12,
    "promedio_confianza": 0.84,
    "tiempo_promedio_ms": 1250,
    "top_preguntas_frecuentes": [
      {
        "pregunta": "¿Qué presión tiene?",
        "frecuencia": 15
      },
      {
        "pregunta": "¿Cuándo fue el mantenimiento?",
        "frecuencia": 8
      }
    ]
  }
}
```

---

## 5. EJEMPLOS DE RESPUESTAS DE ERROR

### 400 Bad Request
```json
{
  "success": false,
  "error": "Datos de entrada inválidos",
  "details": "El campo 'categoria' debe ser uno de: ficha_tecnica, manual_fabricante, reporte_mantenimiento, diagnostico, certificacion"
}
```

### 404 Not Found
```json
{
  "success": false,
  "error": "Documento no encontrado",
  "details": "No existe un documento con ID 999"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "error": "Error interno del servidor",
  "details": "Error procesando la solicitud"
}
```

---

## 6. COMANDOS CURL DE EJEMPLO

### Listar documentos
```bash
curl -X GET "http://localhost:8092/api/v1/documentos" \
  -H "Accept: application/json"
```

### Obtener documento específico
```bash
curl -X GET "http://localhost:8092/api/v1/documentos/1" \
  -H "Accept: application/json"
```

### Crear documento
```bash
curl -X POST "http://localhost:8092/api/v1/documentos" \
  -F "archivo=@documento.pdf" \
  -F "activo_id=1" \
  -F "nombre=Manual de Prueba" \
  -F "categoria=manual_fabricante" \
  -F "subido_por=test@empresa.com"
```

### Actualizar documento
```bash
curl -X PUT "http://localhost:8092/api/v1/documentos/1" \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Documento Actualizado","descripcion":"Nueva descripción"}'
```

### Realizar consulta
```bash
curl -X POST "http://localhost:8092/api/v1/documentos/1/consultar" \
  -H "Content-Type: application/json" \
  -d '{"pregunta":"¿Qué presión de agua tiene esta bomba?"}'
```

### Buscar documentos
```bash
curl -X GET "http://localhost:8092/api/v1/documentos/buscar?q=bomba" \
  -H "Accept: application/json"
```

---

## 7. HEADERS RECOMENDADOS

Para todas las peticiones, es recomendable incluir:

```http
Accept: application/json
Content-Type: application/json (para POST/PUT con JSON)
Content-Type: multipart/form-data (para subida de archivos)
```

---

## 8. CÓDIGOS DE ESTADO ESPERADOS

- **200**: Operación exitosa
- **201**: Recurso creado exitosamente
- **400**: Datos de entrada inválidos
- **404**: Recurso no encontrado
- **500**: Error interno del servidor

---

## 9. NOTAS IMPORTANTES

1. **Subida de archivos**: Solo se permiten archivos PDF y DOCX con un tamaño máximo de 50MB.

2. **Consultas IA**: Las respuestas incluyen un nivel de confianza (0.0-1.0) y referencias a las fuentes del documento.

3. **Búsqueda**: Utiliza búsqueda de texto completo en PostgreSQL con soporte para español.

4. **Paginación**: Usa `limit` y `offset` para controlar la cantidad de resultados.

5. **Filtros**: Combina múltiples filtros para búsquedas más específicas.
