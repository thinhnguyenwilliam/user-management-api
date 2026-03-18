// user-management-api/internal/repository/user_repository.go
package repository

import (
	"context"

	"github.com/google/uuid"
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

func (r *userRepository) CountUsers(ctx context.Context, params db.CountUsersParams) (int64, error) {
	return r.q.CountUsers(ctx, params)
}

func (r *userRepository) DeleteSoft(
	ctx context.Context,
	userUUID uuid.UUID,
) (db.User, error) {

	traceID := ctx.Value(middleware.TraceIDKey).(string)

	log.Info().
		Str("trace_id", traceID).
		Msg("soft deleting user")

	return r.q.DeleteUserSoft(ctx, userUUID)
}

func (r *userRepository) Restore(
	ctx context.Context,
	userUUID uuid.UUID,
) (db.User, error) {

	traceID := ctx.Value(middleware.TraceIDKey).(string)

	log.Info().
		Str("trace_id", traceID).
		Msg("restoring user")

	return r.q.RestoreUser(ctx, userUUID)
}

func (r *userRepository) Update(
	ctx context.Context,
	arg db.UpdateUserParams,
) (db.User, error) {

	traceID := ctx.Value(middleware.TraceIDKey).(string)

	log.Info().
		Str("trace_id", traceID).
		Msg("updating user")

	return r.q.UpdateUser(ctx, arg)
}

func (r *userRepository) ListUsers(ctx context.Context, arg db.ListUsersOrderByCreatedAtDescParams) ([]db.User, error) {
	return r.q.ListUsersOrderByCreatedAtDesc(ctx, arg)
}

func (r *userRepository) Create(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	traceID := ctx.Value(middleware.TraceIDKey).(string)

	log.Info().
		Str("trace_id", traceID).
		Msg("getting user request")
	return r.q.CreateUser(ctx, arg)
}

func (r *userRepository) GetByUUID(ctx context.Context, userUUID string) (db.User, error) {
	traceID := ctx.Value(middleware.TraceIDKey).(string)

	log.Info().
		Str("trace_id", traceID).
		Msg("getting user by uuid")

	parsedUUID, err := uuid.Parse(userUUID)
	if err != nil {
		return db.User{}, err
	}

	return r.q.GetUserByUUID(ctx, parsedUUID)
}
