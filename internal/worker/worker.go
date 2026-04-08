package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/thinhnguyenwilliam/user-management-api/internal/config"
	"github.com/thinhnguyenwilliam/user-management-api/internal/events"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/logger"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/mail"
	"github.com/thinhnguyenwilliam/user-management-api/pkg/rabbitmq"
)

type Worker struct {
	mq          rabbitmq.RabbitMQService
	mailService mail.EmailProviderService
}

func NewWorker(cfg *config.Config) (*Worker, error) {
	logger.InitLogger(logger.LoggerConfig{
		Level:    "info",
		Filename: "worker.log",
	})

	appLogger := logger.Log

	mq, err := rabbitmq.NewRabbitMQService(cfg.RabbitMQ.URL, appLogger)
	if err != nil {
		return nil, err
	}

	mailService, err := mail.NewMailService(cfg, appLogger)
	if err != nil {
		return nil, err
	}

	return &Worker{
		mq:          mq,
		mailService: mailService,
	}, nil
}

func (w *Worker) Start(ctx context.Context) error {
	log.Println("Worker started...")

	return w.mq.Consume(ctx, "send_email", func(body []byte) error {
		var msg events.EmailMessage

		if err := json.Unmarshal(body, &msg); err != nil {
			return err
		}

		return w.handleMessage(ctx, msg)
	})
}
