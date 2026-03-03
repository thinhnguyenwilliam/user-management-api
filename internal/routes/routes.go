// user-management-api/internal/routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thinhnguyenwilliam/user-management-api/internal/middleware"
)

type Route interface {
	Register(rg *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {

	api := r.Group("/api/v1")
	api.Use(middleware.Logger()) // chỉ apply cho api v1

	for _, route := range routes {
		route.Register(api)
	}
}
