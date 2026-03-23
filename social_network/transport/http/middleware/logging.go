package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{
				ResponseWriter: w,
				status: 		http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			logRequest(logger, r, wrapped.status, duration)
		})
	}
}

func logRequest(
	logger *slog.Logger,
	r *http.Request,
	status int,
	duration time.Duration,
) {
	attrs := []any{
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Int("status", status),
		slog.Duration("duration", duration),
	}

	switch {
	case status >= 500:
		logger.Error("http request", attrs...)
	case status >= 400:
		logger.Warn("http request", attrs...)
	default:
		logger.Info("http request", attrs...)
	}
}
