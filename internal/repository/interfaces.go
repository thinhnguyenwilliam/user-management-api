package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
)

type IUserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
}
