package handlers

import (
    "log"
    "github.com/gin-gonic/gin"
    "apigateway/api/middleware"
)

// SpecialHandler representa un handler especial para Gin
type SpecialHandler struct {
    Method      string
    Pattern     string
    Handler     gin.HandlerFunc
    Protected   bool
}

// getSpecialHandlers devuelve todos los handlers especiales
func getSpecialHandlers() []SpecialHandler {
    return []SpecialHandler{
        {
            Method:  "GET",
            Pattern: "/api/activos-completos",
            Handler: ActivosCompletosHandler,
            Protected: true,
        },
        {
            Method:  "GET",
            Pattern: "/api/activos-y-sensores",
            Handler: ActivosYSensores,
            Protected: true,
        },
        {
            Method: "GET",
            Pattern: "/api/obtener-activos",
            Handler: ObtenerActivos,
            Protected: true,    
        },
        {
            Method: "GET",
            Pattern: "/api/generar-reporte/:id",
            Handler: GenerarReporte,   
            Protected: true,
        },
        {
            Method: "GET",
            Pattern: "/api/obtener-activo-id/:id",
            Handler: ObtenerActivoPorID,
            Protected: true,
        },
        { 
            Method:  "POST",
            Pattern: "/api/login",
            Handler: LoginHandler,
            Protected: false,
        },
        // Aquí puedes agregar más handlers especiales fácilmente
        // {
        //     Method:  "POST",
        //     Pattern: "/api/otra-ruta",
        //     Handler: OtroHandler,
        // }, // Importante la coma al final
        // { 
        //     Method:  "",
        //     Pattern: "",
        //     Handler: ,
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
        }

        log.Printf("Registered special handler: %s %s (%s)", handler.Method, handler.Pattern, protectionStatus)
    }
}