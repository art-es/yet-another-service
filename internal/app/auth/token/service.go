//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
package token

import (
	"fmt"
	"time"

	"github.com/art-es/yet-another-service/internal/app/auth"
)

var getCurrentTime = time.Now

type jwtService interface {
	Parse(token string) (*auth.TokenClaims, error)
	Generate(claims *auth.TokenClaims) (string, error)
}

type Service struct {
	jwtService jwtService
}

func NewService(jwtService jwtService) *Service {
	return &Service{
		jwtService: jwtService,
	}
}

func (s *Service) Generate(userID string) (*auth.TokenPair, error) {
	now := getCurrentTime()

	accessToken, err := s.jwtService.Generate(auth.NewAccessTokenClaims(now, userID))
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.Generate(auth.NewRefreshTokenClaims(now, userID))
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &auth.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Refresh(refreshToken string) (string, error) {
	refreshTokenClaims, err := s.jwtService.Parse(refreshToken)
	if err != nil {
		return "", fmt.Errorf("parse refresh token: %w", err)
	}

	accessTokenClaims := refreshTokenClaims.ToAccessToken(getCurrentTime())

	accessToken, err := s.jwtService.Generate(accessTokenClaims)
	if err != nil {
		return "", fmt.Errorf("generate access token: %w", err)
	}

	return accessToken, nil
}
