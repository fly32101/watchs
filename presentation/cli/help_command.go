package cli

import (
	"fmt"

	"github.com/watchs/infrastructure/ui"
)

// HelpCommand 帮助命令
type HelpCommand struct {
	registry *CommandRegistry
}

// NewHelpCommand 创建帮助命令
func NewHelpCommand(registry *CommandRegistry) *HelpCommand {
	return &HelpCommand{
		registry: registry,
	}
}

// Name 返回命令名称
func (c *HelpCommand) Name() string {
	return "help"
}

// Description 返回命令描述
func (c *HelpCommand) Description() string {
	return "显示帮助信息"
}

// Execute 执行命令
func (c *HelpCommand) Execute(args []string) error {
	if len(args) > 0 {
		// 显示特定命令的帮助信息
		cmdName := args[0]
		cmd, ok := c.registry.Get(cmdName)
		if !ok {
			ui.PrintError(fmt.Sprintf("未知命令: %s", cmdName))
			return fmt.Errorf("未知命令: %s", cmdName)
		}

		// 通过传递 --help 参数给命令来显示帮助信息
		return cmd.Execute([]string{"--help"})
	}

	// 显示所有命令的帮助信息
	ui.PrintHeader("Watchs - 文件变更监控工具")
	fmt.Println("\n可用命令:")

	commands := c.registry.ListCommands()
	for _, cmd := range commands {
		fmt.Printf("  %s%s%s\t%s\n", ui.Green, cmd.Name(), ui.Reset, cmd.Description())
	}

	fmt.Println("\n使用 'watchs <命令> --help' 获取更多信息")
	return nil
}
