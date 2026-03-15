// user-management-api/internal/middleware/trace_id.go
package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const TraceIDKey contextKey = "trace_id"

func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {

		traceID := c.GetHeader("X-Trace-ID")

		if traceID == "" {
			traceID = uuid.NewString()
		}

		// add trace id vào request context
		ctx := context.WithValue(c.Request.Context(), TraceIDKey, traceID)
		c.Request = c.Request.WithContext(ctx)

		// set vào gin context
		c.Set("trace_id", traceID)

		// trả về client
		c.Writer.Header().Set("X-Trace-ID", traceID)

		c.Next()
	}
}
