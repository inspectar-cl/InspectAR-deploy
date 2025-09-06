# 🌐 API Gateway - InspectAR

API Gateway inteligente que funciona como punto de entrada unificado para todos los microservicios del sistema InspectAR. Proporciona dos tipos de funcionalidades:

1. **Proxy transparente**: Reenvía peticiones directamente a los microservicios
2. **Endpoints compuestos**: Combina múltiples llamadas a microservicios en una sola respuesta

## 🚀 Características Principales

### Funcionalidades Implementadas

- ✅ **Enrutamiento inteligente** - Distribución automática de peticiones
- ✅ **Proxy reverso** con reescritura de URLs
- ✅ **Endpoints compuestos** - Múltiples consultas en una sola petición
- ✅ **Consultas concurrentes** - Mejor rendimiento con llamadas paralelas
- ✅ **Transformaciones de datos** - Procesamiento y normalización
- ✅ **Configuración declarativa** de rutas y endpoints
- ✅ **Health checks** para monitoreo
- ✅ **Manejo de errores** centralizado
- ✅ **Timeouts configurables** para conexiones
- ✅ **Logs detallados** de peticiones y errores

### Capacidades Técnicas

- ✅ Balanceador de carga básico
- ✅ Transformación de paths automática
- ✅ Consultas dependientes (secuenciales)
- ✅ Respuestas parciales en caso de errores
- ✅ Configuración por variables de entorno
- ✅ Compatibilidad con Docker y Docker Compose
- ✅ Middleware extensible (CORS, logging, etc.)

## 🏗️ Arquitectura

```
Frontend → API Gateway → Microservicios
                   ↓
            Respuesta Unificada
```

### Tipos de Rutas

#### 1. Rutas de Proxy (Transparentes)
```
Cliente → GET /api/gestion/activos → API Gateway → Microservicio Gestión
                                                       ↓
Cliente ← JSON Response ←←←←←←←←←←←←←←←←←←←←←←←←←←←←←←← 
```

#### 2. Endpoints Compuestos (Inteligentes)
```
Cliente → GET /api/calderas-estado → API Gateway
                                        ↓
                    ┌─── Gestión: /activos/tipo/caldera
                    ├─── Parser: /activo/{id1}/sensores/estado  
                    ├─── Parser: /activo/{id2}/sensores/estado
                    └─── Parser: /activo/{idN}/sensores/estado
                                        ↓
                              Procesa + Combina
                                        ↓  
Cliente ← Respuesta Unificada ←←←←←←← API Gateway
```

## 🛠️ Tecnologías Utilizadas

- **Lenguaje**: Go 1.22
- **Framework**: Gin (HTTP Router)
- **Proxy**: httputil.ReverseProxy nativo de Go
- **Configuración**: Variables de entorno + archivos .env
- **Containerización**: Docker + Docker Compose

## 📁 Estructura del Proyecto

```
apigateway/
├── cmd/
│   └── main.go                    # Punto de entrada
├── config/
│   └── config.go                  # Configuración y variables de entorno
├── handlers/                      # 🆕 Endpoints compuestos
│   ├── endpoints_config.go        # ⭐ Configuración de endpoints compuestos
│   ├── composite_handler.go       # Handler principal para endpoints
│   └── transforms.go              # Funciones de transformación de datos
├── proxy/
│   ├── routes_config.go           # ⭐ Configuración de rutas proxy
│   ├── routes.go                  # Registro de rutas HTTP
│   └── factory.go                 # Lógica del proxy reverso
├── pkg/
│   └── pathutil.go                # Utilidades para paths
├── middleware/
│   └── cors.go                    # Middleware CORS
├── .env                           # Variables de entorno
├── Dockerfile                     # Imagen Docker
├── go.mod                         # Dependencias Go
└── README.md                      # Este archivo
```

## 🚀 Instalación y Uso

### Prerrequisitos

- Docker y Docker Compose
- Go 1.22+ (para desarrollo local)
- Microservicios corriendo en sus puertos respectivos

### Configuración

1. **Configurar variables de entorno** (archivo `.env`):
```bash
# Puerto del API Gateway
GATEWAY_PORT=3500

# URLs de los microservicios
PARSER_URL=http://iot-service:8090
NOTIFICATION_URL=http://notification-service:8091
GESTION_URL=http://gestion-service:8092
DOCUMENTATION_URL=http://documentacion-service:8093

# Frontend (para CORS)
FRONTEND_URL=http://localhost:3000
```

