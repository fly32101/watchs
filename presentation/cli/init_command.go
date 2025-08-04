package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/watchs/domain/repository"
)

// InitCommand 初始化命令
type InitCommand struct {
	configRepo repository.ConfigRepository
}

// NewInitCommand 创建初始化命令
func NewInitCommand(configRepo repository.ConfigRepository) *InitCommand {
	return &InitCommand{
		configRepo: configRepo,
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
		fmt.Println("生成配置文件")
		fmt.Println("\n用法: watchs init [选项]")
		fmt.Println("\n选项:")
		initCmd.PrintDefaults()
		return nil
	}

	// 检查配置文件是否已存在
	if _, err := os.Stat(*configPath); err == nil && !*force {
		return fmt.Errorf("配置文件 %s 已存在，使用 -force 参数覆盖", *configPath)
	}

	// 创建配置
	config, err := createConfigFromArgs(*watchDir, *fileTypes, *excludePaths, *command)
	if err != nil {
		return fmt.Errorf("创建配置失败: %v", err)
	}

	// 保存配置
	if err := c.configRepo.SaveConfig(config, *configPath); err != nil {
		return fmt.Errorf("保存配置文件失败: %v", err)
	}

	// 显示配置信息
	fmt.Printf("配置文件已生成: %s\n", *configPath)
	fmt.Printf("监控目录: %s\n", config.WatchDir)
	if len(config.FileTypes) > 0 {
		fmt.Printf("监控的文件类型: %v\n", config.FileTypes)
	} else {
		fmt.Printf("监控所有文件类型\n")
	}
	if len(config.ExcludePaths) > 0 {
		fmt.Printf("排除的路径: %v\n", config.ExcludePaths)
	}
	fmt.Printf("执行命令: %s\n", config.Command)

	return nil
}
