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
