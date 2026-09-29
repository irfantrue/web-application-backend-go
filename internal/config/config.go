package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	App AppConfig `mapstructure:",squash"`
}

type AppConfig struct {
	Port string `mapstructure:"APP_PORT"`
	Env  string `mapstructure:"APP_ENV"`
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
