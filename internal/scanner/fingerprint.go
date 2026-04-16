package scanner

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

// Fingerprint represents a stable hash of a set of port entries.
type Fingerprint struct {
	Hash  string
	Count int
}

// NewFingerprint computes a deterministic fingerprint from a slice of entries.
// Entries are sorted by key before hashing to ensure stability.
func NewFingerprint(entries []Entry) Fingerprint {
	if len(entries) == 0 {
		return Fingerprint{Hash: emptySHA256(), Count: 0}
	}

	keys := make([]string, 0, len(entries))
	for _, e := range entries {
		keys = append(keys, e.Key())
	}
	sort.Strings(keys)

	h := sha256.New()
	h.Write([]byte(strings.Join(keys, "\n")))
	return Fingerprint{
		Hash:  fmt.Sprintf("%x", h.Sum(nil)),
		Count: len(entries),
	}
}

// Equal returns true if two fingerprints match.
func (f Fingerprint) Equal(other Fingerprint) bool {
	return f.Hash == other.Hash && f.Count == other.Count
}

// String returns a short human-readable representation.
func (f Fingerprint) String() string {
	short := f.Hash
	if len(short) > 12 {
		short = short[:12]
	}
	return fmt.Sprintf("%s (%d entries)", short, f.Count)
}

func emptySHA256() string {
	h := sha256.New()
	return fmt.Sprintf("%x", h.Sum(nil))
}
