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

### 6.1 技术选型：Vue 3 + Naive UI

#### 为什么选择 Vue 3？
```
✅ 优势:
- Composition API 提供灵活的逻辑复用
- 响应式系统高效且直观
- 学习曲线平缓，开发效率高
- TypeScript 支持持续改进
- 社区生态成熟 (Vite, Pinia, Vue Router)
```

#### 为什么选择 Naive UI？
```
✅ 核心优势:
1. 完整组件体系 - 90+ 高质量组件，覆盖企业级应用全场景
2. TypeScript 原生 - 全量使用 TypeScript 编写，类型定义完善
3. 主题定制灵活 - 先进的类型安全主题系统，无需 CSS 变量或预处理器
4. 按需导入优化 - 所有组件支持 Tree Shaking，优化打包体积
5. 性能优化设计 - select、tree、table 等组件内置虚拟列表
6. 开发者友好 - 清晰的文档和 API 设计，快速上手
7. 官方推荐 - Vue 作者尤雨溪亲自推荐的 UI 组件库
```

#### Naive UI 关键特性
```typescript
// 1. 主题定制 - 类型安全
import { darkTheme, NConfigProvider } from 'naive-ui'

const themeOverrides = {
  common: {
    primaryColor: '#18A058',
    borderRadius: '6px'
  },
  Button: {
    textColor: '#FF6B00',
    borderHover: '1px solid #FF6B00'
  }
}

// 2. 按需导入 - 自动 Tree Shaking
import { 
  NButton, 
  NInput, 
  NTree, 
  NDataTable,
  NModal 
} from 'naive-ui'

// 3. 全局配置 - i18n & RTL
import { zhCN, dateZhCN } from 'naive-ui'
```

### 6.2 状态管理：Pinia

#### Store 设计
``typescript
// stores/editor.ts
export const useEditorStore = defineStore('editor', {
  state: () => ({
    activeEditor: null as string | null,
    editors: {} as Record<string, EditorInstance>,
    tabs: [] as Tab[],
    sidebarVisible: true,
    sidebarView: 'explorer' as SidebarView
  }),
  
  getters: {
    activeTab: (state) => state.tabs.find(t => t.id === state.activeEditor),
    dirtyTabs: (state) => state.tabs.filter(t => t.dirty)
  },
  
  actions: {
    openFile(path: string) { /* ... */ },
    closeTab(id: string) { /* ... */ },
    saveFile(id: string) { /* ... */ }
  }
})

// stores/git.ts
export const useGitStore = defineStore('git', {
  state: () => ({
    repository: null as Repository | null,
    branches: [] as Branch[],
    currentBranch: '',
    changes: [] as Change[],
    stagedChanges: [] as Change[],
    commitMessage: '',
    isCommitting: false,
    isLoading: false
  }),
  
  actions: {
    async loadRepository() { /* ... */ },
    async fetchChanges() { /* ... */ },
    async commit(message: string) { /* ... */ },
    async switchBranch(branch: string) { /* ... */ }
  }
})

// stores/plugins.ts
export const usePluginStore = defineStore('plugins', {
  state: () => ({
    installedPlugins: [] as PluginInfo[],
    activatedPlugins: new Set<string>(),
    marketplace: null as PluginMarketplace | null,
    isLoading: false
  }),
  
  actions: {
    async installPlugin(id: string) { /* ... */ },
    async uninstallPlugin(id: string) { /* ... */ },
    activatePlugin(id: string) { /* ... */ }
  }
})
```

### 6.3 Naive UI 组件映射

| 编辑器功能 | Naive UI 组件 |
|-----------|--------------|
| 文件树 | `NTree`, `NDirectoryTree` |
| 表格/列表 | `NDataTable`, `NList` |
| 对话框/模态框 | `NModal`, `NDrawer` |
| 消息提示 | `NMessage`, `NNotification` |
| 按钮/图标 | `NButton`, `NIcon` |
| 输入框 | `NInput`, `NInputGroup` |
| 标签页 | `NTabs`, `NTabPane` |
| 菜单 | `NMenu`, `NDropdown` |
| 加载状态 | `NSpin`, `NSkeleton` |
| 工具提示 | `NTooltip`, `NPopover` |
| 表单 | `NForm`, `NFormItem` |
| 选择器 | `NSelect` |
| 进度条 | `NProgress` |
| 分割面板 | `NSplit` |
| 滚动容器 | `NScrollbar` |

### 6.4 组件结构

```
App.vue
├── NLayout (整体布局)
│   ├── NLayoutHeader (标题栏)
│   │   └── CustomTitleBar
│   ├── NLayout (主体区域)
│   │   ├── NLayoutSider (侧边栏)
│   │   │   ├── NMenu (活动栏)
│   │   │   └── SidePanel
│   │   │       ├── FileExplorer (NTree)
│   │   │       ├── SearchPanel
│   │   │       ├── GitPanel
│   │   │       │   ├── ChangesView (NDataTable)
│   │   │       │   ├── HistoryView (NList)
│   │   │       │   └── BranchesView (NSelect)
│   │   │       └── ExtensionsPanel
│   │   ├── NLayoutContent (编辑区)
│   │   │   ├── NTabs (标签页)
│   │   │   ├── Breadcrumbs
│   │   │   └── MonacoEditor
│   │   └── NLayoutFooter (底部面板)
│   │       ├── Terminal
│   │       ├── Output
│   │       └── Problems
│   └── NLayoutFooter (状态栏)
│       └── StatusBar
```

### 6.5 主题系统设计

#### 深色/浅色主题切换
```typescript
// composables/useTheme.ts
export function useTheme() {
  const theme = ref<'light' | 'dark'>('dark')
  const themeOverrides = computed(() => {
    return theme.value === 'dark' ? darkThemeOverrides : lightThemeOverrides
  })
  
  return {
    theme,
    themeOverrides,
    toggleTheme: () => theme.value = theme.value === 'dark' ? 'light' : 'dark'
  }
}

// VSCode 风格深色主题
const darkThemeOverrides = {
  common: {
    primaryColor: '#0E639C',
    primaryColorHover: '#1177BB',
    bodyColor: '#1E1E1E',
    cardColor: '#252526',
    borderColor: '#3E3E42',
    textColor: '#CCCCCC'
  },
  Menu: {
    itemColor: '#CCCCCC',
    itemColorActive: '#FFFFFF',
    itemColorHover: '#2A2D2E'
  },
  Tree: {
    nodeColor: '#CCCCCC',
    nodeColorHover: '#2A2D2E'
  }
}
```

### 6.6 性能优化策略

#### Naive UI 性能优化
```typescript
// 1. 按需导入 - 减小打包体积
import { NTree, NDataTable } from 'naive-ui' // ✅ Tree Shaking

// 2. 虚拟滚动 - 大列表优化
<NVirtualList :items={largeDataSet} :item-size={32}>
  {({ item }) => <div>{item.name}</div>}
</NVirtualList>

// 3. 懒加载 - 路由级别代码分割
const GitPanel = defineAsyncComponent(() => import('./views/GitPanel.vue'))

// 4. Memoization - 缓存计算结果
const filteredFiles = computed(() => {
  return files.value.filter(f => f.name.includes(searchQuery.value))
})
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
