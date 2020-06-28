package protoio_test

import (
	"github.com/sebnyberg/protoio"
	"github.com/sebnyberg/protoio/test"
	"io/ioutil"
	"testing"
)

var m = test.SimpleMessage{
	Message: "hello",
}

func BenchmarkWriter_WriteMsg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		protoWriter := protoio.NewWriter(ioutil.Discard)

		b.StartTimer()
		// Write messages
		for j := 0; j < int(1e4); j++ {
			err := protoWriter.WriteMsg(&m)
			if err != nil {
				b.Error(err)
			}
		}
	}
}

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
