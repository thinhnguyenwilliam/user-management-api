// user-management-api/internal/repository/interfaces.go
package repository

import (
	"context"

	db "github.com/thinhnguyenwilliam/user-management-api/internal/db/sqlc"
)

type IUserRepository interface {
	Create(ctx context.Context, arg db.CreateUserParams) (db.User, error)
}
