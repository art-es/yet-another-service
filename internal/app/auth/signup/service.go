//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
package signup

import (
	"context"
	"fmt"

	"github.com/art-es/yet-another-service/internal/app/auth"
	"github.com/art-es/yet-another-service/internal/app/shared/errors"
	"github.com/art-es/yet-another-service/internal/app/shared/models"
	"github.com/art-es/yet-another-service/internal/core/transaction"
)

type hashGenerator interface {
	Generate(str string) (string, error)
}

type userRepository interface {
	Exists(ctx context.Context, email string) (bool, error)
	Save(ctx context.Context, tx transaction.Transaction, user *models.User) error
}

type activationService interface {
	CreateActivation(ctx context.Context, tx transaction.Transaction, user *models.User) error
}

type Service struct {
	hashGenerator     hashGenerator
	userRepository    userRepository
	activationService activationService
}

func NewService(
	hashGenerateService hashGenerator,
	userRepository userRepository,
	activationService activationService,
) *Service {
	return &Service{
		hashGenerator:     hashGenerateService,
		userRepository:    userRepository,
		activationService: activationService,
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

	if err = s.activationService.CreateActivation(ctx, tx, user); err != nil {
		return fmt.Errorf("create activation: %w", err)
	}

	return nil
}
