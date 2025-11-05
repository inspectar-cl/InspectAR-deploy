# Arquitectura del API Gateway - Diagrama

## 📊 Flujo de Solicitudes

```
┌─────────────────────────────────────────────────────────────────┐
│                         Cliente / Frontend                       │
└───────────────────────────────┬─────────────────────────────────┘
                                │
                                │ HTTP Request
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                          API GATEWAY                             │
│                        (Puerto 8080)                             │
└───────────────────────────────┬─────────────────────────────────┘
                                │
                ┌───────────────┴───────────────┐
                │                               │
                ▼                               ▼
    ┌───────────────────┐           ┌───────────────────┐
    │ CORS Middleware   │           │  Logger / Recovery │
    └───────────────────┘           └───────────────────┘
                │                               │
                └───────────────┬───────────────┘
                                │
                                ▼
        ┌───────────────────────────────────────────┐
        │      RegisterSpecialRoutes()              │
        │   (handlers/requests_config.go)           │
        └───────────────────────────────────────────┘
                                │
        ┌───────────────────────┼───────────────────────┐
        │                       │                       │
        ▼                       ▼                       ▼
┌───────────────┐   ┌─────────────────────┐   ┌──────────────────┐
│   Healthz     │   │ Special Handlers    │   │  Proxy Routes    │
│   /healthz    │   │ (Personalizados)    │   │ (Transparente)   │
└───────────────┘   └─────────────────────┘   └──────────────────┘
                                │                       │
                                │                       │
                    ┌───────────┴───────┐   ┌──────────┴──────────┐
                    │                   │   │                     │
                    ▼                   ▼   ▼                     ▼
            ┌──────────────┐   ┌──────────────┐   ┌──────────────────┐
            │   Handler    │   │   Protected  │   │   Generic Proxy  │
            │Personalizado │   │  Proxy Route │   │   (sin validar)  │
            │              │   │              │   │                  │
            │ (requests.go)│   │              │   │                  │
            └──────────────┘   └──────────────┘   └──────────────────┘
                    │                   │                   │
                    │                   │                   │
        ┌───────────┼───────────────────┼───────────────────┘
        │           │                   │
        │           │         ┌─────────┴─────────┐
        │           │         │                   │
        │           │         ▼                   ▼
        │           │  ┌──────────────┐   ┌──────────────┐
        │           │  │ Auth         │   │  Scope       │
        │           │  │ Middleware   │   │  Middleware  │
        │           │  │              │   │              │
        │           │  │ (Valida JWT) │   │ (Valida Rol) │
        │           │  └──────────────┘   └──────────────┘
        │           │         │                   │
        │           │         └─────────┬─────────┘
        │           │                   │
        │           │                   ▼
        │           │            ┌──────────────┐
        │           │            │ Validation   │
        │           │            │ Middleware   │
        │           │            │ (si aplica)  │
        │           │            │              │
        │           │            │ Valida acceso│
        │           │            │ a recurso    │
        │           │            └──────────────┘
        │           │                   │
        │           │                   │
        │           └───────────────────┤
        │                               │
        ▼                               ▼
┌──────────────────┐           ┌──────────────────┐
│ Lógica en Go     │           │ ReverseProxy     │
│ (requests.go)    │           │ (factory.go)     │
└──────────────────┘           └──────────────────┘
        │                               │
        │                               │
        └───────────────┬───────────────┘
                        │
                        ▼
        ┌───────────────────────────────────┐
        │         Microservicios            │
        │                                   │
        │  ┌─────────┐  ┌─────────┐        │
        │  │ Gestión │  │   ML    │  ...   │
        │  │ :4000   │  │  :5000  │        │
        │  └─────────┘  └─────────┘        │
        └───────────────────────────────────┘
```

---

## 🔍 Tipos de Rutas en Detalle

### 1. Handler Personalizado
```
Cliente → Gateway → [CORS] → [Logger] → [Auth] → [Scope] → [Validation]
                                                               ↓
                                                    Handler Personalizado (Go)
                                                               ↓
                                                    Múltiples Microservicios
                                                               ↓
                                                    Transformación de Datos
                                                               ↓
                                                         Respuesta
```

