// user-management-api/internal/routes/v1/user_routes.go
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

func (r *UserRoutes) RegisterPublic(rg *gin.RouterGroup) {
	userGroup := rg.Group("/users")
	{
		userGroup.POST("", r.userHandler.CreateUser) // 👈 register
	}
}

func (r *UserRoutes) RegisterProtected(rg *gin.RouterGroup) {
	userGroup := rg.Group("/users")
	{
		userGroup.GET("", r.userHandler.ListUsers)
		userGroup.GET("/:uuid", r.userHandler.GetUserByUUID)
		userGroup.PUT("/:uuid", r.userHandler.UpdateUser)
		userGroup.DELETE("/:uuid", r.userHandler.DeleteUserSoft)
		userGroup.PATCH("/:uuid/restore", r.userHandler.RestoreUser)
	}
}
