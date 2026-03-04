// user-management-api/cmd/api/main.go
package main

import (
	"log"

	"github.com/thinhnguyenwilliam/user-management-api/internal/app"
	"github.com/thinhnguyenwilliam/user-management-api/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatal("cannot create application:", err)
	}

	if err := application.Run(); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
