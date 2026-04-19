package monitor

import "net/http"

// NewRouter builds the HTTP mux wiring all API handlers.
func NewRouter(
	events *EventLog,
	metrics *Metrics,
	history *ScanHistory,
	alertHistory *AlertHistory,
	trend *PortTrend,
	snapshot *PortSnapshotStore,
	processMap *ProcessMapStore,
	topPorts *TopPortsStore,
	velocity *PortVelocity,
	age *PortAgeStore,
	risk *PortRiskStore,
) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/health", NewHTTPAPI(events))
	mux.Handle("/events", NewHTTPAPI(events))
	mux.Handle("/metrics", NewMetricsAPI(metrics))
	mux.Handle("/scans", NewScanHistoryAPI(history))
	mux.Handle("/alerts", NewAlertHistoryAPI(alertHistory))
	mux.Handle("/trends", NewPortTrendAPI(trend))
	mux.Handle("/snapshot", NewPortSnapshotAPI(snapshot))
	mux.Handle("/processes", NewProcessMapAPI(processMap))
	mux.Handle("/top", NewTopPortsAPI(topPorts))
	mux.Handle("/velocity", NewPortVelocityAPI(velocity))
	mux.Handle("/age", NewPortAgeAPI(age))
	mux.Handle("/risk", NewPortRiskAPI(risk))
	return mux
}
