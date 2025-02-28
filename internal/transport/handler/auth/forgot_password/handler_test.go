package forgotpassword

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

	mockhttp "github.com/art-es/yet-another-service/internal/core/http/mock"
	mockvalidation "github.com/art-es/yet-another-service/internal/core/validation/mock"
	"github.com/art-es/yet-another-service/internal/driver/zerolog"
	"github.com/art-es/yet-another-service/internal/transport/handler/auth/forgot_password/mock"
)

func TestHandler(t *testing.T) {
	for _, tt := range []struct {
		name   string
		setup  func(recSvc *mock.MockrecoveryService, validator *mockvalidation.MockValidator, req *http.Request)
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
			setup: func(recSvc *mock.MockrecoveryService, validator *mockvalidation.MockValidator, req *http.Request) {
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
			setup: func(recSvc *mock.MockrecoveryService, validator *mockvalidation.MockValidator, req *http.Request) {
				reqBody := `{"email": "dummy@example.com"}`
				req.Body = io.NopCloser(strings.NewReader(reqBody))

				expParsedReq := &request{Email: "dummy@example.com"}
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
			setup: func(recSvc *mock.MockrecoveryService, validator *mockvalidation.MockValidator, req *http.Request) {
				reqBody := `{"email": "dummy@example.com"}`
				req.Body = io.NopCloser(strings.NewReader(reqBody))

				expParsedReq := &request{Email: "dummy@example.com"}
				validator.EXPECT().
					Struct(gomock.Eq(expParsedReq)).
					Return(nil)

				recSvc.EXPECT().
					Create(gomock.Any(), gomock.Eq("dummy@example.com")).
					Return(errors.New("auth service dummy error"))
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusInternalServerError, res.Code)
				expResBody := `{"message": "An unexpected error occurred. Please try again later."}`
				assert.JSONEq(t, expResBody, res.Body.String())

				assert.Len(t, logs, 1)
				expErrorLog := `{"level":"error", "error":"auth service dummy error", "message":"create recovery error on auth service"}`
				assert.JSONEq(t, expErrorLog, logs[0])
			},
		},
		{
			name: "ok",
			setup: func(recSvc *mock.MockrecoveryService, validator *mockvalidation.MockValidator, req *http.Request) {
				reqBody := `{"email": "dummy@example.com"}`
				req.Body = io.NopCloser(strings.NewReader(reqBody))

				expParsedReq := &request{Email: "dummy@example.com"}
				validator.EXPECT().
					Struct(gomock.Eq(expParsedReq)).
					Return(nil)

				recSvc.EXPECT().
					Create(gomock.Any(), gomock.Eq("dummy@example.com")).
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

			recSvc := mock.NewMockrecoveryService(ctrl)
			logger := zerolog.NewLoggerWithWriter(logbuf)
			validator := mockvalidation.NewMockValidator(ctrl)

			ctx := mockhttp.NewMockContext(ctrl)
			req := httptest.NewRequest(http.MethodPost, "/forgot-password", nil)
			res := httptest.NewRecorder()
			ctx.EXPECT().Request().Return(req).AnyTimes()
			ctx.EXPECT().ResponseWriter().Return(res).AnyTimes()

			if tt.setup != nil {
				tt.setup(recSvc, validator, req)
			}

			handler := NewHandler(recSvc, logger, validator)
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
