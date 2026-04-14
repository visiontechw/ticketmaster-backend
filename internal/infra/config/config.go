package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL   string `mapstructure:"DATABASE_URL"`
	ApiServerPort string `mapstructure:"API_SERVER_PORT"`
	JwtKey        string `mapstructure:"JWT_SECRET"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
