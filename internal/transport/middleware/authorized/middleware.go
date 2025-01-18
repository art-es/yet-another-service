//go:generate mockgen -source=middleware.go -destination=mock/middleware.go -package=mock
package authorized

import (
	"context"
	"strings"

	"github.com/art-es/yet-another-service/internal/app/shared/errors"
	contextcore "github.com/art-es/yet-another-service/internal/core/context"
	"github.com/art-es/yet-another-service/internal/core/http"
	httputil "github.com/art-es/yet-another-service/internal/core/http/util"
	"github.com/art-es/yet-another-service/internal/core/log"
)

const headerPrefix = "bearer "

type authService interface {
	Authorize(ctx context.Context, accessToken string) (string, error)
}

type Middleware struct {
	authService authService
	logger      log.Logger
}

func NewMiddleware(
	authService authService,
	logger log.Logger,
) *Middleware {
	return &Middleware{
		authService: authService,
		logger:      logger,
	}
}

func (m *Middleware) Wrap(handle http.Handler) http.Handler {
	return func(ctx http.Context) {
		authHeader := ctx.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), headerPrefix) {
			httputil.RespondUnauthorized(ctx)
			return
		}

		accessToken := authHeader[len(headerPrefix):]

		userID, err := m.authService.Authorize(ctx, accessToken)
		if err != nil {
			if err == errors.ErrInvalidAuthToken {
				httputil.RespondUnauthorized(ctx)
				return
			}

			m.logger.Error().Err(err).Msg("authorize error")
			httputil.RespondInternalError(ctx)
			return
		}

		ctx = ctx.With(contextcore.WithUserID(ctx, userID))
		handle(ctx)
	}
}
