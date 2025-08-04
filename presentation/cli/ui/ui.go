package ui

import (
	"fmt"
	"os"
	"time"

	"github.com/watchs/domain/entity"
)

// ANSI颜色码
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

// Emoji图标
const (
	CheckMark  = "✅"
	CrossMark  = "❌"
	Warning    = "⚠️"
	Info       = "ℹ️"
	Sparkles   = "✨"
	Rocket     = "🚀"
	Gear       = "⚙️"
	Lightning  = "⚡"
	WatchEmoji = "⏰"
)

// PrintSuccess 打印成功信息（绿色）
func PrintSuccess(message string) {
	fmt.Fprintf(os.Stderr, "%s%s%s %s\n", Green, CheckMark, Reset, message)
}

// PrintError 打印错误信息（红色）
func PrintError(message string) {
	fmt.Fprintf(os.Stderr, "%s%s%s %s\n", Red, CrossMark, Reset, message)
}

// PrintWarning 打印警告信息（黄色）
func PrintWarning(message string) {
	fmt.Fprintf(os.Stderr, "%s%s%s %s\n", Yellow, Warning, Reset, message)
}

// PrintInfo 打印信息（蓝色）
func PrintInfo(message string) {
	fmt.Fprintf(os.Stderr, "%s%s%s %s\n", Blue, Info, Reset, message)
}

// PrintHeader 打印标题（紫色）
func PrintHeader(message string) {
	fmt.Fprintf(os.Stderr, "%s%s %s %s\n", Purple, Sparkles, message, Reset)
}

// PrintEvent 打印文件事件信息
func PrintEvent(event *entity.FileEvent) {
	var emoji string
	var color string
	
	switch event.Type {
	case entity.EventCreate:
		emoji = "📄"
		color = Green
	case entity.EventWrite:
		emoji = "✏️"
		color = Blue
	case entity.EventRemove:
		emoji = "🗑️"
		color = Red
	case entity.EventRename:
		emoji = "✏️"
		color = Yellow
	case entity.EventChmod:
		emoji = "🔐"
		color = Purple
	default:
		emoji = "❓"
		color = Gray
	}
	
	fmt.Fprintf(os.Stderr, "%s%s%s %s %s%s%s\n", color, emoji, Reset, event.Path, Gray, event.Type, Reset)
}

// PrintProgressBar 显示进度条
func PrintProgressBar(current, total int, label string) {
	const barLength = 20
	progress := current * barLength / total
	percentage := current * 100 / total
	
	fmt.Fprintf(os.Stderr, "\r%s [%s%s] %d%% %s",
		WatchEmoji,
		Green+string([]rune("████████████████████")[:progress])+Reset,
		Gray+string([]rune("                    ")[progress:])+Reset,
		percentage,
		label)
	
	if current == total {
		fmt.Fprintln(os.Stderr)
	}
}

// SimulateLoading 模拟加载动画
func SimulateLoading(duration time.Duration, message string) {
	chars := []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	endTime := time.Now().Add(duration)
	
	i := 0
	for time.Now().Before(endTime) {
		fmt.Fprintf(os.Stderr, "\r%s %s %s", Lightning, chars[i%len(chars)], message)
		time.Sleep(100 * time.Millisecond)
		i++
	}
	fmt.Fprintf(os.Stderr, "\r%s %s Done!%s\n", CheckMark, Green, Reset)
}