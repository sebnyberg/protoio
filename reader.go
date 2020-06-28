package protoio

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
)

// Reader reads Protobuf messages from a stream
type Reader struct {
	r          io.Reader
	lenBuf     []byte
	msgBuf     []byte
	close      func() error
	msgSizeMax int
	byteOrder  binary.ByteOrder
}

type ReaderOption func(*Reader)

func ReadWithMsgSizeMax(size int) ReaderOption {
	return func(r *Reader) {
		r.msgSizeMax = size
	}
}

func ReadWithBufIO(size int) ReaderOption {
	return func(r *Reader) {
		reader := bufio.NewReaderSize(r.r, size)
		r.close = func() error {
			if closer, ok := r.r.(io.Closer); ok {
				// Ensure that embedded writer is closed even if flush fails
				if err := closer.Close(); err != nil {
					return err
				}
			}
			return nil
		}
		r.r = reader
	}
}

func ReadWithByteOrder(bo binary.ByteOrder) ReaderOption {
	return func(r *Reader) {
		r.byteOrder = bo
	}
}

func NewReader(r io.Reader, options ...ReaderOption) *Reader {
	rr := &Reader{
		r:      r,
		lenBuf: make([]byte, 8),
	}

	// Recommended max Protobuf message size (4MB)
	ReadWithMsgSizeMax(1024 * 1024 * 4)(rr)

	for _, opt := range options {
		opt(rr)
	}

	return rr
}

// ReadMsg reads one protobuf message from the stream.
// If the message is greater than the maximum message size,
// io.ErrShortBuffer is returned.
func (r *Reader) ReadMsg(m proto.Message) error {
	if _, err := io.ReadFull(r.r, r.lenBuf); err != nil {
		return err
	}

	msgLen := int(r.byteOrder.Uint32(r.lenBuf))
	if msgLen < 0 || msgLen > r.msgSizeMax {
		return fmt.Errorf("proto msg > msgSizeMax: %w", io.ErrShortBuffer)
	}

	if msgLen > len(r.msgBuf) {
		r.msgBuf = make([]byte, msgLen)
	}

	_, err := io.ReadFull(r.r, r.msgBuf[:msgLen])
	if err != nil {
		return fmt.Errorf("failed to read message: %w", err)
	}

	return proto.Unmarshal(r.msgBuf[:msgLen], m)
}

func (r *Reader) Close() error {
	if closer, ok := r.r.(io.Closer); ok {
		return closer.Close()
	}
	if r.close != nil {
		return r.close()
	}
	return nil
}
