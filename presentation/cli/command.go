package cli

import (
	"fmt"

	"github.com/watchs/domain/repository"
)

// Command 命令接口
type Command interface {
	// Name 返回命令名称
	Name() string
	// Description 返回命令描述
	Description() string
	// Execute 执行命令
	Execute(args []string) error
}

// CommandRegistry 命令注册表
type CommandRegistry struct {
	commands   map[string]Command
	configRepo repository.ConfigRepository
}

// NewCommandRegistry 创建命令注册表
func NewCommandRegistry(configRepo repository.ConfigRepository) *CommandRegistry {
	return &CommandRegistry{
		commands:   make(map[string]Command),
		configRepo: configRepo,
	}
}

// Register 注册命令
func (r *CommandRegistry) Register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

// Get 获取命令
func (r *CommandRegistry) Get(name string) (Command, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}

// GetDefaultCommand 获取默认命令
func (r *CommandRegistry) GetDefaultCommand() Command {
	return r.commands["watch"]
}

// ListCommands 列出所有命令
func (r *CommandRegistry) ListCommands() []Command {
	var cmds []Command
	for _, cmd := range r.commands {
		cmds = append(cmds, cmd)
	}
	return cmds
}

// ShowHelp 显示帮助信息
func (r *CommandRegistry) ShowHelp() {
	fmt.Println("Watchs - 文件变更监控工具")
	fmt.Println("\n可用命令:")
	for _, cmd := range r.commands {
		fmt.Printf("  %s\t%s\n", cmd.Name(), cmd.Description())
	}
	fmt.Println("\n使用 'watchs <命令> --help' 获取更多信息")
}
