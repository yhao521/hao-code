package backend

import (
	"encoding/json"
	"fmt"
)

// PluginAPI 定义了暴露给插件沙箱的 API 集合
type PluginAPI struct {
	sandbox *PluginSandbox
}

// NewPluginAPI 创建一个新的插件 API 实例
func NewPluginAPI(sandbox *PluginSandbox) *PluginAPI {
	return &PluginAPI{sandbox: sandbox}
}

// HandleRequest 处理来自插件沙箱的请求
func (api *PluginAPI) HandleRequest(method string, params json.RawMessage) (interface{}, error) {
	switch method {
	case "vscode.window.showInformationMessage":
		var msg string
		if err := json.Unmarshal(params, &msg); err != nil {
			return nil, err
		}
		fmt.Printf("[Plugin Message] %s\n", msg)
		return nil, nil
	case "vscode.workspace.openTextDocument":
		var path string
		if err := json.Unmarshal(params, &path); err != nil {
			return nil, err
		}
		// 这里可以调用后端的文件服务
		return map[string]interface{}{"path": path}, nil
	default:
		return nil, fmt.Errorf("method not implemented: %s", method)
	}
}
