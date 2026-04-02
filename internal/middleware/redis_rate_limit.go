package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/ratelimiter"
)

func RedisRateLimitMiddleware(rl ratelimiter.IRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 👉 chọn key: IP hoặc userID
		key := c.ClientIP()

		// nếu có user (JWT)
		if user, exists := c.Get("user"); exists {
			key = user.(map[string]interface{})["user_id"].(string)
		}

		allowed, err := rl.Allow(c.Request.Context(), key)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "rate limit error"})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}

		c.Next()
	}
}
