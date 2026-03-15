// user-management-api/internal/db/connection.go
package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"

	"github.com/thinhnguyenwilliam/user-management-api/internal/config"
	"github.com/thinhnguyenwilliam/user-management-api/internal/middleware"
)

type PgLogger struct {
	logger zerolog.Logger
}

func (l PgLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {

	traceID := ""
	if v := ctx.Value(middleware.TraceIDKey); v != nil {
		if id, ok := v.(string); ok {
			traceID = id
		}
	}

	var event *zerolog.Event

	switch level {
	case tracelog.LogLevelError:
		event = l.logger.Error()
	case tracelog.LogLevelWarn:
		event = l.logger.Warn()
	case tracelog.LogLevelInfo:
		event = l.logger.Info()
	default:
		event = l.logger.Debug()
	}

	event.
		Str("trace_id", traceID).
		Str("event", msg).
		Interface("data", data).
		Msg("pgx")

	if duration, ok := data["time"].(time.Duration); ok {
		if duration > 100*time.Millisecond {
			l.logger.Warn().
				Str("trace_id", traceID).
				Dur("duration", duration).
				Interface("data", data).
				Msg("slow query detected")
		}
	}
}

func newPoolConfig(connString string, logger zerolog.Logger) (*pgxpool.Config, error) {

	conf, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	conf.MaxConns = 20
	conf.MinConns = 5
	conf.MaxConnLifetime = time.Hour
	conf.MaxConnIdleTime = 30 * time.Minute
	conf.HealthCheckPeriod = time.Minute

	conf.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   PgLogger{logger},
		LogLevel: tracelog.LogLevelDebug,
	}

	return conf, nil
}

func InitDB(cfg *config.Config, logger zerolog.Logger) (*pgxpool.Pool, error) {

	conf, err := newPoolConfig(cfg.DatabaseURL, logger)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	logger.Info().Msg("PostgreSQL connected successfully")

	return pool, nil
}
