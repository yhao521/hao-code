package backend

import (
	"github.com/go-delve/delve/service/api"
)

// StartDebug 启动调试
func (a *AppService) StartDebug(program string, args []string) error {
	return a.debug.StartDebugging(program, args)
}

// StopDebug 停止调试
func (a *AppService) StopDebug() error {
	return a.debug.StopDebugging()
}

// SetBreakpoint 设置断点
func (a *AppService) SetBreakpoint(file string, line int) (*api.Breakpoint, error) {
	return a.debug.SetBreakpoint(file, line)
}

// ClearBreakpoint 清除断点
func (a *AppService) ClearBreakpoint(id int) error {
	return a.debug.ClearBreakpoint(id)
}

// DebugContinue 继续执行
func (a *AppService) DebugContinue() (*api.DebuggerState, error) {
	return a.debug.Continue()
}

// DebugNext 单步跳过
func (a *AppService) DebugNext() (*api.DebuggerState, error) {
	return a.debug.Next()
}

// DebugStep 单步进入
func (a *AppService) DebugStep() (*api.DebuggerState, error) {
	return a.debug.Step()
}

// GetDebugVariables 获取变量
func (a *AppService) GetDebugVariables(goroutineID int, frame int) ([]api.Variable, error) {
	return a.debug.GetVariables(api.EvalScope{GoroutineID: int64(goroutineID), Frame: frame})
}

// GetDebugStacktrace 获取调用栈
func (a *AppService) GetDebugStacktrace(depth int) ([]api.Stackframe, error) {
	return a.debug.GetStacktrace(depth)
}