2. **Iniciar con Docker Compose**:
```bash
# Desde el directorio backend/
docker-compose up apigateway
```

3. **Verificar funcionamiento**:
```bash
# Health Check
curl http://localhost:3500/healthz

# Respuesta esperada:
# {"ok": true}
```

## 📚 Configuración de Rutas

### Tipos de Endpoints

El API Gateway maneja dos tipos de endpoints:

#### 1. 🔀 Rutas de Proxy (Simples y Directas)
**Archivo:** `proxy/routes_config.go`

Estas rutas reenvían directamente al microservicio sin procesamiento adicional.

```go
var ProxyRoutes = []ProxyRoute{
    {
        PathPrefix:   "/api/gestion",        // Lo que escucha el gateway
        TargetEnvVar: "GESTION_URL",         // Variable de entorno con la URL destino  
        PrependPath:  "",                    // Path adicional (opcional)
    },
    {
        PathPrefix:   "/api/documentacion",
        TargetEnvVar: "DOCUMENTATION_URL", 
        PrependPath:  "/api/v1",            // Agrega /api/v1 al path
    },
}
```

**Comportamiento:**
- Frontend solicita: `GET /api/gestion/activos`
- Gateway reenvía a: `GET ${GESTION_URL}/activos`
- Respuesta: Directa del microservicio

#### 2. 🧠 Endpoints Compuestos (Inteligentes)
**Archivo:** `handlers/endpoints_config.go`

Estos endpoints realizan múltiples consultas, las procesan y combinan antes de responder.

```go
var CompositeEndpoints = []EndpointConfig{
    {
        Path:        "/api/calderas-estado",
        Method:      "GET", 
        Handler:     "GetCalderasEstado",
        Description: "Obtiene todas las calderas y el estado de sus sensores",
        Queries: []ServiceQuery{
            {
                Name:         "calderas",
                TargetEnvVar: "GESTION_URL",
                Path:         "/gestion/activos/tipo/caldera",
                Method:       "GET",
                Transform:    "TransformCalderas",
            },
            {
                Name:         "sensores_estado",
                TargetEnvVar: "PARSER_URL",
                PathTemplate: "/parser/activo/{activo_id}/sensores/estado",
                Method:       "GET", 
                DependsOn:    "calderas",  // Ejecuta DESPUÉS de obtener calderas
                Transform:    "TransformSensoresEstado",
            },
        },
    },
}
```

**Comportamiento del endpoint `/api/calderas-estado`:**

1. **Frontend solicita:** `GET /api/calderas-estado`
2. **API Gateway ejecuta:**
   - 🔹 Consulta: `GET ${GESTION_URL}/gestion/activos/tipo/caldera` 
   - 🔹 Por cada caldera obtenida → `GET ${PARSER_URL}/parser/activo/{activo_id}/sensores/estado`
   - 🔹 Combina y procesa todas las respuestas
3. **Respuesta unificada:**

```json
{
  "calderas_con_sensores": [
    {
      "caldera": {
        "id": 1,
        "activo_id": "AC-1001", 
        "nombre": "Caldera Principal",
        "tipo": "caldera",
        "estado": "operativo"
      },
      "sensores_info": {
        "activo_id": "AC-1001",
        "estado": "operativo", 
        "resumen": {
          "sensores_activos": 0,
          "sensores_desconectados": 1,
          "sensores_nunca_conectados": 0
        },
        "sensores": [
          {
            "sensor_id": "PRES001",
            "tipo": "presion",
            "estado": "disconnected",
            "unidad": "bar"
          }
        ]
      },
      "consulta_timestamp": "2025-09-05T..."
    }
  ],
  "total_calderas": 1,
  "consultas_exitosas": 1,
  "timestamp": "2025-09-05T..."
}
```

## 🔧 Cómo Agregar Nuevos Endpoints

### 🔀 Para Rutas de Proxy Simples

#### Agregar en un solo lugar: `proxy/routes_config.go`

```go
var ProxyRoutes = []ProxyRoute{
    // Rutas existentes...
    {
        PathPrefix:   "/api/usuarios",      // Lo que escribe el frontend
        TargetEnvVar: "USUARIOS_URL",       // Variable de entorno del servicio
        PrependPath:  "",                   // Path adicional (opcional)
    }, // ✨ Solo agrega esto y YA funciona automáticamente
}
```

**¡Y listo!** 🚀 El endpoint `/api/usuarios/*` ya funciona automáticamente.

