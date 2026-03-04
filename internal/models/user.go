// user-management-api/internal/models/user.go
package models

import "time"

type User struct {
	UUID           string    `json:"uuid"`
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
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Age      int    `json:"age"`
	Password string `json:"password" binding:"required,min=6"`
	Level    int    `json:"level" binding:"required,oneof=1 2"`
}
