package password_recovery

import (
	"context"
	"fmt"
	"net/url"

	"github.com/art-es/yet-another-service/internal/app/shared/errors"
	"github.com/art-es/yet-another-service/internal/app/shared/models"
	"github.com/art-es/yet-another-service/internal/core/mail"
	"github.com/art-es/yet-another-service/internal/core/transaction"
)

func (s *Service) Create(ctx context.Context, email string) error {
	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("find user by email in repository: %w", err)
	}

	if user == nil {
		return errors.ErrUserNotFound
	}

	tx := transaction.New(ctx)

	if err = s.doCreationTransaction(ctx, tx, user); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (s *Service) doCreationTransaction(ctx context.Context, tx transaction.Transaction, user *models.User) error {
	recovery := &models.PasswordRecovery{
		UserID: user.ID,
	}

	if err := s.recoveryRepository.Save(ctx, tx, recovery); err != nil {
		return fmt.Errorf("save recovery in repository: %w", err)
	}

	mailData := mail.PasswordRecoveryData{
		RecoveryURL: newRecoveryURL(s.baseRecoveryURl, recovery.Token),
	}

	if err := s.recoveryMailer.MailTo(ctx, user.Email, mailData); err != nil {
		return fmt.Errorf("mail recovery to user: %w", err)
	}

	return nil
}

func newRecoveryURL(u url.URL, token string) string {
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()
	return u.String()
}