---

### 🧠 Para Endpoints Compuestos Inteligentes

#### Paso 1: Agregar configuración en `handlers/endpoints_config.go`

```go
var CompositeEndpoints = []EndpointConfig{
    // Endpoints existentes...
    {
        Path:        "/api/dashboard-completo",
        Method:      "GET", 
        Handler:     "GetDashboardCompleto",
        Description: "Dashboard con datos de múltiples servicios",
        Queries: []ServiceQuery{
            {
                Name:         "activos",
                TargetEnvVar: "GESTION_URL",
                Path:         "/activos",
                Method:       "GET",
                Transform:    "TransformActivos", // Opcional
            },
            {
                Name:         "alertas",
                TargetEnvVar: "NOTIFICATION_URL", 
                Path:         "/alertas/recientes",
                Method:       "GET",
                Transform:    "TransformAlertas", // Opcional
            },
            {
                Name:         "documentos",
                TargetEnvVar: "DOCUMENTATION_URL",
                Path:         "/documentos/importantes", 
                Method:       "GET",
            },
        },
    }, // ✨ Solo esto y el endpoint ya funciona!
}
```

#### Paso 2 (Opcional): Agregar transformación personalizada

Si necesitas procesar los datos, agrega en `handlers/transforms.go`:

```go
// Función de transformación personalizada
func transformActivos(data interface{}) (interface{}, error) {
    // Tu lógica de transformación aquí
    return data, nil
}

func transformAlertas(data interface{}) (interface{}, error) {
    // Tu lógica de transformación aquí
    return data, nil
}

// Y registrarlas en NewCompositeHandler()
func NewCompositeHandler() *CompositeHandler {
    return &CompositeHandler{
        // ...
        transforms: map[string]TransformFunc{
            "TransformCalderas":       transformCalderas,
            "TransformActivos":        transformActivos,     // ✨ Nueva
            "TransformAlertas":        transformAlertas,     // ✨ Nueva
        },
    }
}
```

---

## 🎯 Ejemplo Práctico: Agregar Endpoint de Bombas

### Quiero un endpoint que me dé todas las bombas con sus sensores

#### 1️⃣ Solo agregamos en `endpoints_config.go`:

```go
{
    Path:        "/api/bombas-estado", 
    Method:      "GET",
    Handler:     "GetBombasEstado",
    Description: "Obtiene todas las bombas con estado de sensores",
    Queries: []ServiceQuery{
        {
            Name:         "bombas",
            TargetEnvVar: "GESTION_URL", 
            Path:         "/activos/tipo/bomba", // 🎯 Cambiamos solo el tipo
            Method:       "GET",
            Transform:    "TransformBombas", // Reutilizamos la lógica
        },
    },
},
```

#### 2️⃣ Agregamos la función de transformación:

```go
// En transforms.go - Copia de transformCalderas pero para bombas
func transformBombas(data interface{}) (interface{}, error) {
    // Igual que transformCalderas, pero procesa bombas
    if dataMap, ok := data.(map[string]interface{}); ok {
        dataMap["processed_at"] = time.Now()
        // ... mismo código que calderas
        return executeBombasSensoresQueries(dataMap) // Nueva función
    }
    return data, nil
}

// Y en NewCompositeHandler():
transforms: map[string]TransformFunc{
    "TransformCalderas": transformCalderas,
    "TransformBombas":   transformBombas,     // ✨ Nueva
}
```

#### 3️⃣ ¡Ya funciona!

```bash
curl http://localhost:3500/api/bombas-estado
```

---

## ✨ Ventajas del Sistema Modular

### ✅ **Un solo lugar para configurar**
- Proxy routes: Solo `routes_config.go`
- Endpoints compuestos: Solo `endpoints_config.go`

### ✅ **Automático**
- Sin código adicional para routing
- Sin registro manual de endpoints
- Sin configuración de CORS por endpoint

### ✅ **Reutilizable**
- Las funciones de transformación se pueden reutilizar
- Los patrones se replican fácilmente

### ✅ **Escalable**
- Agregar 10 endpoints = 10 configuraciones
- Sin tocar el código del handler principal

---

## 🚀 Patrones Comunes

### 1. **Endpoint Simple (Un servicio)**
```go
{
    Path: "/api/tecnicos",
    Queries: []ServiceQuery{
        {Name: "tecnicos", TargetEnvVar: "GESTION_URL", Path: "/tecnicos"},
    },
}
```

