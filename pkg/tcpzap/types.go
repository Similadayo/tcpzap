// pkg/tcpzap/types.go
package tcpzap

import (
	"context"
	"time"

	"github.com/similadayo/tcpzap/internal/metrics"
)

// Config holds client and server options.
type Config struct {
	Timeout     time.Duration         // Connection/send timeout
	Retries     int                   // Max retry attempts
	RetryDelay  time.Duration         // Delay between retries
	MetricsFunc func(metrics.Metrics) // Callback for metrics
}

// DefaultConfig provides sane defaults.
func DefaultConfig() Config {
	return Config{
		Timeout:    5 * time.Second,
		Retries:    3,
		RetryDelay: 100 * time.Millisecond,
	}
}

// Handler processes incoming messages.
type Handler interface {
	Handle(ctx context.Context, msg []byte) ([]byte, error)
}
