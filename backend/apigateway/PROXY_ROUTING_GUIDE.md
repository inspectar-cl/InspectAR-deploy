# Guía de Enrutamiento Unificado - API Gateway

## 📋 Resumen

El API Gateway ahora soporta **dos modos de operación unificados** bajo una misma configuración centralizada en `handlers/requests_config.go`:

1. **Handlers Personalizados**: Lógica de negocio específica en Go
2. **Proxy Transparente con Validación**: Redirección directa a microservicios con control de autorización

**Ambos modos comparten**:
- ✅ Validación de tokens JWT
- ✅ Control de scopes por rol
- ✅ Validación de acceso a recursos (edificios/activos)
- ✅ Configuración centralizada

---

## 🏗️ Arquitectura

### Antes (Duplicado)
```
main.go
├── proxy.RegisterRoutes()        → Sin validación
└── handlers.RegisterSpecialRoutes() → Con validación
```

### Ahora (Unificado)
```
main.go
└── handlers.RegisterSpecialRoutes()
    ├── Handlers personalizados (con validación)
    ├── Proxy protegido específico (con validación)
    └── Proxy transparente genérico (sin validación por defecto)
```

---

## 📂 Estructura de Archivos

### Archivos Principales

| Archivo | Propósito |
|---------|-----------|
| `handlers/requests_config.go` | **Configuración centralizada** de todas las rutas |
| `handlers/requests.go` | Implementación de handlers personalizados |
| `proxy/factory.go` | Creación de handlers de proxy con middlewares |
| `proxy/routes_config.go` | Definición del struct `ProxyRoute` |
| `cmd/main.go` | Punto de entrada simplificado |

### Archivos Deprecados

| Archivo | Estado |
|---------|--------|
| `proxy/routes.go` | ❌ Renombrado a `.deprecated` (no se usa) |

---

## ⚙️ Configuración de Rutas

### 1. Handlers Personalizados

Para rutas con lógica personalizada en Go:

```go
// En handlers/requests_config.go → getSpecialHandlers()

{
    Method:  "POST",
    Pattern: "/api/crear-activo",
    Handler: CrearActivoHandler,  // Función en requests.go
    Protected: true,
    RequiredScopes: rootOnly(),   // Solo Root
    RequiredValidation: "",
}
```

**Casos de uso**:
- Lógica compleja de negocio
- Transformación de datos
- Agregación de múltiples servicios
- Validaciones específicas

---

### 2. Proxy Transparente Genérico

Para redirigir TODO un prefijo de ruta sin protección:

```go
// En handlers/requests_config.go → getProxyRoutes()

{
    PathPrefix:         "/api/ML",
    TargetEnvVar:       "ML_URL",
    PrependPath:        "",
    Protected:          false,  // Sin validación
    RequiredScopes:     nil,
    RequiredValidation: "",
}
```

**Resultado**:
- `/api/ML/anomalies` → `{ML_URL}/anomalies`
- `/api/ML/predictions/activo/123` → `{ML_URL}/predictions/activo/123`
- **Sin validación de token ni scopes**

---

### 3. Proxy Protegido Específico

Para rutas específicas con validación de autorización:

```go
// En handlers/requests_config.go → getProtectedProxyRoutes()

{
    Method:             "GET",
    Pattern:            "/api/ML/anomalies/activo/:activo_id",
    TargetEnvVar:       "ML_URL",
    PrependPath:        "",
    RequiredScopes:     userTypes("Tecnico", "Analista", "Admin"),
    RequiredValidation: "activo",  // Verifica acceso al activo
}
```

**Resultado**:
- ✅ Valida token JWT
- ✅ Verifica que el usuario tenga rol Tecnico, Analista, Admin o Root
- ✅ Verifica que el usuario tenga acceso al `activo_id` especificado
- ✅ Si todo OK → redirige a `{ML_URL}/anomalies/activo/{activo_id}`

---

## 🔒 Niveles de Protección

### Nivel 0: Sin Protección
```go
Protected: false,
RequiredScopes: nil,
RequiredValidation: "",
```
- ✅ Acceso público
- ❌ No valida token
- **Ejemplo**: `/api/login`, rutas de proxy genéricas

---

### Nivel 1: Solo Autenticación
```go
Protected: true,
RequiredScopes: nil,  // O array vacío
RequiredValidation: "",
```
- ✅ Requiere token válido
- ❌ No valida roles específicos
- **Ejemplo**: Endpoints de información general del usuario

---

### Nivel 2: Autenticación + Scopes
```go
Protected: true,
RequiredScopes: userTypes("Tecnico", "Admin"),
RequiredValidation: "",
```
- ✅ Requiere token válido
- ✅ Requiere rol Tecnico, Admin o Root
- ❌ No valida acceso a recursos específicos

---

