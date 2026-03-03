// user-management-api/pkg/rabbitmq/interface.go
package rabbitmq

import "context"

type RabbitMQService interface {
	Publish(ctx context.Context, queue string, body []byte) error
	Consume(ctx context.Context, queue string, handler func([]byte) error) error
	Close() error
}
