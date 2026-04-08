// user-management-api/internal/routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/thinhnguyenwilliam/user-management-api/internal/middleware"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/auth"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/rediscache"
)

type Route interface {
	RegisterPublic(rg *gin.RouterGroup)
	RegisterProtected(rg *gin.RouterGroup)
}

func RegisterRoutes(
	r *gin.Engine,
	middlewares []gin.HandlerFunc,
	tokenService auth.ITokenService,
	publicRoutes []Route,
	protectedRoutes []Route,
	cache rediscache.Cache,
) {
	api := r.Group("/api/v1")

	// public group
	public := api.Group("")
	for _, route := range publicRoutes {
		route.RegisterPublic(public)
	}

	// protected group
	protected := api.Group("")
	if len(middlewares) > 0 {
		protected.Use(middlewares...)
	}
	protected.Use(middleware.AuthMiddleware(tokenService, cache))
	for _, route := range protectedRoutes {
		route.RegisterProtected(protected)
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
