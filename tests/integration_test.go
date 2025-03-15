package test

import (
	"context"
	"testing"
	"time"

	"github.com/similadayo/tcpzap/internal/metrics"
	"github.com/similadayo/tcpzap/pkg/tcpzap"
)

func TestClientServerIntegration(t *testing.T) {
	cfg := tcpzap.DefaultConfig()
	cfg.MetricsFunc = func(m metrics.Metrics) {
		if m.Latency == 0 || !m.Success {
			t.Errorf("Metrics: invalid latency %v or success %v", m.Latency, m.Success)
		}
	}

	ctx := context.Background()
	srv, err := tcpzap.NewServer(":0", cfg) // Use :0 for random port
	if err != nil {
		t.Fatalf("Server init: %v", err)
	}
	defer srv.Close()

	go func() {
		if err := srv.Serve(ctx, testHandler{}); err != nil {
			t.Errorf("Server serve: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond) // Wait for server

	addr := srv.Ln().Addr().String()
	cli, err := tcpzap.NewClient(addr, cfg)
	if err != nil {
		t.Fatalf("Client init: %v", err)
	}
	defer cli.Close()

	resp, err := cli.Send(ctx, []byte("Test"))
	if err != nil {
		t.Fatalf("Client send: %v", err)
	}
	if string(resp) != "Echo: Test" {
		t.Errorf("Expected %q, got %q", "Echo: Test", resp)
	}
}

type testHandler struct{}

func (h testHandler) Handle(_ context.Context, msg []byte) ([]byte, error) {
	return append([]byte("Echo: "), msg...), nil
}
