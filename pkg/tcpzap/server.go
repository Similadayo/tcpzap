package tcpzap

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/similadayo/tcpzap/internal/framing"
	"github.com/similadayo/tcpzap/internal/transport"
)

type Server struct {
	ln     net.Listener
	codec  framing.Codec
	h      Handler
	mu     sync.Mutex
	closed bool
}

// NewServer creates a new server
func NewServer(addr string, timeout time.Duration) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("tcpzap: listen: %w", err)
	}
	return &Server{ln: ln, codec: framing.NewCodec()}, nil
}

// Serve accepts incoming connections on the listener
func (s *Server) Serve(ctx context.Context, h Handler) error {
	s.mu.Lock()
	s.h = h
	s.mu.Unlock()

	for {
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

// handleConn handles a connection
func (s *Server) handleConn(ctx context.Context, conn net.Conn) {
	c := transport.NewConn(conn, s.codec)
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

// Close closes the server
func (s *Server) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return s.ln.Close()
}
