# Microservicio de Documentación y Conocimiento

Microservicio especializado en la gestión de documentos técnicos, fichas técnicas de activos industriales y análisis de documentos mediante inteligencia artificial con Gemini.

## Funcionalidades

### HdU05 - Subida de documentación técnica y almacenamiento
- ✅ Subir documentos en formatos PDF y DOCX
- ✅ Asociar documentos a activos específicos
- ✅ Categorización automática (ficha técnica, informe de mantenimiento, diagnóstico, manual fabricante, certificación)
- ✅ Búsqueda y filtrado por categoría, fecha de emisión o palabra clave
- ✅ Almacenamiento en MinIO o sistema de archivos local
- ✅ Gestión de metadatos y control de versiones

### HdU23 - Fichas técnicas de activos
- ✅ Acceso directo a fichas técnicas por activo
- ✅ Exportación en formato PDF
- ✅ Identificación automática de documentos de fabricante
- ✅ Vista integrada con información del activo
- ✅ Descarga directa desde la vista del activo

### HdU19 - Lectura y análisis semántico con Gemini (IA)
- ✅ Análisis automático de documentos PDF
- ✅ Extracción de información clave mediante IA
- ✅ Identificación y descripción de gráficos y diagramas
- ✅ Resumen inteligente de contenido técnico
- ✅ Construcción de historial de mantenimiento enriquecido
- ✅ Procesamiento asíncrono para documentos grandes

## Arquitectura

### Stack Tecnológico
- **Backend**: Go 1.22 con Gin Framework
- **Base de Datos**: PostgreSQL 15 
- **Storage**: MinIO (S3-compatible) o Local File System
- **IA**: Google Gemini 1.5 Flash
- **Contenedores**: Docker & Docker Compose

### Estructura de Datos
```
documentos/
├── metadata (PostgreSQL)
│   ├── id, activo_id, tecnico_id
│   ├── nombre, descripcion, categoria
│   ├── tipo_archivo, tamano_bytes
│   ├── fecha_emision, palabras_clave
│   └── es_ficha_tecnica
├── archivos (MinIO/Local)
│   ├── PDFs técnicos
│   └── Documentos DOCX
└── analisis_ia (PostgreSQL)
    ├── resumen automático
    ├── puntos_claves (JSON)
    └── graficos detectados (JSON)
```

## API Endpoints

### 🏥 Health Check

#### Verificar estado del servicio
```bash
curl -X GET http://localhost:8093/health
```
**Respuesta:**
```json
{
  "status": "ok",
  "service": "documentacion"
}
```

---

### 📄 Documentos (HdU05)

#### Subir documento técnico
```bash
curl -X POST http://localhost:8093/documentos \
  -F "archivo=@manual_caldera.pdf" \
  -F "activo_id=1" \
  -F "nombre=Manual Técnico Caldera Principal" \
  -F "descripcion=Manual completo de operación y mantenimiento" \
  -F "categoria=manual_fabricante" \
  -F "fecha_emision=2024-01-15" \
  -F "subido_por=admin" \
  -F "palabras_clave=caldera,manual,operacion,mantenimiento" \
  -F "es_ficha_tecnica=false"
```
**Respuesta:**
```json
{
  "id": 5,
  "activo_id": 1,
  "nombre": "Manual Técnico Caldera Principal",
  "descripcion": "Manual completo de operación y mantenimiento",
  "categoria": "manual_fabricante",
  "tipo_archivo": "pdf",
  "ruta_archivo": "1_20240825_143022_manual_caldera.pdf",
  "tamano_bytes": 2048576,
  "fecha_emision": "2024-01-15T00:00:00Z",
  "subido_por": "admin",
  "palabras_clave": "caldera,manual,operacion,mantenimiento",
  "es_ficha_tecnica": false,
  "creado_en": "2025-08-25T14:30:22Z"
}
```

