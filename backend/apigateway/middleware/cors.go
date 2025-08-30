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
        AllowAllOrigins:  true, // Permite cualquier origen
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        AllowCredentials: false, // Debe ser false si AllowAllOrigins es true
        MaxAge:           12 * time.Hour,
    })
}
