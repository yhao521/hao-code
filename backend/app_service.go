package backend

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// AppService 应用主服务，组合所有子服务
type AppService struct {
	fileSystem IFileSystemService
	git        IGitService
	config     IConfigService
	debug      *DebugService
	lsp        *LSPService
	loader     *PluginLoader
	bridge     *PluginBridge
	lifecycle  *PluginLifecycleManager
	tasks      *TaskService
	store      *PluginStore
	ai         *AIService
	terminals  map[string]*TerminalSession
	terminalMu sync.Mutex
}

// NewAppService 创建应用服务
func NewAppService(fs IFileSystemService, git IGitService, config IConfigService) *AppService {
	// 确定插件存储路径，例如用户目录下的 .hao-code/plugins
	homeDir, _ := os.UserHomeDir()
	pluginDir := filepath.Join(homeDir, ".hao-code", "plugins")

	return &AppService{
		fileSystem: fs,
		git:        git,
		config:     config,
		debug:      NewDebugService(),
		lsp:        NewLSPService(),
		loader:     NewPluginLoader(),
		bridge:     NewPluginBridge(),
		lifecycle:  NewPluginLifecycleManager(),
		store:      NewPluginStore(pluginDir),
		ai:         NewAIService(AIConfig{BaseURL: "https://api.openai.com/v1", Model: "gpt-3.5-turbo", MaxTokens: 100}),
		terminals:  make(map[string]*TerminalSession),
	}
}

// Greet 问候方法
func (a *AppService) Greet(name string) string {
	return "Hello " + name + ", Welcome to Hao-Code Editor!"
}

// ReadFile 读取文件（委托给文件系统服务）
func (a *AppService) ReadFile(path string) (string, error) {
	return a.fileSystem.ReadFile(path)
}

// WriteFile 写入文件
func (a *AppService) WriteFile(path string, content string) error {
	return a.fileSystem.WriteFile(path, content)
}

// ListDir 列出目录
func (a *AppService) ListDir(path string) ([]FileInfo, error) {
	return a.fileSystem.ListDir(path)
}

// GetProjectRoot 获取项目根目录
func (a *AppService) GetProjectRoot() string {
	return a.fileSystem.GetProjectRoot()
}

// OpenFolderDialog 打开文件夹选择对话框
func (a *AppService) OpenFolderDialog() (string, error) {
	return a.fileSystem.OpenFolderDialog()
}

// OpenFileDialog 打开文件选择对话框
func (a *AppService) OpenFileDialog() (string, error) {
	return a.fileSystem.OpenFileDialog()
}

// SaveFileDialog 保存文件对话框
func (a *AppService) SaveFileDialog() (string, error) {
	return a.fileSystem.SaveFileDialog()
}

// SetProjectRoot 设置项目根目录
func (a *AppService) SetProjectRoot(path string) error {
	return a.fileSystem.SetProjectRoot(path)
}

// OpenRepository 打开 Git 仓库（委托给 Git 服务）
func (a *AppService) OpenRepository(path string) (*RepoInfo, error) {
	return a.git.OpenRepository(path)
}

// GetGitStatus 获取 Git 状态
func (a *AppService) GetGitStatus(path string) (*GitStatus, error) {
	return a.git.GetGitStatus(path)
}

// GitCommit 提交 Git 更改
func (a *AppService) GitCommit(path, message string) (string, error) {
	return a.git.GitCommit(path, message)
}

// GitGetBranches 获取 Git 分支
func (a *AppService) GitGetBranches(path string) (*BranchInfo, error) {
	return a.git.GitGetBranches(path)
}

// GitGetLog 获取 Git 日志
func (a *AppService) GitGetLog(path string, maxCommits int) ([]CommitInfo, error) {
	return a.git.GitGetLog(path, maxCommits)
}

// GetGitGraph 获取 Git 图谱数据
func (a *AppService) GetGitGraph(path string, maxCommits int) ([]GitGraphNode, error) {
	return a.git.GetGitGraph(path, maxCommits)
}

// GetFileDiff 获取文件差异
func (a *AppService) GetFileDiff(path, filePath string) (*FileDiff, error) {
	return a.git.GetFileDiff(path, filePath)
}

// StartTerminal 启动终端会话（返回会话 ID）
var terminalSessions = make(map[string]*TerminalSession)

func (a *AppService) StartTerminal() (string, error) {
	session, err := NewTerminalSession()
	if err != nil {
		return "", err
	}
	id := fmt.Sprintf("term_%d", time.Now().UnixNano())
	terminalSessions[id] = session
	return id, nil
}

