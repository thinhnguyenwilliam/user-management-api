// user-management-api/internal/app/auth_module.go
package app

import (
	db "github.com/thinhnguyenwilliam/user-management-api/internal/db/sqlc"
	v1handler "github.com/thinhnguyenwilliam/user-management-api/internal/handler/v1"
	"github.com/thinhnguyenwilliam/user-management-api/internal/repository"
	"github.com/thinhnguyenwilliam/user-management-api/internal/routes"
	v1routes "github.com/thinhnguyenwilliam/user-management-api/internal/routes/v1"
	v1service "github.com/thinhnguyenwilliam/user-management-api/internal/service/v1"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/auth"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/rediscache"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthModule struct {
	authRoutes routes.Route
}

func NewAuthModule(
	store *db.Queries,
	pool *pgxpool.Pool,
	cache rediscache.Cache,
) *AuthModule {

	// 1. Repo
	userRepo := repository.NewUserRepository(store, pool, cache)

	// 2. Token service (🔥 bạn đang thiếu cái này)
	// tokenService := auth.NewJWTService(cfg.JWT.Secret)
	tokenService := auth.NewJWTService("your-secret-key")

	// 3. Service
	authService := v1service.NewAuthService(userRepo, tokenService)

	// 4. Handler
	authHandler := v1handler.NewAuthHandler(authService)

	// 5. Routes
	authRoutes := v1routes.NewAuthRoutes(authHandler)

	return &AuthModule{
		authRoutes: authRoutes,
	}
}

func (m *AuthModule) Route() routes.Route {
	return m.authRoutes
}
