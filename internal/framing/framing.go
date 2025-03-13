package framing

import (
	"encoding/binary"
	"io"
)

// Write sends a message by prefixing it with its length
func Write(conn io.Writer, message []byte) error {
	// Write the lenth as a 4-byte unit32
	if err := binary.Write(conn, binary.BigEndian, uint32(len(message))); err != nil {
		return err
	}

	// Write the message
	_, err := conn.Write(message)
	return err
}

// Read reads a message by first reading its length
func Read(conn io.Reader) ([]byte, error) {
	// Read the message length
	var length uint32
	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	// Read the message
	message := make([]byte, length)
	_, err := io.ReadFull(conn, message)
	return message, err
}
