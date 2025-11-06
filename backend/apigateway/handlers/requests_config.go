package handlers

import (
    "log"
    "os"
    "strings"
    "github.com/gin-gonic/gin"
    "apigateway/api/middleware"
    "apigateway/proxy"
)

// SpecialHandler representa un handler especial para Gin
type SpecialHandler struct {
    Method              string
    Pattern             string
    Handler             gin.HandlerFunc
    Protected           bool
    RequiredScopes      []string // Root siempre incluido automáticamente. Otros roles: Analista, Tecnico, Residente, Admin
    RequiredValidation  string   // Puede ser "edificio", "activo" o vacío si no se requiere validación.
}

// userTypes agrega el prefijo "user-type:" a cada rol especificado
// y automáticamente incluye "Root" para acceso total
func userTypes(roles ...string) []string {
    // Siempre incluir Root al inicio
    scopes := make([]string, 0, len(roles)+1)
    scopes = append(scopes, "user-type:Root")
    
    // Agregar los demás roles especificados
    for _, role := range roles {
        scopes = append(scopes, "user-type:"+role)
    }
    return scopes
}

// rootOnly retorna un scope exclusivo para Root
func rootOnly() []string {
    return []string{"user-type:Root"}
}

// getProxyRoutes devuelve todas las rutas de proxy transparente con su configuración de seguridad
// IMPORTANTE: Solo habilitar rutas que NO tengan handlers personalizados específicos registrados
// Si tienes handlers específicos en /api/gestion/*, NO habilites el catch-all de /api/gestion
func getProxyRoutes() []proxy.ProxyRoute {
    return []proxy.ProxyRoute{
        // NOTA: /api/gestion está deshabilitado porque tenemos handlers personalizados
        // Si necesitas proxy para gestion, usa getProtectedProxyRoutes() para rutas específicas
        // {
        //     PathPrefix:         "/api/gestion",
        //     TargetEnvVar:       "GESTION_URL",
        //     PrependPath:        "",
        //     Protected:          false,
        //     RequiredScopes:     nil,
        //     RequiredValidation: "",
        // },
        // {
        //     PathPrefix:         "/api/notificacion",
        //     TargetEnvVar:       "NOTIFICATION_URL",
        //     PrependPath:        "",
        //     Protected:          false,
        //     RequiredScopes:     nil,
        //     RequiredValidation: "",
        // },
        {
            PathPrefix:         "/api/parser",
            TargetEnvVar:       "PARSER_URL",
            PrependPath:        "",
            Protected:          false,
            RequiredScopes:     nil,
            RequiredValidation: "",
        },
        {
            PathPrefix:         "/api/documentacion",
            TargetEnvVar:       "DOCUMENTATION_URL",
            PrependPath:        "/api/v1",
            Protected:          false,
            RequiredScopes:     nil,
            RequiredValidation: "",
        },
        // NOTA: /api/ML está deshabilitado porque tenemos rutas específicas protegidas
        // Si necesitas más rutas de ML, agrégalas en getProtectedProxyRoutes()
        // {
        //     PathPrefix:         "/api/ML",
        //     TargetEnvVar:       "ML_URL",
        //     PrependPath:        "",
        //     Protected:          false,
        //     RequiredScopes:     nil,
        //     RequiredValidation: "",
        // },
    }
}