### Nivel 3: Autenticación + Scopes + Validación de Recurso
```go
Protected: true,
RequiredScopes: userTypes("Tecnico", "Admin"),
RequiredValidation: "activo",  // o "edificio"
```
- ✅ Requiere token válido
- ✅ Requiere rol específico
- ✅ Verifica acceso al `:id_activo` o `:id_edificio` en la ruta
- **Ejemplo**: Solo técnicos con acceso al edificio del activo

---

## 🎯 Ejemplos Prácticos

### Ejemplo 1: Endpoint de ML con Validación de Activo

**Requisito**: Solo técnicos y analistas pueden ver anomalías de activos a los que tienen acceso.

```go
// En getProtectedProxyRoutes()
{
    Method:             "GET",
    Pattern:            "/api/ML/anomalies/activo/:activo_id",
    TargetEnvVar:       "ML_URL",
    PrependPath:        "",
    RequiredScopes:     userTypes("Tecnico", "Analista"),
    RequiredValidation: "activo",
}
```

**Flujo**:
1. Usuario hace `GET /api/ML/anomalies/activo/123`
2. `AuthMiddleware` valida el token JWT
3. `ScopeMiddleware` verifica que el usuario sea Tecnico, Analista o Root
4. `ValidationMiddleware` consulta al servicio de gestión si el usuario tiene acceso al activo 123
5. Si todo OK → proxy redirige a `{ML_URL}/anomalies/activo/123`

---

### Ejemplo 2: Endpoint de Gestión Solo para Admins

**Requisito**: Solo administradores pueden listar todos los edificios.

```go
// En getProtectedProxyRoutes()
{
    Method:             "GET",
    Pattern:            "/api/gestion/edificios",
    TargetEnvVar:       "GESTION_URL",
    PrependPath:        "",
    RequiredScopes:     userTypes("Admin"),
    RequiredValidation: "",  // No valida edificio específico
}
```

**Flujo**:
1. Usuario hace `GET /api/gestion/edificios`
2. `AuthMiddleware` valida el token
3. `ScopeMiddleware` verifica que sea Admin o Root
4. Si OK → proxy redirige a `{GESTION_URL}/edificios`

---

### Ejemplo 3: Proxy Genérico sin Protección

**Requisito**: El servicio ML maneja su propia autenticación interna.

```go
// En getProxyRoutes()
{
    PathPrefix:         "/api/ML",
    TargetEnvVar:       "ML_URL",
    PrependPath:        "",
    Protected:          false,
    RequiredScopes:     nil,
    RequiredValidation: "",
}
```

**Flujo**:
- Cualquier ruta `/api/ML/*` se redirige a `{ML_URL}/*` sin validación
- El microservicio ML es responsable de su propia seguridad

---

## 🛠️ Funciones Auxiliares

### `userTypes(...roles)`
Agrega el prefijo `user-type:` y **siempre incluye Root**.

```go
userTypes("Tecnico", "Admin")
// → ["user-type:Root", "user-type:Tecnico", "user-type:Admin"]
```

### `rootOnly()`
Solo usuarios con rol Root.

```go
rootOnly()
// → ["user-type:Root"]
```

---

## 🔄 Orden de Registro

**Importante**: El orden de registro determina la prioridad de las rutas.

```go
func RegisterSpecialRoutes(r *gin.Engine) {
    // 1. Endpoint de salud
    r.GET("/healthz", ...)

    // 2. Handlers personalizados (más específicos)
    for _, handler := range getSpecialHandlers() { ... }

    // 3. Proxy protegido específico (prioridad sobre genérico)
    for _, route := range getProtectedProxyRoutes() { ... }

    // 4. Proxy transparente genérico (menos específico)
    for _, route := range getProxyRoutes() { ... }
}
```

**Regla**: Rutas más específicas primero, genéricas al final.

---

## 📝 Agregar Nueva Ruta Protegida

### Paso 1: Decidir el Tipo

| ¿Necesitas lógica personalizada? | ¿Validación de autorización? | Tipo Recomendado |
|----------------------------------|------------------------------|-------------------|
| ✅ Sí | ✅ Sí | Handler Personalizado |
| ❌ No | ✅ Sí | Proxy Protegido Específico |
| ❌ No | ❌ No | Proxy Transparente Genérico |

### Paso 2: Agregar Configuración

**Para Handler Personalizado**:
```go
// En getSpecialHandlers()
{
    Method:  "POST",
    Pattern: "/api/mi-endpoint/:id",
    Handler: MiHandlerFunc,  // Implementar en requests.go
    Protected: true,
    RequiredScopes: userTypes("Tecnico"),
    RequiredValidation: "activo",
}
```

**Para Proxy Protegido**:
```go
// En getProtectedProxyRoutes()
{
    Method:             "GET",
    Pattern:            "/api/ML/mi-ruta/:id",
    TargetEnvVar:       "ML_URL",
    PrependPath:        "",
    RequiredScopes:     userTypes("Analista"),
    RequiredValidation: "",
}
```

