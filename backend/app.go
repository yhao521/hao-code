package backend

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ==================== 高级 Git 操作（使用 exec.Command）====================
// 这些方法是对 WailsV2Adapter 的扩展（WailsV2Adapter 在 interfaces.go 中定义）

// GitRebase 执行 rebase 操作
func (a *WailsV2Adapter) GitRebase(path, upstream string) (string, error) {
	return a.execGitCommand(path, "rebase", upstream)
}

// GitCherryPick 执行 cherry-pick 操作
func (a *WailsV2Adapter) GitCherryPick(path, commit string) (string, error) {
	return a.execGitCommand(path, "cherry-pick", commit)
}

// GitMerge 执行 merge 操作
func (a *WailsV2Adapter) GitMerge(path, branch string) (string, error) {
	return a.execGitCommand(path, "merge", branch)
}

// GitReset 执行 reset 操作
func (a *WailsV2Adapter) GitReset(path, mode, target string) (string, error) {
	return a.execGitCommand(path, "reset", fmt.Sprintf("--%s", mode), target)
}

// GitStash 执行 stash 操作
func (a *WailsV2Adapter) GitStash(path, action string, message string) (string, error) {
	args := []string{"stash"}
	if action == "save" && message != "" {
		args = append(args, "save", message)
	} else if action != "" {
		args = append(args, action)
	}
	return a.execGitCommand(path, args...)
}

// GitCheckout 执行 checkout 操作
func (a *WailsV2Adapter) GitCheckout(path, branch string) (string, error) {
	return a.execGitCommand(path, "checkout", branch)
}

// GitPull 执行 pull 操作
func (a *WailsV2Adapter) GitPull(path, remote, branch string) (string, error) {
	return a.execGitCommand(path, "pull", remote, branch)
}

// GitPush 执行 push 操作
func (a *WailsV2Adapter) GitPush(path, remote, branch string) (string, error) {
	return a.execGitCommand(path, "push", remote, branch)
}

// execGitCommand 执行 git 命令并返回输出
func (a *WailsV2Adapter) execGitCommand(path string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errorMsg := stderr.String()
		if errorMsg == "" {
			errorMsg = err.Error()
		}
		return "", fmt.Errorf("git command failed: %s", errorMsg)
	}

	return strings.TrimSpace(stdout.String()), nil
}

// ==================== App 结构体（Wails 绑定入口）====================

// App struct - Wails 绑定的主应用结构
type App struct {
	ctx     context.Context
	adapter *WailsV2Adapter
}

// NewApp 创建 App 实例（Wails 启动时调用）
func NewApp() *App {
	return &App{
		adapter: NewWailsV2Adapter(),
	}
}

// Startup Wails 启动钩子（导出）
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.adapter.Startup(ctx)
}

// ==================== 暴露给前端的方法 ====================

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return a.adapter.Greet(name)
}

// ReadFile reads content from a file
func (a *App) ReadFile(path string) (string, error) {
	return a.adapter.ReadFile(path)
}

// WriteFile writes content to a file
func (a *App) WriteFile(path string, content string) error {
	return a.adapter.WriteFile(path, content)
}

// ListDir lists files in a directory
func (a *App) ListDir(path string) ([]FileInfo, error) {
	return a.adapter.ListDir(path)
}

// GetProjectRoot returns the project root directory
func (a *App) GetProjectRoot() string {
	return a.adapter.GetProjectRoot()
}

// OpenFolderDialog 打开文件夹选择对话框
func (a *App) OpenFolderDialog() (string, error) {
	return a.adapter.OpenFolderDialog()
}

// OpenFileDialog 打开文件选择对话框
func (a *App) OpenFileDialog() (string, error) {
	return a.adapter.OpenFileDialog()
}

// SaveFileDialog 保存文件对话框
func (a *App) SaveFileDialog() (string, error) {
	return a.adapter.SaveFileDialog()
}

// SetProjectRoot 设置项目根目录
func (a *App) SetProjectRoot(path string) error {
	return a.adapter.SetProjectRoot(path)
}

// CreateFile 创建新文件
func (a *App) CreateFile(path string) error {
	return a.adapter.CreateFile(path)
}

// CreateDirectory 创建新目录
func (a *App) CreateDirectory(path string) error {
	return a.adapter.CreateDirectory(path)
}

