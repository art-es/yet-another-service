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
	"github.com/art-es/yet-another-service/internal/driver/bcrypt"
	"github.com/art-es/yet-another-service/internal/driver/gin"
	"github.com/art-es/yet-another-service/internal/driver/jwt"
	"github.com/art-es/yet-another-service/internal/driver/postgres"
	"github.com/art-es/yet-another-service/internal/driver/redis"
	validatord "github.com/art-es/yet-another-service/internal/driver/validator"
	"github.com/art-es/yet-another-service/internal/driver/zerolog"
	pqstorage "github.com/art-es/yet-another-service/internal/storage/postgres"
	rdstorage "github.com/art-es/yet-another-service/internal/storage/redis"
	useractivatetp "github.com/art-es/yet-another-service/internal/transport/handler/auth/activate"
	forgotpasswordtp "github.com/art-es/yet-another-service/internal/transport/handler/auth/forgot_password"
	logintp "github.com/art-es/yet-another-service/internal/transport/handler/auth/login"
	logouttp "github.com/art-es/yet-another-service/internal/transport/handler/auth/logout"
	recoverpasswordtp "github.com/art-es/yet-another-service/internal/transport/handler/auth/recover_password"
	refreshtokentp "github.com/art-es/yet-another-service/internal/transport/handler/auth/refresh"
	signuptp "github.com/art-es/yet-another-service/internal/transport/handler/auth/signup"
)

func main() {
	logger := zerolog.NewLogger()
	config := getAppConfig(logger)

	// Drivers
	pqDB := postgres.Connect(config.postgresURL)
	rdDB := redis.Connect(config.redisAddr)
	validator := validatord.New()
	hashService := bcrypt.NewHashService()
	jwtService := jwt.NewService(config.jwtSecret, logger)

	// Data Layer
	userStorage := pqstorage.NewUserStorage(pqDB)
	userActivationStorage := pqstorage.NewUserActivationStorage(pqDB)
	passwordRecoveryStorage := pqstorage.NewPasswordRecoveryStorage(pqDB)
	mailStorage := pqstorage.NewMailStorage(pqDB)
	authTokenBlackListStorage := rdstorage.NewAuthTokenBlackListStorage(rdDB)

	// Mailers
	userActivationMailer := mail.NewUserActivationMailer(mailStorage)
	passwordRecoveryMailer := mail.NewPasswordRecoveryMailer(mailStorage)

	// App Layer
	userActivationService := useractivation.NewService(config.userActivationURL, userActivationStorage, userStorage, userActivationMailer)
	passwordRecoveryService := passwordrecovery.NewService(config.userPasswordRecoveryURL, userStorage, passwordRecoveryStorage, passwordRecoveryMailer, hashService)
	authTokenService := authtoken.NewService(jwtService, authTokenBlackListStorage)
	signupService := signup.NewService(hashService, userStorage, userActivationService)
	loginService := login.NewService(userStorage, hashService, authTokenService)
	logoutService := logout.NewService(authTokenService, logger)

	// Transport Layer
	signupHandler := signuptp.NewHandler(signupService, logger, validator)
	userActivateHandler := useractivatetp.NewHandler(userActivationService, logger, validator)
	loginHandler := logintp.NewHandler(loginService, logger, validator)
	logoutHandler := logouttp.NewHandler(logoutService, logger, validator)
	refreshHandler := refreshtokentp.NewHandler(authTokenService, logger)
	forgotPasswordHandler := forgotpasswordtp.NewHandler(passwordRecoveryService, logger, validator)
	recoverPasswordHandler := recoverpasswordtp.NewHandler(passwordRecoveryService, logger, validator)

	router := gin.NewRouter()
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
