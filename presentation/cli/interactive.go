package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/watchs/domain/entity"
)

// InteractiveCLI 交互式命令行
type InteractiveCLI struct {
	reader *bufio.Reader
}

// NewInteractiveCLI 创建新的交互式命令行
func NewInteractiveCLI() *InteractiveCLI {
	return &InteractiveCLI{
		reader: bufio.NewReader(os.Stdin),
	}
}

// Run 运行交互式命令行
func (ic *InteractiveCLI) Run() (*entity.WatchConfig, string, error) {
	fmt.Println("欢迎使用 Watchs 文件监控工具配置向导")
	fmt.Println("请回答以下问题来创建配置文件")
	fmt.Println("----------------------------------------")

	// 获取配置文件路径
	configPath := ic.askQuestion("配置文件保存路径 [watchs.json]: ", "watchs.json")

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); err == nil {
		overwrite := ic.askYesNo(fmt.Sprintf("配置文件 %s 已存在，是否覆盖？[y/N]: ", configPath), false)
		if !overwrite {
			return nil, "", fmt.Errorf("操作已取消")
		}
	}

	// 获取监控目录
	watchDir := ic.askQuestion("要监控的目录 [./]: ", "./")

	// 获取文件类型
	fileTypesStr := ic.askQuestion("要监控的文件类型（多个类型用逗号分隔，如 .go,.js，留空监控所有文件）: ", "")
	var fileTypes []string
	if fileTypesStr != "" {
		fileTypes = parseCommaSeparated(fileTypesStr)
	}

	// 获取排除路径
	excludePathsStr := ic.askQuestion("要排除的路径（多个路径用逗号分隔，如 vendor,node_modules）: ", "")
	var excludePaths []string
	if excludePathsStr != "" {
		excludePaths = parseCommaSeparated(excludePathsStr)
	}

	// 获取执行命令
	command := ic.askQuestion("文件变化时执行的命令: ", "echo 文件已更新")

	// 获取防抖时间
	debounceStr := ic.askQuestion("防抖时间（毫秒）[500]: ", "500")

	// 创建配置
	config, err := entity.NewWatchConfig(watchDir, fileTypes, excludePaths, command)
	if err != nil {
		return nil, "", err
	}

	// 显示配置摘要
	fmt.Println("\n----------------------------------------")
	fmt.Println("配置摘要:")
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
	fmt.Printf("防抖时间: %s 毫秒\n", debounceStr)
	fmt.Println("----------------------------------------")

	// 确认配置
	confirm := ic.askYesNo("确认以上配置？[Y/n]: ", true)
	if !confirm {
		return nil, "", fmt.Errorf("操作已取消")
	}

	return config, configPath, nil
}

// askQuestion 询问问题并获取回答
func (ic *InteractiveCLI) askQuestion(question string, defaultValue string) string {
	fmt.Print(question)
	answer, _ := ic.reader.ReadString('\n')
	answer = strings.TrimSpace(answer)

	if answer == "" {
		return defaultValue
	}
	return answer
}

// askYesNo 询问是/否问题
func (ic *InteractiveCLI) askYesNo(question string, defaultValue bool) bool {
	fmt.Print(question)
	answer, _ := ic.reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer == "" {
		return defaultValue
	}

	return answer == "y" || answer == "yes"
}