**Ejemplo**: `POST /api/crear-activo`
- ✅ Validación de JWT
- ✅ Verificación de scope Root
- ✅ Lógica personalizada en `CrearActivoHandler`
- ✅ Llamada a microservicio de gestión
- ✅ Respuesta procesada

---

### 2. Proxy Protegido Específico
```
Cliente → Gateway → [CORS] → [Logger] → [Auth] → [Scope] → [Validation]
                                                               ↓
                                                        ReverseProxy
                                                               ↓
                                                    Microservicio (directo)
                                                               ↓
                                                         Respuesta
```

**Ejemplo**: `GET /api/ML/anomalies/activo/:activo_id`
- ✅ Validación de JWT
- ✅ Verificación de scope (Tecnico/Analista/Root)
- ✅ Validación de acceso al activo
- ✅ Proxy directo a ML service
- ✅ Sin transformación

---

### 3. Proxy Transparente Genérico
```
Cliente → Gateway → [CORS] → [Logger] → ReverseProxy
                                              ↓
                                      Microservicio
                                              ↓
                                         Respuesta
```

**Ejemplo**: `ANY /api/ML/*`
- ❌ Sin validación de JWT
- ❌ Sin validación de scopes
- ✅ Proxy directo a cualquier ruta bajo `/api/ML`
- ✅ Útil para servicios con auth propia

---

## 🗂️ Estructura de Configuración

```
handlers/requests_config.go
│
├── userTypes()               → Helper: Agrega prefijo + Root
├── rootOnly()                → Helper: Solo Root
│
├── getSpecialHandlers()      → 📝 Handlers personalizados
│   └── []SpecialHandler{
│       ├── LoginHandler
│       ├── CrearActivoHandler
│       ├── EditarActivoHandler
│       └── ...
│   }
│
├── getProtectedProxyRoutes() → 🔒 Proxy con validación
│   └── []struct{
│       ├── /api/ML/anomalies/activo/:activo_id
│       ├── /api/gestion/edificios
│       └── ... (comentados por defecto)
│   }
│
├── getProxyRoutes()          → 🌐 Proxy sin validación
│   └── []ProxyRoute{
│       ├── /api/gestion/*
│       ├── /api/ML/*
│       ├── /api/parser/*
│       └── ...
│   }
│
└── RegisterSpecialRoutes()   → 🎯 Registro unificado
    ├── 1. Healthz
    ├── 2. Special Handlers
    ├── 3. Protected Proxy
    └── 4. Generic Proxy
```

---

## 🔐 Cadena de Middlewares

### Para Handler Personalizado con Validación Completa
```
Request
  ↓
[AuthMiddleware]
  ├─ Valida Bearer token
  ├─ Extrae claims (email, scope)
  └─ Guarda en contexto Gin
  ↓
[ScopeMiddleware]
  ├─ Lee scope del contexto
  ├─ Verifica si contiene rol requerido
  └─ 403 si no autorizado
  ↓
[ValidationMiddleware]
  ├─ Extrae :id_edificio o :id_activo
  ├─ Consulta a Gestión: ¿usuario tiene acceso?
  └─ 403 si no tiene acceso
  ↓
[Handler]
  └─ Ejecuta lógica personalizada
```

### Para Proxy Protegido
```
Request
  ↓
[AuthMiddleware]
  ↓
[ScopeMiddleware]
  ↓
[ValidationMiddleware] (opcional)
  ↓
[ReverseProxy]
  └─ Reescribe URL y redirige
```

### Para Proxy Genérico
```
Request
  ↓
[ReverseProxy]
  └─ Redirige directamente
```

---

## 📋 Matriz de Decisión

| ¿Necesitas... | Solución | Configuración |
|--------------|----------|---------------|
| Lógica compleja en Go | Handler Personalizado | `getSpecialHandlers()` |
| Solo redirección + auth | Proxy Protegido | `getProtectedProxyRoutes()` |
| Redirección sin auth | Proxy Genérico | `getProxyRoutes()` |
| Validar acceso a recurso | + `RequiredValidation` | `"edificio"` o `"activo"` |
| Validar rol específico | + `RequiredScopes` | `userTypes("Tecnico")` |
| Solo Root | + `RequiredScopes` | `rootOnly()` |

---

## 🚦 Flujo de Prioridad de Rutas

