package cli

import (
	"fmt"

	"github.com/watchs/application/interfaces"
	"github.com/watchs/infrastructure/ui"
)

// InteractiveCommand 交互式命令
type InteractiveCommand struct {
	configService interfaces.ConfigApplicationService
	watchService  interfaces.WatchApplicationService
}

// NewInteractiveCommand 创建交互式命令
func NewInteractiveCommand(configService interfaces.ConfigApplicationService, watchService interfaces.WatchApplicationService) *InteractiveCommand {
	return &InteractiveCommand{
		configService: configService,
		watchService:  watchService,
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

	// 运行交互式配置
	config, configPath, err := c.configService.RunInteractiveConfig()
	if err != nil {
		ui.PrintError(fmt.Sprintf("交互式配置失败: %v", err))
		return fmt.Errorf("交互式配置失败: %v", err)
	}

	// 保存配置
	if err := c.configService.SaveConfig(config, configPath); err != nil {
		ui.PrintError(fmt.Sprintf("保存配置文件失败: %v", err))
		return fmt.Errorf("保存配置文件失败: %v", err)
	}

	ui.PrintSuccess(fmt.Sprintf("配置文件已保存到: %s", configPath))
	ui.PrintInfo("你可以通过以下命令启动监控:")
	fmt.Printf("watchs -config %s\n", configPath)

	// 询问是否立即启动监控
	interactiveCLI := NewInteractiveCLI()
	startNow := interactiveCLI.AskYesNo("是否立即启动监控？[Y/n]: ", true)
	if startNow {
		// 创建监控配置参数
		watchConfig := &interfaces.WatchConfig{
			ConfigPath:   configPath,
			WatchDir:     "",
			FileTypes:    "",
			ExcludePaths: "",
			Command:      "",
			DebounceMs:   500,
		}

		// 启动监控服务
		return c.watchService.StartWatch(watchConfig)
	}

	return nil
}
