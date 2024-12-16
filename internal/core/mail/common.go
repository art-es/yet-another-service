//go:generate mockgen -source=common.go -destination=mock/common.go -package=mock
package mail

import (
	"context"
	"fmt"

	"github.com/art-es/yet-another-service/internal/domain/shared/models"
)

type mailRepository interface {
	Save(ctx context.Context, mails []models.Mail) error
}

func saveMail(repo mailRepository, ctx context.Context, address, subject, content string) error {
	mail := models.Mail{
		Address: address,
		Subject: subject,
		Content: content,
	}

	if err := repo.Save(ctx, []models.Mail{mail}); err != nil {
		return fmt.Errorf("save mail: %w", err)
	}

	return nil
}