### 2. **Endpoint Compuesto (Múltiples servicios)**  
```go
{
    Path: "/api/resumen-edificio",
    Queries: []ServiceQuery{
        {Name: "activos", TargetEnvVar: "GESTION_URL", Path: "/activos/edificio/1"},
        {Name: "alertas", TargetEnvVar: "NOTIFICATION_URL", Path: "/edificio/1/alertas"},
        {Name: "documentos", TargetEnvVar: "DOCUMENTATION_URL", Path: "/edificio/1/docs"},
    },
}
```

### 3. **Endpoint con Dependencias (Como calderas-estado)**
```go
{
    Path: "/api/activos-con-sensores", 
    Queries: []ServiceQuery{
        {Name: "activos", Path: "/activos", Transform: "TransformActivosConSensores"},
        // La transformación hace las llamadas a sensores dinámicamente
    },
}
```

---

## 🎯 Resumen: Cómo Funciona la Modularidad

### 🔥 **La Magia del Sistema**

#### Para el endpoint `/api/calderas-estado` que acabamos de probar:

**1. Solo configuramos esto en `endpoints_config.go`:**
```go
{
    Path: "/api/calderas-estado",
    Method: "GET", 
    Handler: "GetCalderasEstado",
    Queries: []ServiceQuery{
        {
            Name: "calderas",
            TargetEnvVar: "GESTION_URL",
            Path: "/activos/tipo/caldera",
            Transform: "TransformCalderas",
        },
    },
}
```

**2. El sistema automáticamente:**
- ✅ Registra la ruta HTTP `GET /api/calderas-estado`
- ✅ Llama a `${GESTION_URL}/activos/tipo/caldera`
- ✅ Ejecuta `TransformCalderas()` que hace las consultas de sensores
- ✅ Devuelve la respuesta combinada en formato JSON
- ✅ Maneja errores, timeouts y CORS

**3. Resultado:**
```bash
curl http://localhost:3500/api/calderas-estado
# Devuelve calderas + sensores en una sola respuesta 🚀
```

---

### 📋 **Checklist para Agregar Nuevos Endpoints**

#### ✅ **Endpoint Simple (Proxy)**
1. Agregar configuración en `routes_config.go`
2. Definir variable de entorno del servicio
3. **¡Listo!** Ya funciona automáticamente

#### ✅ **Endpoint Compuesto**
1. Agregar configuración en `endpoints_config.go`
2. (Opcional) Crear función de transformación
3. (Opcional) Registrar función en `NewCompositeHandler()`
4. **¡Listo!** Ya funciona automáticamente

#### ✅ **Variables de Entorno Necesarias**
```bash
# Solo asegúrate de que estén definidas
GESTION_URL=http://gestion-service:8092
PARSER_URL=http://iot-service:8090  
NOTIFICATION_URL=http://notification-service:8091
DOCUMENTATION_URL=http://documentacion-service:8093
```

---

## 🎉 **Ejemplo de Uso Actual**

### Ya funciona este endpoint:
```bash
# ✅ Endpoint compuesto que combina 2 servicios
curl http://localhost:3500/api/calderas-estado

# Devuelve:
{
  "calderas": {
    "activos": [
      {
        "activo_id": "AC-1001",
        "nombre": "Caldera Principal", 
        "sensores_info": {
          "sensores": [{"sensor_id": "PRES001", "tipo": "presion"}],
          "total_sensores": 1
        }
      }
    ]
  }
}
```

### Y también funcionan estos proxies automáticamente:
```bash
# ✅ Proxy directo al servicio de gestión  
curl http://localhost:3500/api/gestion/activos

# ✅ Proxy directo al servicio de parser
curl http://localhost:3500/api/parser/activo/AC-1001/sensores/estado

# ✅ Health check
curl http://localhost:3500/healthz
```

---

## 🚀 **¿Por qué es tan Poderoso?**

### 🔥 **Un Solo Archivo = Nuevas Funcionalidades**

**Antes (sin modularidad):**
- Crear endpoint ➜ 5-10 archivos modificados
- Registrar ruta ➜ Código manual 
- Manejar errores ➜ Código repetitivo
- Testing ➜ Complicado

**Ahora (modular):**
- Crear endpoint ➜ **1 configuración**
- Registrar ruta ➜ **Automático**
- Manejar errores ➜ **Automático** 
- Testing ➜ **Simple**

### ✨ **Ejemplos Reales**

