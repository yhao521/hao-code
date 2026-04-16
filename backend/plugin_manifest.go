package backend

import "encoding/json"

// PluginManifest 定义插件的元数据 (package.json)
type PluginManifest struct {
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	Description string      `json:"description"`
	Main        string      `json:"main"` // 入口文件路径
	Author      string      `json:"author"`
	License     string      `json:"license"`
	Contributes Contributes `json:"contributes"`
}

// Contributes 定义插件对编辑器的扩展点
type Contributes struct {
	Commands  []CommandContribution  `json:"commands,omitempty"`
	Languages []LanguageContribution `json:"languages,omitempty"`
}

// CommandContribution 定义一个命令
type CommandContribution struct {
	Command string `json:"command"`
	Title   string `json:"title"`
}

// LanguageContribution 定义语言支持
type LanguageContribution struct {
	ID         string   `json:"id"`
	Extensions []string `json:"extensions"`
}

// ParseManifest 从 JSON 字节流解析 Manifest
func ParseManifest(data []byte) (*PluginManifest, error) {
	var manifest PluginManifest
	err := json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}
	return &manifest, nil
}
