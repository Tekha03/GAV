package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"gav/dbserver"
	"gav/internal/auth"
	"gav/internal/config"
	"gav/internal/media"
	httptransport "gav/transport/http"
	"gav/transport/http/middleware"

	"gorm.io/gorm"
)

type App struct {
	Server *http.Server
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

	db, err := dbserver.InitDB(cfg.DB.Path, logger)
	if err != nil {
		logger.Error("failed to open database", "error", err)
		return nil, err
	}

	logger.Info("database opened", "path", cfg.DB.Path)

	if os.Getenv("ENV") != "production" {
		if err := dbserver.SeedDatabase(db, logger); err != nil {
			logger.Error("failed to seed database", "error", err)
			return nil, err
		}
		logger.Info("database seeding completed")
	}

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

	logger.Info("jwt cofigured", "ttl", cfg.JWT.TTL.String())

	mediaStorage := media.NewLocalStorage(cfg.Storage.LocalPath)

	// repositories
	repos, err := initRepositories(db);
	if err != nil {
		return nil, err
	}

	// services
	services, err := initServices(repos, jwtConfig, mediaStorage)
	if err != nil {
		return nil, err
	}

	// handlers
	handlers, err := initHandlers(services)
	if err != nil {
		return nil, err
	}

	router := httptransport.NewRouter(
		httptransport.Handlers{
			Auth:			handlers.Auth,
			User: 			handlers.User,
			Profile: 	 	 handlers.Profile,
			Post: 			handlers.Post,
			Feed:			handlers.Feed,
			Comment: 		handlers.Comment,
			Like: 			handlers.Like,
			Follow: 		handlers.Follow,
			Dog: 			handlers.Dog,
			Vaccination:	handlers.Vaccination,
			Stats: 			handlers.Stats,
			Settings: 		handlers.Settings,
			Upload: 		handlers.Upload,
		},
		httptransport.RouterDeps{
			AuthMW: middleware.JWTAuth(jwtConfig),
			PostService: services.Post,
		},
		logger,
	)

	server := &http.Server{
		Addr:	":" + cfg.HTTP.Port,
		Handler: router,
	}

	logger.Info("http server configured", "port", cfg.HTTP.Port)

	return &App{Server: server, sqlDB: db, logger: logger}, nil
}

func (a *App) Run() error {
	if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.logger.Error("http server failed", "error", err)
		return err
	}

	a.logger.Info("http server stopped")

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("shutting down http server")

	if err := a.Server.Shutdown(ctx); err != nil {
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
