package notifier

import (
	"fmt"

	"github.com/user/portwatch/internal/config"
)

// FromConfig builds a Notifier from the application configuration.
// If multiple backends are enabled the result is wrapped in a Multi notifier.
// Returns an error if no backends are enabled or a backend cannot be created.
func FromConfig(cfg *config.Config) (Notifier, error) {
	var notifiers []Notifier

	if cfg.Webhook.URL != "" {
		wh, err := NewWebhookNotifier(cfg.Webhook.URL)
		if err != nil {
			return nil, fmt.Errorf("webhook notifier: %w", err)
		}
		notifiers = append(notifiers, wh)
	}

	if cfg.Desktop.Enabled {
		appName := cfg.Desktop.AppName
		if appName == "" {
			appName = "portwatch"
		}
		notifiers = append(notifiers, NewDesktopNotifier(appName))
	}

	if len(notifiers) == 0 {
		return nil, fmt.Errorf("no notifier configured: set webhook.url or desktop.enabled in config")
	}

	if len(notifiers) == 1 {
		return notifiers[0], nil
	}

	return NewMulti(notifiers...), nil
}
