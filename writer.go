package protoio

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
)

// Writer writes protobuf messages to an io.Writer
type Writer struct {
	w         io.Writer
	lenBuf    []byte
	msgBuf    []byte
	close     func() error
	byteOrder binary.ByteOrder
}

type WriterOption func(*Writer)

// WriteWithBufIO wraps the writer in a BufIO writer, improving performance
// when writing to / from files
func WriteWithBufIO(size int) WriterOption {
	return func(w *Writer) {
		writer := bufio.NewWriterSize(w.w, size)
		w.close = func() error {
			// Flush first, then close the embedded Writer
			flushErr := writer.Flush()
			if closer, ok := w.w.(io.Closer); ok {
				// Ensure that embedded writer is closed even if flush fails
				if err := closer.Close(); err != nil {
					return err
				}
			}
			return flushErr
		}
		w.w = writer
	}
}

func WriteWithByteOrder(bo binary.ByteOrder) WriterOption {
	return func(r *Writer) {
		r.byteOrder = bo
	}
}

func NewWriter(w io.Writer, options ...WriterOption) *Writer {
	wr := &Writer{
		w:      w,
		lenBuf: make([]byte, 8),
	}

	for _, opt := range options {
		opt(wr)
	}

	return wr
}

// WriteMsg writes one protobuf message to the embedded writer
func (w *Writer) WriteMsg(m proto.Message) error {
	msgLen := proto.Size(m)

	// Write length of message
	w.byteOrder.PutUint32(w.lenBuf, uint32(msgLen))
	_, err := w.w.Write(w.lenBuf)
	if err != nil {
		return fmt.Errorf("failed to write msg length, err: %w", err)
	}

	// Expand buffer if it is too small
	if cap(w.msgBuf) < msgLen {
		w.msgBuf = make([]byte, msgLen)
	}

	// Read message into buffer
	msgBytes, err := proto.MarshalOptions{}.MarshalAppend(w.msgBuf[:0], m)
	if err != nil {
		return fmt.Errorf("failed to get msg bytes, err: %w", err)
	}

	// Write message
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
	if w.close != nil {
		return w.close()
	}
	return nil
}
