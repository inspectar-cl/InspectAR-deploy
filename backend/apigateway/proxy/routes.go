package proxy

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Registra las rutas /api/* y las conecta con los proxies
func RegisterRoutes(r *gin.Engine, sensoresURL, activoURL string) {
	// Ej: /healthz
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// /api/sensores/* → http://SENSORES_URL/*
	r.Any("/api/sensores/*proxyPath", func(c *gin.Context) {
		NewReverseProxy(sensoresURL, "/api/sensores", "").ServeHTTP(c.Writer, c.Request)
	})

	// /api/activo/* → http://ACTIVO_URL/activo/*
	// (tu MS espera prefijo /activo en backend)
	r.Any("/api/activo/*proxyPath", func(c *gin.Context) {
		NewReverseProxy(activoURL, "/api/activo", "/activo").ServeHTTP(c.Writer, c.Request)
	})
}
