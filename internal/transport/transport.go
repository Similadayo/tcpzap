package transport

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/similadayo/tcpzap/internal/congestion"
	"github.com/similadayo/tcpzap/internal/framing"
)

type Conn struct {
	net.Conn
	codec framing.Codec
	mu    sync.Mutex
	ctrl  *congestion.Controller
}

// NewConn wraps a net.Conn with a framing codec
func NewConn(conn net.Conn, codec framing.Codec) *Conn {
	return &Conn{
		Conn:  conn,
		codec: codec,
		ctrl:  congestion.NewController(congestion.DefaultConfig()),
	}
}

// Send writes a message to the connection
func (c *Conn) Send(ctx context.Context, msg []byte) error {
	for !c.ctrl.CanSend(ctx) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(c.ctrl.RTT() / 10):
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.codec.Encode(c.Conn, msg); err != nil {
		return err
	}
	return c.codec.Encode(c.Conn, msg)
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
