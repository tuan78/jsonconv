.PHONY: test
test:
	go test -cover -race ./...

.PHONY: cover
cover:
	go test -race -coverprofile=cover.out -coverpkg=./... ./...
	go tool cover -html=cover.out -o cover.html

.PHONY: bench
bench:
	cd ./benchmarks && go test -bench=. -run="^$$" -cpuprofile=cpu.pprof -memprofile=mem.pprof -benchmem