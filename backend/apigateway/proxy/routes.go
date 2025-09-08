package proxy

import (
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
    // Endpoint de salud
    r.GET("/healthz", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"ok": true})
    })

    // Registra todas las rutas definidas en ProxyRoutes
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