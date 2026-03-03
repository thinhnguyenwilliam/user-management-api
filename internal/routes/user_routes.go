// user-management-api/internal/routes/user_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thinhnguyenwilliam/user-management-api/internal/handler"
)

type UserRoutes struct {
	router      *gin.Engine
	userHandler *handler.UserHandler
}

func NewUserRoutes(router *gin.Engine, userHandler *handler.UserHandler) *UserRoutes {
	return &UserRoutes{
		router:      router,
		userHandler: userHandler,
	}
}

func (r *UserRoutes) Register() {
	userGroup := r.router.Group("/users")
	{
		userGroup.GET("/:id", r.userHandler.GetUser)
	}
}
