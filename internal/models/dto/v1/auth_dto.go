// user-management-api/internal/models/dto/v1/auth_dto.go
package v1dto

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	AccessToken string
}
