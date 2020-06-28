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
	getMsgLen  func() (int, error)
}

// ReaderOption changes the behaviour of the reader.
type ReaderOption func(*Reader)

// ReadWithMsgSizeMax sets the maximum protobuf message size.
func ReadWithMsgSizeMax(size int) ReaderOption {
	return func(r *Reader) {
		r.msgSizeMax = size
	}
}

// ReadWithBufIO wraps the reader in a BufIO reader. This is an effective
// way of improving speed when reading to and from files.
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

// ReadWithDelimiter sets the byte order and type for the message delimiter
func ReadWithDelimiter(bo binary.ByteOrder, lenType DelimiterType) ReaderOption {
	return func(r *Reader) {
		switch lenType {
		case DelimiterTypeUint32:
			r.lenBuf = make([]byte, LenBytesUint32)
			r.getMsgLen = func() (int, error) {
				if _, err := io.ReadFull(r.r, r.lenBuf); err != nil {
					return 0, err
				}
				return int(bo.Uint32(r.lenBuf)), nil
			}
		case DelimiterTypeUint64:
			r.lenBuf = make([]byte, LenBytesUint64)
			r.getMsgLen = func() (int, error) {
				if _, err := io.ReadFull(r.r, r.lenBuf); err != nil {
					return 0, err
				}
				return int(bo.Uint64(r.lenBuf)), nil
			}
		default:
			panic("unknown len type")
		}
	}
}

// NewReader returns a streaming Protobuf message reader
// By default, the reader uses binary.BigEndian and Uint32,
// To change this behaviour, pass the ReadWithDelimiter option to this constructor
// If a  read message exceeds the MsgSizeMax, an error will be returned by the reader
func NewReader(r io.Reader, options ...ReaderOption) *Reader {
	rr := &Reader{
		r: r,
	}

	// Recommended max Protobuf message size (4MB)
	ReadWithMsgSizeMax(1024 * 1024 * 4)(rr)
	ReadWithDelimiter(binary.BigEndian, DelimiterTypeUint32)(rr)

	for _, opt := range options {
		opt(rr)
	}

	return rr
}

// ReadMsg reads one protobuf message from the stream.
// If the message is greater than the maximum message size,
// io.ErrShortBuffer is returned.
func (r *Reader) ReadMsg(m proto.Message) error {
	msgLen, err := r.getMsgLen()
	if err != nil {
		return err
	}

	if msgLen < 0 || msgLen > r.msgSizeMax {
		return fmt.Errorf("proto msg > msgSizeMax: %w", io.ErrShortBuffer)
	}

	if msgLen > len(r.msgBuf) {
		r.msgBuf = make([]byte, msgLen)
	}

	_, err = io.ReadFull(r.r, r.msgBuf[:msgLen])
	if err != nil {
		return fmt.Errorf("failed to read message: %w", err)
	}

	return proto.Unmarshal(r.msgBuf[:msgLen], m)
}

// Close closes any embedded io.Closer readers and flushes the bufio reader
// if that option was provided.
func (r *Reader) Close() error {
	if closer, ok := r.r.(io.Closer); ok {
		return closer.Close()
	}
	if r.close != nil {
		return r.close()
	}
	return nil
}
