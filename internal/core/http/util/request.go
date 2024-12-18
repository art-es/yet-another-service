package util

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/art-es/yet-another-service/internal/core/http"
)

var ErrInvalidRequestBody = errors.New("invalid request body")

func EnrichRequestBody(ctx http.Context, out any) error {
	req := ctx.Request()
	if req == nil || req.Body == nil {
		return ErrInvalidRequestBody
	}

	if err := json.NewDecoder(req.Body).Decode(out); err != nil {
		return ErrInvalidRequestBody
	}

	return nil
}

func GetAuthorizationToken(ctx http.Context) (string, bool) {
	req := ctx.Request()
	if req == nil || req.Header == nil {
		return "", false
	}

	s := req.Header.Get("Authorization")
	if !strings.HasPrefix(strings.ToLower(s), "bearer ") {
		return "", false
	}

	s = strings.TrimSpace(s[len("bearer "):])
	return s, len(s) > 0
}
