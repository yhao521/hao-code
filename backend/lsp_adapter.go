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
