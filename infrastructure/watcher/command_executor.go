package watcher

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/watchs/presentation/cli/ui"
)

// CommandExecutorImpl 是命令执行器的实现
type CommandExecutorImpl struct {
	cmd         *exec.Cmd
	mu          sync.Mutex
	lastRunTime time.Time
	debounceMs  int
}

// NewCommandExecutor 创建一个新的命令执行器
func NewCommandExecutor(debounceMs int) *CommandExecutorImpl {
	if debounceMs <= 0 {
		debounceMs = 500 // 默认防抖时间为500毫秒
	}

	return &CommandExecutorImpl{
		debounceMs: debounceMs,
	}
}

// Execute 执行命令
func (e *CommandExecutorImpl) Execute(command string, workDir string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 防抖：如果两次执行间隔小于设定时间，则忽略
	now := time.Now()
	if now.Sub(e.lastRunTime) < time.Duration(e.debounceMs)*time.Millisecond {
		return nil
	}
	e.lastRunTime = now

	// 先终止之前的命令
	e.terminateUnsafe()

	ui.PrintInfo(fmt.Sprintf("执行命令: %s", command))

	// 根据操作系统选择不同的命令执行方式
	var cmd *exec.Cmd
	if os.PathSeparator == '\\' { // Windows
		cmd = exec.Command("cmd", "/c", command)
	} else { // Unix
		cmd = exec.Command("/bin/sh", "-c", command)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = workDir

	e.cmd = cmd
	return cmd.Start()
}

// Terminate 终止正在执行的命令
func (e *CommandExecutorImpl) Terminate() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.terminateUnsafe()
}

// terminateUnsafe 在不加锁的情况下终止命令（内部使用）
func (e *CommandExecutorImpl) terminateUnsafe() error {
	if e.cmd == nil || e.cmd.Process == nil {
		return nil
	}

	// 在 Windows 上使用 taskkill 来终止进程树
	if os.PathSeparator == '\\' { // Windows
		exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", e.cmd.Process.Pid)).Run()
	} else { // Unix
		e.cmd.Process.Kill()
	}

	e.cmd.Wait() // 等待进程结束
	e.cmd = nil
	return nil
}