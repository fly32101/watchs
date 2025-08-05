package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/watchs/domain/entity"
	"github.com/watchs/presentation/cli/ui"
)

// InteractiveCLI 交互式命令行界面
type InteractiveCLI struct {
	reader *bufio.Reader
}

// NewInteractiveCLI 创建新的交互式命令行
func NewInteractiveCLI() *InteractiveCLI {
	return &InteractiveCLI{
		reader: bufio.NewReader(os.Stdin),
	}
}

// Run 运行交互式配置向导
func (cli *InteractiveCLI) Run() (*entity.WatchConfig, string, error) {
	ui.PrintHeader("欢迎使用 Watchs 文件监控工具配置向导")
	ui.PrintInfo("请回答以下问题来创建配置文件")
	fmt.Println("----------------------------------------")

	// 获取配置文件路径
	configPath := cli.askString("配置文件路径", "watchs.json")

	// 获取监控目录
	watchDir := cli.askString("要监控的目录", "./")
	absWatchDir, err := filepath.Abs(watchDir)
	if err != nil {
		ui.PrintError(fmt.Sprintf("无效的目录路径: %v", err))
		return nil, "", err
	}

	// 获取文件类型
	fileTypesStr := cli.askString("要监控的文件类型（以逗号分隔，如 .go,.js）", "")
	var fileTypes []string
	if fileTypesStr != "" {
		for _, t := range strings.Split(fileTypesStr, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				fileTypes = append(fileTypes, t)
			}
		}
	}

	// 获取排除路径
	excludePathsStr := cli.askString("要排除的路径（以逗号分隔）", "")
	var excludePaths []string
	if excludePathsStr != "" {
		for _, p := range strings.Split(excludePathsStr, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				excludePaths = append(excludePaths, p)
			}
		}
	}

	// 获取执行命令
	command := cli.askString("文件变化时执行的命令", "echo 文件已更新")

	// 创建配置
	config, err := entity.NewWatchConfig(absWatchDir, fileTypes, excludePaths, command)
	if err != nil {
		ui.PrintError(fmt.Sprintf("创建配置失败: %v", err))
		return nil, "", err
	}

	// 显示配置摘要
	fmt.Println("\n----------------------------------------")
	ui.PrintHeader("配置摘要:")
	fmt.Printf("配置文件: %s\n", configPath)
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

	return config, configPath, nil
}

// askString 询问字符串输入
func (cli *InteractiveCLI) askString(question, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s%s%s [%s%s%s]: ", ui.Blue, question, ui.Reset, ui.Green, defaultValue, ui.Reset)
	} else {
		fmt.Printf("%s%s%s: ", ui.Blue, question, ui.Reset)
	}

	input, err := cli.reader.ReadString('\n')
	if err != nil {
		// 如果读取失败，返回默认值
		return defaultValue
	}

	// 去除换行符
	input = strings.TrimRight(input, "\r\n")

	// 如果输入为空，返回默认值
	if input == "" {
		return defaultValue
	}

	return input
}

// AskYesNo 询问是/否问题
func (cli *InteractiveCLI) AskYesNo(question string, defaultYes bool) bool {
	defaultValue := "Y/n"
	if !defaultYes {
		defaultValue = "y/N"
	}

	fmt.Printf("%s%s%s [%s]: ", ui.Blue, question, ui.Reset, defaultValue)

	input, err := cli.reader.ReadString('\n')
	if err != nil {
		// 如果读取失败，返回默认值
		return defaultYes
	}

	// 去除换行符并转换为小写
	input = strings.ToLower(strings.TrimRight(input, "\r\n"))

	// 如果输入为空，返回默认值
	if input == "" {
		return defaultYes
	}

	// 检查是否为肯定回答
	return input == "y" || input == "yes"
}
