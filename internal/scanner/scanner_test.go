package scanner

import (
	"os"
	"testing"
)

func TestParsePort(t *testing.T) {
	tests := []struct {
		input   string
		expect  uint16
		wantErr bool
	}{
		{"00000000:0050", 80, false},
		{"00000000:01BB", 443, false},
		{"00000000:270F", 9999, false},
		{"badformat", 0, true},
		{"00000000:ZZZZ", 0, true},
	}

	for _, tt := range tests {
		got, err := parsePort(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("parsePort(%q) expected error, got nil", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("parsePort(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if got != tt.expect {
			t.Errorf("parsePort(%q) = %d, want %d", tt.input, got, tt.expect)
		}
	}
}

func TestParseProcNet_MissingFile(t *testing.T) {
	_, err := parseProcNet("/nonexistent/path", "tcp")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestParseProcNet_ValidFile(t *testing.T) {
	// Write a minimal /proc/net/tcp-style fixture
	content := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 00000000:0050 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 12345 1 0000000000000000 100 0 0 10 0
   1: 00000000:01BB 00000000:0000 06 00000000:00000000 00:00000000 00000000     0        0 12346 1 0000000000000000 100 0 0 10 0
`
	tmp, err := os.CreateTemp("", "procnet")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	tmp.WriteString(content)
	tmp.Close()

	entries, err := parseProcNet(tmp.Name(), "tcp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 LISTEN entry, got %d", len(entries))
	}
	if entries[0].Port != 80 {
		t.Errorf("expected port 80, got %d", entries[0].Port)
	}
}
