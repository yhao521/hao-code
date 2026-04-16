package backend

// InitializeLSP 初始化指定语言的 LSP 服务
func (a *AppService) InitializeLSP(languageID string, rootPath string) error {
	return a.lsp.InitializeLanguage(languageID, rootPath)
}

// GetCompletions 获取代码补全建议
func (a *AppService) GetCompletions(languageID string, uri string, line int, col int) ([]map[string]interface{}, error) {
	return a.lsp.GetCompletions(languageID, uri, line, col)
}

// GetDefinition 获取定义位置
func (a *AppService) GetDefinition(languageID string, uri string, line int, col int) (map[string]interface{}, error) {
	return a.lsp.GetDefinition(languageID, uri, line, col)
}

// GetDocumentSymbols 获取文档符号大纲
func (a *AppService) GetDocumentSymbols(languageID string, uri string) ([]map[string]interface{}, error) {
	return a.lsp.GetDocumentSymbols(languageID, uri)
}

// FindReferences 查找引用
func (a *AppService) FindReferences(languageID string, uri string, line int, col int) ([]map[string]interface{}, error) {
	return a.lsp.FindReferences(languageID, uri, line, col)
}

// RenameSymbol 重命名符号
func (a *AppService) RenameSymbol(languageID string, uri string, line int, col int, newName string) (map[string]interface{}, error) {
	return a.lsp.RenameSymbol(languageID, uri, line, col, newName)
}

// FormatDocument 格式化文档
func (a *AppService) FormatDocument(languageID string, uri string, content string) ([]map[string]interface{}, error) {
	return a.lsp.FormatDocument(languageID, uri, content)
}

// GetDiagnosticsCount 获取 Linter 错误统计
func (a *AppService) GetDiagnosticsCount(languageID string, uri string) (map[string]int, error) {
	return a.lsp.GetDiagnosticsCount(languageID, uri)
}

// GetHoverInfo 获取悬停提示信息
func (a *AppService) GetHoverInfo(languageID string, uri string, line int, col int) (map[string]interface{}, error) {
	return a.lsp.GetHoverInfo(languageID, uri, line, col)
}

// GetSignatureHelp 获取签名帮助
func (a *AppService) GetSignatureHelp(languageID string, uri string, line int, col int) (map[string]interface{}, error) {
	return a.lsp.GetSignatureHelp(languageID, uri, line, col)
}

// GetCodeActions 获取代码动作
func (a *AppService) GetCodeActions(languageID string, uri string, startLine int, startCol int, endLine int, endCol int, diagnostics []map[string]interface{}) ([]map[string]interface{}, error) {
	return a.lsp.GetCodeActions(languageID, uri, startLine, startCol, endLine, endCol, diagnostics)
}

// GetFoldingRanges 获取折叠范围
func (a *AppService) GetFoldingRanges(languageID string, uri string) ([]map[string]interface{}, error) {
	return a.lsp.GetFoldingRanges(languageID, uri)
}

// GetSemanticTokens 获取语义高亮标记
func (a *AppService) GetSemanticTokens(languageID string, uri string) (map[string]interface{}, error) {
	return a.lsp.GetSemanticTokens(languageID, uri)
}

// GetDocumentLinks 获取文档链接
func (a *AppService) GetDocumentLinks(languageID string, uri string) ([]map[string]interface{}, error) {
	return a.lsp.GetDocumentLinks(languageID, uri)
}

// GetCodeLenses 获取代码操作 (Code Lens)
func (a *AppService) GetCodeLenses(languageID string, uri string) ([]map[string]interface{}, error) {
	return a.lsp.GetCodeLenses(languageID, uri)
}

// PrepareCallHierarchy 准备调用层级
func (a *AppService) PrepareCallHierarchy(languageID string, uri string, line int, col int) ([]map[string]interface{}, error) {
	return a.lsp.PrepareCallHierarchy(languageID, uri, line, col)
}

// GetIncomingCalls 获取传入调用
func (a *AppService) GetIncomingCalls(languageID string, item map[string]interface{}) ([]map[string]interface{}, error) {
	return a.lsp.GetIncomingCalls(languageID, item)
}

// GetTypeHierarchy 获取类型层次结构
func (a *AppService) GetTypeHierarchy(languageID string, uri string, line int, col int) ([]map[string]interface{}, error) {
	return a.lsp.GetTypeHierarchy(languageID, uri, line, col)
}

// GetImplementations 获取实现查找
func (a *AppService) GetImplementations(languageID string, uri string, line int, col int) ([]map[string]interface{}, error) {
	return a.lsp.GetImplementations(languageID, uri, line, col)
}

// GetWorkspaceSymbols 获取工作区符号
func (a *AppService) GetWorkspaceSymbols(query string) ([]map[string]interface{}, error) {
	return a.lsp.GetWorkspaceSymbols(query)
}

// ResolveCodeAction 解析代码动作
func (a *AppService) ResolveCodeAction(languageID string, action map[string]interface{}) (map[string]interface{}, error) {
	return a.lsp.ResolveCodeAction(languageID, action)
}
