package waclient

import (
	"fmt"
	"strings"

	waLog "go.mau.fi/whatsmeow/util/log"
	"zpwoot/platform/logger"
)


type WhatsmeowLogger struct {
	zpLogger *logger.Logger
	module   string
}


func NewWhatsmeowLogger(zpLogger *logger.Logger) waLog.Logger {
	return &WhatsmeowLogger{
		zpLogger: zpLogger.WithModule("whatsmeow"),
		module:   "whatsmeow",
	}
}


func (wl *WhatsmeowLogger) Debugf(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	wl.zpLogger.DebugWithFields("WhatsApp Debug", map[string]interface{}{
		"module":  wl.module,
		"message": message,
	})
}


func (wl *WhatsmeowLogger) Infof(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	wl.zpLogger.InfoWithFields("WhatsApp Info", map[string]interface{}{
		"module":  wl.module,
		"message": message,
	})
}


func (wl *WhatsmeowLogger) Warnf(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	wl.zpLogger.WarnWithFields("WhatsApp Warning", map[string]interface{}{
		"module":  wl.module,
		"message": message,
	})
}


func (wl *WhatsmeowLogger) Errorf(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	wl.zpLogger.ErrorWithFields("WhatsApp Error", map[string]interface{}{
		"module":  wl.module,
		"message": message,
	})
}


func (wl *WhatsmeowLogger) Sub(module string) waLog.Logger {
	subModule := fmt.Sprintf("%s.%s", wl.module, module)
	return &WhatsmeowLogger{
		zpLogger: wl.zpLogger.WithModule(subModule),
		module:   subModule,
	}
}


func (wl *WhatsmeowLogger) GetLevel() string {

	return "INFO"
}


func (wl *WhatsmeowLogger) SetLevel(level string) {


	wl.zpLogger.DebugWithFields("WhatsApp log level change requested", map[string]interface{}{
		"requested_level": level,
		"note":           "zpwoot manages log levels globally",
	})
}


type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)


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


func CreateWhatsmeowLoggerWithLevel(zpLogger *logger.Logger, minLevel LogLevel) waLog.Logger {
	return &FilteredWhatsmeowLogger{
		WhatsmeowLogger: &WhatsmeowLogger{
			zpLogger: zpLogger.WithModule("whatsmeow"),
			module:   "whatsmeow",
		},
		minLevel: minLevel,
	}
}


type FilteredWhatsmeowLogger struct {
	*WhatsmeowLogger
	minLevel LogLevel
}


func (fwl *FilteredWhatsmeowLogger) shouldLog(level LogLevel) bool {
	return level >= fwl.minLevel
}


func (fwl *FilteredWhatsmeowLogger) Debugf(msg string, args ...interface{}) {
	if fwl.shouldLog(LogLevelDebug) {
		fwl.WhatsmeowLogger.Debugf(msg, args...)
	}
}


func (fwl *FilteredWhatsmeowLogger) Infof(msg string, args ...interface{}) {
	if fwl.shouldLog(LogLevelInfo) {
		fwl.WhatsmeowLogger.Infof(msg, args...)
	}
}


func (fwl *FilteredWhatsmeowLogger) Warnf(msg string, args ...interface{}) {
	if fwl.shouldLog(LogLevelWarn) {
		fwl.WhatsmeowLogger.Warnf(msg, args...)
	}
}


func (fwl *FilteredWhatsmeowLogger) Errorf(msg string, args ...interface{}) {
	if fwl.shouldLog(LogLevelError) {
		fwl.WhatsmeowLogger.Errorf(msg, args...)
	}
}


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


type LoggerConfig struct {
	Level      string
	Module     string
	EnableSub  bool
	FilterSQL  bool
	FilterHTTP bool
}


func NewWhatsmeowLoggerWithConfig(zpLogger *logger.Logger, config LoggerConfig) waLog.Logger {
	if config.Module == "" {
		config.Module = "whatsmeow"
	}

	baseLogger := &WhatsmeowLogger{
		zpLogger: zpLogger.WithModule(config.Module),
		module:   config.Module,
	}


	if config.FilterSQL || config.FilterHTTP {
		return &FilteringWhatsmeowLogger{
			WhatsmeowLogger: baseLogger,
			filterSQL:       config.FilterSQL,
			filterHTTP:      config.FilterHTTP,
		}
	}

	return baseLogger
}


type FilteringWhatsmeowLogger struct {
	*WhatsmeowLogger
	filterSQL  bool
	filterHTTP bool
}


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


func (fwl *FilteringWhatsmeowLogger) Debugf(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	if !fwl.shouldFilter(message) {
		fwl.WhatsmeowLogger.Debugf(msg, args...)
	}
}


func (fwl *FilteringWhatsmeowLogger) Infof(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	if !fwl.shouldFilter(message) {
		fwl.WhatsmeowLogger.Infof(msg, args...)
	}
}


func (fwl *FilteringWhatsmeowLogger) Warnf(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	if !fwl.shouldFilter(message) {
		fwl.WhatsmeowLogger.Warnf(msg, args...)
	}
}


func (fwl *FilteringWhatsmeowLogger) Errorf(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	if !fwl.shouldFilter(message) {
		fwl.WhatsmeowLogger.Errorf(msg, args...)
	}
}


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
