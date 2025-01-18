package signup

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/yet-another-service/internal/app/shared/dto"
	mockhttp "github.com/art-es/yet-another-service/internal/core/http/mock"
	mockvalidation "github.com/art-es/yet-another-service/internal/core/validation/mock"
	"github.com/art-es/yet-another-service/internal/driver/zerolog"
	"github.com/art-es/yet-another-service/internal/transport/handler/auth/signup/mock"
)

func TestHandler(t *testing.T) {
	for _, tt := range []struct {
		name   string
		setup  func(authSvc *mock.MockauthService, validator *mockvalidation.MockValidator, req *http.Request)
		assert func(t *testing.T, res *httptest.ResponseRecorder, logs []string)
	}{
		{
			name: "no request body",
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				expResBody := `{"message": "invalid request body"}`
				assert.JSONEq(t, expResBody, res.Body.String())

				assert.Len(t, logs, 0)
			},
		},
		{
			name: "invalid content type",
			setup: func(authSvc *mock.MockauthService, validator *mockvalidation.MockValidator, req *http.Request) {
				req.Body = io.NopCloser(strings.NewReader(`foo`))
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				expResBody := `{"message": "invalid request body"}`
				assert.JSONEq(t, expResBody, res.Body.String())

				assert.Len(t, logs, 0)
			},
		},
		{
			name: "validation error",
			setup: func(authSvc *mock.MockauthService, validator *mockvalidation.MockValidator, req *http.Request) {
				reqBody := `{"name": "dummyName", "email": "dummy@example.com", "password": "dummy123"}`
				req.Body = io.NopCloser(strings.NewReader(reqBody))

				expParsedReq := &request{DisplayName: "dummyName", Email: "dummy@example.com", Password: "dummy123"}
				validator.EXPECT().
					Struct(gomock.Eq(expParsedReq)).
					Return(errors.New("dummy validation error"))
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusBadRequest, res.Code)
				expResBody := `{"message": "dummy validation error"}`
				assert.JSONEq(t, expResBody, res.Body.String())

				assert.Len(t, logs, 0)
			},
		},
		{
			name: "auth service error",
			setup: func(authSvc *mock.MockauthService, validator *mockvalidation.MockValidator, req *http.Request) {
				reqBody := `{"name": "dummyName", "email": "dummy@example.com", "password": "dummy123"}`
				req.Body = io.NopCloser(strings.NewReader(reqBody))

				expParsedReq := &request{DisplayName: "dummyName", Email: "dummy@example.com", Password: "dummy123"}
				validator.EXPECT().
					Struct(gomock.Eq(expParsedReq)).
					Return(nil)

				expAuthReq := &dto.SignupIn{DisplayName: "dummyName", Email: "dummy@example.com", Password: "dummy123"}
				authSvc.EXPECT().
					Signup(gomock.Any(), gomock.Eq(expAuthReq)).
					Return(errors.New("auth service dummy error"))
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusInternalServerError, res.Code)
				expResBody := `{"message": "An unexpected error occurred. Please try again later."}`
				assert.JSONEq(t, expResBody, res.Body.String())

				assert.Len(t, logs, 1)
				expErrorLog := `{"level":"error", "error":"auth service dummy error", "message":"signup error on auth service"}`
				assert.JSONEq(t, expErrorLog, logs[0])
			},
		},
		{
			name: "ok",
			setup: func(authSvc *mock.MockauthService, validator *mockvalidation.MockValidator, req *http.Request) {
				reqBody := `{"name": "dummyName", "email": "dummy@example.com", "password": "dummy123"}`
				req.Body = io.NopCloser(strings.NewReader(reqBody))

				expParsedReq := &request{DisplayName: "dummyName", Email: "dummy@example.com", Password: "dummy123"}
				validator.EXPECT().
					Struct(gomock.Eq(expParsedReq)).
					Return(nil)

				expAuthReq := &dto.SignupIn{DisplayName: "dummyName", Email: "dummy@example.com", Password: "dummy123"}
				authSvc.EXPECT().
					Signup(gomock.Any(), gomock.Eq(expAuthReq)).
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
			req := httptest.NewRequest(http.MethodPost, "/signup", nil)
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
