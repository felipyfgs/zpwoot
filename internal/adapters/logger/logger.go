package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"zpwoot/internal/adapters/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var globalLogger zerolog.Logger

// LogFormat represents the log output format
type LogFormat string

const (
	FormatJSON    LogFormat = "json"
	FormatConsole LogFormat = "console"
)

// LogOutput represents the log output destination
type LogOutput string

const (
	OutputStdout LogOutput = "stdout"
	OutputStderr LogOutput = "stderr"
)

// shortCaller extracts just the filename and line number from the full path
func shortCaller(file string, line int) string {
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// extractPackageFromFile extracts the package/module name from file path
func extractPackageFromFile(file string) string {
	// Remove the workspace prefix and get the directory
	if strings.Contains(file, "/workspaces/zpwoot/") {
		parts := strings.Split(file, "/workspaces/zpwoot/")
		if len(parts) > 1 {
			path := parts[1]
			// Get the directory part
			dir := filepath.Dir(path)
			if dir == "." {
				return "main"
			}
			// Convert path to package name
			// e.g., "internal/adapters/database" -> "database"
			// e.g., "internal/adapters/http/handlers" -> "handlers"
			// e.g., "cmd/zpwoot" -> "main"
			pathParts := strings.Split(dir, "/")
			if len(pathParts) > 0 {
				lastPart := pathParts[len(pathParts)-1]
				if lastPart == "zpwoot" || strings.HasPrefix(dir, "cmd/") {
					return "main"
				}
				return lastPart
			}
		}
	}

	// Fallback: try to extract from any path structure
	dir := filepath.Dir(file)
	if dir == "." {
		return "main"
	}

	pathParts := strings.Split(dir, "/")
	if len(pathParts) > 0 {
		return pathParts[len(pathParts)-1]
	}

	return "unknown"
}

// packageHook adds package information to log events
type packageHook struct{}

func (h packageHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	// Try different caller depths to find the right frame
	for depth := 3; depth <= 6; depth++ {
		_, file, _, ok := runtime.Caller(depth)
		if ok {
			// Skip internal logger files
			if strings.Contains(file, "/logger/") ||
				strings.Contains(file, "/zerolog/") ||
				strings.Contains(file, "/runtime/") {
				continue
			}

			pkg := extractPackageFromFile(file)
			if pkg != "unknown" && pkg != "runtime" {
				e.Str("pkg", pkg)
				return
			}
		}
	}

	// Fallback
	e.Str("pkg", "main")
}

// Init initializes the global logger with basic configuration
func Init(level string) {
	// Set log level
	logLevel := parseLogLevel(level)

	// Configure caller to show only filename without prefix
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return shortCaller(file, line)
	}

	// Create console writer for pretty output (default)
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}

	// Initialize global logger with package detection
	globalLogger = zerolog.New(consoleWriter).
		Level(logLevel).
		Hook(packageHook{}).
		With().
		Timestamp().
		Caller().
		Logger()

	// Set as global zerolog logger
	log.Logger = globalLogger
}

// InitWithConfig initializes the logger with full configuration following the exact pattern
func InitWithConfig(cfg *config.Config) {
	// Set log level
	logLevel := parseLogLevel(cfg.LogLevel)

	// Configure caller to show only filename
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return shortCaller(file, line)
	}

	// Determine output destination
	var output *os.File
	switch LogOutput(strings.ToLower(cfg.LogOutput)) {
	case OutputStderr:
		output = os.Stderr
	case OutputStdout:
		output = os.Stdout
	default:
		output = os.Stderr // Default to stderr like the example
	}

	// Configure logger based on format
	var logger zerolog.Logger
	switch LogFormat(strings.ToLower(cfg.LogFormat)) {
	case FormatJSON:
		// JSON structured logging with package detection
		logger = zerolog.New(output).
			Level(logLevel).
			Hook(packageHook{}).
			With().
			Timestamp().
			Caller().
			Str("service", "zpwoot").
			Str("version", "1.0.0").
			Logger()
	case FormatConsole:
		// Console logging with package detection
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: time.RFC3339,
		}).
			Level(logLevel).
			Hook(packageHook{}).
			With().
			Timestamp().
			Caller().
			Logger()
	default:
		// Default to console format with package detection
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: time.RFC3339,
		}).
			Level(logLevel).
			Hook(packageHook{}).
			With().
			Timestamp().
			Caller().
			Logger()
	}

	// Set global logger
	globalLogger = logger
	log.Logger = logger
}

// parseLogLevel parses string log level to zerolog level
func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "disabled":
		return zerolog.Disabled
	default:
		return zerolog.InfoLevel
	}
}

// Logger wraps zerolog.Logger following the established pattern
type Logger struct {
	logger zerolog.Logger
}

// New creates a new logger instance using the global logger
func New() *Logger {
	return &Logger{
		logger: globalLogger,
	}
}

// NewFromAppConfig creates a new logger from app config
func NewFromAppConfig(cfg *config.Config) *Logger {
	return &Logger{
		logger: globalLogger,
	}
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() *Logger {
	return &Logger{
		logger: globalLogger,
	}
}

// WithContext creates a logger with context fields
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		logger: l.logger.With().Ctx(ctx).Logger(),
	}
}

