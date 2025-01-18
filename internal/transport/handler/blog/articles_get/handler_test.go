package articles_get

import (
	_ "embed"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/art-es/yet-another-service/internal/app/shared/dto"
	"github.com/art-es/yet-another-service/internal/core/pointer"
	"github.com/art-es/yet-another-service/internal/testutil"
	"github.com/art-es/yet-another-service/internal/transport/handler/blog/articles_get/mock"
)

var (
	//go:embed testdata/app_error.json
	expectedBodyAppError []byte

	//go:embed testdata/ok.json
	expectedBodyOK []byte
)

func TestHandler(t *testing.T) {
	for _, tt := range []struct {
		name   string
		setup  func(req *http.Request, blog *mock.MockblogService)
		assert func(t *testing.T, res *httptest.ResponseRecorder, logs []string)
	}{
		{
			name: "app error",
			setup: func(req *http.Request, blog *mock.MockblogService) {
				blog.EXPECT().
					GetArticles(gomock.Any(), gomock.Eq(&dto.GetArticlesIn{})).
					Return(nil, errors.New("dummy error"))
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusInternalServerError, res.Code)
				assert.JSONEq(t, string(expectedBodyAppError), res.Body.String())
				assert.Len(t, logs, 1)
				assert.Equal(t, `{"level":"error","error":"dummy error","message":"get articles error on blog service"}`, logs[0])
			},
		},
		{
			name: "ok",
			setup: func(req *http.Request, blog *mock.MockblogService) {
				query := url.Values{}
				query.Set("fromSlug", "foo")
				req.URL.RawQuery = query.Encode()

				blog.EXPECT().
					GetArticles(gomock.Any(), gomock.Eq(&dto.GetArticlesIn{
						FromSlug: pointer.To("foo"),
					})).
					Return(
						&dto.GetArticlesOut{
							Articles: []dto.Article{
								{
									Slug:    "bar",
									Title:   "Bar Title",
									Content: "Bar Content",
									Author: &dto.ArticleAuthor{
										DisplayName: "Bob",
										NickName:    "bob123",
									},
								},
								{
									Slug:    "baz",
									Title:   "Baz Title",
									Content: "Baz Content",
									Author:  nil,
								},
							},
							HasMore: true,
						},
						nil,
					)
			},
			assert: func(t *testing.T, res *httptest.ResponseRecorder, logs []string) {
				assert.Equal(t, http.StatusOK, res.Code)
				assert.JSONEq(t, string(expectedBodyOK), res.Body.String())
				assert.Empty(t, logs)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			blog := mock.NewMockblogService(ctrl)
			logger := testutil.NewLogger()
			ctx, req, res := testutil.NewHTTPContext(ctrl)

			tt.setup(req, blog)

			handler := NewHandler(blog, logger)
			handler.Handle(ctx)

			tt.assert(t, res, logger.Logs())
		})
	}
}
