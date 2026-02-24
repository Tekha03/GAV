package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"gav/internal/auth"
	"gav/internal/comment"
	"gav/internal/config"
	"gav/internal/dog"
	"gav/internal/follow"
	"gav/internal/like"
	"gav/internal/post"
	"gav/internal/profile"
	"gav/internal/settings"
	"gav/internal/stats"
	"gav/internal/user"
	"gav/internal/vaccination"
	"gav/transport/http/handlers"
	"gav/transport/http/middleware"

	gavSqlite "gav/storage/sqlite"
	httptransport "gav/transport/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	server *http.Server
	sqlDB  *gorm.DB
	logger *slog.Logger
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {
	if cfg == nil {
		return nil, ErrConfigNil
	}
	if cfg.DB.Path == "" {
		return nil, ErrDBPathEmpty
	}
	if cfg.JWT.Secret == "" {
		return nil, ErrJWTSecretEmpty
	}

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
	profileRepo := gavSqlite.NewProfileRepository(db)
	postRepo := gavSqlite.NewPostRepository(db)
	commentRepo := gavSqlite.NewCommentRepository(db)
	likeRepo := gavSqlite.NewLikeRepository(db)
	followRepo := gavSqlite.NewFollowRepository(db)
	dogRepo := gavSqlite.NewDogRepository(db)
	vaccinationRepo := gavSqlite.NewVaccinationRepository(db)
	statsRepo := gavSqlite.NewStatsRepository(db)
	settingsRepo := gavSqlite.NewSettingsRepository(db)

	// services
	userService, err := user.NewService(userRepo)
	if err != nil {
		return nil, err
	}

	hasher := &auth.PasswordHasher{}
	authService, err := auth.NewService(userService, jwtConfig, hasher)
	if err != nil {
		return nil, err
	}

	profileService, err := profile.NewService(profileRepo)
	if err != nil {
		return nil, err
	}

	postService, err := post.NewService(postRepo)
	if err != nil {
		return nil, err
	}

	commentService, err := comment.NewService(commentRepo)
	if err != nil {
		return nil, err
	}

	likeService, err := like.NewService(likeRepo)
	if err != nil {
		return nil, err
	}

	followService, err := follow.NewService(followRepo)
	if err != nil {
		return nil, err
	}

	dogService, err := dog.NewService(dogRepo)
	if err != nil {
		return nil, err
	}

	vaccinationService, err := vaccination.NewService(vaccinationRepo)
	if err != nil {
		return nil, err
	}

	statsService, err := stats.NewService(statsRepo)
	if err != nil {
		return nil, err
	}

	settingsService, err := settings.NewService(settingsRepo)
	if err != nil {
		return nil, err
	}

	// handlers
	handlers := httptransport.Handlers{
		Auth: handlers.NewAuthHandler(authService),
		User: handlers.NewUserHandler(userService),
		Profile: handlers.NewProfileHandler(profileService),
		Post: handlers.NewPostHandler(postService),
		Feed: handlers.NewFeedHandler(postService),
		Comment: handlers.NewCommentHandler(commentService),
		Like: handlers.NewLikeHandler(likeService),
		Follow: handlers.NewFollowHandler(followService),
		Dog: handlers.NewDogHandler(dogService),
		Vaccination: handlers.NewVaccinationHandler(vaccinationService),
		Stats: handlers.NewStatsHandler(statsService),
		Settings: handlers.NewSettingsHandler(settingsService),
	}

	router := httptransport.NewRouter(
		handlers,
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
