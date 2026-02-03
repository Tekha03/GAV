package app

import (
	"context"
	"net/http"

	"gav/internal/auth"
	// "gav/internal/comment"
	"gav/internal/post"
	"gav/internal/user"
	"gav/internal/transport/http/handlers"
	"gav/internal/transport/http/middleware"

	gavSqlite 		"gav/storage/sqlite"
	httptransport 	"gav/internal/transport/http"

	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

type App struct {
	server *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	db, err := gorm.Open(sqlite.Open("social.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// repositories
	userRepo := gavSqlite.NewUserRepository(db)
	postRepo := gavSqlite.NewPostRepository(db)
	// commentRepo := gavSqlite.NewCommentRepository(db)

	// services
	authService := auth.NewService(userRepo)
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
		middleware.JWTAuth,
	)

	server := &http.Server{
		Addr:	":8080",
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
