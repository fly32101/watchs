package cli

import (
	"strings"

	"github.com/watchs/domain/entity"
)

// createConfigFromArgs 从命令行参数创建配置
func createConfigFromArgs(watchDir, fileTypes, excludePaths, command string) (*entity.WatchConfig, error) {
	return entity.NewWatchConfig(
		watchDir,
		parseCommaSeparated(fileTypes),
		parseCommaSeparated(excludePaths),
		command,
	)
}

// overrideConfig 使用命令行参数覆盖配置
func overrideConfig(config *entity.WatchConfig, watchDir, fileTypes, excludePaths, command string) (*entity.WatchConfig, error) {
	// 如果命令行参数不为空，则覆盖配置
	dir := config.WatchDir
	if watchDir != "" {
		dir = watchDir
	}

	types := config.FileTypes
	if fileTypes != "" {
		types = parseCommaSeparated(fileTypes)
	}

	exclude := config.ExcludePaths
	if excludePaths != "" {
		exclude = parseCommaSeparated(excludePaths)
	}

	cmd := config.Command
	if command != "" {
		cmd = command
	}

	return entity.NewWatchConfig(dir, types, exclude, cmd)
}

// parseCommaSeparated 解析逗号分隔的字符串
func parseCommaSeparated(s string) []string {
	if s == "" {
		return nil
	}

	var result []string
	for _, item := range strings.Split(s, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}

	return result
}
