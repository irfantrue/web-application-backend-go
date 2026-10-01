// Package config provides centralized application configuration management.
package config

import (
	"github.com/spf13/viper"
)

// Config is the main structure that combines all application configurations (App, DB, and Redis).
type Config struct {
	App      AppConfig      `mapstructure:",squash"`
	DB       DBConfig       `mapstructure:",squash"`
	RabbitMQ RabbitMQConfig `mapstructure:",squash"`
	Redis    RedisConfig    `mapstructure:",squash"`
}

// AppConfig holds specific application settings such as port, environment, and CORS.
type AppConfig struct {
	AppName            string `mapstructure:"APP_NAME"`
	ServerHeader       string `mapstructure:"SERVER_HEADER"`
	Port               string `mapstructure:"APP_PORT"`
	Env                string `mapstructure:"APP_ENV"`
	CompressLevel      int    `mapstructure:"COMPRESS_LEVEL"`
	CorsAllowedOrigins string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	JwtSecret          string `mapstructure:"JWT_SECRET"`
	JwtExpire          int    `mapstructure:"JWT_EXPIRE_HOURS"`
}

// DBConfig holds the credentials and connection settings for the database.
type DBConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	Port     string `mapstructure:"DB_PORT"`
}

// RabbitMQConfig holds the configuration settings for connecting to RabbitMQ.
type RabbitMQConfig struct {
	URL string `mapstructure:"RABBITMQ_URL"`
}

// RedisConfig holds the credentials and connection settings for Redis.
type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
}

// LoadConfig reads the configuration file from the given path along with environment variables,
// and unmarshals them into the Config struct.
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)

	viper.AutomaticEnv()

	var cfg Config

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
