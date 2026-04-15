package backend

import (
	"context"
)

// ==================== 核心业务接口定义 ====================
// 这些接口定义了前端与后端交互的契约
// 无论 Wails v2 还是 v3，只要实现这些接口即可

// IFileSystemService 文件系统服务接口
type IFileSystemService interface {
	// 读取文件内容
	ReadFile(path string) (string, error)
	// 写入文件内容
	WriteFile(path string, content string) error
	// 列出目录内容
	ListDir(path string) ([]FileInfo, error)
	// 获取项目根目录
	GetProjectRoot() string
	// 打开文件夹选择对话框
	OpenFolderDialog() (string, error)
	// 打开文件选择对话框
	OpenFileDialog() (string, error)
	// 保存文件对话框
	SaveFileDialog() (string, error)
	// 设置当前工作目录
	SetProjectRoot(path string) error
	// 创建文件
	CreateFile(path string) error
	// 创建目录
	CreateDirectory(path string) error
	// 删除文件或目录
	DeleteFileOrDirectory(path string) error
	// 重命名文件或目录
	RenameFileOrDirectory(oldPath, newPath string) error
	// 移动文件或目录
	MoveFileOrDirectory(sourcePath, targetPath string) error
	// 复制文件或目录
	CopyFileOrDirectory(sourcePath, targetPath string) error
	// 获取文件统计信息
	GetFileStats(path string) (*FileInfo, error)
	// 搜索文件
	SearchFiles(rootPath, keyword string, maxResults int) ([]FileInfo, error)
	// 判断是否为文本文件
	IsTextFile(path string) bool
	// 获取文件扩展名
	GetFileExtension(path string) string
	// 获取目录树
	GetDirectoryTree(path string, depth int) ([]FileInfo, error)
	// 备份文件
	BackupFile(path string) error
	// Touch文件
	TouchFile(path string) error
}

// IConfigService 配置服务接口
type IConfigService interface {
	// 添加最近打开的文件
	AddRecentFile(path string) error
	// 添加最近打开的文件夹
	AddRecentFolder(path string) error
	// 获取最近打开的文件列表
	GetRecentFiles() []RecentItem
	// 获取最近打开的文件夹列表
	GetRecentFolders() []RecentItem
	// 从最近文件列表中移除指定文件
	RemoveRecentFile(path string) error
	// 从最近文件夹列表中移除指定文件夹
	RemoveRecentFolder(path string) error
	// 清空最近文件列表
	ClearRecentFiles() error
	// 清空最近文件夹列表
	ClearRecentFolders() error
}

// IGitService Git 服务接口
type IGitService interface {
	// 打开 Git 仓库
	OpenRepository(path string) (*RepoInfo, error)
	// 获取仓库状态
	GetGitStatus(path string) (*GitStatus, error)
	// 提交更改
	GitCommit(path, message string) (string, error)
	// 获取分支信息
	GitGetBranches(path string) (*BranchInfo, error)
	// 获取提交日志
	GitGetLog(path string, maxCommits int) ([]CommitInfo, error)
}

// IAppService 应用主服务接口（组合所有服务）
type IAppService interface {
	// 基础问候方法
	Greet(name string) string
	
	// 文件系统服务
	IFileSystemService
	
	// Git 服务
	IGitService
	
	// 配置服务
	IConfigService
}

// ==================== 依赖注入容器 ====================

// ServiceContainer 服务容器，管理所有服务实例
type ServiceContainer struct {
	FileSystem IFileSystemService
	Git        IGitService
	Config     IConfigService
	App        IAppService
}

// NewServiceContainer 创建服务容器
func NewServiceContainer() *ServiceContainer {
	container := &ServiceContainer{}
	
	// 初始化具体服务实现
	fileSystemService := NewFileSystemService()
	gitService := NewGitService()
	configService := NewConfigManager()
	appService := NewAppService(fileSystemService, gitService, configService)
	
	// 注入到容器
	container.FileSystem = fileSystemService
	container.Git = gitService
	container.Config = configService
	container.App = appService
	
	return container
}

// ==================== 适配器模式 - Wails v2 兼容层 ====================

// WailsV2Adapter Wails v2 适配器
// 这个结构体实现了 Wails v2 需要的方法签名
// 当迁移到 v3 时，只需修改此适配器，业务逻辑无需改动
type WailsV2Adapter struct {
	ctx     context.Context
	services *ServiceContainer
}

// NewWailsV2Adapter 创建 Wails v2 适配器
func NewWailsV2Adapter() *WailsV2Adapter {
	return &WailsV2Adapter{
		services: NewServiceContainer(),
	}
}

// Startup Wails v2 启动钩子
func (w *WailsV2Adapter) Startup(ctx context.Context) {
	w.ctx = ctx
	
	// 将 context 注入到需要 runtime API 的服务中
	if fs, ok := w.services.FileSystem.(*FileSystemService); ok {
		fs.SetContext(ctx)
	}
}

// ==================== 暴露给前端的方法（通过 Wails Bridge）====================

// Greet 问候方法（保持向后兼容）
func (w *WailsV2Adapter) Greet(name string) string {
	return w.services.App.Greet(name)
}

// ReadFile 读取文件（代理到服务层）
func (w *WailsV2Adapter) ReadFile(path string) (string, error) {
	return w.services.App.ReadFile(path)
}

// WriteFile 写入文件
func (w *WailsV2Adapter) WriteFile(path string, content string) error {
	return w.services.App.WriteFile(path, content)
}

// ListDir 列出目录
func (w *WailsV2Adapter) ListDir(path string) ([]FileInfo, error) {
	return w.services.App.ListDir(path)
}

// GetProjectRoot 获取项目根目录
func (w *WailsV2Adapter) GetProjectRoot() string {
	return w.services.App.GetProjectRoot()
}

