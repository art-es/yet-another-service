package password_recovery

import (
	"context"
	"fmt"

	"github.com/art-es/yet-another-service/internal/app/shared/dto"
	"github.com/art-es/yet-another-service/internal/app/shared/errors"
	"github.com/art-es/yet-another-service/internal/core/transaction"
)

func (s *Service) Recover(ctx context.Context, in *dto.PasswordRecoverIn) error {
	recovery, err := s.recoveryRepository.Find(ctx, in.Token)
	if err != nil {
		return fmt.Errorf("find recovery in repository: %w", err)
	}

	if recovery == nil {
		return errors.ErrUserPasswordRecoveryNotFound
	}

	user, err := s.userRepository.Find(ctx, recovery.UserID)
	if err != nil {
		return fmt.Errorf("find user in repository: %w", err)
	}

	if user == nil {
		return errors.ErrUserNotFound
	}

	if err = s.hashService.Check(in.OldPassword, user.PasswordHash); err != nil {
		return fmt.Errorf("check old password with hash: %w", err)
	}

	newPasswordHash, err := s.hashService.Generate(in.NewPassword)
	if err != nil {
		return fmt.Errorf("generate new password hash: %w", err)
	}

	user.PasswordHash = newPasswordHash

	tx := transaction.New(ctx)

	if err = s.doRecoverTransaction(ctx, tx, user, recovery); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (s *Service) doRecoverTransaction(
	ctx context.Context,
	tx transaction.Transaction,
	user *dto.User,
	recovery *dto.PasswordRecovery,
) error {
	if err := s.userRepository.Save(ctx, tx, user); err != nil {
		return fmt.Errorf("save user in repository: %w", err)
	}

	if err := s.recoveryRepository.Delete(ctx, tx, recovery.Token); err != nil {
		return fmt.Errorf("delete recovery in repository: %w", err)
	}

	return nil
}
