# Wails3 代码编辑器设计文档

## 1. 项目概述

### 1.1 项目名称
Hao-Code Editor

### 1.2 项目愿景
打造一个轻量级但功能强大的跨平台代码编辑器，对标 VSCode 的扩展能力和 IntelliJ IDEA 的 Git 集成体验。

### 1.3 核心技术栈（最终选型）
- **后端框架**: Wails v3 (Go)
- **前端框架**: Vue 3 + TypeScript
- **UI 组件库**: Naive UI
- **代码编辑核心**: Monaco Editor
- **状态管理**: Pinia
- **构建工具**: Vite
- **插件系统**: 基于 Go plugin + JavaScript SDK

---

## 1.5 技术选型深度分析

### 1.5.1 Vue 3 vs React 18 全面对比

#### 📊 核心维度对比表

| 对比维度 | Vue 3 | React 18 | 胜出者 |
|---------|-------|----------|--------|
| **学习曲线** | ⭐⭐⭐⭐⭐ 平缓 | ⭐⭐⭐ 较陡 | Vue 3 |
| **开发效率** | ⭐⭐⭐⭐⭐ 高 | ⭐⭐⭐⭐ 较高 | Vue 3 |
| **TypeScript 支持** | ⭐⭐⭐⭐ 良好 | ⭐⭐⭐⭐⭐ 完美 | React 18 |
| **生态系统** | ⭐⭐⭐⭐ 丰富 | ⭐⭐⭐⭐⭐ 最丰富 | React 18 |
| **性能表现** | ⭐⭐⭐⭐⭐ 优秀 | ⭐⭐⭐⭐⭐ 优秀 | 平手 |
| **大型项目规范** | ⭐⭐⭐ 需自律 | ⭐⭐⭐⭐⭐ 强制 | React 18 |
| **社区资源** | ⭐⭐⭐⭐ 多 | ⭐⭐⭐⭐⭐ 最多 | React 18 |
| **就业市场** | ⭐⭐⭐⭐ 需求增长 | ⭐⭐⭐⭐⭐ 主流 | React 18 |
| **创新特性** | ⭐⭐⭐⭐⭐ 响应式 | ⭐⭐⭐⭐ Hooks | Vue 3 |
| **企业采用** | ⭐⭐⭐⭐ 增长中 | ⭐⭐⭐⭐⭐ 广泛 | React 18 |

#### 🎯 Vue 3 深度分析

**✅ 核心优势：**

1. **响应式系统革命性创新**
   ```typescript
   // Vue 3 Proxy 响应式 - 更直观
   const count = ref(0)
   count.value++ // 自动追踪依赖
   
   // vs React useState - 需要手动更新
   const [count, setCount] = useState(0)
   setCount(count + 1) // 必须调用 setter
   ```

2. **Composition API 灵活性强**
   ```typescript
   // Vue 3 Composition API
   setup() {
     const { files } = useFileSystem()
     const { commits } = useGit()
     return { files, commits }
   }
   // 逻辑复用更自然，无需自定义 Hook 包装
   ```

3. **单文件组件 (SFC) 开发体验佳**
   ```vue
   <template>
     <div>{{ message }}</div>
   </template>
   
   <script setup lang="ts">
   import { ref } from 'vue'
   const message = ref('Hello')
   </script>
   
   <style scoped>
   div { color: red; }
   </style>
   // HTML/CSS/JS 一体化，适合快速开发
   ```

4. **Naive UI 完美契合**
   - Naive UI 专为 Vue 3 打造
   - 类型定义完善
   - 主题系统先进

5. **学习成本低**
   - 模板语法接近 HTML
   - 指令系统直观 (v-if, v-for)
   - 新手友好

**❌ 劣势与挑战：**

1. **TypeScript 集成略逊**
   - 模板中的类型推导有限
   - 需要额外配置才能获得最佳体验
   
2. **大型项目规范性**
   - 灵活性过高可能导致代码风格不统一
   - 需要团队自行制定规范

3. **生态相对较小**
   - 某些领域库选择较少
   - 国际化程度不如 React

