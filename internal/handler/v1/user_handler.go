// user-management-api/internal/handler/v1/user_handler.go
package v1handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

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

func (h *UserHandler) UpdateUser(c *gin.Context) {
	traceID, _ := c.Get("trace_id")
	log.Info().
		Str("trace_id", traceID.(string)).
		Msg("update user request")

	uuidStr := c.Param("uuid")

	id, err := uuid.Parse(uuidStr)
	if err != nil {
		utils.ResponseError(c, utils.NewError("invalid uuid", utils.ErrCodeBadRequest))
		return
	}

	var req v1dto.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, utils.NewError("invalid request body", utils.ErrCodeBadRequest))
		return
	}

	user, err := h.userService.UpdateUser(
		c.Request.Context(),
		id,
		req,
	)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	resp := mapper.ToUserResponse(user)

	utils.ResponseSuccess(c, 200, resp)
}

func (h *UserHandler) GetUserByUUID(c *gin.Context) {

	traceID, _ := c.Get("trace_id")

	log.Info().
		Str("trace_id", traceID.(string)).
		Msg("getting user by uuid request")

	uuid := c.Param("uuid")

	user, err := h.userService.GetUserByUUID(
		c.Request.Context(),
		uuid,
	)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	resp := mapper.ToUserResponse(user)

	utils.ResponseSuccess(c, 200, resp)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	traceID, _ := c.Get("trace_id")
	log.Info().
		Str("trace_id", traceID.(string)).
		Msg("getting user request")

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
