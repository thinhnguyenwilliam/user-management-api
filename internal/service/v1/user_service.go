// user-management-api/internal/service/v1/user_service.go
package v1service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"

	db "github.com/thinhnguyenwilliam/user-management-api/internal/db/sqlc"
	"github.com/thinhnguyenwilliam/user-management-api/internal/middleware"
	v1dto "github.com/thinhnguyenwilliam/user-management-api/internal/models/dto/v1"
	"github.com/thinhnguyenwilliam/user-management-api/internal/models/mapper"
	"github.com/thinhnguyenwilliam/user-management-api/internal/repository"
	"github.com/thinhnguyenwilliam/user-management-api/internal/utils"
)

type userService struct {
	userRepo repository.IUserRepository
}

func NewUserService(userRepo repository.IUserRepository) IUserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CountUsers(ctx context.Context, deleted *bool, search *string) (int64, error) {

	// ✅ normalize search
	if search != nil && *search == "" {
		search = nil
	}

	count, err := s.userRepo.CountUsers(ctx, db.CountUsersParams{
		Deleted: deleted,
		Search:  search,
	})

	if err != nil {
		return 0, utils.WrapError("failed to count users", utils.ErrCodeDatabase, err)
	}

	return count, nil
}

func (s *userService) ListUsers(ctx context.Context, limit, offset int32, search *string) ([]db.User, int64, error) {

	// validate
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	users, err := s.userRepo.ListUsers(ctx, db.ListUsersOrderByCreatedAtDescParams{
		Limit:  limit,
		Offset: offset,
		Search: search,
	})

	if err != nil {
		return nil, 0, utils.WrapError("failed to list users", utils.ErrCodeDatabase, err)
	}

	total, err := s.userRepo.CountUsers(ctx, db.CountUsersParams{
		Search: search,
	})
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *userService) DeleteUserSoft(
	ctx context.Context,
	id uuid.UUID,
) (db.User, error) {

	traceID := ctx.Value(middleware.TraceIDKey).(string)

	log.Info().
		Str("trace_id", traceID).
		Msg("delete user request")

	user, err := s.userRepo.DeleteSoft(ctx, id)
	if err != nil {
		return db.User{}, utils.WrapError(
			"failed to delete user",
			utils.ErrCodeDatabase,
			err,
		)
	}

	return user, nil
}

func (s *userService) RestoreUser(
	ctx context.Context,
	id uuid.UUID,
) (db.User, error) {

	traceID := ctx.Value(middleware.TraceIDKey).(string)

	log.Info().
		Str("trace_id", traceID).
		Msg("restore user request")

	user, err := s.userRepo.Restore(ctx, id)
	if err != nil {
		return db.User{}, utils.WrapError(
			"failed to restore user",
			utils.ErrCodeDatabase,
			err,
		)
	}

	return user, nil
}

func (s *userService) UpdateUser(
	ctx context.Context,
	id uuid.UUID,
	req v1dto.UpdateUserRequest,
) (db.User, error) {

	traceID := ctx.Value(middleware.TraceIDKey).(string)
	log.Info().
		Str("trace_id", traceID).
		Msg("update user request")

	if req.Fullname != nil && *req.Fullname == "" {
		return db.User{}, utils.NewError(
			"fullname cannot be empty",
			utils.ErrCodeInvalidInput,
		)
	}

	if req.Password != nil {
		if *req.Password == "" {
			return db.User{}, utils.NewError(
				"Password cannot be empty",
				utils.ErrCodeInvalidInput,
			)
		}
	}

	params := mapper.ToUpdateUserParams(req, id)

	user, err := s.userRepo.Update(ctx, params)
	if err != nil {
		return db.User{}, utils.WrapError(
			"failed to update user",
			utils.ErrCodeDatabase,
			err,
		)
	}

	return user, nil
}

func (s *userService) GetUserByUUID(
	ctx context.Context,
	userUUID string,
) (db.User, error) {

	traceID := ctx.Value(middleware.TraceIDKey).(string)

	log.Info().
		Str("trace_id", traceID).
		Msg("getting user by uuid")

	user, err := s.userRepo.GetByUUID(ctx, userUUID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
		}
		return db.User{}, utils.WrapError("failed to get user", utils.ErrCodeDatabase, err)
	}

	return user, nil
}

func isDuplicateKey(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr.Code == "23505"
	}
	return false
}

func (s *userService) CreateUser(
	ctx context.Context,
	req v1dto.CreateUserRequest,
) (db.User, error) {
	traceID := ctx.Value(middleware.TraceIDKey).(string)

	log.Info().
		Str("trace_id", traceID).
		Msg("getting user request")
	if req.Fullname == "" || req.Email == "" || req.Password == "" {
		return db.User{}, utils.NewError("missing required fields", utils.ErrCodeInvalidInput)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return db.User{}, err
	}

	params := mapper.ToCreateUserParams(req, hashedPassword)

	user, err := s.userRepo.Create(ctx, params)

	if err != nil {

		if isDuplicateKey(err) {
			return db.User{}, utils.NewError("email already exists", utils.ErrCodeConflict)
		}

		return db.User{}, utils.WrapError("failed to create user", utils.ErrCodeDatabase, err)
	}

	return user, nil
}
