package scanner

import (
	"testing"
)

func TestSeverity_String(t *testing.T) {
	cases := []struct {
		s    Severity
		want string
	}{
		{SeverityLow, "low"},
		{SeverityMedium, "medium"},
		{SeverityHigh, "high"},
		{Severity(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("Severity(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}

func TestClassifier_Nil(t *testing.T) {
	var c *Classifier
	e := Entry{Port: 22}
	if got := c.Classify(e); got != SeverityLow {
		t.Errorf("nil classifier should return low, got %s", got)
	}
}

func TestClassifier_WellKnownPort_High(t *testing.T) {
	c := NewClassifier()
	for _, port := range []uint16{22, 23, 80, 443, 3389} {
		e := Entry{Port: port}
		if got := c.Classify(e); got != SeverityHigh {
			t.Errorf("port %d: want high, got %s", port, got)
		}
	}
}

func TestClassifier_PrivilegedPort_Medium(t *testing.T) {
	c := NewClassifier()
	// Port 512 is privileged but not in well-known list
	e := Entry{Port: 512}
	if got := c.Classify(e); got != SeverityMedium {
		t.Errorf("port 512: want medium, got %s", got)
	}
}

func TestClassifier_HighPort_Low(t *testing.T) {
	c := NewClassifier()
	e := Entry{Port: 49152}
	if got := c.Classify(e); got != SeverityLow {
		t.Errorf("port 49152: want low, got %s", got)
	}
}

func TestClassifier_CustomHighPort(t *testing.T) {
	c := NewClassifier()
	c.HighPorts[9999] = true
	e := Entry{Port: 9999}
	if got := c.Classify(e); got != SeverityHigh {
		t.Errorf("custom port 9999: want high, got %s", got)
	}
}
