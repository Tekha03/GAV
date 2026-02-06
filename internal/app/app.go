package app

import (
	"context"
	"net/http"

	"gav/internal/auth"
	"gav/internal/config"

	// "gav/internal/comment"
	"gav/internal/post"
	"gav/internal/transport/http/handlers"
	"gav/internal/transport/http/middleware"
	"gav/internal/user"

	httptransport "gav/internal/transport/http"
	gavSqlite "gav/storage/sqlite"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	server *http.Server
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DB.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// repositories
	userRepo := gavSqlite.NewUserRepository(db)
	postRepo := gavSqlite.NewPostRepository(db)
	// commentRepo := gavSqlite.NewCommentRepository(db)

	// services
	authService := auth.NewService(
		userRepo,
		auth.JWTConfig{
			Secret: []byte(cfg.JWT.Secret),
			TTL: cfg.JWT.TTL,
		},
	)
	userService := user.NewService(userRepo)
	postService := post.NewService(postRepo)
	// commentService := comment.NewService(commentRepo)

	// handlers
	authHandler := handlers.NewAuthHandler(*authService)
	userHandler := handlers.NewUserHandler(*userService)
	postHandler := handlers.NewPostHandler(*postService)
	// commentHandler := handlers.NewCommentHandler(*commentService)

	router := httptransport.NewRouter(
		authHandler,
		userHandler,
		postHandler,
		middleware.JWTAuth(auth.JWTConfig{
			Secret: []byte(cfg.JWT.Secret),
			TTL: cfg.JWT.TTL,
		}),
	)

	server := &http.Server{
		Addr:	":" + cfg.HTTP.Port,
		Handler: router,
	}

	return &App{server: server}, nil
}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return  a.server.Shutdown(ctx)
}
