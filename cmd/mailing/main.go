package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/art-es/yet-another-service/internal/app/mailing"
	"github.com/art-es/yet-another-service/internal/core/retrier"
	pq "github.com/art-es/yet-another-service/internal/driver/postgres"
	"github.com/art-es/yet-another-service/internal/driver/smtp"
	"github.com/art-es/yet-another-service/internal/driver/zerolog"
	pqstorage "github.com/art-es/yet-another-service/internal/storage/postgres"
)

func main() {
	logger := zerolog.NewLogger()
	config := getAppConfig(logger)
	pqDB := pq.Connect(config.postgresURL)

	// Dependencies
	processRetrier := retrier.New(logger, config.maxProcessRetries, time.Millisecond*500)
	saveMailRetrier := retrier.New(logger, config.maxSaveMailRetries, time.Millisecond*500)
	mailStorage := pqstorage.NewMailStorage(pqDB)
	smtpService := smtp.NewService(config.smtp)

	mailingService := mailing.NewService(config.mailing, processRetrier, saveMailRetrier, mailStorage, smtpService, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	if err := mailingService.Run(ctx); err != nil {
		logger.Error().Err(err).Msg("mailing run error")
	}
}
