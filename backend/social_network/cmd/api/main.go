// @title           Dog Lovers API (GAV)
// @version         1.0
// @description     Социальная сеть для любителей собак
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                JWT Authorization header using the Bearer scheme. Example: "Bearer {token}"

package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"social_network/internal/app"
	"social_network/internal/config"
	"social_network/internal/kafka"
	grpctransport "social_network/transport/grpc"

	_ "docs"
	grpcpkg "google.golang.org/grpc"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	logger.Info("starting application")

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	application, err := app.NewApp(ctx, cfg)
	if err != nil {
		logger.Error("failed to create application", "error", err)
		os.Exit(1)
	}

	grpcListener, err := net.Listen("tcp", cfg.GRPC.Addr)
	if err != nil {
		logger.Error("failed to listen for gRPC", "addr", cfg.GRPC.Addr, "error", err)
		os.Exit(1)
	}
	grpcServer := grpcpkg.NewServer()
	grpctransport.Register(grpcServer, application.Services.User, application.Services.Profile, application.Services.Auth, application.Services.Notification, cfg.JWT.TTL)
	go func() {
		logger.Info("starting gRPC server", "addr", cfg.GRPC.Addr)
		if err := grpcServer.Serve(grpcListener); err != nil {
			logger.Error("gRPC server failed", "error", err)
			stop()
		}
	}()

	if os.Getenv("KAFKA_ENABLED") == "true" {
		if err := kafka.LaunchKafka(ctx, application.Services.Notification); err != nil {
			logger.Error("failed to launch Kafka", "error", err)
			os.Exit(1)
		}
	} else {
		logger.Info("Kafka consumer disabled", "KAFKA_ENABLED", os.Getenv("KAFKA_ENABLED"))
	}

	go func() {
		logger.Info("starting HTTP server", "addr", application.Server.Addr)
		if err := application.Run(); err != nil && err != http.ErrServerClosed {
			logger.Error("http server failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()

	logger.Info("received shutdown signal")
	grpcServer.GracefulStop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("application stopped gracefully")
}
