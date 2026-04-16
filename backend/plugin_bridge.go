package backend

import (
	"fmt"
	"log"
)

// PluginBridge 负责前后端插件通信的桥接
type PluginBridge struct {
	handlers map[string]func(payload interface{}) interface{}
}

// NewPluginBridge 创建一个新的插件通信桥接
func NewPluginBridge() *PluginBridge {
	return &PluginBridge{
		handlers: make(map[string]func(payload interface{}) interface{}),
	}
}

// RegisterHandler 注册一个处理器，用于处理来自前端插件的请求
func (pb *PluginBridge) RegisterHandler(command string, handler func(payload interface{}) interface{}) {
	pb.handlers[command] = handler
	log.Printf("Plugin command registered: %s", command)
}

// ExecuteCommand 执行插件命令并返回结果
func (pb *PluginBridge) ExecuteCommand(command string, payload interface{}) (interface{}, error) {
	handler, exists := pb.handlers[command]
	if !exists {
		return nil, fmt.Errorf("plugin command not found: %s", command)
	}

	result := handler(payload)
	return result, nil
}

// EmitEvent 模拟向后端发送事件（实际在 Wails v3 中通常由前端直接调用 Bindings）
// 这里主要用于后端主动触发某些逻辑后通知前端，前端通过 window.wails.events.on 监听
func (pb *PluginBridge) EmitEvent(eventName string, data interface{}) {
	log.Printf("Emitting plugin event: %s", eventName)
	// 在实际 Wails v3 实现中，这里会调用 application.EmitEvent
}