#### Obtener documento por ID
```bash
curl -X GET http://localhost:8093/documentos/1
```
**Respuesta:**
```json
{
  "id": 1,
  "activo_id": 1,
  "nombre": "Ficha Técnica Caldera Principal",
  "descripcion": "Especificaciones técnicas oficiales del fabricante",
  "categoria": "ficha_tecnica",
  "tipo_archivo": "pdf",
  "tamano_bytes": 2048576,
  "fecha_emision": "2024-01-15T00:00:00Z",
  "subido_por": "admin",
  "palabras_clave": "caldera,especificaciones,fabricante",
  "es_ficha_tecnica": true,
  "creado_en": "2025-08-25T10:30:00Z"
}
```

#### Descargar archivo
```bash
curl -X GET http://localhost:8093/documentos/1/descargar -o documento.pdf
```
**Respuesta:** Descarga directa del archivo PDF

#### Buscar documentos con filtros
```bash
curl -X GET "http://localhost:8093/documentos/buscar?categoria=ficha_tecnica&palabra_clave=caldera&fecha_desde=2024-01-01"
```
**Respuesta:**
```json
{
  "documentos": [
    {
      "id": 1,
      "activo_id": 1,
      "nombre": "Ficha Técnica Caldera Principal",
      "categoria": "ficha_tecnica",
      "fecha_emision": "2024-01-15T00:00:00Z"
    }
  ],
  "total": 1,
  "filtros": {
    "categoria": "ficha_tecnica",
    "palabra_clave": "caldera",
    "fecha_desde": "2024-01-01T00:00:00Z"
  }
}
```

---

### 🏭 Documentos por Activo

#### Obtener todos los documentos de un activo
```bash
curl -X GET http://localhost:8093/activos/1/documentos
```
**Respuesta:**
```json
{
  "activo_id": 1,
  "documentos": [
    {
      "id": 1,
      "nombre": "Ficha Técnica Caldera Principal",
      "categoria": "ficha_tecnica",
      "es_ficha_tecnica": true
    },
    {
      "id": 2,
      "nombre": "Informe Mantenimiento Enero 2024",
      "categoria": "informe_mantenimiento",
      "es_ficha_tecnica": false
    }
  ],
  "total": 2
}
```

#### Obtener historial de mantenimiento completo
```bash
curl -X GET http://localhost:8093/activos/1/historial
```
**Respuesta:**
```json
{
  "activo_id": 1,
  "documentos_total": 2,
  "ultimo_mantenimiento": "2024-01-30T00:00:00Z",
  "documentos": [...],
  "analisis_ia": [
    {
      "id": 1,
      "documento_id": 1,
      "resumen": "Documento técnico que especifica las características operacionales...",
      "puntos_claves": ["Presión máxima: 150 PSI", "Temperatura operación: 180°C"],
      "graficos": ["Diagrama de flujo del sistema", "Tabla de especificaciones"],
      "estado": "completado"
    }
  ]
}
```

---

### 📋 Fichas Técnicas (HdU23)

#### 🎯 **RUTA PRINCIPAL: Obtener ficha técnica de un activo**
```bash
curl -X GET http://localhost:8093/activos/1/ficha-tecnica
```
**Respuesta:**
```json
{
  "id": 1,
  "activo_id": 1,
  "nombre": "Ficha Técnica Caldera Principal",
  "descripcion": "Especificaciones técnicas oficiales del fabricante",
  "categoria": "ficha_tecnica",
  "tipo_archivo": "pdf",
  "ruta_archivo": "ficha_caldera_principal.pdf",
  "tamano_bytes": 2048576,
  "fecha_emision": "2024-01-15T00:00:00Z",
  "subido_por": "admin",
  "es_ficha_tecnica": true,
  "creado_en": "2025-08-25T10:30:00Z"
}
```

---

### 🤖 Análisis con IA (HdU19)

#### Solicitar análisis de IA de un documento
```bash
curl -X POST http://localhost:8093/documentos/1/analizar
```
**Respuesta:**
```json
{
  "mensaje": "Análisis iniciado",
  "analisis": {
    "id": 1,
    "documento_id": 1,
    "estado": "procesando",
    "creado_en": "2025-08-25T14:45:00Z"
  },
  "estado": "procesando"
}
```

