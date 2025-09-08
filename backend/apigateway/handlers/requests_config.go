package handlers

import (
    "log"
    "github.com/gin-gonic/gin"
)

// SpecialHandler representa un handler especial para Gin
type SpecialHandler struct {
    Method  string
    Pattern string
    Handler gin.HandlerFunc
}

// getSpecialHandlers devuelve todos los handlers especiales
func getSpecialHandlers() []SpecialHandler {
    return []SpecialHandler{
        {
            Method:  "GET",
            Pattern: "/api/activos-completos",
            Handler: ActivosCompletosHandler,
        },
        {
            Method:  "GET",
            Pattern: "/api/sensores",
            Handler: ActivosYSensores,
        },
        // Aquí puedes agregar más handlers especiales fácilmente
        // {
        //     Method:  "POST",
        //     Pattern: "/api/otra-ruta",
        //     Handler: OtroHandler,
        // },
    }
}

// RegisterSpecialRoutes registra todos los handlers especiales en Gin
func RegisterSpecialRoutes(r *gin.Engine) {
    specialHandlers := getSpecialHandlers()
    for _, handler := range specialHandlers {
        switch handler.Method {
        case "GET":
            r.GET(handler.Pattern, handler.Handler)
        case "POST":
            r.POST(handler.Pattern, handler.Handler)
        case "PUT":
            r.PUT(handler.Pattern, handler.Handler)
        case "DELETE":
            r.DELETE(handler.Pattern, handler.Handler)
        }
        log.Printf("Registered special handler: %s %s", handler.Method, handler.Pattern)
    }
}