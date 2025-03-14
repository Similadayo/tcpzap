package congestion

import (
	"context"
	"sync"
	"time"
)

// Config defines congestion control parameters
type Config struct {
	InitialWindow int
	MinWindow     int
	MaxWindow     int
	RTTFactor     float64 // Smoothing factor for RTT
	BackoffFactor float64 //Multiplicative decrease factor (e.g., 0.5 halves window)
}

// DefaultConfig provides sane defaults
func DefaultConfig() Config {
	return Config{
		InitialWindow: 10,
		MinWindow:     1,
		MaxWindow:     100,
		RTTFactor:     0.125, // TCP like smoothing
		BackoffFactor: 0.5,   // Halve window on congestion
	}
}

// Controller manages the congestion window and RTT
type Controller struct {
	cfg       Config
	window    int
	rtt       time.Duration
	unacked   int
	lastSent  time.Time
	mu        sync.Mutex
	congested bool
}

// NewController creates a new controller
func NewController(cfg Config) *Controller {
	return &Controller{
		cfg:    cfg,
		window: cfg.InitialWindow,
		rtt:    100 * time.Millisecond,
	}
}

// CanSend checks if the controller can send a message
func (c *Controller) CanSend(ctx context.Context) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.unacked >= c.window {
		c.congested = true
		return false
	}
	select {
	case <-ctx.Done():
		return false
	default:
		c.unacked++
		c.lastSent = time.Now()
		return true
	}
}

// AckReceived processes an ack
func (c *Controller) AckReceived() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.unacked > 0 {
		c.unacked--
		rttSample := time.Since(c.lastSent)
		c.rtt = time.Duration(float64(c.rtt)*(1-c.cfg.RTTFactor) + float64(rttSample)*c.cfg.RTTFactor)
		c.AdjustWindow()
	}
}

// AdjustWindow adjusts the congestion window
func (c *Controller) AdjustWindow() {
	if c.congested {
		// Multiplicative decrease on congestion
		c.window = int(float64(c.window) * c.cfg.BackoffFactor)
		if c.window < c.cfg.MinWindow {
			c.window = c.cfg.MinWindow
		}
		c.congested = false
	} else if c.unacked < c.window/2 && c.window < c.cfg.MaxWindow {
		// Additive increase on no congestion
		c.window++
	}
}

// Window returns the current window size (for testing or monitoring)
func (c *Controller) Window() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.window
}
