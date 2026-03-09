// user-management-api/internal/app/user_module.go
package app

import (
	"github.com/thinhnguyenwilliam/user-management-api/internal/db"
	"github.com/thinhnguyenwilliam/user-management-api/internal/handler"
	"github.com/thinhnguyenwilliam/user-management-api/internal/repository"
	"github.com/thinhnguyenwilliam/user-management-api/internal/routes"
	"github.com/thinhnguyenwilliam/user-management-api/internal/service"
)

type UserModule struct {
	userRoutes routes.Route
}

func NewUserModule() *UserModule {
	userRepo := repository.NewUserRepository(db.Store)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	userRoutes := routes.NewUserRoutes(userHandler)

	return &UserModule{
		userRoutes: userRoutes,
	}
}

func (m *UserModule) Route() routes.Route {
	return m.userRoutes
}
