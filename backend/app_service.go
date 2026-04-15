package backend

// AppService 应用主服务，组合所有子服务
type AppService struct {
	fileSystem IFileSystemService
	git        IGitService
	config     IConfigService
}

// NewAppService 创建应用服务
func NewAppService(fs IFileSystemService, git IGitService, config IConfigService) *AppService {
	return &AppService{
		fileSystem: fs,
		git:        git,
		config:     config,
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
