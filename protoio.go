package protoio

import (
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
)

type Writer struct {
	w      io.Writer
	lenBuf []byte
}

func WithLenBuf(w *Writer) {
	w.lenBuf = make([]byte, 0, 8)
}

func WithoutLenBuf(w *Writer) {
	w.lenBuf = nil
}

func NewWriter(w io.Writer, options ...func(*Writer)) *Writer {
	wr := &Writer{
		w: w,
	}

	// Default options
	WithLenBuf(wr)

	for _, opt := range options {
		opt(wr)
	}

	return wr
}

// WriteMsg writes one protobuf message to the embedded writer
func (w *Writer) WriteMsg(m proto.Message) error {
	// Write length of message
	lenBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(lenBuf, uint64(proto.Size(m)))
	_, err := w.w.Write(lenBuf)
	if err != nil {
		return fmt.Errorf("failed to write msg length, err: %w", err)
	}

	// Write message
	data, err := proto.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal msg, err: %w", err)
	}
	_, err = w.w.Write(data)
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
