package metrics

import (
	"time"

	"github.com/similadayo/tcpzap/pkg/tcpzap"
)

// Tracker measures operation latencies
type Tracker struct {
	start  time.Time
	cfg    *tcpzap.Config
	target string
}

// NewTracker creates a new tracker
func NewTracker(cfg *tcpzap.Config, target string) *Tracker {
	return &Tracker{
		start:  time.Now(),
		cfg:    cfg,
		target: target,
	}
}

// Reports calculates latency and calls metric callback
func (t *Tracker) Report(success bool) {
	if t.cfg == nil || t.cfg.MetricsFunc == nil {
		return
	}
	latency := time.Since(t.start)
	t.cfg.MetricsFunc(tcpzap.Metrics{
		Latency: latency,
		Target:  t.target,
		Success: success,
	})
}
