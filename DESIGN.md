# Wails3 代码编辑器设计文档

## 1. 项目概述

### 1.1 项目名称
Hao-Code Editor

### 1.2 项目愿景
打造一个轻量级但功能强大的跨平台代码编辑器，对标 VSCode 的扩展能力和 IntelliJ IDEA 的 Git 集成体验。

### 1.3 核心技术栈
- **后端框架**: Wails v3 (Go)
- **前端框架**: Vue 3 + TypeScript
- **UI 组件库**: Naive UI
- **代码编辑核心**: Monaco Editor
- **状态管理**: Pinia
- **构建工具**: Vite
- **插件系统**: 基于 Go plugin + JavaScript SDK

---

## 2. 系统架构设计

### 2.1 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                      前端层 (Vue 3 + TypeScript)              │
├─────────────────────────────────────────────────────────────┤
│  Naive UI  │  Monaco Editor  │  插件运行时  │  Pinia  │  Git可视化  │
├─────────────────────────────────────────────────────────────┤
│                    Wails Bridge (IPC通信)                     │
├─────────────────────────────────────────────────────────────┤
│                      后端层 (Go)                              │
├──────────────────┬──────────────────┬───────────────────────┤
│   插件管理器      │   Git 引擎       │   文件系统             │
│   Plugin Manager │   Git Engine     │   File System         │
├──────────────────┼──────────────────┼───────────────────────┤
│   LSP 客户端     │   搜索索引        │   配置管理             │
│   LSP Client     │   Search Index   │   Config Manager      │
└──────────────────┴──────────────────┴───────────────────────┘
```

### 2.2 模块划分

#### 2.2.1 前端模块
1. **编辑器核心模块**
   - Monaco Editor 集成
   - 多标签页管理
   - 代码高亮与智能提示
   - 快捷键系统

2. **插件UI模块**
   - 插件市场界面
   - 插件管理面板
   - 插件贡献点渲染

3. **Git可视化模块**
   - 提交历史时间线
   - 分支图谱
   - Diff 对比视图
   - 冲突解决界面

4. **工作区模块**
   - 文件树浏览器
   - 快速打开 (Ctrl+P)
   - 面包屑导航

#### 2.2.2 后端模块
1. **插件引擎**
   - Go plugin 加载
   - 插件生命周期管理
   - 插件间通信

2. **Git 引擎**
   - libgit2 绑定
   - Git 操作封装
   - 历史记录分析

3. **文件系统**
   - 虚拟文件系统
   - 文件监听
   - 批量操作

4. **LSP 客户端**
   - 语言服务器协议
   - 诊断信息收集
   - 代码补全代理

---

## 3. 插件系统设计

### 3.1 插件架构

```
┌──────────────────────────────────────────────┐
│              插件类型                          │
├────────────┬─────────┬──────────┬────────────┤
│ 主题插件   │ 语言插件 │ 工具插件 │ UI插件     │
└────────────┴─────────┴──────────┴────────────┘

┌──────────────────────────────────────────────┐
│            插件贡献点 (Contribution Points)    │
├────────────┬─────────┬──────────┬────────────┤
│ commands   │ menus   │ keymaps  │ languages  │
├────────────┼─────────┼──────────┼────────────┤
│ themes     │ views   │ editors  │ git-hooks  │
└────────────┴─────────┴──────────┴────────────┘
```

### 3.2 插件 SDK

#### 3.2.1 JavaScript SDK
```typescript
interface HaoExtensionContext {
  // 订阅生命周期事件
  onActivate: () => void;
  onDeactivate: () => void;
  
  // 注册命令
  registerCommand: (name: string, handler: Function) => void;
  
  // 注册语言支持
  registerLanguageProvider: (language: string, provider: LanguageProvider) => void;
  
  // 访问工作区
  workspace: WorkspaceAPI;
  
  // 访问窗口
  window: WindowAPI;
  
  // 访问 Git
  git: GitAPI;
}
```

#### 3.2.2 Go Plugin 接口
```go
type Extension interface {
    Activate(ctx ExtensionContext) error
    Deactivate() error
    Metadata() ExtensionMetadata
}

