// user-management-api/internal/models/user.go
package models

import "time"

type User struct {
	UserID         int
	Name           string
	Email          string
	HashedPassword string
	PhoneNumber    string
	CreatedAt      time.Time
}

// type User struct {
// 	UserID         int       `gorm:"column:user_id;primaryKey"`
// 	Name           string    `gorm:"column:name"`
// 	Email          string    `gorm:"column:email"`
// 	HashedPassword string    `gorm:"column:hashed_password"`
// 	PhoneNumber    string    `gorm:"column:phone_number"`
// 	CreatedAt      time.Time `gorm:"column:created_at"`
// }
