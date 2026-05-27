package bot

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"telegram-bot/internal/config"
	"telegram-bot/internal/session"
)

type Service struct {
	client             telegramClient
	username           string
	mode               string
	webhookPublicURL   string
	webhookPath        string
	webhookSecretToken string
	webhookDropPending bool
	sessionTTL         time.Duration
	sessions           session.Store
}

type telegramClient interface {
	Start(context.Context)
	StartWebhook(context.Context)
	SetWebhook(context.Context, *tgbot.SetWebhookParams) (bool, error)
	GetWebhookInfo(context.Context) (*models.WebhookInfo, error)
	DeleteWebhook(context.Context, *tgbot.DeleteWebhookParams) (bool, error)
	AnswerInlineQuery(context.Context, *tgbot.AnswerInlineQueryParams) (bool, error)
	GetMe(context.Context) (*models.User, error)
	WebhookHandler() http.HandlerFunc
	RegisterHandler(tgbot.HandlerType, string, tgbot.MatchType, tgbot.HandlerFunc, ...tgbot.Middleware) string
	RegisterHandlerMatchFunc(tgbot.MatchFunc, tgbot.HandlerFunc, ...tgbot.Middleware) string
	SendMessage(context.Context, *tgbot.SendMessageParams) (*models.Message, error)
	AnswerCallbackQuery(context.Context, *tgbot.AnswerCallbackQueryParams) (bool, error)
	SetMyCommands(context.Context, *tgbot.SetMyCommandsParams) (bool, error)
}

func New(cfg *config.Config, sessions session.Store) (*Service, error) {
	service := &Service{
		username:           cfg.TelegramBotUsername,
		mode:               cfg.BotMode,
		webhookPublicURL:   cfg.WebhookPublicURL,
		webhookPath:        cfg.WebhookPath,
		webhookSecretToken: cfg.WebhookSecretToken,
		webhookDropPending: cfg.WebhookDropPendingUpdates,
		sessionTTL:         cfg.SessionTTL,
		sessions:           sessions,
	}

	options := []tgbot.Option{
		tgbot.WithDefaultHandler(service.defaultHandler),
		tgbot.WithCheckInitTimeout(cfg.BotInitTimeout),
	}
	if cfg.BotMode == config.BotModeWebhook && cfg.WebhookSecretToken != "" {
		options = append(options, tgbot.WithWebhookSecretToken(cfg.WebhookSecretToken))
	}

	if cfg.Debug {
		options = append(options, tgbot.WithDebug())
	}

	client, err := tgbot.New(cfg.TelegramBotToken, options...)
	if err != nil {
		return nil, err
	}

	service.client = client
	service.registerHandlers()

	return service, nil
}

func (s *Service) Run(ctx context.Context) error {
	if s.mode == config.BotModeWebhook {
		s.client.StartWebhook(ctx)

		return nil
	}

	s.client.Start(ctx)

	return nil
}

func (s *Service) Prepare(ctx context.Context) error {
	if err := s.syncCommands(ctx); err != nil {
		return fmt.Errorf("sync commands: %w", err)
	}

	if s.mode == config.BotModeWebhook {
		if _, err := s.client.SetWebhook(ctx, &tgbot.SetWebhookParams{
			URL:                s.webhookURL(),
			SecretToken:        s.webhookSecretToken,
			DropPendingUpdates: s.webhookDropPending,
		}); err != nil {
			return fmt.Errorf("set webhook: %w", err)
		}

		slog.Info("telegram webhook configured", "url", s.webhookURL(), "path", s.webhookPath)

		return nil
	}

	webhookInfo, err := s.client.GetWebhookInfo(ctx)
	if err != nil {
		return fmt.Errorf("get webhook info: %w", err)
	}

	if webhookInfo == nil || strings.TrimSpace(webhookInfo.URL) == "" {
		return nil
	}

	if _, err := s.client.DeleteWebhook(ctx, &tgbot.DeleteWebhookParams{
		DropPendingUpdates: s.webhookDropPending,
	}); err != nil {
		return fmt.Errorf("delete webhook: %w", err)
	}

	return nil
}

func (s *Service) Ping(ctx context.Context) error {
	_, err := s.client.GetMe(ctx)

	return err
}

func (s *Service) WebhookHandler() http.Handler {
	if s.mode != config.BotModeWebhook {
		return nil
	}

	return s.client.WebhookHandler()
}

func (s *Service) WebhookPath() string {
	if s.mode != config.BotModeWebhook {
		return ""
	}

	return s.webhookPath
}

