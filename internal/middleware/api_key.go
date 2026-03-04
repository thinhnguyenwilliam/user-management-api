// user-management-api/internal/middleware/api_key.go
package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleware(expectedKey string) gin.HandlerFunc {
	if expectedKey == "" {
		log.Fatal("API_KEY is not configured")
	}

	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-Key")
		if apiKey == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Missing X-API-Key",
			})
			return
		}

		if apiKey != expectedKey {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API Key",
			})
			return
		}

		ctx.Set("username", "willam")
		ctx.Next()
	}
}
