package main

import (
	"log"

	"github.com/watchs/presentation/cli"
	"github.com/watchs/presentation/cli/factory"
)

// 版本信息，将由GoReleaser在构建时注入
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// 设置日志格式
	log.SetFlags(log.Ldate | log.Ltime)

	// 设置版本信息
	cli.Version = version
	cli.Commit = commit
	cli.Date = date

	// 使用工厂创建CLI
	cliFactory := factory.NewCLIFactory()
	cliInstance := cliFactory.CreateCLI()
	cliInstance.Run()
}