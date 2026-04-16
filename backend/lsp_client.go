package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"

	"github.com/sourcegraph/jsonrpc2"
)

// LSPClient 封装了与语言服务器的通信
type LSPClient struct {
	conn   *jsonrpc2.Conn
	cmd    *exec.Cmd
	mu     sync.Mutex
	closed bool
}

// NewLSPClient 启动并连接到一个 LSP 服务器 (如 gopls)
func NewLSPClient(command string, args ...string) (*LSPClient, error) {
	cmd := exec.Command(command, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start LSP server: %w", err)
	}

	diagChan := make(chan map[string]interface{}, 100)
	handler := &noopHandler{DiagnosticsChan: diagChan}

	stream := jsonrpc2.NewBufferedStream(&pipe{stdin, stdout}, jsonrpc2.VSCodeObjectCodec{})
	conn := jsonrpc2.NewConn(context.Background(), stream, handler)

	client := &LSPClient{
		conn: conn,
		cmd:  cmd,
	}

	// 启动 goroutine 监听诊断信息
	go func() {
		for params := range diagChan {
			log.Printf("Received diagnostics for URI: %v", params["uri"])
			// 这里可以进一步处理，比如通过 Wails 事件发送给前端
		}
	}()

	return client, nil
}

// SendRequest 发送 LSP 请求
func (c *LSPClient) SendRequest(method string, params interface{}) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed || c.conn == nil {
		return nil, fmt.Errorf("connection closed")
	}

	var result interface{}
	err := c.conn.Call(context.Background(), method, params, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Close 关闭连接
func (c *LSPClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}
	c.closed = true
	if c.conn != nil {
		c.conn.Close()
	}
	if c.cmd != nil && c.cmd.Process != nil {
		c.cmd.Process.Kill()
	}
	return nil
}

// pipe 实现了 io.ReadWriteCloser，用于连接 stdin/stdout
type pipe struct {
	io.WriteCloser
	io.Reader
}

func (p *pipe) Read(data []byte) (int, error) {
	return p.Reader.Read(data)
}

// noopHandler 处理来自服务器的通知
type noopHandler struct {
	DiagnosticsChan chan map[string]interface{}
}

func (h *noopHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	if req.Method == "textDocument/publishDiagnostics" {
		if h.DiagnosticsChan != nil && req.Params != nil {
			var params map[string]interface{}
			json.Unmarshal(*req.Params, &params)
			h.DiagnosticsChan <- params
		}
	}
}
