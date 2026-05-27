package config

import (
	"strings"
	"testing"
	"time"
)

func TestValidateAcceptsMinimalConfig(t *testing.T) {
	cfg := &Config{
		Port:              3000,
		TelegramBotToken:  "123456:abc",
		BotMode:           BotModePolling,
		BotInitTimeout:    5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		ShutdownTimeout:   10 * time.Second,
		WebhookPath:       "/telegram/webhook",
		SessionTTL:        24 * time.Hour,
	}

	if errs := cfg.Validate(); len(errs) > 0 {
		t.Fatalf("expected valid config, got errors: %v", errs)
	}
}

func TestValidateRejectsInvalidValues(t *testing.T) {
	cfg := &Config{}
	errText := strings.Join(cfg.Validate(), "\n")

	for _, expected := range []string{
		"PORT must be between 1 and 65535",
		"TELEGRAM_BOT_TOKEN is required",
		"BOT_MODE must be polling|webhook",
		"BOT_INIT_TIMEOUT must be greater than 0",
		"READ_HEADER_TIMEOUT must be greater than 0",
		"SHUTDOWN_TIMEOUT must be greater than 0",
		"SESSION_TTL must be greater than 0",
		"WEBHOOK_PATH must start with /",
	} {
		if !strings.Contains(errText, expected) {
			t.Fatalf("expected validation errors to contain %q, got %q", expected, errText)
		}
	}
}

func TestValidateWebhookModeRequiresPublicURL(t *testing.T) {
	cfg := &Config{
		Port:              3000,
		TelegramBotToken:  "123456:abc",
		BotMode:           BotModeWebhook,
		BotInitTimeout:    5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		ShutdownTimeout:   10 * time.Second,
		WebhookPath:       "/telegram/webhook",
		SessionTTL:        24 * time.Hour,
	}

	errText := strings.Join(cfg.Validate(), "\n")
	if !strings.Contains(errText, "WEBHOOK_PUBLIC_URL is required") {
		t.Fatalf("expected webhook validation error, got %q", errText)
	}
}

func TestNormalizeTrimsTelegramBotUsernamePrefix(t *testing.T) {
	cfg := &Config{TelegramBotUsername: "  @adakrei_bot  "}

	cfg.normalize()

	if cfg.TelegramBotUsername != "adakrei_bot" {
		t.Fatalf("expected normalized username without leading @, got %q", cfg.TelegramBotUsername)
	}
}