#### 🔥 React 18 深度分析

**✅ 核心优势：**

1. **TypeScript 原生级支持**
   ```tsx
   // React + TS 天衣无缝
   interface Props {
     files: File[]
     onSelect: (file: File) => void
   }
   
   const FileTree: FC<Props> = ({ files, onSelect }) => {
     // 完整的类型推导和检查
     return <div>{/* ... */}</div>
   }
   ```

2. **生态系统无可匹敌**
   - 最多的第三方库
   - 最丰富的学习资源
   - 最大的社区支持

3. **企业级项目首选**
   - Facebook、Netflix、Airbnb 等大厂背书
   - 长期稳定性有保障
   - 就业机会更多

4. **强约束性带来规范性**
   - JSX 强制使用 JavaScript 表达式
   - Hooks 规则 enforced by linter
   - 大型团队协作更顺畅

5. **并发特性 (React 18)**
   ```tsx
   // Concurrent Features
   <Suspense fallback={<Loading />}>
     <LazyComponent />
   </Suspense>
   
   // 自动批处理、过渡 API 等
   ```

**❌ 劣势与挑战：**

1. **学习曲线陡峭**
   - JSX 需要适应
   - Hooks 心智模型复杂
   - useEffect 依赖陷阱

2. **样板代码较多**
   ```tsx
   // React 需要更多代码
   const [count, setCount] = useState(0)
   const handleClick = useCallback(() => {
     setCount(c => c + 1)
   }, [])
   
   // vs Vue
   const count = ref(0)
   const handleClick = () => count.value++
   ```

3. **性能优化需要手动**
   - 需要记忆化 (useMemo, useCallback)
   - 不当使用会导致反模式

4. **与 Naive UI 不兼容**
   - Naive UI 是 Vue 专属
   - React 需要使用 Ant Design / Material-UI

#### 🎖️ 最终决策：为什么选择 Vue 3？

**针对本项目的关键因素：**

1. **Naive UI 绑定** ✅
   - Naive UI 是目前 Vue 3 生态中最优秀的组件库之一
   - 如果选择 React，将失去 Naive UI，需要改用 Ant Design
   - Naive UI 的主题系统和组件质量非常适合编辑器项目

2. **开发效率优先** ✅
   - 个人/小团队项目，开发速度至关重要
   - Vue 3 的学习曲线平缓，能快速迭代原型
   - SFC 让单文件内聚性更好

3. **性能足够优秀** ✅
   - Vue 3 的响应式系统基于 Proxy，性能优于 React 的虚拟 DOM diff
   - 对于编辑器场景（大文件、频繁更新）表现更好

4. **Monaco Editor 集成** ✅
   - Monaco Editor 是框架无关的
   - Vue 3 中集成同样简单，无需特殊适配

5. **Wails 集成友好** ✅
   - Wails 的前端只是静态资源
   - Vue 3 打包体积更小（~33kb vs React ~42kb）

**何时应该选择 React？**
- 需要招聘大量前端工程师
- 目标进入外企或大厂
- 项目需要极致的 TypeScript 体验
- 已有成熟的 React 技术栈

**结论：**
对于本编辑器项目，**Vue 3 + Naive UI 是最佳组合**。它在开发效率、组件质量、性能表现上达到了最佳平衡点。

---

### 1.5.2 Wails v2 vs Wails v3 深度对比

#### 📊 核心差异对比表

| 对比维度 | Wails v2 | Wails v3 | 推荐度 |
|---------|----------|----------|--------|
| **成熟度** | ⭐⭐⭐⭐⭐ 稳定发布 | ⭐⭐ 早期阶段 | v2 胜 |
| **文档完善度** | ⭐⭐⭐⭐⭐ 完整 | ⭐⭐ 不完善 | v2 胜 |
| **社区支持** | ⭐⭐⭐⭐ 活跃 | ⭐⭐ 较少 | v2 胜 |
| **新特性** | ⭐⭐⭐ 标准 | ⭐⭐⭐⭐⭐ 先进 | v3 胜 |
| **性能** | ⭐⭐⭐⭐ 良好 | ⭐⭐⭐⭐⭐ 优化 | v3 胜 |
| **API 设计** | ⭐⭐⭐ 传统 | ⭐⭐⭐⭐⭐ 现代 | v3 胜 |
| **长期维护** | ⭐⭐⭐ 维护模式 | ⭐⭐⭐⭐⭐ 未来方向 | v3 胜 |
| **风险评估** | ⭐⭐⭐⭐⭐ 低风险 | ⭐⭐ 高风险 | v2 胜 |