type ExtensionMetadata struct {
    Name        string
    Version     string
    Description string
    Author      string
    Contributors map[string]interface{}
}
```

### 3.3 插件市场

#### 功能特性
- 在线浏览插件
- 评分与评论
- 版本管理
- 一键安装/卸载
- 依赖解析

#### 数据结构
```go
type PluginInfo struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Version     string   `json:"version"`
    Description string   `json:"description"`
    Author      string   `json:"author"`
    Tags        []string `json:"tags"`
    Downloads   int      `json:"downloads"`
    Rating      float64  `json:"rating"`
    Repository  string   `json:"repository"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

---

## 4. Git 功能设计 (对标 IDEA)

### 4.1 核心功能

#### 4.1.1 提交历史
- **时间线视图**: 线性展示提交历史
- **分支图谱**: 可视化展示分支合并关系
- **过滤搜索**: 按作者、日期、关键词过滤
- **无限滚动**: 性能优化的懒加载

#### 4.1.2 变更管理
- **本地修改**: 实时显示未暂存更改
- **暂存区**: 可视化选择要提交的文件
- **提交详情**: 完整的提交信息展示
- **回滚操作**: 支持任意提交的回滚

#### 4.1.3 分支操作
- **分支列表**: 本地和远程分支
- **创建分支**: 从任意提交创建
- **合并分支**: 可视化合并流程
- **冲突解决**: 三向合并编辑器

#### 4.1.4 Diff 对比
- **文件对比**: 并排或行内对比
- **目录对比**: 批量文件差异
- **提交对比**: 任意两个提交对比
- **工作区对比**: 与任意提交对比

### 4.2 Git 数据模型

```go
type Commit struct {
    Hash        string    `json:"hash"`
    ShortHash   string    `json:"short_hash"`
    Author      Signature `json:"author"`
    Committer   Signature `json:"committer"`
    Message     string    `json:"message"`
    Parents     []string  `json:"parents"`
    Timestamp   time.Time `json:"timestamp"`
    Stats       DiffStats `json:"stats"`
}

type Branch struct {
    Name       string   `json:"name"`
    FullName   string   `json:"full_name"`
    IsRemote   bool     `json:"is_remote"`
    IsCurrent  bool     `json:"is_current"`
    LastCommit *Commit  `json:"last_commit"`
    Ahead      int      `json:"ahead"`
    Behind     int      `json:"behind"`
}

type Change struct {
    Path       string     `json:"path"`
    Status     ChangeStatus `json:"status"`
    OldPath    string     `json:"old_path,omitempty"`
    Diff       *FileDiff  `json:"diff,omitempty"`
}

type MergeConflict struct {
    FilePath      string   `json:"file_path"`
    ConflictType  string   `json:"conflict_type"`
    Ancestor      string   `json:"ancestor"`
    Current       string   `json:"current"`
    Incoming      string   `json:"incoming"`
}
```

### 4.3 Git 工作流程

#### 4.3.1 提交流程
```
1. 检测文件变更 → 2. 显示变更列表 → 3. 选择暂存文件 
→ 4. 编写提交信息 → 5. 执行提交 → 6. 刷新视图
```

#### 4.3.2 合并流程
```
1. 选择目标分支 → 2. 预检查冲突 → 3. 执行合并 
→ 4. 显示冲突列表 → 5. 逐个解决 → 6. 完成合并
```

#### 4.3.3 Rebase 流程
```
1. 选择基准分支 → 2. 交互式选择提交 → 3. 处理冲突 
→ 4. 继续/跳过/中止 → 5. 完成变基
```

---

## 5. 编辑器核心功能

### 5.1 代码编辑

#### 功能清单
- ✅ 语法高亮 (100+ 语言)
- ✅ 智能补全 (IntelliSense)
- ✅ 错误诊断 (实时)
- ✅ 代码格式化
- ✅ 重构支持 (重命名、提取函数等)
- ✅ 代码折叠
- ✅ 括号匹配
- ✅ 多光标编辑
- ✅ 查找替换 (支持正则)

### 5.2 LSP 集成

```go
type LSPClient struct {
    ServerPath string
    RootURI    string
    Capabilities ServerCapabilities
    Diagnostics  map[string][]Diagnostic
}

type Diagnostic struct {
    Range    Range
    Severity DiagnosticSeverity
    Message  string
    Source   string
    Code     string
}
```

### 5.3 快捷键系统

| 快捷键 | 功能 |
|--------|------|
| Ctrl+P | 快速打开文件 |
| Ctrl+Shift+P | 命令面板 |
| Ctrl+B | 切换侧边栏 |
| Ctrl+` | 切换终端 |
| Ctrl+S | 保存文件 |
| Ctrl+Z | 撤销 |
| Ctrl+Y | 重做 |
| Ctrl+F | 查找 |
| Ctrl+H | 替换 |
| Ctrl+/ | 切换注释 |

---

## 6. 前端架构

### 6.1 技术选型对比

