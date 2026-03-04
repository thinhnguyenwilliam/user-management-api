// user-management-api/internal/routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"
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
}
