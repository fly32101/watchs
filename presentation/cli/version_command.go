package cli

import (
	"fmt"

	"github.com/watchs/presentation/cli/ui"
)

// 版本信息，将由main包传入
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// VersionCommand 版本命令
type VersionCommand struct{}

// NewVersionCommand 创建版本命令
func NewVersionCommand() *VersionCommand {
	return &VersionCommand{}
}

// Name 返回命令名称
func (c *VersionCommand) Name() string {
	return "version"
}

// Description 返回命令描述
func (c *VersionCommand) Description() string {
	return "显示版本信息"
}

// Execute 执行命令
func (c *VersionCommand) Execute(args []string) error {
	ui.PrintHeader("Watchs 版本信息")
	fmt.Printf("%s%s%s version %s%s%s, commit %s%s%s, built at %s%s%s\n",
		ui.Blue, ui.Rocket, ui.Reset,
		ui.Green, Version, ui.Reset,
		ui.Yellow, Commit, ui.Reset,
		ui.Purple, Date, ui.Reset)
	return nil
}
