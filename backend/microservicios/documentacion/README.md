# 📄 Microservicio de Documentación Técnica

Microservicio especializado en la gestión de documentos técnicos asociados a activos industriales, con capacidades de análisis de IA utilizando Google Gemini.

## 🚀 Características Principales

### Funcionalidades Implementadas

- **HdU05**: Subida de documentos técnicos asociados a activos
- **HdU23**: Gestión y consulta de fichas técnicas
- **HdU19**: Análisis de documentos con IA Gemini (bajo demanda)

### Capacidades Técnicas

- ✅ Almacenamiento híbrido (MinIO + PostgreSQL)
- ✅ Búsqueda de texto completo optimizada
- ✅ Análisis de IA con Google Gemini
- ✅ **Consultas interactivas** - Preguntas específicas sobre documentos
- ✅ API REST completa con documentación
- ✅ Validación de tipos de archivo
- ✅ Gestión de metadatos enriquecidos

## 🏗️ Arquitectura

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Cliente Web   │────│   API Gateway    │────│  Documentación  │
│                 │    │   (Puerto 3500)  │    │  (Puerto 8093)  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                                        │
                                            ┌───────────┼───────────┐
                                            │           │           │
                                    ┌───────▼───┐ ┌────▼────┐ ┌────▼────┐
                                    │PostgreSQL │ │  MinIO  │ │ Gemini  │
                                    │(Puerto    │ │(Puerto  │ │   IA    │
                                    │ 5434)     │ │ 9000)   │ │         │
                                    └───────────┘ └─────────┘ └─────────┘
```

## 🛠️ Tecnologías Utilizadas

- **Lenguaje**: Go 1.22
- **Framework**: Gin (HTTP Router)
- **Base de Datos**: PostgreSQL 15 con búsqueda full-text
- **Almacenamiento**: MinIO (compatible S3)
- **IA**: Google Gemini API
- **Containerización**: Docker + Docker Compose

## 📁 Estructura del Proyecto

```
microservicios/documentacion/
├── cmd/
│   └── main.go                 # Punto de entrada
├── internal/
│   ├── models/                 # Modelos de datos
│   │   ├── documento.go
│   │   └── analisis_ia.go
│   ├── handlers/               # Controladores HTTP
│   │   ├── documento_handler.go
│   │   ├── analisis_handler.go
│   │   └── health_handler.go
│   ├── services/               # Lógica de negocio
│   │   ├── documento_service.go
│   │   ├── analisis_service.go
│   │   └── gemini_service.go
│   ├── repository/             # Acceso a datos
│   │   ├── documento_repository.go
│   │   └── analisis_repository.go
│   └── storage/                # Almacenamiento de archivos
│       ├── minio_storage.go
│       ├── local_storage.go
│       └── storage_interface.go
├── api/
│   └── router/                 # Configuración de rutas
│       └── router.go
├── config/                     # Configuración
│   ├── config.yaml
│   └── config.go
├── Dockerfile                  # Imagen Docker
├── go.mod                      # Dependencias Go
├── go.sum
└── README.md                   # Este archivo

database/documentacion/         # Base de datos separada
├── init.sql                    # Schema PostgreSQL
├── Dockerfile                  # Imagen de BD
├── entrypoint.sh              # Script de inicialización
└── README.md                   # Documentación de BD
```

## 🚀 Instalación y Uso

### Prerrequisitos

- Docker y Docker Compose
- Go 1.22+ (para desarrollo local)
- Clave API de Google Gemini

### Configuración

1. **Configurar variables de entorno**:
```bash
export GEMINI_API_KEY="your-gemini-api-key"
export GCP_PROJECT_ID="your-gcp-project-id"
```

2. **Iniciar servicios con Docker Compose**:
```bash
# Desde el directorio backend/
docker-compose up documentacion-service documentacion-db minio
```

3. **Verificar servicios**:
```bash
# API Health Check
curl http://localhost:8093/health

