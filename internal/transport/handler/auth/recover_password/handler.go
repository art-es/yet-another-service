//go:generate mockgen -source=handler.go -destination=mock/handler.go -package=mock
package recoverpassword

import (
	"context"
	"errors"
	nethttp "net/http"

	"github.com/art-es/yet-another-service/internal/app/auth"
	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
	"github.com/art-es/yet-another-service/internal/core/http"
	"github.com/art-es/yet-another-service/internal/core/http/util"
	"github.com/art-es/yet-another-service/internal/core/log"
	"github.com/art-es/yet-another-service/internal/core/validation"
)

type authService interface {
	Recover(ctx context.Context, in *auth.PasswordRecoverIn) error
}

type request struct {
	Token       string `json:"token" validate:"required,uuid"`
	OldPassword string `json:"oldPassword" validate:"required,lte=32"`
	NewPassword string `json:"newPassword" validate:"required,lte=32"`
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
		util.RespondBadRequest(ctx, err.Error())
		return
	}

	err = h.authService.Recover(ctx, &auth.PasswordRecoverIn{
		Token:       req.Token,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})

	switch {
	case err == nil:
		util.Respond(ctx, nethttp.StatusOK, struct{}{})
	case errors.Is(err, apperrors.ErrUserPasswordRecoveryNotFound):
		util.RespondNotFound(ctx)
	default:
		h.logger.Error().Err(err).Msg("recover error on auth service")
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
