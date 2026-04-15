package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// PIDInfo holds process information associated with a port binding.
type PIDInfo struct {
	PID  int
	Name string
}

// LookupPID attempts to find the PID and process name for a given inode.
// It walks /proc/<pid>/fd and matches symlinks to socket:[inode].
func LookupPID(inode uint64) (*PIDInfo, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("read /proc: %w", err)
	}

	target := fmt.Sprintf("socket:[%d]", inode)

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue
		}

		fdDir := fmt.Sprintf("/proc/%d/fd", pid)
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}

		for _, fd := range fds {
			link, err := os.Readlink(filepath.Join(fdDir, fd.Name()))
			if err != nil {
				continue
			}
			if link == target {
				name := readProcessName(pid)
				return &PIDInfo{PID: pid, Name: name}, nil
			}
		}
	}
	return nil, fmt.Errorf("inode %d not found", inode)
}

// readProcessName reads the process name from /proc/<pid>/comm.
func readProcessName(pid int) string {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(data))
}
