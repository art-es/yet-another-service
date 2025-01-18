//go:generate mockgen -source=handler.go -destination=mock/handler.go -package=mock
package articles_get

import (
	"context"
	"net/http"

	"github.com/art-es/yet-another-service/internal/app/shared/dto"
	corehttp "github.com/art-es/yet-another-service/internal/core/http"
	corehttputil "github.com/art-es/yet-another-service/internal/core/http/util"
	"github.com/art-es/yet-another-service/internal/core/log"
)

type articleService interface {
	Get(ctx context.Context, in *dto.GetArticlesIn) (*dto.GetArticlesOut, error)
}

type Handler struct {
	articleService articleService
	logger         log.Logger
}

func NewHandler(
	articleService articleService,
	logger log.Logger,
) *Handler {
	return &Handler{
		articleService: articleService,
		logger:         logger,
	}
}

func (h *Handler) Handle(ctx corehttp.Context) {
	req := parseRequest(ctx.Request())

	out, err := h.articleService.Get(ctx, &dto.GetArticlesIn{
		FromSlug: req.FromSlug,
	})
	if err != nil {
		h.logger.Error().Err(err).Msg("get error on article service")
		corehttputil.RespondInternalError(ctx)
		return
	}

	corehttputil.Respond(ctx, http.StatusOK, convertResponse(out))
}