```
1. Healthz (/healthz)
   ↓
2. Handlers Personalizados (más específicos)
   - /api/crear-activo
   - /api/editar-activo/:id_activo
   - /api/login
   - ...
   ↓
3. Proxy Protegido (específico)
   - /api/ML/anomalies/activo/:activo_id
   - /api/gestion/edificios
   - ...
   ↓
4. Proxy Genérico (catch-all)
   - /api/ML/*
   - /api/gestion/*
   - /api/parser/*
   - ...
```

**Importante**: Gin evalúa rutas en orden de registro. Las más específicas se registran primero.

---

## 🎨 Ejemplo Completo de Solicitud

### Solicitud: `GET /api/ML/anomalies/activo/123`

#### Con Proxy Protegido Configurado:
```
1. Cliente envía:
   GET /api/ML/anomalies/activo/123
   Authorization: Bearer eyJhbGci...

2. Gateway recibe → RegisterSpecialRoutes
   ↓
3. Coincide con Protected Proxy:
   Pattern: /api/ML/anomalies/activo/:activo_id
   ↓
4. AuthMiddleware:
   - Valida token → OK
   - Extrae email: "tecnico@example.com"
   - Extrae scope: "user-type:Tecnico"
   ↓
5. ScopeMiddleware:
   - Required: ["user-type:Root", "user-type:Tecnico", "user-type:Analista"]
   - Usuario tiene: "user-type:Tecnico"
   - Resultado: ✅ AUTORIZADO
   ↓
6. ValidationMiddleware:
   - Tipo: "activo"
   - ID: 123 (de :activo_id)
   - Email: "tecnico@example.com"
   - Consulta: GET {GESTION_URL}/usuarios/tecnico@example.com/activo/123/acceso
   - Respuesta: {"tiene_acceso": true}
   - Resultado: ✅ AUTORIZADO
   ↓
7. ReverseProxy:
   - Strip prefix: /api/ML
   - Ruta relativa: /anomalies/activo/123
   - Target: {ML_URL}/anomalies/activo/123
   - Forward request
   ↓
8. Microservicio ML responde:
   {"anomalies": [...]}
   ↓
9. Gateway reenvía respuesta al cliente
```

#### Sin Proxy Protegido (Catch-all):
```
1. Cliente envía:
   GET /api/ML/anomalies/activo/123

2. Gateway recibe → RegisterSpecialRoutes
   ↓
3. Coincide con Generic Proxy:
   PathPrefix: /api/ML
   Protected: false
   ↓
4. ReverseProxy (sin middlewares):
   - Forward directo a {ML_URL}/anomalies/activo/123
   ↓
5. Microservicio ML responde
```

---

## 📊 Comparación Visual

### Antes de la Restructuración
```
┌────────────────────────────────────┐
│           main.go                  │
│                                    │
│  proxy.RegisterRoutes(r)           │  → Sin validación
│  handlers.RegisterSpecialRoutes(r) │  → Con validación
└────────────────────────────────────┘
                ↓
         PROBLEMA: Duplicación
```

### Después de la Restructuración
```
┌────────────────────────────────────┐
│           main.go                  │
│                                    │
│  handlers.RegisterSpecialRoutes(r) │  → TODO centralizado
└────────────────────────────────────┘
                ↓
         SOLUCIÓN: Unificado
         ✅ Proxy puede tener validación
         ✅ Configuración centralizada
         ✅ Mismo flujo de seguridad
```

---

## 🎓 Resumen Visual

```
┌─────────────────────────────────────────────────────────┐
│              API GATEWAY - Enrutamiento                 │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌──────────────┐   ┌──────────────┐   ┌───────────┐  │
│  │   Handler    │   │   Protected  │   │  Generic  │  │
│  │ Personalizado│   │    Proxy     │   │   Proxy   │  │
│  └──────────────┘   └──────────────┘   └───────────┘  │
│        │                   │                   │       │
│        ├───────────────────┴───────────────────┤       │
│        │         Middlewares Unificados        │       │
│        │  (Auth → Scope → Validation)         │       │
│        └───────────────────────────────────────┘       │
│                                                         │
│  📝 Configuración: handlers/requests_config.go         │
│  🔧 Factory: proxy/factory.go                          │
│  🛡️ Middlewares: api/middleware/                      │
└─────────────────────────────────────────────────────────┘
```
