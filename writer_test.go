package protoio_test

import (
	"fmt"
	qt "github.com/frankban/quicktest"
	"github.com/sebnyberg/protoio"
	"github.com/sebnyberg/protoio/test"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestWriter_WriteMsg(t *testing.T) {
	c := qt.New(t)

	testPerson := test.Person{
		Name:     "Bob",
		Phone:    "07072738293",
		Siblings: 4,
		Spouse:   true,
		Money:    1337,
	}

	f, err := os.OpenFile("out", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	c.Assert(err, qt.IsNil)
	protoW := protoio.NewWriter(f)

	numMsg := 1000
	for i := 0; i < numMsg; i++ {
		c.Assert(protoW.WriteMsg(&testPerson), qt.IsNil)
	}

	c.Assert(protoW.Close(), qt.IsNil)
}

func BenchmarkWriter_WriteMsg(b *testing.B) {
	m := test.Person{
		Name:     "Bob",
		Phone:    "07072738293",
		Siblings: 4,
		Spouse:   true,
		Money:    1337,
	}
	numMsg := 1000

	for _, benchCase := range []struct {
		name         string
		options      []protoio.WriterOption
		getOutWriter func() (io.Writer, func())
	}{
		{
			"InMem",
			nil,
			func() (io.Writer, func()) {
				return ioutil.Discard, func() {}
			},
		},
		{
			"ToFile",
			nil,
			getTempFile,
		},
		{
			"ToFileBufIO",
			[]protoio.WriterOption{protoio.WriterWithBufIO(1024 * 1024)},
			getTempFile,
		},
	} {
		context := fmt.Sprintf("%v_%v", numMsg, benchCase.name)
		b.Run(context, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()

				w, cleanup := benchCase.getOutWriter()
				protoWriter := protoio.NewWriter(w, benchCase.options...)

				b.StartTimer()

				for j := 0; j < numMsg; j++ {
					err := protoWriter.WriteMsg(&m)
					if err != nil {
						b.Error(err)
					}
				}

				b.StopTimer()
				cleanup()
			}
		})
	}
}

func getTempFile() (io.Writer, func()) {
	f, err := ioutil.TempFile("", "proto-test")
	if err != nil {
		panic(err)
	}
	return f, func() {
		err := os.Remove(f.Name())
		if err != nil {
			panic(err)
		}
	}
}
