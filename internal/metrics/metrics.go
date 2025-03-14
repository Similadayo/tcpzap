// internal/metrics/metrics.go
package metrics

import "time"

// Metrics captures performance data.
type Metrics struct {
	Latency time.Duration
	Target  string
	Success bool
}

// ReportFunc is a callback to report metrics.
type ReportFunc func(Metrics)

// Tracker measures operation latency.
type Tracker struct {
	start  time.Time
	target string
	report ReportFunc
}

// NewTracker starts a metrics tracker.
func NewTracker(target string, report ReportFunc) *Tracker {
	return &Tracker{
		start:  time.Now(),
		target: target,
		report: report,
	}
}

// Report calculates latency and calls the callback.
func (t *Tracker) Report(success bool) {
	if t.report == nil {
		return
	}
	latency := time.Since(t.start)
	t.report(Metrics{
		Latency: latency,
		Target:  t.target,
		Success: success,
	})
}