# Dashboard MinIO
open http://localhost:9001
```

## 📚 API Endpoints

### Documentos

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/api/v1/documentos` | Listar documentos con filtros |
| `POST` | `/api/v1/documentos` | Subir nuevo documento |
| `GET` | `/api/v1/documentos/{id}` | Obtener documento específico |
| `PUT` | `/api/v1/documentos/{id}` | Actualizar metadatos |
| `DELETE` | `/api/v1/documentos/{id}` | Eliminar documento |
| `GET` | `/api/v1/documentos/{id}/download` | Descargar archivo |

### Búsqueda

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `GET` | `/api/v1/documentos/buscar` | Búsqueda por texto |
| `GET` | `/api/v1/documentos/activo/{activo_id}` | Documentos por activo |
| `GET` | `/api/v1/documentos/activo/{activo_id}/ficha-tecnica` | Ficha técnica del activo |
| `GET` | `/api/v1/documentos?solo_fichas_tecnicas=true` | Solo fichas técnicas |

### Análisis IA

Nota: Los endpoints de análisis IA están temporalmente deshabilitados por requerimiento.

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/api/v1/documentos/{id}/analizar` | Analizar con IA (deshabilitado) |
| `GET` | `/api/v1/documentos/{id}/analisis` | Obtener análisis (deshabilitado) |
| `POST` | `/api/v1/documentos/{id}/consultar` | Hacer pregunta específica sobre el documento |

## 🔧 Ejemplos de Uso

### Subir Documento

```bash
curl -X POST http://localhost:8093/api/v1/documentos \
  -F "archivo=@ficha_tecnica.pdf" \
  -F "activo_id=1" \
  -F "nombre=Ficha Técnica Caldera Principal" \
  -F "categoria=ficha_tecnica" \
  -F "descripcion=Especificaciones técnicas oficiales" \
  -F "palabras_clave=caldera,bosch,500kw"
```

### Buscar Documentos

```bash
# Búsqueda por texto
curl "http://localhost:8093/api/v1/documentos/buscar?q=caldera mantenimiento"

# Documentos de un activo específico
curl "http://localhost:8093/api/v1/documentos/activo/1"

# Solo fichas técnicas
curl "http://localhost:8093/api/v1/documentos/activo/1/ficha-tecnica"
```

### Analizar con IA

```bash
# Solicitar análisis
curl -X POST http://localhost:8093/api/v1/documentos/1/analizar

# Obtener resultados
curl http://localhost:8093/api/v1/documentos/1/analisis
```

### Consulta Interactiva con IA

```bash
# Hacer pregunta específica sobre el documento
curl -X POST http://localhost:8093/api/v1/documentos/1/consultar \
  -H "Content-Type: application/json" \
  -d '{
    "pregunta": "¿Cuál es la presión de agua de esta bomba?"
  }'

# Pregunta sobre mantenimiento
curl -X POST http://localhost:8093/api/v1/documentos/2/consultar \
  -H "Content-Type: application/json" \
  -d '{
    "pregunta": "¿Cuándo fue la última mantención registrada?"
  }'

# Pregunta sobre especificaciones técnicas
curl -X POST http://localhost:8093/api/v1/documentos/1/consultar \
  -H "Content-Type: application/json" \
  -d '{
    "pregunta": "¿Qué capacidad térmica tiene esta caldera y cuál es su eficiencia?"
  }'
