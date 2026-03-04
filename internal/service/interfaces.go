// user-management-api/internal/service/interfaces.go
package service

import (
	"context"

	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
)

type IUserService interface {
	GetUser(ctx context.Context, id int64) (*models.User, error)
	CreateUser(ctx context.Context, username, email, password string) (*models.User, error)
}
