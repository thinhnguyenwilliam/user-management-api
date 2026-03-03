// user-management-api/cmd/api/main.go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/thinhnguyenwilliam/user-management-api/internal/config"
	"github.com/thinhnguyenwilliam/user-management-api/internal/handler"
	"github.com/thinhnguyenwilliam/user-management-api/internal/repository"
	"github.com/thinhnguyenwilliam/user-management-api/internal/routes"
	"github.com/thinhnguyenwilliam/user-management-api/internal/service"
)

func main() {
	// 1️⃣ Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// 2️⃣ Create Gin engine
	r := gin.Default()

	// 3️⃣ Dependency Injection
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// 4️⃣ Register routes
	userRoutes := routes.NewUserRoutes(r, userHandler)
	userRoutes.Register()

	// 5️⃣ Start server using config port
	log.Println("Server running on port:", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
