// user-management-api/internal/service/interfaces.go
package service

import (
	"context"

	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
)

type IUserService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error)
}
