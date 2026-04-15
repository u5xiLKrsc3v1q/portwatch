package notifier

import (
	"runtime"
	"testing"
)

func TestNewDesktopNotifier_DefaultAppName(t *testing.T) {
	n := NewDesktopNotifier("")
	if n.AppName != "portwatch" {
		t.Errorf("expected default app name 'portwatch', got %q", n.AppName)
	}
}

func TestNewDesktopNotifier_CustomAppName(t *testing.T) {
	n := NewDesktopNotifier("myapp")
	if n.AppName != "myapp" {
		t.Errorf("expected app name 'myapp', got %q", n.AppName)
	}
}

// TestDesktopNotifier_Send_UnsupportedPlatform verifies that Send returns an
// error on platforms where no notification backend is available. We simulate
// this by temporarily patching the GOOS — instead we test the error path
// directly via a helper that accepts an explicit OS string.
func TestDesktopNotifier_Send_UnsupportedPlatform(t *testing.T) {
	n := NewDesktopNotifier("portwatch")
	err := n.sendForOS("plan9", "title", "msg")
	if err == nil {
		t.Error("expected error for unsupported platform, got nil")
	}
}

func TestDesktopNotifier_Send_Linux_MissingBinary(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("linux-only test")
	}
	n := NewDesktopNotifier("portwatch")
	// This may succeed if notify-send is installed; we only check no panic.
	_ = n.sendLinux("Test Title", "Test Message")
}

func TestDesktopNotifier_Send_Darwin_Script(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("darwin-only test")
	}
	n := NewDesktopNotifier("portwatch")
	err := n.sendDarwin("Test Title", "Test Message")
	if err != nil {
		t.Logf("sendDarwin returned error (may be expected in CI): %v", err)
	}
}
