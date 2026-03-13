// user-management-api/internal/repository/user_repository.go
package repository

import (
	"context"

	db "github.com/thinhnguyenwilliam/user-management-api/internal/db/sqlc"
)

type userRepository struct {
	q db.Querier
}

func NewUserRepository(q db.Querier) IUserRepository {
	return &userRepository{
		q: q,
	}
}

func (r *userRepository) Create(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return r.q.CreateUser(ctx, arg)
}
