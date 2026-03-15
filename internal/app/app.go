// user-management-api/internal/app/app.go
package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"

	"github.com/thinhnguyenwilliam/user-management-api/internal/cache"
	"github.com/thinhnguyenwilliam/user-management-api/internal/config"
	sqlc "github.com/thinhnguyenwilliam/user-management-api/internal/db/sqlc"
	"github.com/thinhnguyenwilliam/user-management-api/internal/middleware"
	"github.com/thinhnguyenwilliam/user-management-api/internal/routes"
	"github.com/thinhnguyenwilliam/user-management-api/internal/validation"
)

type Module interface {
	Route() routes.Route
}

type Application struct {
	config     *config.Config
	router     *gin.Engine
	modules    []Module
	rabbitConn *amqp091.Connection
	dbPool     *pgxpool.Pool
}

func NewApplication(cfg *config.Config, pool *pgxpool.Pool) (*Application, error) {
	validateConfig(cfg)

	// ✅ Init validation
	if err := validation.InitValidation(); err != nil {
		return nil, err
	}

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
	log.Println("Connected to RabbitMQ")

	// 3️⃣ Init Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Middlewares
	rateLimiter := middleware.NewRateLimiter(
		5,             // 5 requests/sec
		15,            // burst
		3*time.Minute, // client TTL
	)

	middlewares := []gin.HandlerFunc{
		middleware.Recovery(),
		middleware.TraceID(),
		middleware.LoggerMiddleware(),
		middleware.ApiKeyMiddleware(cfg.ApiKey),
		rateLimiter.Middleware(),
	}

	store := sqlc.New(pool)
	// Load modules (inject dependencies properly)
	modules := []Module{
		NewUserModule(store),
	}
	routeList := collectRoutes(modules)
	routes.RegisterRoutes(r, middlewares, routeList...)

	return &Application{
		config:     cfg,
		router:     r,
		modules:    modules,
		rabbitConn: conn,
		dbPool:     pool,
	}, nil
}

func validateConfig(cfg *config.Config) {
	if cfg.Port == "" {
		log.Fatal("PORT is required")
	}
	if cfg.ApiKey == "" {
		log.Fatal("API_KEY is required")
	}
	if cfg.RabbitMQ.URL == "" {
		log.Fatal("RABBITMQ_URL is required")
	}
}

func collectRoutes(modules []Module) []routes.Route {
	var result []routes.Route

	for _, module := range modules {
		result = append(result, module.Route())
	}

	return result
}

func (a *Application) Run() error {
	addr := ":" + a.config.Port

	server := &http.Server{
		Addr:           addr,
		Handler:        a.router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Println("Server running on", addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}

func (a *Application) Shutdown() {
	if a.rabbitConn != nil {
		_ = a.rabbitConn.Close()
	}
	log.Println("Application shutdown complete")
}
