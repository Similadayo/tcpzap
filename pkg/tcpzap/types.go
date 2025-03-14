package tcpzap

import (
	"context"
	"time"
)

type Config struct {
	Timeout     time.Duration
	Retries     int
	RetryDelay  time.Duration
	MetricsFunc func(Metrics)
}

func DefaultConfig() Config {
	return Config{
		Timeout:    5 * time.Second,
		Retries:    3,
		RetryDelay: 100 * time.Second,
	}
}

type Metrics struct {
	Latency time.Duration
	Target  string
	Success bool
}

type Handler interface {
	Handle(ctx context.Context, msg []byte) ([]byte, error)
}
