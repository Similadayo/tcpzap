package tests

import (
	"bytes"
	"testing"

	"github.com/similadayo/tcpzap/internal/framing"
)

func TestFraming(t *testing.T) {
	//simulate connection with buffer
	buf := new(bytes.Buffer)

	//test data
	msg := []byte("hello, tcpzap")
	err := framing.Write(buf, msg)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// read it back
	data, err := framing.Read(buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if !bytes.Equal(data, msg) {
		t.Fatalf("expected %q, got %q", msg, data)
	}
}