// DeleteFileOrDirectory 删除文件或目录
func (a *App) DeleteFileOrDirectory(path string) error {
	return a.adapter.DeleteFileOrDirectory(path)
}

// RenameFileOrDirectory 重命名文件或目录
func (a *App) RenameFileOrDirectory(oldPath, newPath string) error {
	return a.adapter.RenameFileOrDirectory(oldPath, newPath)
}

// MoveFileOrDirectory 移动文件或目录
func (a *App) MoveFileOrDirectory(sourcePath, targetPath string) error {
	return a.adapter.MoveFileOrDirectory(sourcePath, targetPath)
}

// GetFileStats 获取文件统计信息
func (a *App) GetFileStats(path string) (*FileInfo, error) {
	return a.adapter.GetFileStats(path)
}

// SearchFiles 搜索文件
func (a *App) SearchFiles(rootPath, keyword string, maxResults int) ([]FileInfo, error) {
	return a.adapter.SearchFiles(rootPath, keyword, maxResults)
}

// CopyFileOrDirectory 复制文件或目录
func (a *App) CopyFileOrDirectory(sourcePath, targetPath string) error {
	return a.adapter.CopyFileOrDirectory(sourcePath, targetPath)
}

// IsTextFile 判断是否为文本文件
func (a *App) IsTextFile(path string) bool {
	return a.adapter.IsTextFile(path)
}

// GetFileExtension 获取文件扩展名
func (a *App) GetFileExtension(path string) string {
	return a.adapter.GetFileExtension(path)
}

// GetDirectoryTree 获取目录树
func (a *App) GetDirectoryTree(path string, depth int) ([]FileInfo, error) {
	return a.adapter.GetDirectoryTree(path, depth)
}

// BackupFile 备份文件
func (a *App) BackupFile(path string) error {
	return a.adapter.BackupFile(path)
}

// TouchFile Touch 文件
func (a *App) TouchFile(path string) error {
	return a.adapter.TouchFile(path)
}

// OpenRepository opens a git repository at the given path
func (a *App) OpenRepository(path string) (*RepoInfo, error) {
	return a.adapter.OpenRepository(path)
}

// GetGitStatus gets the current status of the repository
func (a *App) GetGitStatus(path string) (*GitStatus, error) {
	return a.adapter.GetGitStatus(path)
}

// GitCommit creates a new commit
func (a *App) GitCommit(path, message string) (string, error) {
	return a.adapter.GitCommit(path, message)
}

// GitGetBranches gets all branches
func (a *App) GitGetBranches(path string) (*BranchInfo, error) {
	return a.adapter.GitGetBranches(path)
}

// GitGetLog gets commit log
func (a *App) GitGetLog(path string, maxCommits int) ([]CommitInfo, error) {
	return a.adapter.GitGetLog(path, maxCommits)
}

// ==================== 高级 Git 操作 ====================

// GitRebase 执行 rebase 操作
func (a *App) GitRebase(path, upstream string) (string, error) {
	return a.adapter.GitRebase(path, upstream)
}

// GitCherryPick 执行 cherry-pick 操作
func (a *App) GitCherryPick(path, commit string) (string, error) {
	return a.adapter.GitCherryPick(path, commit)
}

// GitMerge 执行 merge 操作
func (a *App) GitMerge(path, branch string) (string, error) {
	return a.adapter.GitMerge(path, branch)
}

// GitReset 执行 reset 操作
func (a *App) GitReset(path, mode, target string) (string, error) {
	return a.adapter.GitReset(path, mode, target)
}

// GitStash 执行 stash 操作
func (a *App) GitStash(path, action string, message string) (string, error) {
	return a.adapter.GitStash(path, action, message)
}

// GitCheckout 执行 checkout 操作
func (a *App) GitCheckout(path, branch string) (string, error) {
	return a.adapter.GitCheckout(path, branch)
}

// GitPull 执行 pull 操作
func (a *App) GitPull(path, remote, branch string) (string, error) {
	return a.adapter.GitPull(path, remote, branch)
}

// GitPush 执行 push 操作
func (a *App) GitPush(path, remote, branch string) (string, error) {
	return a.adapter.GitPush(path, remote, branch)
}