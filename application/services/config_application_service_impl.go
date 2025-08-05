package services

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/watchs/application/interfaces"
	"github.com/watchs/domain/entity"
	"github.com/watchs/domain/repository"
	"github.com/watchs/infrastructure/ui"
	"github.com/watchs/presentation/cli"
)

// ConfigApplicationServiceImpl 配置应用服务实现
type ConfigApplicationServiceImpl struct {
	configRepo repository.ConfigRepository
}

// NewConfigApplicationService 创建配置应用服务
func NewConfigApplicationService(configRepo repository.ConfigRepository) interfaces.ConfigApplicationService {
	return &ConfigApplicationServiceImpl{
		configRepo: configRepo,
	}
}

// LoadOrCreateConfig 加载或创建配置
func (s *ConfigApplicationServiceImpl) LoadOrCreateConfig(configPath, watchDir, fileTypes, excludePaths, command string) (*entity.WatchConfig, error) {
	// 尝试加载配置
	config, err := s.configRepo.LoadConfig(configPath)
	if err != nil {
		// 如果配置文件不存在，尝试使用命令行参数
		if os.IsNotExist(err) && watchDir != "" && command != "" {
			log.Printf("配置文件 %s 不存在，使用命令行参数", configPath)
			return s.createConfigFromArgs(watchDir, fileTypes, excludePaths, command)
		}
		return nil, err
	}

	// 命令行参数覆盖配置文件
	if watchDir != "" || fileTypes != "" || excludePaths != "" || command != "" {
		return s.overrideConfig(config, watchDir, fileTypes, excludePaths, command)
	}

	return config, nil
}

// SaveConfig 保存配置
func (s *ConfigApplicationServiceImpl) SaveConfig(config *entity.WatchConfig, configPath string) error {
	return s.configRepo.SaveConfig(config, configPath)
}

// InitializeConfig 初始化配置文件
func (s *ConfigApplicationServiceImpl) InitializeConfig(params *interfaces.InitConfigParams) error {
	// 检查配置文件是否已存在
	if _, err := os.Stat(params.ConfigPath); err == nil && !params.Force {
		ui.PrintError(fmt.Sprintf("配置文件 %s 已存在，使用 --force 参数强制覆盖", params.ConfigPath))
		return fmt.Errorf("配置文件已存在")
	}

	// 创建配置
	config, err := s.createConfigFromArgs(params.WatchDir, params.FileTypes, params.ExcludePaths, params.Command)
	if err != nil {
		ui.PrintError(fmt.Sprintf("创建配置失败: %v", err))
		return fmt.Errorf("创建配置失败: %v", err)
	}

	// 保存配置
	if err := s.SaveConfig(config, params.ConfigPath); err != nil {
		ui.PrintError(fmt.Sprintf("保存配置文件失败: %v", err))
		return fmt.Errorf("保存配置文件失败: %v", err)
	}

	ui.PrintSuccess(fmt.Sprintf("配置文件已生成: %s", params.ConfigPath))
	ui.PrintInfo("配置内容:")
	fmt.Printf("  监控目录: %s\n", config.WatchDir)
	if len(config.FileTypes) > 0 {
		fmt.Printf("  监控的文件类型: %v\n", config.FileTypes)
	} else {
		fmt.Printf("  监控所有文件类型\n")
	}
	if len(config.ExcludePaths) > 0 {
		fmt.Printf("  排除的路径: %v\n", config.ExcludePaths)
	}
	fmt.Printf("  执行命令: %s\n", config.Command)

	return nil
}

// RunInteractiveConfig 运行交互式配置向导
func (s *ConfigApplicationServiceImpl) RunInteractiveConfig() (*entity.WatchConfig, string, error) {
	// 创建交互式CLI
	interactiveCLI := cli.NewInteractiveCLI()

	// 运行交互式配置
	return interactiveCLI.Run()
}

// createConfigFromArgs 从命令行参数创建配置
func (s *ConfigApplicationServiceImpl) createConfigFromArgs(watchDir, fileTypes, excludePaths, command string) (*entity.WatchConfig, error) {
	return entity.NewWatchConfig(
		watchDir,
		s.parseCommaSeparated(fileTypes),
		s.parseCommaSeparated(excludePaths),
		command,
	)
}

// overrideConfig 用命令行参数覆盖配置
func (s *ConfigApplicationServiceImpl) overrideConfig(config *entity.WatchConfig, watchDir, fileTypes, excludePaths, command string) (*entity.WatchConfig, error) {
	newWatchDir := config.WatchDir
	if watchDir != "" {
		newWatchDir = watchDir
	}

	newFileTypes := config.FileTypes
	if fileTypes != "" {
		newFileTypes = s.parseCommaSeparated(fileTypes)
	}

	newExcludePaths := config.ExcludePaths
	if excludePaths != "" {
		newExcludePaths = s.parseCommaSeparated(excludePaths)
	}

	newCommand := config.Command
	if command != "" {
		newCommand = command
	}

	return entity.NewWatchConfig(newWatchDir, newFileTypes, newExcludePaths, newCommand)
}

// parseCommaSeparated 解析逗号分隔的字符串
func (s *ConfigApplicationServiceImpl) parseCommaSeparated(str string) []string {
	if str == "" {
		return nil
	}

	var result []string
	for _, item := range strings.Split(str, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}
