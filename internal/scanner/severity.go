package scanner

// Severity represents the alert level for a port event.
type Severity int

const (
	SeverityLow Severity = iota
	SeverityMedium
	SeverityHigh
)

var severityNames = map[Severity]string{
	SeverityLow:    "low",
	SeverityMedium: "medium",
	SeverityHigh:   "high",
}

func (s Severity) String() string {
	if name, ok := severityNames[s]; ok {
		return name
	}
	return "unknown"
}

// wellKnownPorts are ports that trigger high severity when unexpectedly opened.
var wellKnownPorts = map[uint16]bool{
	22: true, 23: true, 3389: true, 5900: true,
	80: true, 443: true, 8080: true, 8443: true,
}

// privilegedThreshold — ports below this are elevated to medium by default.
const privilegedThreshold uint16 = 1024

// Classifier assigns a Severity to an Entry.
type Classifier struct {
	HighPorts map[uint16]bool
}

// NewClassifier returns a Classifier with default well-known port rules.
func NewClassifier() *Classifier {
	return &Classifier{HighPorts: wellKnownPorts}
}

// Classify returns the severity for the given entry.
func (c *Classifier) Classify(e Entry) Severity {
	if c == nil {
		return SeverityLow
	}
	if c.HighPorts[e.Port] {
		return SeverityHigh
	}
	if e.Port < privilegedThreshold {
		return SeverityMedium
	}
	return SeverityLow
}