#### Vue 3 方案
```
优势:
- 学习曲线平缓
- Composition API 灵活
- 响应式系统高效
- 生态丰富 (Element Plus, Vite)

劣势:
- 大型项目规范性稍弱
- TypeScript 支持略逊于 React
```

#### React 18 方案
```
优势:
- 大型企业项目首选
- TypeScript 完美支持
- 生态最丰富
- 社区资源最多

劣势:
- 学习曲线较陡
- 样板代码较多
```

### 6.2 状态管理

#### Pinia Store (Vue)
```typescript
interface EditorState {
  activeEditor: string | null;
  editors: Record<string, EditorInstance>;
  tabs: Tab[];
  sidebarVisible: boolean;
  sidebarView: 'explorer' | 'search' | 'git' | 'extensions';
}

interface GitState {
  repository: Repository | null;
  branches: Branch[];
  currentBranch: string;
  changes: Change[];
  stagedChanges: Change[];
  commitMessage: string;
  isCommitting: boolean;
}

interface PluginState {
  installedPlugins: PluginInfo[];
  activatedPlugins: Set<string>;
  marketplace: PluginMarketplace;
}
```

### 6.3 组件结构

```
App.vue
├── TitleBar (标题栏)
├── ActivityBar (活动栏)
│   ├── ExplorerIcon
│   ├── SearchIcon
│   ├── GitIcon
│   └── ExtensionsIcon
├── SideBar (侧边栏)
│   ├── FileExplorer
│   ├── SearchPanel
│   ├── GitPanel
│   │   ├── ChangesView
│   │   ├── HistoryView
│   │   └── BranchesView
│   └── ExtensionsPanel
├── EditorArea (编辑区)
│   ├── TabBar
│   ├── Breadcrumbs
│   └── MonacoEditor
├── Panel (底部面板)
│   ├── Terminal
│   ├── Output
│   └── Problems
└── StatusBar (状态栏)
```

---

## 7. 后端架构

### 7.1 Wails App 结构

```go
type App struct {
    ctx          context.Context
    fs           *FileSystem
    git          *GitEngine
    pluginMgr    *PluginManager
    lspMgr       *LSPManager
    config       *ConfigManager
    searchIdx    *SearchIndexer
}

// 初始化
func NewApp() *App {
    return &App{
        fs:        NewFileSystem(),
        git:       NewGitEngine(),
        pluginMgr: NewPluginManager(),
        lspMgr:    NewLSPManager(),
        config:    NewConfigManager(),
        searchIdx: NewSearchIndexer(),
    }
}
```

### 7.2 文件系统

```go
type FileSystem struct {
    watchers map[string]*fsnotify.Watcher
    cache    *FileCache
}

func (fs *FileSystem) ReadFile(path string) ([]byte, error)
func (fs *FileSystem) WriteFile(path string, content []byte) error
func (fs *FileSystem) Watch(path string) error
func (fs *FileSystem) ListDir(path string) ([]FileInfo, error)
func (fs *FileSystem) Move(oldPath, newPath string) error
func (fs *FileSystem) Delete(path string) error
```

### 7.3 配置管理

```go
type ConfigManager struct {
    userConfig    map[string]interface{}
    workspaceConfig map[string]interface{}
    defaultConfig map[string]interface{}
}

// 配置层级: 默认 < 用户 < 工作区 < 语言特定
func (cm *ConfigManager) Get(key string) interface{}
func (cm *ConfigManager) Set(key string, value interface{}, scope ConfigScope) error
```

---

## 8. 数据库设计

### 8.1 SQLite 表结构

```
-- 插件表
CREATE TABLE plugins (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    version TEXT NOT NULL,
    enabled BOOLEAN DEFAULT 1,
    config JSON,
    installed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 设置表
CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value JSON NOT NULL,
    scope TEXT NOT NULL, -- 'user' | 'workspace'
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 最近文件表
CREATE TABLE recent_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL,
    opened_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 书签表
CREATE TABLE bookmarks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_path TEXT NOT NULL,
    line INTEGER NOT NULL,
    label TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## 9. 性能优化

### 9.1 前端优化
- **虚拟滚动**: 大文件列表、提交历史
- **懒加载**: 按需加载插件和语言支持
- **Web Workers**: LSP 通信、搜索索引
- **Memoization**: 缓存计算结果
- **Code Splitting**: 路由级别代码分割

### 9.2 后端优化
- **文件缓存**: LRU 策略
- **增量搜索**: 后台索引构建
- **异步 Git**: 非阻塞 Git 操作
- **连接池**: LSP 服务器复用
- **流式处理**: 大文件分块读取

### 9.3 IPC 优化
- **批量传输**: 减少桥接调用次数
- **二进制协议**: 大数据使用 Base64
- **事件驱动**: 避免轮询
- **压缩传输**: 大数据启用压缩

---

## 10. 安全设计

### 10.1 插件安全
- **沙箱机制**: 限制插件权限
- **代码签名**: 验证插件来源
- **权限申请**: 敏感操作需授权
- **资源限制**: CPU、内存配额

### 10.2 数据安全
- **加密存储**: 敏感配置加密
- **备份机制**: 自动备份工作区
- **崩溃恢复**: 异常退出恢复
- **操作审计**: 关键操作日志

---

## 11. 扩展性设计

### 11.1 主题系统
```typescript
interface Theme {
  name: string;
  type: 'light' | 'dark';
  colors: ThemeColors;
  tokenColors: TokenColor[];
}

