package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/watchs/application"
	"github.com/watchs/domain/entity"
	"github.com/watchs/domain/repository"
	"github.com/watchs/infrastructure/watcher"
)

// WatchCommand 监控命令
type WatchCommand struct {
	configRepo repository.ConfigRepository
}

// NewWatchCommand 创建监控命令
func NewWatchCommand(configRepo repository.ConfigRepository) *WatchCommand {
	return &WatchCommand{
		configRepo: configRepo,
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
		return nil
	}

	// 加载或创建配置
	config, err := c.loadOrCreateConfig(*configPath, *watchDir, *fileTypes, *excludePaths, *command)
	if err != nil {
		return fmt.Errorf("配置加载失败: %v", err)
	}

	// 创建文件监控服务
	fsWatcher, err := watcher.NewFSNotifyWatcher(config)
	if err != nil {
		return fmt.Errorf("创建文件监控器失败: %v", err)
	}

	// 创建命令执行器
	cmdExecutor := watcher.NewCommandExecutor(*debounceMs)

	// 创建应用服务
	watchService := application.NewWatchService(config, fsWatcher, cmdExecutor)

	// 启动监控
	if err := watchService.Start(); err != nil {
		return fmt.Errorf("启动监控失败: %v", err)
	}

	fmt.Println("监控已启动，按 Ctrl+C 停止...")

	// 等待中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("正在关闭监控...")
	if err := watchService.Stop(); err != nil {
		log.Printf("关闭监控失败: %v", err)
	}

	return nil
}

// loadOrCreateConfig 加载或创建配置
func (c *WatchCommand) loadOrCreateConfig(configPath, watchDir, fileTypes, excludePaths, command string) (*entity.WatchConfig, error) {
	// 尝试加载配置
	config, err := c.configRepo.LoadConfig(configPath)
	if err != nil {
		// 如果配置文件不存在，尝试使用命令行参数
		if os.IsNotExist(err) && watchDir != "" && command != "" {
			log.Printf("配置文件 %s 不存在，使用命令行参数", configPath)
			return createConfigFromArgs(watchDir, fileTypes, excludePaths, command)
		}
		return nil, err
	}

	// 命令行参数覆盖配置文件
	if watchDir != "" || fileTypes != "" || excludePaths != "" || command != "" {
		return overrideConfig(config, watchDir, fileTypes, excludePaths, command)
	}

	return config, nil
}
