package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/art-es/yet-another-service/internal/core/log"
	"github.com/art-es/yet-another-service/internal/domain/mailing"
	"github.com/art-es/yet-another-service/internal/driver/smtp"
)

const (
	appEnvLocal   = "local"
	appEnvStaging = "staging"
	appEnvProd    = "prod"
)

var availableAppEnv = []string{
	appEnvLocal, appEnvStaging, appEnvProd,
}

type appConfig struct {
	appEnv             string
	postgresURL        string
	mailing            mailing.Config
	smtp               smtp.Config
	maxProcessRetries  int
	maxSaveMailRetries int

	logger log.Logger
}

func getAppConfig(logger log.Logger) *appConfig {
	c := &appConfig{logger: logger}
	c.initAppEnv()
	c.initPostgresURL()
	c.initSMTP()
	c.initMaxRetries()
	return c
}

func (c *appConfig) initAppEnv() {
	c.appEnv = os.Getenv("APP_ENV")

	if c.appEnv == "" {
		c.appEnv = appEnvLocal
		c.logger.Warn().
			Msg(fmt.Sprintf("APP_ENV is empty, using: %s", appEnvLocal))
	}

	if !slices.Contains(availableAppEnv, c.appEnv) {
		c.logger.Panic().
			Str("value", c.appEnv).
			Str("available_values", fmt.Sprintf("%v", availableAppEnv)).
			Msg("APP_ENV has unavailable value")
	}
}

func (c *appConfig) initPostgresURL() {
	if c.postgresURL = os.Getenv("POSTGRES_URL"); c.postgresURL != "" {
		return
	}

	if c.appEnv != appEnvLocal {
		c.logger.Panic().Msg("POSTGRES_URL is required")
	}

	c.postgresURL = "http://postgres:postgres@127.0.0.1:5432/master?ssldisabled=true"
}

func (c *appConfig) initMailing() {
	delay, _ := strconv.Atoi(os.Getenv("MAILING_DELAY"))
	if delay < 100 {
		delay = 300
	}

	processTimeout, _ := strconv.Atoi(os.Getenv("MAILING_PROCESS_TIMEOUT"))
	if processTimeout < 0 {
		processTimeout = 300
	}

	c.mailing.Delay = time.Duration(delay) * time.Millisecond
	c.mailing.ProcessTimeout = time.Duration(processTimeout) * time.Millisecond
}

func (c *appConfig) initSMTP() {
	c.smtp.Host = os.Getenv("SMTP_HOST")
	c.smtp.Port, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))
	c.smtp.Username = os.Getenv("SMTP_USERNAME")
	c.smtp.Password = os.Getenv("SMTP_PASSWORD")

	if c.smtp.Host != "" && c.smtp.Username != "" {
		return
	}

	if c.appEnv != appEnvLocal {
		c.logger.Panic().Msg("SMTP_HOST, SMTP_USERNAME are required")
	}
}

func (c *appConfig) initMaxRetries() {
	processMaxRetries, _ := strconv.Atoi(os.Getenv("PROCESS_MAX_RETRIES"))
	saveMailMaxRetries, _ := strconv.Atoi(os.Getenv("SAVE_MAIL_MAX_RETRIES"))

	if processMaxRetries < 1 {
		processMaxRetries = 5
	}
	if saveMailMaxRetries < 1 {
		saveMailMaxRetries = 3
	}
}
