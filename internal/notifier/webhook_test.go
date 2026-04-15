package notifier

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhookNotifier_Send_Success(t *testing.T) {
	var received Event

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &received); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := NewWebhookNotifier(server.URL)
	event := Event{
		Type:  "added",
		Proto: "tcp",
		Port:  8080,
		PID:   1234,
	}

	if err := n.Send(event); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}
	if received.Port != 8080 {
		t.Errorf("expected port 8080, got %d", received.Port)
	}
	if received.Type != "added" {
		t.Errorf("expected type 'added', got %s", received.Type)
	}
	if received.Time == "" {
		t.Error("expected Time to be set automatically")
	}
}

func TestWebhookNotifier_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := NewWebhookNotifier(server.URL)
	err := n.Send(Event{Type: "removed", Proto: "udp", Port: 53})
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestWebhookNotifier_Send_BadURL(t *testing.T) {
	n := NewWebhookNotifier("http://127.0.0.1:0/nonexistent")
	err := n.Send(Event{Type: "added", Proto: "tcp", Port: 9999})
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}
