//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
package password_recovery

import (
	"context"
	"net/url"

	"github.com/art-es/yet-another-service/internal/app/shared/models"
	"github.com/art-es/yet-another-service/internal/core/mail"
	"github.com/art-es/yet-another-service/internal/core/transaction"
)

type userRepository interface {
	Find(ctx context.Context, id string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Save(ctx context.Context, tx transaction.Transaction, user *models.User) error
}

type recoveryRepository interface {
	Find(ctx context.Context, token string) (*models.PasswordRecovery, error)
	Delete(ctx context.Context, tx transaction.Transaction, token string) error
	Save(ctx context.Context, tx transaction.Transaction, recovery *models.PasswordRecovery) error
}

type recoveryMailer interface {
	MailTo(ctx context.Context, address string, data mail.PasswordRecoveryData) error
}

type hashService interface {
	Check(str, hashStr string) error
	Generate(str string) (string, error)
}

type Service struct {
	baseRecoveryURl    url.URL
	userRepository     userRepository
	recoveryRepository recoveryRepository
	recoveryMailer     recoveryMailer
	hashService        hashService
}

func NewService(
	baseRecoveryURl url.URL,
	userRepository userRepository,
	recoveryRepository recoveryRepository,
	recoveryMailer recoveryMailer,
	hashService hashService,
) *Service {
	return &Service{
		baseRecoveryURl:    baseRecoveryURl,
		userRepository:     userRepository,
		recoveryRepository: recoveryRepository,
		recoveryMailer:     recoveryMailer,
		hashService:        hashService,
	}
}
