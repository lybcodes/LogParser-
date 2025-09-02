package profiler

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
)

// Profiler 性能分析器
type Profiler struct {
	cpuProfile *os.File
	memProfile *os.File
	blockProfile *os.File
	mutexProfile *os.File
}

// NewProfiler 创建新的性能分析器
func NewProfiler() *Profiler {
	return &Profiler{}
}

// StartCPUProfile 开始 CPU 性能分析
func (p *Profiler) StartCPUProfile(filename string) error {
	var err error
	p.cpuProfile, err = os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建 CPU profile 文件失败: %v", err)
	}
	
	if err := pprof.StartCPUProfile(p.cpuProfile); err != nil {
		p.cpuProfile.Close()
		return fmt.Errorf("启动 CPU profile 失败: %v", err)
	}
	
	fmt.Printf("CPU 性能分析已启动，输出到: %s\n", filename)
	return nil
}

// StopCPUProfile 停止 CPU 性能分析
func (p *Profiler) StopCPUProfile() {
	if p.cpuProfile != nil {
		pprof.StopCPUProfile()
		p.cpuProfile.Close()
		fmt.Println("CPU 性能分析已停止")
	}
}

// WriteMemProfile 写入内存性能分析
func (p *Profiler) WriteMemProfile(filename string) error {
	var err error
	p.memProfile, err = os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建内存 profile 文件失败: %v", err)
	}
	defer p.memProfile.Close()
	
	if err := pprof.WriteHeapProfile(p.memProfile); err != nil {
		return fmt.Errorf("写入内存 profile 失败: %v", err)
	}
	
	fmt.Printf("内存性能分析已写入: %s\n", filename)
	return nil
}

// WriteBlockProfile 写入阻塞性能分析
func (p *Profiler) WriteBlockProfile(filename string) error {
	var err error
	p.blockProfile, err = os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建阻塞 profile 文件失败: %v", err)
	}
	defer p.blockProfile.Close()
	
	if err := pprof.Lookup("block").WriteTo(p.blockProfile, 0); err != nil {
		return fmt.Errorf("写入阻塞 profile 失败: %v", err)
	}
	
	fmt.Printf("阻塞性能分析已写入: %s\n", filename)
	return nil
}

// WriteMutexProfile 写入互斥锁性能分析
func (p *Profiler) WriteMutexProfile(filename string) error {
	var err error
	p.mutexProfile, err = os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建互斥锁 profile 文件失败: %v", err)
	}
	defer p.mutexProfile.Close()
	
	if err := pprof.Lookup("mutex").WriteTo(p.mutexProfile, 0); err != nil {
		return fmt.Errorf("写入互斥锁 profile 失败: %v", err)
	}
	
	fmt.Printf("互斥锁性能分析已写入: %s\n", filename)
	return nil
}

// StartHTTPProfiler 启动 HTTP 性能分析服务器
func (p *Profiler) StartHTTPProfiler(addr string) {
	go func() {
		fmt.Printf("HTTP 性能分析服务器已启动，访问: http://%s/debug/pprof/\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Printf("HTTP 性能分析服务器启动失败: %v", err)
		}
	}()
}

// SetupSignalHandlers 设置信号处理器
func (p *Profiler) SetupSignalHandlers() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		fmt.Println("\n收到退出信号，正在生成性能分析文件...")
		
		// 生成各种性能分析文件
		p.WriteMemProfile("mem.prof")
		p.WriteBlockProfile("block.prof")
		p.WriteMutexProfile("mutex.prof")
		
		fmt.Println("性能分析文件生成完成，程序退出")
		os.Exit(0)
	}()
}

// PrintGoroutines 打印当前 goroutine 数量
func (p *Profiler) PrintGoroutines() {
	fmt.Printf("当前 goroutine 数量: %d\n", runtime.NumGoroutine())
}

// PrintMemoryStats 打印内存统计信息
func (p *Profiler) PrintMemoryStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	fmt.Printf("内存统计:\n")
	fmt.Printf("  已分配内存: %d bytes (%.2f MB)\n", m.Alloc, float64(m.Alloc)/1024/1024)
	fmt.Printf("  总分配内存: %d bytes (%.2f MB)\n", m.TotalAlloc, float64(m.TotalAlloc)/1024/1024)
	fmt.Printf("  系统内存: %d bytes (%.2f MB)\n", m.Sys, float64(m.Sys)/1024/1024)
	fmt.Printf("  GC 次数: %d\n", m.NumGC)
}
