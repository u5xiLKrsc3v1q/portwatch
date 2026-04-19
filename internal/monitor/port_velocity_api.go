package monitor

import (
	"encoding/json"
	"net/http"
)

// PortVelocityAPI exposes velocity data over HTTP.
type PortVelocityAPI struct {
	velocity *PortVelocity
}

// NewPortVelocityAPI creates a new PortVelocityAPI.
func NewPortVelocityAPI(v *PortVelocity) *PortVelocityAPI {
	return &PortVelocityAPI{velocity: v}
}

func (a *PortVelocityAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	snap := a.velocity.Snapshot()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"added":   snap.Added,
		"removed": snap.Removed,
		"net":     snap.Net,
		"window":  snap.Window.String(),
	})
}
