.PHONY: build test benchmark clean profile help



all: build

build:
	@echo "Build Log Parser..."
	go build -o logparser main.go
	go build -o logparser-profiler cmd/profiler/main.go
	@echo "build finished！"

test:
	@echo "Run unit tests..."
	go test -v

benchmark:
	@echo "Run benchmark tests..."
	go test -bench=. -benchmem

benchmark-report:
	@echo "Generate benchmark report..."
	go test -bench=. -benchmem > benchmark_results.txt
	@echo "benchmark report saved to benchmark_results.txt"

coverage:
	@echo "Run test coverage..."
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Test coverage report saved to coverage.html"

profile:
	@echo "Run cpu performance analysis..."
	go test -cpuprofile=cpu.prof -bench=.
	@echo "CPU profile: cpu.prof"

memprofile:
	@echo "Run memory performance analysis..."
	go test -memprofile=mem.prof -bench=.
	@echo "memory profile: mem.prof"

clean:
	@echo "Clean genetated files..."
	rm -f logparser logparser-profiler
	rm -f *.prof
	rm -f coverage.out coverage.html
	rm -f benchmark_results.txt
	@echo "Clean finished! "

example:
	@echo "Test basic functionality..."
	echo "2024-01-01 12:00:00 [INFO] test log info" | ./logparser

example-profiler:
	@echo "Test with profiler..."
	echo "2024-01-01 12:00:00 [INFO] test log info" | ./logparser-profiler -profile -stats

help:
	@echo "Log Parser Makefile Commands:"
	@echo "  build            - build log parser"
	@echo "  test             - Run unit tests"
	@echo "  benchmark        - Run benchmark tests"
	@echo "  benchmark-report - Generate benchmark report"
	@echo "  coverage         - Run test coverage"
	@echo "  profile          - Generate CPU profile"
	@echo "  memprofile       - Generate memory profile"
	@echo "  clean            - Clean generated files"
	@echo "  example          - Test basic functionality"
	@echo "  example-profiler - Test with profiler"
	@echo "  help             - Display this help message"
