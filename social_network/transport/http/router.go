package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"social_network/internal/post"
	"social_network/transport/http/handlers"
	"social_network/transport/http/middleware"
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
	Upload		 *handlers.UploadHandler
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

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logging(logger))
	r.Use(middleware.Recover(logger))
	r.Use(middleware.Timeout)

	r.Get("/", func (w http.ResponseWriter, r *http.Request)  {
		w.Header().Set("Content-Type", "text/html")
	html := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>GAV API</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 40px auto; line-height: 1.6; }
        h1 { color: #2c3e50; }
        pre { background: #f8f9fa; padding: 15px; border-radius: 8px; }
        code { background: #e9ecef; padding: 2px 6px; border-radius: 4px; }
    </style>
</head>
<body>
    <h1>GAV API запущен!</h1>
    <p>Все эндпоинты находятся по адресу: <code>/api/v1</code></p>

    <h3>Примеры использования:</h3>
    <pre>
POST /api/v1/auth/register   - регистрация
POST /api/v1/auth/login      - вход
GET  /api/v1/feed            - лента постов (требуется токен)
POST /api/v1/posts           - создать пост
    </pre>

    <p>Swagger UI доступен по адресу: <a href="/swagger/index.html">/swagger/index.html</a> (если подключён)</p>
</body>
</html>
    `
		w.Write([]byte(html))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// r.Route("/admin", func(r chi.Router) {
		// r.Use(middleware.RequireRole("admin"))
		// r.Get("/users", h.Admin.ListUsers)
		// r.Post("/stats/reset", h.Admin.ResetStats)
		// })

		r.Group(func(r chi.Router) {
			r.Use(deps.AuthMW)

			// ---- Auth ----
			r.Route("/auth", func(r chi.Router) {
				r.Post("/register", h.Auth.Register)
				r.Post("/login", h.Auth.Login)

				r.Group(func(r chi.Router) {
					r.Use(deps.AuthMW)

					r.Get("/me", h.Auth.Me)
					r.Post("/refresh", h.Auth.Refresh)
					r.Post("/logout", h.Auth.Logout)
				})
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

			// ---- Upload ----
			r.Route("/upload", func(r chi.Router) {
				r.Post("/avatar", h.Upload.UploadAvatar)
				r.Post("/post-image", h.Upload.UploadPostImage)
			})
		})
	})

	return r
}
