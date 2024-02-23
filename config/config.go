package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config represents the overall configuration structure.
type Config struct {
	PostgresDB     *Database
	HTTP           *Endpoint
	UserService    *Endpoint
	ProductService *Endpoint
	CouponService  *Endpoint
	GatewayService *Endpoint
	SymetricKey    string
	FileLogOutPut  string
}

// config is a private structure used for unmarshaling the configuration from Viper.
type config struct {
	DBHost          string `mapstructure:"DB_HOST"`
	DBPort          string `mapstructure:"DB_PORT"`
	DBUser          string `mapstructure:"DB_USER"`
	DBPassword      string `mapstructure:"DB_PASSWORD"`
	DBDatabase      string `mapstructure:"DB_NAME"`
	UserGRPCHost    string `mapstructure:"USER_GRPC_HOST"`
	UserGRPCPort    string `mapstructure:"USER_GRPC_PORT"`
	ProductGRPCHost string `mapstructure:"PRODUCT_GRPC_HOST"`
	ProductGRPCPort string `mapstructure:"PRODUCT_GRPC_PORT"`
	CouponGRPCHost  string `mapstructure:"COUPON_GRPC_HOST"`
	CouponGRPCPort  string `mapstructure:"COUPON_GRPC_PORT"`
	GatewayGRPCHost string `mapstructure:"GATEWAY_GRPC_HOST"`
	GatewayGRPCPort string `mapstructure:"GATEWAY_GRPC_PORT"`
	SymetricKey     string `mapstructure:"SYMETRIC_KEY"`
	FileLogOutPut   string `mapstructure:"FILE_LOG_OUTPUT"`
}

// LoadConfig loads the configuration from the specified file path and environment.
func LoadConfig(path string, env string) (*Config, error) {
	// Initialize an instance of the private config structure.
	var cfg config

	// Add the configuration path and name to Viper.
	viper.AddConfigPath(path)
	viper.SetConfigName(env)
	viper.SetConfigType("env")

	// Automatically read environment variables that match the structure field names.
	viper.AutomaticEnv()

	// Read the configuration from the file.
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("unable to read config file: %w", err)
	}

	// Unmarshal the configuration into the private config structure.
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config file: %w", err)
	}

	// Create and return the public Config structure based on the private config.
	return &Config{
		PostgresDB: &Database{
			Host:     cfg.DBHost,
			Port:     cfg.DBPort,
			User:     cfg.DBUser,
			Password: cfg.DBPassword,
			Database: cfg.DBDatabase,
		},
		UserService: &Endpoint{
			Host: cfg.UserGRPCHost,
			Port: cfg.UserGRPCPort,
		},
		ProductService: &Endpoint{
			Host: cfg.ProductGRPCHost,
			Port: cfg.ProductGRPCPort,
		},
		CouponService: &Endpoint{
			Host: cfg.CouponGRPCHost,
			Port: cfg.CouponGRPCPort,
		},
		GatewayService: &Endpoint{
			Host: cfg.GatewayGRPCHost,
			Port: cfg.GatewayGRPCPort,
		},
		SymetricKey:   cfg.SymetricKey,
		FileLogOutPut: cfg.FileLogOutPut,
	}, nil
}
