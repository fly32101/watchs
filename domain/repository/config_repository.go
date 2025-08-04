package repository

import "github.com/watchs/domain/entity"

// ConfigRepository 定义配置仓储接口
type ConfigRepository interface {
	// LoadConfig 从指定路径加载配置
	LoadConfig(path string) (*entity.WatchConfig, error)
	// SaveConfig 保存配置到指定路径
	SaveConfig(config *entity.WatchConfig, path string) error
}
