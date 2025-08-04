package entity

import (
	"fmt"
	"os"
	"path/filepath"
)

// WatchConfig 表示文件监控的配置实体
type WatchConfig struct {
	// 要监控的目录
	WatchDir string
	// 要监控的文件类型，如 [".go", ".js"]，为空则监控所有文件
	FileTypes []string
	// 要排除的目录或文件
	ExcludePaths []string
	// 文件变化时要执行的命令
	Command string
}

// NewWatchConfig 创建一个新的监控配置
func NewWatchConfig(watchDir string, fileTypes []string, excludePaths []string, command string) (*WatchConfig, error) {
	if watchDir == "" {
		return nil, fmt.Errorf("监控目录不能为空")
	}

	// 将相对路径转换为绝对路径
	absPath, err := filepath.Abs(watchDir)
	if err != nil {
		return nil, fmt.Errorf("获取绝对路径失败: %w", err)
	}

	// 检查目录是否存在
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("监控目录不存在: %s", absPath)
	}

	if command == "" {
		return nil, fmt.Errorf("执行命令不能为空")
	}

	return &WatchConfig{
		WatchDir:     absPath,
		FileTypes:    fileTypes,
		ExcludePaths: excludePaths,
		Command:      command,
	}, nil
}

// ShouldWatch 判断给定文件是否应该被监控
func (c *WatchConfig) ShouldWatch(path string) bool {
	// 检查是否在排除列表中
	for _, excludePath := range c.ExcludePaths {
		absExcludePath, err := filepath.Abs(excludePath)
		if err == nil && (path == absExcludePath || filepath.HasPrefix(path, absExcludePath+string(os.PathSeparator))) {
			return false
		}

		// 支持通配符匹配
		matched, err := filepath.Match(excludePath, filepath.Base(path))
		if err == nil && matched {
			return false
		}
	}

	// 如果没有指定文件类型，监控所有文件
	if len(c.FileTypes) == 0 {
		return true
	}

	// 检查文件类型是否匹配
	ext := filepath.Ext(path)
	for _, fileType := range c.FileTypes {
		if ext == fileType {
			return true
		}
	}

	return false
}
