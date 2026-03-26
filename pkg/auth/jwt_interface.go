// user-management-api/pkg/auth/jwt_interface.go
package auth

type ITokenService interface {
	GenerateAccessToken(payload TokenPayload) (string, error)
	ParseAccessToken(tokenStr string) (*TokenPayload, error)
	// GenerateRefreshToken(userID string) (string, error)
	// ValidateToken(token string) (*Claims, error)
	// GenerateRefreshToken(user sqlc.User) (RefreshToken, error)
	// DecryptAccessTokenPayload(tokenString string) (*EncryptedPayload, error)
	// StoreRefreshToken(token RefreshToken) error
	// ValidateRefreshToken(token string) (RefreshToken, error)
	// RevokeRefreshToken(token string) error
}