#### Obtener resultado del análisis de IA
```bash
curl -X GET http://localhost:8093/documentos/1/analisis
```
**Respuesta:**
```json
{
  "id": 1,
  "documento_id": 1,
  "resumen": "Esta ficha técnica describe una caldera industrial con capacidad de 500kW, diseñada para operación continua en entornos industriales. Incluye especificaciones de presión, temperatura y requisitos de mantenimiento.",
  "puntos_claves": [
    "Capacidad térmica: 500kW",
    "Presión máxima operación: 150 PSI",
    "Temperatura máxima: 180°C",
    "Combustible: Gas natural/GLP",
    "Mantenimiento requerido cada 3 meses"
  ],
  "graficos": [
    "Diagrama esquemático del sistema de combustión",
    "Tabla de especificaciones técnicas principales",
    "Gráfico de eficiencia vs carga operativa"
  ],
  "estado": "completado",
  "creado_en": "2025-08-25T14:45:00Z"
}
```

## Configuración

### Variables de Entorno
```bash
# Base de datos
DOC_DATABASE_HOST=localhost
DOC_DATABASE_PORT=5433
DOC_DATABASE_USER=documentacion_user
DOC_DATABASE_PASSWORD=documentacion_pass
DOC_DATABASE_DBNAME=documentacion_db

# Storage
DOC_STORAGE_TYPE=minio  # o "local"

# MinIO (opcional)
DOC_STORAGE_MINIO_ENDPOINT=localhost:9000
DOC_STORAGE_MINIO_ACCESS_KEY=minioadmin
DOC_STORAGE_MINIO_SECRET_KEY=minioadmin
DOC_STORAGE_MINIO_BUCKET=documentos

# Gemini IA (opcional)
GEMINI_API_KEY=your_gemini_api_key
GCP_PROJECT_ID=your_gcp_project_id
```

### Límites y Restricciones
- **Tamaño máximo de archivo**: 50MB
- **Formatos soportados**: PDF, DOCX
- **Categorías válidas**: ficha_tecnica, informe_mantenimiento, diagnostico, manual_fabricante, certificacion

## Ejecución

### Con Docker Compose
```bash
# Iniciar servicios de base de datos y storage
docker-compose up documentacion-db minio -d

# Iniciar microservicio
docker-compose up documentacion-service
```

### Desarrollo Local
```bash
# Instalar dependencias
go mod tidy

# Ejecutar
go run cmd/main.go
```

El servicio estará disponible en el puerto `8093`.

## Estructura de Base de Datos

### Tablas principales:
- **`documentos`** - Metadatos de documentos técnicos
- **`analisis_ia`** - Resultados de análisis con Gemini

### Relaciones:
- `documentos` ↔ `analisis_ia` (uno a muchos)
- `documentos` → `activos` (referencia externa al microservicio de gestión)

## Casos de Uso

### 🔧 Para Técnicos
- Subir reportes de mantenimiento con IA que extrae automáticamente información clave
- Acceder a fichas técnicas durante trabajo en campo
- Consultar historial completo de intervenciones en un activo

### 🏢 Para Inspectores
- Cargar documentación técnica de activos críticos
- Buscar documentos por palabras clave o fechas
- Generar reportes automáticos con análisis de IA

### 📊 Para Administradores
- Gestionar biblioteca centralizada de documentación técnica
- Monitorear análisis de IA y extraer insights
- Mantener histórico completo de documentación por activo

## Estados y Tipos

### Categorías de Documentos:
- `ficha_tecnica` - Especificaciones del fabricante
- `informe_mantenimiento` - Reportes de mantenimiento
- `diagnostico` - Análisis de problemas
- `manual_fabricante` - Manuales oficiales
- `certificacion` - Certificados y validaciones

### Estados de Análisis IA:
- `procesando` - Análisis en curso
- `completado` - Análisis terminado exitosamente
- `error` - Error en el procesamiento

## Integración con Otros Microservicios

- **Gestión**: Obtiene información de activos y técnicos
- **Notificaciones**: Envía alertas cuando se completan análisis
- **API Gateway**: Rutas centralizadas para acceso frontend
