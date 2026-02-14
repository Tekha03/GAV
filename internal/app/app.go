package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"gav/internal/auth"
	"gav/internal/config"
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
	sqlDB  *gorm.DB
	logger *slog.Logger
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("initializing application")

	db, err := gorm.Open(sqlite.Open(cfg.DB.Path), &gorm.Config{})
	if err != nil {
		logger.Error("failed to open database", "error", err)
		return nil, err
	}

	logger.Info("database opened", "path", cfg.DB.Path)

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("failed to get sql.DB", "error", err)
		return nil, err
	}

	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	logger.Info("database pool configured", "max_open_conns", 1)

	jwtConfig := auth.JWTConfig{
		Secret: []byte(cfg.JWT.Secret),
		TTL: 	cfg.JWT.TTL,
	}

	logger.Info("jwt cofidured", "ttl", cfg.JWT.TTL.String())

	// repositories
	userRepo := gavSqlite.NewUserRepository(db)
	postRepo := gavSqlite.NewPostRepository(db)

	// services
	authService := auth.NewService(userRepo, jwtConfig)
	userService := user.NewService(userRepo)
	postService := post.NewService(postRepo)

	// handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService)

	router := httptransport.NewRouter(
		authHandler,
		userHandler,
		postHandler,
		middleware.JWTAuth(jwtConfig),
		logger,
	)

	server := &http.Server{
		Addr:	":" + cfg.HTTP.Port,
		Handler: router,
	}

	logger.Info("http server configured", "port", cfg.HTTP.Port)

	return &App{server: server, sqlDB: db, logger: logger}, nil
}

func (a *App) Run() error {
	a.logger.Info("starting http server", "addr", a.server.Addr)

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.logger.Error("http server failed", "error", err)
		return err
	}

	a.logger.Info("http server stopped")

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("shutting down http server")

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("http server shutdown failed", "error", err)
		return err
	}

	sqlDB, err := a.sqlDB.DB()
	if err != nil {
		a.logger.Error("failed to get sql.DB for shutdown", "error", err)
		return err
	}

	if err := sqlDB.Close(); err != nil {
			a.logger.Error("failed to close database", "error", err)
			return err
	}

	a.logger.Info("application shutdown completed")

	return nil
}
