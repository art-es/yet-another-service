package activate

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
	mockvalidation "github.com/art-es/yet-another-service/internal/core/validation/mock"
	"github.com/art-es/yet-another-service/internal/domain/auth"
	"github.com/art-es/yet-another-service/internal/driver/zerolog"
	"github.com/art-es/yet-another-service/internal/transport/handler/auth/activate/mock"
)

func TestHandler(t *testing.T) {
	const token = "18d440f5-2664-42b1-bfaa-1c15f1687885"

	for _, tt := range []struct {
		name   string
		setup  func(authSvc *mock.MockauthService, validator *mockvalidation.MockValidator, req *http.Request)
		assert func(t *testing.T, res *httptest.ResponseRecorder, logs []string)
	}{
		{
			name: "validation error",
			setup: func(authSvc *mock.MockauthService, validator *mockvalidation.MockValidator, req *http.Request) {
				req.URL.RawQuery = "token=foo"

				validator.EXPECT().
					Var(gomock.Eq("foo"), gomock.Eq("required,uuid")).
					Return(errors.New("dummy validation error"))
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusNotFound, res.Code)
				expResBody := `{"message": "Not found."}`
				assert.JSONEq(t, expResBody, res.Body.String())

				assert.Len(t, logs, 0)
			},
		},
		{
			name: "activation not found",
			setup: func(authSvc *mock.MockauthService, validator *mockvalidation.MockValidator, req *http.Request) {
				req.URL.RawQuery = "token=" + token

				validator.EXPECT().
					Var(gomock.Eq(token), gomock.Eq("required,uuid")).
					Return(nil)

				authSvc.EXPECT().
					Activate(gomock.Any(), gomock.Eq(token)).
					Return(auth.ErrActivationNotFound)
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusNotFound, res.Code)
				expResBody := `{"message": "Not found."}`
				assert.JSONEq(t, expResBody, res.Body.String())

				assert.Len(t, logs, 0)
			},
		},
		{
			name: "auth service error",
			setup: func(authSvc *mock.MockauthService, validator *mockvalidation.MockValidator, req *http.Request) {
				req.URL.RawQuery = "token=" + token

				validator.EXPECT().
					Var(gomock.Eq(token), gomock.Eq("required,uuid")).
					Return(nil)

				authSvc.EXPECT().
					Activate(gomock.Any(), gomock.Eq(token)).
					Return(errors.New("auth service dummy error"))
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusInternalServerError, res.Code)
				expResBody := `{"message": "An unexpected error occurred. Please try again later."}`
				assert.JSONEq(t, expResBody, res.Body.String())

				assert.Len(t, logs, 1)
				expErrorLog := `{"level":"error", "error":"auth service dummy error", "message":"activate error on auth service"}`
				assert.JSONEq(t, expErrorLog, logs[0])
			},
		},
		{
			name: "ok",
			setup: func(authSvc *mock.MockauthService, validator *mockvalidation.MockValidator, req *http.Request) {
				req.URL.RawQuery = "token=" + token

				validator.EXPECT().
					Var(gomock.Eq(token), gomock.Eq("required,uuid")).
					Return(nil)

				authSvc.EXPECT().
					Activate(gomock.Any(), gomock.Eq(token)).
					Return(nil)
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.JSONEq(t, "{}", res.Body.String())

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
			validator := mockvalidation.NewMockValidator(ctrl)

			ctx := mockhttp.NewMockContext(ctrl)
			req := httptest.NewRequest(http.MethodGet, "/activate", nil)
			res := httptest.NewRecorder()
			ctx.EXPECT().Request().Return(req).AnyTimes()
			ctx.EXPECT().ResponseWriter().Return(res).AnyTimes()

			if tt.setup != nil {
				tt.setup(authSvc, validator, req)
			}

			handler := NewHandler(authSvc, logger, validator)
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
