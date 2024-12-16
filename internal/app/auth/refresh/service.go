//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
package refresh

import (
	"fmt"
	"time"

	"github.com/art-es/yet-another-service/internal/app/auth"
)

var getCurrentTime = time.Now

type tokenService interface {
	Parse(token string) (*auth.TokenClaims, error)
	Generate(claims *auth.TokenClaims) (string, error)
}

type Service struct {
	tokenService tokenService
}

func NewService(tokenService tokenService) *Service {
	return &Service{
		tokenService: tokenService,
	}
}

func (s *Service) Refresh(refreshToken string) (string, error) {
	refreshTokenClaims, err := s.tokenService.Parse(refreshToken)
	if err != nil {
		return "", fmt.Errorf("parse refresh token: %w", err)
	}

	accessTokenClaims := refreshTokenClaims.ToAccessToken(getCurrentTime())

	accessToken, err := s.tokenService.Generate(accessTokenClaims)
	if err != nil {
		return "", fmt.Errorf("generate access token: %w", err)
	}

	return accessToken, nil
}
