package cli

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/watchs/infrastructure/ui"
	"github.com/watchs/infrastructure/utils"
)

// MemoryCommand 内存信息命令
type MemoryCommand struct{}

// NewMemoryCommand 创建内存信息命令
func NewMemoryCommand() *MemoryCommand {
	return &MemoryCommand{}
}

// Name 返回命令名称
func (c *MemoryCommand) Name() string {
	return "memory"
}

// Description 返回命令描述
func (c *MemoryCommand) Description() string {
	return "显示内存使用信息"
}

// Execute 执行命令
func (c *MemoryCommand) Execute(args []string) error {
	// 定义命令参数
	memCmd := flag.NewFlagSet("memory", flag.ExitOnError)
	detailed := memCmd.Bool("detailed", false, "显示详细的内存信息")
	monitor := memCmd.Bool("monitor", false, "启动内存监控模式")
	interval := memCmd.Int("interval", 5, "监控间隔（秒）")
	gc := memCmd.Bool("gc", false, "执行垃圾回收后显示内存信息")
	help := memCmd.Bool("help", false, "显示帮助信息")

	// 解析参数
	if err := memCmd.Parse(args); err != nil {
		return err
	}

	// 显示帮助信息
	if *help {
		ui.PrintHeader("内存信息命令")
		fmt.Println("\n用法: watchs memory [选项]")
		fmt.Println("\n选项:")
		memCmd.PrintDefaults()
		fmt.Println("\n示例:")
		fmt.Println("  watchs memory                    # 显示当前内存信息")
		fmt.Println("  watchs memory --detailed         # 显示详细内存信息")
		fmt.Println("  watchs memory --gc               # 执行GC后显示内存信息")
		fmt.Println("  watchs memory --monitor          # 启动内存监控")
		fmt.Println("  watchs memory --monitor --interval 10  # 每10秒监控一次")
		return nil
	}

	// 执行垃圾回收
	if *gc {
		ui.PrintInfo("正在执行垃圾回收...")
		utils.ForceGC()
		ui.PrintSuccess("垃圾回收完成")
	}

	// 获取内存统计信息
	stats := utils.GetMemoryStats()

	if *monitor {
		// 监控模式
		ui.PrintHeader("内存监控模式")
		ui.PrintInfo(fmt.Sprintf("监控间隔: %d秒 (按 Ctrl+C 停止)", *interval))
		fmt.Println()

		// 显示初始状态
		if *detailed {
			utils.PrintDetailedMemoryStats(stats)
		} else {
			utils.PrintMemoryStats(stats)
		}

		// 启动监控
		stopCh := utils.StartMemoryMonitor(time.Duration(*interval)*time.Second, func(stats utils.MemoryStats) {
			if *detailed {
				fmt.Println() // 换行
				utils.PrintDetailedMemoryStats(stats)
			} else {
				utils.PrintMemoryStats(stats)
			}
		})

		// 等待用户中断
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		// 停止监控
		close(stopCh)
		signal.Stop(sigCh)
		close(sigCh)

		ui.PrintWarning("内存监控已停止")
	} else {
		// 单次显示模式
		if *detailed {
			utils.PrintDetailedMemoryStats(stats)
		} else {
			utils.PrintMemoryStats(stats)
		}
	}

	return nil
}
