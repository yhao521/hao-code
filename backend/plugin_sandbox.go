package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
)

// PluginSandbox 代表一个隔离的插件运行环境
type PluginSandbox struct {
	Name      string
	Path      string
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	mu        sync.Mutex
	responses map[string]chan interface{}
}

// SandboxRequest 发送给沙箱的请求
type SandboxRequest struct {
	ID      string      `json:"id"`
	Method  string      `json:"method"`
	Payload interface{} `json:"payload"`
}

// SandboxResponse 沙箱返回的响应
type SandboxResponse struct {
	ID     string      `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// NewPluginSandbox 创建一个新的插件沙箱实例
func NewPluginSandbox(name, path string) *PluginSandbox {
	return &PluginSandbox{
		Name:      name,
		Path:      path,
		responses: make(map[string]chan interface{}),
	}
}

// Start 启动沙箱进程（这里以 Node.js 为例）
func (ps *PluginSandbox) Start() error {
	ps.cmd = exec.Command("node", ps.Path+"/main.js") // 假设插件入口是 main.js

	var err error
	ps.stdin, err = ps.cmd.StdinPipe()
	if err != nil {
		return err
	}

	ps.stdout, err = ps.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	go ps.listenResponses()

	if err := ps.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start sandbox for %s: %v", ps.Name, err)
	}

	log.Printf("Sandbox started for plugin: %s", ps.Name)
	return nil
}

// listenResponses 监听来自子进程的 JSON-RPC 响应
func (ps *PluginSandbox) listenResponses() {
	decoder := json.NewDecoder(ps.stdout)
	for {
		var resp SandboxResponse
		if err := decoder.Decode(&resp); err != nil {
			log.Printf("Sandbox decode error for %s: %v", ps.Name, err)
			return
		}

		ps.mu.Lock()
		if ch, ok := ps.responses[resp.ID]; ok {
			ch <- resp.Result
			close(ch)
			delete(ps.responses, resp.ID)
		}
		ps.mu.Unlock()
	}
}

// Execute 在沙箱中执行方法
func (ps *PluginSandbox) Execute(method string, payload interface{}) (interface{}, error) {
	reqID := fmt.Sprintf("%s-%d", ps.Name, len(ps.responses))
	req := SandboxRequest{
		ID:      reqID,
		Method:  method,
		Payload: payload,
	}

	respChan := make(chan interface{}, 1)
	ps.mu.Lock()
	ps.responses[reqID] = respChan
	ps.mu.Unlock()

	data, _ := json.Marshal(req)
	if _, err := ps.stdin.Write(append(data, '\n')); err != nil {
		return nil, err
	}

	result := <-respChan
	return result, nil
}

// Stop 停止沙箱进程
func (ps *PluginSandbox) Stop() error {
	if ps.cmd != nil && ps.cmd.Process != nil {
		return ps.cmd.Process.Kill()
	}
	return nil
}
