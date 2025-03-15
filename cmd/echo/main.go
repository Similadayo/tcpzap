package main

import (
	"context"
	"log"

	"github.com/similadayo/tcpzap/internal/metrics"
	"github.com/similadayo/tcpzap/pkg/tcpzap"
)

func main() {
	cfg := tcpzap.DefaultConfig()
	cfg.MetricsFunc = func(m metrics.Metrics) {
		log.Printf("Latency: %v, Target: %s", m.Latency, m.Target)
	}

	ctx := context.Background()
	srv, _ := tcpzap.NewServer(":8080", cfg)
	go srv.Serve(ctx, echoHandler{})

	cli, _ := tcpzap.NewClient("localhost:8080", cfg)
	resp, _ := cli.Send(ctx, []byte("Hello"))
	log.Printf("Response: %q", resp)
}

type echoHandler struct{}

func (h echoHandler) Handle(_ context.Context, msg []byte) ([]byte, error) {
	return append([]byte("Echo: "), msg...), nil
}
