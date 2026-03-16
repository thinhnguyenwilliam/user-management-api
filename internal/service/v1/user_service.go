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
