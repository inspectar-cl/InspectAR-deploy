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
    RequiredScopes      []string // Pueden ser 4 hasta el momento, user-type:Analista, user-type:Tecnico, user-type:Residente y user-type:Admin
    RequiredValidation  string   // Puede ser "edificio", "activo" o vacío si no se requiere validación.
}

// userTypes agrega el prefijo "user-type:" a cada rol especificado
func userTypes(roles ...string) []string {
    scopes := make([]string, len(roles))
    for i, role := range roles {
        scopes[i] = "user-type:" + role
    }
    return scopes
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
            Method:  "GET",
            Pattern: "/api/lista-activos-edificio/:id_edificio",
            Handler: ListaActivosHandler,
            Protected: true,
            RequiredScopes: userTypes("Tecnico", "Analista", "Admin"),
            RequiredValidation: "edificio",
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
            RequiredScopes: userTypes("Root"),
        },
        { 
            Method:  "PUT",
            Pattern: "/api/editar-activo/:id_activo",
            Handler: EditarActivoHandler,
            Protected: true,
            RequiredScopes: userTypes("Root"),
        },
        {
            Method:  "DELETE",
            Pattern: "/api/eliminar-activo/:id_activo",
            Handler: EliminarActivoHandler,
            Protected: true,
            RequiredScopes: userTypes("Root"),
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