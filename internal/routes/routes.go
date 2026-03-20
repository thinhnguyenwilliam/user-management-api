// user-management-api/internal/routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Route interface {
	Register(rg *gin.RouterGroup)
}

func RegisterRoutes(
	r *gin.Engine,
	middlewares []gin.HandlerFunc,
	routes ...Route,
) {
	api := r.Group("/api/v1")

	// Apply all middlewares to api/v1
	if len(middlewares) > 0 {
		api.Use(middlewares...)
	}

	for _, route := range routes {
		route.Register(api)
	}

	r.NoRoute(func(c *gin.Context) {
		log.Warn().
			Str("path", c.Request.URL.Path).
			Str("method", c.Request.Method).
			Msg("route not found")

		c.JSON(404, gin.H{
			"code":    "ROUTE_NOT_FOUND",
			"message": "API endpoint not found",
		})
	})

	r.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{
			"code":    "METHOD_NOT_ALLOWED",
			"message": "Method not allowed",
		})
	})
}
