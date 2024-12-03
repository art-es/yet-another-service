//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
package activate

import (
	"context"
	"fmt"

	"github.com/art-es/yet-another-service/internal/core/transaction"
	"github.com/art-es/yet-another-service/internal/domain/auth"
)

type activationRepository interface {
	FindByToken(ctx context.Context, token string) (*auth.Activation, error)
	Delete(ctx context.Context, tx transaction.Transaction, token string) error
}

type userRepository interface {
	Activate(ctx context.Context, tx transaction.Transaction, userID string) error
}

type Service struct {
	activationRepository activationRepository
	userRepository       userRepository
}

func NewService(
	activationRepository activationRepository,
	userRepository userRepository,
) *Service {
	return &Service{
		activationRepository: activationRepository,
		userRepository:       userRepository,
	}
}

func (s *Service) Activate(ctx context.Context, token string) error {
	activation, err := s.activationRepository.FindByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("find activation by token in repository: %w", err)
	}

	if activation == nil {
		return auth.ErrActivationNotFound
	}

	tx := transaction.New(ctx)

	if err = s.doTransaction(ctx, tx, activation); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (s *Service) doTransaction(ctx context.Context, tx transaction.Transaction, activation *auth.Activation) error {
	if err := s.userRepository.Activate(ctx, tx, activation.UserID); err != nil {
		return fmt.Errorf("activate user in repository: %w", err)
	}

	if err := s.activationRepository.Delete(ctx, tx, activation.Token); err != nil {
		return fmt.Errorf("delete activation by token in repository: %w", err)
	}

	return nil
}
