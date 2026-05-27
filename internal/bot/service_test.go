package bot

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"telegram-bot/internal/config"
	"telegram-bot/internal/session"
)

type stubTelegramClient struct {
	setMyCommandsCalled  bool
	getWebhookInfoCalled bool
	deleteWebhookCalled  bool
	answerInlineCalled   bool
	matchFuncCalled      bool
	registeredHandlers   []registeredHandler
	inlineAnswerParams   *tgbot.AnswerInlineQueryParams
	webhookInfo          *models.WebhookInfo
	setMyCommandsErr     error
	getWebhookInfoErr    error
	deleteWebhookErr     error
}

type registeredHandler struct {
	handlerType tgbot.HandlerType
	pattern     string
	matchType   tgbot.MatchType
}

func (s *stubTelegramClient) Start(context.Context) {}

func (s *stubTelegramClient) StartWebhook(context.Context) {}

func (s *stubTelegramClient) SetWebhook(context.Context, *tgbot.SetWebhookParams) (bool, error) {
	return true, nil
}

func (s *stubTelegramClient) GetWebhookInfo(context.Context) (*models.WebhookInfo, error) {
	s.getWebhookInfoCalled = true
	if s.getWebhookInfoErr != nil {
		return nil, s.getWebhookInfoErr
	}
	if s.webhookInfo == nil {
		return &models.WebhookInfo{}, nil
	}
	return s.webhookInfo, nil
}

func (s *stubTelegramClient) DeleteWebhook(context.Context, *tgbot.DeleteWebhookParams) (bool, error) {
	s.deleteWebhookCalled = true
	return s.deleteWebhookErr == nil, s.deleteWebhookErr
}

func (s *stubTelegramClient) AnswerInlineQuery(_ context.Context, params *tgbot.AnswerInlineQueryParams) (bool, error) {
	s.answerInlineCalled = true
	s.inlineAnswerParams = params
	return true, nil
}

func (s *stubTelegramClient) GetMe(context.Context) (*models.User, error) {
	return &models.User{}, nil
}

func (s *stubTelegramClient) WebhookHandler() http.HandlerFunc {
	return func(http.ResponseWriter, *http.Request) {}
}

func (s *stubTelegramClient) RegisterHandler(handlerType tgbot.HandlerType, pattern string, matchType tgbot.MatchType, _ tgbot.HandlerFunc, _ ...tgbot.Middleware) string {
	s.registeredHandlers = append(s.registeredHandlers, registeredHandler{
		handlerType: handlerType,
		pattern:     pattern,
		matchType:   matchType,
	})
	return ""
}

func (s *stubTelegramClient) RegisterHandlerMatchFunc(_ tgbot.MatchFunc, _ tgbot.HandlerFunc, _ ...tgbot.Middleware) string {
	s.matchFuncCalled = true
	return ""
}

func (s *stubTelegramClient) SendMessage(context.Context, *tgbot.SendMessageParams) (*models.Message, error) {
	return &models.Message{}, nil
}

func (s *stubTelegramClient) AnswerCallbackQuery(context.Context, *tgbot.AnswerCallbackQueryParams) (bool, error) {
	return true, nil
}

func (s *stubTelegramClient) SetMyCommands(context.Context, *tgbot.SetMyCommandsParams) (bool, error) {
	s.setMyCommandsCalled = true
	return s.setMyCommandsErr == nil, s.setMyCommandsErr
}

func TestPreparePollingSkipsDeleteWebhookWhenNoWebhookConfigured(t *testing.T) {
	service := newTestService(&stubTelegramClient{webhookInfo: &models.WebhookInfo{URL: ""}})

	if err := service.Prepare(context.Background()); err != nil {
		t.Fatalf("expected prepare to succeed, got %v", err)
	}

	client := service.client.(*stubTelegramClient)
	if !client.setMyCommandsCalled {
		t.Fatal("expected prepare to sync commands")
	}
	if !client.getWebhookInfoCalled {
		t.Fatal("expected prepare to inspect webhook state")
	}
	if client.deleteWebhookCalled {
		t.Fatal("expected prepare not to delete webhook when no webhook is configured")
	}
}

