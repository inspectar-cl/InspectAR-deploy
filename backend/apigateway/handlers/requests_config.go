package handlers

import (
    "log"
    "github.com/gin-gonic/gin"
    "apigateway/api/middleware"
)

// SpecialHandler representa un handler especial para Gin
type SpecialHandler struct {
    Method              string
    Pattern             string
    Handler             gin.HandlerFunc
    Protected           bool
    RequiredScopes      []string // Pueden ser 3 hasta el momento, user-type:Analista, user-type:Tecnico, user-type:Residente
    // RequiredValidation  bool   // Si es true, se valida que el usuario tenga acceso al edificio/activo especificado
}

// getSpecialHandlers devuelve todos los handlers especiales
func getSpecialHandlers() []SpecialHandler {
    return []SpecialHandler{
        {
            Method:  "GET",
            Pattern: "/api/activos-completos",
            Handler: ActivosCompletosHandler,
            Protected: true,
            RequiredScopes: []string{"user-type:Tecnico", "user-type:Residente"},
        },
        {
            Method:  "GET",
            Pattern: "/api/activos-y-sensores",
            Handler: ActivosYSensores,
            Protected: true,
            RequiredScopes: []string{"user-type:Tecnico", "user-type:Residente"},
        },
        {
            Method: "GET",
            Pattern: "/api/obtener-activos",
            Handler: ObtenerActivos,
            Protected: true,
            RequiredScopes: []string{"user-type:Tecnico", "user-type:Residente", "user-type:Analista", "user-type:Admin"},
        },
        {
            Method: "GET",
            Pattern: "/api/generar-reporte/:id_reporte",
            Handler: GenerarReporte,   
            Protected: true,
            RequiredScopes: []string{"user-type:Tecnico", "user-type:Residente"},
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
            RequiredScopes: []string{"user-type:Tecnico", "user-type:Residente", "user-type:Analista", "user-type:Admin"},
        },
        // A partir de acá serían los handlers para el sprint 2
        { 
            Method:  "GET",
            Pattern: "/api/obtener-activos-id/:id_edificio",
            Handler: ObtenerActivosPorEdificio,
            Protected: true,
            RequiredScopes: []string{"user-type:Residente", "user-type:Admin"},
        },
        {
            Method:  "GET",
            Pattern: "/api/obtener-contactos-id/:id_edificio",
            Handler: ObtenerContactosPorEdificio,
            Protected: true,
            RequiredScopes: []string{"user-type:Residente", "user-type:Admin"},
        },
        { 
            Method:  "GET",
            Pattern: "/api/foro-edificio/:id_edificio",
            Handler: ForoEdificioHandler,
            Protected: true,
            RequiredScopes: []string{"user-type:Residente", "user-type:Admin"},
        },
        { 
            Method:  "POST",
            Pattern: "/api/publicacion-foro-id/:id_edificio",
            Handler: PublicacionForoHandler,
            Protected: true,
            RequiredScopes: []string{"user-type:Residente", "user-type:Admin"},
        },
        { 
            Method:  "POST",
            Pattern: "/api/comentarios-foro-id/:id_publicacion",
            Handler: ComentariosForoHandler,
            Protected: true,
            RequiredScopes: []string{"user-type:Residente", "user-type:Admin"},
        },
        { 
            Method:  "GET",
            Pattern: "/api/lista-activos-edificio/:id_edificio",
            Handler: ListaActivosHandler,
            Protected: true,
            RequiredScopes: []string{"user-type:Tecnico", "user-type:Analista", "user-type:Admin"},
        },
        {
            Method: "GET",
            Pattern: "/api/obtener-activo-id/:id_activo",
            Handler: ObtenerActivoPorID,
            Protected: true,
            RequiredScopes: []string{"user-type:Tecnico", "user-type:Residente"},
        },
        { 
            Method:  "POST",
            Pattern: "/api/pdf-reporte/:id_activo",
            Handler: GenerarPDFReporte,
            Protected: true,
            RequiredScopes: []string{"user-type:Tecnico", "user-type:Admin"},
        },
        {
            Method: "GET",
            Pattern: "/api/obtener-acciones",
            Handler: ObtenerAcciones,
            Protected: true,
            RequiredScopes: []string{"user-type:Tecnico", "user-type:Admin"},
        },
        { 
            Method:  "POST",
            Pattern: "/api/subir-documento",
            Handler: SubirDocumentoHandler,
            Protected: true,
            RequiredScopes: []string{"user-type:Tecnico", "user-type:Admin"},
        },
        // Aquí se agregan más handlers de manera fácil
        // {
        //     Method:  "POST",
        //     Pattern: "/api/otra-ruta",
        //     Handler: OtroHandler,
        //     Protected: true,
        //     RequiredScopes: []string{"user-type:Residente"},
        // }, // Importante la coma al final
        // { 
        //     Method:  "",
        //     Pattern: "",
        //     Handler: ,
        //     Protected: true,
        //     RequiredScopes: []string{"user-type:Analista", "user-type:Admin"},
        // },
    }
}

// RegisterSpecialRoutes registra todos los handlers especiales en Gin
func RegisterSpecialRoutes(r *gin.Engine) {
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