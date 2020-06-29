# ProtoIO

Helper functions for writing length-delimited Protobuf messages in Go.

## Why do I need this?

Because of Protobuf's binary format, it is not possible to write multiple messages to a single file or stream without losing the ability to decode the messages again. 

To circumvent this problem, people tend to write massive wrapper messages instead. Since each message needs to be decoded in its entirety in memory before it can be accessed, this is both dangerous and inefficient.

To enable streaming Protobuf messages, this library prefixes each message by its length before writing it to the stream.

## Basic usage

__NOTE__: This package assumes the use of the Protobuf V2 API.

### In-memory

```go
// Write one message
protoWriter := protoio.NewWriter(&buf)

m := pb.User{ Name: "Seb" }
err := protoWriter.WriteMsg(&m)
if err != nil {
    log.Fatal(err)
}

err = protoWriter.Close()
if err != nil {
    log.Fatal(err)
}

// Read one message
protoReader := protoio.NewReader(&buf)

var out pb.User
err = protoReader.ReadMsg(&out)
if err != nil {
    log.Fatal(err)
}
```

### To/from a file

```go
// Open output file
outFile, err := os.OpenFile("out.ldproto", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
if err != nil {
    log.Fatal(err)
}

// Create a writer with 1MB buffer
protoWriter := protoio.NewWriter(outFile, protoio.WriteWithBufIO(1024*1024))

// Write the same message 10k times
m := pb.User{ Name: "Seb" }
for i := 0; i < 10000; i++ {
    err := protoWriter.WriteMsg(&m)
    if err != nil {
        log.Fatal(err)
    }
}

// Remember to close! 
err = protoWriter.Close()
if err != nil {
    log.Fatal(err)
}

// Open the file again
inFile, err := os.OpenFile("out.ldproto", os.O_RWONLY, 0644)
if err != nil {
    log.Fatal(err)
}

// Read a ton of Sebs
protoReader := protoio.NewReader(inFile, protoio.ReadWithBufIO(1024*1024))
aTonOfSebs := make([]pb.User, 10000)

for i := 0; i < 10000; i++ {
    err = protoReader.ReadMsg(&aTonOfSebs[i])
    if err != nil {
        log.Fatal(err)
    }
}
if err := protoReader.Close(); err != nil {
    log.Fatal(err)
}
```

## Options

| Option | Default | Description | 
| --- | --- | --- |
| `ReaderMsgSizeMax(size int)` | 4MB (1024 * 1024 * 4) | Maximum allowed message size when reading |
| `ReaderWithBufIO(bufSize int)` | disabled | Wrap in BufIO reader. Useful for reading from files |
| `ReaderDelimiter(bo binary.ByteOrder, delimType DelimiterType)` | binary.BigEndian, DelimiterTypeUint32 | Length delimiter type |
| `WriterWithBufIO(bufSize int)` | disabled | Wrap in BufIO writer. Useful for writing to files |
| `WriterDelimiter(bo binary.ByteOrder, delimType DelimiterType)` | binary.BigEndian, DelimiterTypeUint32 | Length delimiter type |

## Benchmarks

Benchmark for reading/writing 1k messages.

When working with files, provide `protoio.ReadWithBufIO(size)` and `protoio.WriteWithBufIO(size)` as options to the reader and write respectively.

As can be seen in the benchmarks below, using BufIO when reading to/from improves performance significantly.

```bash
$ go test . -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/sebnyberg/protoio
BenchmarkReader_ReadMsg_FromMemoryBuffer_1000-6             5055            237270 ns/op           16032 B/op       2001 allocs/op
BenchmarkReader_ReadMsg_FromFile/Direct_1000-6               612           1864892 ns/op           16032 B/op       2001 allocs/op
BenchmarkReader_ReadMsg_FromFile/WithBufIO_1000-6           4934            243272 ns/op           16032 B/op       2001 allocs/op
BenchmarkWriter_WriteMsg/1000_InMem-6                       6044            198506 ns/op              32 B/op          1 allocs/op
BenchmarkWriter_WriteMsg/1000_ToFile-6                       210           5658680 ns/op              32 B/op          1 allocs/op
BenchmarkWriter_WriteMsg/1000_ToFileBufIO-6                 5760            207997 ns/op              32 B/op          1 allocs/op
```
