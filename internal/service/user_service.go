// user-management-api/internal/service/user_service.go
package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
	"github.com/thinhnguyenwilliam/user-management-api/internal/repository"
	"github.com/thinhnguyenwilliam/user-management-api/internal/utils"
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
	req models.CreateUserRequest,
) (*models.User, error) {

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return nil, utils.NewError("missing required fields", utils.ErrCodeInvalidInput)
	}

	email := utils.NormalizeEmail(req.Email)

	// check email exists
	existingUser, _ := s.userRepo.FindByEmail(ctx, email)
	if existingUser != nil {
		return nil, utils.NewError("email already exists", utils.ErrCodeConflict)
	}

	// hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, utils.WrapError("failed to hash password", utils.ErrCodeInternal, err)
	}

	user := &models.User{
		UUID:           uuid.New(),
		Name:           req.Name,
		Email:          email,
		HashedPassword: hashedPassword,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, utils.WrapError("failed to create user", utils.ErrCodeDatabase, err)
	}

	return user, nil
}
