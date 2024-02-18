package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	PostgresDB *Database
	HTTP       *Endpoint

	SymetricKey        string
	SuperAdminUsername string
	SuperAdminPassword string
}

type config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBDatabase string `mapstructure:"DB_NAME"`

	HTTPHost string `mapstructure:"HTTP_HOST"`
	HTTPPort string `mapstructure:"HTTP_PORT"`

	SymetricKey string `mapstructure:"SYMETRIC_KEY"`

	SuperAdminUsername string `mapstructure:"SUPER_ADMIN_USERNAME"`
	SuperAdminPassword string `mapstructure:"SUPER_ADMIN_PASSWORD"`
}

func LoadConfig(path string, env string) (*Config, error) {
	var cfg config
	viper.AddConfigPath(path)
	viper.SetConfigName(env)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config file: %w", err)
	}

	return &Config{
		PostgresDB: &Database{
			Host:     cfg.DBHost,
			Port:     cfg.DBPort,
			User:     cfg.DBUser,
			Password: cfg.DBPassword,
			Database: cfg.DBDatabase,
		},
		HTTP: &Endpoint{
			Host: cfg.HTTPHost,
			Port: cfg.HTTPPort,
		},
		SymetricKey:        cfg.SymetricKey,
		SuperAdminUsername: cfg.SuperAdminUsername,
		SuperAdminPassword: cfg.SuperAdminPassword,
	}, nil
}
