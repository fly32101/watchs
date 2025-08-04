package cli

import "fmt"

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
			return fmt.Errorf("未知命令: %s", cmdName)
		}

		// 通过传递 --help 参数给命令来显示帮助信息
		return cmd.Execute([]string{"--help"})
	}

	// 显示所有命令的帮助信息
	c.registry.ShowHelp()
	return nil
}
