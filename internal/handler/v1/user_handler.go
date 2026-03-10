// user-management-api/internal/handler/user_handler.go
package v1handler

import (
	"github.com/gin-gonic/gin"
	v1dto "github.com/thinhnguyenwilliam/user-management-api/internal/models/dto/v1"
	"github.com/thinhnguyenwilliam/user-management-api/internal/models/mapper"
	v1service "github.com/thinhnguyenwilliam/user-management-api/internal/service/v1"
	"github.com/thinhnguyenwilliam/user-management-api/internal/utils"
)

type UserHandler struct {
	userService v1service.IUserService
}

func NewUserHandler(userService v1service.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req v1dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, utils.NewError("invalid request", utils.ErrCodeBadRequest))
		return
	}

	user, err := h.userService.CreateUser(
		c.Request.Context(),
		req,
	)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	resp := mapper.ToUserResponse(user)

	utils.ResponseSuccess(c, 200, resp)
}
