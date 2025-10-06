package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	ServerHost string
	LogLevel   string
	LogFormat  string
	LogOutput  string
	APIKey     string

	Database DatabaseConfig

	Postgres PostgresConfig

	WALogLevel string

	GlobalWebhookURL string

	Environment string
}

type DatabaseConfig struct {
	URL string
}

type PostgresConfig struct {
	DB       string
	User     string
	Password string
	Port     int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	cfg := &Config{

		Port:       getEnv("PORT", "8080"),
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		LogFormat:  getEnv("LOG_FORMAT", "console"),
		LogOutput:  getEnv("LOG_OUTPUT", "stdout"),
		APIKey:     getEnv("ZP_API_KEY", ""),

		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},

		Postgres: PostgresConfig{
			DB:       getEnv("POSTGRES_DB", ""),
			User:     getEnv("POSTGRES_USER", ""),
			Password: getEnv("POSTGRES_PASSWORD", ""),
			Port:     getEnvAsInt("POSTGRES_PORT", 5432),
		},

		WALogLevel: getEnv("WA_LOG_LEVEL", "INFO"),

		GlobalWebhookURL: getEnv("GLOBAL_WEBHOOK_URL", ""),

		Environment: getEnv("NODE_ENV", "development"),
	}

	if err := cfg.Validate(); err != nil {
		panic("Configuration validation failed: " + err.Error())
	}

	return cfg
}

func (c *Config) Validate() error {
	if c.APIKey == "" {
		return errors.New("ZP_API_KEY is required")
	}

	if c.Database.URL == "" {
		return errors.New("DATABASE_URL is required")
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}

	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}

	return fallback
}

func (c *Config) GetServerAddress() string {
	return c.ServerHost + ":" + c.Port
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
