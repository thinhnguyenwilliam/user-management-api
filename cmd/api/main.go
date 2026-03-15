// user-management-api/cmd/api/main.go
// sudo lsof -i :8086
// sudo kill -9 130765
package main

import (
	"log"
	"os"

	"github.com/rs/zerolog"
	"github.com/thinhnguyenwilliam/user-management-api/internal/app"
	"github.com/thinhnguyenwilliam/user-management-api/internal/config"
	"github.com/thinhnguyenwilliam/user-management-api/internal/db"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	os.MkdirAll("logs", 0755)

	logFile, err := os.OpenFile(
		"logs/app.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}

	multi := zerolog.MultiLevelWriter(os.Stdout, logFile)

	logger := zerolog.New(multi).
		With().
		Timestamp().
		Logger()

	dbPool, err := db.InitDB(cfg, logger)
	if err != nil {
		log.Fatal(err)
	}

	application, err := app.NewApplication(cfg, dbPool)
	if err != nil {
		log.Fatal("cannot create application:", err)
	}

	if err := application.Run(); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
