package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type context struct {
	*gin.Context
}

func newContext(ctx *gin.Context) *context {
	return &context{
		Context: ctx,
	}
}

func (c *context) Request() *http.Request {
	return c.Context.Request
}

func (c *context) ResponseWriter() http.ResponseWriter {
	return c.Context.Writer
}
