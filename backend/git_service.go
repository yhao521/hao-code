package backend

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GitService Git 服务实现
type GitService struct{}

// NewGitService 创建 Git 服务
func NewGitService() *GitService {
	return &GitService{}
}

// OpenRepository 打开 Git 仓库
func (g *GitService) OpenRepository(path string) (*RepoInfo, error) {
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

// GetGitStatus 获取仓库状态
func (g *GitService) GetGitStatus(path string) (*GitStatus, error) {
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
		change := g.parseFileStatus(s, file)
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

// GitCommit 创建提交
func (g *GitService) GitCommit(path, message string) (string, error) {
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

// GitGetBranches 获取分支信息
func (g *GitService) GitGetBranches(path string) (*BranchInfo, error) {
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

// GitGetLog 获取提交日志
func (g *GitService) GitGetLog(path string, maxCommits int) ([]CommitInfo, error) {
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

func (g *GitService) parseFileStatus(s *git.FileStatus, file string) *Change {
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




