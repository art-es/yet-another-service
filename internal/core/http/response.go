package http

import (
	"encoding/json"
	"net/http"
)

type errorResponseBody struct {
	Message string `json:"message,omitempty"`
}

func Respond(ctx Context, code int, body any) {
	w := ctx.ResponseWriter()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

func RespondBadRequest(ctx Context, msg string) {
	Respond(ctx, http.StatusBadRequest, errorResponseBody{Message: msg})
}

func RespondUnauthorized(ctx Context) {
	Respond(ctx, http.StatusUnauthorized, errorResponseBody{Message: "Unauthorized."})
}

func RespondNotFound(ctx Context) {
	Respond(ctx, http.StatusNotFound, errorResponseBody{Message: "Not found."})
}

func RespondInternalError(ctx Context) {
	Respond(ctx, http.StatusInternalServerError, errorResponseBody{
		Message: "An unexpected error occurred. Please try again later.",
	})
}
