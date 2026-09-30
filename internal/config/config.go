package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig   `mapstructure:",squash"`
	DB    DBConfig    `mapstructure:",squash"`
	Redis RedisConfig `mapstructure:",squash"`
}

type AppConfig struct {
	Port               string `mapstructure:"APP_PORT"`
	Env                string `mapstructure:"APP_ENV"`
	CompressLevel      int    `mapstructure:"COMPRESS_LEVEL"`
	CorsAllowedOrigins string `mapstructure:"CORS_ALLOWED_ORIGINS"`
}

type DBConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	Port     string `mapstructure:"DB_PORT"`
}

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
}

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
