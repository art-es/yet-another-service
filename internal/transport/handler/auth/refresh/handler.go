//go:generate mockgen -source=handler.go -destination=mock/handler.go -package=mock
package refresh

import (
	"errors"
	nethttp "net/http"

	"github.com/art-es/yet-another-service/internal/core/http"
	"github.com/art-es/yet-another-service/internal/core/log"
	errorsd "github.com/art-es/yet-another-service/internal/domain/shared/errors"
)

type authService interface {
	Refresh(refreshToken string) (string, error)
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
	refreshToken, ok := http.GetAuthorizationToken(ctx)
	if !ok {
		http.RespondUnauthorized(ctx)
		return
	}

	accessToken, err := h.authService.Refresh(refreshToken)

	switch {
	case err == nil:
		http.Respond(ctx, nethttp.StatusOK, response{AccessToken: accessToken})
	case errors.Is(err, errorsd.ErrInvalidAuthToken):
		http.RespondUnauthorized(ctx)
	default:
		h.logger.Error().Err(err).Msg("refresh error on auth service")
		http.RespondInternalError(ctx)
	}
}
