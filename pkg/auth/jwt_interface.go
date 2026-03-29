// user-management-api/pkg/auth/jwt_interface.go
package auth

type ITokenService interface {
	GenerateAccessToken(payload TokenPayload) (string, error)
	ParseAccessToken(tokenStr string) (*TokenPayload, error)
}
