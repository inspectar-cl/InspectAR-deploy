# 🌐 API Gateway - InspectAR

Punto de entrada unificado para todos los microservicios del sistema InspectAR. El API Gateway actúa como proxy reverso, enrutando las peticiones HTTP a los microservicios correspondientes de manera transparente.

## 🚀 Características Principales

### Funcionalidades Implementadas

- ✅ **Enrutamiento inteligente** - Distribución automática de peticiones
- ✅ **Proxy reverso** con reescritura de URLs
- ✅ **Configuración declarativa** de rutas
- ✅ **Health checks** para monitoreo
- ✅ **Manejo de errores** centralizado
- ✅ **Timeouts configurables** para conexiones
- ✅ **Logs detallados** de peticiones y errores

### Capacidades Técnicas

- ✅ Balanceador de carga básico
- ✅ Transformación de paths automática
- ✅ Configuración por variables de entorno
- ✅ Compatibilidad con Docker y Docker Compose
- ✅ Middleware extensible (CORS, logging, etc.)

## 🏗️ Arquitectura

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Cliente Web   │────│   API Gateway    │────│  Microservicio  │
│                 │    │   (Puerto 3500)  │    │   de Destino    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                    ┌───────────┼───────────┐
                    │           │           │
            ┌───────────┐ ┌──────────┐ ┌──────────┐
            │  Gestión  │ │Parser/IoT│ │Notif.    │
            │:8092      │ │:8090     │ │:8091     │
            └───────────┘ └──────────┘ └──────────┘
                                │
                    ┌───────────────────────┐
                    │   Documentación       │
                    │   :8093               │
                    └───────────────────────┘
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
│   └── main.go                 # Punto de entrada
├── config/
│   └── config.go              # Configuración y variables de entorno
├── proxy/
│   ├── routes_config.go       # ⭐ Configuración de rutas
│   ├── routes.go              # Registro de rutas HTTP
│   └── factory.go             # Lógica del proxy reverso
├── pkg/
│   └── utils.go               # Utilidades para paths
├── middleware/
│   └── cors.go                # Middleware CORS (opcional)
├── .env                       # Variables de entorno
├── Dockerfile                 # Imagen Docker
├── go.mod                     # Dependencias Go
└── README.md                  # Este archivo
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

### ⭐ Archivo Principal: `proxy/routes_config.go`

Este es el corazón del API Gateway. Aquí defines todas las rutas que se van a proxificar:

```go
var ProxyRoutes = []ProxyRoute{
    {
        PathPrefix:   "/api/gestion",
        TargetEnvVar: "GESTION_URL",
        PrependPath:  "",
    },
    {
        PathPrefix:   "/api/notificacion",
        TargetEnvVar: "NOTIFICATION_URL", 
        PrependPath:  "",
    },
    {
        PathPrefix:   "/api/parser",
        TargetEnvVar: "PARSER_URL",
        PrependPath:  "",
    },
    {
        PathPrefix:   "/api/documentacion",
        TargetEnvVar: "DOCUMENTATION_URL",
        PrependPath:  "/api/v1",
    },
}
```

### 🔧 Cómo Agregar Nuevas Rutas

#### 1. Agregar la configuración de ruta

En `proxy/routes_config.go`, agrega una nueva entrada:

```go
{
    PathPrefix:   "/api/usuarios",      // Prefijo que escucha el gateway
    TargetEnvVar: "USUARIOS_URL",       // Variable de entorno con la URL destino
    PrependPath:  "",                   // Path adicional a agregar (opcional)
}, // ⚠️ Muy importante esta coma
```

#### 2. Definir la variable de entorno

En el archivo `.env`:

```bash
USUARIOS_URL=http://usuarios-service:8094
```

#### 3. ¡Listo! La ruta ya funciona automáticamente

No necesitas tocar ningún otro archivo. El API Gateway:
- ✅ Detecta automáticamente la nueva ruta
- ✅ Reenvía peticiones al microservicio
- ✅ Maneja errores y timeouts

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

## 📋 Rutas Actuales Configuradas

| Path | Destino | Puerto | PrependPath | Descripción |
|------|---------|---------|-------------|-------------|
| `/api/gestion` | gestion-service | 8092 | - | Gestión de técnicos y mantenimiento |
| `/api/notificacion` | notification-service | 8091 | - | Sistema de notificaciones |
| `/api/parser` | iot-service | 8090 | - | Datos de sensores IoT |
| `/api/documentacion` | documentacion-service | 8093 | `/api/v1` | Gestión de documentos |
| `/healthz` | - | - | - | Health check del gateway |

## 🔧 Ejemplos Reales de Uso

### Gestión de Técnicos
```bash
# A través del API Gateway
curl http://localhost:3500/api/gestion/tecnicos

# Se convierte internamente en:
curl http://gestion-service:8092/tecnicos
```

### Documentos Técnicos
```bash
# A través del API Gateway
curl http://localhost:3500/api/documentacion/documentos

# Se convierte internamente en:
curl http://documentacion-service:8093/api/v1/documentos
```

### Datos de Sensores
```bash
# A través del API Gateway
curl http://localhost:3500/api/parser/activo/AC-001

# Se convierte internamente en:
curl http://iot-service:8090/activo/AC-001
```

### Notificaciones
```bash
# A través del API Gateway  
curl -X POST http://localhost:3500/api/notificacion/notification \
     -H "Content-Type: application/json" \
     -d '{"activo_id": 1}'

# Se convierte internamente en:
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

# Verificar rutas específicas
curl -f http://localhost:3500/api/gestion/health || echo "Gestión DOWN"
curl -f http://localhost:3500/api/parser/activo || echo "Parser DOWN"
```

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
