//go:generate mockgen -source=handler.go -destination=mock/handler.go -package=mock
package login

import (
	"context"
	"errors"
	nethttp "net/http"

	"github.com/art-es/yet-another-service/internal/app/shared/dto"
	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
	"github.com/art-es/yet-another-service/internal/core/http"
	"github.com/art-es/yet-another-service/internal/core/http/util"
	"github.com/art-es/yet-another-service/internal/core/log"
	"github.com/art-es/yet-another-service/internal/core/validation"
)

const tokenType = "Bearer"

type authService interface {
	Login(ctx context.Context, req *dto.LoginIn) (*dto.LoginOut, error)
}

type request struct {
	Email    string `json:"email" validate:"required,email,lte=255"`
	Password string `json:"password" validate:"required,lte=32"`
}

type response struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
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

	out, err := h.authService.Login(ctx, &dto.LoginIn{
		Email:    req.Email,
		Password: req.Password,
	})

	switch {
	case err == nil:
		util.Respond(ctx, nethttp.StatusOK, response{
			AccessToken:  out.AccessToken,
			RefreshToken: out.RefreshToken,
			TokenType:    tokenType,
		})
	case errors.Is(err, apperrors.ErrUserNotFound) || errors.Is(err, apperrors.ErrWrongPassword):
		util.RespondBadRequest(ctx, "Wrong credentials.")
	default:
		h.logger.Error().Err(err).Msg("login error on auth service")
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
