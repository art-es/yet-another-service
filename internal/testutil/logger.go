package testutil

import (
	"bufio"
	"bytes"

	"github.com/art-es/yet-another-service/internal/core/log"
	"github.com/art-es/yet-another-service/internal/driver/zerolog"
)

type Logger struct {
	log.Logger
	buf *bytes.Buffer
}

func NewLogger() *Logger {
	buf := &bytes.Buffer{}

	return &Logger{
		Logger: zerolog.NewLoggerWithWriter(buf),
		buf:    buf,
	}
}

func (l *Logger) Logs() []string {
	var logs []string
	for s := bufio.NewScanner(l.buf); s.Scan(); {
		logs = append(logs, s.Text())
	}
	return logs
}