func getProtectedProxyRoutes() []struct {
    Method             string
    Pattern            string
    TargetEnvVar       string
    PrependPath        string
    RequiredScopes     []string
    RequiredValidation string
} {
    return []struct {
        Method             string
        Pattern            string
        TargetEnvVar       string
        PrependPath        string
        RequiredScopes     []string
        RequiredValidation string
    }{
        // Rutas ML protegidas Sprint 3
        {
            Method:             "GET",
            Pattern:            "/api/ML/anomalies/activo/:activo_id",
            TargetEnvVar:       "ML_URL",
            PrependPath:        "",
            RequiredScopes:     userTypes("Tecnico", "Analista", "Admin"),
            RequiredValidation: "activo",
        },
        {
            Method:             "GET",
            Pattern:            "/api/ML/anomalies/sensor/:sensor_id",
            TargetEnvVar:       "ML_URL",
            PrependPath:        "",
            RequiredScopes:     userTypes("Tecnico", "Analista", "Admin"),
            RequiredValidation: "",
        },
        {
            Method:             "GET",
            Pattern:            "/api/ML/anomalies/:anomaly_id",
            TargetEnvVar:       "ML_URL",
            PrependPath:        "",
            RequiredScopes:     userTypes("Tecnico", "Analista", "Admin"),
            RequiredValidation: "",
        },
        // Rutas Gestion protegidas Sprint 3
        {
            Method:             "GET",
            Pattern:            "/api/gestion/edificios",
            TargetEnvVar:       "GESTION_URL",
            PrependPath:        "",
            RequiredScopes:     rootOnly(),
            RequiredValidation: "",
        },
        {
            Method:             "GET",
            Pattern:            "/api/gestion/tecnicos",
            TargetEnvVar:       "GESTION_URL",
            PrependPath:        "",
            RequiredScopes:     rootOnly(),
            RequiredValidation: "",
        },
        {
            Method:             "PUT",
            Pattern:            "/api/gestion/tecnicos/:id_tecnico",
            TargetEnvVar:       "GESTION_URL",
            PrependPath:        "",
            RequiredScopes:     rootOnly(),
            RequiredValidation: "",
        },
        {
            Method:             "DELETE",
            Pattern:            "/api/gestion/tecnicos/:id_tecnico",
            TargetEnvVar:       "GESTION_URL",
            PrependPath:        "",
            RequiredScopes:     rootOnly(),
            RequiredValidation: "",
        },
        {
            Method:             "POST",
            Pattern:            "/api/gestion/tecnicos",
            TargetEnvVar:       "GESTION_URL",
            PrependPath:        "",
            RequiredScopes:     rootOnly(),
            RequiredValidation: "",
        },
        {
            Method:             "GET",
            Pattern:            "/api/gestion/tecnicos/:id_tecnico",
            TargetEnvVar:       "GESTION_URL",
            PrependPath:        "",
            RequiredScopes:     rootOnly(),
            RequiredValidation: "",
        },
        // Rutas Notificaciones protegidas Sprint 3
        {
            Method:             "GET",
            Pattern:            "/api/notification/tickets",
            TargetEnvVar:       "NOTIFICATION_URL",
            PrependPath:        "",
            RequiredScopes:     rootOnly(),
            RequiredValidation: "",
        },
    
    }
}

