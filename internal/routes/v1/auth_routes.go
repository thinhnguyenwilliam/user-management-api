// user-management-api/internal/routes/v1/auth_routes.go
package v1routes

import (
	"github.com/gin-gonic/gin"
	v1handler "github.com/thinhnguyenwilliam/user-management-api/internal/handler/v1"
)

type AuthRoutes struct {
	authHandler *v1handler.AuthHandler
}

func NewAuthRoutes(authHandler *v1handler.AuthHandler) *AuthRoutes {
	return &AuthRoutes{
		authHandler: authHandler,
	}
}

func (r *AuthRoutes) RegisterProtected(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/logout", r.authHandler.Logout)
	}
}

func (r *AuthRoutes) RegisterPublic(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/login", r.authHandler.Login)
		authGroup.POST("/refresh-token", r.authHandler.RefreshToken)
		authGroup.POST("/forgot-password", r.authHandler.ForgotPassword)
		authGroup.POST("/reset-password", r.authHandler.ResetPassword)
	}
}