#### 🔍 Wails v2 深度分析

**✅ 核心优势：**

1. **生产就绪**
   ```go
   // Wails v2 成熟的项目结构
   func main() {
       app := NewApp()
       err := wails.Run(&options.App{
           Title:  "My App",
           Width:  1024,
           Height: 768,
           Frontend: &frontend.Assets{
               Handler: http.FS(embedded),
           },
       })
   }
   // 稳定可靠，已有大量生产案例
   ```

2. **完善的文档和教程**
   - 官方文档详尽
   - 社区教程丰富
   - Stack Overflow 上有大量问答

3. **成熟的生态系统**
   - 脚手架工具完善 (`wails init`)
   - 第三方插件和模板
   - CI/CD 流程成熟

4. **已知的问题都有解决方案**
   - 跨平台兼容性经过验证
   - 常见问题有社区支持
   - 调试工具成熟

**❌ 劣势与限制：**

1. **架构相对老旧**
   - 基于传统的 WebView
   - 某些平台需要 polyfill
   
2. **性能瓶颈**
   - IPC 通信有开销
   - 大数据传输效率低

3. **API 不够优雅**
   ```go
   // v2 的事件系统较为繁琐
   runtime.EventsEmit(a.ctx, "event-name", data)
   runtime.EventsOn(a.ctx, "event-name", func(data interface{}) {
       // 处理逻辑
   })
   ```

#### 🚀 Wails v3 深度分析

**✅ 核心优势：**

1. **现代化架构**
   ```go
   // Wails v3 简化的 API
   app := NewApp()
   app.Run() // 更简洁的入口
   ```

2. **性能优化**
   - 改进的 IPC 通信机制
   - 更小的运行时体积
   - 更快的启动速度

3. **更好的类型安全**
   - 自动生成 TypeScript 绑定
   - 编译时类型检查

4. **面向未来**
   - 支持最新的 Web 技术
   - 更好的跨平台能力
   - 长期演进方向

**❌ 风险与挑战：**

1. **早期阶段不稳定**
   - API 可能还会变化
   - Bug 修复周期不确定
   - 生产环境风险较高

2. **文档不完善**
   - 官方文档还在编写中
   - 社区资源稀缺
   - 遇到问题需要自己摸索

3. **生态系统未成熟**
   - 缺少最佳实践
   - 第三方库适配滞后
   - 工具链不完善

#### 🎯 决策建议

**方案 A：选择 Wails v2（保守策略）**

**适用场景：**
- 需要快速上线 MVP
- 团队缺乏 Go 桌面应用经验
- 追求稳定性和可预测性
- 商业项目，风险控制重要

**实施计划：**
```
Phase 1 (1-2周): 
  - wails init 创建项目
  - 集成 Vue 3 + Naive UI
  - 实现基础文件浏览
  
Phase 2 (2-3月):
  - 实现核心编辑器功能
  - Git 基础功能
  - 发布可用版本
  
Phase 3 (未来):
  - 待 v3 成熟后迁移
  - Wails 提供迁移指南
```

**方案 B：选择 Wails v3（激进策略）**

**适用场景：**
- 个人项目/开源项目
- 想尝试最新技术
- 有时间应对不确定性
- 技术探索性质

**实施计划：**
```
Phase 1 (2-3周):
  - 深入研究 v3 文档
  - 搭建基础项目
  - 评估可行性
  
Phase 2 (持续):
  - 边开发边反馈问题
  - 参与社区建设
  - 可能需要降级到 v2
```

#### 🏆 最终推荐：**混合策略**

