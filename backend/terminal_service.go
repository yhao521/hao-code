package backend

import (
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
)

// TerminalSession 终端会话
type TerminalSession struct {
	Cmd    *exec.Cmd
	Pty    *os.File
	Output io.ReadCloser
	mu     sync.Mutex
}

// NewTerminalSession 创建新的终端会话
func NewTerminalSession() (*TerminalSession, error) {
	// 获取当前用户的 shell
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	cmd := exec.Command(shell)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	// 启动 PTY
	ptyFile, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	return &TerminalSession{
		Cmd: cmd,
		Pty: ptyFile,
	}, nil
}

// Write 向终端写入数据
func (ts *TerminalSession) Write(data []byte) (int, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.Pty.Write(data)
}

// Read 从终端读取数据
func (ts *TerminalSession) Read(p []byte) (int, error) {
	return ts.Pty.Read(p)
}

// Resize 调整终端大小
func (ts *TerminalSession) Resize(cols, rows uint32) error {
	return pty.Setsize(ts.Pty, &pty.Winsize{
		Cols: uint16(cols),
		Rows: uint16(rows),
	})
}

// Close 关闭终端会话
func (ts *TerminalSession) Close() error {
	if ts.Pty != nil {
		ts.Pty.Close()
	}
	if ts.Cmd != nil && ts.Cmd.Process != nil {
		ts.Cmd.Process.Kill()
	}
	return nil
}
