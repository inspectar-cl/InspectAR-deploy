package router

import (
	"oauth2/api/middleware"
	"oauth2/internal/handler"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(loginHandler *handler.LoginHandler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//Rutas publicas
	r.POST("/login", loginHandler.Login)
	r.POST("/refresh", loginHandler.RefreshToken)

	// Rutas protegidas con middleware
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(loginHandler.AuthService))
	protected.GET("/profile", func(c *gin.Context) {
		userID, _ := c.Get("urername")
		exp, _ := c.Get("exp")
		c.JSON(200, gin.H{"message": "Bienvenido", "user_id": userID, "exp": exp})
	})

	protected.POST("/logout", loginHandler.Logout)
	protected.POST("/change-tokens", loginHandler.ChangePasswordToken)

	// Rutas protegidas con middleware de scope

	/*
		protected.GET("/users", middleware.ScopeMiddleware("read:users"), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Lista de usuarios"})
		})

		protected.POST("/users", middleware.ScopeMiddleware("write:users"), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Usuario creado"})
		})
	*/

	return r
}
