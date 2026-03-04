// user-management-api/internal/app/app.go
package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"

	"github.com/thinhnguyenwilliam/user-management-api/internal/cache"
	"github.com/thinhnguyenwilliam/user-management-api/internal/config"
	"github.com/thinhnguyenwilliam/user-management-api/internal/routes"
)

type Module interface {
	Route() routes.Route
}

type Application struct {
	config *config.Config
	router *gin.Engine
}

func NewApplication(cfg *config.Config) (*Application, error) {
	// 1️⃣ Init Redis
	rdb, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		return nil, err
	}
	log.Println("Redis Addr:", cfg.Redis.Addr)
	log.Println("Redis Pass:", cfg.Redis.Password)
	_ = rdb

	// 2️⃣ Init RabbitMQ
	conn, err := amqp091.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		return nil, err
	}
	_ = conn
	log.Println("Connected to RabbitMQ")

	// 3️⃣ Init Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.SetTrustedProxies(nil)

	// Load modules
	modules := []Module{
		NewUserModule(),
	}
	routeList := collectRoutes(modules)
	routes.RegisterRoutes(r, routeList...)

	return &Application{
		config: cfg,
		router: r,
	}, nil
}

func collectRoutes(modules []Module) []routes.Route {
	var result []routes.Route

	for _, module := range modules {
		result = append(result, module.Route())
	}

	return result
}

func (a *Application) Run() error {
	log.Println("Server running on port:", a.config.Port)
	return a.router.Run(":" + a.config.Port)
}
