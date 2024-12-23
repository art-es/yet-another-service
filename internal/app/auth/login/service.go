//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
package login

import (
	"context"
	"fmt"

	"github.com/art-es/yet-another-service/internal/app/shared/dto"
	"github.com/art-es/yet-another-service/internal/app/shared/errors"
)

type userRepository interface {
	FindByEmail(ctx context.Context, email string) (*dto.User, error)
}

type hashChecker interface {
	Check(str, hashStr string) error
}

type tokenGenerator interface {
	Generate(userID string) (*dto.AuthTokenPair, error)
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

func (s *Service) Login(ctx context.Context, req *dto.LoginIn) (*dto.LoginOut, error) {
	user, err := s.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("find user by email in repository: %w", err)
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	if err = s.hashChecker.Check(req.Password, user.PasswordHash); err != nil {
		if err == errors.ErrHashMismatched {
			return nil, errors.ErrWrongPassword
		}

		return nil, fmt.Errorf("check password by hash: %w", err)
	}

	tokenPair, err := s.tokenGenerator.Generate(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return &dto.LoginOut{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
