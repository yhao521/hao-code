package main

// AppService 应用主服务，组合所有子服务
type AppService struct {
	fileSystem IFileSystemService
	git        IGitService
}

// NewAppService 创建应用服务
func NewAppService(fs IFileSystemService, git IGitService) *AppService {
	return &AppService{
		fileSystem: fs,
		git:        git,
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
