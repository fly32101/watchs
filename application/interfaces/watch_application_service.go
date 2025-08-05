package interfaces

import "github.com/watchs/domain/entity"

// WatchApplicationService 定义监控应用服务接口
type WatchApplicationService interface {
	// StartWatch 启动文件监控
	StartWatch(config *WatchConfig) error
	// StopWatch 停止文件监控
	StopWatch() error
	// CreateWatchConfigFromArgs 从命令行参数创建监控配置
	CreateWatchConfigFromArgs(watchDir, fileTypes, excludePaths, command string) (*entity.WatchConfig, error)
	// IsRunning 检查监控是否正在运行
	IsRunning() bool
}

// WatchConfig 监控配置参数
type WatchConfig struct {
	ConfigPath     string
	WatchDir       string
	FileTypes      string
	ExcludePaths   string
	Command        string
	DebounceMs     int
	ShowMemory     bool
	MemoryInterval int
}
