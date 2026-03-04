// user-management-api/internal/routes/user_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thinhnguyenwilliam/user-management-api/internal/handler"
)

type UserRoutes struct {
	userHandler *handler.UserHandler
}

func NewUserRoutes(userHandler *handler.UserHandler) *UserRoutes {
	return &UserRoutes{
		userHandler: userHandler,
	}
}

func (r *UserRoutes) Register(rg *gin.RouterGroup) {
	userGroup := rg.Group("/users")
	{
		// userGroup.GET("/:id", r.userHandler.GetUser)
		userGroup.POST("/", r.userHandler.CreateUser)
	}
}