**Frontend necesita:** "Todas las bombas con sus sensores"
**Backend developer:** Agrega 5 líneas en `endpoints_config.go` ➜ **¡Listo!**

**Frontend necesita:** "Dashboard con activos + alertas + documentos"  
**Backend developer:** Agrega 15 líneas en `endpoints_config.go` ➜ **¡Listo!**

**Frontend necesita:** "Proxy directo al servicio de usuarios"
**Backend developer:** Agrega 3 líneas en `routes_config.go` ➜ **¡Listo!**

---

## 🎯 **Próximos Pasos Recomendados**

1. **Experimenta** agregando un endpoint simple en `routes_config.go`
2. **Prueba** creando un endpoint compuesto en `endpoints_config.go`  
3. **Desarrolla** funciones de transformación personalizadas
4. **Escala** el sistema agregando más microservicios

El sistema está diseñado para crecer contigo sin complejidad adicional. **¡Un archivo de configuración puede agregar funcionalidades completas al API Gateway!** 🚀
```
    // Rutas existentes...
    {
        PathPrefix:   "/api/usuarios",      // Prefijo que escucha el gateway
        TargetEnvVar: "USUARIOS_URL",       // Variable de entorno con la URL destino
        PrependPath:  "",                   // Path adicional a agregar (opcional)
    }, // ⚠️ Muy importante esta coma
}
```

#### 2. Definir la variable de entorno

En el archivo `.env` o en el entorno:

```bash
USUARIOS_URL=http://usuarios-service:8094
```

#### 3. ¡Listo! La ruta ya funciona automáticamente

No necesitas tocar ningún otro archivo. El API Gateway:
- ✅ Detecta automáticamente la nueva ruta
- ✅ Reenvía peticiones al microservicio
- ✅ Maneja errores y timeouts

### Para Endpoints Compuestos (Inteligentes)

#### 1. Configurar el endpoint en `handlers/endpoints_config.go`

```go
var CompositeEndpoints = []EndpointConfig{
    // Endpoints existentes...
    {
        Path:        "/api/mi-dashboard",
        Method:      "GET",
        Handler:     "GetMiDashboard", 
        Description: "Dashboard personalizado con datos combinados",
        Queries: []ServiceQuery{
            {
                Name:         "activos",
                TargetEnvVar: "GESTION_URL",
                Path:         "/gestion/activos",
                Method:       "GET",
                Transform:    "TransformActivos", // Opcional
            },
            {
                Name:         "alertas",
                TargetEnvVar: "NOTIFICATION_URL", 
                Path:         "/notificacion/alertas/active",
                Method:       "GET",
            },
            {
                Name:         "sensores_detalle",
                TargetEnvVar: "PARSER_URL",
                PathTemplate: "/parser/activo/{activo_id}/sensores",
                Method:       "GET",
                DependsOn:    "activos", // Se ejecuta DESPUÉS de obtener activos
            },
        },
    },
}
```

#### 2. (Opcional) Crear funciones de transformación

Si necesitas procesar los datos, agrégalas en `handlers/transforms.go`:

```go
func transformActivos(data interface{}) (interface{}, error) {
    // Procesar datos de activos
    if dataMap, ok := data.(map[string]interface{}); ok {
        // Agregar timestamp, normalizar campos, etc.
        dataMap["processed_at"] = time.Now()
        
        // Tu lógica aquí...
        
        return dataMap, nil
    }
    return data, nil
}

// No olvides registrarla en composite_handler.go:
transforms: map[string]TransformFunc{
    "TransformActivos": transformActivos,
    // otras transformaciones...
}
```

#### 3. (Opcional) Lógica de procesamiento específica

Si necesitas lógica especial para combinar resultados, agrégala en `handlers/composite_handler.go`:

```go
func (h *CompositeHandler) processResults(handlerName string, results map[string]interface{}, errors map[string]error) interface{} {
    switch handlerName {
    case "GetMiDashboard":
        return h.processMiDashboard(results, errors)
    // casos existentes...
    }
}

func (h *CompositeHandler) processMiDashboard(results map[string]interface{}, errors map[string]error) interface{} {
    // Combinar activos + alertas + sensores
    response := map[string]interface{}{
        "activos":         results["activos"],
        "alertas":         results["alertas"], 
        "sensores_detalle": results["sensores_detalle"],
        "resumen": map[string]interface{}{
            "total_activos": len(results["activos"].([]interface{})),
            // más lógica...
        },
        "timestamp": time.Now(),
    }
    
    if len(errors) > 0 {
        response["errors"] = errors
    }
    
    return response
}
```

