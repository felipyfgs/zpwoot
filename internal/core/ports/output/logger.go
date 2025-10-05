package output

import (
	"github.com/rs/zerolog"
)

type Logger interface {
	Trace() *zerolog.Event
	Debug() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
	Fatal() *zerolog.Event
	Panic() *zerolog.Event
	WithLevel(level zerolog.Level) *zerolog.Event

	TraceMsg(msg string)
	DebugMsg(msg string)
	InfoMsg(msg string)
	WarnMsg(msg string)
	ErrorMsg(msg string)
	FatalMsg(msg string)
	PanicMsg(msg string)
}
