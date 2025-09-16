# Log Parser

An efficient unified log format parser implemented in Go language without third-party dependencies.

## Features

- Supports multiple log formats：
  - Standard format: `2024-01-01 12:00:00 [INFO] content`
  - ISO 8601 format: `2024-01-01T12:00:00Z [INFO] content`
  - JSON fromat: `{"timestamp": "2024-01-01T12:00:00Z", "level": "INFO", "message": "content"}`
  - Simple format: `[INFO] content`
- Read from stdin and output to stdout
- High-performance design supporting large file prcessing
- Built-in performance profiling tools
- Complete benchmark test suite

## Quick Start

### Basic Usage

```bash
# Compile the program
go build -o logparser

# Read logs from file

cat logfile.txt | ./logparser
# Read logs from pipe
echo "2024-01-01 12:00:00 [INFO] 测试日志" | ./logparser
```

### Usage with Performance Profiling

```bash
# Enable performance profiling
./logparser -profile -cpuprofile=cpu.prof -httpprof=:6060 -stats

# Access performance profiling interface
# Open browser to http://localhost:6060/debug/pprof/
```

## Performance Testing

### Running Benchmarks

```bash
# Run all benchmark tests
go test -bench=.

# Run specific benchmark test
go test -bench=BenchmarkParseStandardFormat

# Generate benchmark report
go test -bench=. -benchmem > benchmark_results.txt
```

### Performance Analysis

```bash
# Generate CPU profile
go test -cpuprofile=cpu.prof -bench=.

# Generate memory profile
go test -memprofile=mem.prof -bench=.

# Analyze profile files
go tool pprof cpu.prof
go tool pprof mem.prof
```

### Flame Graph Generation

```bash
# Install go-torch (requires graphviz installation first)
go install github.com/uber/go-torch@latest

# Generate flame graph
go-torch cpu.prof
```

## Project Structure

```
.
├── main.go                 # Main program entry
├── cmd/
│   └── profiler/
│       └── main.go        # Main program with performance profiling
├── parser/
│   ├── parser.go          # Core parser implementation
│   ├── parser_test.go     # Unit tests
│   └── benchmark_test.go  # Benchmark tests
├── profiler/
│   └── profiler.go        # Performance profiling tools
├── go.mod                 # Go module file
├── Makefile               # Build and test scripts
├── performance_test.sh    # Performance testing script
├── sample_logs.txt        # Sample log file
└── README.md              # Project documentation
```

## Supported Log Formats

### 1. Standard Format

```
2024-01-01 12:00:00 [INFO] This is an info log
2024-01-01 12:00:00.123 [ERROR] This is an error log
```

### 2. ISO 8601 Format

```
2024-01-01T12:00:00Z [WARN] This is an warning log
2024-01-01T12:00:00.000 [DEBUG] This is an debug log
```

### 3. JSON Format

```
{"timestamp": "2024-01-01T12:00:00Z", "level": "INFO", "message": "JSON format log"}
```

### 4. Simple Format

```
[INFO] Simple format log
[ERROR] Error log
```

## Performance Optimization Features

- Uses pre-compiled regular expressions
- Efficient string processing
- Minimizes memory allocations
- Supports large buffer reading
- Built-in performance monitoring

## Command Line Arguments

- `-profile`: Enable performance profiling
- `-cpuprofile`: Specify CPU profile output file (default: cpu.prof)
- `-httpprof`: Start HTTP performance profiling server (e.g., :6060)
- `-stats`: Display runtime statistics


Actual performance depends on hardware configuration and log line length.
