// user-management-api/pkg/logger/logger.go
package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerConfig struct {
	Level      string
	FileName   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

func NewLogger(config LoggerConfig) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	// file rotation
	fileWriter := &lumberjack.Logger{
		Filename:   config.FileName,
		MaxSize:    config.MaxSize, // MB
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge, // days
		Compress:   config.Compress,
	}

	multiWriter := io.MultiWriter(os.Stdout, fileWriter)

	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	logger := zerolog.New(multiWriter).
		Level(level).
		With().
		Timestamp().
		Logger()

	return logger
}
