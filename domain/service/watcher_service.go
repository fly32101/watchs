package service

import "github.com/watchs/domain/entity"

// WatcherService 定义文件监控服务的接口
type WatcherService interface {
	// Start 开始监控文件
	Start() error
	// Stop 停止监控
	Stop() error
	// OnFileEvent 注册文件事件处理函数
	OnFileEvent(handler func(event *entity.FileEvent) error)
}

// CommandExecutor 定义命令执行服务的接口
type CommandExecutor interface {
	// Execute 执行命令
	Execute(command string, workDir string) error
	// Terminate 终止正在执行的命令
	Terminate() error
}
