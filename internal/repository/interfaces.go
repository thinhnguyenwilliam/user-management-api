// user-management-api/internal/repository/interfaces.go
package repository

import (
	"context"

	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
)

type IUserRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	Delete(ctx context.Context, id int) error
}
