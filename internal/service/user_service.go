// user-management-api/internal/service/user_service.go
package service

import (
	"context"
	"errors"
	"time"

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
