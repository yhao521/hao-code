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

// RenameSymbol 重命名符号
func (s *LSPService) RenameSymbol(languageID string, uri string, line int, col int, newName string) (map[string]interface{}, error) {
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
		"newName": newName,
	}

	result, err := client.SendRequest("textDocument/rename", params)
	if err != nil {
		return nil, err
	}

	if m, ok := result.(map[string]interface{}); ok {
		return m, nil
	}

	return nil, nil
}

// FormatDocument 格式化文档
func (s *LSPService) FormatDocument(languageID string, uri string, content string) ([]map[string]interface{}, error) {
	client, ok := s.clients[languageID]
	if !ok {
		return nil, fmt.Errorf("LSP client for %s not found", languageID)
	}

	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": uri,
		},
		"options": map[string]interface{}{
			"tabSize":                4,
			"insertSpaces":           true,
			"trimTrailingWhitespace": true,
			"insertFinalNewline":     true,
			"trimFinalNewlines":      true,
		},
	}

	result, err := client.SendRequest("textDocument/formatting", params)
	if err != nil {
		return nil, err
	}

	if items, ok := result.([]interface{}); ok {
		var edits []map[string]interface{}
		for _, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				edits = append(edits, m)
			}
		}
		return edits, nil
	}

	return nil, nil
}

// GetHoverInfo 获取悬停提示信息
func (s *LSPService) GetHoverInfo(languageID string, uri string, line int, col int) (map[string]interface{}, error) {
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

	result, err := client.SendRequest("textDocument/hover", params)
	if err != nil {
		return nil, err
	}

	if m, ok := result.(map[string]interface{}); ok {
		return m, nil
	}

	return nil, nil
}

// GetSignatureHelp 获取签名帮助
func (s *LSPService) GetSignatureHelp(languageID string, uri string, line int, col int) (map[string]interface{}, error) {
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

	result, err := client.SendRequest("textDocument/signatureHelp", params)
	if err != nil {
		return nil, err
	}

	if m, ok := result.(map[string]interface{}); ok {
		return m, nil
	}

	return nil, nil
}

// GetCodeActions 获取代码动作
func (s *LSPService) GetCodeActions(languageID string, uri string, startLine int, startCol int, endLine int, endCol int, diagnostics []map[string]interface{}) ([]map[string]interface{}, error) {
	client, ok := s.clients[languageID]
	if !ok {
		return nil, fmt.Errorf("LSP client for %s not found", languageID)
	}

	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": uri,
		},
		"range": map[string]interface{}{
			"start": map[string]int{"line": startLine, "character": startCol},
			"end":   map[string]int{"line": endLine, "character": endCol},
		},
		"context": map[string]interface{}{
			"diagnostics": diagnostics,
		},
	}

	result, err := client.SendRequest("textDocument/codeAction", params)
	if err != nil {
		return nil, err
	}

	if items, ok := result.([]interface{}); ok {
		var actions []map[string]interface{}
		for _, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				actions = append(actions, m)
			}
		}
		return actions, nil
	}

	return nil, nil
}

// GetFoldingRanges 获取折叠范围
func (s *LSPService) GetFoldingRanges(languageID string, uri string) ([]map[string]interface{}, error) {
	client, ok := s.clients[languageID]
	if !ok {
		return nil, fmt.Errorf("LSP client for %s not found", languageID)
	}

	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": uri,
		},
	}

	result, err := client.SendRequest("textDocument/foldingRange", params)
	if err != nil {
		return nil, err
	}

	if items, ok := result.([]interface{}); ok {
		var ranges []map[string]interface{}
		for _, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				ranges = append(ranges, m)
			}
		}
		return ranges, nil
	}

	return nil, nil
}

// GetSemanticTokens 获取语义高亮标记
func (s *LSPService) GetSemanticTokens(languageID string, uri string) (map[string]interface{}, error) {
	client, ok := s.clients[languageID]
	if !ok {
		return nil, fmt.Errorf("LSP client for %s not found", languageID)
	}

	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": uri,
		},
	}

	result, err := client.SendRequest("textDocument/semanticTokens/full", params)
	if err != nil {
		return nil, err
	}

	if m, ok := result.(map[string]interface{}); ok {
		return m, nil
	}

	return nil, nil
}

