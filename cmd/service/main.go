package main

import (
	"net/http"

	domain_activate "github.com/art-es/yet-another-service/internal/domain/auth/activate"
	domain_login "github.com/art-es/yet-another-service/internal/domain/auth/login"
	domain_logout "github.com/art-es/yet-another-service/internal/domain/auth/logout"
	domain_refresh "github.com/art-es/yet-another-service/internal/domain/auth/refresh"
	domain_signup "github.com/art-es/yet-another-service/internal/domain/auth/signup"
	driver_bcrypt "github.com/art-es/yet-another-service/internal/driver/bcrypt"
	driver_gin "github.com/art-es/yet-another-service/internal/driver/gin"
	driver_jwt "github.com/art-es/yet-another-service/internal/driver/jwt"
	driver_postgres "github.com/art-es/yet-another-service/internal/driver/postgres"
	driver_validator "github.com/art-es/yet-another-service/internal/driver/validator"
	driver_zerolog "github.com/art-es/yet-another-service/internal/driver/zerolog"
	storage_postgres "github.com/art-es/yet-another-service/internal/storage/postgres"
	storage_redis "github.com/art-es/yet-another-service/internal/storage/redis"
	transport_activate "github.com/art-es/yet-another-service/internal/transport/handler/auth/activate"
	transport_login "github.com/art-es/yet-another-service/internal/transport/handler/auth/login"
	transport_logout "github.com/art-es/yet-another-service/internal/transport/handler/auth/logout"
	transport_refresh "github.com/art-es/yet-another-service/internal/transport/handler/auth/refresh"
	transport_signup "github.com/art-es/yet-another-service/internal/transport/handler/auth/signup"
)

func main() {
	// Drivers
	logger := driver_zerolog.NewLogger()
	validator := driver_validator.New()
	hashService := driver_bcrypt.NewHashService()
	pgDB := driver_postgres.Connect("")
	jwtService := driver_jwt.NewService("")

	// Data Layer
	userStorage := storage_postgres.NewUserStorage(pgDB)
	userActivationStorage := storage_postgres.NewUserActivationStorage(pgDB)
	authTokenBlackListStorage := storage_redis.NewAuthTokenBlackListStorage()

	// App Layer
	signupService := domain_signup.NewService(hashService, userStorage, userActivationStorage)
	activateService := domain_activate.NewService(userActivationStorage, userStorage)
	loginService := domain_login.NewService(userStorage, hashService, jwtService)
	logoutService := domain_logout.NewService(jwtService, authTokenBlackListStorage, logger)
	refreshService := domain_refresh.NewService(jwtService)

	// Transport Layer
	signupHandler := transport_signup.NewHandler(signupService, logger, validator)
	activateHandler := transport_activate.NewHandler(activateService, logger, validator)
	loginHandler := transport_login.NewHandler(loginService, logger, validator)
	logoutHandler := transport_logout.NewHandler(logoutService, logger, validator)
	refreshHandler := transport_refresh.NewHandler(refreshService, logger)

	router := driver_gin.NewRouter()
	router.Register(http.MethodPost, "/auth/signup", signupHandler.Handle)
	router.Register(http.MethodGet, "/auth/activate", activateHandler.Handle)
	router.Register(http.MethodPost, "/auth/login", loginHandler.Handle)
	router.Register(http.MethodPost, "/auth/logout", logoutHandler.Handle)
	router.Register(http.MethodPost, "/auth/refresh", refreshHandler.Handle)
	router.Register(http.MethodPost, "/auth/forgot-password", nil)
	router.Register(http.MethodPost, "/auth/recover-password", nil)

	if err := router.Run(); err != nil {
		logger.Panic().Err(err).Msg("router run error")
	}
}
