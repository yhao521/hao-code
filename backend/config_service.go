package backend

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ConfigManager 配置管理器，负责持久化用户配置
type ConfigManager struct {
	configPath string
	config     *AppConfig
	mu         sync.RWMutex
}

// AppConfig 应用配置
type AppConfig struct {
	RecentFiles  []RecentItem `json:"recentFiles"`
	RecentFolders []RecentItem `json:"recentFolders"`
}

// NewConfigManager 创建配置管理器
func NewConfigManager() *ConfigManager {
	configDir := getConfigDir()
	configPath := filepath.Join(configDir, "config.json")
	
	cm := &ConfigManager{
		configPath: configPath,
		config: &AppConfig{
			RecentFiles:   []RecentItem{},
			RecentFolders: []RecentItem{},
		},
	}
	
	// 加载配置
	cm.loadConfig()
	
	return cm
}

// getConfigDir 获取配置目录
func getConfigDir() string {
	homeDir, _ := os.UserHomeDir()
	
	// 根据操作系统选择配置目录
	configDir := filepath.Join(homeDir, ".hao-code")
	
	// 确保目录存在
	os.MkdirAll(configDir, 0755)
	
	return configDir
}

// loadConfig 从文件加载配置
func (cm *ConfigManager) loadConfig() {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		// 配置文件不存在或读取失败，使用默认配置
		return
	}
	
	if err := json.Unmarshal(data, cm.config); err != nil {
		// JSON 解析失败，使用默认配置
		fmt.Printf("Failed to parse config: %v\n", err)
		return
	}
}

// saveConfig 保存配置到文件
func (cm *ConfigManager) saveConfig() error {
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}
	
	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %v", err)
	}
	
	return nil
}

// AddRecentFile 添加最近打开的文件
func (cm *ConfigManager) AddRecentFile(path string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	// 移除已存在的路径
	cm.config.RecentFiles = filterRecentItems(cm.config.RecentFiles, path)
	
	// 添加到开头
	item := RecentItem{
		Path:     path,
		Name:     filepath.Base(path),
		OpenedAt: time.Now().Format(time.RFC3339),
	}
	cm.config.RecentFiles = append([]RecentItem{item}, cm.config.RecentFiles...)
	
	// 限制数量（最多20个）
	if len(cm.config.RecentFiles) > 20 {
		cm.config.RecentFiles = cm.config.RecentFiles[:20]
	}
	
	return cm.saveConfig()
}

// AddRecentFolder 添加最近打开的文件夹
func (cm *ConfigManager) AddRecentFolder(path string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	// 移除已存在的路径
	cm.config.RecentFolders = filterRecentItems(cm.config.RecentFolders, path)
	
	// 添加到开头
	item := RecentItem{
		Path:     path,
		Name:     filepath.Base(path),
		OpenedAt: time.Now().Format(time.RFC3339),
	}
	cm.config.RecentFolders = append([]RecentItem{item}, cm.config.RecentFolders...)
	
	// 限制数量（最多10个）
	if len(cm.config.RecentFolders) > 10 {
		cm.config.RecentFolders = cm.config.RecentFolders[:10]
	}
	
	return cm.saveConfig()
}

// GetRecentFiles 获取最近打开的文件列表
func (cm *ConfigManager) GetRecentFiles() []RecentItem {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	// 过滤掉不存在的文件
	result := make([]RecentItem, 0, len(cm.config.RecentFiles))
	for _, item := range cm.config.RecentFiles {
		if _, err := os.Stat(item.Path); err == nil {
			result = append(result, item)
		}
	}
	
	return result
}

// GetRecentFolders 获取最近打开的文件夹列表
func (cm *ConfigManager) GetRecentFolders() []RecentItem {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	// 过滤掉不存在的文件夹
	result := make([]RecentItem, 0, len(cm.config.RecentFolders))
	for _, item := range cm.config.RecentFolders {
		if info, err := os.Stat(item.Path); err == nil && info.IsDir() {
			result = append(result, item)
		}
	}
	
	return result
}

// RemoveRecentFile 从最近文件列表中移除指定文件
func (cm *ConfigManager) RemoveRecentFile(path string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.config.RecentFiles = filterRecentItems(cm.config.RecentFiles, path)
	return cm.saveConfig()
}

// RemoveRecentFolder 从最近文件夹列表中移除指定文件夹
func (cm *ConfigManager) RemoveRecentFolder(path string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.config.RecentFolders = filterRecentItems(cm.config.RecentFolders, path)
	return cm.saveConfig()
}

// ClearRecentFiles 清空最近文件列表
func (cm *ConfigManager) ClearRecentFiles() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.config.RecentFiles = []RecentItem{}
	return cm.saveConfig()
}

// ClearRecentFolders 清空最近文件夹列表
func (cm *ConfigManager) ClearRecentFolders() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.config.RecentFolders = []RecentItem{}
	return cm.saveConfig()
}

// filterRecentItems 过滤掉指定路径的项目
func filterRecentItems(items []RecentItem, pathToRemove string) []RecentItem {
	result := make([]RecentItem, 0, len(items))
	for _, item := range items {
		if item.Path != pathToRemove {
			result = append(result, item)
		}
	}
	return result
}
