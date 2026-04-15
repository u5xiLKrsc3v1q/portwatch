package config

import (
	"os"
	"time"

	"github.com/pelletier/go-toml/v2"
)

// Config holds all portwatch runtime configuration.
type Config struct {
	Interval      time.Duration `toml:"interval"`
	WebhookURL    string        `toml:"webhook_url"`
	Desktop       bool          `toml:"desktop"`
	AppName       string        `toml:"app_name"`
	BlockedPorts  []int         `toml:"blocked_ports"`
	BlockedAddrs  []string      `toml:"blocked_addrs"`
	WhitelistFile string        `toml:"whitelist_file"`
	BaselineFile  string        `toml:"baseline_file"`
	SaveBaseline  bool          `toml:"save_baseline"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval:     30 * time.Second,
		AppName:      "portwatch",
		BaselineFile: "baseline.json",
		SaveBaseline: false,
	}
}

// Load reads a TOML config file from path and merges it with defaults.
// Missing fields retain their default values.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	var raw struct {
		Interval      string   `toml:"interval"`
		WebhookURL    string   `toml:"webhook_url"`
		Desktop       *bool    `toml:"desktop"`
		AppName       string   `toml:"app_name"`
		BlockedPorts  []int    `toml:"blocked_ports"`
		BlockedAddrs  []string `toml:"blocked_addrs"`
		WhitelistFile string   `toml:"whitelist_file"`
		BaselineFile  string   `toml:"baseline_file"`
		SaveBaseline  *bool    `toml:"save_baseline"`
	}
	if err := toml.Unmarshal(data, &raw); err != nil {
		return cfg, err
	}
	if raw.Interval != "" {
		d, err := time.ParseDuration(raw.Interval)
		if err != nil {
			return cfg, err
		}
		cfg.Interval = d
	}
	if raw.WebhookURL != "" {
		cfg.WebhookURL = raw.WebhookURL
	}
	if raw.Desktop != nil {
		cfg.Desktop = *raw.Desktop
	}
	if raw.AppName != "" {
		cfg.AppName = raw.AppName
	}
	if len(raw.BlockedPorts) > 0 {
		cfg.BlockedPorts = raw.BlockedPorts
	}
	if len(raw.BlockedAddrs) > 0 {
		cfg.BlockedAddrs = raw.BlockedAddrs
	}
	if raw.WhitelistFile != "" {
		cfg.WhitelistFile = raw.WhitelistFile
	}
	if raw.BaselineFile != "" {
		cfg.BaselineFile = raw.BaselineFile
	}
	if raw.SaveBaseline != nil {
		cfg.SaveBaseline = *raw.SaveBaseline
	}
	return cfg, nil
}
