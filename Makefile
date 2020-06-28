.PHONY: proto profile

proto:
	@protoc --go_out=. --go_opt=paths=source_relative test/test.proto

profile:
	@go test ./protoio_test.go \
		-bench=. -benchmem -memprofile=mem.pprof \
		-cpuprofile=cpu.pprof

	@go tool pprof -web protoio.test cpu.pprof
	@go tool pprof -web protoio.test mem.pprof
