package zerolog

import (
	"io"
	"os"

	"github.com/rs/zerolog"

	"github.com/art-es/yet-another-service/internal/core/log"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	logger zerolog.Logger
}

func NewLogger() *Logger {
	return NewLoggerWithWriter(os.Stdout)
}

func NewLoggerWithWriter(writer io.Writer) *Logger {
	return newLogger(zerolog.New(writer))
}

func newLogger(l zerolog.Logger) *Logger {
	return &Logger{logger: l}
}

func (l *Logger) Info() log.Event {
	return newEvent(l.logger.Info())
}

func (l *Logger) Warn() log.Event {
	return newEvent(l.logger.Warn())
}

func (l *Logger) Error() log.Event {
	return newEvent(l.logger.Error())
}

func (l *Logger) Panic() log.Event {
	return newEvent(l.logger.Panic())
}

func (l *Logger) With() log.Context {
	return newContext(l.logger.With())
}
