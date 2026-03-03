// user-management-api/internal/config/config.go
package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string      `mapstructure:"PORT"`
	DatabaseURL string      `mapstructure:"DATABASE_URL"`
	JWTSecret   string      `mapstructure:"JWT_SECRET"`
	Redis       RedisConfig `mapstructure:",squash"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using system env")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
