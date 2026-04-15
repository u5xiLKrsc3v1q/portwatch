package scanner

// Enricher attaches process information (PID and process name) to scan entries
// using inode-to-PID mapping and /proc filesystem lookups.
type Enricher struct {
	inodePIDMap InodePIDMap
}

// NewEnricher creates an Enricher populated with a fresh inode→PID map.
// Errors building the map are non-fatal; the enricher will simply leave
// PID/Process fields at their zero values.
func NewEnricher() *Enricher {
	m, _ := BuildInodePIDMap()
	return &Enricher{inodePIDMap: m}
}

// NewEnricherWithMap creates an Enricher using a pre-built InodePIDMap.
// Useful for testing or when the caller already holds a fresh map.
func NewEnricherWithMap(m InodePIDMap) *Enricher {
	return &Enricher{inodePIDMap: m}
}

// Enrich fills the PID and Process fields of each Entry in the slice.
// Entries whose inode cannot be resolved are left unchanged.
func (e *Enricher) Enrich(entries []Entry) {
	for i := range entries {
		pid, ok := e.inodePIDMap[entries[i].Inode]
		if !ok {
			continue
		}
		entries[i].PID = pid
		name, err := readProcessName(pid)
		if err == nil {
			entries[i].Process = name
		}
	}
}

// EnrichOne fills the PID and Process fields of a single Entry.
func (e *Enricher) EnrichOne(entry *Entry) {
	pid, ok := e.inodePIDMap[entry.Inode]
	if !ok {
		return
	}
	entry.PID = pid
	if name, err := readProcessName(pid); err == nil {
		entry.Process = name
	}
}
