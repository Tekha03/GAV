package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"gav/internal/post"
	"gav/transport/http/handlers"
	"gav/transport/http/middleware"
)

type Handlers struct {
	Auth         *handlers.AuthHandler
	User         *handlers.UserHandler
	Profile       *handlers.ProfileHandler
	Post         *handlers.PostHandler
	Feed		 *handlers.FeedHandler
	Comment      *handlers.CommentHandler
	Like         *handlers.LikeHandler
	Follow       *handlers.FollowHandler
	Dog          *handlers.DogHandler
	Vaccination  *handlers.VaccinationHandler
	Stats        *handlers.StatsHandler
	Settings     *handlers.SettingsHandler
}

type RouterDeps struct {
	AuthMW		func(http.Handler) http.Handler
	PostService post.PostService
}

func NewRouter(
	h Handlers,
	deps RouterDeps,
	logger *slog.Logger,
) http.Handler {

	r := chi.NewRouter()
	r.Use(deps.AuthMW)

	r.Use(middleware.RequireRole("admin"))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logging(logger))
	r.Use(middleware.Recover(logger))
	r.Use(middleware.Timeout)

	r.Route("/api/v1", func(r chi.Router) {
		// r.Route("/admin", func(r chi.Router) {
		// 	r.Get("/users", h.Admin.ListUsers)
		// 	r.Post("/stats/reset", h.Admin.ResetStats)
		// })

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.Auth.Register)
			r.Post("/login", h.Auth.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(deps.AuthMW)

			// ---- Auth ----
			r.Route("/auth", func(r chi.Router) {
				r.Get("/me", h.Auth.Me)
				r.Post("/refresh", h.Auth.Refresh)
				r.Post("/logout", h.Auth.Logout)
			})

			// ---- Users ----
			r.Route("/users", func(r chi.Router) {
				r.Get("/{id}", h.User.GetByID)
				r.Get("/{userID}/profile", h.Profile.GetByUserID)
				r.Put("/{id}", h.User.Update)
				r.Delete("/{id}", h.User.Delete)
			})

			// ---- Profiles ----
			r.Route("/profiles", func(r chi.Router) {
				r.Put("/{userID}", h.Profile.Update)
				r.Delete("/{userID}", h.Profile.Delete)
			})

			// ---- Posts ----
			r.Route("/posts", func(r chi.Router) {
				r.Post("/", h.Post.Create)
				r.Get("/{id}", h.Post.GetByID)
				r.Get("/", h.Post.ListByUser)

				r.With(middleware.RequirePostOwner(deps.PostService)).Delete("/{id}", h.Post.Delete)

				r.Route("/{id}/comments", func(r chi.Router) {
					r.Post("/", h.Comment.Create)
					r.Get("/", h.Comment.ListByPostID)
				})

				r.Route("/{id}/likes", func(r chi.Router) {
					r.Post("/", h.Like.Add)
					r.Delete("/", h.Like.Remove)
				})
			})

			// ---- Feed ----
			r.Get("/feed", h.Feed.GetFeed)

			// ---- Comments ----
			r.Route("/comments", func(r chi.Router) {
				r.Delete("/{id}", h.Comment.Delete)
			})

			// ---- Follows ----
			r.Route("/follows", func(r chi.Router) {
				r.Post("/{userID}", h.Follow.Follow)
				r.Delete("/{userID}", h.Follow.Unfollow)
				r.Get("/followers/{userID}", h.Follow.GetFollowers)
				r.Get("/following/{userID}", h.Follow.GetFollowing)
			})

			// ---- Dogs ----
			r.Route("/dogs", func(r chi.Router) {
				r.Post("/", h.Dog.Create)
				r.Get("/{id}", h.Dog.GetPrivate)
				r.Get("/{id}", h.Dog.GetPublic)
				r.Put("/{id}", h.Dog.Update)
				r.Delete("/{id}", h.Dog.Delete)

				r.Route("/{id}/vaccinations", func(r chi.Router) {
					r.Post("/", h.Vaccination.Create)
					r.Get("/", h.Vaccination.ListByDogID)
					r.Put("/{vaccinationID}", h.Vaccination.Update)
				})
			})

			// ---- Vaccinations ----
			r.Route("/vaccinations", func(r chi.Router) {
				r.Delete("/{id}", h.Vaccination.Delete)
			})

			// ---- Stats ----
			r.Route("/stats", func(r chi.Router) {
				r.Get("/user/{userID}", h.Stats.UserStats)
				r.Get("/post/{postID}", h.Stats.PostStats)
				r.Get("/profile/{userID}", h.Stats.ProfileStats)
			})

			// ---- Settings ----
			r.Route("/settings", func(r chi.Router) {
				r.Get("/", h.Settings.Get)
				r.Put("/", h.Settings.Update)
			})
		})
	})

	return r
}
