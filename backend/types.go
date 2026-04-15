package backend

// ==================== 文件系统相关类型 ====================

// FileInfo 文件信息
type FileInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Size    int64  `json:"size,omitempty"`
	IsDir   bool   `json:"isDir"`
	ModTime int64  `json:"modTime,omitempty"`
}

// ==================== Git 相关类型 ====================

// RepoInfo Git 仓库信息
type RepoInfo struct {
	Path           string `json:"path"`
	CurrentBranch  string `json:"currentBranch"`
	IsRepository   bool   `json:"isRepository"`
}

// GitStatus Git 状态
type GitStatus struct {
	StagedChanges []Change `json:"stagedChanges"`
	Changes       []Change `json:"changes"`
}

// Change 文件变更
type Change struct {
	Path    string `json:"path"`
	Status  string `json:"status"`
	OldPath string `json:"oldPath,omitempty"`
}

// BranchInfo 分支信息
type BranchInfo struct {
	Local         []string `json:"local"`
	Remote        []string `json:"remote"`
	CurrentBranch string   `json:"currentBranch"`
}

// CommitInfo 提交信息
type CommitInfo struct {
	Hash      string `json:"hash"`
	ShortHash string `json:"shortHash"`
	Message   string `json:"message"`
	Author    string `json:"author"`
	Email     string `json:"email"`
	Timestamp string `json:"timestamp"`
}

// RecentItem 最近打开的项目（文件或文件夹）
type RecentItem struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	OpenedAt string `json:"openedAt"`
}
