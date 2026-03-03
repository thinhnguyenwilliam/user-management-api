package config

type RabbitMQConfig struct {
	URL      string `mapstructure:"RABBITMQ_URL"`
	Exchange string `mapstructure:"RABBITMQ_EXCHANGE"`
	Queue    string `mapstructure:"RABBITMQ_QUEUE"`
}
