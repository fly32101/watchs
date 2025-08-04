package application

import (
	"log"
	"time"

	"github.com/watchs/domain/entity"
	"github.com/watchs/domain/service"
)

// WatchService 是应用层的文件监控服务
type WatchService struct {
	config          *entity.WatchConfig
	watcherService  service.WatcherService
	commandExecutor service.CommandExecutor
	isRunning       bool
}

// NewWatchService 创建一个新的应用层文件监控服务
func NewWatchService(
	config *entity.WatchConfig,
	watcherService service.WatcherService,
	commandExecutor service.CommandExecutor,
) *WatchService {
	return &WatchService{
		config:          config,
		watcherService:  watcherService,
		commandExecutor: commandExecutor,
		isRunning:       false,
	}
}

// Start 开始监控文件
func (s *WatchService) Start() error {
	if s.isRunning {
		return nil
	}
	s.isRunning = true

	// 注册文件事件处理器
	s.watcherService.OnFileEvent(func(event *entity.FileEvent) error {
		// 只处理写入、创建和删除事件
		if event.Type == entity.EventWrite || event.Type == entity.EventCreate || event.Type == entity.EventRemove {
			// 延迟执行，避免文件正在写入
			time.Sleep(100 * time.Millisecond)
			return s.commandExecutor.Execute(s.config.Command, s.config.WatchDir)
		}
		return nil
	})

	// 启动监控服务
	if err := s.watcherService.Start(); err != nil {
		s.isRunning = false
		return err
	}

	// 执行初始命令
	if err := s.commandExecutor.Execute(s.config.Command, s.config.WatchDir); err != nil {
		log.Printf("执行初始命令失败: %v", err)
	}

	return nil
}

// Stop 停止监控
func (s *WatchService) Stop() error {
	if !s.isRunning {
		return nil
	}
	s.isRunning = false

	// 终止命令
	if err := s.commandExecutor.Terminate(); err != nil {
		log.Printf("终止命令失败: %v", err)
	}

	// 停止监控服务
	return s.watcherService.Stop()
}
