//go:generate mockgen -source=handler.go -destination=mock/handler.go -package=mock
package forgotpassword

import (
	"context"
	"errors"
	nethttp "net/http"

	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
	"github.com/art-es/yet-another-service/internal/core/http"
	"github.com/art-es/yet-another-service/internal/core/http/util"
	"github.com/art-es/yet-another-service/internal/core/log"
	"github.com/art-es/yet-another-service/internal/core/validation"
)

type recoveryService interface {
	Create(ctx context.Context, email string) error
}

type request struct {
	Email string `json:"email" validate:"required,email,lte=255"`
}

type Handler struct {
	recoveryService recoveryService
	logger          log.Logger
	validator       validation.Validator
}

func NewHandler(
	recoveryService recoveryService,
	logger log.Logger,
	validator validation.Validator,
) *Handler {
	return &Handler{
		recoveryService: recoveryService,
		logger:          logger,
		validator:       validator,
	}
}

func (h *Handler) Handle(ctx http.Context) {
	req, err := h.parseRequest(ctx)
	if err != nil {
		util.RespondBadRequest(ctx, err.Error())
		return
	}

	err = h.recoveryService.Create(ctx, req.Email)

	switch {
	case err == nil:
		util.Respond(ctx, nethttp.StatusOK, struct{}{})
	case errors.Is(err, apperrors.ErrUserNotFound):
		util.RespondBadRequest(ctx, "user with this email not found")
	default:
		h.logger.Error().Err(err).Msg("create recovery error on auth service")
		util.RespondInternalError(ctx)
	}
}

func (h *Handler) parseRequest(ctx http.Context) (*request, error) {
	req := &request{}

	if err := util.EnrichRequestBody(ctx, req); err != nil {
		return nil, err
	}

	if err := h.validator.Struct(req); err != nil {
		return nil, err
	}

	return req, nil
}
