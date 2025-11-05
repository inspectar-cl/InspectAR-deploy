# API Gateway - Enrutamiento Unificado

## 🎯 Resumen de Cambios

El API Gateway ahora soporta **proxy transparente con validación de autorización**, manteniendo compatibilidad con handlers personalizados.

### ¿Qué cambió?

- ✅ **Proxy transparente ahora puede tener validación de tokens y scopes**
- ✅ **Configuración centralizada en `handlers/requests_config.go`**
- ✅ **Sin cambios en rutas existentes** (100% compatible)
- ✅ **`main.go` simplificado** (solo un punto de registro)

---

## 🚀 Inicio Rápido

### Agregar Ruta Protegida de Proxy

**Ejemplo**: Proteger `/api/ML/anomalies/activo/:activo_id`

1. **Editar `handlers/requests_config.go`**
2. **Descomentar el ejemplo en `getProtectedProxyRoutes()`**:

```go
{
    Method:             "GET",
    Pattern:            "/api/ML/anomalies/activo/:activo_id",
    TargetEnvVar:       "ML_URL",
    PrependPath:        "",
    RequiredScopes:     userTypes("Tecnico", "Analista", "Admin"),
    RequiredValidation: "activo",
}
```

3. **Reiniciar el gateway**

**Resultado**:
- ✅ Requiere token JWT válido
- ✅ Requiere rol Tecnico, Analista, Admin o Root
- ✅ Valida que el usuario tenga acceso al activo especificado
- ✅ Redirige a `{ML_URL}/anomalies/activo/{activo_id}`

---

## 📖 Tipos de Rutas

### 1. Handler Personalizado
Para lógica compleja en Go.

```go
// En getSpecialHandlers()
{
    Method:  "POST",
    Pattern: "/api/crear-activo",
    Handler: CrearActivoHandler,
    Protected: true,
    RequiredScopes: rootOnly(),
}
```

### 2. Proxy Protegido Específico
Para ruta específica con validación.

```go
// En getProtectedProxyRoutes()
{
    Method:             "GET",
    Pattern:            "/api/ML/predictions/:id",
    TargetEnvVar:       "ML_URL",
    RequiredScopes:     userTypes("Tecnico"),
    RequiredValidation: "",
}
```

### 3. Proxy Transparente Genérico
Para TODO un prefijo sin validación.

```go
// En getProxyRoutes()
{
    PathPrefix:         "/api/ML",
    TargetEnvVar:       "ML_URL",
    Protected:          false,
}
```

---

## 🔒 Niveles de Protección

| Configuración | Validación | Caso de Uso |
|--------------|------------|-------------|
| `Protected: false` | Ninguna | Rutas públicas |
| `Protected: true` | Solo token | Info general del usuario |
| `+ RequiredScopes` | Token + rol | Acceso por rol |
| `+ RequiredValidation` | Token + rol + recurso | Acceso granular |

---

## 🛠️ Funciones Auxiliares

```go
// Root + roles especificados
userTypes("Tecnico", "Admin")
// → ["user-type:Root", "user-type:Tecnico", "user-type:Admin"]

// Solo Root
rootOnly()
// → ["user-type:Root"]
```

---

## 📝 Ejemplos Completos

### Proteger Endpoint de ML

```go
// En getProtectedProxyRoutes()
{
    Method:             "GET",
    Pattern:            "/api/ML/anomalies/activo/:activo_id",
    TargetEnvVar:       "ML_URL",
    PrependPath:        "",
    RequiredScopes:     userTypes("Tecnico", "Analista"),
    RequiredValidation: "activo", // Valida acceso al activo
}
```

### Proteger Lista de Edificios

```go
{
    Method:             "GET",
    Pattern:            "/api/gestion/edificios",
    TargetEnvVar:       "GESTION_URL",
    PrependPath:        "",
    RequiredScopes:     userTypes("Admin"),
    RequiredValidation: "",
}
```

### Proxy Genérico (Sin Protección)

```go
// En getProxyRoutes()
{
    PathPrefix:         "/api/ML",
    TargetEnvVar:       "ML_URL",
    PrependPath:        "",
    Protected:          false,
}
// Todas las rutas /api/ML/* son públicas
```

---

## ⚙️ Variables de Entorno

Asegúrate de configurar:

```bash
export GESTION_URL=http://localhost:4000
export ML_URL=http://localhost:5000
export PARSER_URL=http://localhost:4001
export DOCUMENTATION_URL=http://localhost:4002
export NOTIFICATION_URL=http://localhost:4003
export OAUTH2_URL=http://localhost:4004
```

---

## 🧪 Testing

### Verificar Logs de Registro

Al iniciar, verás:
```
Registered special handler: POST /api/crear-activo (protected)
Registered protected proxy: GET /api/ML/anomalies/activo/:activo_id → http://ml:5000 (protected)
Registered proxy route: ANY /api/ML/*proxyPath → http://ml:5000 (public)
```

### Probar Autenticación

```bash
# Sin token (401)
curl http://localhost:8080/api/ML/anomalies/activo/123

# Con token válido (200)
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://localhost:8080/api/ML/anomalies/activo/123
```

---

## 🎓 Migración

### Antes
```go
// main.go
proxy.RegisterRoutes(r)          // Sin validación
handlers.RegisterSpecialRoutes(r) // Con validación
```

### Ahora
```go
// main.go
handlers.RegisterSpecialRoutes(r) // Todo unificado
```

**No hay cambios necesarios en rutas existentes.**

---

## 📚 Documentación Completa

Para más detalles, consulta:
- **[PROXY_ROUTING_GUIDE.md](./PROXY_ROUTING_GUIDE.md)** - Guía completa
- **[handlers/requests_config.go](./handlers/requests_config.go)** - Configuración

---

## ⚠️ Importante

1. **Root siempre tiene acceso** a todas las rutas protegidas
2. **ValidationMiddleware** busca `:id_edificio` o `:id_activo` en la ruta
3. **Orden importa**: Rutas específicas antes que genéricas

---

## 📞 Soporte

¿Preguntas? Revisa:
- Código: `handlers/requests_config.go`
- Middlewares: `api/middleware/`
- Guía completa: `PROXY_ROUTING_GUIDE.md`
