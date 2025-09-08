package proxy

import (
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "apigateway/handlers"
)

func RegisterRoutes(r *gin.Engine) {
    // Endpoint de salud
    r.GET("/healthz", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"ok": true})
    })

    // Registrar endpoints compuestos usando la configuración modular
    compositeHandler := handlers.NewCompositeHandler()
    
    for _, endpoint := range handlers.CompositeEndpoints {
        switch endpoint.Method {
        case "GET":
            r.GET(endpoint.Path, compositeHandler.ExecuteCompositeEndpoint(endpoint))
        case "POST":
            r.POST(endpoint.Path, compositeHandler.ExecuteCompositeEndpoint(endpoint))
        case "PUT":
            r.PUT(endpoint.Path, compositeHandler.ExecuteCompositeEndpoint(endpoint))
        case "DELETE":
            r.DELETE(endpoint.Path, compositeHandler.ExecuteCompositeEndpoint(endpoint))
        case "PATCH":
            r.PATCH(endpoint.Path, compositeHandler.ExecuteCompositeEndpoint(endpoint))
        }
    }

    // Mantener las rutas de proxy existentes para casos específicos
    for _, route := range ProxyRoutes {
        target := os.Getenv(route.TargetEnvVar)
        if target == "" {
            continue // Si no está definida la variable, ignora la ruta
        }
        // Closure para capturar variables correctamente
        func(pathPrefix, target, prependPath string) {
            r.Any(pathPrefix+"/*proxyPath", func(c *gin.Context) {
                NewReverseProxy(target, pathPrefix, prependPath).ServeHTTP(c.Writer, c.Request)
            })
        }(route.PathPrefix, target, route.PrependPath)
    }
}