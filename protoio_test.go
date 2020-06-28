package protoio_test

//var m = test.AllFields{
//	Field1:  10.2,
//	Field2:  20.32,
//	Field3:  1,
//	Field4:  2,
//	Field5:  3,
//	Field6:  4,
//	Field7:  5,
//	Field8:  6,
//	Field9:  7,
//	Field10: 8,
//	Field11: 9,
//	Field12: 10,
//	Field13: false,
//	Field14: "test",
//	Field15: []byte("test"),
//}
//
//func BenchmarkSimpleWriter(b *testing.B) {
//	for _, want := range []struct {
//		outCount           int
//		filename, testname string
//		newWriter          func(io.Writer) protoio.Writer
//	}{
//		{
//			outCount: 1e5,
//			filename: "out.proto",
//			testname: "Naive",
//			newWriter: func(w io.Writer) protoio.Writer {
//				return writer{w: w}
//			},
//		},
//		{
//			outCount: 1e5,
//			filename: "out.proto",
//			testname: "WithLenBuf",
//			newWriter: func(w io.Writer) protoio.Writer {
//				return writer{
//					w:      w,
//					lenBuf: make([]byte, 8),
//				}
//			},
//		},
//		{
//			outCount: 1e5,
//			filename: "out.proto",
//			testname: "WithMsgBuf",
//			newWriter: func(w io.Writer) protoio.Writer {
//				return writer{
//					w:      w,
//					lenBuf: make([]byte, 8),
//					buf:    &bytes.Buffer{},
//				}
//			},
//		},
//		{
//			outCount: 1e5,
//			filename: "out.proto",
//			testname: "WithMsgBuf_And_Maxlen",
//			newWriter: func(w io.Writer) protoio.Writer {
//				return writer{
//					w:      w,
//					lenBuf: make([]byte, 8),
//					buf:    &bytes.Buffer{},
//					maxLen: 1024,
//				}
//			},
//		},
//		{
//			outCount: 1e5,
//			filename: "out.proto",
//			testname: "WithMsgBufLarge_And_Maxlen",
//			newWriter: func(w io.Writer) protoio.Writer {
//				return writer{
//					w:      w,
//					lenBuf: make([]byte, 8),
//					buf:    bytes.NewBuffer(make([]byte, 0, 3*1024*1024)),
//					maxLen: 1024,
//				}
//			},
//		},
//		{
//			outCount: 1e6,
//			filename: "out.proto",
//			testname: "WithBuf_Maxlen_LenToBuf",
//			newWriter: func(w io.Writer) protoio.Writer {
//				return writer{
//					w:             w,
//					lenBuf:        make([]byte, 8),
//					buf:           bytes.NewBuffer(make([]byte, 0, 3*1024*1024)),
//					maxLen:        1024,
//					writeLenToBuf: true,
//				}
//			},
//		},
//	} {
//		context := fmt.Sprintf("Write_%v_%v", want.testname, want.outCount)
//		b.Run(context, func(b *testing.B) {
//			for i := 0; i < b.N; i++ {
//				b.StopTimer()
//
//				err := os.Remove(want.filename)
//				if err != nil {
//					b.Error("failed, err: ", err)
//				}
//				f, err := os.OpenFile(want.filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
//				if err != nil {
//					b.Error("Failed, err: ", err)
//				}
//				protoWriter := want.newWriter(f)
//				if err != nil {
//					b.Error("Failed, err: ", err)
//				}
//
//				b.StartTimer()
//				for j := 0; j < want.outCount; j++ {
//					err := protoWriter.WriteMsg(&m)
//					if err != nil {
//						b.Error("failed, err: ", err)
//					}
//				}
//				err = protoWriter.Close()
//				if err != nil {
//					b.Error("failed, err: ", err)
//				}
//			}
//		})
//	}
//}
//
//// writer does not use any buffer, it simply writes to the output writer immediately
//type writer struct {
//	w             io.Writer
//	lenBuf        []byte
//	fixedBuf      []byte
//	buf           *bytes.Buffer
//	maxLen        int
//	writeLenToBuf bool
//}
//
//func (w writer) WriteMsg(m proto.Message) error {
//	// Write length of message (delimiter)
//	var (
//		n   int
//		err error
//	)
//
//	if w.lenBuf == nil {
//		lenBuf := make([]byte, 8)
//		binary.BigEndian.PutUint64(lenBuf, uint64(proto.Size(m)))
//		n, err = w.w.Write(lenBuf[:])
//	} else {
//		binary.BigEndian.PutUint64(w.lenBuf, uint64(proto.Size(m)))
//		if w.writeLenToBuf {
//			n, err = w.buf.Write(w.lenBuf)
//		} else {
//			n, err = w.w.Write(w.lenBuf[:])
//		}
//	}
//
//	if err != nil {
//		return fmt.Errorf("failed to write message length, err: %v", err)
//	}
//	if n == 0 {
//		return errors.New("unable to write delimiter")
//	}
//
//	// Write message
//	if w.buf != nil {
//		written, err := proto.MarshalOptions{}.MarshalAppend(w.buf.Bytes(), m)
//		if err != nil {
//			return fmt.Errorf("failed to append protobuf message to buffer bytes, err: %v", err)
//		}
//		n, err := w.buf.Write(written)
//		if err != nil {
//			return fmt.Errorf("failed to append protobuf message to buffer bytes, err: %v", err)
//		}
//		if n == 0 {
//			return fmt.Errorf("wrote 0 bytes to file...")
//		}
//		// Use maxLen to determine when to flush the buffer
//		if w.maxLen != 0 {
//			if w.buf.Len() > w.maxLen {
//				_, err = w.buf.WriteTo(w.w)
//			}
//		} else {
//			_, err = w.buf.WriteTo(w.w)
//		}
//
//		if err != nil {
//			return fmt.Errorf("failed to append protobuf message to buffer bytes, err: %v", err)
//		}
//	} else {
//		data, err := proto.Marshal(m)
//		if err != nil {
//			return fmt.Errorf("failed to marshal protobuf message, err: %v", err)
//		}
//		n, err = w.w.Write(data)
//		if err != nil {
//			return fmt.Errorf("failed to write protobuf message, err: %v", err)
//		}
//	}
//
//	return nil
//}
//
//func (w writer) Close() error {
//	if w.buf != nil {
//		w.buf.WriteTo(w.w)
//	}
//	if closer, ok := w.w.(io.Closer); ok {
//		return closer.Close()
//	}
//	return nil
//}
