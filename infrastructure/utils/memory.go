package utils

import (
	"fmt"
	"github.com/watchs/infrastructure/ui"
	"runtime"
	"runtime/debug"
	"time"
)

// MemoryStats 内存统计信息
type MemoryStats struct {
	Alloc      uint64 // 当前分配的内存 (bytes)
	TotalAlloc uint64 // 总分配的内存 (bytes)
	Sys        uint64 // 系统内存 (bytes)
	NumGC      uint32 // GC次数
	Goroutines int    // 当前goroutine数量
}

// GetMemoryStats 获取当前内存统计信息
func GetMemoryStats() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemoryStats{
		Alloc:      m.Alloc,
		TotalAlloc: m.TotalAlloc,
		Sys:        m.Sys,
		NumGC:      m.NumGC,
		Goroutines: runtime.NumGoroutine(),
	}
}

// ForceGC 强制执行垃圾回收
func ForceGC() {
	runtime.GC()
	debug.FreeOSMemory()
}

// formatBytes 格式化字节数为人类可读的格式
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// PrintMemoryStats 打印内存统计信息
func PrintMemoryStats(stats MemoryStats) {
	ui.PrintInfo(fmt.Sprintf("内存使用: %s | 系统内存: %s | Goroutines: %d | GC次数: %d",
		formatBytes(stats.Alloc),
		formatBytes(stats.Sys),
		stats.Goroutines,
		stats.NumGC))
}

// PrintDetailedMemoryStats 打印详细的内存统计信息
func PrintDetailedMemoryStats(stats MemoryStats) {
	ui.PrintHeader("内存统计信息")
	fmt.Printf("  当前分配内存: %s\n", formatBytes(stats.Alloc))
	fmt.Printf("  累计分配内存: %s\n", formatBytes(stats.TotalAlloc))
	fmt.Printf("  系统内存使用: %s\n", formatBytes(stats.Sys))
	fmt.Printf("  Goroutine数量: %d\n", stats.Goroutines)
	fmt.Printf("  垃圾回收次数: %d\n", stats.NumGC)
}

// StartMemoryMonitor 启动内存监控（用于调试）
func StartMemoryMonitor(interval time.Duration, callback func(stats MemoryStats)) chan struct{} {
	stopCh := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-stopCh:
				return
			case <-ticker.C:
				stats := GetMemoryStats()
				if callback != nil {
					callback(stats)
				}
			}
		}
	}()

	return stopCh
}

// SetGCPercent 设置GC触发百分比
func SetGCPercent(percent int) int {
	return debug.SetGCPercent(percent)
}

// SetMaxProcs 设置最大CPU核心数
func SetMaxProcs(n int) int {
	return runtime.GOMAXPROCS(n)
}
