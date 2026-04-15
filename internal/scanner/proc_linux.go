//go:build linux

package scanner

import (
	"fmt"
	"os"
)

// procNetPaths returns the paths to /proc/net files for TCP and UDP on Linux.
func procNetPaths() []ProcNetSource {
	return []ProcNetSource{
		{Path: "/proc/net/tcp", Proto: "tcp"},
		{Path: "/proc/net/tcp6", Proto: "tcp6"},
		{Path: "/proc/net/udp", Proto: "udp"},
		{Path: "/proc/net/udp6", Proto: "udp6"},
	}
}

// readProcNet reads the content of a /proc/net file.
func readProcNet(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	return data, nil
}
