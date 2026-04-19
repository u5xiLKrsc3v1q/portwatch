package monitor

import "time"

// AlertHistoryHook records alert events into an AlertHistory after each scan cycle.
type AlertHistoryHook struct {
	history *AlertHistory
}

// NewAlertHistoryHook creates a hook that appends alerts to the given history.
func NewAlertHistoryHook(h *AlertHistory) *AlertHistoryHook {
	return &AlertHistoryHook{history: h}
}

// OnAlert records the alert event into history.
func (h *AlertHistoryHook) OnAlert(event AlertEvent) {
	if !event.HasChanges() {
		return
	}
	h.history.Append(AlertRecord{
		Timestamp: time.Now(),
		Added:     len(event.Added),
		Removed:   len(event.Removed),
		Summary:   event.Summary(),
	})
}
