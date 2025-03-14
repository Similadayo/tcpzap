package tcpzap

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/similadayo/tcpzap/internal/framing"
	"github.com/similadayo/tcpzap/internal/transport"
)

type Client struct {
	conn *transport.Conn
}

// NewClient creates a new client
func NewClient(addr string, timeout time.Duration) (*Client, error) {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, fmt.Errorf("tcpzap: dial: %w", err)
	}
	return &Client{conn: transport.NewConn(conn, framing.NewCodec())}, nil
}

// Send sends a message to the server
func (c *Client) Send(ctx context.Context, msg []byte) ([]byte, error) {
	if err := c.conn.Send(ctx, msg); err != nil {
		return nil, fmt.Errorf("tcpzap: send: %w", err)
	}
	resp, err := c.conn.Receive(ctx)
	if err != nil {
		return nil, fmt.Errorf("tcpzap: receive: %w", err)
	}
	return resp, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	return c.conn.Close()
}
