package persistence

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/watchs/domain/entity"
)

// configDTO 是配置的数据传输对象
type configDTO struct {
	WatchDir     string   `json:"watch_dir"`
	FileTypes    []string `json:"file_types"`
	ExcludePaths []string `json:"exclude_paths"`
	Command      string   `json:"command"`
}

// JsonConfigRepository 是基于JSON文件的配置仓储实现
type JsonConfigRepository struct{}

// NewJsonConfigRepository 创建一个新的JSON配置仓储
func NewJsonConfigRepository() *JsonConfigRepository {
	return &JsonConfigRepository{}
}

// LoadConfig 从JSON文件加载配置
func (r *JsonConfigRepository) LoadConfig(path string) (*entity.WatchConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var dto configDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 转换为领域实体
	config, err := entity.NewWatchConfig(dto.WatchDir, dto.FileTypes, dto.ExcludePaths, dto.Command)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// SaveConfig 保存配置到JSON文件
func (r *JsonConfigRepository) SaveConfig(config *entity.WatchConfig, path string) error {
	// 转换为DTO
	dto := configDTO{
		WatchDir:     config.WatchDir,
		FileTypes:    config.FileTypes,
		ExcludePaths: config.ExcludePaths,
		Command:      config.Command,
	}

	data, err := json.MarshalIndent(dto, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}
