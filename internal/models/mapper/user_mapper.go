// user-management-api/internal/models/mapper/user_mapper.go
package mapper

import (
	"github.com/google/uuid"
	db "github.com/thinhnguyenwilliam/user-management-api/internal/db/sqlc"
	v1dto "github.com/thinhnguyenwilliam/user-management-api/internal/models/dto/v1"
)

func ToUpdateUserParams(
	req v1dto.UpdateUserRequest,
	id uuid.UUID,
) db.UpdateUserParams {

	return db.UpdateUserParams{
		UserUuid:     id,
		UserPassword: req.Password,
		UserFullname: req.Fullname,
		UserAge:      req.Age,
		UserStatus:   req.Status,
		UserLevel:    req.Level,
	}
}

func ToCreateUserParams(req v1dto.CreateUserRequest, hashedPassword string) db.CreateUserParams {
	return db.CreateUserParams{
		UserFullname: req.Fullname,
		UserEmail:    req.Email,
		UserPassword: hashedPassword,
		UserAge:      req.Age,
		UserStatus:   1,
		UserLevel:    3,
	}
}

func ToUserResponse(u db.User) v1dto.UserResponse {
	return v1dto.UserResponse{
		UUID:      u.UserUuid,
		Fullname:  u.UserFullname,
		Email:     u.UserEmail,
		Age:       u.UserAge,
		Status:    u.UserStatus,
		Level:     u.UserLevel,
		CreatedAt: u.UserCreatedAt.Format("2006-01-02 15:04:05"),
	}
}
