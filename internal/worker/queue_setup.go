package worker

import "github.com/rabbitmq/amqp091-go"

const (
	QueueEmail      = "send_email"
	QueueEmailRetry = "send_email_retry"
	QueueEmailDLQ   = "send_email_dlq"
)

func SetupQueues(ch *amqp091.Channel) error {

	// main queue
	_, err := ch.QueueDeclare(
		"send_email",
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-dead-letter-routing-key": "send_email_retry",
		},
	)
	if err != nil {
		return err
	}

	// retry queue (delay)
	_, err = ch.QueueDeclare(
		"send_email_retry",
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-message-ttl":             int32(5000),
			"x-dead-letter-routing-key": "send_email",
		},
	)
	if err != nil {
		return err
	}

	// DLQ
	_, err = ch.QueueDeclare(
		"send_email_dlq",
		true,
		false,
		false,
		false,
		nil,
	)

	return err
}
