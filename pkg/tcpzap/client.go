package tcpzap

import (
	"context"
	"fmt"
	"net"
	"time"

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
		return nil, fmt.Errorf("tcpzap: failed to connect to %s: %w", addr, err)
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
	tracker := metrics.NewTracker(c.addr, metrics.ReportFunc(c.cfg.MetricsFunc))
	defer tracker.Report(true)

	if err := c.conn.Send(ctx, msg); err != nil {
		tracker.Report(false)
		return nil, fmt.Errorf("tcpzap: failed to send to %s after %d retries: %w", c.addr, c.cfg.Retries, err)
	}

	var resp []byte
	var err error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		resp, err = c.conn.Receive(ctx)
		if err == nil {
			break
		}
		isTemp := false
		if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
			isTemp = true
		} else if oe, ok := err.(*net.OpError); ok && oe.Err != nil {
			if te, ok := oe.Err.(interface{ Temporary() bool }); ok && te.Temporary() {
				isTemp = true
			}
		}
		if !isTemp {
			tracker.Report(false)
			return nil, fmt.Errorf("tcpzap: failed to receive from %s: %w", c.addr, err)
		}
		if attempt < c.cfg.Retries {
			time.Sleep(c.cfg.RetryDelay)
			// Resend the message since the previous receive failed
			if err := c.conn.Send(ctx, msg); err != nil {
				tracker.Report(false)
				return nil, fmt.Errorf("tcpzap: resend failed to %s: %w", c.addr, err)
			}
			continue
		}
		tracker.Report(false)
		return nil, fmt.Errorf("tcpzap: receive retries exhausted for %s: %w", c.addr, err)
	}
	return resp, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
