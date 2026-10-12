package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/mkh1117/agentops/internal/config"
	"github.com/mkh1117/agentops/internal/platform/database"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	db, err := database.NewPostgres(context.Background(), cfg.DB)
	if err != nil {
		slog.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("failed to get database connection pool", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err := sqlDB.Close(); err != nil {
			slog.Error("failed to close postgres connection", "error", err)
		}
	}()

	slog.Info(
		"postgres connected successfully",
		"host", cfg.DB.Host,
		"port", cfg.DB.Port,
		"database", cfg.DB.Name,
	)
}