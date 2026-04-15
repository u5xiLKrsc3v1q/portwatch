package scanner

import "testing"

func TestEntry_Key(t *testing.T) {
	tests := []struct {
		name  string
		entry Entry
		want  string
	}{
		{
			name:  "tcp entry",
			entry: Entry{Protocol: "tcp", LocalAddress: "0.0.0.0:80"},
			want:  "tcp|0.0.0.0:80",
		},
		{
			name:  "udp6 entry",
			entry: Entry{Protocol: "udp6", LocalAddress: ":::53"},
			want:  "udp6|:::53",
		},
		{
			name:  "empty entry",
			entry: Entry{},
			want:  "|",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.entry.Key()
			if got != tt.want {
				t.Errorf("Key() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEntry_Fields(t *testing.T) {
	e := Entry{
		Protocol:    "tcp",
		LocalAddress: "127.0.0.1:8080",
		State:       "LISTEN",
		PID:         1234,
		ProcessName: "nginx",
	}
	if e.Protocol != "tcp" {
		t.Errorf("Protocol = %q", e.Protocol)
	}
	if e.PID != 1234 {
		t.Errorf("PID = %d", e.PID)
	}
	if e.ProcessName != "nginx" {
		t.Errorf("ProcessName = %q", e.ProcessName)
	}
}
