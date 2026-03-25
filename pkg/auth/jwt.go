// user-management-api/pkg/auth/jwt.go
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenTTL = 15 * time.Minute
)

type JWTService struct {
	secretKey string
}

func NewJWTService(secret string) ITokenService {
	return &JWTService{
		secretKey: secret,
	}
}

type TokenPayload struct {
	UserID string
	Role   string
	Email  string
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (js *JWTService) GenerateAccessToken(payload TokenPayload) (string, error) {
	claims := Claims{
		UserID: payload.UserID,
		Email:  payload.Email,
		Role:   payload.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   payload.UserID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user-management-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(js.secretKey))
}
