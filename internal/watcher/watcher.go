package watcher

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/watchs/config"
)

// FileWatcher 文件监控器
type FileWatcher struct {
	config      *config.Config
	watcher     *fsnotify.Watcher
	lastEventAt time.Time
	mu          sync.Mutex
	cmd         *exec.Cmd
}

// NewFileWatcher 创建新的文件监控器
func NewFileWatcher(cfg *config.Config) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("创建文件监控器失败: %w", err)
	}

	return &FileWatcher{
		config:      cfg,
		watcher:     watcher,
		lastEventAt: time.Now(),
	}, nil
}

// Start 开始监控文件
func (fw *FileWatcher) Start() error {
	// 添加初始监控目录
	if err := fw.addWatchDir(fw.config.WatchDir); err != nil {
		return err
	}

	log.Printf("开始监控目录: %s", fw.config.WatchDir)
	if len(fw.config.FileTypes) > 0 {
		log.Printf("监控的文件类型: %s", strings.Join(fw.config.FileTypes, ", "))
	} else {
		log.Printf("监控所有文件类型")
	}
	if len(fw.config.ExcludePaths) > 0 {
		log.Printf("排除的路径: %s", strings.Join(fw.config.ExcludePaths, ", "))
	}
	log.Printf("执行命令: %s", fw.config.Command)

	// 执行初始命令
	if err := fw.executeCommand(); err != nil {
		log.Printf("执行初始命令失败: %v", err)
	}

	// 监听事件
	go fw.watchEvents()

	return nil
}

// Close 关闭监控器
func (fw *FileWatcher) Close() {
	fw.watcher.Close()
	fw.killCommand()
}

// 添加监控目录（递归）
func (fw *FileWatcher) addWatchDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 如果是目录，检查是否应该排除
		if info.IsDir() {
			for _, excludePath := range fw.config.ExcludePaths {
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
			if err := fw.watcher.Add(path); err != nil {
				log.Printf("添加监控目录失败 %s: %v", path, err)
			}
		}
		return nil
	})
}

// 监听文件变化事件
func (fw *FileWatcher) watchEvents() {
	debounceTime := 500 * time.Millisecond

	for {
		select {
		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}

			// 检查是否是创建目录事件
			if event.Has(fsnotify.Create) {
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					// 如果是新创建的目录，添加到监控
					fw.addWatchDir(event.Name)
				}
			}

			// 检查是否应该监控此文件
			if !fw.config.ShouldWatch(event.Name) {
				continue
			}

			// 只处理写入、创建和删除事件
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
				fw.mu.Lock()
				now := time.Now()
				// 防抖：如果两次事件间隔小于设定时间，则忽略
				if now.Sub(fw.lastEventAt) > debounceTime {
					fw.lastEventAt = now
					fw.mu.Unlock()

					log.Printf("检测到文件变化: %s", event.Name)
					// 延迟执行，避免文件正在写入
					time.Sleep(100 * time.Millisecond)
					if err := fw.executeCommand(); err != nil {
						log.Printf("执行命令失败: %v", err)
					}
				} else {
					fw.lastEventAt = now
					fw.mu.Unlock()
				}
			}

		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("监控错误: %v", err)
		}
	}
}

// 执行配置的命令
func (fw *FileWatcher) executeCommand() error {
	fw.killCommand() // 先终止之前的命令

	log.Printf("执行命令: %s", fw.config.Command)

	// 在 Windows 上使用 cmd /c，在 Unix 上使用 /bin/sh -c
	var cmd *exec.Cmd
	if os.PathSeparator == '\\' { // Windows
		cmd = exec.Command("cmd", "/c", fw.config.Command)
	} else { // Unix
		cmd = exec.Command("/bin/sh", "-c", fw.config.Command)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = fw.config.WatchDir

	fw.cmd = cmd
	return cmd.Start()
}

// 终止当前正在执行的命令
func (fw *FileWatcher) killCommand() {
	if fw.cmd != nil && fw.cmd.Process != nil {
		// 在 Windows 上使用 taskkill 来终止进程树
		if os.PathSeparator == '\\' { // Windows
			exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", fw.cmd.Process.Pid)).Run()
		} else { // Unix
			fw.cmd.Process.Kill()
		}
		fw.cmd.Wait() // 等待进程结束
		fw.cmd = nil
	}
}
