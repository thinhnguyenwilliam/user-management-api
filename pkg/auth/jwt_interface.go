// user-management-api/pkg/auth/jwt_interface.go
package auth

type ITokenService interface {
	GenerateAccessToken(payload TokenPayload) (string, error)
	ParseAccessToken(tokenStr string) (*TokenPayload, error)
	GenerateRefreshToken(userID string) (string, string, error)
	ParseRefreshToken(tokenStr string) (*RefreshClaims, error)
}
