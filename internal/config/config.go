package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds the application configuration.
type Config struct {
	ScanInterval  time.Duration `json:"scan_interval"`
	WebhookURL    string        `json:"webhook_url"`
	DesktopNotify bool          `json:"desktop_notify"`
	AppName       string        `json:"app_name"`
	LogLevel      string        `json:"log_level"`
}

// scanIntervalRaw is used for JSON unmarshalling of duration as string.
type rawConfig struct {
	ScanInterval  string `json:"scan_interval"`
	WebhookURL    string `json:"webhook_url"`
	DesktopNotify bool   `json:"desktop_notify"`
	AppName       string `json:"app_name"`
	LogLevel      string `json:"log_level"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		ScanInterval:  15 * time.Second,
		DesktopNotify: false,
		AppName:       "portwatch",
		LogLevel:      "info",
	}
}

// Load reads and parses a JSON config file from the given path.
// Missing optional fields fall back to defaults.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var raw rawConfig
	if err := json.NewDecoder(f).Decode(&raw); err != nil {
		return nil, fmt.Errorf("config: decode: %w", err)
	}

	cfg := DefaultConfig()

	if raw.ScanInterval != "" {
		d, err := time.ParseDuration(raw.ScanInterval)
		if err != nil {
			return nil, fmt.Errorf("config: invalid scan_interval %q: %w", raw.ScanInterval, err)
		}
		cfg.ScanInterval = d
	}
	if raw.WebhookURL != "" {
		cfg.WebhookURL = raw.WebhookURL
	}
	cfg.DesktopNotify = raw.DesktopNotify
	if raw.AppName != "" {
		cfg.AppName = raw.AppName
	}
	if raw.LogLevel != "" {
		cfg.LogLevel = raw.LogLevel
	}

	return cfg, nil
}
