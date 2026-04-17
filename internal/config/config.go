package config

import (
	"errors"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// Config holds all portwatch configuration.
type Config struct {
	Interval      time.Duration `toml:"interval"`
	BaselineFile  string        `toml:"baseline_file"`
	WebhookURL    string        `toml:"webhook_url"`
	Desktop       bool          `toml:"desktop"`
	BlockedPorts  []int         `toml:"blocked_ports"`
	BlockedAddrs  []string      `toml:"blocked_addrs"`
	RateCooldown  time.Duration `toml:"rate_cooldown"`
	DebounceWindow time.Duration `toml:"debounce_window"`
	HTTPAddr      string        `toml:"http_addr"`
	MaxEvents     int           `toml:"max_events"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval:       15 * time.Second,
		BaselineFile:   "baseline.json",
		RateCooldown:   60 * time.Second,
		DebounceWindow: 5 * time.Second,
		HTTPAddr:       "",
		MaxEvents:      200,
	}
}

type rawConfig struct {
	Interval       string   `toml:"interval"`
	BaselineFile   string   `toml:"baseline_file"`
	WebhookURL     string   `toml:"webhook_url"`
	Desktop        bool     `toml:"desktop"`
	BlockedPorts   []int    `toml:"blocked_ports"`
	BlockedAddrs   []string `toml:"blocked_addrs"`
	RateCooldown   string   `toml:"rate_cooldown"`
	DebounceWindow string   `toml:"debounce_window"`
	HTTPAddr       string   `toml:"http_addr"`
	MaxEvents      int      `toml:"max_events"`
}

// Load reads a TOML config file, falling back to defaults for missing fields.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, err
	}
	var raw rawConfig
	if _, err := toml.Decode(string(data), &raw); err != nil {
		return cfg, err
	}
	if raw.Interval != "" {
		d, err := time.ParseDuration(raw.Interval)
		if err != nil {
			return cfg, err
		}
		cfg.Interval = d
	}
	if raw.RateCooldown != "" {
		d, err := time.ParseDuration(raw.RateCooldown)
		if err != nil {
			return cfg, err
		}
		cfg.RateCooldown = d
	}
	if raw.DebounceWindow != "" {
		d, err := time.ParseDuration(raw.DebounceWindow)
		if err != nil {
			return cfg, err
		}
		cfg.DebounceWindow = d
	}
	if raw.BaselineFile != "" { cfg.BaselineFile = raw.BaselineFile }
	if raw.WebhookURL != "" { cfg.WebhookURL = raw.WebhookURL }
	if raw.HTTPAddr != "" { cfg.HTTPAddr = raw.HTTPAddr }
	if raw.MaxEvents != 0 { cfg.MaxEvents = raw.MaxEvents }
	cfg.Desktop = raw.Desktop
	cfg.BlockedPorts = raw.BlockedPorts
	cfg.BlockedAddrs = raw.BlockedAddrs
	return cfg, nil
}
