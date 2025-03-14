package transport

import (
	"context"
	"net"
	"sync"

	"github.com/similadayo/tcpzap/internal/framing"
)

type Conn struct {
	net.Conn
	codec framing.Codec
	mu    sync.Mutex
}

// NewConn wraps a net.Conn with a framing codec
func NewConn(conn net.Conn, codec framing.Codec) *Conn {
	return &Conn{Conn: conn, codec: codec}
}

// Send writes a message to the connection
func (c *Conn) Send(ctx context.Context, msg []byte) error {
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
	return c.codec.Decode(c.Conn)
}

func (c *Conn) SetDeadline(ctx context.Context) error {
	if deadline, ok := ctx.Deadline(); ok {
		return c.Conn.SetDeadline(deadline)
	}
	return nil
}
