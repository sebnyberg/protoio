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
	w           io.Writer
	lenBuf      []byte
	msgBuf      []byte
	writeMsgLen func(int) error
	close       func() error
}

type WriterOption func(*Writer)

// WriterWithBufIO wraps the writer in a BufIO writer, improving performance
// when writing to / from files
func WriterWithBufIO(size int) WriterOption {
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

// WriterDelimiter sets the byte order and type for the message delimiter
func WriterDelimiter(bo binary.ByteOrder, delimType DelimiterType) WriterOption {
	return func(w *Writer) {
		switch delimType {
		case DelimiterTypeUint32:
			w.lenBuf = make([]byte, LenBytesUint32)
			w.writeMsgLen = func(msgLen int) error {
				bo.PutUint32(w.lenBuf, uint32(msgLen))
				_, err := w.w.Write(w.lenBuf)
				return err
			}
		case DelimiterTypeUint64:
			w.lenBuf = make([]byte, LenBytesUint64)
			w.writeMsgLen = func(msgLen int) error {
				bo.PutUint64(w.lenBuf, uint64(msgLen))
				_, err := w.w.Write(w.lenBuf)
				return err
			}
		default:
			panic("unknown lentype")
		}
	}
}

// NewWriter returns a streaming Protobuf writer.
// By default, the byte order is set to BigEndian and the length
// delimiter to uint32. This behaviour can be changed by providing
// the WriterDelimiter() option to this constructor
func NewWriter(w io.Writer, options ...WriterOption) *Writer {
	wr := &Writer{
		w:      w,
		lenBuf: make([]byte, 8),
	}

	// Default options
	WriterDelimiter(binary.BigEndian, DelimiterTypeUint32)(wr)

	for _, opt := range options {
		opt(wr)
	}

	return wr
}

// WriteMsg writes one protobuf message to the embedded writer
func (w *Writer) WriteMsg(m proto.Message) error {
	msgLen := proto.Size(m)

	// Write length of message
	if err := w.writeMsgLen(msgLen); err != nil {
		return fmt.Errorf("failed to write message length: %w", err)
	}

	// Expand buffer if it is too small
	if cap(w.msgBuf) < msgLen {
		w.msgBuf = make([]byte, msgLen)
	}

	// Read message into buffer
	msgBytes, err := proto.MarshalOptions{}.MarshalAppend(w.msgBuf[:0], m)
	if err != nil {
		return fmt.Errorf("failed to marshal msg bytes: %w", err)
	}

	// Write message
	_, err = w.w.Write(msgBytes)
	if err != nil {
		return fmt.Errorf("failed to write msg: %w", err)
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
