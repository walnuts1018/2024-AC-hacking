package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/walnuts1018/2024-AC-hacking/config"
	"github.com/walnuts1018/2024-AC-hacking/psql"
	"github.com/walnuts1018/2024-AC-hacking/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to load config: %v", err))
	}

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		TimeFormat: time.RFC3339,
		Level:      cfg.LogLevel,
	}))
	slog.SetDefault(logger)

	psqlClient, err := psql.NewClient(cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create psql client: %v", err))
		os.Exit(1)
	}
	defer psqlClient.Close()

	if err := psqlClient.Init(); err != nil {
		slog.Error(fmt.Sprintf("Failed to initialize psql client: %v", err))
		os.Exit(1)
	}

	r, err := router.NewRouter(cfg, psqlClient)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create router: %v", err))
		os.Exit(1)
	}

	if err := r.Run(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
		slog.Error(fmt.Sprintf("Failed to start server: %v", err))
		os.Exit(1)
	}
}
