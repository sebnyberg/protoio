package protoio_test

import (
	"github.com/sebnyberg/protoio"
	"github.com/sebnyberg/protoio/test"
	"io/ioutil"
	"testing"
)

func BenchmarkWriter_WriteMsg(b *testing.B) {
	messages := make([]*test.SimpleMessage, 10)
	for i := 0; i < 10; i++ {
		messages[i] = &test.SimpleMessage{Message: int32(i)}
	}
	m := test.SimpleMessageList{Messages: messages}

	for _, benchCase := range []struct {
		name    string
		options []protoio.Option
	}{
		{"WithoutLenBuf", []protoio.Option{protoio.WithoutLenBuf, protoio.WithoutMsgBuf}},
		{"WithLenBuf", []protoio.Option{protoio.WithLenBuf, protoio.WithoutMsgBuf}},
		{"WithMsgBuf", []protoio.Option{protoio.WithLenBuf, protoio.WithMsgBuf}},
	} {
		b.Run(benchCase.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()

				protoWriter := protoio.NewWriter(ioutil.Discard, benchCase.options...)

				b.StartTimer()
				// Write messages
				for j := 0; j < int(1e4); j++ {
					err := protoWriter.WriteMsg(&m)
					if err != nil {
						b.Error(err)
					}
				}
			}
		})
	}
}

//type Writer struct {
//	w io.Writer
//}
//
//func NewWriter(w io.Writer) *Writer {
//	return &Writer{w: w}
//}
//
//// WriteMsg writes one length-delimited protobuf message to an io.Writer
//func (w *Writer) WriteMsg(m proto.Message) error {
//	// Write message length
//	lenBuf := make([]byte, 8)
//	binary.BigEndian.PutUint64(lenBuf, uint64(proto.Size(m)))
//	_, err := w.w.Write(lenBuf)
//	if err != nil {
//		return err
//	}
//
//	// Write message
//	data, err := proto.Marshal(m)
//	if err != nil {
//		return err
//	}
//	_, err = w.w.Write(data)
//	return err
//}

//func BenchmarkWriter_WriteMsgFile(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		b.StopTimer()
//
//		f, err := ioutil.TempFile("", "proto-test")
//		if err != nil {
//			b.Error("failed to create test file, err: ", err)
//		}
//
//		protoWriter := protoio.NewWriter(f)
//
//		b.StartTimer()
//		// Write messages
//		for j := 0; j < int(1e4); j++ {
//			err := protoWriter.WriteMsg(&m)
//			if err != nil {
//				b.Error(err)
//			}
//		}
//		err = protoWriter.Close()
//		if err != nil {
//			b.Error("failed to close proto writer, err: ", err)
//		}
//		b.StopTimer()
//
//		err = os.Remove(f.Name())
//		if err != nil {
//			b.Error("failed to remove test file")
//		}
//	}
//}
