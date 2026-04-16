# Hao-Code vs VSCode 功能对比报告

## 1. 项目概述

**Hao-Code** 是一个基于 **Wails v3 + Vue 3 + Monaco Editor** 构建的跨平台代码编辑器。其目标是复刻 VSCode 的核心体验，并集成前沿的 AI 辅助与自动化开发特性。

---

## 2. 核心功能对比表

| 功能模块     | VSCode (原生/插件)                  | Hao-Code (当前实现)                        | 状态        |
| :----------- | :---------------------------------- | :----------------------------------------- | :---------- |
| **基础编辑** | Monaco Editor, 多光标, 智能缩进     | Monaco Editor, 分屏布局 (Split Editor)     | ✅ 已实现   |
| **文件管理** | 资源管理器, 拖拽上传, 最近打开      | 树形文件浏览, 最近文件/文件夹记录          | ✅ 已实现   |
| **版本控制** | Git Graph, GitLens, Staging Area    | Git 图谱可视化, Blame 悬浮窗, **逐行暂存** | ✅ 已实现   |
| **终端集成** | Integrated Terminal, Split Terminal | **WebSocket 双向通信**, 多实例分屏         | ✅ 已实现   |
| **搜索替换** | 全局搜索, 正则匹配, 排除文件        | 递归内容扫描, 编辑器内跳转定位             | ✅ 已实现   |
| **LSP 支持** | 自动补全, 错误检查, 悬停提示        | Go/TS 语言服务器, **Error Lens**, 格式化   | ✅ 已实现   |
| **AI 辅助**  | GitHub Copilot, Cursor              | **Ghost Text**, @file 引用, 聊天侧边栏     | ✅ 已实现   |
| **API 测试** | Thunder Client, Postman             | **API Tester**, 历史记录, **环境变量**     | ✅ 已实现   |
| **任务运行** | Task Runner, npm scripts 检测       | 自动识别 package.json/Makefile, 终端联动   | ✅ 已实现   |
| **调试功能** | Breakpoints, Call Stack, Variables  | 断点管理, 变量监视, 步进执行               | 🚧 基础实现 |
| **扩展系统** | Marketplace, Extension API          | 插件沙箱, JSON-RPC 通信, 动态安装          | 🚧 基础架构 |
| **协同编辑** | Live Share                          | -                                          | ❌ 待开发   |

---

## 3. 深度功能解析

### 3.1 增强型 Git 工作流

- **VSCode**: 依赖 GitLens 等插件实现复杂的 Diff 查看和 Blame 信息展示。
- **Hao-Code**:
  - **逐行暂存 (Staging Area)**: 实现了类似 `git add -p` 的 UI 交互，支持通过复选框精准选择变更行。
  - **Blame 模态框**: 点击行号旁的颜色块即可弹出包含作者、时间、Hash 的详细提交记录。
  - **分支图谱**: 使用 Canvas 绘制的交互式 Git 历史分支图。

### 3.2 智能化 LSP 集成

- **VSCode**: 默认提供基础诊断，需安装 Error Lens 插件才能在行尾显示错误。
- **Hao-Code**:
  - **内置 Error Lens**: 直接将 LSP 诊断信息（Errors/Warnings）以行内装饰器形式渲染在代码末尾。
  - **自动格式化**: 监听保存事件，自动触发 Prettier/gofmt 进行代码美化。

### 3.3 生产力工具套件

- **API Tester**:
  - 集成了类似 Thunder Client 的功能，支持 Method/URL/Headers/Body 配置。
  - **环境变量支持**: 允许使用 `{{base_url}}` 语法快速切换测试环境。
  - **持久化历史**: 自动保存最近 50 条请求记录，支持一键重放。
- **Task Runner**:
  - 自动扫描项目根目录的 `package.json` 和 `Makefile`。
  - **终端联动**: 点击任务直接在底部 WebSocket 终端中执行并实时查看输出。

### 3.4 终端与协同

- **真实 PTY 交互**: 摒弃了模拟输出，通过 `gorilla/websocket` 实现了真正的双向字节流传输，支持 `vim`, `top` 等交互式命令。
- **分屏布局**: 支持水平/垂直方向的终端拆分，并可动态调整布局方向。

---

## 4. 技术栈差异

| 维度         | VSCode                        | Hao-Code                        |
| :----------- | :---------------------------- | :------------------------------ |
| **核心框架** | Electron (Chromium + Node.js) | **Wails v3** (Go + WebView)     |
| **前端技术** | React, TypeScript             | **Vue 3**, TypeScript, Naive UI |
| **后端逻辑** | JavaScript/TypeScript (Node)  | **Go** (原生高性能)             |
| **资源占用** | 较高 (通常 >300MB)            | **极低** (通常 <50MB)           |
| **打包体积** | 较大 (~100MB+)                | **极小** (~15MB)                |

---

## 5. 下一步演进建议

根据对比分析，建议优先完善以下领域以缩小与 VSCode 的差距：

1.  **实时协同编辑 (Real-time Co-editing)**: 引入 CRDT (如 Yjs) 和 WebSocket 广播机制，实现多人同时编辑同一文件。
2.  **插件市场 (Extension Marketplace)**: 完善插件沙箱的安全性，建立简单的插件索引和一键安装流程。
3.  **高级调试器 (Advanced Debugger)**: 增强对 Go/JS 运行时状态的监控，支持条件断点和表达式求值。
4.  **自定义主题引擎**: 允许用户导入 VSCode 的 `.json` 主题配置文件，实现无缝迁移。

---

_生成时间: 2026-04-15_
_版本: Hao-Code v1.0 Alpha_
