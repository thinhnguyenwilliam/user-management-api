// user-management-api/internal/service/v1/auth_service.go
package v1service

import (
	"context"
	"errors"

	domain "github.com/thinhnguyenwilliam/user-management-api/internal/domain/user"
	v1dto "github.com/thinhnguyenwilliam/user-management-api/internal/models/dto/v1"
	"github.com/thinhnguyenwilliam/user-management-api/internal/repository"
	"github.com/thinhnguyenwilliam/user-management-api/internal/utils"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/auth"
)

type authService struct {
	userRepo     repository.IUserRepository
	tokenService auth.ITokenService
}

func NewAuthService(
	userRepo repository.IUserRepository,
	tokenService auth.ITokenService,
) IAuthService {
	return &authService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *authService) Login(ctx context.Context, req v1dto.LoginRequest) (*v1dto.LoginResponse, error) {
	// 1. Find user by email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// 2. Compare password
	if err := utils.CheckPassword(req.Password, user.UserPassword); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// 3. Generate JWT
	token, err := s.tokenService.GenerateAccessToken(auth.TokenPayload{
		UserID: user.UserUuid.String(),
		Email:  user.UserEmail,
		Role:   string(domain.MapRole(user.UserLevel)),
	})
	if err != nil {
		return nil, err
	}

	// 4. Return response
	return &v1dto.LoginResponse{
		AccessToken: token,
	}, nil
}
