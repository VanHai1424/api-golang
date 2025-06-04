package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
	Name string `mapstructure:"name"`
}

type Config struct {
	DB DBConfig `mapstructure:"db"` // db
}

func LoadConfig() (*Config, error) {
	mode := os.Getenv("MODE")
	if mode == "" {
		mode = "local"
	}

	viper.SetConfigName(fmt.Sprintf("%s-config", mode))
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}
