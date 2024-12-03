package refresh

import (
	"bufio"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mockhttp "github.com/art-es/yet-another-service/internal/core/http/mock"
	"github.com/art-es/yet-another-service/internal/domain/auth"
	"github.com/art-es/yet-another-service/internal/driver/zerolog"
	"github.com/art-es/yet-another-service/internal/transport/handler/auth/refresh/mock"
)

func TestHandler(t *testing.T) {
	for _, tt := range []struct {
		name   string
		setup  func(authSvc *mock.MockauthService, req *http.Request)
		assert func(t *testing.T, res *httptest.ResponseRecorder, logs []string)
	}{
		{
			name: "no token",
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusUnauthorized, res.Code)
				expResBody := `{"message": "Unauthorized."}`
				assert.JSONEq(t, expResBody, res.Body.String())

				assert.Len(t, logs, 0)
			},
		},
		{
			name: "invalid token",
			setup: func(authSvc *mock.MockauthService, req *http.Request) {
				req.Header.Set("Authorization", "Bearer dummy refresh token")

				authSvc.EXPECT().
					Refresh(gomock.Eq("dummy refresh token")).
					Return("", auth.ErrInvalidToken)
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusUnauthorized, res.Code)
				expResBody := `{"message": "Unauthorized."}`
				assert.JSONEq(t, expResBody, res.Body.String())

				assert.Len(t, logs, 0)
			},
		},
		{
			name: "auth service error",
			setup: func(authSvc *mock.MockauthService, req *http.Request) {
				req.Header.Set("Authorization", "Bearer dummy refresh token")

				authSvc.EXPECT().
					Refresh(gomock.Eq("dummy refresh token")).
					Return("", errors.New("auth service dummy error"))
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusInternalServerError, res.Code)
				expResBody := `{"message": "An unexpected error occurred. Please try again later."}`
				assert.JSONEq(t, expResBody, res.Body.String())

				assert.Len(t, logs, 1)
				expErrorLog := `{"level":"error", "error":"auth service dummy error", "message":"refresh error on auth service"}`
				assert.JSONEq(t, expErrorLog, logs[0])
			},
		},
		{
			name: "ok",
			setup: func(authSvc *mock.MockauthService, req *http.Request) {
				req.Header.Set("Authorization", "Bearer dummy refresh token")

				authSvc.EXPECT().
					Refresh(gomock.Eq("dummy refresh token")).
					Return("dummy access token", nil)
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusOK, res.Code)
				expResBody := `{"accessToken": "dummy access token"}`
				assert.JSONEq(t, expResBody, res.Body.String())

				assert.Len(t, logs, 0)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logbuf := &bytes.Buffer{}

			authSvc := mock.NewMockauthService(ctrl)
			logger := zerolog.NewLoggerWithWriter(logbuf)

			ctx := mockhttp.NewMockContext(ctrl)
			req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
			res := httptest.NewRecorder()
			ctx.EXPECT().Request().Return(req).AnyTimes()
			ctx.EXPECT().ResponseWriter().Return(res).AnyTimes()

			if tt.setup != nil {
				tt.setup(authSvc, req)
			}

			handler := NewHandler(authSvc, logger)
			handler.Handle(ctx)

			var logs []string
			for s := bufio.NewScanner(logbuf); s.Scan(); {
				logs = append(logs, s.Text())
			}

			if tt.assert != nil {
				tt.assert(t, res, logs)
			}
		})
	}
}
