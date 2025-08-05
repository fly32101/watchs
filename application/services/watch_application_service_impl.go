package services

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/watchs/application"
	"github.com/watchs/application/interfaces"
	"github.com/watchs/domain/entity"
	"github.com/watchs/infrastructure/ui"
	"github.com/watchs/infrastructure/watcher"
)

// WatchApplicationServiceImpl 监控应用服务实现
type WatchApplicationServiceImpl struct {
	configService interfaces.ConfigApplicationService
	watchService  *application.WatchService
	isRunning     bool
}

// NewWatchApplicationService 创建监控应用服务
func NewWatchApplicationService(configService interfaces.ConfigApplicationService) interfaces.WatchApplicationService {
	return &WatchApplicationServiceImpl{
		configService: configService,
		isRunning:     false,
	}
}

// StartWatch 启动文件监控
func (s *WatchApplicationServiceImpl) StartWatch(params *interfaces.WatchConfig) error {
	if s.isRunning {
		return fmt.Errorf("监控已在运行中")
	}

	// 加载或创建配置
	config, err := s.configService.LoadOrCreateConfig(
		params.ConfigPath,
		params.WatchDir,
		params.FileTypes,
		params.ExcludePaths,
		params.Command,
	)
	if err != nil {
		ui.PrintError(fmt.Sprintf("配置加载失败: %v", err))
		return fmt.Errorf("配置加载失败: %v", err)
	}

	ui.PrintInfo("正在初始化监控服务...")

	// 模拟加载动画
	ui.SimulateLoading(2*time.Second, "初始化监控器")

	// 创建文件监控服务
	fsWatcher, err := watcher.NewFSNotifyWatcher(config)
	if err != nil {
		ui.PrintError(fmt.Sprintf("创建文件监控器失败: %v", err))
		return fmt.Errorf("创建文件监控器失败: %v", err)
	}

	// 创建命令执行器
	cmdExecutor := watcher.NewCommandExecutor(params.DebounceMs)

	// 创建应用服务
	s.watchService = application.NewWatchService(config, fsWatcher, cmdExecutor)

	// 启动监控
	if err := s.watchService.Start(); err != nil {
		ui.PrintError(fmt.Sprintf("启动监控失败: %v", err))
		return fmt.Errorf("启动监控失败: %v", err)
	}

	s.isRunning = true
	ui.PrintSuccess(fmt.Sprintf("监控已启动，正在监控目录: %s", config.WatchDir))
	ui.PrintInfo("按 Ctrl+C 停止监控...")

	// 等待中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	return s.StopWatch()
}

// StopWatch 停止文件监控
func (s *WatchApplicationServiceImpl) StopWatch() error {
	if !s.isRunning {
		return nil
	}

	ui.PrintWarning("正在关闭监控...")
	if err := s.watchService.Stop(); err != nil {
		ui.PrintError(fmt.Sprintf("关闭监控失败: %v", err))
		return err
	}

	s.isRunning = false
	ui.PrintSuccess("监控已成功关闭!")
	return nil
}

// CreateWatchConfigFromArgs 从命令行参数创建监控配置
func (s *WatchApplicationServiceImpl) CreateWatchConfigFromArgs(watchDir, fileTypes, excludePaths, command string) (*entity.WatchConfig, error) {
	return s.configService.LoadOrCreateConfig("", watchDir, fileTypes, excludePaths, command)
}

// IsRunning 检查监控是否正在运行
func (s *WatchApplicationServiceImpl) IsRunning() bool {
	return s.isRunning
}
