// user-management-api/pkg/rabbitmq/rabbitmq.go
package rabbitmq

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type rabbitMQService struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	logger  *zerolog.Logger
}

func NewRabbitMQService(amqpURL string, logger *zerolog.Logger) (RabbitMQService, error) {
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect to rabbitmq")
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		logger.Error().Err(err).Msg("failed to open channel")
		return nil, err
	}

	logger.Info().Msg("rabbitmq connected")
	return &rabbitMQService{
		conn:    conn,
		channel: ch,
		logger:  logger,
	}, nil
}

func (r *rabbitMQService) Close() error {
	if err := r.channel.Close(); err != nil {
		r.logger.Error().Err(err).Msg("failed to close channel")
		return err
	}

	if err := r.conn.Close(); err != nil {
		r.logger.Error().Err(err).Msg("failed to close connection")
		return err
	}

	r.logger.Info().Msg("rabbitmq connection closed")
	return nil
}

func (r *rabbitMQService) Publish(
	ctx context.Context,
	queue string,
	body []byte,
) error {

	_, err := r.channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return r.channel.PublishWithContext(
		ctx,
		"",
		queue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *rabbitMQService) Consume(
	ctx context.Context,
	queue string,
	handler func([]byte) error,
) error {

	_, err := r.channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := r.channel.Consume(
		queue,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				msg.Nack(false, true)
				continue
			}
			msg.Ack(false)
		}
	}()

	return nil
}
