package transport

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/similadayo/tcpzap/internal/congestion"
	"github.com/similadayo/tcpzap/internal/framing"
)

type Config struct {
	Retries    int
	RetryDelay time.Duration
}

type Conn struct {
	net.Conn
	codec framing.Codec
	mu    sync.Mutex
	ctrl  *congestion.Controller
	cfg   Config
}

// NewConn wraps a net.Conn with a framing codec
func NewConn(conn net.Conn, codec framing.Codec, cfg Config) *Conn {
	return &Conn{
		Conn:  conn,
		codec: codec,
		ctrl:  congestion.NewController(congestion.DefaultConfig()),
		cfg:   cfg,
	}
}

// Send writes a message to the connection
func (c *Conn) Send(ctx context.Context, msg []byte) error {
	var lastErr error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if !c.ctrl.CanSend(ctx) {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(c.ctrl.RTT() / 10):
				continue
			}
		}
		c.mu.Lock()
		if err := c.SetDeadline(ctx); err != nil {
			c.mu.Unlock()
			return err
		}
		err := c.codec.Encode(c.Conn, msg)
		c.mu.Unlock()
		if err == nil {
			return nil
		}
		lastErr = err
		if attempt < c.cfg.Retries {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(c.cfg.RetryDelay):
				continue
			}

		}
	}
	return lastErr
}

// Receive reads a message from the connection
func (c *Conn) Receive(ctx context.Context) ([]byte, error) {
	if err := c.SetDeadline(ctx); err != nil {
		return nil, err
	}
	data, err := c.codec.Decode(c.Conn)
	if err == nil {
		c.ctrl.AckReceived()
	}
	return data, err
}

func (c *Conn) SetDeadline(ctx context.Context) error {
	if deadline, ok := ctx.Deadline(); ok {
		return c.Conn.SetDeadline(deadline)
	}
	return nil
}

// RTT exposes the current RTT
func (c *Conn) RTT() time.Duration {
	return c.ctrl.RTT()
}
