package scanner

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PortEntry represents a single listening port binding.
type PortEntry struct {
	Protocol string
	LocalAddr string
	Port     uint16
	P      int
}

// Scan reads current TCP and UDP listeners from /proc/net.
func Scan() ([]PortEntry, error) {
	var entries []PortEntry

	for _, proto := range []string{"tcp", "tcp6", "udp", "udp6"} {
		path := fmt/%s", proto)
		pe, err := parseProcNet(path, proto)
		if err != nil {
			continue // non-Linux or file absent
		}
		entries = append(entries, pe...)
	}

	return entries, nil
}

func parseProcNet(path, proto string) ([]PortEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []PortEntry
	scanner := bufio.NewScanner(f)
	scanner.Scan() // skip header

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}
		// state 0A = LISTEN for TCP; UDP is stateless
		state := fields[3]
		if strings.HasPrefix(proto, "tcp") && state != "0A" {
			continue
		}

		local := fields[1]
		port, err := parsePort(local)
		if err != nil {
			continue
		}

		entries = append(entries, PortEntry{
			Protocol: proto,
			LocalAddr: local,
			Port:     port,
		})
	}

	return entries, scanner.Err()
}

func parsePort(hexAddr string) (uint16, error) {
	parts := strings.Split(hexAddr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid address: %s", hexAddr)
	}
	v, err := strconv.ParseUint(parts[1], 16, 16)
	if err != nil {
		return 0, err
	}
	return uint16(v), nil
}
