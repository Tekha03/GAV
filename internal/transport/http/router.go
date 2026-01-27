package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"gav/internal/transport/http/handlers"
	"gav/internal/transport/http/middleware"
)

func NewHandler(
	authH *handlers.AuthHandler,
	userH *handlers.UserHandler,
	postH *handlers.PostHandler,
	authMW func(http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logging)
	r.Use(middleware.Recover)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("auth/register", authH.Register)
		r.Post("auth/login", authH.Login)

		r.Group(func(r chi.Router) {
			r.Use(authMW)

			r.Get("auth/me", userH.Me)

			r.Route("/posts", func(r chi.Router) {
				r.Post("/", postH.Create)
				r.Get("/{id}", postH.GetByID)
				r.Delete("/{id}", postH.Delete)
			})
		})
	})

	return r
}
