// user-management-api/internal/middleware/auth.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/auth"
)

func AuthMiddleware(tokenService auth.ITokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing token",
			})
			return
		}

		// chuẩn format: Bearer <token>
		tokenStr := extractToken(authHeader)
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token format",
			})
			return
		}

		// 🔥 parse token
		payload, err := tokenService.ParseAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		// 🔥 inject user vào context
		c.Set("user", payload)

		c.Next()
	}
}

func extractToken(authHeader string) string {
	const prefix = "Bearer "
	if len(authHeader) > len(prefix) && authHeader[:len(prefix)] == prefix {
		return authHeader[len(prefix):]
	}
	return ""
}
