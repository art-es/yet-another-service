//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
package activation

import (
	"context"
	"net/url"

	"github.com/art-es/yet-another-service/internal/app/shared/models"
	"github.com/art-es/yet-another-service/internal/core/mail"
	"github.com/art-es/yet-another-service/internal/core/transaction"
)

type activationRepository interface {
	Find(ctx context.Context, token string) (*models.UserActivation, error)
	Delete(ctx context.Context, tx transaction.Transaction, token string) error
	Save(ctx context.Context, tx transaction.Transaction, activation *models.UserActivation) error
}

type userRepository interface {
	Activate(ctx context.Context, tx transaction.Transaction, userID string) error
}

type activationMailer interface {
	MailTo(ctx context.Context, address string, data mail.UserActivationData) error
}

type Service struct {
	baseActivationURL    url.URL
	activationRepository activationRepository
	userRepository       userRepository
	activationMailer     activationMailer
}

func NewService(
	baseActivationURL url.URL,
	activationRepository activationRepository,
	userRepository userRepository,
	activationMailer activationMailer,
) *Service {
	return &Service{
		baseActivationURL:    baseActivationURL,
		activationRepository: activationRepository,
		userRepository:       userRepository,
		activationMailer:     activationMailer,
	}
}
