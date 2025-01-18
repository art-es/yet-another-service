package testutil

import (
	"net/http"
	"net/http/httptest"

	"go.uber.org/mock/gomock"

	mockhttp "github.com/art-es/yet-another-service/internal/core/http/mock"
)

func NewHTTPContext(ctrl *gomock.Controller) (*mockhttp.MockContext, *http.Request, *httptest.ResponseRecorder) {
	ctx := mockhttp.NewMockContext(ctrl)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	ctx.EXPECT().Request().Return(req).AnyTimes()
	ctx.EXPECT().ResponseWriter().Return(res).AnyTimes()
	return ctx, req, res
}
