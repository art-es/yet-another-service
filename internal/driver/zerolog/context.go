package zerolog

import (
	"github.com/art-es/yet-another-service/internal/core/log"
	"github.com/rs/zerolog"
)

var _ log.Context = (*context)(nil)

type context struct {
	context zerolog.Context
}

func newContext(c zerolog.Context) *context {
	return &context{context: c}
}

func (c *context) Err(err error) log.Context {
	c.context = c.context.Err(err)
	return c
}

func (c *context) Str(key, val string) log.Context {
	c.context = c.context.Str(key, val)
	return c
}

func (c *context) Logger() log.Logger {
	return newLogger(c.context.Logger())
}
