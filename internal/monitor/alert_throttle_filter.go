package monitor

import (
	"github.com/rgzr/portwatch/internal/scanner"
)

// AlertThrottleFilter wraps AlertThrottle to filter scanner.Entry slices.
type AlertThrottleFilter struct {
	throttle *AlertThrottle
}

// NewAlertThrottleFilter creates a filter backed by the given throttle.
func NewAlertThrottleFilter(th *AlertThrottle) *AlertThrottleFilter {
	return &AlertThrottleFilter{throttle: th}
}

// FilterAdded returns only entries that pass the throttle check.
func (f *AlertThrottleFilter) FilterAdded(entries []scanner.Entry) []scanner.Entry {
	if f.throttle == nil {
		return entries
	}
	out := entries[:0:0]
	for _, e := range entries {
		if f.throttle.Allow(e.Key()) {
			out = append(out, e)
		}
	}
	return out
}

// FilterRemoved always passes removed entries through without throttling.
func (f *AlertThrottleFilter) FilterRemoved(entries []scanner.Entry) []scanner.Entry {
	return entries
}
