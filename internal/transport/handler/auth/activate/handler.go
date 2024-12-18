//go:generate mockgen -source=handler.go -destination=mock/handler.go -package=mock
package activate

import (
	"context"
	"errors"
	nethttp "net/http"

	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
	"github.com/art-es/yet-another-service/internal/core/http"
	httputil "github.com/art-es/yet-another-service/internal/core/http/util"
	"github.com/art-es/yet-another-service/internal/core/log"
	"github.com/art-es/yet-another-service/internal/core/validation"
)

type authService interface {
	Activate(ctx context.Context, token string) error
}

type Handler struct {
	authService authService
	logger      log.Logger
	validator   validation.Validator
}

func NewHandler(
	authService authService,
	logger log.Logger,
	validator validation.Validator,
) *Handler {
	return &Handler{
		authService: authService,
		logger:      logger,
		validator:   validator,
	}
}

func (h *Handler) Handle(ctx http.Context) {
	token, err := h.parseToken(ctx)
	if err != nil {
		httputil.RespondNotFound(ctx)
		return
	}

	err = h.authService.Activate(ctx, token)

	switch {
	case err == nil:
		httputil.Respond(ctx, nethttp.StatusOK, struct{}{})
	case errors.Is(err, apperrors.ErrUserActivationNotFound):
		httputil.RespondNotFound(ctx)
	default:
		h.logger.Error().Err(err).Msg("activate error on auth service")
		httputil.RespondInternalError(ctx)
	}
}

func (h *Handler) parseToken(ctx http.Context) (string, error) {
	token := ctx.Request().URL.Query().Get("token")

	if err := h.validator.Var(token, "required,uuid"); err != nil {
		return "", err
	}

	return token, nil
}