func (s *Service) registerHandlers() {
	s.client.RegisterHandler(tgbot.HandlerTypeMessageText, "start", tgbot.MatchTypeCommandStartOnly, s.handleStart)
	s.client.RegisterHandler(tgbot.HandlerTypeMessageText, "help", tgbot.MatchTypeCommandStartOnly, s.handleHelp)
	s.client.RegisterHandler(tgbot.HandlerTypeMessageText, "ping", tgbot.MatchTypeCommandStartOnly, s.handlePing)
	s.client.RegisterHandler(tgbot.HandlerTypeMessageText, "echo", tgbot.MatchTypeCommand, s.handleEcho)
	s.client.RegisterHandler(tgbot.HandlerTypeMessageText, "keyboard", tgbot.MatchTypeCommandStartOnly, s.handleKeyboard)
	s.client.RegisterHandler(tgbot.HandlerTypeMessageText, "hidekeyboard", tgbot.MatchTypeCommandStartOnly, s.handleHideKeyboard)
	s.client.RegisterHandler(tgbot.HandlerTypeMessageText, "menu", tgbot.MatchTypeCommandStartOnly, s.handleMenu)
	s.client.RegisterHandler(tgbot.HandlerTypeMessageText, "session", tgbot.MatchTypeCommandStartOnly, s.handleSession)
	s.client.RegisterHandler(tgbot.HandlerTypeCallbackQueryData, "menu:", tgbot.MatchTypePrefix, s.handleMenuAction)
	s.client.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		return update != nil && update.InlineQuery != nil
	}, s.handleInlineQuery)
}

func (s *Service) defaultHandler(ctx context.Context, _ *tgbot.Bot, update *models.Update) {
	message := update.Message
	if message == nil || strings.TrimSpace(message.Text) == "" {
		return
	}

	if err := s.sendText(ctx, message.Chat.ID, s.helpMessage()); err != nil {
		slog.Error("send default help message", "error", err)
	}
}

func (s *Service) handleStart(ctx context.Context, _ *tgbot.Bot, update *models.Update) {
	message := update.Message
	if message == nil {
		return
	}

	name := "there"
	if message.From != nil && strings.TrimSpace(message.From.FirstName) != "" {
		name = strings.TrimSpace(message.From.FirstName)
	}

	visits, err := s.sessions.Increment(ctx, message.Chat.ID, "visits", s.sessionTTL)
	if err != nil {
		slog.Error("increment start visits", "error", err)
		visits = 1
	}
	_ = s.sessions.Set(ctx, message.Chat.ID, "last_command", "/start", s.sessionTTL)

	text := fmt.Sprintf(
		"Hello, %s.\nThis repository is a Go Telegram Bot template.\nUse /help to see the available commands.\nUse /keyboard to reopen the reply keyboard later.\nSession visits: %d.",
		name, visits,
	)

	if err := s.sendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:      message.Chat.ID,
		Text:        text,
		ReplyMarkup: s.replyKeyboardMarkup(),
	}); err != nil {
		slog.Error("send start message", "error", err)
	}
}

func (s *Service) handleHelp(ctx context.Context, _ *tgbot.Bot, update *models.Update) {
	message := update.Message
	if message == nil {
		return
	}

	if err := s.sendText(ctx, message.Chat.ID, s.helpMessage()); err != nil {
		slog.Error("send help message", "error", err)
	}
	_ = s.sessions.Set(ctx, message.Chat.ID, "last_command", "/help", s.sessionTTL)
}

func (s *Service) handlePing(ctx context.Context, _ *tgbot.Bot, update *models.Update) {
	message := update.Message
	if message == nil {
		return
	}

	if err := s.sendText(ctx, message.Chat.ID, "pong"); err != nil {
		slog.Error("send ping message", "error", err)
	}
	_ = s.sessions.Set(ctx, message.Chat.ID, "last_command", "/ping", s.sessionTTL)
}

func (s *Service) handleEcho(ctx context.Context, _ *tgbot.Bot, update *models.Update) {
	message := update.Message
	if message == nil {
		return
	}

	args := commandArgs(message.Text)
	if args == "" {
		if err := s.sendText(ctx, message.Chat.ID, "Usage: /echo <text>"); err != nil {
			slog.Error("send echo usage", "error", err)
		}

		return
	}

	if err := s.sendText(ctx, message.Chat.ID, args); err != nil {
		slog.Error("send echo reply", "error", err)
	}
	_ = s.sessions.Set(ctx, message.Chat.ID, "last_command", "/echo", s.sessionTTL)
}

