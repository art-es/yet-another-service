package authorized

import (
	"context"
	_ "embed"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	apperrors "github.com/art-es/yet-another-service/internal/app/shared/errors"
	contextcore "github.com/art-es/yet-another-service/internal/core/context"
	corehttp "github.com/art-es/yet-another-service/internal/core/http"
	mockcorehttp "github.com/art-es/yet-another-service/internal/core/http/mock"
	corehttputil "github.com/art-es/yet-another-service/internal/core/http/util"
	"github.com/art-es/yet-another-service/internal/testutil"
	"github.com/art-es/yet-another-service/internal/transport/middleware/authorized/mock"
)

var (
	//go:embed testdata/unauthorized.json
	expectedUnauthorizedBody []byte

	//go:embed testdata/internal_error.json
	expectedInternalErrorBody []byte

	//go:embed testdata/ok.json
	expectedOKBody []byte
)

func TestMiddleware(t *testing.T) {
	for _, tt := range []struct {
		name   string
		setup  func(t *testing.T, ctx *mockcorehttp.MockContext, req *http.Request, authSvc *mock.MockauthService)
		assert func(t *testing.T, res *httptest.ResponseRecorder, logs []string)
	}{
		{
			name: "no auth header",
			setup: func(t *testing.T, ctx *mockcorehttp.MockContext, req *http.Request, authSvc *mock.MockauthService) {
				req.Header.Del("Authorization")
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusUnauthorized, res.Code)
				assert.JSONEq(t, string(expectedUnauthorizedBody), res.Body.String())
				assert.Empty(t, logs)
			},
		},
		{
			name: "no bearer prefix",
			setup: func(t *testing.T, ctx *mockcorehttp.MockContext, req *http.Request, authSvc *mock.MockauthService) {
				req.Header.Set("Authorization", "dummy token")
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusUnauthorized, res.Code)
				assert.JSONEq(t, string(expectedUnauthorizedBody), res.Body.String())
				assert.Empty(t, logs)
			},
		},
		{
			name: "invalid token",
			setup: func(t *testing.T, ctx *mockcorehttp.MockContext, req *http.Request, authSvc *mock.MockauthService) {
				req.Header.Set("Authorization", "bearer dummy token")

				authSvc.EXPECT().
					Authorize(gomock.Any(), gomock.Eq("dummy token")).
					Return("", apperrors.ErrInvalidAuthToken)
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusUnauthorized, res.Code)
				assert.JSONEq(t, string(expectedUnauthorizedBody), res.Body.String())
				assert.Empty(t, logs)
			},
		},
		{
			name: "internal error",
			setup: func(t *testing.T, ctx *mockcorehttp.MockContext, req *http.Request, authSvc *mock.MockauthService) {
				req.Header.Set("Authorization", "bearer dummy token")

				authSvc.EXPECT().
					Authorize(gomock.Any(), gomock.Eq("dummy token")).
					Return("", errors.New("dummy error"))
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusInternalServerError, res.Code)
				assert.JSONEq(t, string(expectedInternalErrorBody), res.Body.String())
				assert.Len(t, logs, 1)
				assert.Equal(t, `{"level":"error","error":"dummy error","message":"authorize error"}`, logs[0])
			},
		},
		{
			name: "ok",
			setup: func(t *testing.T, ctx *mockcorehttp.MockContext, req *http.Request, authSvc *mock.MockauthService) {
				req.Header.Set("Authorization", "bearer dummy token")

				authSvc.EXPECT().
					Authorize(gomock.Any(), gomock.Eq("dummy token")).
					Return("dummy user ID", nil)

				ctx.EXPECT().
					With(gomock.Any()).
					DoAndReturn(func(newCtx context.Context) corehttp.Context {
						userID, ok := contextcore.UserID(newCtx)
						assert.True(t, ok)
						assert.Equal(t, userID, "dummy user ID")

						return ctx
					})
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.JSONEq(t, string(expectedOKBody), res.Body.String())
				assert.Empty(t, logs)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx, req, res := testutil.NewHTTPContext(ctrl)
			authSvc := mock.NewMockauthService(ctrl)
			logger := testutil.NewLogger()
			tt.setup(t, ctx, req, authSvc)

			handle := NewMiddleware(authSvc, logger).Wrap(func(ctx corehttp.Context) {
				corehttputil.Respond(ctx, http.StatusOK, map[string]any{"message": "OK."})
			})

			handle(ctx)
			tt.assert(t, res, logger.Logs())
		})
	}
}