// CreateFile 创建新文件
func (a *AppService) CreateFile(path string) error {
	return a.fileSystem.CreateFile(path)
}

// CreateDirectory 创建新目录
func (a *AppService) CreateDirectory(path string) error {
	return a.fileSystem.CreateDirectory(path)
}

// DeleteFileOrDirectory 删除文件或目录
func (a *AppService) DeleteFileOrDirectory(path string) error {
	return a.fileSystem.DeleteFileOrDirectory(path)
}

// RenameFileOrDirectory 重命名文件或目录
func (a *AppService) RenameFileOrDirectory(oldPath, newPath string) error {
	return a.fileSystem.RenameFileOrDirectory(oldPath, newPath)
}

// MoveFileOrDirectory 移动文件或目录
func (a *AppService) MoveFileOrDirectory(sourcePath, targetPath string) error {
	return a.fileSystem.MoveFileOrDirectory(sourcePath, targetPath)
}

// GetFileStats 获取文件统计信息
func (a *AppService) GetFileStats(path string) (*FileInfo, error) {
	return a.fileSystem.GetFileStats(path)
}

// SearchFiles 搜索文件
func (a *AppService) SearchFiles(rootPath, keyword string, maxResults int) ([]FileInfo, error) {
	return a.fileSystem.SearchFiles(rootPath, keyword, maxResults)
}

// CopyFileOrDirectory 复制文件或目录
func (a *AppService) CopyFileOrDirectory(sourcePath, targetPath string) error {
	return a.fileSystem.CopyFileOrDirectory(sourcePath, targetPath)
}

// IsTextFile 判断是否为文本文件
func (a *AppService) IsTextFile(path string) bool {
	return a.fileSystem.IsTextFile(path)
}

// GetFileExtension 获取文件扩展名
func (a *AppService) GetFileExtension(path string) string {
	return a.fileSystem.GetFileExtension(path)
}

// GetDirectoryTree 获取目录树
func (a *AppService) GetDirectoryTree(path string, depth int) ([]FileInfo, error) {
	return a.fileSystem.GetDirectoryTree(path, depth)
}

// BackupFile 备份文件
func (a *AppService) BackupFile(path string) error {
	return a.fileSystem.BackupFile(path)
}

// TouchFile Touch 文件
func (a *AppService) TouchFile(path string) error {
	return a.fileSystem.TouchFile(path)
}

// ==================== 配置服务方法 ====================

// AddRecentFile 添加最近打开的文件
func (a *AppService) AddRecentFile(path string) error {
	return a.config.AddRecentFile(path)
}

// AddRecentFolder 添加最近打开的文件夹
func (a *AppService) AddRecentFolder(path string) error {
	return a.config.AddRecentFolder(path)
}

// GetRecentFiles 获取最近打开的文件列表
func (a *AppService) GetRecentFiles() []RecentItem {
	return a.config.GetRecentFiles()
}

// GetRecentFolders 获取最近打开的文件夹列表
func (a *AppService) GetRecentFolders() []RecentItem {
	return a.config.GetRecentFolders()
}

// RemoveRecentFile 从最近文件列表中移除指定文件
func (a *AppService) RemoveRecentFile(path string) error {
	return a.config.RemoveRecentFile(path)
}

// RemoveRecentFolder 从最近文件夹列表中移除指定文件夹
func (a *AppService) RemoveRecentFolder(path string) error {
	return a.config.RemoveRecentFolder(path)
}

// ClearRecentFiles 清空最近文件列表
func (a *AppService) ClearRecentFiles() error {
	return a.config.ClearRecentFiles()
}

// ClearRecentFolders 清空最近文件夹列表
func (a *AppService) ClearRecentFolders() error {
	return a.config.ClearRecentFolders()
}