func (s *Service) handleKeyboard(ctx context.Context, _ *tgbot.Bot, update *models.Update) {
	message := update.Message
	if message == nil {
		return
	}

	if err := s.sendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:      message.Chat.ID,
		Text:        "Reply keyboard enabled.",
		ReplyMarkup: s.replyKeyboardMarkup(),
	}); err != nil {
		slog.Error("send keyboard message", "error", err)
	}
	_ = s.sessions.Set(ctx, message.Chat.ID, "last_command", "/keyboard", s.sessionTTL)
}

func (s *Service) handleHideKeyboard(ctx context.Context, _ *tgbot.Bot, update *models.Update) {
	message := update.Message
	if message == nil {
		return
	}

	if err := s.sendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:      message.Chat.ID,
		Text:        "Reply keyboard hidden. Use /keyboard to show it again.",
		ReplyMarkup: s.replyKeyboardRemoveMarkup(),
	}); err != nil {
		slog.Error("send hide keyboard message", "error", err)
	}
	_ = s.sessions.Set(ctx, message.Chat.ID, "last_command", "/hidekeyboard", s.sessionTTL)
}

func (s *Service) handleMenu(ctx context.Context, _ *tgbot.Bot, update *models.Update) {
	message := update.Message
	if message == nil {
		return
	}

	if err := s.sendMessage(ctx, &tgbot.SendMessageParams{
		ChatID:      message.Chat.ID,
		Text:        "Template menu: choose an action.",
		ReplyMarkup: s.inlineMenuMarkup(),
	}); err != nil {
		slog.Error("send menu message", "error", err)
	}
	_ = s.sessions.Set(ctx, message.Chat.ID, "last_command", "/menu", s.sessionTTL)
}

func (s *Service) handleSession(ctx context.Context, _ *tgbot.Bot, update *models.Update) {
	message := update.Message
	if message == nil {
		return
	}

	visits, err := s.sessions.Increment(ctx, message.Chat.ID, "session_hits", s.sessionTTL)
	if err != nil {
		slog.Error("increment session counter", "error", err)
		return
	}

	lastCommand, err := s.sessions.Get(ctx, message.Chat.ID, "last_command")
	if err != nil && err != session.ErrNotFound {
		slog.Error("get last command", "error", err)
	}
	_ = s.sessions.Set(ctx, message.Chat.ID, "last_command", "/session", s.sessionTTL)

	text := fmt.Sprintf("Session counter: %d\nBackend: %s", visits, s.sessions.Backend())
	if lastCommand != "" {
		text += fmt.Sprintf("\nPrevious command: %s", lastCommand)
	}

	if err := s.sendText(ctx, message.Chat.ID, text); err != nil {
		slog.Error("send session message", "error", err)
	}
}

func (s *Service) handleMenuAction(ctx context.Context, _ *tgbot.Bot, update *models.Update) {
	callback := update.CallbackQuery
	if callback == nil {
		return
	}

	action := strings.TrimPrefix(callback.Data, "menu:")
	chatID := callback.From.ID

	var text string
	switch action {
	case "hello":
		text = "Hello from the inline keyboard demo."
	case "help":
		text = s.helpMessage()
	case "session":
		count, err := s.sessions.Increment(ctx, chatID, "session_hits", s.sessionTTL)
		if err != nil {
			slog.Error("increment callback session counter", "error", err)
			text = "Session update failed."
		} else {
			text = fmt.Sprintf("Session counter via inline keyboard: %d\nBackend: %s", count, s.sessions.Backend())
		}
	default:
		text = "Unknown menu action."
	}
	_ = s.sessions.Set(ctx, chatID, "last_command", "menu:"+action, s.sessionTTL)

	if _, err := s.client.AnswerCallbackQuery(ctx, &tgbot.AnswerCallbackQueryParams{
		CallbackQueryID: callback.ID,
		Text:            "Done",
	}); err != nil {
		slog.Error("answer callback query", "error", err)
	}

	if err := s.sendText(ctx, chatID, text); err != nil {
		slog.Error("send callback response", "error", err)
	}
}

func (s *Service) handleInlineQuery(ctx context.Context, _ *tgbot.Bot, update *models.Update) {
	inlineQuery := update.InlineQuery
	if inlineQuery == nil {
		return
	}

	if _, err := s.client.AnswerInlineQuery(ctx, &tgbot.AnswerInlineQueryParams{
		InlineQueryID: inlineQuery.ID,
		Results:       s.inlineQueryResults(strings.TrimSpace(inlineQuery.Query)),
		CacheTime:     1,
		IsPersonal:    true,
	}); err != nil {
		slog.Error("answer inline query", "error", err)
	}
}

