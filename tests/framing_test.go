package tests

import (
	"bytes"
	"testing"

	"github.com/similadayo/tcpzap/internal/framing"
)

func TestLengthPrefixCodec(t *testing.T) {
	codec := framing.NewCodec()
	buf := new(bytes.Buffer)

	tests := []struct {
		input   []byte
		wantErr bool
		wantMsg []byte
	}{
		{[]byte("Hello"), false, []byte("Hello")},
		{[]byte(""), false, []byte("")},
	}

	for _, tt := range tests {
		if err := codec.Encode(buf, tt.input); (err != nil) != tt.wantErr {
			t.Errorf("Encode(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
		if tt.wantErr {
			continue
		}
		got, err := codec.Decode(buf)
		if err != nil {
			t.Errorf("Decode() error = %v", err)
			continue
		}
		if !bytes.Equal(got, tt.wantMsg) {
			t.Errorf("Decode() got = %v, want %v", got, tt.wantMsg)
		}
		buf.Reset()
	}
}