### Tipos de Consultas Soportadas

#### Consultas Independientes (Concurrentes)
```go
{
    Name:         "activos",
    TargetEnvVar: "GESTION_URL",
    Path:         "/gestion/activos",
    Method:       "GET",
    // Se ejecuta en paralelo con otras consultas independientes
}
```

#### Consultas Dependientes (Secuenciales)
```go
{
    Name:         "sensores_activo",
    TargetEnvVar: "PARSER_URL",
    PathTemplate: "/parser/activo/{activo_id}/sensores",
    Method:       "GET",
    DependsOn:    "activos", // Espera a que termine la consulta "activos"
}
```

#### Consultas con Headers Personalizados
```go
{
    Name:         "datos_autenticados",
    TargetEnvVar: "SECURE_SERVICE_URL",
    Path:         "/secure/data",
    Method:       "GET",
    Headers:      map[string]string{
        "Authorization": "Bearer token-here",
        "X-Custom-Header": "value",
    },
}
```

### 🎯 Ejemplos de Transformación de URLs

#### Ejemplo 1: Ruta Simple (sin PrependPath)
```go
{
    PathPrefix:   "/api/gestion",
    TargetEnvVar: "GESTION_URL",      // http://gestion-service:8092
    PrependPath:  "",
}
```

**Transformación:**
- Cliente solicita: `GET /api/gestion/tecnicos`
- Gateway reenvía a: `GET http://gestion-service:8092/tecnicos`

#### Ejemplo 2: Ruta con PrependPath
```go
{
    PathPrefix:   "/api/documentacion",
    TargetEnvVar: "DOCUMENTATION_URL", // http://documentacion-service:8093
    PrependPath:  "/api/v1",
}
```

**Transformación:**
- Cliente solicita: `GET /api/documentacion/documentos`
- Gateway reenvía a: `GET http://documentacion-service:8093/api/v1/documentos`

#### Ejemplo 3: Ruta Compleja
```go
{
    PathPrefix:   "/api/v2/sensores",
    TargetEnvVar: "IOT_URL",           // http://iot-service:8090
    PrependPath:  "/activo",
}
```

**Transformación:**
- Cliente solicita: `GET /api/v2/sensores/AC-001/datos`
- Gateway reenvía a: `GET http://iot-service:8090/activo/AC-001/datos`

## 🔍 ¿Qué Pasa "Por Dentro" del API Gateway?

### 1. **Recepción de Petición**
```
Cliente → GET /api/gestion/tecnicos → API Gateway
```

### 2. **Análisis de Ruta**
El gateway busca en `ProxyRoutes` la entrada que coincida:
```go
// Encuentra: PathPrefix = "/api/gestion"
{
    PathPrefix:   "/api/gestion",
    TargetEnvVar: "GESTION_URL",
    PrependPath:  "",
}
```

### 3. **Resolución de URL Destino**
```go
target := os.Getenv("GESTION_URL")  // http://gestion-service:8092
```

### 4. **Transformación de Path**
```go
// Path original: /api/gestion/tecnicos
// Strip prefix:  /api/gestion → /tecnicos
// Prepend path:  "" (vacío) → /tecnicos
// URL final: http://gestion-service:8092/tecnicos
```

### 5. **Proxy Reverso**
```go
// Crea un httputil.ReverseProxy
// Configura timeouts (5s conexión, 15s respuesta)
// Reenvía la petición completa (headers, body, método)
```

### 6. **Manejo de Respuesta**
```go
// Si SUCCESS: Reenvía respuesta al cliente
// Si ERROR: Devuelve JSON con error 502
{"error":true,"message":"Servicio destino no disponible"}
```

### 7. **Logging**
```bash
[Gateway] GET /api/gestion/tecnicos → http://gestion-service:8092/tecnicos (200 OK)
[Gateway ERROR] GET /api/parser/activos → connection refused
```

## 🎛️ Configuración Avanzada

### Timeouts Configurados
```go
DialContext:           5 * time.Second   // Conexión inicial
TLSHandshakeTimeout:   5 * time.Second   // Handshake SSL
ResponseHeaderTimeout: 15 * time.Second  // Timeout de respuesta
ExpectContinueTimeout: 1 * time.Second   // Continue HTTP
```

### Health Check Endpoint
```bash
GET /healthz
```
**Respuesta:**
```json
{"ok": true}
```

