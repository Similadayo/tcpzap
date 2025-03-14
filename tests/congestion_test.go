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

	//Fill the window
	for i := 0; i < cfg.InitialWindow; i++ {
		if !ctrl.CanSend(ctx) {
			t.Fatalf("CanSend failed at %d, expected true", i)
		}
	}

	//should block now
	if ctrl.CanSend(ctx) {
		t.Fatalf("CanSend should be false when window is full")
	}

	//simulate acks
	for i := 0; i < 5; i++ {
		ctrl.AckReceived()
		time.Sleep(10 * time.Millisecond) //simulate RTT
	}

	//Window should adjust
	if ctrl.Window() >= cfg.InitialWindow {
		t.Fatalf("Window should be less than initial window")
	}

	//should allow sending again
	if !ctrl.CanSend(ctx) {
		t.Fatalf("CanSend failed after acks, expected true")
	}

}
