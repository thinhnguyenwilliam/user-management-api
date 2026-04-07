// user-management-api/internal/repository/user_repository_interface.go
package repository

import (
	"context"

	"github.com/google/uuid"
	db "github.com/thinhnguyenwilliam/user-management-api/internal/db/sqlc"
)

type IUserRepository interface {
	UpdatePassword(ctx context.Context, userID string, password string) error
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
	ListUsersV2(ctx context.Context, in ListUsersInput) ([]db.User, error)
	CountUsers(ctx context.Context, params db.CountUsersParams) (int64, error)
	ListUsers(ctx context.Context, arg db.ListUsersOrderByCreatedAtDescParams) ([]db.User, error)
	Create(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetByUUID(ctx context.Context, uuid string) (db.User, error)
	Update(ctx context.Context, arg db.UpdateUserParams) (db.User, error)
	DeleteSoft(ctx context.Context, userUUID uuid.UUID) (db.User, error)
	Restore(ctx context.Context, userUUID uuid.UUID) (db.User, error)
}
