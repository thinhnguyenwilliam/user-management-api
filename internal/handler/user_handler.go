// user-management-api/internal/handler/user_handler.go
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
	"github.com/thinhnguyenwilliam/user-management-api/internal/service"
	"github.com/thinhnguyenwilliam/user-management-api/internal/validation"
)

type UserHandler struct {
	userService service.IUserService
}

func NewUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, validation.HandleValidationErrors(err))
		return
	}

	user, err := h.userService.CreateUser(
		c.Request.Context(),
		req.Name,
		req.Email,
		req.Password,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// func (h *UserHandler) GetUser(c *gin.Context) {
// 	idParam := c.Param("id")

// 	id, err := strconv.ParseInt(idParam, 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
// 		return
// 	}

// 	user, err := h.userService.GetUser(c.Request.Context(), id)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, user)
// }
