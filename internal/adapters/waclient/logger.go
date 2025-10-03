package waclient

import (
	"fmt"
	"strings"

	waLog "go.mau.fi/whatsmeow/util/log"
	"zpwoot/platform/logger"
)

// WhatsmeowLogger implements whatsmeow's Logger interface using zpwoot's logger
type WhatsmeowLogger struct {
	zpLogger *logger.Logger
	module   string
}

// NewWhatsmeowLogger creates a new whatsmeow logger that integrates with zpwoot's logging system
func NewWhatsmeowLogger(zpLogger *logger.Logger) waLog.Logger {
	return &WhatsmeowLogger{
		zpLogger: zpLogger.WithModule("whatsmeow"),
		module:   "whatsmeow",
	}
}

// Debugf logs a debug message with formatting
func (wl *WhatsmeowLogger) Debugf(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	wl.zpLogger.DebugWithFields("WhatsApp Debug", map[string]interface{}{
		"module":  wl.module,
		"message": message,
	})
}

// Infof logs an info message with formatting
func (wl *WhatsmeowLogger) Infof(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	wl.zpLogger.InfoWithFields("WhatsApp Info", map[string]interface{}{
		"module":  wl.module,
		"message": message,
	})
}

// Warnf logs a warning message with formatting
func (wl *WhatsmeowLogger) Warnf(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	wl.zpLogger.WarnWithFields("WhatsApp Warning", map[string]interface{}{
		"module":  wl.module,
		"message": message,
	})
}

// Errorf logs an error message with formatting
func (wl *WhatsmeowLogger) Errorf(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	wl.zpLogger.ErrorWithFields("WhatsApp Error", map[string]interface{}{
		"module":  wl.module,
		"message": message,
	})
}

// Sub creates a sub-logger with additional context
func (wl *WhatsmeowLogger) Sub(module string) waLog.Logger {
	subModule := fmt.Sprintf("%s.%s", wl.module, module)
	return &WhatsmeowLogger{
		zpLogger: wl.zpLogger.WithModule(subModule),
		module:   subModule,
	}
}

// GetLevel returns the current log level (placeholder implementation)
func (wl *WhatsmeowLogger) GetLevel() string {
	// Return a simple string level since whatsmeow doesn't define Level type
	return "INFO"
}

// SetLevel sets the log level (placeholder implementation)
func (wl *WhatsmeowLogger) SetLevel(level string) {
	// This is a no-op since zpwoot manages log levels globally
	// We could potentially update the zpwoot logger level here if needed
	wl.zpLogger.DebugWithFields("WhatsApp log level change requested", map[string]interface{}{
		"requested_level": level,
		"note":           "zpwoot manages log levels globally",
	})
}

// LogLevel represents whatsmeow log levels for internal use
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// String returns string representation of log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// CreateWhatsmeowLoggerWithLevel creates a whatsmeow logger with specific level filtering
func CreateWhatsmeowLoggerWithLevel(zpLogger *logger.Logger, minLevel LogLevel) waLog.Logger {
	return &FilteredWhatsmeowLogger{
		WhatsmeowLogger: &WhatsmeowLogger{
			zpLogger: zpLogger.WithModule("whatsmeow"),
			module:   "whatsmeow",
		},
		minLevel: minLevel,
	}
}

// FilteredWhatsmeowLogger wraps WhatsmeowLogger with level filtering
type FilteredWhatsmeowLogger struct {
	*WhatsmeowLogger
	minLevel LogLevel
}

// shouldLog checks if a message should be logged based on level
func (fwl *FilteredWhatsmeowLogger) shouldLog(level LogLevel) bool {
	return level >= fwl.minLevel
}

// Debugf logs a debug message with level filtering
func (fwl *FilteredWhatsmeowLogger) Debugf(msg string, args ...interface{}) {
	if fwl.shouldLog(LogLevelDebug) {
		fwl.WhatsmeowLogger.Debugf(msg, args...)
	}
}

// Infof logs an info message with level filtering
func (fwl *FilteredWhatsmeowLogger) Infof(msg string, args ...interface{}) {
	if fwl.shouldLog(LogLevelInfo) {
		fwl.WhatsmeowLogger.Infof(msg, args...)
	}
}

