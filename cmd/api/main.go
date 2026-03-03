// user-management-api/cmd/api/main.go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	"github.com/thinhnguyenwilliam/user-management-api/internal/cache"
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

	//
	rdb, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatal("cannot connect to redis:", err)
	}

	_ = rdb // inject vào service sau
	log.Println("Redis Addr:", cfg.Redis.Addr)
	log.Println("Redis Pass:", cfg.Redis.Password)

	//
	conn, err := amqp091.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatal("cannot connect to rabbitmq:", err)
	}
	defer conn.Close()
	log.Println("Connected to RabbitMQ")

	// 2️⃣ Create Gin engine
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.SetTrustedProxies(nil)

	// 3️⃣ Dependency Injection
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// 4️⃣ Register routes
	userRoutes := routes.NewUserRoutes(userHandler)
	routes.RegisterRoutes(r, userRoutes)

	// 5️⃣ Start server using config port
	log.Println("Server running on port:", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
