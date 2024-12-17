package authorized

import (
	"context"
	"strings"

	"github.com/art-es/yet-another-service/internal/app/shared/errors"
	"github.com/art-es/yet-another-service/internal/core/http"
	"github.com/art-es/yet-another-service/internal/core/log"
	contextutil "github.com/art-es/yet-another-service/internal/util/context"
)

const headerPrefix = "bearer "

type authService interface {
	Authorize(ctx context.Context, accessToken string) (string, error)
}

type Middleware struct {
	authService authService
	logger      log.Logger
}

func (m *Middleware) Wrap(handle http.Handler) http.Handler {
	return func(ctx http.Context) {
		authHeader := ctx.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), headerPrefix) {
			http.RespondUnauthorized(ctx)
			return
		}

		accessToken := authHeader[len(headerPrefix):]

		userID, err := m.authService.Authorize(ctx, accessToken)
		if err != nil {
			if err == errors.ErrInvalidAuthToken {
				http.RespondUnauthorized(ctx)
			}

			m.logger.Error().Err(err).Msg("authorize error")
			http.RespondInternalError(ctx)
			return
		}

		ctx = ctx.With(contextutil.WithUserID(ctx, userID))
		handle(ctx)
	}
}