// Warnf logs a warning message with level filtering
func (fwl *FilteredWhatsmeowLogger) Warnf(msg string, args ...interface{}) {
	if fwl.shouldLog(LogLevelWarn) {
		fwl.WhatsmeowLogger.Warnf(msg, args...)
	}
}

// Errorf logs an error message with level filtering
func (fwl *FilteredWhatsmeowLogger) Errorf(msg string, args ...interface{}) {
	if fwl.shouldLog(LogLevelError) {
		fwl.WhatsmeowLogger.Errorf(msg, args...)
	}
}

// Sub creates a filtered sub-logger
func (fwl *FilteredWhatsmeowLogger) Sub(module string) waLog.Logger {
	subModule := fmt.Sprintf("%s.%s", fwl.module, module)
	return &FilteredWhatsmeowLogger{
		WhatsmeowLogger: &WhatsmeowLogger{
			zpLogger: fwl.zpLogger.WithModule(subModule),
			module:   subModule,
		},
		minLevel: fwl.minLevel,
	}
}

// LoggerConfig holds configuration for whatsmeow logger
type LoggerConfig struct {
	Level      string
	Module     string
	EnableSub  bool
	FilterSQL  bool // Filter out SQL-related logs
	FilterHTTP bool // Filter out HTTP-related logs
}

// NewWhatsmeowLoggerWithConfig creates a whatsmeow logger with custom configuration
func NewWhatsmeowLoggerWithConfig(zpLogger *logger.Logger, config LoggerConfig) waLog.Logger {
	if config.Module == "" {
		config.Module = "whatsmeow"
	}

	baseLogger := &WhatsmeowLogger{
		zpLogger: zpLogger.WithModule(config.Module),
		module:   config.Module,
	}

	// Apply filtering if needed
	if config.FilterSQL || config.FilterHTTP {
		return &FilteringWhatsmeowLogger{
			WhatsmeowLogger: baseLogger,
			filterSQL:       config.FilterSQL,
			filterHTTP:      config.FilterHTTP,
		}
	}

	return baseLogger
}

// FilteringWhatsmeowLogger filters specific types of messages
type FilteringWhatsmeowLogger struct {
	*WhatsmeowLogger
	filterSQL  bool
	filterHTTP bool
}

// shouldFilter checks if a message should be filtered out
func (fwl *FilteringWhatsmeowLogger) shouldFilter(msg string) bool {
	msgLower := strings.ToLower(msg)
	
	if fwl.filterSQL && (strings.Contains(msgLower, "sql") || strings.Contains(msgLower, "database")) {
		return true
	}
	
	if fwl.filterHTTP && (strings.Contains(msgLower, "http") || strings.Contains(msgLower, "request")) {
		return true
	}
	
	return false
}

// Debugf logs a debug message with content filtering
func (fwl *FilteringWhatsmeowLogger) Debugf(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	if !fwl.shouldFilter(message) {
		fwl.WhatsmeowLogger.Debugf(msg, args...)
	}
}

// Infof logs an info message with content filtering
func (fwl *FilteringWhatsmeowLogger) Infof(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	if !fwl.shouldFilter(message) {
		fwl.WhatsmeowLogger.Infof(msg, args...)
	}
}

// Warnf logs a warning message with content filtering
func (fwl *FilteringWhatsmeowLogger) Warnf(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	if !fwl.shouldFilter(message) {
		fwl.WhatsmeowLogger.Warnf(msg, args...)
	}
}

// Errorf logs an error message with content filtering
func (fwl *FilteringWhatsmeowLogger) Errorf(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	if !fwl.shouldFilter(message) {
		fwl.WhatsmeowLogger.Errorf(msg, args...)
	}
}

// Sub creates a filtering sub-logger
func (fwl *FilteringWhatsmeowLogger) Sub(module string) waLog.Logger {
	subModule := fmt.Sprintf("%s.%s", fwl.module, module)
	return &FilteringWhatsmeowLogger{
		WhatsmeowLogger: &WhatsmeowLogger{
			zpLogger: fwl.zpLogger.WithModule(subModule),
			module:   subModule,
		},
		filterSQL:  fwl.filterSQL,
		filterHTTP: fwl.filterHTTP,
	}
}
