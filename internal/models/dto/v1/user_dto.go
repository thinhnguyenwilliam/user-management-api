// user-management-api/internal/models/dto/v1/user_dto.go
package v1dto

import (
	"github.com/google/uuid"
)

type UpdateUserRequest struct {
	Password *string `json:"password"`
	Fullname *string `json:"fullname"`
	Age      *int32  `json:"age"`
	Status   *int32  `json:"status"`
	Level    *int32  `json:"level"`
}

type CreateUserRequest struct {
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Age      *int32 `json:"age"`
}

type UserResponse struct {
	UUID      uuid.UUID `json:"uuid"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	Age       *int32    `json:"age"`
	Status    int32     `json:"status"`
	Level     int32     `json:"level"`
	CreatedAt string    `json:"created_at"`
}
