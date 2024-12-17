package main

import (
	"net/http"

	"github.com/art-es/yet-another-service/internal/app/auth/login"
	"github.com/art-es/yet-another-service/internal/app/auth/logout"
	"github.com/art-es/yet-another-service/internal/app/auth/signup"
	authtoken "github.com/art-es/yet-another-service/internal/app/auth/token"
	useractivation "github.com/art-es/yet-another-service/internal/app/user/activation"
	passwordrecovery "github.com/art-es/yet-another-service/internal/app/user/password_recovery"
	"github.com/art-es/yet-another-service/internal/core/mail"
	bcryptDriver "github.com/art-es/yet-another-service/internal/driver/bcrypt"
	driverGin "github.com/art-es/yet-another-service/internal/driver/gin"
	jwtDriver "github.com/art-es/yet-another-service/internal/driver/jwt"
	pqDriver "github.com/art-es/yet-another-service/internal/driver/postgres"
	validatorDriver "github.com/art-es/yet-another-service/internal/driver/validator"
	"github.com/art-es/yet-another-service/internal/driver/zerolog"
	pqstorage "github.com/art-es/yet-another-service/internal/storage/postgres"
	redisstorage "github.com/art-es/yet-another-service/internal/storage/redis"
	useractivateTransport "github.com/art-es/yet-another-service/internal/transport/handler/auth/activate"
	forgotpasswordTransport "github.com/art-es/yet-another-service/internal/transport/handler/auth/forgot_password"
	loginTransport "github.com/art-es/yet-another-service/internal/transport/handler/auth/login"
	logoutTransport "github.com/art-es/yet-another-service/internal/transport/handler/auth/logout"
	recoverpasswordTransport "github.com/art-es/yet-another-service/internal/transport/handler/auth/recover_password"
	refreshtokenTransport "github.com/art-es/yet-another-service/internal/transport/handler/auth/refresh"
	signupTransport "github.com/art-es/yet-another-service/internal/transport/handler/auth/signup"
)

func main() {
	logger := zerolog.NewLogger()
	config := getAppConfig(logger)

	// Drivers
	pqDB := pqDriver.Connect(config.postgresURL)
	validator := validatorDriver.New()
	hashService := bcryptDriver.NewHashService()
	jwtService := jwtDriver.NewService(config.jwtSecret, logger)

	// Data Layer
	userStorage := pqstorage.NewUserStorage(pqDB)
	userActivationStorage := pqstorage.NewUserActivationStorage(pqDB)
	passwordRecoveryStorage := pqstorage.NewPasswordRecoveryStorage(pqDB)
	mailStorage := pqstorage.NewMailStorage(pqDB)
	authTokenBlackListStorage := redisstorage.NewAuthTokenBlackListStorage()

	// Mailers
	userActivationMailer := mail.NewUserActivationMailer(mailStorage)
	passwordRecoveryMailer := mail.NewPasswordRecoveryMailer(mailStorage)

	// App Layer
	userActivationService := useractivation.NewService(config.userActivationURL, userActivationStorage, userStorage, userActivationMailer)
	passwordRecoveryService := passwordrecovery.NewService(config.userPasswordRecoveryURL, userStorage, passwordRecoveryStorage, passwordRecoveryMailer, hashService)
	authTokenService := authtoken.NewService(jwtService)
	signupService := signup.NewService(hashService, userStorage, userActivationService)
	loginService := login.NewService(userStorage, hashService, authTokenService)
	logoutService := logout.NewService(jwtService, authTokenBlackListStorage, logger)

	// Transport Layer
	signupHandler := signupTransport.NewHandler(signupService, logger, validator)
	userActivateHandler := useractivateTransport.NewHandler(userActivationService, logger, validator)
	loginHandler := loginTransport.NewHandler(loginService, logger, validator)
	logoutHandler := logoutTransport.NewHandler(logoutService, logger, validator)
	refreshHandler := refreshtokenTransport.NewHandler(authTokenService, logger)
	forgotPasswordHandler := forgotpasswordTransport.NewHandler(passwordRecoveryService, logger, validator)
	recoverPasswordHandler := recoverpasswordTransport.NewHandler(passwordRecoveryService, logger, validator)

	router := driverGin.NewRouter()
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
