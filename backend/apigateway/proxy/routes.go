package proxy

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Registra las rutas /api/* y las conecta con los proxies
func RegisterRoutes(r *gin.Engine, GestionURL string) {
	// Ej: /healthz
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// Redirige cualquier ruta /api/gestion/* al microservicio de gestión
	r.Any("/api/gestion/*proxyPath", func(c *gin.Context) {
		NewReverseProxy(GestionURL, "/api/gestion", "").ServeHTTP(c.Writer, c.Request)
	})
}
