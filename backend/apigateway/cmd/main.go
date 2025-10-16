package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"apigateway/config"
	"apigateway/proxy"
	"apigateway/handlers"
	"apigateway/api/middleware"
)

func main() {
	cfg := config.Load()

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORS())
	
	// Rutas de proxy
	proxy.RegisterRoutes(r)

	// Registrar handlers
	handlers.RegisterSpecialRoutes(r)

	log.Printf("API Gateway escuchando en :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
