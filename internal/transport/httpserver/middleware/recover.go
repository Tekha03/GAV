package middleware

import (
	"log"
	"net/http"

	"gav/internal/transport/response"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				response.Error(
					w,
					http.StatusInternalServerError,
					err.(error),
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