// SearchInFiles 在所有文件中搜索文本
func (a *AppService) SearchInFiles(opts SearchOptions) ([]SearchResult, error) {
	var results []SearchResult
	rootPath := opts.RootPath
	searchText := opts.Query
	caseSensitive := opts.CaseSensitive

	err := filepath.WalkDir(rootPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil || len(results) >= 100 { // 限制最大结果数
			return nil
		}

		// 跳过隐藏文件和 node_modules
		name := entry.Name()
		if strings.HasPrefix(name, ".") || name == "node_modules" || name == ".git" {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 排除特定模式 (简单实现)
		if opts.Exclude != "" {
			for _, pattern := range strings.Split(opts.Exclude, ",") {
				pattern = strings.TrimSpace(pattern)
				if matched, _ := filepath.Match(pattern, name); matched {
					return nil
				}
			}
		}

		// 只搜索文本文件
		if !entry.IsDir() && a.fileSystem.IsTextFile(path) {
			content, _ := a.fileSystem.ReadFile(path)

			var searchFunc func(string, string) bool
			if caseSensitive {
				searchFunc = strings.Contains
			} else {
				searchFunc = func(s, substr string) bool {
					return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
				}
			}

			if searchFunc(content, searchText) {
				// 找到匹配的行
				lines := strings.Split(content, "\n")
				for i, line := range lines {
					if searchFunc(line, searchText) {
						results = append(results, SearchResult{
							FilePath:    path,
							LineNumber:  i + 1,
							LineContent: line,
						})

						if len(results) >= 100 {
							break
						}
					}
				}
			}
		}

		return nil
	})

	return results, err
}

// ==================== 插件系统方法 ====================

// GetInstalledPlugins 获取已安装的插件列表
func (a *AppService) GetInstalledPlugins() []PluginManifest {
	var manifests []PluginManifest
	for _, p := range a.loader.Plugins {
		manifests = append(manifests, *p.Manifest)
	}
	return manifests
}

// ActivatePlugin 激活指定插件
func (a *AppService) ActivatePlugin(name string) error {
	// 1. 从加载器中获取 Manifest
	if err := a.loader.ScanAndLoad(); err != nil {
		return err
	}
	plugin, exists := a.loader.Plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	// 2. 通过生命周期管理器进行状态管理
	return a.lifecycle.ActivatePlugin(name, plugin.Manifest)
}

// DeactivatePlugin 停用指定插件
func (a *AppService) DeactivatePlugin(name string) error {
	return a.lifecycle.DeactivatePlugin(name)
}

// ExecutePluginCommand 执行插件命令
func (a *AppService) ExecutePluginCommand(command string, payload interface{}) (interface{}, error) {
	return a.bridge.ExecuteCommand(command, payload)
}

// ==================== 任务运行器方法 ====================

// GetTasks 获取工作区任务列表
func (a *AppService) GetTasks(rootPath string) ([]TaskItem, error) {
	return a.tasks.GetTasks(rootPath)
}

// RunTask 运行指定任务
func (a *AppService) RunTask(rootPath string, command string) error {
	return a.tasks.RunTask(rootPath, command)
}

// ==================== 插件商店方法 ====================

// InstallPluginFromURL 从 URL 安装插件
func (a *AppService) InstallPluginFromURL(url string) (*PluginManifest, error) {
	manifest, err := a.store.InstallFromURL(url)
	if err != nil {
		return nil, err
	}
	// 安装后自动加载
	a.loader.ScanAndLoad()
	return manifest, nil
}

// UninstallPlugin 卸载插件
func (a *AppService) UninstallPlugin(name string) error {
	// 先停用
	a.lifecycle.DeactivatePlugin(name)
	return a.store.Uninstall(name)
}

// GetInstalledPluginNames 获取已安装插件名称列表
func (a *AppService) GetInstalledPluginNames() ([]string, error) {
	return a.store.ListInstalled()
}

// ==================== AI 助手方法 ====================

// GetGhostText 获取行内代码建议
func (a *AppService) GetGhostText(prefix string, suffix string, language string, filePath string) (*GhostTextResponse, error) {
	req := GhostTextRequest{
		Prefix:   prefix,
		Suffix:   suffix,
		Language: language,
		FilePath: filePath,
	}
	return a.ai.GetGhostText(req)
}

// SetAIConfig 设置 AI 配置
func (a *AppService) SetAIConfig(apiKey string, baseURL string, model string) {
	cfg := AIConfig{
		APIKey:    apiKey,
		BaseURL:   baseURL,
		Model:     model,
		MaxTokens: 100,
	}
	a.ai.UpdateConfig(cfg)
}

// ChatWithAI 与 AI 聊天
func (a *AppService) ChatWithAI(messages []ChatMessage, context string) (*ChatResponse, error) {
	req := ChatRequest{
		Messages: messages,
		Context:  context,
	}
	return a.ai.ChatWithAI(req)
}

// GetAIConfig 获取 AI 配置
func (a *AppService) GetAIConfig() AIConfig {
	return a.ai.GetAIConfig()
}

// ==================== Git Staging 增强方法 ====================

