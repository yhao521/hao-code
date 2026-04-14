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
}

// ==================== 依赖注入容器 ====================

// ServiceContainer 服务容器，管理所有服务实例
type ServiceContainer struct {
	FileSystem IFileSystemService
	Git        IGitService
	App        IAppService
}

// NewServiceContainer 创建服务容器
func NewServiceContainer() *ServiceContainer {
	container := &ServiceContainer{}
	
	// 初始化具体服务实现
	fileSystemService := NewFileSystemService()
	gitService := NewGitService()
	appService := NewAppService(fileSystemService, gitService)
	
	// 注入到容器
	container.FileSystem = fileSystemService
	container.Git = gitService
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
