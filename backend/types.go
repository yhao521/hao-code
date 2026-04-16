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
	Path          string `json:"path"`
	CurrentBranch string `json:"currentBranch"`
	IsRepository  bool   `json:"isRepository"`
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

// SearchResult 搜索结果
type SearchResult struct {
	FilePath    string `json:"filePath"`
	LineNumber  int    `json:"lineNumber"`
	LineContent string `json:"lineContent"`
}

// SearchOptions 搜索选项
type SearchOptions struct {
	RootPath       string `json:"rootPath"`
	Query          string `json:"query"`
	CaseSensitive  bool   `json:"caseSensitive"`
	MatchWholeWord bool   `json:"matchWholeWord"`
	UseRegex       bool   `json:"useRegex"`
	Exclude        string `json:"exclude"` // 排除的文件模式，如 "*.log,node_modules"
}

// GitGraphNode Git 图谱节点
type GitGraphNode struct {
	Hash      string   `json:"hash"`
	ShortHash string   `json:"shortHash"`
	Message   string   `json:"message"`
	Author    string   `json:"author"`
	Timestamp int64    `json:"timestamp"`
	Branches  []string `json:"branches"` // 该节点所在的分支名
	Parents   []string `json:"parents"`  // 父节点 Hash
	Color     string   `json:"color"`    // 分支颜色
}

// DiffLine 差异行信息
type DiffLine struct {
	Type    string `json:"type"` // added, deleted, modified, unchanged
	Content string `json:"content"`
	OldNum  int    `json:"oldNum"` // 原始文件行号
	NewNum  int    `json:"newNum"` // 新文件行号
}

// FileDiff 文件差异信息
type FileDiff struct {
	Path       string     `json:"path"`
	OldContent string     `json:"oldContent"`
	NewContent string     `json:"newContent"`
	Status     string     `json:"status"` // added, modified, deleted, renamed
	Lines      []DiffLine `json:"lines"`  // 详细的行级差异
}

// LineRange 行范围
type LineRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// BlameInfo 代码行责任信息
type BlameInfo struct {
	Line      int    `json:"line"`
	Hash      string `json:"hash"`
	Author    string `json:"author"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

// TaskItem 任务项
type TaskItem struct {
	Name        string `json:"name"`
	Command     string `json:"command"`
	Type        string `json:"type"` // npm, make, custom
	Description string `json:"description,omitempty"`
}
