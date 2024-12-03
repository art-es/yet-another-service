//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
package login

import (
	"context"
	"fmt"
	"time"

	"github.com/art-es/yet-another-service/internal/domain/auth"
	"github.com/art-es/yet-another-service/internal/domain/hash"
)

var getCurrentTime = time.Now

type userRepository interface {
	FindByEmail(ctx context.Context, email string) (*auth.User, error)
}

type hashChecker interface {
	Check(str, hashStr string) error
}

type tokenGenerator interface {
	Generate(claims *auth.TokenClaims) (string, error)
}

type Service struct {
	userRepository userRepository
	hashChecker    hashChecker
	tokenGenerator tokenGenerator
}

func NewService(
	userRepository userRepository,
	hashChecker hashChecker,
	tokenGenerator tokenGenerator,
) *Service {
	return &Service{
		userRepository: userRepository,
		hashChecker:    hashChecker,
		tokenGenerator: tokenGenerator,
	}
}

func (s *Service) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResult, error) {
	user, err := s.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("find user by email in repository: %w", err)
	}

	if user == nil {
		return nil, auth.ErrUserNotFound
	}

	if err = s.hashChecker.Check(req.Password, user.PasswordHash); err != nil {
		if err == hash.ErrMismatched {
			return nil, auth.ErrWrongPassword
		}

		return nil, fmt.Errorf("check password by hash: %w", err)
	}

	now := getCurrentTime()

	accessToken, err := s.tokenGenerator.Generate(auth.NewAccessTokenClaims(now, user.ID))
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.tokenGenerator.Generate(auth.NewRefreshTokenClaims(now, user.ID))
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &auth.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
