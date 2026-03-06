// internal/mapper/user_mapper.go
package mapper

import (
	"github.com/thinhnguyenwilliam/user-management-api/internal/models"
	"github.com/thinhnguyenwilliam/user-management-api/internal/models/dto"
)

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
	}
}
