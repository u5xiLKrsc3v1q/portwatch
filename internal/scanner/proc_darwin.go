//go:build darwin

package scanner

import (
	"bytes"
	"fmt"
	"os/exec"
)

// ProcNetSource represents a source of port binding information.
type ProcNetSource struct {
	Path  string
	Proto string
}

// procNetPaths returns sources for macOS using netstat output files.
func procNetPaths() []ProcNetSource {
	return []ProcNetSource{
		{Path: "netstat", Proto: "tcp"},
		{Path: "netstat", Proto: "udp"},
	}
}

// readProcNet on Darwin runs netstat and returns its output.
func readProcNet(path string) ([]byte, error) {
	cmd := exec.Command("netstat", "-an", "-p", "tcp")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("running netstat: %w", err)
	}
	return out.Bytes(), nil
}
