package activation

import (
	"context"
	"fmt"
	"net/url"

	"github.com/art-es/yet-another-service/internal/app/shared/models"
	"github.com/art-es/yet-another-service/internal/core/mail"
	"github.com/art-es/yet-another-service/internal/core/transaction"
)

func (s *Service) Create(ctx context.Context, tx transaction.Transaction, user *models.User) error {
	activation := &models.UserActivation{
		UserID: user.ID,
	}

	if err := s.activationRepository.Save(ctx, tx, activation); err != nil {
		return fmt.Errorf("save activation in repository: %w", err)
	}

	mailData := mail.UserActivationData{
		ActivationURL: newActivationURL(s.baseActivationURL, activation.Token),
	}

	if err := s.activationMailer.MailTo(ctx, user.Email, mailData); err != nil {
		return fmt.Errorf("mail activation to user: %w", err)
	}

	return nil
}

func newActivationURL(u url.URL, token string) string {
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()
	return u.String()
}
