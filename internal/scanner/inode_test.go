package scanner

import (
	"testing"
)

func TestParseSocketInode_Valid(t *testing.T) {
	tests := []struct {
		link  string
		want  uint64
		ok    bool
	}{
		{"socket:[12345]", 12345, true},
		{"socket:[0]", 0, true},
		{"socket:[9999999]", 9999999, true},
	}
	for _, tc := range tests {
		got, ok := parseSocketInode(tc.link)
		if ok != tc.ok {
			t.Errorf("parseSocketInode(%q) ok=%v, want %v", tc.link, ok, tc.ok)
		}
		if got != tc.want {
			t.Errorf("parseSocketInode(%q) = %d, want %d", tc.link, got, tc.want)
		}
	}
}

func TestParseSocketInode_Invalid(t *testing.T) {
	cases := []string{
		"/dev/null",
		"pipe:[999]",
		"socket:[abc]",
		"socket:[",
		"",
		"socket:[]notclosed",
	}
	for _, link := range cases {
		_, ok := parseSocketInode(link)
		if ok {
			t.Errorf("parseSocketInode(%q) expected false, got true", link)
		}
	}
}

func TestBuildInodePIDMap_ReturnsMap(t *testing.T) {
	// This test runs against the real /proc filesystem.
	// On non-Linux systems it will simply return an error, which is acceptable.
	m, err := BuildInodePIDMap()
	if err != nil {
		// Non-Linux or permission issue — acceptable in CI.
		t.Skipf("BuildInodePIDMap error (may be non-Linux): %v", err)
	}
	// The map may be empty if no sockets are open, but it should not be nil.
	if m == nil {
		t.Fatal("expected non-nil InodePIDMap")
	}
}

func TestInodePIDMap_Type(t *testing.T) {
	m := make(InodePIDMap)
	m[1234] = 42
	if m[1234] != 42 {
		t.Errorf("expected pid 42, got %d", m[1234])
	}
	if _, exists := m[9999]; exists {
		t.Error("unexpected entry for inode 9999")
	}
}

func TestInodePIDMap_Overwrite(t *testing.T) {
	// Verify that inserting a new PID for an existing inode overwrites the old value.
	m := make(InodePIDMap)
	m[1234] = 10
	m[1234] = 20
	if m[1234] != 20 {
		t.Errorf("expected pid 20 after overwrite, got %d", m[1234])
	}
}
