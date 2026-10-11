package httpx

import (
	"log/slog"
	"net/http"
)

func LogError(message string, err error, r *http.Request) {
	slog.Error(
		message,
		"error", err,
		"endpoint", r.URL.Path,
		"method", r.Method,
	)
}
