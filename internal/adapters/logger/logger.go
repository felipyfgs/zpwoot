package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"zpwoot/internal/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var globalLogger zerolog.Logger


type LogFormat string

const (
	FormatJSON    LogFormat = "json"
	FormatConsole LogFormat = "console"
)


type LogOutput string

const (
	OutputStdout LogOutput = "stdout"
	OutputStderr LogOutput = "stderr"
)


func shortCaller(file string, line int) string {
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}


func extractPackageFromFile(file string) string {

	if strings.Contains(file, "/workspaces/zpwoot/") {
		parts := strings.Split(file, "/workspaces/zpwoot/")
		if len(parts) > 1 {
			path := parts[1]

			dir := filepath.Dir(path)
			if dir == "." {
				return "main"
			}




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


type packageHook struct{}

func (h packageHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {

	for depth := 3; depth <= 6; depth++ {
		_, file, _, ok := runtime.Caller(depth)
		if ok {

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


	e.Str("pkg", "main")
}


func Init(level string) {

	logLevel := parseLogLevel(level)


	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return shortCaller(file, line)
	}


	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}


	globalLogger = zerolog.New(consoleWriter).
		Level(logLevel).
		Hook(packageHook{}).
		With().
		Timestamp().
		Caller().
		Logger()


	log.Logger = globalLogger
}


func InitWithConfig(cfg *config.Config) {

	logLevel := parseLogLevel(cfg.LogLevel)


	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return shortCaller(file, line)
	}


	var output *os.File
	switch LogOutput(strings.ToLower(cfg.LogOutput)) {
	case OutputStderr:
		output = os.Stderr
	case OutputStdout:
		output = os.Stdout
	default:
		output = os.Stderr
	}


	var logger zerolog.Logger
	switch LogFormat(strings.ToLower(cfg.LogFormat)) {
	case FormatJSON:

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


	globalLogger = logger
	log.Logger = logger
}


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


type Logger struct {
	logger zerolog.Logger
}


func New() *Logger {
	return &Logger{
		logger: globalLogger,
	}
}


func NewFromAppConfig(cfg *config.Config) *Logger {
	return &Logger{
		logger: globalLogger,
	}
}


func GetGlobalLogger() *Logger {
	return &Logger{
		logger: globalLogger,
	}
}


func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		logger: l.logger.With().Ctx(ctx).Logger(),
	}
}


func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		logger: l.logger.With().Interface(key, value).Logger(),
	}
}


func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	logContext := l.logger.With()
	for key, value := range fields {
		logContext = logContext.Interface(key, value)
	}
	return &Logger{
		logger: logContext.Logger(),
	}
}


func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		logger: l.logger.With().Err(err).Logger(),
	}
}


func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{
		logger: l.logger.With().Str("component", component).Logger(),
	}
}


func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{
		logger: l.logger.With().Str("request_id", requestID).Logger(),
	}
}


func (l *Logger) WithSessionID(sessionID string) *Logger {
	return &Logger{
		logger: l.logger.With().Str("session_id", sessionID).Logger(),
	}
}




func (l *Logger) Trace() *zerolog.Event {
	return l.logger.Trace()
}


func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}


func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}


func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}


func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}


func (l *Logger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}


func (l *Logger) Panic() *zerolog.Event {
	return l.logger.Panic()
}


func (l *Logger) WithLevel(level zerolog.Level) *zerolog.Event {
	return l.logger.WithLevel(level)
}




func (l *Logger) TraceMsg(msg string) {
	l.logger.Trace().Msg(msg)
}


func (l *Logger) DebugMsg(msg string) {
	l.logger.Debug().Msg(msg)
}


func (l *Logger) InfoMsg(msg string) {
	l.logger.Info().Msg(msg)
}


func (l *Logger) WarnMsg(msg string) {
	l.logger.Warn().Msg(msg)
}


func (l *Logger) ErrorMsg(msg string) {
	l.logger.Error().Msg(msg)
}


func (l *Logger) FatalMsg(msg string) {
	l.logger.Fatal().Msg(msg)
}


func (l *Logger) PanicMsg(msg string) {
	l.logger.Panic().Msg(msg)
}




func Trace() *zerolog.Event {
	return globalLogger.Trace()
}


func Debug() *zerolog.Event {
	return globalLogger.Debug()
}


func Info() *zerolog.Event {
	return globalLogger.Info()
}


func Warn() *zerolog.Event {
	return globalLogger.Warn()
}


func Error() *zerolog.Event {
	return globalLogger.Error()
}


func Fatal() *zerolog.Event {
	return globalLogger.Fatal()
}


func Panic() *zerolog.Event {
	return globalLogger.Panic()
}


func WithLevel(level zerolog.Level) *zerolog.Event {
	return globalLogger.WithLevel(level)
}




func WithFields(fields map[string]interface{}) *Logger {
	logContext := globalLogger.With()
	for key, value := range fields {
		logContext = logContext.Interface(key, value)
	}
	return &Logger{
		logger: logContext.Logger(),
	}
}


func WithError(err error) *Logger {
	return &Logger{
		logger: globalLogger.With().Err(err).Logger(),
	}
}


func WithComponent(component string) *Logger {
	return &Logger{
		logger: globalLogger.With().Str("component", component).Logger(),
	}
}


func WithRequestID(requestID string) *Logger {
	return &Logger{
		logger: globalLogger.With().Str("request_id", requestID).Logger(),
	}
}


func WithSessionID(sessionID string) *Logger {
	return &Logger{
		logger: globalLogger.With().Str("session_id", sessionID).Logger(),
	}
}


func GetZerologLogger() zerolog.Logger {
	return globalLogger
}
