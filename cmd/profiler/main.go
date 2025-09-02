package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"logparser/parser"
	"logparser/profiler"
)

func main() {
	var (
		enableProfiler = flag.Bool("profile", false, "启用性能分析")
		cpuProfile     = flag.String("cpuprofile", "cpu.prof", "CPU profile 输出文件")
		httpProfiler   = flag.String("httpprof", "", "HTTP 性能分析服务器地址 (例如: :6060)")
		showStats      = flag.Bool("stats", false, "显示运行时统计信息")
	)
	flag.Parse()

	var prof *profiler.Profiler
	if *enableProfiler {
		prof = profiler.NewProfiler()

		// 启动 CPU 性能分析
		if err := prof.StartCPUProfile(*cpuProfile); err != nil {
			fmt.Fprintf(os.Stderr, "启动 CPU profile 失败: %v\n", err)
			os.Exit(1)
		}
		defer prof.StopCPUProfile()

		// 启动 HTTP 性能分析服务器
		if *httpProfiler != "" {
			prof.StartHTTPProfiler(*httpProfiler)
		}

		// 设置信号处理器
		prof.SetupSignalHandlers()

		fmt.Println("性能分析已启用")
	}

	// 创建日志解析器
	logParser := parser.NewLogParser()

	// 创建带缓冲的读取器
	scanner := bufio.NewScanner(os.Stdin)

	// 设置更大的缓冲区以处理长行
	const maxCapacity = 1024 * 1024 // 1MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	lineCount := 0
	startTime := time.Now()

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parsed, err := logParser.Parse(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "解析错误: %v\n", err)
			continue
		}

		fmt.Println(parsed.String())
		lineCount++

		// 定期显示统计信息
		if *showStats && lineCount%1000 == 0 {
			elapsed := time.Since(startTime)
			rate := float64(lineCount) / elapsed.Seconds()
			fmt.Fprintf(os.Stderr, "已处理 %d 行，速率: %.2f 行/秒\n", lineCount, rate)

			if prof != nil {
				prof.PrintGoroutines()
				prof.PrintMemoryStats()
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "读取错误: %v\n", err)
		os.Exit(1)
	}

	// 最终统计
	elapsed := time.Since(startTime)
	fmt.Fprintf(os.Stderr, "处理完成: 共处理 %d 行，耗时 %.2f 秒，平均速率 %.2f 行/秒\n",
		lineCount, elapsed.Seconds(), float64(lineCount)/elapsed.Seconds())
}
