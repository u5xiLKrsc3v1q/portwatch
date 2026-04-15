package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-config-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.ScanInterval != 15*time.Second {
		t.Errorf("expected 15s, got %v", cfg.ScanInterval)
	}
	if cfg.AppName != "portwatch" {
		t.Errorf("expected 'portwatch', got %q", cfg.AppName)
	}
	if cfg.DesktopNotify {
		t.Error("expected desktop_notify to be false by default")
	}
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTempConfig(t, `{
		"scan_interval": "30s",
		"webhook_url": "https://example.com/hook",
		"desktop_notify": true,
		"app_name": "mywatch",
		"log_level": "debug"
	}`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScanInterval != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.ScanInterval)
	}
	if cfg.WebhookURL != "https://example.com/hook" {
		t.Errorf("unexpected webhook_url: %q", cfg.WebhookURL)
	}
	if !cfg.DesktopNotify {
		t.Error("expected desktop_notify true")
	}
	if cfg.AppName != "mywatch" {
		t.Errorf("unexpected app_name: %q", cfg.AppName)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("unexpected log_level: %q", cfg.LogLevel)
	}
}

func TestLoad_Defaults_WhenFieldsMissing(t *testing.T) {
	path := writeTempConfig(t, `{}`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScanInterval != 15*time.Second {
		t.Errorf("expected default 15s, got %v", cfg.ScanInterval)
	}
}

func TestLoad_InvalidDuration(t *testing.T) {
	path := writeTempConfig(t, `{"scan_interval": "not-a-duration"}`)
	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid duration, got nil")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/config.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
