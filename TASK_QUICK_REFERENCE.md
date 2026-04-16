# Task 命令速查表

## 🚀 常用命令

### 开发模式

```bash
task d              # 启动开发模式
task df             # 快速开发模式
```

### 构建应用

```bash
task b              # 生产构建
task build          # 构建当前平台
```

### 📦 打包应用（新增）

#### 一键全平台打包

```bash
task package:all    # ⭐ 一次性打包所有平台和架构
```

生成的文件：

- **Windows**: `.exe`, `-installer.exe`, `-windows-amd64.zip`, `-windows-arm64.zip`
- **macOS**: `.app`, `-macos-amd64.dmg`, `-macos-amd64.zip`, `-macos-arm64.dmg`, `-macos-arm64.zip`, `-macos-universal.dmg`, `-macos-universal.zip`
- **Linux**: `AppImage`, `.deb`, `.rpm`, `.pkg.tar.zst` (amd64 + arm64)

#### 单独平台打包

```bash
# Windows
task windows:package           # NSIS 安装程序
task windows:create:zip        # ZIP 压缩包

# macOS
task darwin:package            # .app 包
task darwin:create:dmg         # DMG 磁盘映像
task darwin:create:zip         # ZIP 压缩包

# Linux
task linux:package             # 所有 Linux 包格式
```

### 清理和帮助

```bash
task c              # 清理
task help           # 查看所有任务
task check          # 检查环境
```

## 📦 完整任务列表

### 主要任务

| 命令            | 别名 | 说明     |
| --------------- | ---- | -------- |
| `task dev`      | `d`  | 开发模式 |
| `task dev:fast` | `df` | 快速开发 |
| `task build`    | `b`  | 生产构建 |
| `task clean`    | `c`  | 清理     |
| `task rebuild`  | -    | 完整重建 |

### 前端任务

| 命令                    | 说明       |
| ----------------------- | ---------- |
| `task frontend:install` | 安装依赖   |
| `task frontend:build`   | 构建前端   |
| `task frontend:dev`     | 开发服务器 |
| `task frontend:lint`    | 代码检查   |

### Go 任务

| 命令           | 说明       |
| -------------- | ---------- |
| `task go:tidy` | 整理模块   |
| `task go:fmt`  | 格式化代码 |
| `task go:vet`  | 代码检查   |
| `task go:test` | 运行测试   |

### Wails 任务

| 命令                  | 别名   | 说明       |
| --------------------- | ------ | ---------- |
| `task wails:bindings` | `bind` | 生成绑定   |
| `task wails:assets`   | -      | 生成资源   |
| `task wails:update`   | -      | 更新 Wails |

### 平台构建

| 命令                         | 说明             |
| ---------------------------- | ---------------- |
| `task build:macos`           | macOS 构建       |
| `task build:macos:universal` | Universal Binary |
| `task build:windows`         | Windows 构建     |
| `task build:linux`           | Linux 构建       |

### 工具任务

| 命令         | 说明       |
| ------------ | ---------- |
| `task help`  | 显示帮助   |
| `task check` | 检查环境   |
| `task init`  | 初始化项目 |

---

## 💡 典型工作流

### 日常开发

```bash
task d              # 开始开发
# ... 编码 ...
task b              # 构建发布版本
```

### 提交前检查

```bash
task go:fmt         # 格式化 Go 代码
task go:vet         # 检查 Go 代码
task frontend:lint  # 检查前端代码
task b              # 确保能构建
```

### 跨平台发布

```bash
# ⭐ 推荐：一键打包所有平台
task package:all

# 或者单独打包
task build:macos:universal   # macOS
task build:windows           # Windows
task build:linux             # Linux
```

### 维护

```bash
task wails:update     # 更新 Wails
task init             # 重新初始化
task c                # 清理空间
```

---

## 🔗 相关链接

- [Task 官方文档](https://taskfile.dev/)
- [Wails v3 文档](https://v3.wails.io/)
- [详细指南](./TASKFILE_GUIDE.md)

---

**提示**: 使用 `task help` 查看最新的任务列表
