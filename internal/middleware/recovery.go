// user-management-api/internal/middleware/recovery.go
package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {

		defer func() {
			if err := recover(); err != nil {

				log.Printf("PANIC is: %v\nDebug stack is: %s", err, debug.Stack())

				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}
