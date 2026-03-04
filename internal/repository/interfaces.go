package repository

import (
	"context"

	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
)

type IUserRepository interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
}
