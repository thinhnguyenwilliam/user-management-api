// user-management-api/internal/service/v1/user_service_interface.go
package v1service

import (
	"context"

	"github.com/google/uuid"
	db "github.com/thinhnguyenwilliam/user-management-api/internal/db/sqlc"
	v1dto "github.com/thinhnguyenwilliam/user-management-api/internal/models/dto/v1"
)

type IUserService interface {
	CreateUser(ctx context.Context, req v1dto.CreateUserRequest) (db.User, error)
	GetUserByUUID(ctx context.Context, uuid string) (db.User, error)
	UpdateUser(
		ctx context.Context,
		id uuid.UUID,
		req v1dto.UpdateUserRequest,
	) (db.User, error)
}
