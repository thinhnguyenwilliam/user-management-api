// user-management-api/internal/models/user.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID           uuid.UUID `json:"uuid"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Age            int       `json:"age"`
	Status         int       `json:"status"`
	Level          int       `json:"level"`
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Password string `json:"password" binding:"required,min=6"`
	Level    int    `json:"level" binding:"required,oneof=1 2"`
}
