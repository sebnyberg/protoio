package protoio_test

import (
	"bytes"
	"fmt"
	"github.com/sebnyberg/protoio"
	"github.com/sebnyberg/protoio/test"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

var m = test.Person{
	Name:     "Bob",
	Phone:    "07072738293",
	Siblings: 4,
	Spouse:   true,
	Money:    1337,
}

func BenchmarkReader_ReadMsg_FromMemoryBuffer_1000(b *testing.B) {
	numMsg := 1000
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		// Write messages to the in-memory buffer
		var buf bytes.Buffer
		writeMessages(numMsg, &buf)
		protoR := protoio.NewReader(&buf)
		messages := make([]test.Person, numMsg)

		b.StartTimer()

		// Read from buffer
		for j := 0; j < numMsg; j++ {
			if err := protoR.ReadMsg(&messages[j]); err != nil {
				b.Fatal(err)
			}
		}

		b.StopTimer()
		if err := protoR.Close(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReader_ReadMsg_FromFile(b *testing.B) {
	numMsg := 1000
	for _, benchCase := range []struct {
		name string
		opts []protoio.ReaderOption
	}{
		{"Direct", nil},
		{"WithBufIO", []protoio.ReaderOption{protoio.ReadWithBufIO(1024 * 1024)}},
	} {
		context := fmt.Sprintf("%v_%v", benchCase.name, numMsg)

		// Write messages to a file
		writeFile, err := ioutil.TempFile("", "proto-test")
		if err != nil {
			panic(err)
		}
		writeMessages(numMsg, writeFile)

		b.Run(context, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()

				readFile, err := os.OpenFile(writeFile.Name(), os.O_RDONLY, 0644)
				if err != nil {
					panic(err)
				}
				protoR := protoio.NewReader(readFile, benchCase.opts...)
				messages := make([]test.Person, numMsg)

				b.StartTimer()

				// Read from buffer
				for j := 0; j < numMsg; j++ {
					if err := protoR.ReadMsg(&messages[j]); err != nil {
						b.Fatal(err)
					}
				}
				b.StopTimer()
			}
		})
	}
}

func writeMessages(numMsg int, w io.Writer) {
	protoW := protoio.NewWriter(w, protoio.WriteWithBufIO(1024*1024))
	var err error
	for i := 0; i < numMsg; i++ {
		err = protoW.WriteMsg(&m)
		if err != nil {
			panic(err)
		}
	}
	if err := protoW.Close(); err != nil {
		panic(err)
	}
}
