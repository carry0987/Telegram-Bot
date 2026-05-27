package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type ReadinessChecker interface {
	Ping(ctx context.Context) error
}

type Handler struct {
	readiness      ReadinessChecker
	mode           string
	sessionBackend string
	webhookPath    string
	webhookHandler http.Handler
}

type Options struct {
	Readiness      ReadinessChecker
	Mode           string
	SessionBackend string
	WebhookPath    string
	WebhookHandler http.Handler
}

type healthResponse struct {
	OK      bool   `json:"ok"`
	Status  string `json:"status"`
	Service string `json:"service"`
}

type readinessResponse struct {
	OK      bool   `json:"ok"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func New(options Options) *Handler {
	return &Handler{
		readiness:      options.Readiness,
		mode:           options.Mode,
		sessionBackend: options.SessionBackend,
		webhookPath:    options.WebhookPath,
		webhookHandler: options.WebhookHandler,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	if h.webhookHandler != nil && h.webhookPath != "" {
		mux.Handle(h.webhookPath, h.webhookHandler)
	}

	mux.HandleFunc("/", h.root)
	mux.HandleFunc("/healthz", h.healthz)
	mux.HandleFunc("/readyz", h.readyz)
}

func (h *Handler) root(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"name":        "telegram-bot",
		"description": "Go Telegram Bot with health checks, polling, webhook, and session demos",
		"mode":        h.mode,
		"session":     h.sessionBackend,
	})
}

func (h *Handler) healthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, healthResponse{
		OK:      true,
		Status:  "healthy",
		Service: "telegram-bot",
	})
}

func (h *Handler) readyz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if h.readiness == nil {
		writeJSON(w, http.StatusServiceUnavailable, readinessResponse{
			OK:      false,
			Status:  "not_ready",
			Message: "readiness checker is not configured",
		})

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.readiness.Ping(ctx); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, readinessResponse{
			OK:      false,
			Status:  "not_ready",
			Message: err.Error(),
		})

		return
	}

	writeJSON(w, http.StatusOK, readinessResponse{
		OK:      true,
		Status:  "ready",
		Message: "telegram bot client is ready",
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
