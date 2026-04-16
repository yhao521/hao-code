package backend

import (
	"fmt"
	"os"
	"path/filepath"
)

// PluginInstance 代表一个已加载的插件实例
type PluginInstance struct {
	Manifest *PluginManifest
	Path     string
	IsActive bool
}

// PluginLoader 负责插件的发现、加载和管理
type PluginLoader struct {
	ExtensionsDir string
	Plugins       map[string]*PluginInstance
}

// NewPluginLoader 创建一个新的插件加载器
func NewPluginLoader() *PluginLoader {
	home, _ := os.UserHomeDir()
	return &PluginLoader{
		ExtensionsDir: filepath.Join(home, ".hao-code", "extensions"),
		Plugins:       make(map[string]*PluginInstance),
	}
}

// ScanAndLoad 扫描扩展目录并加载所有合法插件
func (pl *PluginLoader) ScanAndLoad() error {
	if _, err := os.Stat(pl.ExtensionsDir); os.IsNotExist(err) {
		os.MkdirAll(pl.ExtensionsDir, 0755)
		return nil
	}

	entries, err := os.ReadDir(pl.ExtensionsDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			pluginPath := filepath.Join(pl.ExtensionsDir, entry.Name())
			manifestPath := filepath.Join(pluginPath, "package.json")

			data, err := os.ReadFile(manifestPath)
			if err != nil {
				fmt.Printf("Failed to read manifest for %s: %v\n", entry.Name(), err)
				continue
			}

			manifest, err := ParseManifest(data)
			if err != nil {
				fmt.Printf("Invalid manifest for %s: %v\n", entry.Name(), err)
				continue
			}

			pl.Plugins[manifest.Name] = &PluginInstance{
				Manifest: manifest,
				Path:     pluginPath,
				IsActive: false,
			}
			fmt.Printf("Plugin loaded: %s@%s\n", manifest.Name, manifest.Version)
		}
	}
	return nil
}

// ActivatePlugin 激活指定插件
func (pl *PluginLoader) ActivatePlugin(name string) error {
	plugin, exists := pl.Plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	if plugin.IsActive {
		return nil
	}

	// TODO: 在此处执行实际的沙箱初始化逻辑
	plugin.IsActive = true
	fmt.Printf("Plugin activated: %s\n", name)
	return nil
}
