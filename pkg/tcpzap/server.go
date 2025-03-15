package tcpzap

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/similadayo/tcpzap/internal/framing"
	"github.com/similadayo/tcpzap/internal/transport"
)

// Server manages TCP connections and message handling.
type Server struct {
	ln     net.Listener
	codec  framing.Codec
	h      Handler
	mu     sync.Mutex
	closed bool
	cfg    Config
}

// NewServer initializes a TCP server with the given address and config.
func NewServer(addr string, cfg Config) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("tcpzap: listen: %w", err)
	}
	return &Server{
		ln:    ln,
		codec: framing.NewCodec(),
		cfg:   cfg,
	}, nil
}

// Serve accepts connections and processes them with the handler.
func (s *Server) Serve(ctx context.Context, h Handler) error {
	s.mu.Lock()
	s.h = h
	s.mu.Unlock()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			conn, err := s.ln.Accept()
			if err != nil {
				s.mu.Lock()
				if s.closed {
					s.mu.Unlock()
					return nil
				}
				s.mu.Unlock()
				log.Printf("tcpzap: accept: %v", err)
				continue
			}
			go s.handleConn(ctx, conn)
		}
	}
}

// handleConn processes a single client connection.
func (s *Server) handleConn(ctx context.Context, conn net.Conn) {
	tCfg := transport.Config{
		Retries:    s.cfg.Retries,
		RetryDelay: s.cfg.RetryDelay,
	}
	c := transport.NewConn(conn, s.codec, tCfg)
	defer c.Close()

	for {
		msg, err := c.Receive(ctx)
		if err != nil {
			log.Printf("tcpzap: receive: %v", err)
			return
		}
		resp, err := s.h.Handle(ctx, msg)
		if err != nil {
			log.Printf("tcpzap: handle: %v", err)
			return
		}
		if err := c.Send(ctx, resp); err != nil {
			log.Printf("tcpzap: send: %v", err)
			return
		}
	}
}

// Close shuts down the server gracefully.
func (s *Server) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return s.ln.Close()
}

// Ln returns the underlying listener (for testing).
func (s *Server) Ln() net.Listener {
	return s.ln
}
