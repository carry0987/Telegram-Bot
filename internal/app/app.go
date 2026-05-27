package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"golang.org/x/sync/errgroup"

	appbot "telegram-bot/internal/bot"
	"telegram-bot/internal/config"
	"telegram-bot/internal/handler"
	"telegram-bot/internal/session"
)

type App struct {
	cfg        *config.Config
	bot        *appbot.Service
	httpServer *http.Server
	sessions   session.Store
}

func New(cfg *config.Config) (*App, error) {
	sessionStore, err := session.NewStore(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("init session store: %w", err)
	}

	botService, err := appbot.New(cfg, sessionStore)
	if err != nil {
		_ = sessionStore.Close()
		return nil, fmt.Errorf("init telegram bot: %w", err)
	}

	h := handler.New(handler.Options{
		Readiness:      botService,
		Mode:           cfg.BotMode,
		SessionBackend: sessionStore.Backend(),
		WebhookPath:    botService.WebhookPath(),
		WebhookHandler: botService.WebhookHandler(),
	})
	mux := http.NewServeMux()
	h.Register(mux)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           mux,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	return &App{
		cfg:        cfg,
		bot:        botService,
		httpServer: server,
		sessions:   sessionStore,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	if err := a.bot.Prepare(ctx); err != nil {
		return fmt.Errorf("prepare bot transport: %w", err)
	}

	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		slog.Info("telegram bot transport started", "mode", a.cfg.BotMode)
		if err := a.bot.Run(groupCtx); err != nil {
			return fmt.Errorf("run bot: %w", err)
		}

		return nil
	})

	group.Go(func() error {
		slog.Info("http server listening", "addr", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve http: %w", err)
		}

		return nil
	})

	group.Go(func() error {
		<-groupCtx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
		defer cancel()

		if err := a.httpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("shutdown http server: %w", err)
		}

		if err := a.sessions.Close(); err != nil {
			return fmt.Errorf("close session store: %w", err)
		}

		return nil
	})

	return group.Wait()
}
