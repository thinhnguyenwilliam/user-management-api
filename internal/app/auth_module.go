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
	tokenService auth.ITokenService,
) *AuthModule {

	userRepo := repository.NewUserRepository(store, pool, cache)
	authService := v1service.NewAuthService(userRepo, tokenService)
	authHandler := v1handler.NewAuthHandler(authService)
	authRoutes := v1routes.NewAuthRoutes(authHandler)

	return &AuthModule{
		authRoutes: authRoutes,
	}
}

func (m *AuthModule) PublicRoutes() []routes.Route {
	return []routes.Route{
		m.authRoutes, // ✅ dùng lại
	}
}

func (m *AuthModule) ProtectedRoutes() []routes.Route {
	return nil
}
