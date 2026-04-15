//go:build linux

package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcNetPaths_Linux(t *testing.T) {
	sources := procNetPaths()
	if len(sources) == 0 {
		t.Fatal("expected non-empty proc net paths on Linux")
	}
	protos := map[string]bool{}
	for _, s := range sources {
		protos[s.Proto] = true
	}
	for _, expected := range []string{"tcp", "tcp6", "udp", "udp6"} {
		if !protos[expected] {
			t.Errorf("expected proto %q in sources", expected)
		}
	}
}

func TestReadProcNet_MissingFile(t *testing.T) {
	_, err := readProcNet("/nonexistent/path/to/file")
	if err == nil {
		t.Fatal("expected error reading missing file")
	}
}

func TestReadProcNet_ValidFile(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "tcp")
	content := []byte("sl local_address rem_address st\n")
	if err := os.WriteFile(p, content, 0o644); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	data, err := readProcNet(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != string(content) {
		t.Errorf("expected %q, got %q", content, data)
	}
}
