//go:generate mockgen -source=handler.go -destination=mock/handler.go -package=mock
package forgot_password

import (
	"context"
	"errors"
	nethttp "net/http"

	"github.com/art-es/yet-another-service/internal/core/http"
	"github.com/art-es/yet-another-service/internal/core/log"
	"github.com/art-es/yet-another-service/internal/core/validation"
	errorsd "github.com/art-es/yet-another-service/internal/domain/shared/errors"
)

type authService interface {
	CreateRecovery(ctx context.Context, email string) error
}

type request struct {
	Email string `json:"email" validate:"required,email,lte=255"`
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
	req, err := h.parseRequest(ctx)
	if err != nil {
		http.RespondBadRequest(ctx, err.Error())
		return
	}

	err = h.authService.CreateRecovery(ctx, req.Email)

	switch {
	case err == nil:
		http.Respond(ctx, nethttp.StatusOK, struct{}{})
	case errors.Is(err, errorsd.ErrUserNotFound):
		http.RespondBadRequest(ctx, "user with this email not found")
	default:
		h.logger.Error().Err(err).Msg("create recovery error on auth service")
		http.RespondInternalError(ctx)
	}
}

func (h *Handler) parseRequest(ctx http.Context) (*request, error) {
	req := &request{}

	if err := http.EnrichRequestBody(ctx, req); err != nil {
		return nil, err
	}

	if err := h.validator.Struct(req); err != nil {
		return nil, err
	}

	return req, nil
}
