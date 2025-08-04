package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 表示监控配置
type Config struct {
	// 要监控的目录
	WatchDir string `json:"watch_dir"`
	// 要监控的文件类型，如 [".go", ".js"]，为空则监控所有文件
	FileTypes []string `json:"file_types"`
	// 要排除的目录或文件
	ExcludePaths []string `json:"exclude_paths"`
	// 文件变化时要执行的命令
	Command string `json:"command"`
}

// LoadConfig 从指定路径加载配置文件
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if config.WatchDir == "" {
		return nil, fmt.Errorf("监控目录不能为空")
	}

	// 将相对路径转换为绝对路径
	absPath, err := filepath.Abs(config.WatchDir)
	if err != nil {
		return nil, fmt.Errorf("获取绝对路径失败: %w", err)
	}
	config.WatchDir = absPath

	// 检查目录是否存在
	if _, err := os.Stat(config.WatchDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("监控目录不存在: %s", config.WatchDir)
	}

	if config.Command == "" {
		return nil, fmt.Errorf("执行命令不能为空")
	}

	return &config, nil
}

// ShouldWatch 判断给定文件是否应该被监控
func (c *Config) ShouldWatch(path string) bool {
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
