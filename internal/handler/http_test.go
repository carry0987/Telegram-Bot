package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubReadiness struct {
	err error
}

func (s stubReadiness) Ping(context.Context) error {
	return s.err
}

func TestHealthz(t *testing.T) {
	h := New(Options{Readiness: stubReadiness{}})
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestReadyzReturnsServiceUnavailableWhenProbeFails(t *testing.T) {
	h := New(Options{Readiness: stubReadiness{err: errors.New("telegram unavailable")}})
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", recorder.Code)
	}
}

func TestReadyzReturnsOKWhenProbeSucceeds(t *testing.T) {
	h := New(Options{Readiness: stubReadiness{}})
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
}

func TestRegisterAddsWebhookRouteWhenConfigured(t *testing.T) {
	h := New(Options{
		Readiness:      stubReadiness{},
		WebhookPath:    "/telegram/webhook",
		WebhookHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusAccepted) }),
	})
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodPost, "/telegram/webhook", nil)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusAccepted {
		t.Fatalf("expected status 202, got %d", recorder.Code)
	}
}
