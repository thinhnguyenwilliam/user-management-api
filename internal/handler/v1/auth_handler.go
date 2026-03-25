// user-management-api/internal/handler/v1/auth_handler.go
package v1handler

import (
	"github.com/gin-gonic/gin"

	v1dto "github.com/thinhnguyenwilliam/user-management-api/internal/models/dto/v1"
	v1service "github.com/thinhnguyenwilliam/user-management-api/internal/service/v1"
)

type AuthHandler struct {
	authService v1service.IAuthService
}

func NewAuthHandler(authService v1service.IAuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req v1dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 👇 convert từ gin.Context → context.Context
	ctx := c.Request.Context()

	res, err := h.authService.Login(ctx, req)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, res)
}
