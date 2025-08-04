package cli

import (
	"fmt"
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
	fmt.Printf("watchs version %s, commit %s, built at %s\n", Version, Commit, Date)
	return nil
}
