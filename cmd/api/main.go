package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gav/dbserver"
	"gav/internal/app"
	"gav/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	logger.Info("starting application")

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := dbserver.InitDB(cfg.DB.Path, logger)
	if err != nil {
		logger.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}

	defer func () {
		if err := dbserver.CloseDB(db); err != nil {
			logger.Warn("failed to close database", "error", err)
		} else {
			logger.Info("database closed successfully")
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := app.NewApp(ctx, cfg)
	if err != nil {
		logger.Error("failed to create application", "error", err)
		os.Exit(1)
	}

	go func () {
		logger.Info("starting HTTP server", "addr", app.Server.Addr)
		if err := app.Run(); err != nil && err != http.ErrServerClosed {
			logger.Error("http server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit

	logger.Info("received shutdown signal")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer shutdownCancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("application stopped gracefully")
}
