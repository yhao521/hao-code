package backend

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sergi/go-diff/diffmatchpatch"
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

// GetGitGraph 获取 Git 分支图谱数据
func (g *GitService) GetGitGraph(path string, maxCommits int) ([]GitGraphNode, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, err
	}
	defer commitIter.Close()

	var nodes []GitGraphNode
	count := 0

	// 简单的颜色分配逻辑
	colors := []string{"#E57373", "#64B5F6", "#81C784", "#FFD54F", "#BA68C8"}
	branchColors := make(map[string]string)
	colorIdx := 0

	err = commitIter.ForEach(func(c *object.Commit) error {
		if count >= maxCommits {
			return nil
		}

		// 获取该提交所在的分支
		var branches []string
		refs, _ := repo.References()
		refs.ForEach(func(r *plumbing.Reference) error {
			if r.Hash() == c.Hash && !r.Name().IsRemote() && r.Name().Short() != "HEAD" {
				branches = append(branches, r.Name().Short())
			}
			return nil
		})

		// 分配颜色
		mainBranch := "master"
		if len(branches) > 0 {
			mainBranch = branches[0]
		}
		if _, ok := branchColors[mainBranch]; !ok {
			branchColors[mainBranch] = colors[colorIdx%len(colors)]
			colorIdx++
		}

		var parents []string
		for _, p := range c.ParentHashes {
			parents = append(parents, p.String())
		}

		nodes = append(nodes, GitGraphNode{
			Hash:      c.Hash.String(),
			ShortHash: c.Hash.String()[:7],
			Message:   strings.TrimSpace(c.Message),
			Author:    c.Author.Name,
			Timestamp: c.Author.When.Unix(),
			Branches:  branches,
			Parents:   parents,
			Color:     branchColors[mainBranch],
		})

		count++
		return nil
	})

	return nodes, err
}

// GetFileDiff 获取文件差异
func (g *GitService) GetFileDiff(path, filePath string) (*FileDiff, error) {
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

	fileStatus := status.File(filePath)
	if fileStatus == nil {
		return nil, fmt.Errorf("file %s not found in git status", filePath)
	}

	var oldContent, newContent string
	var statusStr string

	// 获取工作区内容 (NewContent)
	fullPath := filepath.Join(path, filePath)
	data, err := os.ReadFile(fullPath)
	if err == nil {
		newContent = string(data)
	}

	// 获取暂存区或 HEAD 内容 (OldContent)
	if fileStatus.Staging != git.Unmodified {
		// 如果已暂存，对比暂存区与 HEAD
		statusStr = "staged"
		headRef, err := repo.Head()
		if err == nil {
			commit, err := repo.CommitObject(headRef.Hash())
			if err == nil {
				tree, err := commit.Tree()
				if err == nil {
					file, err := tree.File(filePath)
					if err == nil {
						oldContent, _ = file.Contents()
					}
				}
			}
		}
	} else {
		// 如果未暂存，对比工作区与暂存区/HEAD
		statusStr = "modified"
		// 尝试从暂存区获取
		index, _ := repo.Storer.Index()
		if index != nil {
			for _, entry := range index.Entries {
				if entry.Name == filePath {
					blob, _ := repo.BlobObject(entry.Hash)
					if blob != nil {
						reader, _ := blob.Reader()
						if reader != nil {
							bytes, _ := io.ReadAll(reader)
							oldContent = string(bytes)
						}
					}
					break
				}
			}
		}
		// 如果暂存区也没有，说明是新增文件或者从未提交过
		if oldContent == "" && fileStatus.Worktree == git.Added {
			statusStr = "added"
		}
	}

	switch fileStatus.Worktree {
	case git.Deleted:
		statusStr = "deleted"
		newContent = ""
	}

	return &FileDiff{
		Path:       filePath,
		OldContent: oldContent,
		NewContent: newContent,
		Status:     statusStr,
		Lines:      generateLineDiffs(oldContent, newContent),
	}, nil
}

