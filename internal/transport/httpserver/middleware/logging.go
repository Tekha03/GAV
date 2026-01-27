package middleware

import (
	"log"
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

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			status: http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		log.Printf(
			"%s %s %d %s",
			r.Method,
			r.URL.Path,
			wrapped.status,
			time.Since(start),
		)
	})
}
