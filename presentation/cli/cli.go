package cli

import (
	"log"
	"os"

	"github.com/watchs/infrastructure/persistence"
)

// CLI 表示命令行界面
type CLI struct {
	registry *CommandRegistry
}

// NewCLI 创建一个新的命令行界面
func NewCLI() *CLI {
	// 创建配置仓储
	configRepo := persistence.NewJsonConfigRepository()

	// 创建命令注册表
	registry := NewCommandRegistry(configRepo)

	// 创建CLI
	cli := &CLI{
		registry: registry,
	}

	// 注册命令
	registry.Register(NewWatchCommand(configRepo))
	registry.Register(NewInitCommand(configRepo))
	registry.Register(NewInteractiveCommand(configRepo))
	registry.Register(NewVersionCommand())

	// 注册帮助命令（需要在其他命令注册后）
	registry.Register(NewHelpCommand(registry))

	return cli
}

// Run 运行命令行界面
func (c *CLI) Run() {
	// 如果没有提供参数，显示帮助信息
	if len(os.Args) < 2 {
		// 默认执行watch命令
		cmd := c.registry.GetDefaultCommand()
		if err := cmd.Execute(nil); err != nil {
			log.Printf("错误: %v", err)
		}
		return
	}

	// 获取命令名称
	cmdName := os.Args[1]

	// 查找命令
	cmd, ok := c.registry.Get(cmdName)
	if !ok {
		// 如果命令不存在，尝试将参数作为watch命令的参数
		cmd = c.registry.GetDefaultCommand()
		if err := cmd.Execute(os.Args[1:]); err != nil {
			log.Printf("错误: %v", err)
		}
		return
	}

	// 执行命令
	if err := cmd.Execute(os.Args[2:]); err != nil {
		log.Printf("错误: %v", err)
		os.Exit(1)
	}
}
