//go:generate mockgen -source=handler.go -destination=mock/handler.go -package=mock
package logout

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
	Logout(ctx context.Context, req *auth.LogoutIn) error
}

type request struct {
	AccessToken  *string `json:"-" validate:"omitnil,len=70"`
	RefreshToken string  `json:"refreshToken" validate:"required,len=70"`
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

	err = h.authService.Logout(ctx, &auth.LogoutIn{
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
	})

	switch {
	case err == nil:
		util.Respond(ctx, nethttp.StatusOK, struct{}{})
	case errors.Is(err, apperrors.ErrInvalidAuthToken):
		util.RespondBadRequest(ctx, err.Error())
	default:
		h.logger.Error().Err(err).Msg("logout error on auth service")
		util.RespondInternalError(ctx)
	}
}

func (h *Handler) parseRequest(ctx http.Context) (*request, error) {
	req := &request{}

	if err := util.EnrichRequestBody(ctx, req); err != nil {
		return nil, err
	}

	if accessToken, ok := util.GetAuthorizationToken(ctx); ok {
		req.AccessToken = &accessToken
	}

	if err := h.validator.Struct(req); err != nil {
		return nil, err
	}

	return req, nil
}
