package backend

import (
	"fmt"
	"log"
	"net"

	"github.com/go-delve/delve/service"
	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/rpc2"
	"github.com/go-delve/delve/service/rpccommon"
)

// DebugService 调试服务
type DebugService struct {
	client      *rpc2.RPCClient
	breakpoints map[int]*api.Breakpoint
}

// NewDebugService 创建调试服务
func NewDebugService() *DebugService {
	return &DebugService{
		breakpoints: make(map[int]*api.Breakpoint),
	}
}

// StartDebugging 启动调试会话
func (ds *DebugService) StartDebugging(program string, args []string) error {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("failed to create listener: %v", err)
	}

	config := &service.Config{
		Listener:    listener,
		ProcessArgs: append([]string{program}, args...),
		APIVersion:  2,
	}

	server := rpccommon.NewServer(config)
	go server.Run()

	// 连接客户端
	ds.client = rpc2.NewClient(listener.Addr().String())

	// 初始化断点映射
	ds.breakpoints = make(map[int]*api.Breakpoint)

	log.Println("Debugging session started for:", program)
	return nil
}

// SetBreakpoint 设置断点
func (ds *DebugService) SetBreakpoint(file string, line int) (*api.Breakpoint, error) {
	if ds.client == nil {
		return nil, fmt.Errorf("debug session not started")
	}

	bp, err := ds.client.CreateBreakpoint(&api.Breakpoint{
		File: file,
		Line: line,
	})
	if err != nil {
		return nil, err
	}

	ds.breakpoints[bp.ID] = bp
	return bp, nil
}

// ClearBreakpoint 清除断点
func (ds *DebugService) ClearBreakpoint(id int) error {
	if ds.client == nil {
		return fmt.Errorf("debug session not started")
	}

	_, err := ds.client.ClearBreakpoint(id)
	if err != nil {
		return err
	}

	delete(ds.breakpoints, id)
	return nil
}

// Continue 继续执行
func (ds *DebugService) Continue() (*api.DebuggerState, error) {
	if ds.client == nil {
		return nil, fmt.Errorf("debug session not started")
	}
	stateChan := ds.client.Continue()
	_ = stateChan // 异步处理，暂不阻塞
	return nil, nil
}

// Next 单步跳过
func (ds *DebugService) Next() (*api.DebuggerState, error) {
	if ds.client == nil {
		return nil, fmt.Errorf("debug session not started")
	}
	return ds.client.Next()
}

// Step 单步进入
func (ds *DebugService) Step() (*api.DebuggerState, error) {
	if ds.client == nil {
		return nil, fmt.Errorf("debug session not started")
	}
	return ds.client.Step()
}

// GetVariables 获取当前作用域变量
func (ds *DebugService) GetVariables(scope api.EvalScope) ([]api.Variable, error) {
	if ds.client == nil {
		return nil, fmt.Errorf("debug session not started")
	}
	vars, err := ds.client.ListLocalVariables(scope, api.LoadConfig{})
	return vars, err
}

// GetStacktrace 获取调用栈
func (ds *DebugService) GetStacktrace(depth int) ([]api.Stackframe, error) {
	if ds.client == nil {
		return nil, fmt.Errorf("debug session not started")
	}
	frames, err := ds.client.Stacktrace(0, depth, 0, &api.LoadConfig{})
	return frames, err
}

// StopDebugging 停止调试
func (ds *DebugService) StopDebugging() error {
	if ds.client != nil {
		ds.client.Detach(true)
		ds.client = nil
	}
	log.Println("Debugging session stopped")
	return nil
}
