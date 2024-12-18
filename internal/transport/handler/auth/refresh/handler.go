//go:generate mockgen -source=handler.go -destination=mock/handler.go -package=mock
package refresh

import (
	"context"
	"errors"
	nethttp "net/http"

	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
	"github.com/art-es/yet-another-service/internal/core/http"
	"github.com/art-es/yet-another-service/internal/core/http/util"
	"github.com/art-es/yet-another-service/internal/core/log"
)

type authService interface {
	Refresh(ctx context.Context, refreshToken string) (string, error)
}

type response struct {
	AccessToken string `json:"accessToken"`
}

type Handler struct {
	authService authService
	logger      log.Logger
}

func NewHandler(
	authService authService,
	logger log.Logger,
) *Handler {
	return &Handler{
		authService: authService,
		logger:      logger,
	}
}

func (h *Handler) Handle(ctx http.Context) {
	refreshToken, ok := util.GetAuthorizationToken(ctx)
	if !ok {
		util.RespondUnauthorized(ctx)
		return
	}

	accessToken, err := h.authService.Refresh(ctx, refreshToken)

	switch {
	case err == nil:
		util.Respond(ctx, nethttp.StatusOK, response{AccessToken: accessToken})
	case errors.Is(err, apperrors.ErrInvalidAuthToken):
		util.RespondUnauthorized(ctx)
	default:
		h.logger.Error().Err(err).Msg("refresh error on auth service")
		util.RespondInternalError(ctx)
	}
}
