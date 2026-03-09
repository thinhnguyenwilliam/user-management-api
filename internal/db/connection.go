// user-management-api/internal/db/connection.go
package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thinhnguyenwilliam/user-management-api/internal/config"
	sqlc "github.com/thinhnguyenwilliam/user-management-api/internal/db/sqlc"
)

var (
	DBPool *pgxpool.Pool
	Store  *sqlc.Queries
)

func InitDBV2(cfg *config.Config) error {
	connString := cfg.DatabaseURL

	conf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return err
	}

	// pool settings
	conf.MaxConns = 20
	conf.MinConns = 5
	conf.MaxConnLifetime = time.Hour
	conf.MaxConnIdleTime = 30 * time.Minute
	conf.HealthCheckPeriod = time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return err
	}

	// test connection
	if err := pool.Ping(ctx); err != nil {
		return err
	}

	DBPool = pool
	Store = sqlc.New(pool)

	log.Println("PostgreSQL connected successfully")

	return nil
}

// var DB *gorm.DB

// func InitDB(cfg *config.Config) error {

// 	dsn := fmt.Sprintf(
// 		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
// 		cfg.DB.Host,
// 		cfg.DB.User,
// 		cfg.DB.Password,
// 		cfg.DB.Name,
// 		cfg.DB.Port,
// 	)

// 	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		return err
// 	}

// 	sqlDB, err := database.DB()
// 	if err != nil {
// 		return err
// 	}

// 	// Connection Pool
// 	sqlDB.SetMaxOpenConns(10)
// 	sqlDB.SetMaxIdleConns(5)

// 	// Lifetime settings
// 	sqlDB.SetConnMaxLifetime(30 * time.Minute)
// 	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

// 	// Context timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	// Ping database
// 	if err := sqlDB.PingContext(ctx); err != nil {
// 		return err
// 	}

// 	DB = database

// 	log.Println("Database connected successfully")

// 	return nil
// }
