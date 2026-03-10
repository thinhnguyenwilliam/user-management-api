// user-management-api/internal/service/interfaces.go
package v1service

import (
	"context"

	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
	v1dto "github.com/thinhnguyenwilliam/user-management-api/internal/models/dto/v1"
)

type IUserService interface {
	CreateUser(ctx context.Context, req v1dto.CreateUserRequest) (*models.User, error)
}

// type IUserService interface {
// 	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*models.User, error)
// }
