package interfaces

import "github.com/watchs/domain/entity"

// ConfigApplicationService 定义配置应用服务接口
type ConfigApplicationService interface {
	// LoadOrCreateConfig 加载或创建配置
	LoadOrCreateConfig(configPath, watchDir, fileTypes, excludePaths, command string) (*entity.WatchConfig, error)
	// SaveConfig 保存配置
	SaveConfig(config *entity.WatchConfig, configPath string) error
	// InitializeConfig 初始化配置文件
	InitializeConfig(params *InitConfigParams) error
	// RunInteractiveConfig 运行交互式配置向导
	RunInteractiveConfig() (*entity.WatchConfig, string, error)
}

// InitConfigParams 初始化配置参数
type InitConfigParams struct {
	ConfigPath   string
	WatchDir     string
	FileTypes    string
	ExcludePaths string
	Command      string
	Force        bool
}
