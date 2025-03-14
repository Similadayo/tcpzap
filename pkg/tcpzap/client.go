// pkg/tcpzap/client.go
package tcpzap

import (
	"context"
	"fmt"
	"net"

	"github.com/similadayo/tcpzap/internal/framing"
	"github.com/similadayo/tcpzap/internal/metrics"
	"github.com/similadayo/tcpzap/internal/transport"
)

type Client struct {
	conn *transport.Conn
	addr string
	cfg  Config
}

func NewClient(addr string, cfg Config) (*Client, error) {
	conn, err := net.DialTimeout("tcp", addr, cfg.Timeout)
	if err != nil {
		return nil, fmt.Errorf("tcpzap: dial: %w", err)
	}
	tCfg := transport.Config{
		Retries:    cfg.Retries,
		RetryDelay: cfg.RetryDelay,
	}
	return &Client{
		conn: transport.NewConn(conn, framing.NewCodec(), tCfg),
		addr: addr,
		cfg:  cfg,
	}, nil
}

func (c *Client) Send(ctx context.Context, msg []byte) ([]byte, error) {
	tracker := metrics.NewTracker(c.addr, c.cfg.MetricsFunc)
	defer tracker.Report(true)

	if err := c.conn.Send(ctx, msg); err != nil {
		tracker.Report(false)
		return nil, fmt.Errorf("tcpzap: send: %w", err)
	}
	resp, err := c.conn.Receive(ctx)
	if err != nil {
		tracker.Report(false)
		return nil, fmt.Errorf("tcpzap: receive: %w", err)
	}
	return resp, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
