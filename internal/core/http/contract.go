//go:generate mockgen -source=contract.go -destination=mock/contract.go -package=mock
package http

import (
	"context"
	"net/http"
)

type Context interface {
	context.Context
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
}
