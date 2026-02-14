package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"gav/internal/transport/response"
)

func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
	return  func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				if rec := recover(); rec != nil {
					var err error
					switch e := rec.(type) {
					case error:
						err = e
					default:
						err = fmt.Errorf("panic: %v", e)
					}

					logger.Error(
						"panic recovered",
						slog.Any("error", err),
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
						slog.String("stack", string(debug.Stack())),
					)

					response.InternalError(w)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
