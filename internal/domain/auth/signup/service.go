//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
package signup

import (
	"context"
	"fmt"
	"net/url"

	"github.com/art-es/yet-another-service/internal/core/mail"
	"github.com/art-es/yet-another-service/internal/core/transaction"
	"github.com/art-es/yet-another-service/internal/domain/auth"
	"github.com/art-es/yet-another-service/internal/domain/shared/errors"
	"github.com/art-es/yet-another-service/internal/domain/shared/models"
)

type hashGenerator interface {
	Generate(str string) (string, error)
}

type userRepository interface {
	Exists(ctx context.Context, email string) (bool, error)
	Save(ctx context.Context, tx transaction.Transaction, user *models.User) error
}

type activationRepository interface {
	Save(ctx context.Context, tx transaction.Transaction, activation *models.UserActivation) error
}

type activationMailer interface {
	MailTo(address string, data mail.UserActivationData) error
}

type Service struct {
	baseActivationURL    url.URL
	hashGenerator        hashGenerator
	userRepository       userRepository
	activationRepository activationRepository
	activationMailer     activationMailer
}

func NewService(
	activationURL url.URL,
	hashGenerateService hashGenerator,
	userRepository userRepository,
	activationRepository activationRepository,
	activationMailer activationMailer,
) *Service {
	return &Service{
		baseActivationURL:    activationURL,
		hashGenerator:        hashGenerateService,
		userRepository:       userRepository,
		activationRepository: activationRepository,
		activationMailer:     activationMailer,
	}
}

func (s *Service) Signup(ctx context.Context, in *auth.SignupIn) error {
	userExists, err := s.userRepository.Exists(ctx, in.Email)
	if err != nil {
		return fmt.Errorf("check user exists in repository: %w", err)
	}

	if userExists {
		return errors.ErrEmailAlreadyTaken
	}

	tx := transaction.New(ctx)

	if err = s.doTransaction(ctx, tx, in); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (s *Service) doTransaction(ctx context.Context, tx transaction.Transaction, in *auth.SignupIn) error {
	passwordHash, err := s.hashGenerator.Generate(in.Password)
	if err != nil {
		return fmt.Errorf("generate password hash: %w", err)
	}

	user := &models.User{
		Name:         in.Name,
		Email:        in.Email,
		PasswordHash: passwordHash,
	}

	if err = s.userRepository.Save(ctx, tx, user); err != nil {
		return fmt.Errorf("save user in repository: %w", err)
	}

	activation := &models.UserActivation{
		UserID: user.ID,
	}

	if err = s.activationRepository.Save(ctx, tx, activation); err != nil {
		return fmt.Errorf("save activation in repository: %w", err)
	}

	activationMailData := mail.UserActivationData{
		ActivationURL: newActivationURL(s.baseActivationURL, activation.Token),
	}

	if err = s.activationMailer.MailTo(user.Email, activationMailData); err != nil {
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
