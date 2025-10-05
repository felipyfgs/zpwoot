package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Application
	Port       string
	ServerHost string
	LogLevel   string
	LogFormat  string
	LogOutput  string
	APIKey     string

	// Database
	Database DatabaseConfig

	// PostgreSQL (for Docker services)
	Postgres PostgresConfig

	// WhatsApp/Wameow
	WALogLevel string

	// Webhooks
	GlobalWebhookURL string

	// Environment
	Environment string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL string
}

// PostgresConfig holds PostgreSQL configuration for Docker services
type PostgresConfig struct {
	DB       string
	User     string
	Password string
	Port     int
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	godotenv.Load()

	return &Config{
		// Application
		Port:       getEnv("PORT", "8080"),
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
		LogFormat:  getEnv("LOG_FORMAT", "console"),
		LogOutput:  getEnv("LOG_OUTPUT", "stdout"),
		APIKey:     getEnv("ZP_API_KEY", "a0b1125a0eb3364d98e2c49ec6f7d6ba"),

		// Database
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://zpwoot:zpwoot123@localhost:5432/zpwoot?sslmode=disable"),
		},

		// PostgreSQL (for Docker services)
		Postgres: PostgresConfig{
			DB:       getEnv("POSTGRES_DB", "zpwoot"),
			User:     getEnv("POSTGRES_USER", "zpwoot"),
			Password: getEnv("POSTGRES_PASSWORD", "zpwoot123"),
			Port:     getEnvAsInt("POSTGRES_PORT", 5432),
		},

		// WhatsApp/Wameow
		WALogLevel: getEnv("WA_LOG_LEVEL", "INFO"),

		// Webhooks
		GlobalWebhookURL: getEnv("GLOBAL_WEBHOOK_URL", ""),

		// Environment
		Environment: getEnv("NODE_ENV", "development"),
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// getEnvAsBool gets an environment variable as boolean with a fallback value
func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return fallback
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return c.ServerHost + ":" + c.Port
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