// WithField adds a field to the logger context
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		logger: l.logger.With().Interface(key, value).Logger(),
	}
}

// WithFields adds multiple fields to the logger context
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	logContext := l.logger.With()
	for key, value := range fields {
		logContext = logContext.Interface(key, value)
	}
	return &Logger{
		logger: logContext.Logger(),
	}
}

// WithError adds an error field to the logger context
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		logger: l.logger.With().Err(err).Logger(),
	}
}

// WithComponent adds a component field to the logger context
func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{
		logger: l.logger.With().Str("component", component).Logger(),
	}
}

// WithRequestID adds a request ID field to the logger context
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{
		logger: l.logger.With().Str("request_id", requestID).Logger(),
	}
}

// WithSessionID adds a session ID field to the logger context
func (l *Logger) WithSessionID(sessionID string) *Logger {
	return &Logger{
		logger: l.logger.With().Str("session_id", sessionID).Logger(),
	}
}

// Basic logging methods following the zerolog pattern

// Trace logs a trace message
func (l *Logger) Trace() *zerolog.Event {
	return l.logger.Trace()
}

// Debug logs a debug message
func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

// Info logs an info message
func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}

// Warn logs a warning message
func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}

// Error logs an error message
func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

// Panic logs a panic message and panics
func (l *Logger) Panic() *zerolog.Event {
	return l.logger.Panic()
}

// WithLevel logs with a specific level (like the example)
func (l *Logger) WithLevel(level zerolog.Level) *zerolog.Event {
	return l.logger.WithLevel(level)
}

// Convenience methods for simple string messages

// TraceMsg logs a simple trace message
func (l *Logger) TraceMsg(msg string) {
	l.logger.Trace().Msg(msg)
}

// DebugMsg logs a simple debug message
func (l *Logger) DebugMsg(msg string) {
	l.logger.Debug().Msg(msg)
}

// InfoMsg logs a simple info message
func (l *Logger) InfoMsg(msg string) {
	l.logger.Info().Msg(msg)
}

// WarnMsg logs a simple warning message
func (l *Logger) WarnMsg(msg string) {
	l.logger.Warn().Msg(msg)
}

// ErrorMsg logs a simple error message
func (l *Logger) ErrorMsg(msg string) {
	l.logger.Error().Msg(msg)
}

// FatalMsg logs a simple fatal message and exits
func (l *Logger) FatalMsg(msg string) {
	l.logger.Fatal().Msg(msg)
}

// PanicMsg logs a simple panic message and panics
func (l *Logger) PanicMsg(msg string) {
	l.logger.Panic().Msg(msg)
}

// Global logger functions following the zerolog pattern

// Trace returns a trace event from the global logger
func Trace() *zerolog.Event {
	return globalLogger.Trace()
}

// Debug returns a debug event from the global logger
func Debug() *zerolog.Event {
	return globalLogger.Debug()
}

// Info returns an info event from the global logger
func Info() *zerolog.Event {
	return globalLogger.Info()
}

// Warn returns a warn event from the global logger
func Warn() *zerolog.Event {
	return globalLogger.Warn()
}

// Error returns an error event from the global logger
func Error() *zerolog.Event {
	return globalLogger.Error()
}

// Fatal returns a fatal event from the global logger
func Fatal() *zerolog.Event {
	return globalLogger.Fatal()
}

// Panic returns a panic event from the global logger
func Panic() *zerolog.Event {
	return globalLogger.Panic()
}

// WithLevel returns an event with the specified level from the global logger
func WithLevel(level zerolog.Level) *zerolog.Event {
	return globalLogger.WithLevel(level)
}

// Convenience functions for creating contextual loggers from global logger

// WithFields creates a logger with structured fields using the global logger
func WithFields(fields map[string]interface{}) *Logger {
	logContext := globalLogger.With()
	for key, value := range fields {
		logContext = logContext.Interface(key, value)
	}
	return &Logger{
		logger: logContext.Logger(),
	}
}

// WithError creates a logger with error field using the global logger
func WithError(err error) *Logger {
	return &Logger{
		logger: globalLogger.With().Err(err).Logger(),
	}
}

// WithComponent creates a logger with component field using the global logger
func WithComponent(component string) *Logger {
	return &Logger{
		logger: globalLogger.With().Str("component", component).Logger(),
	}
}

// WithRequestID creates a logger with request ID field using the global logger
func WithRequestID(requestID string) *Logger {
	return &Logger{
		logger: globalLogger.With().Str("request_id", requestID).Logger(),
	}
}

// WithSessionID creates a logger with session ID field using the global logger
func WithSessionID(sessionID string) *Logger {
	return &Logger{
		logger: globalLogger.With().Str("session_id", sessionID).Logger(),
	}
}

// GetZerologLogger returns the underlying zerolog.Logger for direct access
func GetZerologLogger() zerolog.Logger {
	return globalLogger
}
