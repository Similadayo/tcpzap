package test

import (
	"context"
	"testing"
	"time"

	"github.com/similadayo/tcpzap/internal/congestion"
)

func TestController(t *testing.T) {
	cfg := congestion.DefaultConfig()
	ctrl := congestion.NewController(cfg)

	ctx := context.Background()

	for i := 0; i < cfg.InitialWindow; i++ {
		if !ctrl.CanSend(ctx) {
			t.Fatalf("CanSend failed at %d, expected true", i)
		}
	}

	if ctrl.CanSend(ctx) {
		t.Error("CanSend should be false when window is full")
	}

	for i := 0; i < 5; i++ {
		ctrl.AckReceived()
		time.Sleep(10 * time.Millisecond)
	}

	if ctrl.Window() >= cfg.InitialWindow {
		t.Errorf("Window %d should have decreased after congestion", ctrl.Window())
	}

	if !ctrl.CanSend(ctx) {
		t.Error("CanSend should be true after acks")
	}
}
