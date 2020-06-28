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
	r     io.Reader
	close func() error
}

type ReaderOption func(*Reader)

func NewReader(r io.Reader, options ...ReaderOption) *Reader {
	rr := &Reader{
		r: r,
	}

	for _, opt := range options {
		opt(rr)
	}

	return rr
}

// ReadMsg reads one protobuf message from the stream
func (r *Reader) ReadMsg(m proto.Message) error {
	msgLen := proto.Size(m)

	return nil
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
