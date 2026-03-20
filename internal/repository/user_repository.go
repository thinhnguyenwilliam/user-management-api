// user-management-api/internal/repository/user_repository.go
package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	db "github.com/thinhnguyenwilliam/user-management-api/internal/db/sqlc"
	"github.com/thinhnguyenwilliam/user-management-api/internal/middleware"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/rediscache"
)

type userRepository struct {
	q     db.Querier
	pool  *pgxpool.Pool
	cache rediscache.Cache
}

func NewUserRepository(
	q db.Querier,
	pool *pgxpool.Pool,
	cache rediscache.Cache,
) IUserRepository {
	return &userRepository{
		q:     q,
		pool:  pool,
		cache: cache,
	}
}

func buildSearchKey(s *string) string {
	if s == nil {
		return "no-filter"
	}

	val := strings.ToLower(strings.TrimSpace(*s))
	if val == "" {
		return "empty"
	}

	return val
}

func (r *userRepository) ListUsers(ctx context.Context, arg db.ListUsersOrderByCreatedAtDescParams) ([]db.User, error) {
	key := fmt.Sprintf("users:list:%d:%d:%s",
		arg.Limit,
		arg.Offset,
		buildSearchKey(arg.Search),
	)

	var users []db.User

	// 1. try cache
	err := r.cache.Get(ctx, key, &users)
	if err == nil {
		log.Info().Msg("CACHE HIT")
		return users, nil
	}

	// 2. fallback DB
	log.Info().Err(err).Msg("CACHE MISS")
	users, err = r.q.ListUsersOrderByCreatedAtDesc(ctx, arg)
	if err != nil {
		return nil, err
	}

	// 3. set cache (TTL ngắn thôi)
	log.Info().Msg("SET CACHE")
	_ = r.cache.Set(ctx, key, users, 30*time.Second)

	return users, nil
}

type ListUsersInput struct {
	Limit  int32
	Offset int32
	Search *string
	Sort   string
	Order  string
}

func safeSort(sort string) string {
	switch sort {
	case "user_created_at", "user_email", "user_fullname":
		return sort
	default:
		return "user_created_at"
	}
}

func safeOrder(order string) string {
	if order == "asc" {
		return "ASC"
	}
	return "DESC"
}

func (r *userRepository) ListUsersV2(ctx context.Context, in ListUsersInput) ([]db.User, error) {

	sort := safeSort(in.Sort)
	order := safeOrder(in.Order)

	query := fmt.Sprintf(`
	SELECT *
	FROM users
	WHERE user_deleted_at IS NULL
	AND (
		$3::text IS NULL
		OR $3 = ''
		OR user_email ILIKE '%%' || $3 || '%%'
		OR user_fullname ILIKE '%%' || $3 || '%%'
	)
	ORDER BY %s %s
	LIMIT $1 OFFSET $2;
	`, sort, order)

	rows, err := r.pool.Query(ctx, query, in.Limit, in.Offset, in.Search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []db.User
	for rows.Next() {
		var u db.User
		if err := rows.Scan(
			&u.UserUuid,
			&u.UserFullname,
			&u.UserEmail,
			&u.UserPassword,
			&u.UserAge,
			&u.UserStatus,
			&u.UserLevel,
			&u.UserDeletedAt,
			&u.UserCreatedAt,
			&u.UserUpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *userRepository) CountUsers(ctx context.Context, arg db.CountUsersParams) (int64, error) {
	key := fmt.Sprintf("users:count:%v", buildSearchKey(arg.Search))

	var total int64

	if err := r.cache.Get(ctx, key, &total); err == nil {
		log.Info().Msg("CACHE HIT COUNT")
		return total, nil
	}

	log.Info().Msg("CACHE MISS COUNT")
	total, err := r.q.CountUsers(ctx, arg)
	if err != nil {
		return 0, err
	}

	_ = r.cache.Set(ctx, key, total, 1*time.Minute)

	return total, nil
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

func (r *userRepository) Create(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	traceID, _ := ctx.Value(middleware.TraceIDKey).(string)

	log.Info().
		Str("trace_id", traceID).
		Msg("creating user")

	user, err := r.q.CreateUser(ctx, arg)
	if err != nil {
		return user, err
	}

	// invalidate cache AFTER success
	if err := r.cache.DeleteByPattern(ctx, "users:list:*"); err != nil {
		log.Error().Err(err).Msg("failed to delete cache")
	}
	_ = r.cache.DeleteByPattern(ctx, "users:count:*")

	return user, nil
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
