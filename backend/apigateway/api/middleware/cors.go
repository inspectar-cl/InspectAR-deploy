package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// func CORS(origin string) gin.HandlerFunc {
//     return cors.New(cors.Config{
//         AllowOrigins:     []string{origin},
//         AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
//         AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
//         AllowCredentials: true,
//         MaxAge:           12 * time.Hour,
//     })
// }

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization",
			"ngrok-skip-browser-warning",
		},
		// Si quieres permitir cualquier header:
		// AllowHeaders: []string{"*"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
}
