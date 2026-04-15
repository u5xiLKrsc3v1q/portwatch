//go:build !linux && !darwin

package scanner

import "testing"

func TestProcNetPaths_Other(t *testing.T) {
	sources := procNetPaths()
	if len(sources) != 0 {
		t.Errorf("expected empty proc net paths on unsupported platform, got %d", len(sources))
	}
}

func TestReadProcNet_Other(t *testing.T) {
	_, err := readProcNet("anything")
	if err == nil {
		t.Fatal("expected error on unsupported platform")
	}
}