**我的建议：以 Wails v2 为主，预留 v3 迁移路径**

**理由：**

1. **时间成本考量** ⏰
   - 编辑器本身已经是很复杂的项目
   - 不应在基础设施上消耗过多精力
   - v2 足以支撑 MVP 验证

2. **风险控制** 🛡️
   - v2 的稳定性保证项目进度
   - 避免因框架问题导致项目停滞
   - 有明确的社区支持

3. **未来可期** 🔮
   - Wails v3 正式发布后，迁移成本可控
   - 前端代码（Vue 3）无需改动
   - 只需调整 Go 后端部分

4. **实际执行策略**
   ```
   短期 (0-3个月): 使用 Wails v2 快速开发
   中期 (3-6个月): 关注 v3 进展，评估迁移
   长期 (6个月+): v3 稳定后迁移，享受新特性
   ```

**技术架构解耦设计：**
```go
// 设计原则：前后端通过接口交互，降低耦合

// backend/internal/app/interface.go
type AppInterface interface {
    ReadFile(path string) ([]byte, error)
    WriteFile(path string, content []byte) error
    // ... 其他方法
}

// 这样无论 v2 还是 v3，前端调用方式不变
```

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
``typescript
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
│   │   └── main.go         # Wails 应用入口
│   ├── internal/
│   │   ├── app/            # 主应用逻辑
│   │   ├── filesystem/     # 文件系统模块
│   │   ├── git/            # Git 引擎
│   │   ├── plugin/         # 插件管理器
│   │   ├── lsp/            # LSP 客户端
│   │   ├── config/         # 配置管理
│   │   └── search/         # 搜索索引
│   ├── pkg/
│   │   ├── types/          # 共享类型定义
│   │   └── utils/          # 工具函数
│   └── go.mod
│
├── frontend/                # 前端 (Vue 3 + TypeScript + Naive UI)
│   ├── src/
│   │   ├── components/     # 通用组件
│   │   │   ├── editor/     # 编辑器相关
│   │   │   │   ├── MonacoEditor.vue
│   │   │   │   ├── TabBar.vue
│   │   │   │   └── Breadcrumbs.vue
│   │   │   ├── layout/     # 布局组件
│   │   │   │   ├── TitleBar.vue
│   │   │   │   ├── ActivityBar.vue
│   │   │   │   ├── SideBar.vue
│   │   │   │   ├── StatusBar.vue
│   │   │   │   └── BottomPanel.vue
│   │   │   └── common/     # Naive UI 封装
│   │   │       ├── FileTree.vue      (NTree)
│   │   │       ├── DataTable.vue     (NDataTable)
│   │   │       └── SearchBox.vue     (NInput)
│   │   ├── views/          # 页面视图
│   │   │   ├── ExplorerView.vue
│   │   │   ├── SearchView.vue
│   │   │   ├── GitView/
│   │   │   │   ├── ChangesPanel.vue
│   │   │   │   ├── HistoryPanel.vue
│   │   │   │   └── BranchesPanel.vue
│   │   │   └── ExtensionsView.vue
│   │   ├── stores/         # Pinia 状态管理
│   │   │   ├── editor.ts
│   │   │   ├── git.ts
│   │   │   ├── plugins.ts
│   │   │   └── theme.ts
│   │   ├── composables/    # Vue Composition API
│   │   │   ├── useTheme.ts
│   │   │   ├── useKeybindings.ts
│   │   │   └── useFileSystem.ts
│   │   ├── services/       # 业务服务
│   │   │   ├── wails.service.ts
│   │   │   ├── plugin.service.ts
│   │   │   └── lsp.service.ts
│   │   ├── types/          # TypeScript 类型定义
│   │   │   ├── editor.types.ts
│   │   │   ├── git.types.ts
│   │   │   └── plugin.types.ts
│   │   ├── styles/         # 全局样式
│   │   │   ├── variables.css
│   │   │   └── theme-overrides.css
│   │   ├── utils/          # 工具函数
│   │   │   ├── format.ts
│   │   │   └── validation.ts
│   │   ├── App.vue         # 根组件
│   │   └── main.ts         # 应用入口
│   ├── public/
│   │   └── icons/
│   ├── package.json
│   ├── tsconfig.json
│   ├── vite.config.ts
│   └── wailsjs/            # Wails 生成的 TS 绑定
│       ├── runtime/
│       └── main/
│
├── extensions/              # 内置插件
│   ├── git/
│   ├── markdown/
│   ├── json/
│   └── ...
│
├── docs/                    # 文档
│   ├── DESIGN.md           # 设计文档
│   ├── PLUGIN_API.md       # 插件开发指南
│   └── CONTRIBUTING.md     # 贡献指南
│
├── build/                   # 构建输出
├── wails.json              # Wails 配置
└── README.md
```

---

## 15. Naive UI 快速上手指南

### 15.1 安装依赖

```bash
cd frontend
npm install naive-ui vfonts @vicons/ionicons5
```

### 15.2 基础配置

```typescript
// vite.config.ts
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  optimizeDeps: {
    include: ['naive-ui']
  }
})
```

### 15.3 主题提供者

```
<!-- App.vue -->
<template>
  <NConfigProvider
    :theme="theme"
    :theme-overrides="themeOverrides"
    :locale="zhCN"
    :date-locale="dateZhCN"
  >
    <NMessageProvider>
      <NNotificationProvider>
        <NDialogProvider>
          <RouterView />
        </NDialogProvider>
      </NNotificationProvider>
    </NMessageProvider>
  </NConfigProvider>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { 
  NConfigProvider,
  NMessageProvider,
  NNotificationProvider,
  NDialogProvider,
  darkTheme,
  zhCN,
  dateZhCN
} from 'naive-ui'

