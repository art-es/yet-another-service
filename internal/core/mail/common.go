//go:generate mockgen -source=common.go -destination=mock/common.go -package=mock
package mail

import (
	"context"
	"fmt"

	"github.com/art-es/yet-another-service/internal/app/shared/dto"
)

type mailRepository interface {
	Save(ctx context.Context, mails []dto.Mail) error
}

func saveMail(repo mailRepository, ctx context.Context, address, subject, content string) error {
	mail := dto.Mail{
		Address: address,
		Subject: subject,
		Content: content,
	}

	if err := repo.Save(ctx, []dto.Mail{mail}); err != nil {
		return fmt.Errorf("save mail: %w", err)
	}

	return nil
}
