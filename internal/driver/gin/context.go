package gin

import (
	basecontext "context"
	"net/http"

	corehttp "github.com/art-es/yet-another-service/internal/core/http"
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

func (c *context) With(ctx basecontext.Context) corehttp.Context {
	out := c.Context.Copy()
	out.Request = out.Request.WithContext(ctx)
	return newContext(out)
}
