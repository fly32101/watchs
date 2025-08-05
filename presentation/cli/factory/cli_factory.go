package factory

import (
	"github.com/watchs/infrastructure/di"
	"github.com/watchs/presentation/cli"
)

// CLIFactory CLI工厂
type CLIFactory struct {
	container *di.Container
}

// NewCLIFactory 创建CLI工厂
func NewCLIFactory() *CLIFactory {
	return &CLIFactory{
		container: di.NewContainer(),
	}
}

// CreateCLI 创建CLI实例
func (f *CLIFactory) CreateCLI() *cli.CLI {
	// 创建命令注册表
	registry := cli.NewCommandRegistry(f.container.GetConfigRepository())

	// 创建CLI
	cliInstance := cli.NewCLIWithRegistry(registry)

	// 注册命令
	registry.Register(cli.NewWatchCommand(f.container.GetWatchApplicationService()))
	registry.Register(cli.NewInitCommand(f.container.GetConfigApplicationService()))
	registry.Register(cli.NewInteractiveCommand(f.container.GetConfigApplicationService(), f.container.GetWatchApplicationService()))
	registry.Register(cli.NewVersionCommand())
	registry.Register(cli.NewMemoryCommand())

	// 注册帮助命令（需要在其他命令注册后）
	registry.Register(cli.NewHelpCommand(registry))

	return cliInstance
}
