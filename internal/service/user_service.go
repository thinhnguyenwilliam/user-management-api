// user-management-api/internal/service/user_service.go
package service

import (
	"context"
	"errors"
	"time"

	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
	"github.com/thinhnguyenwilliam/user-management-api/internal/repository"
)

type UserService interface {
	GetUser(ctx context.Context, id int64) (*models.User, error)
	CreateUser(ctx context.Context, username, email, password string) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUser(ctx context.Context, id int64) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) CreateUser(
	ctx context.Context,
	username, email, password string,
) (*models.User, error) {

	if username == "" || email == "" || password == "" {
		return nil, errors.New("missing required fields")
	}

	user := &models.User{
		ID:       time.Now().UnixNano(), // demo ID
		Username: username,
		Email:    email,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
