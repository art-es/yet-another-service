package zerolog

import (
	"github.com/art-es/yet-another-service/internal/core/log"
	"github.com/rs/zerolog"
)

var _ log.Event = (*event)(nil)

type event struct {
	event *zerolog.Event
}

func newEvent(e *zerolog.Event) *event {
	return &event{event: e}
}

func (e *event) Err(err error) log.Event {
	e.event = e.event.Err(err)
	return e
}

func (e *event) Str(key, val string) log.Event {
	e.event = e.event.Str(key, val)
	return e
}

func (e *event) Msg(msg string) {
	e.event.Msg(msg)
}

func (e *event) Send() {
	e.event.Send()
}
