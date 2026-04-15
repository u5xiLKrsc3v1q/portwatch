//go:build !linux && !darwin

package scanner

import "errors"

// ProcNetSource represents a source of port binding information.
type ProcNetSource struct {
	Path  string
	Proto string
}

// procNetPaths returns an empty list on unsupported platforms.
func procNetPaths() []ProcNetSource {
	return nil
}

// readProcNet returns an error on unsupported platforms.
func readProcNet(_ string) ([]byte, error) {
	return nil, errors.New("unsupported platform: cannot read proc net")
}