func TestPreparePollingDeletesWebhookWhenConfigured(t *testing.T) {
	service := newTestService(&stubTelegramClient{webhookInfo: &models.WebhookInfo{URL: "https://example.com/telegram/webhook"}})

	if err := service.Prepare(context.Background()); err != nil {
		t.Fatalf("expected prepare to succeed, got %v", err)
	}

	client := service.client.(*stubTelegramClient)
	if !client.deleteWebhookCalled {
		t.Fatal("expected prepare to delete webhook when one is configured")
	}
}

func TestPreparePollingReturnsDeleteWebhookError(t *testing.T) {
	service := newTestService(&stubTelegramClient{
		webhookInfo:      &models.WebhookInfo{URL: "https://example.com/telegram/webhook"},
		deleteWebhookErr: errors.New("unexpected end of JSON input"),
	})

	err := service.Prepare(context.Background())
	if err == nil {
		t.Fatal("expected prepare to return delete webhook error")
	}
}

func TestRegisterHandlersUseCommandNamesWithoutLeadingSlash(t *testing.T) {
	client := &stubTelegramClient{}
	service := newTestService(client)

	service.registerHandlers()

	if !client.matchFuncCalled {
		t.Fatal("expected inline query match handler to be registered")
	}

	for _, handler := range client.registeredHandlers {
		if handler.matchType != tgbot.MatchTypeCommand && handler.matchType != tgbot.MatchTypeCommandStartOnly {
			continue
		}
		if len(handler.pattern) > 0 && handler.pattern[0] == '/' {
			t.Fatalf("expected command pattern without leading slash, got %q", handler.pattern)
		}
	}
}

func TestHandleInlineQueryReturnsDefaultResultsForEmptyQuery(t *testing.T) {
	client := &stubTelegramClient{}
	service := newTestService(client)

	service.handleInlineQuery(context.Background(), nil, &models.Update{
		InlineQuery: &models.InlineQuery{ID: "inline-empty", Query: "   "},
	})

	if !client.answerInlineCalled {
		t.Fatal("expected inline query to be answered")
	}
	if client.inlineAnswerParams == nil {
		t.Fatal("expected inline answer params to be captured")
	}
	if client.inlineAnswerParams.InlineQueryID != "inline-empty" {
		t.Fatalf("expected inline query id to be propagated, got %q", client.inlineAnswerParams.InlineQueryID)
	}
	if !client.inlineAnswerParams.IsPersonal {
		t.Fatal("expected inline answers to be marked personal")
	}
	if len(client.inlineAnswerParams.Results) != 3 {
		t.Fatalf("expected 3 default inline results, got %d", len(client.inlineAnswerParams.Results))
	}

	firstResult, ok := client.inlineAnswerParams.Results[0].(*models.InlineQueryResultArticle)
	if !ok {
		t.Fatalf("expected first result to be an article, got %T", client.inlineAnswerParams.Results[0])
	}
	if firstResult.Title != "Bot Help" {
		t.Fatalf("expected help result first, got %q", firstResult.Title)
	}
}

func TestHandleInlineQueryEchoesQueryText(t *testing.T) {
	client := &stubTelegramClient{}
	service := newTestService(client)

	service.handleInlineQuery(context.Background(), nil, &models.Update{
		InlineQuery: &models.InlineQuery{ID: "inline-query", Query: "hello inline"},
	})

	if client.inlineAnswerParams == nil {
		t.Fatal("expected inline answer params to be captured")
	}
	if len(client.inlineAnswerParams.Results) != 3 {
		t.Fatalf("expected 3 inline results, got %d", len(client.inlineAnswerParams.Results))
	}

	firstResult, ok := client.inlineAnswerParams.Results[0].(*models.InlineQueryResultArticle)
	if !ok {
		t.Fatalf("expected first result to be an article, got %T", client.inlineAnswerParams.Results[0])
	}
	if firstResult.Title != "Echo: hello inline" {
		t.Fatalf("expected echo result title, got %q", firstResult.Title)
	}

	content, ok := firstResult.InputMessageContent.(*models.InputTextMessageContent)
	if !ok {
		t.Fatalf("expected text message content, got %T", firstResult.InputMessageContent)
	}
	if content.MessageText != "hello inline" {
		t.Fatalf("expected echoed text, got %q", content.MessageText)
	}
}

func newTestService(client telegramClient) *Service {
	return &Service{
		client:             client,
		mode:               config.BotModePolling,
		webhookPath:        "/telegram/webhook",
		webhookDropPending: false,
		sessionTTL:         24 * time.Hour,
		sessions:           session.NewMemoryStore(),
	}
}
