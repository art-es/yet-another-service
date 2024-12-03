//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
package signup

import (
	"context"
	"fmt"
	"net/url"

	"github.com/art-es/yet-another-service/internal/core/transaction"
	"github.com/art-es/yet-another-service/internal/domain/auth"
)

type hashGenerator interface {
	Generate(str string) (string, error)
}

type userRepository interface {
	Add(ctx context.Context, tx transaction.Transaction, user *auth.User) error
	EmailExists(ctx context.Context, email string) (bool, error)
}

type activationCreator interface {
	Create(ctx context.Context, tx transaction.Transaction, userID string) (*auth.Activation, error)
}

type mailSender interface {
	SendMail(address string, subject, content string) error
}

type Service struct {
	activationURL     url.URL
	hashGenerator     hashGenerator
	userRepository    userRepository
	activationCreator activationCreator
	mailSender        mailSender
}

func NewService(
	activationURL url.URL,
	hashGenerateService hashGenerator,
	userRepository userRepository,
	activationCreator activationCreator,
	mailSender mailSender,
) *Service {
	return &Service{
		activationURL:     activationURL,
		hashGenerator:     hashGenerateService,
		userRepository:    userRepository,
		activationCreator: activationCreator,
		mailSender:        mailSender,
	}
}

func (s *Service) Signup(ctx context.Context, req *auth.SignupRequest) error {
	emailExists, err := s.userRepository.EmailExists(ctx, req.Email)
	if err != nil {
		return fmt.Errorf("check user email exists in repository: %w", err)
	}

	if emailExists {
		return auth.ErrEmailAlreadyTaken
	}

	passwordHash, err := s.hashGenerator.Generate(req.Password)
	if err != nil {
		return fmt.Errorf("generate password hash: %w", err)
	}

	user := &auth.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}

	tx := transaction.New(ctx)

	if err = s.doTransaction(ctx, tx, user); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (s *Service) doTransaction(ctx context.Context, tx transaction.Transaction, user *auth.User) error {
	if err := s.userRepository.Add(ctx, tx, user); err != nil {
		return fmt.Errorf("add user to repository: %w", err)
	}

	activation, err := s.activationCreator.Create(ctx, tx, user.ID)
	if err != nil {
		return fmt.Errorf("create activation: %w", err)
	}

	if err = s.sendActivationMail(user.Email, activation.Token); err != nil {
		return err
	}

	return nil
}

func (s *Service) sendActivationMail(address, token string) error {
	content, err := buildActivationMailContent(s.activationURL, token)
	if err != nil {
		return err
	}

	if err = s.mailSender.SendMail(address, activationMailSubject, content); err != nil {
		return fmt.Errorf("send activation mail: %w", err)
	}

	return nil
}
