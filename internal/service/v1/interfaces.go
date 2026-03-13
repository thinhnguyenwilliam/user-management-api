// user-management-api/internal/service/interfaces.go
package v1service

import (
	"context"

	db "github.com/thinhnguyenwilliam/user-management-api/internal/db/sqlc"
	v1dto "github.com/thinhnguyenwilliam/user-management-api/internal/models/dto/v1"
)

type IUserService interface {
	CreateUser(ctx context.Context, req v1dto.CreateUserRequest) (db.User, error)
}
