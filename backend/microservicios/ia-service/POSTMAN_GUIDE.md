# 📮 Guía de Uso - Colección Postman IA-Service

## 🚀 Inicio Rápido

### 1. Importar la Colección

1. Abrir Postman Desktop o Web
2. Click en **"Import"** (esquina superior izquierda)
3. Seleccionar el archivo `IA-Service.postman_collection.json`
4. Click en **"Import"**

### 2. Configurar Variables de Entorno

Crear un nuevo Environment en Postman con las siguientes variables:

| Variable | Valor (Local) | Valor (Docker) | Descripción |
|----------|---------------|----------------|-------------|
| `base_url` | `http://localhost:8095` | `http://ia-service:8095` | URL base del IA-Service |
| `ml_engine_host` | `localhost` | `ml_engine_python` | Host del ML Engine |

**Pasos:**
1. Click en el ícono de engranaje (⚙️) en la esquina superior derecha
2. Click en **"Add"** para crear un nuevo environment
3. Nombre: `IA-Service Local` o `IA-Service Docker`
4. Agregar las variables mencionadas
5. Click en **"Save"**
6. Seleccionar el environment desde el dropdown superior

---

## 📋 Estructura de la Colección

La colección está organizada en **3 carpetas principales**:

### 1️⃣ Health & Status
Endpoints para verificar el estado del servicio:
- ✅ **Health Check**: Verifica que el servicio esté corriendo
- 📊 **Get Latest Anomaly**: Última anomalía detectada

### 2️⃣ Anomalies Management
Endpoints principales para gestionar anomalías:
- 💾 **Store Anomaly**: Guardar anomalía manualmente
- 🔍 **Get Anomalies by Sensor**: Consultar por sensor
- 🆕 **Get Anomalies by Activo**: Consultar por activo (con paginación)

### 3️⃣ ML Engine Integration
Endpoints para testing directo del ML Engine:
- 🤖 **Simulate ML Prediction**: Enviar datos al ML Engine directamente

---

## 🎯 Flujo de Testing Recomendado

### Paso 1: Verificar el Servicio
```
GET /health
```
**Resultado esperado:** `200 OK` con `"status": "healthy"`

### Paso 2: Verificar Anomalías Existentes
```
GET /status/
```
**Posibles resultados:**
- `200 OK`: Hay anomalías, retorna la más reciente
- `404 Not Found`: No hay anomalías aún

### Paso 3: Consultar Anomalías por Activo
```
GET /anomalies/activo/2?limit=10&offset=0
```
**Resultado esperado:** Lista de anomalías del activo 2 con paginación

### Paso 4: Guardar una Anomalía de Prueba
```
POST /anomalies/store
```
**Body:**
```json
{
  "activo_id": 2,
  "sensor_id": "test_sensor_postman",
  "timestamp": "2025-10-16T10:00:00Z",
  "anomaly_score": 0.95,
  "anomaly_likelihood": 0.88,
  "severidad": "critica",
  "descripcion": "Test desde Postman",
  "threshold": 0.75,
  "is_anomaly": 1
}
```
**Resultado esperado:** `201 Created` con el objeto guardado

### Paso 5: Consultar por Sensor Específico
```
GET /anomalies/sensor/test_sensor_postman?limit=5
```
**Resultado esperado:** Lista de anomalías del sensor de prueba

---

## 🧪 Tests Automáticos

Cada endpoint incluye **tests automáticos** que se ejecutan después de cada request:

### Tests Globales (todos los endpoints)
- ✅ Tiempo de respuesta < 5000ms
- ✅ Content-Type es JSON

### Tests Específicos

#### Health Check
- ✅ Status code es 200
- ✅ Campo `status` es "healthy"
- ✅ Campo `service` es "ia-service"

#### Get Latest Anomaly
- ✅ Status code es 200 o 404
- ✅ Si 200, valida estructura de datos

#### Store Anomaly
- ✅ Status code es 201
- ✅ Mensaje de éxito presente
- ✅ Guarda el ID en variable de entorno

#### Get Anomalies by Activo
- ✅ Status code es 200
- ✅ Respuesta tiene array `data`
- ✅ Tiene parámetros de paginación

**Ver resultados:**
- Los tests se ejecutan automáticamente
- Ver resultados en la pestaña "Test Results" después de cada request
- ✅ = Test pasado
- ❌ = Test fallido

---

## 📊 Ejemplos de Respuesta

### ✅ Éxito: Get Anomalies by Activo

