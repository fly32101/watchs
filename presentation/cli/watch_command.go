package cli

import (
	"flag"
	"fmt"

	"github.com/watchs/application/interfaces"
)

// WatchCommand 监控命令
type WatchCommand struct {
	watchService interfaces.WatchApplicationService
}

// NewWatchCommand 创建监控命令
func NewWatchCommand(watchService interfaces.WatchApplicationService) *WatchCommand {
	return &WatchCommand{
		watchService: watchService,
	}
}

// Name 返回命令名称
func (c *WatchCommand) Name() string {
	return "watch"
}

// Description 返回命令描述
func (c *WatchCommand) Description() string {
	return "监控文件变化并执行命令"
}

// Execute 执行命令
func (c *WatchCommand) Execute(args []string) error {
	// 定义命令参数
	watchCmd := flag.NewFlagSet("watch", flag.ExitOnError)
	configPath := watchCmd.String("config", "watchs.json", "配置文件路径")
	watchDir := watchCmd.String("dir", "", "要监控的目录 (覆盖配置文件)")
	fileTypes := watchCmd.String("types", "", "要监控的文件类型，以逗号分隔，如 '.go,.js' (覆盖配置文件)")
	excludePaths := watchCmd.String("exclude", "", "要排除的路径，以逗号分隔 (覆盖配置文件)")
	command := watchCmd.String("cmd", "", "文件变化时执行的命令 (覆盖配置文件)")
	debounceMs := watchCmd.Int("debounce", 500, "防抖时间（毫秒）")
	showMemory := watchCmd.Bool("memory", false, "显示内存使用信息")
	memoryInterval := watchCmd.Int("memory-interval", 30, "内存信息显示间隔（秒）")
	help := watchCmd.Bool("help", false, "显示帮助信息")

	// 解析参数
	if err := watchCmd.Parse(args); err != nil {
		return err
	}

	// 显示帮助信息
	if *help {
		fmt.Println("监控文件变化并执行命令")
		fmt.Println("\n用法: watchs watch [选项]")
		fmt.Println("\n选项:")
		watchCmd.PrintDefaults()
		fmt.Println("\n示例:")
		fmt.Println("  watchs watch                           # 使用默认配置监控")
		fmt.Println("  watchs watch --memory                  # 监控时显示内存信息")
		fmt.Println("  watchs watch --memory --memory-interval 60  # 每60秒显示内存信息")
		return nil
	}

	// 创建监控配置参数
	watchConfig := &interfaces.WatchConfig{
		ConfigPath:     *configPath,
		WatchDir:       *watchDir,
		FileTypes:      *fileTypes,
		ExcludePaths:   *excludePaths,
		Command:        *command,
		DebounceMs:     *debounceMs,
		ShowMemory:     *showMemory,
		MemoryInterval: *memoryInterval,
	}

	// 启动监控服务
	return c.watchService.StartWatch(watchConfig)
}
