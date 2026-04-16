# Hao-Code 功能清单与开发指南

Hao-Code 是一个基于 **Wails v3** 和 **Vue 3** 构建的高性能、跨平台代码编辑器，旨在复刻 VSCode 的核心体验并扩展高级功能。

## 🚀 核心功能版图

### 1. 核心编辑器 (Core Editor)

- **Monaco 集成**: 深度集成 Microsoft Monaco Editor，支持多标签页管理与拖拽排序。
- **Diff 视图**: 支持行内（Inline）与并排（Side-by-side）差异对比。
- **断点交互**: 在编辑器左侧边距点击即可设置/取消断点，并与调试器实时同步。

### 2. 智能辅助 (LSP Intelligence)

- **自动补全**: 基于 LSP 的智能代码建议。
- **实时诊断**: 语法错误与警告的实时标记（Markers）。
- **代码跳转**: 定义查找（Go to Definition）、引用查找（Find References）。
- **重构工具**: 符号重命名（Rename Symbol）、文档格式化（Format Document）。
- **智能感知**: 悬停提示（Hover Info）、签名帮助（Signature Help）。
- **高级导航**:
  - **调用层级 (Call Hierarchy)**: 递归展示函数的调用树。
  - **类型层次结构 (Type Hierarchy)**: 可视化展示类的继承关系（父类/子类）。
  - **实现查找 (Find Implementations)**: 快速定位接口或抽象方法的具体实现。
- **视觉增强**: 语义高亮（Semantic Highlighting）、文档链接解析（Document Links）。

### 3. 文件与搜索 (File & Search)

- **资源管理器**: 树形文件浏览，支持右键菜单操作。
- **全局搜索**: 集成 Ripgrep，支持正则表达式与大小写匹配。
- **工作区符号**: 通过命令面板输入 `@` 触发全项目符号检索。
- **面包屑导航**: 实时显示当前文件路径，支持点击跳转。

### 4. Git 集成

- **分支图谱**: 使用 Canvas 绘制 Git 提交历史与分支走向。
- **差异查看**: 直接在编辑器中查看暂存区与工作区的文件变动。
- **状态同步**: 状态栏实时显示当前分支名称与变更统计。

### 5. 交互终端与调试

- **PTY 终端**: 嵌入真实的系统 Shell，支持双向数据流通信。
- **Delve 调试**: 后端集成 Go Delve，前端提供变量监视、调用栈查看及断点管理 UI。

### 6. 插件系统 (Plugin System)

- **沙箱隔离**: 基于子进程（Subprocess）的隔离环境，确保插件崩溃不影响主程序。
- **生命周期管理**: 完整的激活、停用与销毁逻辑。
- **API 桥接**: 通过 JSON-RPC over Stdio 实现前后端安全通信。
- **示例插件**: 内置 `hello-world` 插件用于测试沙箱环境。

### 7. UI/UX 细节

- **欢迎页 (Welcome Page)**: 提供无工作区时的引导界面，展示最近项目与常用快捷键。
- **活动栏优化**: 动态图标缩放与选中态发光效果。
- **状态栏深化**: 实时同步 Git 分支、错误计数与语言模式。
- **全局美化**: 统一的深色主题滚动条与字体渲染优化。

### 8. 高级项目管理

- **任务运行器 (Task Runner)**: 自动识别 `package.json` scripts 并提供一键执行面板。
- **多工作区支持**: 允许同时打开多个文件夹并在资源管理器中并列显示。

---

## 📈 开发进度 (2026 Q2)

| 模块 | 进度 | 说明 |
| :--- | :--- | :--- |
| **核心编辑器** | 100% | 支持分屏布局、模型懒加载、Diff 视图 |
| **智能辅助 (LSP)** | 95% | 涵盖补全、跳转、重构及类型层次结构 |
| **文件与搜索** | 90% | 已实现虚拟滚动优化，支持大规模文件树 |
| **Git 集成** | 85% | 分支图谱、差异查看、状态同步 |
| **插件系统** | 80% | 沙箱隔离已完成，后端加载器待完善 |
| **性能优化** | 90% | 内存管理优化，大文件渲染流畅度提升 |
| **UI/UX** | 95% | 像素级复刻 VSCode 风格，增加欢迎页 |

---

## 🛠️ 开发指南

### 环境要求

- **Go**: 1.21+
- **Node.js**: 18+
- **Wails CLI**: `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`

### 常用命令

| 命令                       | 描述                                 |
| :------------------------- | :----------------------------------- |
| `wails3 dev`               | 启动开发模式（支持热重载）           |
| `wails3 build`             | 执行生产环境打包                     |
| `wails3 generate bindings` | 重新生成 Go 到前端的 TypeScript 绑定 |

### 目录结构

```text
hao-code/
├── backend/          # Go 后端逻辑 (LSP, Git, Plugin Loader)
├── frontend/         # Vue 3 前端界面
│   ├── src/
│   │   ├── components/ # UI 组件 (Editor, SideBar, Panels)
│   │   ├── stores/     # Pinia 状态管理
│   │   └── utils/      # LSP Manager, Wails Bridge
├── plugins/          # 插件示例目录
└── wails.json        # Wails 项目配置
```

### 贡献指南

1.  **Fork** 本仓库。
2.  创建你的特性分支 (`git checkout -b feature/AmazingFeature`)。
3.  提交你的改动 (`git commit -m 'Add some AmazingFeature'`)。
4.  推送到分支 (`git push origin feature/AmazingFeature`)。
5.  开启一个 **Pull Request**。

---

**Hao-Code** - 为开发者打造的下一代编辑体验。
