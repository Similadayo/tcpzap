// cmd/echo/main.go
package main

import (
	"context"
	"log"
	"time"

	"github.com/similadayo/tcpzap/internal/metrics"
	"github.com/similadayo/tcpzap/pkg/tcpzap"
)

type echoHandler struct{}

func (h echoHandler) Handle(_ context.Context, msg []byte) ([]byte, error) {
	return append([]byte("Echo: "), msg...), nil
}

func main() {
	cfg := tcpzap.DefaultConfig()
	cfg.MetricsFunc = func(m metrics.Metrics) {
		log.Printf("Latency: %v, Target: %s, Success: %v", m.Latency, m.Target, m.Success)
	}

	ctx := context.Background()
	srv, err := tcpzap.NewServer(":8080", cfg)
	if err != nil {
		log.Fatalf("Server init: %v", err)
	}
	go func() {
		if err := srv.Serve(ctx, echoHandler{}); err != nil {
			log.Fatalf("Server: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	cli, err := tcpzap.NewClient("localhost:8080", cfg)
	if err != nil {
		log.Fatalf("Client init: %v", err)
	}
	defer cli.Close()

	resp, err := cli.Send(ctx, []byte("Hello, tcpzap!"))
	if err != nil {
		log.Fatalf("Client send: %v", err)
	}
	log.Printf("Response: %q", resp)
}
