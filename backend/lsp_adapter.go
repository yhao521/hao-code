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
