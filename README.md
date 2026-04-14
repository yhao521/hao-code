# Hao-Code Editor

<div align="center">

![License](https://img.shields.io/badge/license-AGPL--3.0-blue.svg)
![Wails](https://img.shields.io/badge/Wails-v2.10+-blue.svg)
![Vue](https://img.shields.io/badge/Vue-3.x-green.svg)
![TypeScript](https://img.shields.io/badge/TypeScript-5.x-blue.svg)
![Naive UI](https://img.shields.io/badge/Naive%20UI-latest-18a058.svg)

**一个基于 Wails + Vue 3 的现代化跨平台代码编辑器**

对标 VSCode 的插件系统 · Idea 级别的 Git 可视化 · 轻量级高性能

[设计文档](DESIGN.md) · [许可证](LICENSE)

</div>

---

## 📖 简介

Hao-Code Editor 是一个功能强大的跨平台代码编辑器，旨在提供媲美 VSCode 的开发体验和 IntelliJ IDEA 级别的 Git 可视化功能。

### ✨ 核心特性

- 🎯 **跨平台支持** - 基于 Wails，一套代码运行在 Windows、macOS、Linux
- ⚡ **高性能编辑** - Monaco Editor 内核，智能补全、语法高亮、代码折叠
- 🔌 **插件系统** - 可扩展的插件架构，对标 VSCode
- 🌿 **Git 可视化** - 分支图谱、提交历史、冲突解决，媲美 IDEA
- 🎨 **主题定制** - Naive UI 先进主题系统，深色/浅色无缝切换
- 🔍 **全局搜索** - 快速文件搜索、内容搜索、符号搜索
- 💻 **现代界面** - Vue 3 + TypeScript + Naive UI 打造优雅用户体验

## 🚀 技术栈

### 前端
- **框架**: Vue 3 + TypeScript
- **UI 库**: Naive UI
- **状态管理**: Pinia
- **构建工具**: Vite
- **编辑器**: Monaco Editor

### 后端
- **框架**: Wails v2 (Go)
- **Git 引擎**: libgit2
- **数据库**: SQLite
- **插件系统**: Go plugin + JavaScript SDK

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
│   │   │   ├── layout/     # 布局组件
│   │   │   │   ├── TitleBar.vue
│   │   │   │   ├── SideBar.vue
│   │   │   │   ├── StatusBar.vue
│   │   │   │   └── ...
│   │   │   └── editor/     # 编辑器组件
│   │   │       └── EditorArea.vue
│   │   ├── stores/         # Pinia 状态
│   │   │   ├── editor.ts
│   │   │   └── git.ts
│   │   ├── App.vue         # 根组件
│   │   └── main.ts         # 入口文件
│   └── package.json
├── app.go                  # Go 后端逻辑
├── main.go                 # 应用入口
├── DESIGN.md              # 设计文档
└── LICENSE                # AGPL-3.0 许可证
```

## 🎯 当前状态

✅ **已完成：**
- 项目脚手架搭建
- Wails v2 + Vue 3 + TypeScript 集成
- Naive UI 主题配置
- Pinia 状态管理
- Monaco Editor 集成
- 基础 UI 组件（标题栏、侧边栏、状态栏）
- 文件浏览器
- Git 面板框架

🚧 **进行中：**
- Git 功能实现
- 插件系统开发
- LSP 集成

📋 **计划中：**
- 完整 Git 可视化
- 插件市场
- 终端集成
- 调试器支持

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