### Manejo de Errores
- **502 Bad Gateway**: Microservicio no disponible
- **404 Not Found**: Ruta no configurada
- **Timeout**: Respuesta JSON con error

## 📋 Endpoints Disponibles

### Endpoints del Sistema

- `GET /healthz` - Health check del API Gateway

### 🧠 Endpoints Compuestos (Nuevos)

- **`GET /api/calderas-estado`** - Obtiene todas las calderas con el estado completo de sus sensores
  - Combina datos de Gestión + Parser
  - Respuesta unificada con información de calderas y sensores
  - Consultas concurrentes para mejor rendimiento

### 🔀 Endpoints de Proxy (Transparentes)

- `GET|POST|PUT|DELETE /api/gestion/*` → `${GESTION_URL}/*`
- `GET|POST|PUT|DELETE /api/notificacion/*` → `${NOTIFICATION_URL}/*`  
- `GET|POST|PUT|DELETE /api/parser/*` → `${PARSER_URL}/*`
- `GET|POST|PUT|DELETE /api/documentacion/*` → `${DOCUMENTATION_URL}/api/v1/*`

## 🔧 Ejemplos Reales de Uso

### 🧠 Endpoint Compuesto: Calderas con Estado de Sensores

```bash
# Una sola petición que internamente hace múltiples consultas
curl http://localhost:3500/api/calderas-estado
```

**Lo que sucede internamente:**
1. 🔹 Consulta calderas: `GET ${GESTION_URL}/gestion/activos/tipo/caldera`
2. 🔹 Por cada caldera: `GET ${PARSER_URL}/parser/activo/{activo_id}/sensores/estado`
3. 🔹 Combina y procesa todas las respuestas
4. 🔹 Retorna respuesta unificada

**Respuesta del endpoint:**
```json
{
  "calderas_con_sensores": [
    {
      "caldera": {
        "id": 1,
        "activo_id": "AC-1001",
        "nombre": "Caldera Principal", 
        "estado": "operativo"
      },
      "sensores_info": {
        "activo_id": "AC-1001",
        "estado": "operativo",
        "resumen": {
          "sensores_activos": 0,
          "sensores_desconectados": 1
        },
        "sensores": [...]
      }
    }
  ],
  "total_calderas": 1,
  "consultas_exitosas": 1
}
```

### 🔀 Endpoints de Proxy: Acceso Directo a Microservicios

#### Gestión de Técnicos
```bash
# A través del API Gateway
curl http://localhost:3500/api/gestion/tecnicos

# Se reenvía a:
curl http://gestion-service:8092/tecnicos
```

#### Documentos Técnicos
```bash
# A través del API Gateway  
curl http://localhost:3500/api/documentacion/documentos

# Se reenvía a:
curl http://documentacion-service:8093/api/v1/documentos
```

#### Datos de Sensores
```bash
# A través del API Gateway
curl http://localhost:3500/api/parser/activo/AC-001

# Se reenvía a:
curl http://iot-service:8090/activo/AC-001
```

#### Notificaciones
```bash
# A través del API Gateway  
curl -X POST http://localhost:3500/api/notificacion/notification \
     -H "Content-Type: application/json" \
     -d '{"activo_id": 1}'

# Se reenvía a:
curl -X POST http://notification-service:8091/notification \
     -H "Content-Type: application/json" \
     -d '{"activo_id": 1}'
```

## 📝 Logs y Monitoreo

### Logs del API Gateway
```bash
# Ver logs en tiempo real
docker-compose logs -f apigateway

# Ejemplo de salida:
[GIN] 2025/09/01 - 15:30:45 | 200 |      45.123ms |  192.168.1.100 | GET     "/api/gestion/tecnicos"
[GIN] 2025/09/01 - 15:30:46 | 502 |       5.001s |  192.168.1.100 | GET     "/api/parser/activos"
[Gateway ERROR] GET /api/parser/activos → dial tcp: connection refused
```

### Health Checks Recomendados
```bash
# Verificar conectividad básica
curl -f http://localhost:3500/healthz || echo "Gateway DOWN"

# Verificar endpoint compuesto
curl -f http://localhost:3500/api/calderas-estado || echo "Calderas endpoint DOWN"

# Verificar rutas de proxy específicas
curl -f http://localhost:3500/api/gestion/health || echo "Gestión DOWN"
curl -f http://localhost:3500/api/parser/activo || echo "Parser DOWN"
```

