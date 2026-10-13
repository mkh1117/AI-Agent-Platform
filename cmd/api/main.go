
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mkh1117/agentops/internal/config"
	"github.com/mkh1117/agentops/internal/platform/database"
	"github.com/mkh1117/agentops/internal/platform/httpserver"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	db, err := database.NewPostgres(context.Background(), cfg.DB)
	if err != nil {
		logger.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("failed to get database pool", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err := sqlDB.Close(); err != nil {
			logger.Error("failed to close postgres", "error", err)
		}
	}()

	api, err := httpserver.New(cfg.HTTP.Addr, db, logger)
	if err != nil {
		logger.Error("failed to initialize HTTP server", "error", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:              api.Addr(),
		Handler:           api.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("starting API server", "addr", server.Addr)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}

		close(serverErrors)
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")

	case err := <-serverErrors:
		if err != nil {
			logger.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
		return
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)

		if err := server.Close(); err != nil {
			logger.Error("forced server close failed", "error", err)
		}
	}

	logger.Info("API server stopped")
}
