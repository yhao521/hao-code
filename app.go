package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	git "github.com/libgit2/git2go/v34"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, Welcome to Hao-Code Editor!", name)
}

// ==================== 文件系统 API ====================

// ReadFile reads content from a file
func (a *App) ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WriteFile writes content to a file
func (a *App) WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// ListDir lists files in a directory
func (a *App) ListDir(path string) ([]FileInfo, error) {
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
func (a *App) GetProjectRoot() string {
	dir, _ := os.Getwd()
	return dir
}

// FileInfo represents information about a file
type FileInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Size    int64  `json:"size,omitempty"`
	IsDir   bool   `json:"isDir"`
	ModTime int64  `json:"modTime,omitempty"`
}

// ==================== Git API ====================

// OpenRepository opens a git repository at the given path
func (a *App) OpenRepository(path string) (*RepoInfo, error) {
	repo, err := git.OpenRepository(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %v", err)
	}
	defer repo.Free()

	// 获取当前分支
	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get head: %v", err)
	}
	defer head.Free()

	branchName := head.Shorthand()

	return &RepoInfo{
		Path:           path,
		CurrentBranch:  branchName,
		IsRepository:   true,
	}, nil
}

// GetGitStatus gets the current status of the repository
func (a *App) GetGitStatus(path string) (*GitStatus, error) {
	repo, err := git.OpenRepository(path)
	if err != nil {
		return nil, err
	}
	defer repo.Free()

	statusList, err := repo.StatusList(&git.StatusOptions{
		Show:  git.StatusShowIndexAndWorkdir,
		Flags: git.StatusOptIncludeUntracked,
	})
	if err != nil {
		return nil, err
	}
	defer statusList.Free()

	var stagedChanges []Change
	var changes []Change

	entryCount, err := statusList.EntryCount()
	if err != nil {
		return nil, err
	}

	for i := 0; i < entryCount; i++ {
		entry, err := statusList.ByIndex(i)
		if err != nil {
			continue
		}

		change := a.parseStatusEntry(entry)
		if change != nil {
			if entry.Status&git.StatusIndexNew != 0 || 
			   entry.Status&git.StatusIndexModified != 0 ||
			   entry.Status&git.StatusIndexDeleted != 0 {
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
func (a *App) GitCommit(path, message string) (string, error) {
	repo, err := git.OpenRepository(path)
	if err != nil {
		return "", err
	}
	defer repo.Free()

	// 获取 HEAD
	head, err := repo.Head()
	if err != nil {
		return "", err
	}
	defer head.Free()

	// 获取签名
	sig, err := repo.DefaultSignature()
	if err != nil {
		return "", fmt.Errorf("failed to get signature: %v", err)
	}

	// 创建树
	index, err := repo.Index()
	if err != nil {
		return "", err
	}
	defer index.Free()

	treeId, err := index.WriteTree()
	if err != nil {
		return "", err
	}

	tree, err := repo.LookupTree(treeId)
	if err != nil {
		return "", err
	}
	defer tree.Free()

	// 获取父提交
	var parents []*git.Commit
	if !head.IsUnborn() {
		parentCommit, err := repo.LookupCommit(head.Target())
		if err != nil {
			return "", err
		}
		defer parentCommit.Free()
		parents = append(parents, parentCommit)
	}

	// 创建提交
	commitId, err := repo.CreateCommit(
		"HEAD",
		&sig,
		&sig,
		message,
		tree,
		len(parents) > 0,
		parents...,
	)
	if err != nil {
		return "", err
	}

	return commitId.String(), nil
}

// GitGetBranches gets all branches
func (a *App) GitGetBranches(path string) (*BranchInfo, error) {
	repo, err := git.OpenRepository(path)
	if err != nil {
		return nil, err
	}
	defer repo.Free()

	var localBranches []string
	var remoteBranches []string

	// 遍历本地分支
	err = repo.WalkReferences(git.ReferenceTypeDirect, func(ref *git.Reference) error {
		if strings.HasPrefix(ref.Name(), "refs/heads/") {
			name := strings.TrimPrefix(ref.Name(), "refs/heads/")
			localBranches = append(localBranches, name)
		} else if strings.HasPrefix(ref.Name(), "refs/remotes/") {
			name := strings.TrimPrefix(ref.Name(), "refs/remotes/")
			remoteBranches = append(remoteBranches, name)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 获取当前分支
	currentBranch := ""
	head, err := repo.Head()
	if err == nil && !head.IsUnborn() {
		currentBranch = head.Shorthand()
	}

	return &BranchInfo{
		Local:          localBranches,
		Remote:         remoteBranches,
		CurrentBranch:  currentBranch,
	}, nil
}

// GitGetLog gets commit log
func (a *App) GitGetLog(path string, maxCommits int) ([]CommitInfo, error) {
	repo, err := git.OpenRepository(path)
	if err != nil {
		return nil, err
	}
	defer repo.Free()

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}
	defer head.Free()

	walk, err := repo.Walk()
	if err != nil {
		return nil, err
	}
	defer walk.Free()

	err = walk.PushHead()
	if err != nil {
		return nil, err
	}

	var commits []CommitInfo
	count := 0

	walk.Iterate(func(commit *git.Commit) bool {
		if count >= maxCommits {
			return false
		}

		author := commit.Author()
		
		commits = append(commits, CommitInfo{
			Hash:      commit.Id().String(),
			ShortHash: commit.Id().String()[:7],
			Message:   commit.Message(),
			Author:    author.Name,
			Email:     author.Email,
			Timestamp: time.Unix(author.When.Unix(), 0).Format("2006-01-02 15:04:05"),
		})

		count++
		return true
	})

	return commits, nil
}

// Helper functions

func (a *App) parseStatusEntry(entry git.StatusEntry) *Change {
	var status string
	var path string

	switch {
	case entry.Status&git.StatusIndexNew != 0 || entry.Status&git.StatusWTNew != 0:
		status = "added"
		path = entry.HeadToIndex.NewFile.Path
	case entry.Status&git.StatusIndexDeleted != 0 || entry.Status&git.StatusWTDeleted != 0:
		status = "deleted"
		path = entry.HeadToIndex.OldFile.Path
	case entry.Status&git.StatusIndexModified != 0 || entry.Status&git.StatusWTModified != 0:
		status = "modified"
		if entry.Status&git.StatusIndexModified != 0 {
			path = entry.HeadToIndex.OldFile.Path
		} else {
			path = entry.WorkdirToIndex.OldFile.Path
		}
	default:
		return nil
	}

	return &Change{
		Path:   path,
		Status: status,
	}
}

// Git types

type RepoInfo struct {
	Path           string `json:"path"`
	CurrentBranch  string `json:"currentBranch"`
	IsRepository   bool   `json:"isRepository"`
}

type GitStatus struct {
	StagedChanges []Change `json:"stagedChanges"`
	Changes       []Change `json:"changes"`
}

type Change struct {
	Path    string `json:"path"`
	Status  string `json:"status"`
	OldPath string `json:"oldPath,omitempty"`
}

type BranchInfo struct {
	Local         []string `json:"local"`
	Remote        []string `json:"remote"`
	CurrentBranch string   `json:"currentBranch"`
}

type CommitInfo struct {
	Hash      string `json:"hash"`
	ShortHash string `json:"shortHash"`
	Message   string `json:"message"`
	Author    string `json:"author"`
	Email     string `json:"email"`
	Timestamp string `json:"timestamp"`
}
