// user-management-api/internal/db/connection.go
package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thinhnguyenwilliam/user-management-api/internal/config"
)

func InitDB(cfg *config.Config) (*pgxpool.Pool, error) {
	connString := cfg.DatabaseURL

	conf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	conf.MaxConns = 20
	conf.MinConns = 5
	conf.MaxConnLifetime = time.Hour
	conf.MaxConnIdleTime = 30 * time.Minute
	conf.HealthCheckPeriod = time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("PostgreSQL connected successfully")

	return pool, nil
}
