// user-management-api/internal/service/interfaces.go
package service

import (
	"context"

	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
	"github.com/thinhnguyenwilliam/user-management-api/internal/models/dto"
)

type IUserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*models.User, error)
}
