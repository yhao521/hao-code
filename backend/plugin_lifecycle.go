package backend

import (
	"fmt"
	"log"
	"sync"
)

// PluginState 定义插件的生命周期状态
type PluginState int

const (
	StateInactive PluginState = iota
	StateActivating
	StateActive
	StateDeactivating
)

// PluginLifecycleInstance 代表一个运行中的插件实例
// (Renamed to avoid conflict with plugin_loader.go)
type PluginLifecycleInstance struct {
	Manifest *PluginManifest
	State    PluginState
	Context  *PluginContext // 沙箱上下文
}

// PluginContext 提供插件运行的受限环境
type PluginContext struct {
	StoragePath string
	Config      map[string]interface{}
}

// PluginLifecycleManager 管理插件的生命周期
type PluginLifecycleManager struct {
	instances map[string]*PluginLifecycleInstance
	mu        sync.RWMutex
}

// NewPluginLifecycleManager 创建生命周期管理器
func NewPluginLifecycleManager() *PluginLifecycleManager {
	return &PluginLifecycleManager{
		instances: make(map[string]*PluginLifecycleInstance),
	}
}

// ActivatePlugin 激活插件
func (plm *PluginLifecycleManager) ActivatePlugin(name string, manifest *PluginManifest) error {
	plm.mu.Lock()
	defer plm.mu.Unlock()

	if inst, exists := plm.instances[name]; exists {
		if inst.State == StateActive {
			return fmt.Errorf("plugin %s is already active", name)
		}
	}

	log.Printf("Activating plugin: %s v%s", name, manifest.Version)

	inst := &PluginLifecycleInstance{
		Manifest: manifest,
		State:    StateActivating,
		Context: &PluginContext{
			StoragePath: fmt.Sprintf("/tmp/hao-code-plugins/%s", name),
			Config:      make(map[string]interface{}),
		},
	}

	// 模拟沙箱初始化逻辑
	// 在实际生产中，这里会启动一个隔离的 V8 引擎或子进程
	plm.instances[name] = inst
	inst.State = StateActive

	log.Printf("Plugin %s activated successfully.", name)
	return nil
}

// DeactivatePlugin 停用插件
func (plm *PluginLifecycleManager) DeactivatePlugin(name string) error {
	plm.mu.Lock()
	defer plm.mu.Unlock()

	inst, exists := plm.instances[name]
	if !exists || inst.State != StateActive {
		return fmt.Errorf("plugin %s is not active", name)
	}

	log.Printf("Deactivating plugin: %s", name)
	inst.State = StateDeactivating

	// 清理资源：关闭文件句柄、断开网络连接等
	// 模拟清理过程
	inst.State = StateInactive

	log.Printf("Plugin %s deactivated.", name)
	return nil
}

// GetPluginState 获取插件当前状态
func (plm *PluginLifecycleManager) GetPluginState(name string) PluginState {
	plm.mu.RLock()
	defer plm.mu.RUnlock()

	if inst, exists := plm.instances[name]; exists {
		return inst.State
	}
	return StateInactive
}
