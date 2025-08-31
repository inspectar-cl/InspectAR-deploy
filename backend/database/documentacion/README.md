# Base de Datos - Microservicio Documentación

Este directorio contiene la definición de la base de datos PostgreSQL para el microservicio de documentación técnica.

## Estructura de la Base de Datos

### Tabla `documentos`
Almacena metadatos de documentos técnicos asociados a activos industriales.

**Campos principales:**
- `id`: Identificador único del documento
- `activo_id`: Referencia al activo (del microservicio gestión)
- `tecnico_id`: Técnico que subió el documento (opcional)
- `nombre`: Nombre del archivo/documento
- `categoria`: Tipo de documento (ficha_tecnica, manual_fabricante, etc.)
- `ruta_archivo`: Ubicación del archivo en MinIO o storage local
- `es_ficha_tecnica`: Marca especial para fichas técnicas oficiales
- `texto_busqueda`: Vector tsvector para búsqueda de texto completo

### Tabla `analisis_ia`
Resultados de análisis realizados por Gemini IA sobre los documentos.

**Campos principales:**
- `documento_id`: Referencia al documento analizado
- `resumen`: Resumen generado por IA
- `puntos_claves`: JSON con información técnica clave
- `graficos_detectados`: JSON con descripciones de gráficos encontrados

## Características de Optimización

### Índices de Rendimiento
- **GIN Index** en `texto_busqueda` para búsqueda full-text súper rápida
- **Índices compuestos** para consultas frecuentes
- **Índices por categoría** para filtros específicos

### Búsqueda de Texto Completo
- **Idioma**: Configurado para español
- **Trigger automático** que actualiza el vector de búsqueda
- **Función `buscar_documentos_texto()`** para búsquedas optimizadas

## Datos de Ejemplo

El script incluye datos de prueba para:
- Ficha técnica de caldera Bosch 500kW
- Manual de operación y mantenimiento
- Reportes de mantenimiento preventivo
- Certificaciones de seguridad
- Documentos de bomba hidráulica
- Análisis de IA completados

## Puerto de Base de Datos

**Puerto asignado**: `5434`

## Uso

```bash
# Ejecutar desde docker-compose
docker-compose up documentacion-db

# Conectar directamente
psql -h localhost -p 5434 -U documentacion_user -d documentacion_db
```

## Consultas de Ejemplo

```sql
-- Buscar documentos por texto
SELECT * FROM buscar_documentos_texto('caldera mantenimiento');

-- Obtener fichas técnicas de un activo
SELECT * FROM documentos 
WHERE activo_id = 1 AND es_ficha_tecnica = true;

-- Documentos con análisis de IA completado
SELECT d.nombre, a.resumen 
FROM documentos d 
JOIN analisis_ia a ON d.id = a.documento_id 
WHERE a.estado = 'completado';
```

## Migración y Mantenimiento

- Los triggers mantienen automáticamente el vector de búsqueda
- La columna `actualizado_en` se actualiza automáticamente
- Los análisis de IA se pueden regenerar eliminando registros de `analisis_ia`
