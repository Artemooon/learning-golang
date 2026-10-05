package frame

import (
    "bytes"
    "io"
    "net"
    "errors"
    "testing"
)

func TestWriteEmptyPayloadRejects(t *testing.T) {
    err := WriteFrame(io.Discard, []byte{})

    if !errors.Is(err, ErrEmptyPayload) {
        t.Fatalf("WriteFrame(empty) error = %v, want ErrEmptyPayload", err)
    }
}


func TestWriteEmptyPayloadLargerThen64KibRejects(t *testing.T) {
    payload := bytes.Repeat([]byte("s"), 70*1024)
    err := WriteFrame(io.Discard, payload)

    if !errors.Is(err, ErrTooLarge) {
        t.Fatalf("WriteFrame(empty) error = %v, want ErrTooLarge", err)
    }
}

func TestTwoFramesRemainSeparate(t *testing.T) {
    client, server := net.Pipe()
    defer client.Close()
    defer server.Close()

    first := []byte("first message")
    second := []byte("second message")

    writeErr := make(chan error, 1)

    go func() {
            if err := WriteFrame(client, first); err != nil {
                    writeErr <- err
                    return
            }

            writeErr <- WriteFrame(client, second)
    }()

    gotFirst, err := ReadFrame(server)
    if err != nil {
            t.Fatalf("read first frame: %v", err)
    }
    if !bytes.Equal(gotFirst, first) {
            t.Fatalf("first frame = %q, want %q", gotFirst, first)
    }

    gotSecond, err := ReadFrame(server)
    if err != nil {
            t.Fatalf("read second frame: %v", err)
    }
    if !bytes.Equal(gotSecond, second) {
            t.Fatalf("second frame = %q, want %q", gotSecond, second)
    }

    if err := <-writeErr; err != nil {
            t.Fatalf("write frames: %v", err)
    }
}