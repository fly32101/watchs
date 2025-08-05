package watcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/watchs/domain/entity"
	"github.com/watchs/infrastructure/ui"
)

// FSNotifyWatcher 是基于fsnotify的文件监控服务实现
type FSNotifyWatcher struct {
	config        *entity.WatchConfig
	watcher       *fsnotify.Watcher
	eventHandlers []func(event *entity.FileEvent) error
	mu            sync.RWMutex
	isRunning     bool
}

// NewFSNotifyWatcher 创建一个新的fsnotify文件监控器
func NewFSNotifyWatcher(config *entity.WatchConfig) (*FSNotifyWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("创建文件监控器失败: %w", err)
	}

	return &FSNotifyWatcher{
		config:    config,
		watcher:   watcher,
		isRunning: false,
	}, nil
}

// Start 开始监控文件
func (w *FSNotifyWatcher) Start() error {
	w.mu.Lock()
	if w.isRunning {
		w.mu.Unlock()
		return nil
	}
	w.isRunning = true
	w.mu.Unlock()

	// 添加初始监控目录
	if err := w.addWatchDir(w.config.WatchDir); err != nil {
		return err
	}

	ui.PrintSuccess(fmt.Sprintf("开始监控目录: %s", w.config.WatchDir))
	if len(w.config.FileTypes) > 0 {
		ui.PrintInfo(fmt.Sprintf("监控的文件类型: %v", w.config.FileTypes))
	} else {
		ui.PrintInfo("监控所有文件类型")
	}
	if len(w.config.ExcludePaths) > 0 {
		ui.PrintInfo(fmt.Sprintf("排除的路径: %v", w.config.ExcludePaths))
	}

	// 监听事件
	go w.watchEvents()

	return nil
}

// Stop 停止监控
func (w *FSNotifyWatcher) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.isRunning {
		return nil
	}

	w.isRunning = false
	return w.watcher.Close()
}

// OnFileEvent 注册文件事件处理函数
func (w *FSNotifyWatcher) OnFileEvent(handler func(event *entity.FileEvent) error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.eventHandlers = append(w.eventHandlers, handler)
}

// 添加监控目录（递归）
func (w *FSNotifyWatcher) addWatchDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 如果是目录，检查是否应该排除
		if info.IsDir() {
			for _, excludePath := range w.config.ExcludePaths {
				absExcludePath, err := filepath.Abs(excludePath)
				if err == nil && (path == absExcludePath || filepath.HasPrefix(path, absExcludePath+string(os.PathSeparator))) {
					return filepath.SkipDir
				}

				// 支持通配符匹配目录名
				matched, err := filepath.Match(excludePath, filepath.Base(path))
				if err == nil && matched {
					return filepath.SkipDir
				}
			}

			// 添加目录到监控
			if err := w.watcher.Add(path); err != nil {
				log.Printf("添加监控目录失败 %s: %v", path, err)
			}
		}
		return nil
	})
}

// 监听文件变化事件
func (w *FSNotifyWatcher) watchEvents() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			// 检查是否是创建目录事件
			if event.Has(fsnotify.Create) {
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					// 如果是新创建的目录，添加到监控
					w.addWatchDir(event.Name)
				}
			}

			// 检查是否应该监控此文件
			if !w.config.ShouldWatch(event.Name) {
				continue
			}

			// 转换为领域事件
			var eventType entity.EventType
			if event.Has(fsnotify.Create) {
				eventType = entity.EventCreate
			} else if event.Has(fsnotify.Write) {
				eventType = entity.EventWrite
			} else if event.Has(fsnotify.Remove) {
				eventType = entity.EventRemove
			} else if event.Has(fsnotify.Rename) {
				eventType = entity.EventRename
			} else if event.Has(fsnotify.Chmod) {
				eventType = entity.EventChmod
			} else {
				continue
			}

			fileEvent := entity.NewFileEvent(event.Name, eventType)
			ui.PrintEvent(fileEvent)

			// 通知所有处理器
			w.mu.RLock()
			handlers := w.eventHandlers
			w.mu.RUnlock()

			for _, handler := range handlers {
				if err := handler(fileEvent); err != nil {
					log.Printf("处理文件事件失败: %v", err)
				}
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			ui.PrintError(fmt.Sprintf("监控错误: %v", err))
		}
	}
}