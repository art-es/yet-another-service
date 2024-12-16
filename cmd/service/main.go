package main

import (
	"net/http"

	app_login "github.com/art-es/yet-another-service/internal/app/auth/login"
	app_logout "github.com/art-es/yet-another-service/internal/app/auth/logout"
	app_password_recovery "github.com/art-es/yet-another-service/internal/app/auth/password_recovery"
	app_refresh "github.com/art-es/yet-another-service/internal/app/auth/refresh"
	app_signup "github.com/art-es/yet-another-service/internal/app/auth/signup"
	app_user_activation "github.com/art-es/yet-another-service/internal/app/auth/user_activation"
	"github.com/art-es/yet-another-service/internal/core/mail"
	driver_bcrypt "github.com/art-es/yet-another-service/internal/driver/bcrypt"
	driver_gin "github.com/art-es/yet-another-service/internal/driver/gin"
	driver_jwt "github.com/art-es/yet-another-service/internal/driver/jwt"
	driver_postgres "github.com/art-es/yet-another-service/internal/driver/postgres"
	driver_validator "github.com/art-es/yet-another-service/internal/driver/validator"
	driver_zerolog "github.com/art-es/yet-another-service/internal/driver/zerolog"
	storage_postgres "github.com/art-es/yet-another-service/internal/storage/postgres"
	storage_redis "github.com/art-es/yet-another-service/internal/storage/redis"
	transport_user_activate "github.com/art-es/yet-another-service/internal/transport/handler/auth/activate"
	transport_forgot_password "github.com/art-es/yet-another-service/internal/transport/handler/auth/forgot_password"
	transport_login "github.com/art-es/yet-another-service/internal/transport/handler/auth/login"
	transport_logout "github.com/art-es/yet-another-service/internal/transport/handler/auth/logout"
	transport_recover_password "github.com/art-es/yet-another-service/internal/transport/handler/auth/recover_password"
	transport_refresh "github.com/art-es/yet-another-service/internal/transport/handler/auth/refresh"
	transport_signup "github.com/art-es/yet-another-service/internal/transport/handler/auth/signup"
)

func main() {
	logger := driver_zerolog.NewLogger()
	config := getAppConfig(logger)

	// Drivers
	validator := driver_validator.New()
	hashService := driver_bcrypt.NewHashService()
	postgresDB := driver_postgres.Connect(config.postgresURL)
	jwtService := driver_jwt.NewService(config.jwtSecret)

	// Data Layer
	userStorage := storage_postgres.NewUserStorage(postgresDB)
	userActivationStorage := storage_postgres.NewUserActivationStorage(postgresDB)
	passwordRecoveryStorage := storage_postgres.NewPasswordRecoveryStorage(postgresDB)
	mailStorage := storage_postgres.NewMailStorage(postgresDB)
	authTokenBlackListStorage := storage_redis.NewAuthTokenBlackListStorage()

	// Mailers
	userActivationMailer := mail.NewUserActivationMailer(mailStorage)
	passwordRecoveryMailer := mail.NewPasswordRecoveryMailer(mailStorage)

	// App Layer
	userActivationService := app_user_activation.NewService(config.userActivationURL, userActivationStorage, userStorage, userActivationMailer)
	signupService := app_signup.NewService(hashService, userStorage, userActivationService)
	loginService := app_login.NewService(userStorage, hashService, jwtService)
	logoutService := app_logout.NewService(jwtService, authTokenBlackListStorage, logger)
	refreshService := app_refresh.NewService(jwtService)
	passwordRecoveryService := app_password_recovery.NewService(config.userPasswordRecoveryURL, userStorage, passwordRecoveryStorage, passwordRecoveryMailer, hashService)

	// Transport Layer
	signupHandler := transport_signup.NewHandler(signupService, logger, validator)
	userActivateHandler := transport_user_activate.NewHandler(userActivationService, logger, validator)
	loginHandler := transport_login.NewHandler(loginService, logger, validator)
	logoutHandler := transport_logout.NewHandler(logoutService, logger, validator)
	refreshHandler := transport_refresh.NewHandler(refreshService, logger)
	forgotPasswordHandler := transport_forgot_password.NewHandler(passwordRecoveryService, logger, validator)
	recoverPasswordHandler := transport_recover_password.NewHandler(passwordRecoveryService, logger, validator)

	router := driver_gin.NewRouter()
	router.Register(http.MethodPost, "/auth/signup", signupHandler.Handle)
	router.Register(http.MethodGet, "/auth/activate", userActivateHandler.Handle)
	router.Register(http.MethodPost, "/auth/login", loginHandler.Handle)
	router.Register(http.MethodPost, "/auth/logout", logoutHandler.Handle)
	router.Register(http.MethodPost, "/auth/refresh", refreshHandler.Handle)
	router.Register(http.MethodPost, "/auth/forgot-password", forgotPasswordHandler.Handle)
	router.Register(http.MethodPost, "/auth/recover-password", recoverPasswordHandler.Handle)

	if err := router.Run(); err != nil {
		logger.Panic().Err(err).Msg("router run error")
	}
}
