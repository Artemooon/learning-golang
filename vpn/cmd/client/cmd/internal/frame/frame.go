package frame

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const MaxPayload = 64 * 1024
const HeaderSize = 4

var (
	ErrEmptyPayload = errors.New("frame payload is empty")
	ErrTooLarge     = errors.New("frame payload exceeds maximum size")
)

func WriteFrame(w io.Writer, payload []byte) error {
	if err := validatePayloadLength(len(payload)); err != nil {
		return err
	}

	var header [HeaderSize]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(payload)))

	if err := writeAll(w, header[:]); err != nil {
		return err
	}
	return writeAll(w, payload)
}

func ReadFrame(r io.Reader) ([]byte, error) {
	var header [HeaderSize]byte

	if _, err := io.ReadFull(r, header[:]); err != nil {
		return nil, err
	}

	length := int(binary.BigEndian.Uint32(header[:]))
	if err := validatePayloadLength(length); err != nil {
		return nil, err
	}

	payload := make([]byte, length)

	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}

	return payload, nil
}

func validatePayloadLength(length int) error {
	if length == 0 {
		return ErrEmptyPayload
	}
	if length > MaxPayload {
		return fmt.Errorf(
			"%w: got %d bytes, maximum is %d bytes",
			ErrTooLarge,
			length,
			MaxPayload,
		)
	}
	return nil
}

func writeAll(w io.Writer, data []byte) error {
	for len(data) > 0 {
		n, err := w.Write(data)
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrShortWrite
		}
		data = data[n:]
	}
	return nil
}