// getSpecialHandlers devuelve todos los handlers especiales
func getSpecialHandlers() []SpecialHandler {
    return []SpecialHandler{
        {
            Method: "GET",
            Pattern: "/api/generar-reporte/:id_reporte",
            Handler: GenerarReporte,   
            Protected: true,
            RequiredScopes: userTypes("Tecnico", "Residente", "Admin"),
        },
        { 
            Method:  "POST",
            Pattern: "/api/login",
            Handler: LoginHandler,
            Protected: false,
            RequiredScopes: nil, // No requiere scopes ya que es público
        },
        {
            Method:  "POST",
            Pattern: "/api/logout",
            Handler: LogoutHandler,
            Protected: false,
            RequiredScopes: nil,
        },
        { 
            Method:  "POST",
            Pattern: "/api/refresh",
            Handler: RefreshHandler,
            Protected: false,
            RequiredScopes: nil,
        },
        { 
            Method:  "GET",
            Pattern: "/api/rol",
            Handler: ObtenerRol,
            Protected: true,
            RequiredScopes: userTypes("Tecnico", "Residente", "Analista", "Admin"),
        },
        // A partir de acá serían los handlers para el sprint 2
        { 
            Method:  "GET",
            Pattern: "/api/obtener-activos-id/:id_edificio",
            Handler: ObtenerActivosPorEdificio,
            Protected: true,
            RequiredScopes: userTypes("Residente", "Tecnico", "Analista", "Admin"),
            RequiredValidation: "edificio",
        },
        {
            Method:  "GET",
            Pattern: "/api/obtener-contactos-id/:id_edificio",
            Handler: ObtenerContactosPorEdificio,
            Protected: true,
            RequiredScopes: userTypes("Residente", "Tecnico", "Analista", "Admin"),
            RequiredValidation: "edificio",
        },
        { 
            Method:  "GET",
            Pattern: "/api/foro-edificio/:id_edificio",
            Handler: ForoEdificioHandler,
            Protected: true,
            RequiredScopes: userTypes("Residente", "Admin"),
            RequiredValidation: "edificio",
        },
        { 
            Method:  "POST",
            Pattern: "/api/publicacion-foro-id/:id_edificio",
            Handler: PublicacionForoHandler,
            Protected: true,
            RequiredScopes: userTypes("Residente", "Admin"),
            RequiredValidation: "edificio",
        },
        { 
            Method:  "POST",
            Pattern: "/api/comentarios-foro-id/:id_publicacion",
            Handler: ComentariosForoHandler,
            Protected: true,
            RequiredScopes: userTypes("Residente", "Admin"),
        },
        {
            Method: "GET",
            Pattern: "/api/obtener-activo-id/:id_activo",
            Handler: ObtenerActivoPorID,
            Protected: true,
            RequiredScopes: userTypes("Tecnico", "Residente", "Analista", "Admin"),
            RequiredValidation: "activo",
        },
        { 
            Method:  "POST",
            Pattern: "/api/pdf-reporte/:id_activo",
            Handler: GenerarPDFReporte,
            Protected: true,
            RequiredScopes: userTypes("Tecnico", "Admin"),
            RequiredValidation: "activo",
        },
        {
            Method: "GET",
            Pattern: "/api/obtener-acciones/:id_tecnico",
            Handler: ObtenerAcciones,
            Protected: true,
            RequiredScopes: userTypes("Tecnico", "Admin"),
        },
        { 
            Method:  "POST",
            Pattern: "/api/subir-documento",
            Handler: SubirDocumentoHandler,
            Protected: true,
            RequiredScopes: userTypes("Tecnico", "Admin"),
        },
        { 
            Method:  "POST",
            Pattern: "/api/obtener-documentos",
            Handler: ObtenerDocumentosHandler,
            Protected: true,
            RequiredScopes: userTypes("Analista", "Tecnico", "Admin"),
        },
        {
            Method: "PUT",
            Pattern: "/api/actualizar-estado-contacto/:id_contacto",
            Handler: ActualizarEstadoContactoHandler,
            Protected: true,
            RequiredScopes: userTypes("Analista", "Tecnico", "Admin"),
        },
        {
            Method: "GET",
            Pattern: "/api/obtener-todos-activos",
            Handler: ObtenerTodosActivos,
            Protected: true,
            RequiredScopes: userTypes("Tecnico", "Analista", "Admin"),
        },
        { 
            Method:  "POST",
            Pattern: "/api/accion-mantenimiento-id/:id_activo",
            Handler: AccionMantenimientoHandler,
            Protected: true,
            RequiredScopes: userTypes("Tecnico", "Analista", "Admin"),
            RequiredValidation: "activo",
        },
        { 
            Method:  "PUT",
            Pattern: "/api/actualizar-estado-accion/:id_accion",
            Handler: ActualizarEstadoAccionHandler,
            Protected: true,
            RequiredScopes: userTypes("Tecnico", "Analista", "Admin"),
        },
        { 
            Method:  "GET",
            Pattern: "/api/obtener-firmas-usuario",
            Handler: ObtenerFirmasUsuarioHandler,
            Protected: true,
            RequiredScopes: userTypes("Tecnico", "Admin"),
        },
        { 
            Method:  "POST",
            Pattern: "/api/subir-firma-usuario",
            Handler: SubirFirmaUsuarioHandler,
            Protected: true,
            RequiredScopes: userTypes("Tecnico", "Admin"),
        },
        { 
            Method:  "DELETE",
            Pattern: "/api/eliminar-firma-usuario/:id_firma",
            Handler: EliminarFirmaUsuarioHandler,
            Protected: true,
            RequiredScopes: userTypes("Tecnico", "Admin"),
        },
        { 
            Method:  "PUT",
            Pattern: "/api/actualizar-firma-usuario/:id_firma",
            Handler: ActualizarFirmaUsuarioHandler,
            Protected: true,
            RequiredScopes: userTypes("Tecnico", "Admin"),
        },
        // A partir de acá serían los handlers para el sprint 3
        { 
            Method:  "POST",
            Pattern: "/api/crear-activo",
            Handler: CrearActivoHandler,
            Protected: true,
            RequiredScopes: rootOnly(),
        },
        { 
            Method:  "PUT",
            Pattern: "/api/editar-activo/:id_activo",
            Handler: EditarActivoHandler,
            Protected: true,
            RequiredScopes: rootOnly(),
        },
        {
            Method:  "DELETE",
            Pattern: "/api/eliminar-activo/:id_activo",
            Handler: EliminarActivoHandler,
            Protected: true,
            RequiredScopes: rootOnly(),
        },
        { 
            Method:  "POST",
            Pattern: "/api/crear-sensor",
            Handler: CrearSensorHandler,
            Protected: true,
            RequiredScopes: rootOnly(),
        },
        { 
            Method:  "PUT",
            Pattern: "/api/editar-sensor/:id_sensor",
            Handler: EditarSensorHandler,
            Protected: true,
            RequiredScopes: rootOnly(),
        },
        { 
            Method:  "GET",
            Pattern: "/api/anomalia-activo/:id_activo",
            Handler: ObtenerAnomaliaActivoHandler,
            Protected: true,
            RequiredScopes: rootOnly(),
        },
        { 
            Method:  "POST",
            Pattern: "/api/crear-edificio",
            Handler: CrearEdificioHandler,
            Protected: true,
            RequiredScopes: rootOnly(),
        },
        { 
            Method:  "DELETE",
            Pattern: "/api/eliminar-edificio/:id_edificio",
            Handler: EliminarEdificioHandler,
            Protected: true,
            RequiredScopes: rootOnly(),
        },
        { 
            Method:  "PUT",
            Pattern: "/api/actualizar-edificio/:id_edificio",
            Handler: ActualizarEdificioHandler,
            Protected: true,
            RequiredScopes: rootOnly(),
        },
        { 
            Method:  "GET",
            Pattern: "/api/codigo_qr/:codigo_activo",
            Handler: ObtenerInfoQRHandler,
            Protected: false,
            RequiredScopes: nil, // Acceso abierto
        },
        { 
            Method:  "POST",
            Pattern: "/api/crear-ticket",
            Handler: CrearTicketHandler,
            Protected: true,
            RequiredScopes: userTypes("Administrador", "Admin"),
        },
        { 
            Method:  "PUT",
            Pattern: "/api/resolver-ticket/:id_ticket",
            Handler: ResolverTicketHandler,
            Protected: true,
            RequiredScopes: userTypes("Administrador", "Admin"),
        },
        // Aquí se agregan más handlers de manera fácil
        // {
        //     Method:  "POST",
        //     Pattern: "/api/otra-ruta",
        //     Handler: OtroHandler,
        //     Protected: true,
        //     RequiredScopes: userTypes("Residente"),
        //     RequiredValidation: "edificio" o "activo"
        // }, // Importante la coma al final
        // { 
        //     Method:  "",
        //     Pattern: "",
        //     Handler: ,
        //     Protected: true,
        //     RequiredScopes: userTypes("root"),
        //     RequiredValidation: "edificio",
        // },
        // { 
        //     Method:  "",
        //     Pattern: "",
        //     Handler: ,
        //     Protected: true,
        //     RequiredScopes: rootOnly(),
        // },
    }
}

