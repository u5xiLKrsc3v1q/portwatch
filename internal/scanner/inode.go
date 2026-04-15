package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// InodePIDMap maps socket inodes to PIDs.
type InodePIDMap map[uint64]int

// BuildInodePIDMap scans /proc/<pid>/fd to build a mapping from socket
// inode numbers to the PID that owns them.
func BuildInodePIDMap() (InodePIDMap, error) {
	result := make(InodePIDMap)

	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("reading /proc: %w", err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue // not a PID directory
		}
		fdDir := filepath.Join("/proc", e.Name(), "fd")
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue // process may have exited
		}
		for _, fd := range fds {
			link, err := os.Readlink(filepath.Join(fdDir, fd.Name()))
			if err != nil {
				continue
			}
			inode, ok := parseSocketInode(link)
			if !ok {
				continue
			}
			result[inode] = pid
		}
	}
	return result, nil
}

// parseSocketInode extracts the inode number from a symlink target like
// "socket:[12345]".
func parseSocketInode(link string) (uint64, bool) {
	const prefix = "socket:["
	if !strings.HasPrefix(link, prefix) || !strings.HasSuffix(link, "]") {
		return 0, false
	}
	inodeStr := link[len(prefix) : len(link)-1]
	inode, err := strconv.ParseUint(inodeStr, 10, 64)
	if err != nil {
		return 0, false
	}
	return inode, true
}
