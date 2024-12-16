//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
package mailing

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/art-es/yet-another-service/internal/core/log"
	"github.com/art-es/yet-another-service/internal/core/retrier"
	"github.com/art-es/yet-another-service/internal/domain/shared/models"
)

type mailRepository interface {
	Get(ctx context.Context) ([]models.Mail, error)
	Save(ctx context.Context, mails []models.Mail) error
}

type mailer interface {
	MailTo(address, subject, content string) error
}

type Config struct {
	Delay          time.Duration
	ProcessTimeout time.Duration
}

type Service struct {
	config         Config
	processRetrier retrier.Retrier
	saveRetrier    retrier.Retrier
	mailRepository mailRepository
	mailer         mailer
	logger         log.Logger
	done           bool
}

func NewService(
	config Config,
	processRetrier retrier.Retrier,
	saveRetrier retrier.Retrier,
	mailRepository mailRepository,
	mailer mailer,
	logger log.Logger,
) *Service {
	return &Service{
		config:         config,
		processRetrier: processRetrier,
		saveRetrier:    saveRetrier,
		mailRepository: mailRepository,
		mailer:         mailer,
		logger:         logger,
	}
}

func (s *Service) Run(ctx context.Context) error {
	for !s.done {
		err := func() error {
			pctx, cancel := context.WithTimeout(ctx, s.config.ProcessTimeout)
			defer cancel()

			return s.processRetrier.Process(func() error { return s.process(pctx) })
		}()

		if err != nil {
			return fmt.Errorf("process: %w", err)
		}

		time.Sleep(s.config.Delay)
	}

	return nil
}

func (s *Service) process(ctx context.Context) error {
	mails, err := s.mailRepository.Get(ctx)
	if err != nil {
		return fmt.Errorf("get mails from repository: %w", err)
	}

	if len(mails) == 0 {
		s.done = true
		return nil
	}

	mailed := 0
	for i := range mails {
		if err = s.mailer.MailTo(mails[i].Address, mails[i].Subject, mails[i].Content); err != nil {
			s.logger.Error().Err(err).Str("mail_id", mails[i].ID).Msg("mail error")
			continue
		}

		mailed++
		mails[i].SetMailed()
	}

	if mailed == 0 {
		return errors.New("all mails finished with error")
	}

	err = s.saveRetrier.Process(func() error { return s.mailRepository.Save(ctx, mails) })
	if err != nil {
		return fmt.Errorf("save in repository: %w", err)
	}

	return nil
}