interface ThemeColors {
  'editor.background': string;
  'editor.foreground': string;
  'activityBar.background': string;
  // ... 更多颜色
}
```

### 11.2  snippets 系统
```json
{
  "Print to console": {
    "prefix": "log",
    "body": [
      "console.log('$1');",
      "$2"
    ],
    "description": "Log output to console"
  }
}
```

### 11.3 Keymap 系统
```json
{
  "keybindings": [
    {
      "key": "ctrl+s",
      "command": "file.save",
      "when": "editorTextFocus"
    },
    {
      "key": "ctrl+shift+f",
      "command": "search.open",
      "when": "!inputFocus"
    }
  ]
}
```

---

## 12. 开发路线图

### Phase 1: MVP (2-3个月)
- [x] 项目初始化
- [ ] 基础文件编辑
- [ ] 基本 Git 操作 (commit, push, pull)
- [ ] 简单插件加载
- [ ] 文件树浏览

### Phase 2: 核心功能 (2-3个月)
- [ ] 完整 Git 可视化
- [ ] LSP 集成
- [ ] 插件市场
- [ ] 主题系统
- [ ] 搜索功能

### Phase 3: 增强功能 (2-3个月)
- [ ] 终端集成
- [ ] 调试器支持
- [ ] 高级重构
- [ ] 协作编辑
- [ ] 云同步

### Phase 4: 优化完善 (1-2个月)
- [ ] 性能优化
- [ ] 用户体验改进
- [ ] 文档完善
- [ ] 插件生态

---

## 13. 技术风险

### 13.1 Wails v3 成熟度
- **风险**: Wails v3 处于早期阶段
- **应对**: 准备降级到 v2 方案，关注官方进展

### 13.2 性能瓶颈
- **风险**: 大文件、多文件场景性能
- **应对**: 充分测试，优化算法，引入虚拟滚动

### 13.3 插件兼容性
- **风险**: 跨平台兼容性问题
- **应对**: 完善测试矩阵，提供模拟环境

### 13.4 Git 复杂度
- **风险**: Git 操作复杂性超出预期
- **应对**: 使用成熟库 (libgit2)，逐步实现

---

## 14. 项目结构

```
hao-code/
├── backend/                 # Go 后端
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── app/
│   │   ├── filesystem/
│   │   ├── git/
│   │   ├── plugin/
│   │   ├── lsp/
│   │   ├── config/
│   │   └── search/
│   ├── pkg/
│   │   ├── types/
│   │   └── utils/
│   └── go.mod
│
├── frontend/                # 前端 (Vue 或 React)
│   ├── src/
│   │   ├── components/
│   │   ├── views/
│   │   ├── stores/ (Vue) 或 store/ (React)
│   │   ├── services/
│   │   ├── types/
│   │   ├── hooks/ (React) 或 composables/ (Vue)
│   │   ├── styles/
│   │   └── App.tsx
│   ├── public/
│   ├── package.json
│   └── vite.config.ts
│
├── extensions/              # 内置插件
│   ├── git/
│   ├── markdown/
│   ├── json/
│   └── ...
│
├── docs/                    # 文档
│   ├── architecture.md
│   ├── plugin-api.md
│   └── contributing.md
│
└── wails.json
```

---

## 15. 总结

本项目旨在打造一个现代化的代码编辑器，核心特点：

1. **跨平台**: 基于 Wails3，一套代码运行在三端
2. **可扩展**: 完善的插件系统，对标 VSCode
3. **Git 强大**: Idea 级别的 Git 可视化和管理
4. **高性能**: 多层次优化，流畅处理大项目
5. **易用性**: 直观界面，降低新手门槛

技术挑战与机遇并存，需要平衡功能完整性与性能表现。
