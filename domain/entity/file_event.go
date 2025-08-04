package entity

import "time"

// EventType 表示文件事件类型
type EventType int

const (
	// EventCreate 文件创建事件
	EventCreate EventType = iota
	// EventWrite 文件写入事件
	EventWrite
	// EventRemove 文件删除事件
	EventRemove
	// EventRename 文件重命名事件
	EventRename
	// EventChmod 文件权限变更事件
	EventChmod
)

// FileEvent 表示文件变更事件
type FileEvent struct {
	// 文件路径
	Path string
	// 事件类型
	Type EventType
	// 事件发生时间
	Timestamp time.Time
}

// NewFileEvent 创建一个新的文件事件
func NewFileEvent(path string, eventType EventType) *FileEvent {
	return &FileEvent{
		Path:      path,
		Type:      eventType,
		Timestamp: time.Now(),
	}
}
