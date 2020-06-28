package protoio

import (
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
)

// Writer writes protobuf messages to an io.Writer
type Writer struct {
	w      io.Writer
	lenBuf []byte
	getLen func(m proto.Message, msgLen int) []byte
	msgBuf []byte
	getMsg func(m proto.Message, msgLen int) ([]byte, error)
}

type Option func(*Writer)

// WithLenBuf re-uses an internal buffer for writing the message length.
func WithLenBuf(w *Writer) {
	w.lenBuf = make([]byte, 8)
	w.getLen = func(m proto.Message, msgLen int) []byte {
		binary.BigEndian.PutUint64(w.lenBuf, uint64(msgLen))
		return w.lenBuf
	}
}

// WithoutLenBuf creates a buffer every time for writing the message length.
func WithoutLenBuf(w *Writer) {
	w.getLen = func(m proto.Message, msgLen int) []byte {
		lenBuf := make([]byte, 8)
		binary.BigEndian.PutUint64(lenBuf, uint64(msgLen))
		return lenBuf
	}
}

// WithMsgBuf re-uses an internal buffer for writing messages.
func WithMsgBuf(w *Writer) {
	w.getMsg = func(m proto.Message, msgLen int) ([]byte, error) {
		// Expand buffer if it is too small
		if cap(w.msgBuf) < msgLen {
			w.msgBuf = make([]byte, msgLen)
		}

		// Read message into buffer
		return proto.MarshalOptions{}.MarshalAppend(w.msgBuf[:0], m)
	}
}

// WithoutMsgBuf creates a buffer every time a message is written.
func WithoutMsgBuf(w *Writer) {
	w.getMsg = func(m proto.Message, msgLen int) ([]byte, error) {
		return proto.Marshal(m)
	}
}

func NewWriter(w io.Writer, options ...Option) *Writer {
	wr := &Writer{
		w: w,
	}

	// Default options
	WithLenBuf(wr)
	WithMsgBuf(wr)

	// Overrides
	for _, opt := range options {
		opt(wr)
	}

	return wr
}

// WriteMsg writes one protobuf message to the embedded writer
func (w *Writer) WriteMsg(m proto.Message) error {
	msgLen := proto.Size(m)

	// Write length of message
	lenBytes := w.getLen(m, msgLen)
	_, err := w.w.Write(lenBytes)
	if err != nil {
		return fmt.Errorf("failed to write msg length, err: %w", err)
	}

	// Write message
	msgBytes, err := w.getMsg(m, msgLen)
	if err != nil {
		return fmt.Errorf("failed to get msg bytes, err: %w", err)
	}
	_, err = w.w.Write(msgBytes)
	if err != nil {
		return fmt.Errorf("failed to write msg, err: %w", err)
	}

	return nil
}

func (w *Writer) Close() error {
	if closer, ok := w.w.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
