package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	redis "github.com/redis/go-redis/v9"
)

const (
	BotModePolling = "polling"
	BotModeWebhook = "webhook"
)

type Config struct {
	Port                      int           `env:"PORT" envDefault:"3000"`
	Debug                     bool          `env:"DEBUG" envDefault:"false"`
	TelegramBotToken          string        `env:"TELEGRAM_BOT_TOKEN,required"`
	TelegramBotUsername       string        `env:"TELEGRAM_BOT_USERNAME"`
	BotMode                   string        `env:"BOT_MODE" envDefault:"polling"`
	BotInitTimeout            time.Duration `env:"BOT_INIT_TIMEOUT" envDefault:"5s"`
	ReadHeaderTimeout         time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"5s"`
	ShutdownTimeout           time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
	WebhookPublicURL          string        `env:"WEBHOOK_PUBLIC_URL"`
	WebhookPath               string        `env:"WEBHOOK_PATH" envDefault:"/telegram/webhook"`
	WebhookSecretToken        string        `env:"WEBHOOK_SECRET_TOKEN"`
	WebhookDropPendingUpdates bool          `env:"WEBHOOK_DROP_PENDING_UPDATES" envDefault:"false"`
	RedisURL                  string        `env:"REDIS_URL"`
	SessionTTL                time.Duration `env:"SESSION_TTL" envDefault:"24h"`
}

func Load() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	cfg.normalize()

	return &cfg, nil
}

func (c *Config) Validate() []string {
	var errs []string

	if c.Port < 1 || c.Port > 65535 {
		errs = append(errs, "PORT must be between 1 and 65535")
	}

	if strings.TrimSpace(c.TelegramBotToken) == "" {
		errs = append(errs, "TELEGRAM_BOT_TOKEN is required")
	}

	switch c.BotMode {
	case BotModePolling, BotModeWebhook:
	default:
		errs = append(errs, fmt.Sprintf("BOT_MODE must be %s|%s", BotModePolling, BotModeWebhook))
	}

	if c.BotInitTimeout <= 0 {
		errs = append(errs, "BOT_INIT_TIMEOUT must be greater than 0")
	}

	if c.ReadHeaderTimeout <= 0 {
		errs = append(errs, "READ_HEADER_TIMEOUT must be greater than 0")
	}

	if c.ShutdownTimeout <= 0 {
		errs = append(errs, "SHUTDOWN_TIMEOUT must be greater than 0")
	}

	if c.SessionTTL <= 0 {
		errs = append(errs, "SESSION_TTL must be greater than 0")
	}

	if c.WebhookPath == "" || !strings.HasPrefix(c.WebhookPath, "/") {
		errs = append(errs, "WEBHOOK_PATH must start with /")
	}

	if c.BotMode == BotModeWebhook {
		if err := validateHTTPURL(c.WebhookPublicURL, "WEBHOOK_PUBLIC_URL"); err != "" {
			errs = append(errs, err)
		}
		if c.WebhookPath == "/" {
			errs = append(errs, "WEBHOOK_PATH cannot be / in webhook mode")
		}
	}

	if c.RedisURL != "" {
		if _, err := redis.ParseURL(c.RedisURL); err != nil {
			errs = append(errs, fmt.Sprintf("REDIS_URL must be a valid redis URL: %v", err))
		}
	}

	return errs
}

func (c *Config) normalize() {
	c.TelegramBotUsername = strings.TrimPrefix(strings.TrimSpace(c.TelegramBotUsername), "@")
	c.TelegramBotToken = strings.TrimSpace(c.TelegramBotToken)
	c.BotMode = strings.TrimSpace(strings.ToLower(c.BotMode))
	c.WebhookPublicURL = strings.TrimRight(strings.TrimSpace(c.WebhookPublicURL), "/")
	c.WebhookPath = strings.TrimSpace(c.WebhookPath)
	if c.WebhookPath == "" {
		c.WebhookPath = "/telegram/webhook"
	}
	c.WebhookSecretToken = strings.TrimSpace(c.WebhookSecretToken)
	c.RedisURL = strings.TrimSpace(c.RedisURL)
}

func validateHTTPURL(raw, name string) string {
	if raw == "" {
		return fmt.Sprintf("%s is required", name)
	}

	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Sprintf("%s must be a valid URL: %v", name, err)
	}

	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Sprintf("%s must use http:// or https://", name)
	}

	if u.Host == "" {
		return fmt.Sprintf("%s must include a host", name)
	}

	return ""
}
