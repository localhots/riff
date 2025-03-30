.PHONY: test bench pprof

test:
	go test -v .

bench:
	go test -v . -bench=. -benchmem -run=Benchmark


pprof:
	go build -o /tmp/clog prof/main.go
	/tmp/clog -cpuprofile=/tmp/clog.prof
	go tool pprof /tmp/clog /tmp/clog.prof
