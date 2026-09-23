package http

import (
	"log/slog"
	"net/http"
	"time"

	"shared/ratelimit"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	"social_network/internal/post"
	"social_network/transport/http/handlers"
	"social_network/transport/http/middleware"
)

type Handlers struct {
	Auth        *handlers.AuthHandler
	User        *handlers.UserHandler
	Walk        *handlers.WalkHandler
	Profile     *handlers.ProfileHandler
	Post        *handlers.PostHandler
	Feed        *handlers.FeedHandler
	Comment     *handlers.CommentHandler
	Like        *handlers.LikeHandler
	Follow      *handlers.FollowHandler
	Dog         *handlers.DogHandler
	Vaccination *handlers.VaccinationHandler
	Stats       *handlers.StatsHandler
	Settings    *handlers.SettingsHandler
	Upload      *handlers.UploadHandler
	WS          *handlers.NotificationHandler
	Device      *handlers.DeviceHandler
}

type RouterDeps struct {
	AuthMW      func(http.Handler) http.Handler
	PostService post.PostService
	StoragePath string
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
	r.Use(ratelimit.Middleware(ratelimit.Config{
		Requests: 300,
		Window:   time.Minute,
		Burst:    300,
		KeyFunc:  ratelimit.IPKey,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
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

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	if deps.StoragePath != "" {
		r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(deps.StoragePath))))
	}

	r.Route("/api/v1", func(r chi.Router) {
		authLimiter := ratelimit.Middleware(ratelimit.Config{
			Requests: 10,
			Window:   time.Minute,
			Burst:    10,
			KeyFunc:  ratelimit.IPKey,
		})
		userLimiter := ratelimit.Middleware(ratelimit.Config{
			Requests: 120,
			Window:   time.Minute,
			Burst:    120,
			KeyFunc:  ratelimit.UserIDKey(middleware.UserIDKey),
		})
		geoLimiter := ratelimit.Middleware(ratelimit.Config{
			Requests: 30,
			Window:   time.Minute,
			Burst:    30,
			KeyFunc:  ratelimit.UserIDKey(middleware.UserIDKey),
		})

		// r.Route("/admin", func(r chi.Router) {
		// r.Use(middleware.RequireRole("admin"))
		// r.Get("/users", h.Admin.ListUsers)
		// r.Post("/stats/reset", h.Admin.ResetStats)
		// })

		// A new client has no access token yet. Refresh uses its refresh token.
		r.With(authLimiter).Post("/auth/register", h.Auth.Register)
		r.With(authLimiter).Post("/auth/login", h.Auth.Login)
		r.With(authLimiter).Post("/auth/refresh", h.Auth.Refresh)

		r.Group(func(r chi.Router) {
			r.Use(deps.AuthMW)
			r.Route("/walks", func(r chi.Router) {
				r.With(geoLimiter).Post("/start", h.Walk.Start)
				r.Get("/current", h.Walk.Current)
				r.With(geoLimiter).Patch("/current/location", h.Walk.UpdateLocation)
				r.Patch("/current/visibility", h.Walk.UpdateVisibility)
				r.Post("/current/stop", h.Walk.Stop)
				r.With(geoLimiter).Get("/nearby", h.Walk.Nearby)
			})

			r.Get("/auth/me", h.Auth.Me)
			r.Post("/auth/logout", h.Auth.Logout)
			r.With(userLimiter).Post("/devices/tokens", h.Device.Register)
			r.With(userLimiter).Delete("/devices/tokens", h.Device.Unregister)

			r.Route("/users", func(r chi.Router) {
				r.Use(userLimiter)
				r.Put("/me/password", h.User.ChangePassword)
				r.Get("/{id}", h.User.GetByID)
				r.Put("/{id}", h.User.Update)
				r.Put("/{id}/role", h.User.ChangeRole)
				r.Delete("/{id}", h.User.Delete)
				r.Get("/{userID}/profile", h.Profile.GetByUserID)

				r.Post("/{userID}/profile", h.Profile.Create)
				r.Get("/{userID}/profile", h.Profile.GetByUserID)
				r.Put("/{userID}/profile", h.Profile.Update)
				r.Delete("/{userID}/profile", h.Profile.Delete)

				r.With(geoLimiter).Put("/{id}/location", h.User.UpdateLocation)
				r.Put("/{id}/location-visibility", h.User.SetLocationVisibility)
				r.With(geoLimiter).Get("/dogs/nearby", h.Walk.NearbyDogs)
			})

			r.With(userLimiter).Get("/profiles/search", h.Profile.Search)

			r.Route("/posts", func(r chi.Router) {
				r.Use(userLimiter)
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

			r.With(userLimiter).Get("/feed", h.Feed.GetFeed)

			r.Route("/comments", func(r chi.Router) {
				r.Use(userLimiter)
				r.Delete("/{id}", h.Comment.Delete)
			})

			r.Route("/follows", func(r chi.Router) {
				r.Use(userLimiter)
				r.Post("/{userID}", h.Follow.Follow)
				r.Delete("/{userID}", h.Follow.Unfollow)
				r.Get("/followers/{userID}", h.Follow.GetFollowers)
				r.Get("/following/{userID}", h.Follow.GetFollowing)
			})

			r.Route("/dogs", func(r chi.Router) {
				r.Use(userLimiter)
				r.Post("/", h.Dog.Create)
				r.Get("/", h.Dog.ListMine)
				r.With(geoLimiter).Get("/nearby", h.Walk.NearbyDogs)
				r.Get("/{id}", h.Dog.GetPrivate)
				r.Put("/{id}", h.Dog.Update)
				r.Delete("/{id}", h.Dog.Delete)

				r.Route("/{id}/vaccinations", func(r chi.Router) {
					r.Post("/", h.Vaccination.Create)
					r.Get("/", h.Vaccination.ListByDogID)
					r.Put("/{vaccinationID}", h.Vaccination.Update)
					r.Delete("/{vaccinationID}", h.Vaccination.Delete)
				})
			})

			r.Route("/vaccinations", func(r chi.Router) {
				r.Use(userLimiter)
				r.Delete("/{id}", h.Vaccination.Delete)
			})

			r.Route("/stats", func(r chi.Router) {
				r.Use(userLimiter)
				r.Get("/user/{userID}", h.Stats.UserStats)
				r.Get("/post/{postID}", h.Stats.PostStats)
				r.Get("/profile/{userID}", h.Stats.ProfileStats)
			})

			r.Route("/settings", func(r chi.Router) {
				r.Use(userLimiter)
				r.Get("/", h.Settings.Get)
				r.Put("/", h.Settings.Update)
			})

			r.Route("/upload", func(r chi.Router) {
				r.Use(ratelimit.Middleware(ratelimit.Config{
					Requests: 20,
					Window:   time.Minute,
					Burst:    20,
					KeyFunc:  ratelimit.UserIDKey(middleware.UserIDKey),
				}))
				r.Post("/avatar", h.Upload.UploadAvatar)
				r.Post("/post-image", h.Upload.UploadPostImage)
				r.Post("/dog-image", h.Upload.UploadDogImage)
			})

			r.With(ratelimit.Middleware(ratelimit.Config{
				Requests: 10,
				Window:   time.Minute,
				Burst:    10,
				KeyFunc:  ratelimit.UserIDKey(middleware.UserIDKey),
			})).Get("/ws", h.WS.ServeWS)
		})
	})

	return r
}
