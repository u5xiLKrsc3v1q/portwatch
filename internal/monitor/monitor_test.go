package monitor

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

type mockNotifier struct {
	calls []string
}

func (m *mockNotifier) Send(title, message string) error {
	m.calls = append(m.calls, title+":"+message)
	return nil
}

func TestNew_Defaults(t *testing.T) {
	mon := New(time.Second, nil)
	if mon == nil {
		t.Fatal("expected non-nil monitor")
	}
	if mon.interval != time.Second {
		t.Errorf("expected interval 1s, got %v", mon.interval)
	}
	if mon.logger == nil {
		t.Error("expected default logger")
	}
}

func TestNew_WithNotifiers(t *testing.T) {
	n := &mockNotifier{}
	logger := log.New(os.Discard, "", 0)
	mon := New(time.Second, logger, n)
	if len(mon.notifiers) != 1 {
		t.Errorf("expected 1 notifier, got %d", len(mon.notifiers))
	}
}

func TestFormatEntry(t *testing.T) {
	e := scanner.Entry{Protocol: "tcp", LocalAddress: "0.0.0.0:8080"}
	got := formatEntry("new listener detected", e)
	want := "new listener detected: tcp port 0.0.0.0:8080"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRun_CancelImmediately(t *testing.T) {
	logger := log.New(os.Discard, "", 0)
	mon := New(time.Minute, logger)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := mon.Run(ctx)
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestNotify_CallsAllNotifiers(t *testing.T) {
	n1 := &mockNotifier{}
	n2 := &mockNotifier{}
	logger := log.New(os.Discard, "", 0)
	mon := New(time.Second, logger, n1, n2)
	mon.notify("Title", "Message")
	if len(n1.calls) != 1 {
		t.Errorf("n1: expected 1 call, got %d", len(n1.calls))
	}
	if len(n2.calls) != 1 {
		t.Errorf("n2: expected 1 call, got %d", len(n2.calls))
	}
}
