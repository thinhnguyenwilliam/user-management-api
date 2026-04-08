// user-management-api/internal/handler/v1/user_handler.go
package v1handler

import (
	"net/http"
	"strconv"

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

func GetString(c *gin.Context, key string) string {
	v, ok := c.Get(key)
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	traceID := GetString(c, "trace_id")

	log.Info().
		Str("trace_id", traceID).
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

	utils.ResponseSuccess(c, 200, "user updated", resp)
}

func (h *UserHandler) ListUsers(c *gin.Context) {

	// query params
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	search := c.Query("search")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	var searchPtr *string
	if search != "" {
		searchPtr = &search
	}

	users, total, err := h.userService.ListUsers(
		c,
		int32(limit),
		int32(offset),
		searchPtr,
	)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	resp := mapper.ToUserResponseList(users)

	// ✅ tính page từ offset
	page := int32(1)
	if limit > 0 {
		page = int32(offset/int(limit) + 1)
	}

	// ✅ dùng PaginationResponse
	result := utils.NewPaginationResponse(
		resp,
		page,
		int32(limit),
		int32(total),
	)

	c.JSON(http.StatusOK, result)
}

func (h *UserHandler) RestoreUser(c *gin.Context) {

	traceID, _ := c.Get("trace_id")

	log.Info().
		Str("trace_id", traceID.(string)).
		Msg("restore user request")

	uuidStr := c.Param("uuid")

	id, err := uuid.Parse(uuidStr)
	if err != nil {
		utils.ResponseError(c, utils.NewError("invalid uuid", utils.ErrCodeBadRequest))
		return
	}

	user, err := h.userService.RestoreUser(
		c.Request.Context(),
		id,
	)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	resp := mapper.ToUserResponse(user)

	utils.ResponseSuccess(c, 200, "user restore", resp)
}

func (h *UserHandler) DeleteUserSoft(c *gin.Context) {

	traceID, _ := c.Get("trace_id")

	log.Info().
		Str("trace_id", traceID.(string)).
		Msg("soft delete user request")

	uuidStr := c.Param("uuid")

	id, err := uuid.Parse(uuidStr)
	if err != nil {
		utils.ResponseError(c, utils.NewError("invalid uuid", utils.ErrCodeBadRequest))
		return
	}

	user, err := h.userService.DeleteUserSoft(
		c.Request.Context(),
		id,
	)

	if err != nil {
		utils.ResponseError(c, err)
		return
	}

	resp := mapper.ToUserResponse(user)

	utils.ResponseSuccess(c, 200, "user deleted", resp)
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

	utils.ResponseSuccess(c, 200, "user updated", resp)
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

	utils.ResponseSuccess(c, 200, "user updated", resp)
}
