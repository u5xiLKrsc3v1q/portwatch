package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Event represents a port binding change to be reported.
type Event struct {
	Type    string `json:"type"` // "added" or "removed"
	Proto   string `json:"proto"`
	Port    uint16 `json:"port"`
	PID     int    `json:"pid"`
	Time    string `json:"time"`
}

// WebhookNotifier sends port change events to an HTTP endpoint.
type WebhookNotifier struct {
	URL    string
	Client *http.Client
}

// NewWebhookNotifier creates a WebhookNotifier with a default timeout.
func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{
		URL: url,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send serializes the event as JSON and POSTs it to the configured URL.
func (w *WebhookNotifier) Send(event Event) error {
	if event.Time == "" {
		event.Time = time.Now().UTC().Format(time.RFC3339)
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("webhook: marshal event: %w", err)
	}

	resp, err := w.Client.Post(w.URL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("webhook: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}

	return nil
}
