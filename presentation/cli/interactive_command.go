package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/watchs/application"
	"github.com/watchs/domain/repository"
	"github.com/watchs/infrastructure/watcher"
	"github.com/watchs/presentation/cli/ui"
)

// InteractiveCommand 交互式命令
type InteractiveCommand struct {
	configRepo repository.ConfigRepository
}

// NewInteractiveCommand 创建交互式命令
func NewInteractiveCommand(configRepo repository.ConfigRepository) *InteractiveCommand {
	return &InteractiveCommand{
		configRepo: configRepo,
	}
}

// Name 返回命令名称
func (c *InteractiveCommand) Name() string {
	return "interactive"
}

// Description 返回命令描述
func (c *InteractiveCommand) Description() string {
	return "交互式配置向导"
}

// Execute 执行命令
func (c *InteractiveCommand) Execute(args []string) error {
	ui.PrintHeader("欢迎使用 Watchs 文件监控工具配置向导")
	ui.PrintInfo("请回答以下问题来创建配置文件")
	fmt.Println("----------------------------------------")

	// 创建交互式CLI
	interactiveCLI := NewInteractiveCLI()

	// 运行交互式配置
	config, configPath, err := interactiveCLI.Run()
	if err != nil {
		ui.PrintError(fmt.Sprintf("交互式配置失败: %v", err))
		return fmt.Errorf("交互式配置失败: %v", err)
	}

	// 保存配置
	if err := c.configRepo.SaveConfig(config, configPath); err != nil {
		ui.PrintError(fmt.Sprintf("保存配置文件失败: %v", err))
		return fmt.Errorf("保存配置文件失败: %v", err)
	}

	ui.PrintSuccess(fmt.Sprintf("配置文件已保存到: %s", configPath))
	ui.PrintInfo("你可以通过以下命令启动监控:")
	fmt.Printf("watchs -config %s\n", configPath)

	// 询问是否立即启动监控
	startNow := interactiveCLI.askYesNo("是否立即启动监控？[Y/n]: ", true)
	if startNow {
		ui.PrintInfo("正在初始化监控服务...")

		// 模拟加载动画
		ui.SimulateLoading(2*time.Second, "初始化监控器")

		// 创建文件监控服务
		fsWatcher, err := watcher.NewFSNotifyWatcher(config)
		if err != nil {
			ui.PrintError(fmt.Sprintf("创建文件监控器失败: %v", err))
			return fmt.Errorf("创建文件监控器失败: %v", err)
		}

		// 创建命令执行器
		cmdExecutor := watcher.NewCommandExecutor(500) // 默认防抖时间

		// 创建应用服务
		watchService := application.NewWatchService(config, fsWatcher, cmdExecutor)

		// 启动监控
		if err := watchService.Start(); err != nil {
			ui.PrintError(fmt.Sprintf("启动监控失败: %v", err))
			return fmt.Errorf("启动监控失败: %v", err)
		}

		ui.PrintSuccess(fmt.Sprintf("监控已启动，正在监控目录: %s", config.WatchDir))
		ui.PrintInfo("按 Ctrl+C 停止监控...")

		// 等待中断信号
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		ui.PrintWarning("正在关闭监控...")
		if err := watchService.Stop(); err != nil {
			ui.PrintError(fmt.Sprintf("关闭监控失败: %v", err))
		}

		ui.PrintSuccess("监控已成功关闭!")
	}

	return nil
}
