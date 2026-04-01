// user-management-api/pkg/auth/jwt_interface.go
package auth

import "github.com/golang-jwt/jwt/v5"

type ITokenService interface {
	GenerateAccessToken(payload TokenPayload) (string, error)
	ParseAccessToken(tokenStr string) (*TokenPayload, error)
	ParseAccessTokenRaw(tokenStr string) (*jwt.RegisteredClaims, error)
	GenerateRefreshToken(userID string) (string, string, error)
	ParseRefreshToken(tokenStr string) (*RefreshClaims, error)
}
