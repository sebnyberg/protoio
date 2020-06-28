package protoio

type DelimiterType int

const (
	DelimiterTypeUint32 DelimiterType = 0
	DelimiterTypeUint64 DelimiterType = 1
)

const (
	LenBytesUint32 int = 4
	LenBytesUint64 int = 8
)
