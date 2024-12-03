package gin

import (
	"github.com/gin-gonic/gin"

	"github.com/art-es/yet-another-service/internal/core/http"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter() *Router {
	return &Router{
		engine: gin.New(),
	}
}

func (r *Router) Register(method, path string, handle func(ctx http.Context)) {
	r.engine.Handle(method, path, func(ctx *gin.Context) {
		handle(newContext(ctx))
	})
}

func (r *Router) Run() error {
	return r.engine.Run(":8080")
}
