package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"

	"telegram-bot/internal/app"
	"telegram-bot/internal/config"
)

var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	_ = godotenv.Overload(".env", ".env.local")

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to load config: %v\n", err)
		os.Exit(1)
	}

	if errs := cfg.Validate(); len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "ERROR: invalid configuration:\n")
		for _, entry := range errs {
			fmt.Fprintf(os.Stderr, "  - %s\n", entry)
		}
		os.Exit(1)
	}

	setupLogger(cfg.Debug)
	slog.Info("starting telegram-bot", "version", version, "commit", commit)

	application, err := app.New(cfg)
	if err != nil {
		slog.Error("build app", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := application.Run(ctx); err != nil {
		slog.Error("run app", "error", err)
		os.Exit(1)
	}

	slog.Info("telegram-bot stopped")
}

func setupLogger(debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.TimeKey {
				attr.Value = slog.StringValue(strings.TrimSpace(attr.Value.String()))
			}

			return attr
		},
	}))

	slog.SetDefault(logger)
}
