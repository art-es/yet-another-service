package logout

import (
	"context"
	"fmt"
	"time"

	"github.com/art-es/yet-another-service/internal/core/log"
	"github.com/art-es/yet-another-service/internal/domain/auth"
)

type tokenParser interface {
	Parse(token string) (*auth.TokenClaims, error)
}

type tokenBlackList interface {
	Add(ctx context.Context, token string, ttl time.Duration) error
}

type Service struct {
	tokenParser    tokenParser
	tokenBlackList tokenBlackList
	logger         log.Logger
}

func NewService(
	tokenParser tokenParser,
	tokenBlackList tokenBlackList,
	logger log.Logger,
) *Service {
	return &Service{
		tokenParser:    tokenParser,
		tokenBlackList: tokenBlackList,
		logger:         logger,
	}
}

func (s *Service) Logout(ctx context.Context, req *auth.LogoutRequest) error {
	now := time.Now()

	if err := s.invalidate(ctx, req.RefreshToken, now); err != nil {
		return fmt.Errorf("invalidate refresh token: %w", err)
	}

	if req.AccessToken != nil {
		if err := s.invalidate(ctx, *req.AccessToken, now); err != nil {
			s.logger.Warn().Err(err).Msg("acccess token invalidate error")
		}
	}

	return nil
}

func (s *Service) invalidate(ctx context.Context, token string, now time.Time) error {
	tokenClaims, err := s.tokenParser.Parse(token)
	if err != nil {
		return fmt.Errorf("verify token: %w", err)
	}

	blackListTTL := tokenClaims.ExpiresAt.Sub(now)

	if err = s.tokenBlackList.Add(ctx, token, blackListTTL); err != nil {
		return fmt.Errorf("add to black list: %w", err)
	}

	return nil
}
