// user-management-api/internal/repository/user_repository.go
package repository

import (
	"context"

	"github.com/rs/zerolog/log"
	db "github.com/thinhnguyenwilliam/user-management-api/internal/db/sqlc"
	"github.com/thinhnguyenwilliam/user-management-api/internal/middleware"
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
	traceID := ctx.Value(middleware.TraceIDKey).(string)

	log.Info().
		Str("trace_id", traceID).
		Msg("getting user request")
	return r.q.CreateUser(ctx, arg)
}