// generateLineDiffs 使用 diff-match-patch 生成行级差异
func generateLineDiffs(oldText, newText string) []DiffLine {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(oldText, newText, false)

	var lines []DiffLine
	oldLineNum := 1
	newLineNum := 1

	for _, diff := range diffs {
		text := diff.Text
		// 处理换行符，按行分割
		textLines := strings.Split(text, "\n")
		// 如果最后一个元素是空字符串（因为以换行符结尾），则移除
		if len(textLines) > 0 && textLines[len(textLines)-1] == "" {
			textLines = textLines[:len(textLines)-1]
		}

		switch diff.Type {
		case diffmatchpatch.DiffEqual:
			for _, line := range textLines {
				lines = append(lines, DiffLine{
					Type:    "unchanged",
					Content: line,
					OldNum:  oldLineNum,
					NewNum:  newLineNum,
				})
				oldLineNum++
				newLineNum++
			}
		case diffmatchpatch.DiffInsert:
			for _, line := range textLines {
				lines = append(lines, DiffLine{
					Type:    "added",
					Content: line,
					OldNum:  0,
					NewNum:  newLineNum,
				})
				newLineNum++
			}
		case diffmatchpatch.DiffDelete:
			for _, line := range textLines {
				lines = append(lines, DiffLine{
					Type:    "deleted",
					Content: line,
					OldNum:  oldLineNum,
					NewNum:  0,
				})
				oldLineNum++
			}
		}
	}
	return lines
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

// StageSelectedRanges 暂存文件的指定行范围
func (g *GitService) StageSelectedRanges(path, filePath string, ranges []LineRange) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	// 1. 获取文件当前内容
	fullPath := filepath.Join(path, filePath)
	_, err = os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	// 2. 简单实现：先重置该文件，然后根据 ranges 生成 patch 并应用
	// 注意：由于 go-git 对部分暂存支持较弱，这里先调用 Add 作为基础
	_, err = worktree.Add(filePath)
	return err
}

// UnstageFile 取消暂存文件
func (g *GitService) UnstageFile(path, filePath string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	// 使用 Reset 命令将文件从暂存区移回工作区
	// git reset HEAD <file>
	return worktree.Reset(&git.ResetOptions{
		Mode:  git.MixedReset,
		Files: []string{filePath},
	})
}

// GetFileBlame 获取文件每一行的 Blame 信息
func (g *GitService) GetFileBlame(path, filePath string) ([]BlameInfo, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	// 获取 HEAD 中的文件内容
	headRef, _ := repo.Head()
	if headRef == nil {
		return nil, fmt.Errorf("no head found")
	}

	commit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	file, err := tree.File(filePath)
	if err != nil {
		return nil, err
	}

	// go-git 不直接支持 blame，我们需要通过遍历日志来模拟或调用 git 命令
	// 为了保持纯 Go 实现，我们这里简化处理：返回文件最后一次修改的提交信息
	// 真正的 blame 需要复杂的 diff 算法，这里先提供一个基础版本
	var blames []BlameInfo
	content, _ := file.Contents()
	lines := strings.Split(content, "\n")

	// 获取该文件的最新一次提交
	logIter, _ := repo.Log(&git.LogOptions{FileName: &filePath})
	latestCommit, _ := logIter.Next()

	if latestCommit != nil {
		for i := range lines {
			blames = append(blames, BlameInfo{
				Line:      i + 1,
				Hash:      latestCommit.Hash.String()[:7],
				Author:    latestCommit.Author.Name,
				Timestamp: latestCommit.Author.When.Format("2006-01-02"),
				Message:   strings.TrimSpace(latestCommit.Message),
			})
		}
	}

	return blames, nil
}

// GetFileHistory 获取文件的历史记录
func (g *GitService) GetFileHistory(path, filePath string) ([]CommitInfo, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	logIter, err := repo.Log(&git.LogOptions{FileName: &filePath})
	if err != nil {
		return nil, err
	}

	var commits []CommitInfo
	err = logIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, CommitInfo{
			Hash:      c.Hash.String(),
			ShortHash: c.Hash.String()[:7],
			Message:   strings.TrimSpace(c.Message),
			Author:    c.Author.Name,
			Email:     c.Author.Email,
			Timestamp: c.Author.When.Format("2006-01-02 15:04:05"),
		})
		return nil
	})

	return commits, err
}
