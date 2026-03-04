// user-management-api/internal/service/interfaces.go
package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
)

type IUserService interface {
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, name, email, password string) (*models.User, error)
}
