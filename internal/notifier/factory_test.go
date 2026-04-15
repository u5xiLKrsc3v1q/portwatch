package notifier_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/notifier"
)

func TestFromConfig_NoNotifiers(t *testing.T) {
	cfg := &config.Config{
		WebhookURL: "",
		DesktopNotify: false,
	}
	n := notifier.FromConfig(cfg)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestFromConfig_WebhookOnly(t *testing.T) {
	cfg := &config.Config{
		WebhookURL:    "http://example.com/hook",
		DesktopNotify: false,
	}
	n := notifier.FromConfig(cfg)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestFromConfig_DesktopOnly(t *testing.T) {
	cfg := &config.Config{
		WebhookURL:    "",
		DesktopNotify: true,
	}
	n := notifier.FromConfig(cfg)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestFromConfig_BothNotifiers(t *testing.T) {
	cfg := &config.Config{
		WebhookURL:    "http://example.com/hook",
		DesktopNotify: true,
	}
	n := notifier.FromConfig(cfg)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
