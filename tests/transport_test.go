package test

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"

	"github.com/similadayo/tcpzap/internal/framing"
	"github.com/similadayo/tcpzap/internal/transport"
)

func TestTransportSendReceive(t *testing.T) {
	// Setup: Create a pipe to simulate a connection
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	cfg := transport.Config{
		Retries:    1,
		RetryDelay: 10 * time.Millisecond,
	}
	client := transport.NewConn(clientConn, framing.NewCodec(), cfg)
	server := transport.NewConn(serverConn, framing.NewCodec(), cfg)

	ctx := context.Background()

	// Test sending and receiving
	go func() {
		data, err := server.Receive(ctx)
		if err != nil {
			t.Errorf("Server receive: %v", err)
			return
		}
		if err := server.Send(ctx, append([]byte("Echo: "), data...)); err != nil {
			t.Errorf("Server send: %v", err)
		}
	}()

	msg := []byte("Hello")
	if err := client.Send(ctx, msg); err != nil {
		t.Fatalf("Client send: %v", err)
	}
	resp, err := client.Receive(ctx)
	if err != nil {
		t.Fatalf("Client receive: %v", err)
	}
	if !bytes.Equal(resp, []byte("Echo: Hello")) {
		t.Errorf("Expected %q, got %q", "Echo: Hello", resp)
	}
}
