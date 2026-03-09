// internal/mapper/user_mapper.go
package mapper

import (
	sqlc "github.com/thinhnguyenwilliam/user-management-api/internal/db/sqlc"
	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
	"github.com/thinhnguyenwilliam/user-management-api/internal/models/dto"
)

func MapUserFromDB(src sqlc.User, dst *models.User) {
	dst.UserID = int(src.UserID)
	dst.Name = src.Name
	dst.Email = src.Email
	dst.PhoneNumber = src.PhoneNumber
	dst.CreatedAt = src.CreatedAt
}

func ToUserModelFromDB(u sqlc.User) *models.User {
	return &models.User{
		UserID:      int(u.UserID),
		Name:        u.Name,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		CreatedAt:   u.CreatedAt,
	}
}

func ToCreateUserParams(user *models.User) sqlc.CreateUserParams {
	return sqlc.CreateUserParams{
		Name:           user.Name,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
	}
}

func ToUserModel(req dto.CreateUserRequest, hashedPassword string) *models.User {
	return &models.User{
		Name:           req.Name,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}
}

func ToUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		Id:          user.UserID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func mapUser(u sqlc.User) *models.User {
	return &models.User{
		Name:        u.Name,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		CreatedAt:   u.CreatedAt,
	}
}