const theme = ref(darkTheme)
const themeOverrides = {
  common: {
    primaryColor: '#0E639C',
    borderRadius: '6px'
  }
}
</script>
```

### 15.4 常用组件示例

```
<!-- 文件树示例 -->
<template>
  <NTree
    :data="fileTreeData"
    :default-expand-all="true"
    block-line
    selectable
    @update:selected-keys="handleFileSelect"
  />
</template>

<script setup lang="ts">
import { NTree } from 'naive-ui'
import type { TreeOption } from 'naive-ui'

const fileTreeData: TreeOption[] = [
  {
    key: 'src',
    label: 'src',
    children: [
      { key: 'main.ts', label: 'main.ts' },
      { key: 'App.vue', label: 'App.vue' }
    ]
  }
]

const handleFileSelect = (keys: string[]) => {
  console.log('Selected:', keys)
}
</script>
```

---

## 16. 总结

本项目旨在打造一个现代化的代码编辑器，核心特点：

1. **跨平台**: 基于 Wails3，一套代码运行在三端
2. **可扩展**: 完善的插件系统，对标 VSCode
3. **Git 强大**: Idea 级别的 Git 可视化和管理
4. **高性能**: 多层次优化，流畅处理大项目
5. **易用性**: Naive UI 提供优雅的用户界面
6. **TypeScript**: 全链路类型安全，提升开发体验

### 技术栈总览

| 层级 | 技术选型 | 说明 |
|------|---------|------|
| 后端 | Wails v3 + Go | 跨平台桌面应用框架 |
| 前端 | Vue 3 + TypeScript | 现代化前端框架 |
| UI | Naive UI | 企业级组件库 |
| 编辑器 | Monaco Editor | VSCode 同款编辑器 |
| 状态管理 | Pinia | Vue 官方推荐 |
| 构建工具 | Vite | 极速开发体验 |
| Git | libgit2 | 成熟稳定的 Git 库 |

技术挑战与机遇并存，需要平衡功能完整性与性能表现。选择 Naive UI 将为项目带来：
- 🎨 **优雅的视觉设计** - 默认主题已非常美观
- 🔧 **强大的定制能力** - 类型安全的主题系统
- ⚡ **优秀的性能** - 内置虚拟滚动等优化
- 📦 **完整的组件体系** - 90+ 组件覆盖全场景

让我们开始构建吧！🚀
