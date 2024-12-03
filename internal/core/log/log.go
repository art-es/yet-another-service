//go:generate mockgen -source=log.go -destination=mock/log.go -package=mock
package log

type Logger interface {
	Info() Event
	Warn() Event
	Error() Event
	Panic() Event
	With() Context
}

type Event interface {
	Err(err error) Event
	Str(key, val string) Event
	Msg(msg string)
	Send()
}

type Context interface {
	Err(err error) Context
	Str(key, val string) Context
	Logger() Logger
}