```json
{
  "message": "Anomalías obtenidas exitosamente",
  "count": 15,
  "activo_id": 2,
  "limit": 50,
  "offset": 0,
  "data": [
    {
      "id": 130,
      "activo_id": 2,
      "sensor_id": "combined",
      "timestamp": "2025-10-16T03:32:00Z",
      "anomaly_score": 0.9456,
      "anomaly_likelihood": 0.9123,
      "severidad": "critica",
      "descripcion": "Patrón temporal crítico detectado por HTM",
      "threshold": 0.75,
      "is_anomaly": 1,
      "created_at": "2025-10-16T03:32:05Z",
      "updated_at": "2025-10-16T03:32:05Z"
    }
  ]
}
```

### ❌ Error: Invalid Data

```json
{
  "error": "Datos inválidos",
  "details": "Key: 'StoreAnomalyRequest.SensorID' Error:Field validation for 'SensorID' failed on the 'required' tag"
}
```

---

## 🔧 Variables Dinámicas

La colección usa **variables dinámicas de Postman**:

| Variable | Uso | Ejemplo |
|----------|-----|---------|
| `{{$isoTimestamp}}` | Timestamp actual en ISO 8601 | `2025-10-16T10:30:45.123Z` |
| `{{$randomInt}}` | Número aleatorio | `12345` |
| `{{last_anomaly_id}}` | ID de última anomalía guardada | `456` (guardado automáticamente) |

**Ejemplo de uso en Body:**
```json
{
  "timestamp": "{{$isoTimestamp}}",
  "sensor_id": "sensor_{{$randomInt}}"
}
```

---

## 📝 Casos de Uso Comunes

### 1. Testing Completo del Servicio
Ejecutar **toda la carpeta** "Anomalies Management":
1. Click derecho en la carpeta
2. Seleccionar "Run folder"
3. Ver resultados en el Collection Runner

### 2. Paginación de Anomalías
```
# Primera página
GET /anomalies/activo/2?limit=10&offset=0

# Segunda página
GET /anomalies/activo/2?limit=10&offset=10

# Tercera página
GET /anomalies/activo/2?limit=10&offset=20
```

### 3. Filtrar Anomalías Críticas
Primero obtener todas:
```
GET /anomalies/activo/2?limit=100
```
Luego filtrar en código o usar Postman Visualizer.

### 4. Simular Carga de Datos al ML Engine
```
POST http://localhost:8085/predict_anomaly
```
Con 180+ registros en el body.

---

## 🐛 Troubleshooting

### Error: "Could not get any response"
**Solución:**
- Verificar que el servicio esté corriendo: `docker ps | grep ia-service`
- Verificar la URL: `http://localhost:8095` (no `https`)
- Probar: `curl http://localhost:8095/health`

### Error: "404 Not Found"
**Solución:**
- Verificar el endpoint en la documentación
- Verificar que el path esté correcto
- Algunos endpoints requieren trailing slash: `/status/`

### Error: "400 Bad Request - Datos inválidos"
**Solución:**
- Verificar que todos los campos requeridos estén presentes
- Verificar tipos de datos (int, string, float)
- Ver campo `details` en la respuesta para más información

### Tests Fallan
**Solución:**
- Verificar que el servicio esté corriendo correctamente
- Verificar que haya datos en la BD (algunas queries retornan 404 si no hay datos)
- Ver consola de Postman para logs detallados

---

## 📚 Recursos Adicionales

- **README Principal**: `./README.md`
- **Documentación API**: Ver sección "Endpoints" en README.md
- **Código Fuente**: `./internal/handler/anomaly_handler.go`
- **Esquema de BD**: `../../database/ia/README.md`

---

## 🎓 Tips y Mejores Prácticas

### 1. Usar Pre-request Scripts
Los scripts se ejecutan antes de cada request:
```javascript
console.log('Ejecutando request a:', pm.request.url);
```

### 2. Guardar Valores en Variables
```javascript
// En el Test Script
var jsonData = pm.response.json();
pm.environment.set("anomaly_id", jsonData.data.id);
```

### 3. Organizar Requests en Carpetas
Crear carpetas personalizadas:
- "Development Tests"
- "Production Validation"
- "Performance Tests"

### 4. Exportar Resultados
Usar Collection Runner para:
- Generar reportes HTML
- Exportar a JSON
- Integrar con CI/CD

---

## 🤝 Contribuir

Si encuentras issues o quieres agregar más tests:
1. Editar el archivo JSON directamente
2. Importar en Postman para probar
3. Exportar y hacer commit

---

## 📞 Soporte

Para issues o preguntas:
- Ver logs del servicio: `docker logs ia-service -f`
- Reportar en el repositorio de InspectAR
- Consultar documentación técnica en README.md

---

**Última actualización:** 2025-10-16  
**Versión de la colección:** 2.0.0  
**Autor:** InspectAR Team
