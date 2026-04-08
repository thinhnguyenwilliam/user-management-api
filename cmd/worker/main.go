// user-management-api/cmd/worker/main.go
package main

import (
	"context"
	"log"

	"github.com/thinhnguyenwilliam/user-management-api/internal/config"
	"github.com/thinhnguyenwilliam/user-management-api/internal/worker"
)

func main() {
	cfg, _ := config.LoadConfig()

	ctx := context.Background()

	w, _ := worker.NewWorker(cfg)

	if err := w.Start(ctx); err != nil {
		log.Fatal(err)
	}

	// 👇 block process
	select {}
}

// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"log"

// 	"github.com/thinhnguyenwilliam/user-management-api/internal/config"
// 	"github.com/thinhnguyenwilliam/user-management-api/internal/events"
// 	"github.com/thinhnguyenwilliam/user-management-api/pkg/logger"
// 	"github.com/thinhnguyenwilliam/user-management-api/pkg/mail"
// 	"github.com/thinhnguyenwilliam/user-management-api/pkg/rabbitmq"
// )

// func main() {
// 	cfg, _ := config.LoadConfig()
// 	ctx := context.Background()

// 	// init logger
// 	logger.InitLogger(logger.LoggerConfig{
// 		Level:      "info",
// 		Filename:   "app.log",
// 		MaxSize:    10,
// 		MaxBackups: 5,
// 		MaxAge:     30,
// 		Compress:   true,
// 		IsDev:      "development",
// 	})

// 	appLogger := logger.Log

// 	mq, _ := rabbitmq.NewRabbitMQService(cfg.RabbitMQ.URL, appLogger)
// 	mailService, _ := mail.NewMailService(cfg, appLogger)

// 	log.Println("Worker started...")

// 	err := mq.Consume(ctx, "send_email", func(body []byte) error {
// 		var msg events.EmailMessage
// 		if err := json.Unmarshal(body, &msg); err != nil {
// 			return err
// 		}

// 		return handleEmail(ctx, msg, mailService)
// 	})

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	select {}
// }

// func handleEmail(
// 	ctx context.Context,
// 	msg events.EmailMessage,
// 	mailService mail.EmailProviderService,
// ) error {

// 	switch msg.Type {

// 	case "reset_password":
// 		return sendResetPassword(ctx, msg, mailService)

// 	// case "welcome_email":
// 	// 	return sendWelcomeEmail(ctx, msg, mailService)

// 	default:
// 		return nil // ignore unknown type
// 	}
// }

// func sendResetPassword(
// 	ctx context.Context,
// 	msg events.EmailMessage,
// 	mailService mail.EmailProviderService,
// ) error {

// 	mailContent := &mail.Email{
// 		To: []mail.Address{
// 			{Email: msg.To},
// 		},
// 		Subject: msg.Subject,
// 		Text:    msg.Body,
// 	}

// 	return mailService.SendMail(ctx, mailContent)
// }
