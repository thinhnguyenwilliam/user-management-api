// user-management-api/internal/config/redis.go
package config

type RedisConfig struct {
	Addr     string `mapstructure:"REDIS_ADDR"`
	User     string `mapstructure:"REDIS_USER"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}
