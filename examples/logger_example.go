package main

import (
	"context"
	"errors"
	"time"

	"zpwoot/internal/adapters/config"
	"zpwoot/internal/adapters/logger"

	"github.com/rs/zerolog"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger with full configuration
	logger.InitWithConfig(cfg)

	// =============================================================================
	// BASIC LOGGING - Using global logger functions
	// =============================================================================
	
	logger.Info().Msg("Application started")
	logger.Debug().Msg("Debug information")
	logger.Warn().Msg("Warning message")
	logger.Error().Msg("Error occurred")

	// =============================================================================
	// STRUCTURED LOGGING - Adding fields to logs
	// =============================================================================
	
	logger.Info().
		Str("version", "1.0.0").
		Int("port", 8080).
		Bool("debug", true).
		Msg("Server configuration")

	// With error
	err := errors.New("database connection failed")
	logger.Error().
		Err(err).
		Str("component", "database").
		Str("host", "localhost").
		Int("port", 5432).
		Msg("Failed to connect to database")

	// =============================================================================
	// CONTEXTUAL LOGGING - Creating loggers with context
	// =============================================================================
	
	// Logger with component context
	apiLogger := logger.WithComponent("api")
	apiLogger.Info().
		Str("method", "GET").
		Str("path", "/health").
		Int("status", 200).
		Dur("duration", 15*time.Millisecond).
		Msg("API request processed")

	// Logger with request ID
	requestLogger := logger.WithRequestID("req-12345")
	requestLogger.Info().
		Str("user_id", "user-456").
		Str("action", "login").
		Msg("User authentication")

	// Logger with session ID
	sessionLogger := logger.WithSessionID("sess-789")
	sessionLogger.Info().
		Str("event", "session_created").
		Msg("New session established")

	// Logger with multiple fields
	contextLogger := logger.WithFields(map[string]interface{}{
		"request_id": "req-123",
		"user_id":    "user-456",
		"session_id": "sess-789",
		"ip_address": "192.168.1.100",
	})
	contextLogger.Info().
		Str("endpoint", "/api/v1/users").
		Int("response_size", 1024).
		Msg("Processing request")

	// =============================================================================
	// INSTANCE LOGGING - Using logger instances
	// =============================================================================
	
	log := logger.New()
	log.Info().
		Str("service", "zpwoot").
		Msg("Service instance created")

	// Instance with context
	dbLogger := logger.NewFromAppConfig(cfg).WithComponent("database")
	dbLogger.Info().
		Str("operation", "migration").
		Int("version", 5).
		Msg("Running database migration")

	// =============================================================================
	// DIFFERENT LOG LEVELS
	// =============================================================================
	
	logger.Trace().Msg("Very detailed trace information")
	logger.Debug().Msg("Debug information for development")
	logger.Info().Msg("General information")
	logger.Warn().Msg("Warning - something might be wrong")
	logger.Error().Msg("Error - something went wrong")
	
	// Using WithLevel for dynamic levels
	logger.WithLevel(zerolog.InfoLevel).Msg("Dynamic level logging")

	// =============================================================================
	// COMPLEX STRUCTURED DATA
	// =============================================================================
	
	// Nested objects
	logger.Info().
		Dict("user", zerolog.Dict().
			Str("name", "John Doe").
			Int("age", 30).
			Str("email", "john@example.com")).
		Dict("request", zerolog.Dict().
			Str("method", "POST").
			Str("url", "/api/users").
			Int("status", 201)).
		Msg("User created successfully")

	// Arrays
	logger.Info().
		Strs("tags", []string{"api", "user", "create"}).
		Ints("response_codes", []int{200, 201, 204}).
		Msg("Operation completed")

	// =============================================================================
	// PERFORMANCE LOGGING
	// =============================================================================
	
	start := time.Now()
	
	// Simulate some work
	time.Sleep(10 * time.Millisecond)
	
	logger.Info().
		Str("operation", "data_processing").
		Dur("duration", time.Since(start)).
		Int("records_processed", 1000).
		Msg("Batch processing completed")

	// =============================================================================
	// ERROR HANDLING WITH CONTEXT
	// =============================================================================
	
	func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error().
					Interface("panic", r).
					Str("function", "example_function").
					Msg("Panic recovered")
			}
		}()
		
		// Simulate error scenarios
		if err := processData(); err != nil {
			logger.Error().
				Err(err).
				Str("component", "data_processor").
				Int("retry_count", 3).
				Msg("Failed to process data after retries")
		}
	}()

	// =============================================================================
	// CONDITIONAL LOGGING
	// =============================================================================
	
	debugMode := cfg.LogLevel == "debug"
	if debugMode {
		logger.Debug().
			Bool("debug_mode", debugMode).
			Str("config_file", ".env").
			Msg("Debug mode enabled")
	}

	// =============================================================================
	// HIGH-FREQUENCY LOGGING EXAMPLE
	// =============================================================================

	highFreqLogger := logger.WithComponent("high_frequency")

	// Example of logging in a loop (be careful with performance)
	for i := 0; i < 5; i++ {
		highFreqLogger.Debug().
			Int("iteration", i).
			Msg("High frequency operation")
	}

	// =============================================================================
	// FINAL MESSAGE
	// =============================================================================
	
	logger.Info().
		Str("example", "completed").
		Dur("total_runtime", time.Since(time.Now().Add(-100*time.Millisecond))).
		Msg("Logger example completed successfully")
}

// processData simulates a function that might return an error
func processData() error {
	// Simulate processing
	time.Sleep(5 * time.Millisecond)
	
	// Return an error for demonstration
	return errors.New("simulated processing error")
}

// demonstrateContextLogging shows how to use context-aware logging
func demonstrateContextLogging(ctx context.Context) {
	// Extract request ID from context (if available)
	if requestID, ok := ctx.Value("request_id").(string); ok {
		logger.WithRequestID(requestID).Info().
			Msg("Processing request with context")
	}
}
