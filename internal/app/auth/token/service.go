//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/art-es/yet-another-service/internal/app/auth"
	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
)

var getCurrentTime = time.Now

type jwtService interface {
	Parse(token string) (*auth.TokenClaims, error)
	Generate(claims *auth.TokenClaims) (string, error)
}

type blackList interface {
	Add(ctx context.Context, token string, ttl time.Duration) error
	Has(ctx context.Context, token string) (bool, error)
}

type Service struct {
	jwtService jwtService
	blackList  blackList
}

func NewService(jwtService jwtService, blackList blackList) *Service {
	return &Service{
		jwtService: jwtService,
		blackList:  blackList,
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

func (s *Service) Refresh(ctx context.Context, refreshToken string) (string, error) {
	claims, err := s.jwtService.Parse(refreshToken)
	if err != nil {
		return "", fmt.Errorf("parse refresh token: %w", err)
	}

	blacklisted, err := s.blackList.Has(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("check refresh token in black list: %w", err)
	}

	if blacklisted {
		return "", apperrors.ErrInvalidAuthToken
	}

	accessTokenClaims := claims.ToAccessToken(getCurrentTime())

	accessToken, err := s.jwtService.Generate(accessTokenClaims)
	if err != nil {
		return "", fmt.Errorf("generate access token: %w", err)
	}

	return accessToken, nil
}

func (s *Service) Authorize(ctx context.Context, accessToken string) (string, error) {
	claims, err := s.jwtService.Parse(accessToken)
	if err != nil {
		return "", fmt.Errorf("parse access token: %w", err)
	}

	blacklisted, err := s.blackList.Has(ctx, accessToken)
	if err != nil {
		return "", fmt.Errorf("check access token in black list: %w", err)
	}

	if blacklisted {
		return "", apperrors.ErrInvalidAuthToken
	}

	return claims.UserID, nil
}

func (s *Service) Invalidate(ctx context.Context, token string) error {
	now := getCurrentTime()

	tokenClaims, err := s.jwtService.Parse(token)
	if err != nil {
		return fmt.Errorf("parse token: %w", err)
	}

	blackListTTL := tokenClaims.ExpiresAt.Sub(now) * -1
	if blackListTTL < 0 {
		return errors.New("black list TTL is negative")
	}

	if err = s.blackList.Add(ctx, token, blackListTTL); err != nil {
		return fmt.Errorf("add to black list: %w", err)
	}

	return nil
}
