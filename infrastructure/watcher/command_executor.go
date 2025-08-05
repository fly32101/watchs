package watcher

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/watchs/infrastructure/ui"
)

// CommandExecutorImpl 是命令执行器的实现
type CommandExecutorImpl struct {
	cmd         *exec.Cmd
	mu          sync.Mutex
	lastRunTime time.Time
	debounceMs  int
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewCommandExecutor 创建一个新的命令执行器
func NewCommandExecutor(debounceMs int) *CommandExecutorImpl {
	if debounceMs <= 0 {
		debounceMs = 500 // 默认防抖时间为500毫秒
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &CommandExecutorImpl{
		debounceMs: debounceMs,
		ctx:        ctx,
		cancel:     cancel,
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

	// 根据操作系统选择不同的命令执行方式，使用context进行管理
	var cmd *exec.Cmd
	if os.PathSeparator == '\\' { // Windows
		cmd = exec.CommandContext(e.ctx, "cmd", "/c", command)
	} else { // Unix
		cmd = exec.CommandContext(e.ctx, "/bin/sh", "-c", command)
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

	var err error

	// 在 Windows 上使用 taskkill 来终止进程树
	if os.PathSeparator == '\\' { // Windows
		killCmd := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", e.cmd.Process.Pid))
		killCmd.Run() // 忽略taskkill的错误，因为进程可能已经结束
	} else { // Unix
		err = e.cmd.Process.Kill()
	}

	// 等待进程结束，避免僵尸进程
	waitErr := e.cmd.Wait()
	if err == nil {
		err = waitErr
	}

	e.cmd = nil
	return err
}

// Close 清理资源
func (e *CommandExecutorImpl) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 取消context，这会终止所有使用该context的命令
	if e.cancel != nil {
		e.cancel()
	}

	// 终止当前命令
	return e.terminateUnsafe()
}