// OpenFolderDialog 打开文件夹选择对话框
func (w *WailsV2Adapter) OpenFolderDialog() (string, error) {
	return w.services.App.OpenFolderDialog()
}

// OpenFileDialog 打开文件选择对话框
func (w *WailsV2Adapter) OpenFileDialog() (string, error) {
	return w.services.App.OpenFileDialog()
}

// SaveFileDialog 保存文件对话框
func (w *WailsV2Adapter) SaveFileDialog() (string, error) {
	return w.services.App.SaveFileDialog()
}

// SetProjectRoot 设置项目根目录
func (w *WailsV2Adapter) SetProjectRoot(path string) error {
	return w.services.App.SetProjectRoot(path)
}

// CreateFile 创建新文件
func (w *WailsV2Adapter) CreateFile(path string) error {
	return w.services.App.CreateFile(path)
}

// CreateDirectory 创建新目录
func (w *WailsV2Adapter) CreateDirectory(path string) error {
	return w.services.App.CreateDirectory(path)
}

// DeleteFileOrDirectory 删除文件或目录
func (w *WailsV2Adapter) DeleteFileOrDirectory(path string) error {
	return w.services.App.DeleteFileOrDirectory(path)
}

// RenameFileOrDirectory 重命名文件或目录
func (w *WailsV2Adapter) RenameFileOrDirectory(oldPath, newPath string) error {
	return w.services.App.RenameFileOrDirectory(oldPath, newPath)
}

// MoveFileOrDirectory 移动文件或目录
func (w *WailsV2Adapter) MoveFileOrDirectory(sourcePath, targetPath string) error {
	return w.services.App.MoveFileOrDirectory(sourcePath, targetPath)
}

// GetFileStats 获取文件统计信息
func (w *WailsV2Adapter) GetFileStats(path string) (*FileInfo, error) {
	return w.services.App.GetFileStats(path)
}

// SearchFiles 搜索文件
func (w *WailsV2Adapter) SearchFiles(rootPath, keyword string, maxResults int) ([]FileInfo, error) {
	return w.services.App.SearchFiles(rootPath, keyword, maxResults)
}

// CopyFileOrDirectory 复制文件或目录
func (w *WailsV2Adapter) CopyFileOrDirectory(sourcePath, targetPath string) error {
	return w.services.App.CopyFileOrDirectory(sourcePath, targetPath)
}

// IsTextFile 判断是否为文本文件
func (w *WailsV2Adapter) IsTextFile(path string) bool {
	return w.services.App.IsTextFile(path)
}

// GetFileExtension 获取文件扩展名
func (w *WailsV2Adapter) GetFileExtension(path string) string {
	return w.services.App.GetFileExtension(path)
}

// GetDirectoryTree 获取目录树
func (w *WailsV2Adapter) GetDirectoryTree(path string, depth int) ([]FileInfo, error) {
	return w.services.App.GetDirectoryTree(path, depth)
}

// BackupFile 备份文件
func (w *WailsV2Adapter) BackupFile(path string) error {
	return w.services.App.BackupFile(path)
}

// TouchFile Touch 文件
func (w *WailsV2Adapter) TouchFile(path string) error {
	return w.services.App.TouchFile(path)
}

// OpenRepository 打开 Git 仓库
func (w *WailsV2Adapter) OpenRepository(path string) (*RepoInfo, error) {
	return w.services.App.OpenRepository(path)
}

// GetGitStatus 获取 Git 状态
func (w *WailsV2Adapter) GetGitStatus(path string) (*GitStatus, error) {
	return w.services.App.GetGitStatus(path)
}

// GitCommit 提交 Git 更改
func (w *WailsV2Adapter) GitCommit(path, message string) (string, error) {
	return w.services.App.GitCommit(path, message)
}

// GitGetBranches 获取 Git 分支
func (w *WailsV2Adapter) GitGetBranches(path string) (*BranchInfo, error) {
	return w.services.App.GitGetBranches(path)
}

// GitGetLog 获取 Git 日志
func (w *WailsV2Adapter) GitGetLog(path string, maxCommits int) ([]CommitInfo, error) {
	return w.services.App.GitGetLog(path, maxCommits)
}

// ==================== 配置服务方法 ====================

// AddRecentFile 添加最近打开的文件
func (w *WailsV2Adapter) AddRecentFile(path string) error {
	return w.services.App.AddRecentFile(path)
}

// AddRecentFolder 添加最近打开的文件夹
func (w *WailsV2Adapter) AddRecentFolder(path string) error {
	return w.services.App.AddRecentFolder(path)
}

// GetRecentFiles 获取最近打开的文件列表
func (w *WailsV2Adapter) GetRecentFiles() []RecentItem {
	return w.services.App.GetRecentFiles()
}

// GetRecentFolders 获取最近打开的文件夹列表
func (w *WailsV2Adapter) GetRecentFolders() []RecentItem {
	return w.services.App.GetRecentFolders()
}

// RemoveRecentFile 从最近文件列表中移除指定文件
func (w *WailsV2Adapter) RemoveRecentFile(path string) error {
	return w.services.App.RemoveRecentFile(path)
}

// RemoveRecentFolder 从最近文件夹列表中移除指定文件夹
func (w *WailsV2Adapter) RemoveRecentFolder(path string) error {
	return w.services.App.RemoveRecentFolder(path)
}

// ClearRecentFiles 清空最近文件列表
func (w *WailsV2Adapter) ClearRecentFiles() error {
	return w.services.App.ClearRecentFiles()
}

// ClearRecentFolders 清空最近文件夹列表
func (w *WailsV2Adapter) ClearRecentFolders() error {
	return w.services.App.ClearRecentFolders()
}
