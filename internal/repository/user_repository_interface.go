// user-management-api/internal/repository/user_repository_interface.go
package repository

import (
	"context"

	db "github.com/thinhnguyenwilliam/user-management-api/internal/db/sqlc"
)

type IUserRepository interface {
	Create(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetByUUID(ctx context.Context, uuid string) (db.User, error)
	Update(ctx context.Context, arg db.UpdateUserParams) (db.User, error)
}
