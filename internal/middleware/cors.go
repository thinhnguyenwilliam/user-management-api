// user-management-api/internal/middleware/cors.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")

		allowedOrigins := map[string]bool{
			"http://127.0.0.1:5500": true,
			"http://localhost:3000": true,
		}

		if allowedOrigins[origin] {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		ctx.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, X-API-Key, Authorization",
		)
		ctx.Writer.Header().Set("Access-Control-Allow-Methods",
			"GET, POST, PUT, PATCH, DELETE, OPTIONS",
		)
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Expose-Headers",
			"Content-Length, Authorization",
		)
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
