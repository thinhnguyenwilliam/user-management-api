// user-management-api/internal/service/v1/auth_service.go
package v1service

import (
	"context"
	"errors"
	"time"

	domain "github.com/thinhnguyenwilliam/user-management-api/internal/domain/user"
	v1dto "github.com/thinhnguyenwilliam/user-management-api/internal/models/dto/v1"
	"github.com/thinhnguyenwilliam/user-management-api/internal/repository"
	"github.com/thinhnguyenwilliam/user-management-api/internal/utils"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/auth"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/rediscache"
)

type authService struct {
	userRepo     repository.IUserRepository
	tokenService auth.ITokenService
	cache        rediscache.Cache
}

func NewAuthService(
	userRepo repository.IUserRepository,
	tokenService auth.ITokenService,
	cache rediscache.Cache,
) IAuthService {
	return &authService{
		userRepo:     userRepo,
		tokenService: tokenService,
		cache:        cache,
	}
}

func (s *authService) Logout(ctx context.Context, refreshToken string, accessToken string) error {
	// 1. parse refresh token
	refreshClaims, err := s.tokenService.ParseRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	// 2. delete refresh token (như cũ)
	refreshKey := "refresh_token:" + refreshClaims.ID
	_ = s.cache.Delete(ctx, refreshKey)

	// 3. parse access token
	accessClaims, err := s.tokenService.ParseAccessTokenRaw(accessToken)
	if err != nil {
		return err
	}

	// 4. tính TTL còn lại
	ttl := time.Until(accessClaims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil // token sắp hết hạn rồi
	}

	// 5. add vào blacklist
	blacklistKey := "blacklist_access_token:" + accessClaims.ID

	err = s.cache.Set(ctx, blacklistKey, true, ttl)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*v1dto.LoginResponse, error) {
	// 1. parse JWT
	claims, err := s.tokenService.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// 2. check Redis
	key := "refresh_token:" + claims.ID

	var userID string
	err = s.cache.Get(ctx, key, &userID)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// 🔥 chống giả mạo
	if userID != claims.UserID {
		return nil, errors.New("token mismatch")
	}

	// 3. rotate token (quan trọng)
	_ = s.cache.Delete(ctx, key)

	newRefreshToken, _, err := s.tokenService.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	// 4. new access token
	accessToken, err := s.tokenService.GenerateAccessToken(auth.TokenPayload{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return &v1dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
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

	// 🔥 4. Generate Refresh Token
	refreshToken, _, err := s.tokenService.GenerateRefreshToken(user.UserUuid.String())
	if err != nil {
		return nil, err
	}

	// 4. Return response
	return &v1dto.LoginResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}
