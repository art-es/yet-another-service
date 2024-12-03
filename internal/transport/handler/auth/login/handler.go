//go:generate mockgen -source=handler.go -destination=mock/handler.go -package=mock
package login

import (
	"context"
	"errors"
	nethttp "net/http"

	"github.com/art-es/yet-another-service/internal/core/http"
	"github.com/art-es/yet-another-service/internal/core/log"
	"github.com/art-es/yet-another-service/internal/core/validation"
	"github.com/art-es/yet-another-service/internal/domain/auth"
)

const tokenType = "Bearer"

type authService interface {
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResult, error)
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
		http.RespondBadRequest(ctx, err.Error())
		return
	}

	res, err := h.authService.Login(ctx, &auth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	switch {
	case err == nil:
		http.Respond(ctx, nethttp.StatusOK, response{
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
			TokenType:    tokenType,
		})
	case errors.Is(err, auth.ErrUserNotFound) || errors.Is(err, auth.ErrWrongPassword):
		http.RespondBadRequest(ctx, "Wrong credentials.")
	default:
		h.logger.Error().Err(err).Msg("login error on auth service")
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
