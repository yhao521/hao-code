# Hao-Code Editor

<div align="center">

![License](https://img.shields.io/badge/license-AGPL--3.0-blue.svg)
![Wails](https://img.shields.io/badge/Wails-v3.0--alpha-blue.svg)
![Vue](https://img.shields.io/badge/Vue-3.x-green.svg)
![TypeScript](https://img.shields.io/badge/TypeScript-5.x-blue.svg)
![Naive UI](https://img.shields.io/badge/Naive%20UI-latest-18a058.svg)

**一个基于 Wails v3 + Vue 3 的现代化跨平台代码编辑器**

对标 VSCode 的全量功能 · Idea 级别的 Git 可视化 · LSP 智能辅助

[设计文档](DESIGN.md) · [许可证](LICENSE)

</div>

---

## 📖 简介

Hao-Code Editor 是一个功能强大的跨平台代码编辑器，旨在提供媲美 VSCode 的开发体验和 IntelliJ IDEA 级别的 Git 可视化功能。

### ✨ 核心特性

- 🎯 **跨平台支持** - 基于 Wails，一套代码运行在 Windows、macOS、Linux
- ⚡ **高性能编辑** - Monaco Editor 内核，支持多标签、Diff 视图与断点交互
- 🔌 **LSP 智能辅助** - 自动补全、实时诊断、代码跳转、重命名与格式化
- 🌿 **Git 可视化** - 分支图谱、提交历史、差异对比、冲突解决，媲美 IDEA
- 💻 **交互终端** - 底部嵌入真正的系统 Shell (PTY)，支持命令执行与日志查看
- 🔍 **高级导航** - 全局搜索、符号大纲 (Outline)、查找引用 (Find References)
- 🛠️ **调试支持** - 集成 Delve 调试器，支持前端断点管理与变量监控

## 🚀 技术栈

### 前端

- **框架**: Vue 3 + TypeScript
- **UI 库**: Naive UI
- **状态管理**: Pinia
- **构建工具**: Vite
- **编辑器**: Monaco Editor

### 后端

- **框架**: Wails v3 (Go)
- **Git 引擎**: libgit2 / go-git
- **调试器**: Delve (dlv)
- **LSP 服务**: gopls, typescript-language-server 等
- **终端**: creack/pty (伪终端实现)

## 📦 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- Wails CLI

### 安装依赖

```bash
# 进入项目目录
cd hao-code

# 安装前端依赖
cd frontend
npm install
```

### 开发模式

```bash
# 在项目根目录运行
wails dev
```

这将启动开发服务器并打开应用程序窗口。

### 构建生产版本

```bash
# 构建应用
wails build

# 构建结果位于 build/bin/ 目录
```

## 🏗️ 项目结构

```
hao-code/
├── frontend/                # Vue 3 前端
│   ├── src/
│   │   ├── components/     # 组件
│   │   │   ├── layout/     # 布局组件 (SideBar, TitleBar, StatusBar...)
│   │   │   ├── editor/     # 编辑器核心 (Monaco 集成)
│   │   │   └── ...         # 功能面板 (SearchPanel, GitPanel, OutlinePanel...)
│   │   ├── stores/         # Pinia 状态管理
│   │   ├── utils/          # 工具类 (LSPManager, DiagnosticsManager...)
│   │   └── main.ts         # 入口文件
│   └── package.json
├── backend/                 # Go 后端逻辑
│   ├── app_service.go      # 应用服务主入口
│   ├── lsp_service.go      # LSP 协议处理中心
│   ├── git_service.go      # Git 业务逻辑
│   └── ...                 # 其他服务模块
├── build/                   # 构建配置与产物
├── wails.json              # Wails v3 配置文件
└── Taskfile.yml            # 任务自动化脚本
```

## 🎯 当前状态

✅ **已完成：**

- [x] Wails v3 + Vue 3 + TypeScript 深度集成
- [x] VSCode 风格五栏布局 (ActivityBar, SideBar, Editor, Panel, Status)
- [x] Monaco Editor 深度定制 (多标签、Diff、断点、LSP 提供者)
- [x] 完整 LSP 链路 (补全、诊断、跳转、大纲、引用、重命名、格式化)
- [x] Git 可视化面板 (分支图谱、时间线、差异查看)
- [x] 交互式终端 (PTY 双向通信)
- [x] Delve 调试器集成 (后端控制 + 前端 UI)
- [x] 高级交互 (命令面板、拖拽排序、高度调整、面包屑导航)

📋 **计划中：**

- [ ] 插件市场与扩展管理系统
- [ ] 更多语言服务器的预置支持
- [ ] 性能优化与内存泄漏排查
- [ ] 跨平台打包测试 (Windows/Linux)

## 🤝 参与贡献

我们欢迎所有形式的贡献！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 **GNU Affero General Public License v3.0 (AGPL-3.0)** 开源许可证。

这是一个强 copyleft 许可证，确保：

- ✅ 源代码必须公开
- ✅ 修改后的版本也必须使用 AGPL-3.0
- ✅ 网络使用也视为分发（触发源码公开义务）
- ⚠️ 商业使用需谨慎（必须开放整个应用源代码）

详细条款请查看 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- [Wails](https://wails.io/) - 优秀的跨平台桌面框架
- [Vue.js](https://vuejs.org/) - 渐进式 JavaScript 框架
- [Naive UI](https://www.naiveui.com/) - 高质量的 Vue 3 组件库
- [Monaco Editor](https://microsoft.github.io/monaco-editor/) - VSCode 同款编辑器
- [VSCode](https://code.visualstudio.com/) - 微软开发的代码编辑器
- [IntelliJ IDEA](https://www.jetbrains.com/idea/) - JetBrains 开发的 Java IDE

## 📧 联系方式

- 项目主页: [GitHub Repository](https://github.com/yourusername/hao-code)
- 问题反馈: [Issues](https://github.com/yourusername/hao-code/issues)
- 讨论交流: [Discussions](https://github.com/yourusername/hao-code/discussions)

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给我们一个 Star！**

Made with ❤️ by Hao-Code Team

</div>
