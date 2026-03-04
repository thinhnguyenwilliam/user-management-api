// user-management-api/internal/service/user_service.go
package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
	"github.com/thinhnguyenwilliam/user-management-api/internal/repository"
)

type userService struct {
	userRepo repository.IUserRepository
}

func NewUserService(userRepo repository.IUserRepository) IUserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(
	ctx context.Context,
	name, email, password string,
) (*models.User, error) {

	if name == "" || email == "" || password == "" {
		return nil, errors.New("missing required fields")
	}

	user := &models.User{
		UUID:  uuid.NewString(),
		Name:  name,
		Email: email,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}