```

**Respuesta típica:**
```json
{
  "pregunta": "¿Cuál es la presión de agua de esta bomba?",
  "respuesta": "Según la ficha técnica, la bomba hidráulica Grundfos tiene una presión máxima de operación de 150 PSI (10.3 bar). La presión de trabajo recomendada es de 120-140 PSI para óptimo rendimiento.",
  "documento_id": 1,
  "confianza": 0.95,
  "fuentes": ["Tabla de especificaciones técnicas", "Sección 3.2 - Parámetros operacionales"],
  "procesado_en": "2025-08-29T10:30:00Z"
}
```

## 🗃️ Base de Datos

### Esquema Principal

```sql
-- Tabla de documentos
CREATE TABLE documentos (
    id SERIAL PRIMARY KEY,
    activo_id INTEGER NOT NULL,
    nombre VARCHAR(255) NOT NULL,
    categoria VARCHAR(50) NOT NULL,
    ruta_archivo VARCHAR(500) NOT NULL,
    texto_busqueda tsvector, -- Para búsqueda full-text
    es_ficha_tecnica BOOLEAN DEFAULT FALSE,
    creado_en TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla de análisis IA
CREATE TABLE analisis_ia (
    id SERIAL PRIMARY KEY,
    documento_id INTEGER REFERENCES documentos(id),
    resumen TEXT NOT NULL,
    puntos_claves JSONB,
    estado VARCHAR(20) DEFAULT 'procesando'
);

-- Tabla de consultas interactivas (NUEVA)
CREATE TABLE consultas_ia (
    id SERIAL PRIMARY KEY,
    documento_id INTEGER REFERENCES documentos(id),
    pregunta TEXT NOT NULL,
    respuesta TEXT NOT NULL,
    confianza DECIMAL(3,2),
    fuentes JSONB,
    creado_en TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## 🤖 Consultas Interactivas con IA

### Funcionalidad Nueva: Preguntas Específicas

El endpoint `/api/v1/documentos/{id}/consultar` permite hacer preguntas específicas sobre el contenido de cualquier documento, ideal para:

- **Consultas técnicas**: "¿Cuál es la presión máxima de esta bomba?"
- **Información de mantenimiento**: "¿Cuándo fue la última mantención?"
- **Especificaciones**: "¿Qué capacidad tiene este equipo?"
- **Procedimientos**: "¿Cómo se calibra este instrumento?"

### Ejemplos Prácticos

#### 1. Consultas sobre Especificaciones Técnicas
```bash
curl -X POST http://localhost:8093/api/v1/documentos/1/consultar \
  -H "Content-Type: application/json" \
  -d '{
    "pregunta": "¿Cuál es la capacidad térmica y eficiencia de esta caldera?"
  }'
```

**Respuesta:**
```json
{
  "pregunta": "¿Cuál es la capacidad térmica y eficiencia de esta caldera?",
  "respuesta": "La caldera Bosch tiene una capacidad térmica de 500kW y una eficiencia energética del 92%. Opera con gas natural o GLP y está certificada ISO 9001.",
  "confianza": 0.98,
  "fuentes": ["Tabla de especificaciones principales", "Sección 2.1 - Características técnicas"]
}
```

#### 2. Consultas sobre Mantenimiento
```bash
curl -X POST http://localhost:8093/api/v1/documentos/3/consultar \
  -H "Content-Type: application/json" \
  -d '{
    "pregunta": "¿Cuándo fue la última mantención y qué se hizo?"
  }'
```

**Respuesta:**
```json
{
  "pregunta": "¿Cuándo fue la última mantención y qué se hizo?",
  "respuesta": "La última mantención fue el 30 de enero de 2024. Se realizó limpieza de quemadores, calibración de termostatos, reemplazo de filtros de aire y verificación de sistemas de seguridad. Estado general: Satisfactorio.",
  "confianza": 0.95,
  "fuentes": ["Reporte de mantención enero 2024", "Checklist de inspección"]
}
```

#### 3. Consultas sobre Presiones y Parámetros
```bash
curl -X POST http://localhost:8093/api/v1/documentos/5/consultar \
  -H "Content-Type: application/json" \
  -d '{
    "pregunta": "¿Cuánta presión de agua maneja esta bomba hidráulica?"
  }'
```

**Respuesta:**
```json
{
  "pregunta": "¿Cuánta presión de agua maneja esta bomba hidráulica?",
  "respuesta": "La bomba hidráulica Grundfos de 50HP maneja una presión máxima de operación de 150 PSI. La presión de trabajo recomendada es entre 120-140 PSI para óptimo rendimiento y vida útil del equipo.",
  "confianza": 0.97,
  "fuentes": ["Ficha técnica Grundfos", "Tabla de parámetros operacionales"]
}
```

### Implementación Técnica

#### Flujo de Procesamiento
1. **Recepción**: El endpoint recibe la pregunta en JSON
2. **Validación**: Verifica que el documento existe y es accesible
3. **Optimización**: Busca consultas similares previas (cache inteligente)
4. **Procesamiento IA**: Analiza pregunta + contexto del documento
5. **Respuesta**: Procesa y estructura la respuesta con metadatos
6. **Persistencia**: Guarda la consulta para optimizar futuras preguntas

#### Prompt Optimizado para IA
El sistema construye prompts específicos basados en:
- **Contexto del documento**: Metadatos y contenido extraído
- **Tipo de pregunta**: Técnica, mantenimiento, especificaciones
- **Categoría del documento**: Ficha técnica, manual, reporte
- **Palabras clave**: Términos relevantes del documento

#### Cache Inteligente
- Detecta preguntas similares con > 80% de coincidencia
- Reutiliza respuestas para optimizar tiempo de respuesta
- Header `X-Consulta-Cache: similar` indica respuesta desde cache

### Casos de Uso Reales

#### Consulta sobre Bomba Hidráulica
```bash
curl -X POST http://localhost:8092/api/v1/documentos/5/consultar \
  -H "Content-Type: application/json" \
  -d '{
    "pregunta": "¿Cuánta presión de agua maneja esta bomba hidráulica Grundfos?"
  }'
```

**Respuesta esperada:**
```json
{
  "pregunta": "¿Cuánta presión de agua maneja esta bomba hidráulica Grundfos?",
  "respuesta": "La bomba hidráulica Grundfos de 50HP maneja una presión máxima de operación de 150 PSI (10.3 bar). La presión de trabajo recomendada está entre 120-140 PSI para óptimo rendimiento y vida útil del equipo.",
  "documento_id": 5,
  "confianza": 0.97,
  "fuentes": ["Ficha técnica Grundfos", "Tabla de parámetros operacionales"],
  "tiempo_respuesta_ms": 1250,
  "procesado_en": "2025-08-29T10:30:00Z"
}
```

#### Consulta sobre Mantenimiento
```bash
curl -X POST http://localhost:8092/api/v1/documentos/3/consultar \
  -H "Content-Type: application/json" \
  -d '{
    "pregunta": "¿Cuándo fue la última mantención de esta caldera y qué trabajos se realizaron?"
  }'
```

**Respuesta esperada:**
```json
{
  "pregunta": "¿Cuándo fue la última mantención de esta caldera y qué trabajos se realizaron?",
  "respuesta": "La última mantención preventiva fue realizada el 30 de enero de 2024. Se ejecutaron las siguientes tareas: limpieza de quemadores, calibración de termostatos, reemplazo de filtros de aire y verificación de sistemas de seguridad. El estado general del equipo fue evaluado como satisfactorio con próxima inspección programada para abril 2024.",
  "documento_id": 3,
  "confianza": 0.95,
  "fuentes": ["Reporte de mantención enero 2024", "Checklist de actividades realizadas"],
  "tiempo_respuesta_ms": 1450,
  "procesado_en": "2025-08-29T10:31:00Z"
}
```

#### Consulta sobre Especificaciones
```bash
curl -X POST http://localhost:8092/api/v1/documentos/1/consultar \
  -H "Content-Type: application/json" \
  -d '{
    "pregunta": "¿Cuál es la capacidad térmica y eficiencia de esta caldera Bosch?"
  }'
```

### Endpoints Adicionales

#### Historial de Consultas
```bash
# Obtener últimas 10 consultas de un documento
curl "http://localhost:8093/api/v1/documentos/1/consultas?limit=10"

# Paginación
curl "http://localhost:8093/api/v1/documentos/1/consultas?limit=5&offset=10"
```

#### Estadísticas Globales
```bash
curl "http://localhost:8093/api/v1/consultas/estadisticas"
```

**Respuesta de estadísticas:**
```json
{
  "total_consultas": 156,
  "confianza_promedio": 0.89,
  "tiempo_promedio_ms": 1350,
  "consultas_hoy": 23,
  "documentos_mas_consultados": [
    {
      "documento_id": 1,
      "nombre_documento": "Ficha Técnica Caldera Bosch",
      "total_consultas": 45
    },
    {
      "documento_id": 5,
      "nombre_documento": "Bomba Hidráulica Grundfos",
      "total_consultas": 32
    }
  ]
}
```

### Optimizaciones

- **Índice GIN** para búsqueda de texto completo
- **Índices compuestos** para consultas frecuentes
- **Configuración en español** para stemming
- **Triggers automáticos** para actualización de vectores

## 🔒 Seguridad

- Validación estricta de tipos de archivo (.pdf, .docx)
- Límite de tamaño de archivo (50MB por defecto)
- Sanitización de nombres de archivo
- Usuario no-root en contenedor Docker
- Variables de entorno para credenciales

## 📊 Monitoreo

### Health Checks

```bash
# Health check básico
curl http://localhost:8093/health

# Health check con detalles
curl http://localhost:8093/health/detailed
```

### Métricas Disponibles

- Estado de conexión a base de datos
- Estado de conexión a MinIO
- Estadísticas de uso de almacenamiento
- Tiempo de respuesta de IA

## 🚀 Desarrollo Local

### Configuración del Entorno

```bash
# Clonar y navegar al directorio
cd microservicios/documentacion

# Instalar dependencias
go mod tidy

# Ejecutar tests
go test ./...

# Ejecutar en modo desarrollo
go run cmd/main.go
```

### Variables de Entorno para Desarrollo

```bash
export GIN_MODE=debug
export DB_HOST=localhost
export DB_PORT=5434
export MINIO_ENDPOINT=localhost:9000
export GEMINI_API_KEY=your-api-key
```

## 🐛 Troubleshooting

### Problemas Comunes

1. **Error de conexión a base de datos**:
   ```bash
   # Verificar que PostgreSQL esté corriendo
   docker-compose ps documentacion-db
   ```

2. **MinIO no accesible**:
   ```bash
   # Reiniciar servicio MinIO
   docker-compose restart minio
   ```

3. **API Gemini no responde**:
   ```bash
   # Verificar clave API
   echo $GEMINI_API_KEY
   ```

### Logs

```bash
# Ver logs del microservicio
docker-compose logs -f documentacion-service

# Ver logs de la base de datos
docker-compose logs -f documentacion-db
```

## 📈 Roadmap

<!-- ### Próximas Características

- [ ] Versioning de documentos
- [ ] Comentarios y anotaciones
- [ ] Notificaciones de nuevos documentos
- [ ] API GraphQL
- [ ] Integración con más proveedores de IA
- [ ] Dashboard de analytics

### Mejoras de Rendimiento

- [ ] Cache Redis para búsquedas frecuentes
- [ ] Compresión automática de archivos
- [ ] CDN para distribución de contenido
- [ ] Paginación optimizada -->

## 🤝 Contribución

1. Fork del repositorio
2. Crear branch de feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit de cambios (`git commit -am 'Añadir nueva funcionalidad'`)
4. Push al branch (`git push origin feature/nueva-funcionalidad`)
5. Crear Pull Request

## 📄 Licencia

Este proyecto es parte del sistema InspectAR y está sujeto a las políticas de licencia del proyecto principal.

---

**Contacto**: Equipo de Desarrollo InspectAR
**Versión**: 1.0.0
**Puerto del Microservicio**: 8093
**Puerto de Base de Datos**: 5434
**Última actualización**: Agosto 2025
