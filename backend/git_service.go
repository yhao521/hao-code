package main

import (
	"fmt"
	"strings"
	"time"

	git "github.com/libgit2/git2go/v34"
)

// GitService Git 服务实现
type GitService struct{}

// NewGitService 创建 Git 服务
func NewGitService() *GitService {
	return &GitService{}
}

// OpenRepository 打开 Git 仓库
func (g *GitService) OpenRepository(path string) (*RepoInfo, error) {
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

// GetGitStatus 获取仓库状态
func (g *GitService) GetGitStatus(path string) (*GitStatus, error) {
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

		change := g.parseStatusEntry(entry)
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

// GitCommit 创建提交
func (g *GitService) GitCommit(path, message string) (string, error) {
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
	if head.Target() != nil {
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
		sig,
		sig,
		message,
		tree,
		parents...,
	)
	if err != nil {
		return "", err
	}

	return commitId.String(), nil
}

// GitGetBranches 获取分支信息
func (g *GitService) GitGetBranches(path string) (*BranchInfo, error) {
	repo, err := git.OpenRepository(path)
	if err != nil {
		return nil, err
	}
	defer repo.Free()

	var localBranches []string
	var remoteBranches []string

	// 创建引用迭代器
	iter, err := repo.NewReferenceIterator()
	if err != nil {
		return nil, err
	}

	// 遍历所有引用
	for {
		ref, err := iter.Next()
		if err != nil {
			break // 迭代完成
		}

		if strings.HasPrefix(ref.Name(), "refs/heads/") {
			name := strings.TrimPrefix(ref.Name(), "refs/heads/")
			localBranches = append(localBranches, name)
		} else if strings.HasPrefix(ref.Name(), "refs/remotes/") {
			name := strings.TrimPrefix(ref.Name(), "refs/remotes/")
			remoteBranches = append(remoteBranches, name)
		}
	}

	// 获取当前分支
	currentBranch := ""
	head, err := repo.Head()
	if err == nil && head.Target() != nil {
		currentBranch = head.Shorthand()
	}

	return &BranchInfo{
		Local:          localBranches,
		Remote:         remoteBranches,
		CurrentBranch:  currentBranch,
	}, nil
}

// GitGetLog 获取提交日志
func (g *GitService) GitGetLog(path string, maxCommits int) ([]CommitInfo, error) {
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

// parseStatusEntry 解析 Git 状态条目
func (g *GitService) parseStatusEntry(entry git.StatusEntry) *Change {
	var status string
	var path string

	switch {
	case entry.Status&git.StatusIndexNew != 0 || entry.Status&git.StatusWtNew != 0:
		status = "added"
		path = entry.HeadToIndex.NewFile.Path
	case entry.Status&git.StatusIndexDeleted != 0 || entry.Status&git.StatusWtDeleted != 0:
		status = "deleted"
		path = entry.HeadToIndex.OldFile.Path
	case entry.Status&git.StatusIndexModified != 0 || entry.Status&git.StatusWtModified != 0:
		status = "modified"
		if entry.Status&git.StatusIndexModified != 0 {
			path = entry.HeadToIndex.OldFile.Path
		} else {
			path = entry.IndexToWorkdir.OldFile.Path
		}
	default:
		return nil
	}

	return &Change{
		Path:   path,
		Status: status,
	}
}
