package cli

import (
	"flag"
	"fmt"

	"github.com/watchs/application/interfaces"
	"github.com/watchs/infrastructure/ui"
)

// InitCommand 初始化命令
type InitCommand struct {
	configService interfaces.ConfigApplicationService
}

// NewInitCommand 创建初始化命令
func NewInitCommand(configService interfaces.ConfigApplicationService) *InitCommand {
	return &InitCommand{
		configService: configService,
	}
}

// Name 返回命令名称
func (c *InitCommand) Name() string {
	return "init"
}

// Description 返回命令描述
func (c *InitCommand) Description() string {
	return "生成配置文件"
}

// Execute 执行命令
func (c *InitCommand) Execute(args []string) error {
	// 定义命令参数
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	configPath := initCmd.String("config", "watchs.json", "配置文件路径")
	watchDir := initCmd.String("dir", "./", "要监控的目录")
	fileTypes := initCmd.String("types", "", "要监控的文件类型，以逗号分隔，如 '.go,.js'")
	excludePaths := initCmd.String("exclude", "", "要排除的路径，以逗号分隔")
	command := initCmd.String("cmd", "echo 文件已更新", "文件变化时执行的命令")
	force := initCmd.Bool("force", false, "是否强制覆盖已存在的配置文件")
	help := initCmd.Bool("help", false, "显示帮助信息")

	// 解析参数
	if err := initCmd.Parse(args); err != nil {
		return err
	}

	// 显示帮助信息
	if *help {
		ui.PrintHeader("生成配置文件")
		fmt.Println("\n用法: watchs init [选项]")
		fmt.Println("\n选项:")
		initCmd.PrintDefaults()
		return nil
	}

	// 创建初始化参数
	params := &interfaces.InitConfigParams{
		ConfigPath:   *configPath,
		WatchDir:     *watchDir,
		FileTypes:    *fileTypes,
		ExcludePaths: *excludePaths,
		Command:      *command,
		Force:        *force,
	}

	// 调用配置服务初始化配置
	return c.configService.InitializeConfig(params)
}