// RegisterSpecialRoutes registra todos los handlers especiales en Gin
func RegisterSpecialRoutes(r *gin.Engine) {
    // Registrar el endpoint de salud
    r.GET("/healthz", func(c *gin.Context) {
        c.JSON(200, gin.H{"ok": true})
    })

    // Registrar handlers personalizados
    specialHandlers := getSpecialHandlers()
    for _, handler := range specialHandlers {
        var routeHandlers []gin.HandlerFunc

        // Si la ruta está protegida esta pasa por el middleware de autenticación
        if handler.Protected {
            routeHandlers = append(routeHandlers, middleware.AuthMiddleware())

            // Si se requieren scopes específicos, agregar al middleware de scopes
            if len(handler.RequiredScopes) > 0 {
                routeHandlers = append(routeHandlers, middleware.ScopeMiddleware(handler.RequiredScopes))
            }

            // Si se requiere validación de acceso a edificio o activo
            if handler.RequiredValidation != "" {
                routeHandlers = append(routeHandlers, middleware.ValidationMiddleware(handler.RequiredValidation))
            }
        }
        
        // Agregar el handler principal
        routeHandlers = append(routeHandlers, handler.Handler)

        // Registrar la ruta en Gin según el método HTTP
        switch handler.Method {
        case "GET":
            r.GET(handler.Pattern, routeHandlers...)
        case "POST":
            r.POST(handler.Pattern, routeHandlers...)
        case "PUT":
            r.PUT(handler.Pattern, routeHandlers...)
        case "DELETE":
            r.DELETE(handler.Pattern, routeHandlers...)
        }

        protectionStatus := "public"
        if handler.Protected {
            protectionStatus = "protected"
            if len(handler.RequiredScopes) > 0 {
                protectionStatus += " (scopes: " + joinScopes(handler.RequiredScopes) + ")"
            }
        }

        log.Printf("Registered special handler: %s %s (%s)", handler.Method, handler.Pattern, protectionStatus)
    }

    // Registrar rutas de proxy protegidas específicas (tienen prioridad sobre las genéricas)
    protectedProxyRoutes := getProtectedProxyRoutes()
    for _, route := range protectedProxyRoutes {
        target := os.Getenv(route.TargetEnvVar)
        if target == "" {
            log.Printf("WARNING: %s not set, skipping protected proxy route %s", route.TargetEnvVar, route.Pattern)
            continue
        }

        var routeHandlers []gin.HandlerFunc

        // Siempre aplicar autenticación para rutas protegidas
        routeHandlers = append(routeHandlers, middleware.AuthMiddleware())

        // Si se requieren scopes específicos
        if len(route.RequiredScopes) > 0 {
            routeHandlers = append(routeHandlers, middleware.ScopeMiddleware(route.RequiredScopes))
        }

        // Si se requiere validación de acceso a edificio o activo
        if route.RequiredValidation != "" {
            routeHandlers = append(routeHandlers, middleware.ValidationMiddleware(route.RequiredValidation))
        }

        // Calcular el stripPrefix desde el Pattern
        stripPrefix := ""
        if strings.HasPrefix(route.Pattern, "/api/") {
            // Extraer hasta el segundo slash después de /api/
            parts := strings.SplitN(route.Pattern, "/", 4) // ["", "api", "servicio", "resto..."]
            if len(parts) >= 3 {
                stripPrefix = "/" + parts[1] + "/" + parts[2] // "/api/servicio"
            }
        }

        // Agregar el handler de proxy
        proxyHandler := proxy.CreateProxyHandler(target, stripPrefix, route.PrependPath)
        routeHandlers = append(routeHandlers, proxyHandler)

        // Registrar la ruta específica con el método indicado
        switch route.Method {
        case "GET":
            r.GET(route.Pattern, routeHandlers...)
        case "POST":
            r.POST(route.Pattern, routeHandlers...)
        case "PUT":
            r.PUT(route.Pattern, routeHandlers...)
        case "DELETE":
            r.DELETE(route.Pattern, routeHandlers...)
        default:
            r.Any(route.Pattern, routeHandlers...)
        }

        protectionStatus := "protected"
        if len(route.RequiredScopes) > 0 {
            protectionStatus += " (scopes: " + joinScopes(route.RequiredScopes) + ")"
        }

        log.Printf("Registered protected proxy: %s %s → %s (%s)", route.Method, route.Pattern, target, protectionStatus)
    }

    // Registrar rutas de proxy transparente genéricas
    proxyRoutes := getProxyRoutes()
    for _, route := range proxyRoutes {
        target := os.Getenv(route.TargetEnvVar)
        if target == "" {
            log.Printf("WARNING: %s not set, skipping proxy route %s", route.TargetEnvVar, route.PathPrefix)
            continue
        }

        var routeHandlers []gin.HandlerFunc

        // Si la ruta está protegida, aplicar middlewares
        if route.Protected {
            routeHandlers = append(routeHandlers, middleware.AuthMiddleware())

            // Si se requieren scopes específicos
            if len(route.RequiredScopes) > 0 {
                routeHandlers = append(routeHandlers, middleware.ScopeMiddleware(route.RequiredScopes))
            }

            // Si se requiere validación de acceso a edificio o activo
            if route.RequiredValidation != "" {
                routeHandlers = append(routeHandlers, middleware.ValidationMiddleware(route.RequiredValidation))
            }
        }

        // Agregar el handler de proxy
        proxyHandler := proxy.CreateProxyHandler(target, route.PathPrefix, route.PrependPath)
        routeHandlers = append(routeHandlers, proxyHandler)

        // Registrar la ruta con comodín para capturar todas las subrutas
        pattern := route.PathPrefix + "/*proxyPath"
        r.Any(pattern, routeHandlers...)

        protectionStatus := "public"
        if route.Protected {
            protectionStatus = "protected"
            if len(route.RequiredScopes) > 0 {
                protectionStatus += " (scopes: " + joinScopes(route.RequiredScopes) + ")"
            }
        }

        log.Printf("Registered proxy route: ANY %s → %s (%s)", pattern, target, protectionStatus)
    }
}

// joinScopes une los scopes en una cadena para logging
func joinScopes(scopes []string) string {
    if len(scopes) == 0 {
        return ""
    }
    result := scopes[0]
    for i := 1; i < len(scopes); i++ {
        result += ", " + scopes[i]
    }
    return result
}