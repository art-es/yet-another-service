package logout

import (
	"context"
	"fmt"

	"github.com/art-es/yet-another-service/internal/app/auth"
	"github.com/art-es/yet-another-service/internal/core/log"
)

type tokenService interface {
	Invalidate(ctx context.Context, token string) error
}

type Service struct {
	tokenService tokenService
	logger       log.Logger
}

func NewService(
	tokenService tokenService,
	logger log.Logger,
) *Service {
	return &Service{
		tokenService: tokenService,
		logger:       logger,
	}
}

func (s *Service) Logout(ctx context.Context, req *auth.LogoutIn) error {
	if err := s.tokenService.Invalidate(ctx, req.RefreshToken); err != nil {
		return fmt.Errorf("invalidate refresh token: %w", err)
	}

	if req.AccessToken != nil {
		if err := s.tokenService.Invalidate(ctx, *req.AccessToken); err != nil {
			s.logger.Warn().Err(err).Msg("acccess token invalidate error")
		}
	}

	return nil
}