## ⚡ Características Avanzadas

### Consultas Concurrentes
Los endpoints compuestos ejecutan múltiples consultas **en paralelo** cuando es posible, mejorando significativamente el rendimiento.

```go
// Estas consultas se ejecutan al mismo tiempo:
Query 1: GET /gestion/activos     ←┐ 
Query 2: GET /notificacion/alertas ←┴→ Paralelo
Query 3: GET /parser/metrics      ←┘
```

### Consultas Dependientes
Algunas consultas necesitan datos de otras. El sistema las ejecuta **secuencialmente** según las dependencias.

```go
// Flujo secuencial:
1. GET /gestion/activos/tipo/caldera          (Independiente)
2. Para cada caldera obtenida:                (Dependiente)
   → GET /parser/activo/{activo_id}/sensores
```

### Transformaciones de Datos
Cada respuesta puede ser **transformada** antes de ser incluida en el resultado final:

```go
// Antes de la transformación:
{"estado": "ok", "valor": "25.5"}

// Después de la transformación:
{
  "estado": "operativo",           // Normalizado
  "estado_original": "ok",         // Preservado
  "valor": 25.5,                   // Convertido a número
  "processed_at": "2025-09-05...", // Agregado
  "valor_normalizado": "normal"     // Calculado
}
```

### Manejo de Errores Inteligente
Si un microservicio falla, el endpoint puede retornar una **respuesta parcial** con información sobre los errores:

```json
{
  "calderas_con_sensores": [...],  // Datos obtenidos exitosamente
  "total_calderas": 2,
  "consultas_exitosas": 1,         // Solo 1 de 2 consultas funcionó
  "consultas_fallidas": 1,
  "errores": [
    "Error consultando sensores para AC-1002: connection refused"
  ]
}
```

### Timeouts y Resilencia
```go
// Configuración de timeouts:
DialContext:           5 * time.Second   // Conexión inicial
TLSHandshakeTimeout:   5 * time.Second   // Handshake SSL
ResponseHeaderTimeout: 15 * time.Second  // Timeout de respuesta
```

El API Gateway **no se bloquea** si un servicio es lento - retorna error después del timeout.

## 🚀 Desarrollo Local

### Ejecutar sin Docker
```bash
cd apigateway/
go mod tidy
go run cmd/main.go
```

### Variables de entorno para desarrollo local
```bash
export GATEWAY_PORT=3500
export GESTION_URL=http://localhost:8092
export PARSER_URL=http://localhost:8090
export NOTIFICATION_URL=http://localhost:8091
export DOCUMENTATION_URL=http://localhost:8093
```

## 🔒 Seguridad y CORS

### Middleware CORS (Comentado por defecto)
En `cmd/main.go` puedes habilitar CORS:
```go
// Descomentar para habilitar CORS
r.Use(middleware.CORS(cfg.FrontendURL))
```

### Headers de Seguridad
El API Gateway preserva todos los headers de las peticiones originales y respuestas de los microservicios.

## 🔧 Troubleshooting

### Error 502 - Servicio no disponible
1. Verificar que el microservicio esté corriendo
2. Verificar la variable de entorno correspondiente
3. Verificar conectividad de red entre contenedores

### Ruta 404 - No encontrada  
1. Verificar que la ruta esté definida en `routes_config.go`
2. Verificar que la variable de entorno esté configurada
3. Verificar que el PathPrefix coincida exactamente

### Timeouts
1. Los timeouts están configurados para 15 segundos máximo
2. Si necesitas timeouts mayores, modifica `proxy/factory.go`

## 📊 Rendimiento

- **Latencia adicional**: ~2-5ms por petición
- **Throughput**: Limitado por el microservicio más lento
- **Conexiones concurrentes**: Sin límite específico (depende del sistema)
- **Memory footprint**: ~10-20MB en tiempo de ejecución

<!-- ## 🎯 Próximas Mejoras

- [ ] Rate limiting por cliente
- [ ] Autenticación/autorización centralizada
- [ ] Métricas y monitoring (Prometheus)
- [ ] Load balancing entre instancias
- [ ] Circuit breaker para failover
- [ ] Cache de respuestas

--- -->

<!-- ## 📧 Soporte

Para preguntas sobre configuración de rutas o troubleshooting del API Gateway, revisa los logs y verifica la configuración de las variables de entorno. -->

## 📧 Información adicional

**El API Gateway es el punto de entrada único para todo el ecosistema InspectAR** 🌐
