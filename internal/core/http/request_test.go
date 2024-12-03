package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mockhttp "github.com/art-es/yet-another-service/internal/core/http/mock"
)

func TestGetAuthorizationToken(t *testing.T) {
	for _, tt := range []struct {
		name   string
		setup  func(ctx *mockhttp.MockContext)
		assert func(t *testing.T, res string, ok bool)
	}{
		{
			name: "no request",
			setup: func(ctx *mockhttp.MockContext) {
				ctx.EXPECT().Request().Return(nil)
			},
			assert: func(t *testing.T, res string, ok bool) {
				assert.False(t, ok)
				assert.Empty(t, res)
			},
		},
		{
			name: "no header",
			setup: func(ctx *mockhttp.MockContext) {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				ctx.EXPECT().Request().Return(req)
			},
			assert: func(t *testing.T, res string, ok bool) {
				assert.False(t, ok)
				assert.Empty(t, res)
			},
		},
		{
			name: "empty string in header",
			setup: func(ctx *mockhttp.MockContext) {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "")
				ctx.EXPECT().Request().Return(req)
			},
			assert: func(t *testing.T, res string, ok bool) {
				assert.False(t, ok)
				assert.Empty(t, res)
			},
		},
		{
			name: "invalid string in header",
			setup: func(ctx *mockhttp.MockContext) {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "foo")
				ctx.EXPECT().Request().Return(req)
			},
			assert: func(t *testing.T, res string, ok bool) {
				assert.False(t, ok)
				assert.Empty(t, res)
			},
		},
		{
			name: "empty token in header",
			setup: func(ctx *mockhttp.MockContext) {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer  ")
				ctx.EXPECT().Request().Return(req)
			},
			assert: func(t *testing.T, res string, ok bool) {
				assert.False(t, ok)
				assert.Empty(t, res)
			},
		},
		{
			name: "ok",
			setup: func(ctx *mockhttp.MockContext) {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer  Foo ")
				ctx.EXPECT().Request().Return(req)
			},
			assert: func(t *testing.T, res string, ok bool) {
				assert.True(t, ok)
				assert.Equal(t, "Foo", res)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := mockhttp.NewMockContext(ctrl)
			tt.setup(ctx)

			res, ok := GetAuthorizationToken(ctx)
			tt.assert(t, res, ok)
		})
	}
}
