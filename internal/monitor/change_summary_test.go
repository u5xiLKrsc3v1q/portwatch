package monitor

import (
	"strings"
	"testing"
	"time"

	"github.com/jwhittle933/portwatch/internal/scanner"
)

func makeSummaryEntry(port uint16, addr string) scanner.Entry {
	return scanner.Entry{
		LocalAddress: addr,
		LocalPort:    port,
		Protocol:     scanner.TCP,
	}
}

func TestChangeSummary_HasChanges_Empty(t *testing.T) {
	s := NewChangeSummary(nil, nil, time.Minute)
	if s.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestChangeSummary_HasChanges_WithAdded(t *testing.T) {
	added := []scanner.Entry{makeSummaryEntry(8080, "0.0.0.0")}
	s := NewChangeSummary(added, nil, time.Minute)
	if !s.HasChanges() {
		t.Error("expected changes")
	}
}

func TestChangeSummary_HasChanges_WithRemoved(t *testing.T) {
	removed := []scanner.Entry{makeSummaryEntry(443, "0.0.0.0")}
	s := NewChangeSummary(nil, removed, time.Minute)
	if !s.HasChanges() {
		t.Error("expected changes")
	}
}

func TestChangeSummary_TotalChanges(t *testing.T) {
	added := []scanner.Entry{makeSummaryEntry(80, "0.0.0.0"), makeSummaryEntry(443, "0.0.0.0")}
	removed := []scanner.Entry{makeSummaryEntry(8080, "127.0.0.1")}
	s := NewChangeSummary(added, removed, time.Minute)
	if got := s.TotalChanges(); got != 3 {
		t.Errorf("expected 3, got %d", got)
	}
}

func TestChangeSummary_String_ContainsLabels(t *testing.T) {
	added := []scanner.Entry{makeSummaryEntry(9090, "0.0.0.0")}
	removed := []scanner.Entry{makeSummaryEntry(22, "0.0.0.0")}
	s := NewChangeSummary(added, removed, 30*time.Second)
	str := s.String()
	if !strings.Contains(str, "Added") {
		t.Error("expected 'Added' in summary")
	}
	if !strings.Contains(str, "Removed") {
		t.Error("expected 'Removed' in summary")
	}
	if !strings.Contains(str, "30s") {
		t.Error("expected window duration in summary")
	}
}

func TestChangeSummary_String_Empty(t *testing.T) {
	s := NewChangeSummary(nil, nil, time.Minute)
	str := s.String()
	if !strings.Contains(str, "ChangeSummary") {
		t.Error("expected header in empty summary")
	}
}
