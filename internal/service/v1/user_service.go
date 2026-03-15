// user-management-api/internal/service/v1/user_service.go
package v1service

import (
	"context"

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
