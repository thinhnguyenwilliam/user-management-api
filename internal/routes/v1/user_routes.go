// user-management-api/internal/routes/user_routes.go
package v1routes

import (
	"github.com/gin-gonic/gin"
	v1handler "github.com/thinhnguyenwilliam/user-management-api/internal/handler/v1"
)

type UserRoutes struct {
	userHandler *v1handler.UserHandler
}

func NewUserRoutes(userHandler *v1handler.UserHandler) *UserRoutes {
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
