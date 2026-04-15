package monitor

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeAlertEvent(added, removed []scanner.Entry) AlertEvent {
	return AlertEvent{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Added:     added,
		Removed:   removed,
	}
}

func TestAlertEvent_HasChanges_Empty(t *testing.T) {
	e := makeAlertEvent(nil, nil)
	if e.HasChanges() {
		t.Error("expected no changes for empty event")
	}
}

func TestAlertEvent_HasChanges_WithAdded(t *testing.T) {
	e := makeAlertEvent([]scanner.Entry{{Port: 8080}}, nil)
	if !e.HasChanges() {
		t.Error("expected changes when entries added")
	}
}

func TestAlertEvent_HasChanges_WithRemoved(t *testing.T) {
	e := makeAlertEvent(nil, []scanner.Entry{{Port: 9090}})
	if !e.HasChanges() {
		t.Error("expected changes when entries removed")
	}
}

func TestAlertEvent_Summary_ContainsPort(t *testing.T) {
	e := makeAlertEvent([]scanner.Entry{{Port: 8080, Protocol: "tcp", Address: "0.0.0.0"}}, nil)
	summary := e.Summary()
	if !strings.Contains(summary, "ADDED") {
		t.Errorf("expected summary to contain ADDED, got: %s", summary)
	}
	if !strings.Contains(summary, "8080") {
		t.Errorf("expected summary to contain port 8080, got: %s", summary)
	}
}

func TestAlertEvent_Title_AddedOnly(t *testing.T) {
	e := makeAlertEvent([]scanner.Entry{{Port: 8080}}, nil)
	title := e.Title()
	if !strings.Contains(title, "new") {
		t.Errorf("expected title to contain 'new', got: %s", title)
	}
}

func TestAlertEvent_Title_BothAddedAndRemoved(t *testing.T) {
	e := makeAlertEvent([]scanner.Entry{{Port: 8080}}, []scanner.Entry{{Port: 9090}})
	title := e.Title()
	if !strings.Contains(title, "new") || !strings.Contains(title, "removed") {
		t.Errorf("expected title to contain both 'new' and 'removed', got: %s", title)
	}
}
