// user-management-api/internal/service/v1/auth_service_interface.go
package v1service

import (
	"context"

	v1dto "github.com/thinhnguyenwilliam/user-management-api/internal/models/dto/v1"
)

type IAuthService interface {
	Login(ctx context.Context, req v1dto.LoginRequest) (*v1dto.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*v1dto.LoginResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}
