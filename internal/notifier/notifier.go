package notifier

// Notifier is the interface implemented by all notification backends.
type Notifier interface {
	// Send delivers a notification with the given title and message body.
	Send(title, message string) error
}

// Multi fans out a single notification to multiple Notifier implementations.
type Multi struct {
	notifiers []Notifier
}

// NewMulti creates a Multi notifier that broadcasts to all provided notifiers.
func NewMulti(notifiers ...Notifier) *Multi {
	return &Multi{notifiers: notifiers}
}

// Send calls Send on every underlying notifier, collecting any errors.
// All notifiers are attempted even if one fails; the last non-nil error
// is returned so callers know at least one delivery failed.
func (m *Multi) Send(title, message string) error {
	var lastErr error
	for _, n := range m.notifiers {
		if err := n.Send(title, message); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
