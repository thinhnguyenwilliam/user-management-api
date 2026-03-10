// user-management-api/internal/service/v1/user_service.go
package v1service

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
	v1dto "github.com/thinhnguyenwilliam/user-management-api/internal/models/dto/v1"
	"github.com/thinhnguyenwilliam/user-management-api/internal/models/mapper"
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

func isDuplicateKey(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr.Code == "23505"
	}
	return false
}

func (s *userService) CreateUser(
	ctx context.Context,
	req v1dto.CreateUserRequest,
) (*models.User, error) {

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return nil, utils.NewError("missing required fields", utils.ErrCodeInvalidInput)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := mapper.ToUserModel(req, hashedPassword)

	err = s.userRepo.Create(ctx, user)
	if err != nil {

		if isDuplicateKey(err) {
			return nil, utils.NewError("email already exists", utils.ErrCodeConflict)
		}

		return nil, utils.WrapError("failed to create user", utils.ErrCodeDatabase, err)
	}

	return user, nil
}

// type userService struct {
// 	userRepo repository.IUserRepository
// }

// func NewUserService(userRepo repository.IUserRepository) IUserService {
// 	return &userService{
// 		userRepo: userRepo,
// 	}
// }

// func (s *userService) CreateUser(
// 	ctx context.Context,
// 	req dto.CreateUserRequest,
// ) (*models.User, error) {

// 	if req.Name == "" || req.Email == "" || req.Password == "" {
// 		return nil, utils.NewError("missing required fields", utils.ErrCodeInvalidInput)
// 	}

// 	email := utils.NormalizeEmail(req.Email)

// 	// check email exists
// 	existingUser, _ := s.userRepo.FindByEmail(ctx, email)
// 	if existingUser != nil {
// 		return nil, utils.NewError("email already exists", utils.ErrCodeConflict)
// 	}

// 	// hash password
// 	hashedPassword, err := utils.HashPassword(req.Password)
// 	if err != nil {
// 		return nil, utils.WrapError("failed to hash password", utils.ErrCodeInternal, err)
// 	}

// 	user := mapper.ToUserModel(req, hashedPassword)

// 	if err := s.userRepo.Create(ctx, user); err != nil {
// 		return nil, utils.WrapError("failed to create user", utils.ErrCodeDatabase, err)
// 	}

// 	return user, nil
// }
