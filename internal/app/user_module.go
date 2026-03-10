// user-management-api/internal/app/user_module.go
package app

import (
	"github.com/thinhnguyenwilliam/user-management-api/internal/db"
	v1handler "github.com/thinhnguyenwilliam/user-management-api/internal/handler/v1"
	"github.com/thinhnguyenwilliam/user-management-api/internal/repository"
	"github.com/thinhnguyenwilliam/user-management-api/internal/routes"
	v1routes "github.com/thinhnguyenwilliam/user-management-api/internal/routes/v1"
	v1service "github.com/thinhnguyenwilliam/user-management-api/internal/service/v1"
)

type UserModule struct {
	userRoutes routes.Route
}

func NewUserModule() *UserModule {
	userRepo := repository.NewUserRepository(db.Store)
	userService := v1service.NewUserService(userRepo)
	userHandler := v1handler.NewUserHandler(userService)
	userRoutes := v1routes.NewUserRoutes(userHandler)

	return &UserModule{
		userRoutes: userRoutes,
	}
}

func (m *UserModule) Route() routes.Route {
	return m.userRoutes
}
