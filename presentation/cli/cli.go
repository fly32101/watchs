package cli

import (
	"log"
	"os"
)

// CLI 表示命令行界面
type CLI struct {
	registry *CommandRegistry
}

// NewCLIWithRegistry 使用指定的命令注册表创建CLI
func NewCLIWithRegistry(registry *CommandRegistry) *CLI {
	return &CLI{
		registry: registry,
	}
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
