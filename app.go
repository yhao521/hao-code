package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// WailsV2Adapter handles the business logic
type WailsV2Adapter struct {
	ctx context.Context
}

// NewWailsV2Adapter creates a new adapter instance
func NewWailsV2Adapter() *WailsV2Adapter {
	return &WailsV2Adapter{}
}

// Startup initializes the adapter with context
func (a *WailsV2Adapter) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *WailsV2Adapter) Greet(name string) string {
	return fmt.Sprintf("Hello %s, Welcome to Hao-Code Editor!", name)
}

// ==================== 文件系统 API ====================

// ReadFile reads content from a file
func (a *WailsV2Adapter) ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WriteFile writes content to a file
func (a *WailsV2Adapter) WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// ListDir lists files in a directory
func (a *WailsV2Adapter) ListDir(path string) ([]FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	for _, entry := range entries {
		// 跳过隐藏文件和 node_modules
		if strings.HasPrefix(entry.Name(), ".") || entry.Name() == "node_modules" {
			continue
		}

		info := FileInfo{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			Path:  filepath.Join(path, entry.Name()),
		}

		if !entry.IsDir() {
			fileInfo, _ := entry.Info()
			info.Size = fileInfo.Size()
			info.ModTime = fileInfo.ModTime().Unix()
		}

		files = append(files, info)
	}

	return files, nil
}

// GetProjectRoot returns the project root directory
func (a *WailsV2Adapter) GetProjectRoot() string {
	dir, _ := os.Getwd()
	return dir
}

// ==================== Git API ====================

// OpenRepository opens a git repository at the given path
func (a *WailsV2Adapter) OpenRepository(path string) (*RepoInfo, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %v", err)
	}

	// 获取当前分支
	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get head: %v", err)
	}

	branchName := head.Name().Short()

	return &RepoInfo{
		Path:          path,
		CurrentBranch: branchName,
		IsRepository:  true,
	}, nil
}

// GetGitStatus gets the current status of the repository
func (a *WailsV2Adapter) GetGitStatus(path string) (*GitStatus, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, err
	}

	var stagedChanges []Change
	var changes []Change

	for file, s := range status {
		change := a.parseFileStatus(s, file)
		if change != nil {
			// go-git 中 Staging 表示暂存区状态，Worktree 表示工作区状态
			if s.Staging != git.Unmodified {
				stagedChanges = append(stagedChanges, *change)
			} else {
				changes = append(changes, *change)
			}
		}
	}

	return &GitStatus{
		StagedChanges: stagedChanges,
		Changes:       changes,
	}, nil
}

// GitCommit creates a new commit
func (a *WailsV2Adapter) GitCommit(path, message string) (string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	// 获取用户信息（用于提交）
	config, err := repo.Config()
	if err != nil {
		return "", fmt.Errorf("failed to get config: %v", err)
	}

	author := &object.Signature{
		Name:  config.User.Name,
		Email: config.User.Email,
		When:  time.Now(),
	}

	// 添加所有更改到暂存区
	err = worktree.AddGlob(".")
	if err != nil {
		return "", err
	}

	// 创建提交
	commitHash, err := worktree.Commit(message, &git.CommitOptions{
		Author: author,
	})
	if err != nil {
		return "", err
	}

	return commitHash.String(), nil
}

// GitGetBranches gets all branches
func (a *WailsV2Adapter) GitGetBranches(path string) (*BranchInfo, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	// 获取本地分支
	localBranches := []string{}
	branchIter, err := repo.Branches()
	if err != nil {
		return nil, err
	}

	err = branchIter.ForEach(func(ref *plumbing.Reference) error {
		name := ref.Name().Short()
		localBranches = append(localBranches, name)
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 获取远程分支
	remoteBranches := []string{}
	refs, err := repo.References()
	if err != nil {
		return nil, err
	}

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsRemote() && ref.Name().Short() != "HEAD" {
			remoteBranches = append(remoteBranches, ref.Name().Short())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 获取当前分支
	currentBranch := ""
	head, err := repo.Head()
	if err == nil {
		currentBranch = head.Name().Short()
	}

	return &BranchInfo{
		Local:         localBranches,
		Remote:        remoteBranches,
		CurrentBranch: currentBranch,
	}, nil
}

// GitGetLog gets commit log
func (a *WailsV2Adapter) GitGetLog(path string, maxCommits int) ([]CommitInfo, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	// 获取 HEAD 引用
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	// 创建提交迭代器
	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}
	defer commitIter.Close()

	var commits []CommitInfo
	count := 0

	err = commitIter.ForEach(func(c *object.Commit) error {
		if count >= maxCommits {
			return nil // 停止迭代
		}

		commits = append(commits, CommitInfo{
			Hash:      c.Hash.String(),
			ShortHash: c.Hash.String()[:7],
			Message:   strings.TrimSpace(c.Message),
			Author:    c.Author.Name,
			Email:     c.Author.Email,
			Timestamp: c.Author.When.Format("2006-01-02 15:04:05"),
		})

		count++
		return nil
	})
	if err != nil {
		return nil, err
	}

	return commits, nil
}

// ==================== 高级 Git 操作（使用 exec.Command）====================

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

// FileInfo represents information about a file
type FileInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Size    int64  `json:"size,omitempty"`
	IsDir   bool   `json:"isDir"`
	ModTime int64  `json:"modTime,omitempty"`
}

// RepoInfo represents information about a git repository
type RepoInfo struct {
	Path           string `json:"path"`
	CurrentBranch  string `json:"currentBranch"`
	IsRepository   bool   `json:"isRepository"`
}

// GitStatus represents the status of a git repository
type GitStatus struct {
	StagedChanges []Change `json:"stagedChanges"`
	Changes       []Change `json:"changes"`
}

// Change represents a single change in a git repository
type Change struct {
	Path    string `json:"path"`
	Status  string `json:"status"`
	OldPath string `json:"oldPath,omitempty"`
}

// BranchInfo represents information about git branches
type BranchInfo struct {
	Local         []string `json:"local"`
	Remote        []string `json:"remote"`
	CurrentBranch string   `json:"currentBranch"`
}

// CommitInfo represents information about a git commit
type CommitInfo struct {
	Hash      string `json:"hash"`
	ShortHash string `json:"shortHash"`
	Message   string `json:"message"`
	Author    string `json:"author"`
	Email     string `json:"email"`
	Timestamp string `json:"timestamp"`
}

// Helper functions

func (a *WailsV2Adapter) parseFileStatus(s *git.FileStatus, file string) *Change {
	var status string

	switch s.Worktree {
	case git.Added:
		status = "added"
	case git.Deleted:
		status = "deleted"
	case git.Modified:
		status = "modified"
	case git.Renamed:
		status = "renamed"
	case git.Copied:
		status = "copied"
	default:
		if s.Staging == git.Added {
			status = "staged"
		} else {
			return nil
		}
	}

	return &Change{
		Path:   file,
		Status: status,
	}
}

// App struct - 保持向后兼容，但内部使用新的服务架构
type App struct {
	ctx     context.Context
	adapter *WailsV2Adapter
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		adapter: NewWailsV2Adapter(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.adapter.Startup(ctx)
}

// ==================== 暴露给前端的方法 ====================
// 这些方法直接代理到适配器层

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