// GetDocumentLinks 获取文档链接
func (s *LSPService) GetDocumentLinks(languageID string, uri string) ([]map[string]interface{}, error) {
	client, ok := s.clients[languageID]
	if !ok {
		return nil, fmt.Errorf("LSP client for %s not found", languageID)
	}

	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": uri,
		},
	}

	result, err := client.SendRequest("textDocument/documentLink", params)
	if err != nil {
		return nil, err
	}

	if items, ok := result.([]interface{}); ok {
		var links []map[string]interface{}
		for _, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				links = append(links, m)
			}
		}
		return links, nil
	}

	return nil, nil
}

// GetCodeLenses 获取代码操作 (Code Lens)
func (s *LSPService) GetCodeLenses(languageID string, uri string) ([]map[string]interface{}, error) {
	client, ok := s.clients[languageID]
	if !ok {
		return nil, fmt.Errorf("LSP client for %s not found", languageID)
	}

	params := map[string]interface{}{
		"textDocument": map[string]string{
			"uri": uri,
		},
	}

	result, err := client.SendRequest("textDocument/codeLens", params)
	if err != nil {
		return nil, err
	}

	if items, ok := result.([]interface{}); ok {
		var lenses []map[string]interface{}
		for _, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				lenses = append(lenses, m)
			}
		}
		return lenses, nil
	}

	return nil, nil
}

// PrepareCallHierarchy 准备调用层级
func (s *LSPService) PrepareCallHierarchy(languageID string, uri string, line int, col int) ([]map[string]interface{}, error) {
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

	result, err := client.SendRequest("textDocument/prepareCallHierarchy", params)
	if err != nil {
		return nil, err
	}

	if items, ok := result.([]interface{}); ok {
		var itemsResult []map[string]interface{}
		for _, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				itemsResult = append(itemsResult, m)
			}
		}
		return itemsResult, nil
	}

	return nil, nil
}

// GetIncomingCalls 获取传入调用
func (s *LSPService) GetIncomingCalls(languageID string, item map[string]interface{}) ([]map[string]interface{}, error) {
	client, ok := s.clients[languageID]
	if !ok {
		return nil, fmt.Errorf("LSP client for %s not found", languageID)
	}

	params := map[string]interface{}{
		"item": item,
	}

	result, err := client.SendRequest("callHierarchy/incomingCalls", params)
	if err != nil {
		return nil, err
	}

	if items, ok := result.([]interface{}); ok {
		var calls []map[string]interface{}
		for _, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				calls = append(calls, m)
			}
		}
		return calls, nil
	}

	return nil, nil
}

// GetTypeHierarchy 获取类型层次结构
func (s *LSPService) GetTypeHierarchy(languageID string, uri string, line int, col int) ([]map[string]interface{}, error) {
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

	result, err := client.SendRequest("textDocument/prepareTypeHierarchy", params)
	if err != nil {
		return nil, err
	}

	if items, ok := result.([]interface{}); ok {
		var itemsResult []map[string]interface{}
		for _, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				itemsResult = append(itemsResult, m)
			}
		}
		return itemsResult, nil
	}

	return nil, nil
}

// GetImplementations 获取实现查找
func (s *LSPService) GetImplementations(languageID string, uri string, line int, col int) ([]map[string]interface{}, error) {
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

	result, err := client.SendRequest("textDocument/implementation", params)
	if err != nil {
		return nil, err
	}

	if items, ok := result.([]interface{}); ok {
		var impls []map[string]interface{}
		for _, item := range items {
			if m, ok := item.(map[string]interface{}); ok {
				impls = append(impls, m)
			}
		}
		return impls, nil
	}

	return nil, nil
}

// GetWorkspaceSymbols 获取工作区符号
func (s *LSPService) GetWorkspaceSymbols(query string) ([]map[string]interface{}, error) {
	// 这里简单起见，只使用第一个活跃的客户端。在实际项目中，可能需要遍历所有客户端或根据上下文选择。
	if len(s.clients) == 0 {
		return nil, fmt.Errorf("no active LSP clients")
	}

	// 获取第一个客户端作为示例
	var client *LSPClient
	for _, c := range s.clients {
		client = c
		break
	}

	params := map[string]interface{}{
		"query": query,
	}

	result, err := client.SendRequest("workspace/symbol", params)
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

// ResolveCodeAction 解析代码动作（深度定制）
func (s *LSPService) ResolveCodeAction(languageID string, action map[string]interface{}) (map[string]interface{}, error) {
	client, ok := s.clients[languageID]
	if !ok {
		return nil, fmt.Errorf("LSP client for %s not found", languageID)
	}

	result, err := client.SendRequest("codeAction/resolve", action)
	if err != nil {
		return nil, err
	}

	if m, ok := result.(map[string]interface{}); ok {
		return m, nil
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
