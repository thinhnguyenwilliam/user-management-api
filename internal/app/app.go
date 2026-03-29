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
	"github.com/thinhnguyenwilliam/user-management-api/pkg/auth"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/rediscache"
)

type Module interface {
	PublicRoutes() []routes.Route
	ProtectedRoutes() []routes.Route
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

	if err := validation.InitValidation(); err != nil {
		return nil, err
	}

	tokenService := auth.NewJWTService(
		"your-signing-secret",
		[]byte("12345678901234567890123456789012"),
	)

	// Redis
	rdb, _ := cache.NewRedisClient(cfg.Redis)
	cacheService := rediscache.New(rdb)

	// RabbitMQ
	conn, err := amqp091.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		return nil, err
	}

	// Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Middlewares
	rateLimiter := middleware.NewRateLimiter(5, 15, 3*time.Minute)

	middlewares := []gin.HandlerFunc{
		middleware.Recovery(),
		middleware.TraceID(),
		middleware.LoggerMiddleware(),
		middleware.ApiKeyMiddleware(cfg.ApiKey),
		rateLimiter.Middleware(),
	}

	store := sqlc.New(pool)

	// ✅ modules
	modules := []Module{
		NewUserModule(store, pool, cacheService),
		NewAuthModule(store, pool, cacheService, tokenService),
	}

	// ✅ collect routes đúng cách
	var publicRoutes []routes.Route
	var protectedRoutes []routes.Route

	for _, m := range modules {
		publicRoutes = append(publicRoutes, m.PublicRoutes()...)
		protectedRoutes = append(protectedRoutes, m.ProtectedRoutes()...)
	}

	// ✅ register
	routes.RegisterRoutes(
		r,
		middlewares,
		tokenService,
		publicRoutes,
		protectedRoutes,
	)

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
