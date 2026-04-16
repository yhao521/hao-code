package backend

import (
	"fmt"
	"log"
)

// LSPService 管理语言服务器实例
type LSPService struct {
	clients map[string]*LSPClient // key: languageId
}

func NewLSPService() *LSPService {
	return &LSPService{
		clients: make(map[string]*LSPClient),
	}
}

// InitializeLanguage 初始化指定语言的 LSP 服务器
func (s *LSPService) InitializeLanguage(languageID string, rootPath string) error {
	if _, exists := s.clients[languageID]; exists {
		return nil // 已经初始化
	}

	var cmd string
	var args []string

	switch languageID {
	case "go":
		cmd = "gopls"
		args = []string{"-rpc.trace"}
	case "typescript", "javascript":
		cmd = "typescript-language-server"
		args = []string{"--stdio"}
	default:
		return fmt.Errorf("unsupported language: %s", languageID)
	}

	client, err := NewLSPClient(cmd, args...)
	if err != nil {
		return err
	}

	// 发送 initialize 请求
	params := map[string]interface{}{
		"processId":             nil,
		"rootUri":               "file://" + rootPath,
		"initializationOptions": map[string]interface{}{},
		"capabilities":          map[string]interface{}{},
	}

	_, err = client.SendRequest("initialize", params)
	if err != nil {
		client.Close()
		return fmt.Errorf("failed to initialize LSP: %w", err)
	}

	// 发送 initialized 通知
	client.SendRequest("initialized", map[string]interface{}{})

	s.clients[languageID] = client
	log.Printf("LSP server for %s initialized successfully", languageID)
	return nil
}

// GetCompletions 获取代码补全建议
func (s *LSPService) GetCompletions(languageID string, uri string, line int, col int) ([]map[string]interface{}, error) {
	client, ok := s.clients[languageID]
	if !ok {
		return nil, fmt.Errorf("LSP client for %s not found", languageID)
	}

	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": uri,
		},
		"position": map[string]int{
			"line":      line,
			"character": col,
		},
	}

	result, err := client.SendRequest("textDocument/completion", params)
	if err != nil {
		return nil, err
	}

	// 解析结果 (简化处理)
	if items, ok := result.([]interface{}); ok {
		var completions []map[string]interface{}
		for _, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				completions = append(completions, m)
			}
		}
		return completions, nil
	}

	return nil, nil
}

// GetDefinition 获取定义位置
func (s *LSPService) GetDefinition(languageID string, uri string, line int, col int) (map[string]interface{}, error) {
	client, ok := s.clients[languageID]
	if !ok {
		return nil, fmt.Errorf("LSP client for %s not found", languageID)
	}

	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": uri,
		},
		"position": map[string]int{
			"line":      line,
			"character": col,
		},
	}

	result, err := client.SendRequest("textDocument/definition", params)
	if err != nil {
		return nil, err
	}

	if m, ok := result.(map[string]interface{}); ok {
		return m, nil
	}

	return nil, nil
}

// GetDocumentSymbols 获取文档符号大纲
func (s *LSPService) GetDocumentSymbols(languageID string, uri string) ([]map[string]interface{}, error) {
	client, ok := s.clients[languageID]
	if !ok {
		return nil, fmt.Errorf("LSP client for %s not found", languageID)
	}

	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": uri,
		},
	}

	result, err := client.SendRequest("textDocument/documentSymbol", params)
	if err != nil {
		return nil, err
	}

	if items, ok := result.([]interface{}); ok {
		var symbols []map[string]interface{}
		for _, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				symbols = append(symbols, m)
			}
		}
		return symbols, nil
	}

	return nil, nil
}

// FindReferences 查找引用
func (s *LSPService) FindReferences(languageID string, uri string, line int, col int) ([]map[string]interface{}, error) {
	client, ok := s.clients[languageID]
	if !ok {
		return nil, fmt.Errorf("LSP client for %s not found", languageID)
	}

	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": uri,
		},
		"position": map[string]int{
			"line":      line,
			"character": col,
		},
		"context": map[string]bool{
			"includeDeclaration": true,
		},
	}

	result, err := client.SendRequest("textDocument/references", params)
	if err != nil {
		return nil, err
	}

	if items, ok := result.([]interface{}); ok {
		var refs []map[string]interface{}
		for _, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				refs = append(refs, m)
			}
		}
		return refs, nil
	}

	return nil, nil
}

// Shutdown 关闭所有 LSP 服务
func (s *LSPService) Shutdown() {
	for lang, client := range s.clients {
		log.Printf("Shutting down LSP for %s", lang)
		client.Close()
	}
}