// StageSelectedRanges 暂存指定行范围
func (a *AppService) StageSelectedRanges(path, filePath string, ranges []LineRange) error {
	return a.git.StageSelectedRanges(path, filePath, ranges)
}

// UnstageFile 取消暂存文件
func (a *AppService) UnstageFile(path, filePath string) error {
	return a.git.UnstageFile(path, filePath)
}

// GetAIContextFromFiles 获取指定文件的 AI 上下文内容
func (a *AppService) GetAIContextFromFiles(rootPath string, references []string) string {
	return a.ai.BuildContextFromReferences(rootPath, references)
}

// ==================== 终端分屏管理方法 ====================

// CreateTerminal 创建一个新的终端会话并返回 ID
func (a *AppService) CreateTerminal() (string, error) {
	a.terminalMu.Lock()
	defer a.terminalMu.Unlock()

	session, err := NewTerminalSession()
	if err != nil {
		return "", err
	}

	id := fmt.Sprintf("term-%d", time.Now().UnixNano())
	a.terminals[id] = session
	return id, nil
}

// WriteToTerminal 向指定终端写入数据
func (a *AppService) WriteToTerminal(id string, data string) error {
	a.terminalMu.Lock()
	session, ok := a.terminals[id]
	a.terminalMu.Unlock()

	if !ok {
		return fmt.Errorf("terminal %s not found", id)
	}
	_, err := session.Write([]byte(data))
	return err
}

// ResizeTerminal 调整指定终端的大小
func (a *AppService) ResizeTerminal(id string, cols, rows int) error {
	a.terminalMu.Lock()
	session, ok := a.terminals[id]
	a.terminalMu.Unlock()

	if !ok {
		return fmt.Errorf("terminal %s not found", id)
	}
	return session.Resize(uint32(cols), uint32(rows))
}

// CloseTerminal 关闭指定终端
func (a *AppService) CloseTerminal(id string) error {
	a.terminalMu.Lock()
	session, ok := a.terminals[id]
	if ok {
		session.Close()
		delete(a.terminals, id)
	}
	a.terminalMu.Unlock()
	return nil
}

// GetFileBlame 获取文件 Blame 信息
func (a *AppService) GetFileBlame(path, filePath string) ([]BlameInfo, error) {
	return a.git.GetFileBlame(path, filePath)
}

// GetFileHistory 获取文件历史记录
func (a *AppService) GetFileHistory(path, filePath string) ([]CommitInfo, error) {
	return a.git.GetFileHistory(path, filePath)
}

// GetDiagnostics 获取 LSP 诊断信息
func (a *AppService) GetDiagnostics(languageID string, uri string) ([]map[string]interface{}, error) {
	return a.lsp.GetDiagnostics(languageID, uri)
}

// ScanTodos 扫描项目中的 TODO/FIXME 注释
func (a *AppService) ScanTodos(rootPath string) ([]SearchResult, error) {
	opts := SearchOptions{
		RootPath:      rootPath,
		Query:         "(TODO|FIXME|HACK|XXX|BUG):",
		UseRegex:      true,
		CaseSensitive: true,
		Exclude:       "node_modules,.git,vendor,dist,build",
	}
	return a.fileSystem.SearchFilesWithOptions(opts)
}

// SearchFilesWithOptions 代理文件系统的高级搜索
func (a *AppService) SearchFilesWithOptions(opts SearchOptions) ([]SearchResult, error) {
	return a.fileSystem.SearchFilesWithOptions(opts)
}

// SendHTTPRequest 发送 HTTP 请求 (API Tester)
func (a *AppService) SendHTTPRequest(req APIRequest) (*APIResponse, error) {
	// 发送成功后自动保存历史记录
	resp, err := SendHTTPRequest(req)
	if err == nil {
		SaveApiHistory(req)
	}
	return resp, err
}

// GetApiHistory 获取 API 历史记录
func (a *AppService) GetApiHistory() ([]APIHistoryItem, error) {
	return GetApiHistory()
}

// DeleteApiHistory 删除指定的历史记录
func (a *AppService) DeleteApiHistory(id string) error {
	return DeleteApiHistory(id)
}

// GetEnvVariables 获取环境变量
func (a *AppService) GetEnvVariables() (map[string]string, error) {
	return GetEnvVariables()
}

// SaveEnvVariables 保存环境变量
func (a *AppService) SaveEnvVariables(vars map[string]string) error {
	return SaveEnvVariables(vars)
}
