package logger

import (
	"log"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

var globalLogger zerolog.Logger

// Init initializes the global logger
func Init(level string) {
	// Set log level
	logLevel := parseLogLevel(level)
	zerolog.SetGlobalLevel(logLevel)

	// Create console writer for pretty output
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	
	// Initialize global logger
	globalLogger = zerolog.New(consoleWriter).With().Timestamp().Logger()

	// Set standard log to use zerolog
	log.SetOutput(globalLogger)
}

// parseLogLevel parses string log level to zerolog level
func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
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
	default:
		return zerolog.InfoLevel
	}
}

// Logger wraps zerolog.Logger with additional methods
type Logger struct {
	logger zerolog.Logger
}

// NewFromAppConfig creates a new logger from app config
func NewFromAppConfig(cfg interface{}) *Logger {
	return &Logger{
		logger: globalLogger,
	}
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Error().Msgf(format, args...)
}

// Fatalf logs a formatted fatal message and exits
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal().Msgf(format, args...)
}

// InfoWithFields logs an info message with fields
func (l *Logger) InfoWithFields(msg string, fields map[string]interface{}) {
	event := l.logger.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// WarnWithFields logs a warning message with fields
func (l *Logger) WarnWithFields(msg string, fields map[string]interface{}) {
	event := l.logger.Warn()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// ErrorWithFields logs an error message with fields
func (l *Logger) ErrorWithFields(msg string, fields map[string]interface{}) {
	event := l.logger.Error()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// Global logger functions for backward compatibility
func Info(msg string) {
	globalLogger.Info().Msg(msg)
}

func Debug(msg string) {
	globalLogger.Debug().Msg(msg)
}

func Warn(msg string) {
	globalLogger.Warn().Msg(msg)
}

func Error(msg string) {
	globalLogger.Error().Msg(msg)
}

func Fatal(msg string) {
	globalLogger.Fatal().Msg(msg)
}

func Infof(format string, args ...interface{}) {
	globalLogger.Info().Msgf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	globalLogger.Debug().Msgf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	globalLogger.Warn().Msgf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	globalLogger.Error().Msgf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatal().Msgf(format, args...)
}