### Paso 3: Reiniciar el Gateway
```bash
cd backend/apigateway
go run cmd/main.go
```

---

## 🧪 Testing

### Verificar Rutas Registradas

Al iniciar el gateway, verás logs como:
```
Registered special handler: POST /api/crear-activo (protected, scopes: user-type:Root)
Registered protected proxy: GET /api/ML/anomalies/activo/:activo_id → http://ml:5000 (protected, scopes: user-type:Root, user-type:Tecnico)
Registered proxy route: ANY /api/ML/*proxyPath → http://ml:5000 (public)
```

### Probar Autenticación

```bash
# Sin token (debe fallar)
curl http://localhost:8080/api/ML/anomalies/activo/123

# Con token inválido (debe fallar)
curl -H "Authorization: Bearer invalid_token" \
     http://localhost:8080/api/ML/anomalies/activo/123

# Con token válido (debe funcionar)
curl -H "Authorization: Bearer YOUR_VALID_TOKEN" \
     http://localhost:8080/api/ML/anomalies/activo/123
```

---

## ⚠️ Consideraciones Importantes

### 1. Variables de Entorno
Asegúrate de configurar las URLs de los microservicios:
```bash
export GESTION_URL=http://localhost:4000
export ML_URL=http://localhost:5000
export PARSER_URL=http://localhost:4001
export DOCUMENTATION_URL=http://localhost:4002
export NOTIFICATION_URL=http://localhost:4003
export OAUTH2_URL=http://localhost:4004
```

### 2. Root Siempre Tiene Acceso
El usuario Root (scope `user-type:Root`) **siempre tiene acceso** a todas las rutas protegidas, incluso sin estar explícito en `RequiredScopes`.

### 3. Validación de Recursos
El `ValidationMiddleware` busca `:id_edificio` o `:id_activo` en los parámetros de la ruta. Asegúrate de usar estos nombres exactos.

### 4. Orden de Middlewares
El orden es importante:
1. `AuthMiddleware` → Valida token
2. `ScopeMiddleware` → Valida rol
3. `ValidationMiddleware` → Valida acceso al recurso
4. Handler/Proxy → Ejecuta la lógica

---

## 🎓 Migración desde el Sistema Anterior

### Cambios Necesarios

1. **main.go**: Eliminar `proxy.RegisterRoutes(r)`
2. **Configuración**: Mover configuración de proxy a `handlers/requests_config.go`
3. **Rutas protegidas**: Usar `getProtectedProxyRoutes()` en lugar de crear handlers personalizados

### Compatibilidad

✅ **Todas las rutas existentes funcionan igual**
- Handlers personalizados no cambiaron
- Proxy genérico funciona como antes
- Solo se agregó capacidad de proteger rutas de proxy

---

## 📚 Referencias

- **Middlewares**: `api/middleware/`
  - `auth_middleware.go`: Validación JWT
  - `scope_middleware.go`: Validación de roles
  - `validation_middleware.go`: Validación de acceso a recursos

- **Proxy**: `proxy/`
  - `factory.go`: Creación de handlers de proxy
  - `routes_config.go`: Definición de structs

- **Handlers**: `handlers/`
  - `requests_config.go`: Configuración centralizada
  - `requests.go`: Implementación de handlers

---

## 🚀 Beneficios del Nuevo Sistema

| Beneficio | Descripción |
|-----------|-------------|
| **Centralización** | Toda la configuración en un solo archivo |
| **Seguridad** | Proxy transparente puede tener validación de autorización |
| **Flexibilidad** | Mezclar handlers personalizados y proxy según necesidad |
| **Mantenibilidad** | Más fácil de entender y modificar |
| **Escalabilidad** | Agregar nuevos microservicios es trivial |

---

## ❓ Preguntas Frecuentes

### ¿Puedo tener un microservicio con rutas públicas y protegidas?

**Sí**. Usa `getProtectedProxyRoutes()` para rutas específicas protegidas y `getProxyRoutes()` para el catch-all.

```go
// Rutas específicas protegidas primero
{
    Pattern: "/api/gestion/admin/users",
    RequiredScopes: rootOnly(),
}

// Catch-all sin protección después
{
    PathPrefix: "/api/gestion",
    Protected: false,
}
```

### ¿Cómo deshabilito temporalmente la protección?

Cambia `Protected: true` a `Protected: false` y reinicia el gateway.

### ¿Root puede acceder a todo sin estar en RequiredScopes?

**Sí**. `userTypes()` automáticamente incluye Root, y `ValidationMiddleware` omite validación para Root.

---

## 📞 Soporte

Para más información, revisa:
- Código en `handlers/requests_config.go`
- Middlewares en `api/middleware/`
- Tests en `tests/`
