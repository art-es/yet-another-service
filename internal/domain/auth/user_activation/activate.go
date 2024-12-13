package user_activation

import (
	"context"
	"fmt"

	"github.com/art-es/yet-another-service/internal/core/transaction"
	"github.com/art-es/yet-another-service/internal/domain/shared/errors"
	"github.com/art-es/yet-another-service/internal/domain/shared/models"
)

func (s *Service) Activate(ctx context.Context, token string) error {
	activation, err := s.activationRepository.Find(ctx, token)
	if err != nil {
		return fmt.Errorf("find activation in repository: %w", err)
	}

	if activation == nil {
		return errors.ErrUserActivationNotFound
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

func (s *Service) doTransaction(ctx context.Context, tx transaction.Transaction, activation *models.UserActivation) error {
	if err := s.userRepository.Activate(ctx, tx, activation.UserID); err != nil {
		return fmt.Errorf("activate user in repository: %w", err)
	}

	if err := s.activationRepository.Delete(ctx, tx, activation.Token); err != nil {
		return fmt.Errorf("delete activation by token in repository: %w", err)
	}

	return nil
}
