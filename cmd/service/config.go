package main

import (
	"fmt"
	"net/url"
	"os"
	"slices"

	"github.com/art-es/yet-another-service/internal/core/log"
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
	appEnv                  string
	postgresURL             string
	redisAddr               string
	jwtSecret               string
	userActivationURL       url.URL
	userPasswordRecoveryURL url.URL

	logger log.Logger
}

func getAppConfig(logger log.Logger) *appConfig {
	c := &appConfig{logger: logger}
	c.initAppEnv()
	c.initPostgresURL()
	c.initJWTSecret()
	c.initUserActivationURL()
	c.initUserPasswordRecoveryURL()
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

func (c *appConfig) initRedisAddr() {
	if c.postgresURL = os.Getenv("REDIS_ADDR"); c.postgresURL != "" {
		return
	}

	if c.appEnv != appEnvLocal {
		c.logger.Panic().Msg("REDIS_ADDR is required")
	}

	c.postgresURL = "127.0.0.1:6379"
}

func (c *appConfig) initJWTSecret() {
	if c.jwtSecret = os.Getenv("JWT_SECRET"); c.jwtSecret != "" {
		return
	}

	if c.appEnv == appEnvProd {
		c.logger.Panic().Msg("JWT_SECRET is required in prod")
	}

	c.appEnv = "secret"
}

func (c *appConfig) initUserActivationURL() {
	rawURL := os.Getenv("USER_ACTIVATION_URL")

	if rawURL == "" {
		if c.appEnv != appEnvLocal {
			c.logger.Panic().Msg("USER_ACTIVATION_URL is required")
		}

		rawURL = "http://127.0.0.1/activate-user"
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		c.logger.Panic().Msg("USER_ACTIVATION_URL has invalid URL")
	}

	c.userActivationURL = *u
}

func (c *appConfig) initUserPasswordRecoveryURL() {
	rawURL := os.Getenv("USER_PASSWORD_RECOVERY_URL")

	if rawURL == "" {
		if c.appEnv != appEnvLocal {
			c.logger.Panic().Msg("USER_PASSWORD_RECOVERY_URL is required")
		}

		rawURL = "http://127.0.0.1/recover-password"
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		c.logger.Panic().Msg("USER_PASSWORD_RECOVERY_URL has invalid URL")
	}

	c.userPasswordRecoveryURL = *u
}
