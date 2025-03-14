package framing

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Codec interface {
	Encode(w io.Writer, msg []byte) error
	Decode(r io.Reader) ([]byte, error)
}

// LengthPrefixCodec is a framing codec that prefixes each message with its length
type LengthPrefixCodec struct{}

// encode writes a message with 32-bit legnth prefix
func (c LengthPrefixCodec) Encode(w io.Writer, msg []byte) error {
	if len(msg) > 1<<32-1 {
		return fmt.Errorf("framing: message too long: %d", len(msg))
	}
	if err := binary.Write(w, binary.BigEndian, uint32(len(msg))); err != nil {
		return fmt.Errorf("framing: failed to write message length: %v", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("framing: failed to write message: %v", err)
	}
	return nil
}

// decode reads a message with 32-bit length prefix
func (c LengthPrefixCodec) Decode(r io.Reader) ([]byte, error) {
	var length uint32
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return nil, fmt.Errorf("framing: failed to read message length: %v", err)
	}
	msg := make([]byte, length)
	if _, err := io.ReadFull(r, msg); err != nil {
		return nil, fmt.Errorf("framing: failed to read message: %v", err)
	}
	return msg, nil
}

func NewCodec() Codec {
	return LengthPrefixCodec{}
}