func (s *Service) sendText(ctx context.Context, chatID int64, text string) error {
	return s.sendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
}

func (s *Service) sendMessage(ctx context.Context, params *tgbot.SendMessageParams) error {
	_, err := s.client.SendMessage(ctx, params)

	return err
}

func (s *Service) helpMessage() string {
	var builder strings.Builder
	builder.WriteString("Available commands:\n")
	builder.WriteString("/start - introduce the bot\n")
	builder.WriteString("/help - show this help\n")
	builder.WriteString("/ping - respond with pong\n")
	builder.WriteString("/echo <text> - echo the text back\n")
	builder.WriteString("/keyboard - show reply keyboard demo\n")
	builder.WriteString("/hidekeyboard - hide reply keyboard\n")
	builder.WriteString("/menu - show inline keyboard demo\n")
	builder.WriteString("/session - inspect session state")

	if s.username != "" {
		_, _ = fmt.Fprintf(&builder, "\n\nGroup chats can also use /echo@%s <text>.", s.username)
	}

	return builder.String()
}

func commandArgs(text string) string {
	parts := strings.Fields(strings.TrimSpace(text))
	if len(parts) <= 1 {
		return ""
	}

	return strings.Join(parts[1:], " ")
}

func (s *Service) syncCommands(ctx context.Context) error {
	_, err := s.client.SetMyCommands(ctx, &tgbot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{Command: "start", Description: "Introduce the bot"},
			{Command: "help", Description: "Show available commands"},
			{Command: "ping", Description: "Reply with pong"},
			{Command: "echo", Description: "Echo text back"},
			{Command: "keyboard", Description: "Show reply keyboard demo"},
			{Command: "hidekeyboard", Description: "Hide reply keyboard"},
			{Command: "menu", Description: "Show inline keyboard demo"},
			{Command: "session", Description: "Inspect session state"},
		},
	})

	return err
}

func (s *Service) webhookURL() string {
	return s.webhookPublicURL + s.webhookPath
}

func (s *Service) inlineQueryResults(query string) []models.InlineQueryResult {
	if query == "" {
		return []models.InlineQueryResult{
			s.inlineArticleResult(
				"inline-help",
				"Bot Help",
				"Insert a summary of the available bot commands",
				s.helpMessage(),
			),
			s.inlineArticleResult(
				"inline-ping",
				"Pong",
				"Insert a simple pong message",
				"pong",
			),
			s.inlineArticleResult(
				"inline-greeting",
				"Greeting",
				"Insert a greeting message from the template bot",
				"Hello from the Go Telegram Bot template.",
			),
		}
	}

	return []models.InlineQueryResult{
		s.inlineArticleResult(
			"inline-echo",
			fmt.Sprintf("Echo: %s", query),
			"Insert the inline query text as a message",
			query,
		),
		s.inlineArticleResult(
			"inline-help",
			"Bot Help",
			"Insert a summary of the available bot commands",
			s.helpMessage(),
		),
		s.inlineArticleResult(
			"inline-menu",
			"Menu Demo",
			"Insert a short note about the inline keyboard demo",
			"Try /menu in a direct chat with the bot to open the inline keyboard demo.",
		),
	}
}

func (s *Service) inlineArticleResult(id, title, description, messageText string) *models.InlineQueryResultArticle {
	return &models.InlineQueryResultArticle{
		ID:          id,
		Title:       title,
		Description: description,
		InputMessageContent: &models.InputTextMessageContent{
			MessageText: messageText,
		},
	}
}

func (s *Service) inlineMenuMarkup() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Hello", CallbackData: "menu:hello"},
				{Text: "Session", CallbackData: "menu:session"},
			},
			{
				{Text: "Help", CallbackData: "menu:help"},
			},
		},
	}
}

func (s *Service) replyKeyboardMarkup() *models.ReplyKeyboardMarkup {
	return &models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: "/help"},
				{Text: "/ping"},
			},
			{
				{Text: "/menu"},
				{Text: "/session"},
			},
			{
				{Text: "/hidekeyboard"},
			},
		},
		IsPersistent:          true,
		ResizeKeyboard:        true,
		InputFieldPlaceholder: "Choose a template command",
	}
}

func (s *Service) replyKeyboardRemoveMarkup() *models.ReplyKeyboardRemove {
	return &models.ReplyKeyboardRemove{RemoveKeyboard: true}
}
