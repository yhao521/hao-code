# Taskfile 使用指南

**Hao-Code Editor** 使用 [Task](https://taskfile.dev/) 作为任务运行器，提供了简化的开发工作流。

---

## 📋 前置要求

### 安装 Task

```bash
# macOS
brew install go-task

# Linux
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin

# Windows (Scoop)
scoop install task

# Go 安装
go install github.com/go-task/task/v3/cmd/task@latest
```

验证安装：

```bash
task --version
```

---

## 🚀 常用命令

### 开发

```bash
# 启动开发模式（推荐）
task dev
# 或简写
task d

# 快速开发（跳过绑定生成）
task dev:fast
# 或简写
task df
```

### 构建

```bash
# 生产构建
task build
# 或简写
task b

# macOS Universal Binary
task build:macos:universal

# Windows 构建
task build:windows

# Linux 构建
task build:linux
```

### 清理

```bash
# 清理构建产物
task clean
# 或简写
task c
```

### 其他

```bash
# 查看所有可用任务
task help

# 检查开发环境
task check

# 初始化项目
task init

# 完整重建
task rebuild
```

---

## 📚 完整任务列表

### 主要任务

| 命令            | 别名      | 说明         |
| --------------- | --------- | ------------ |
| `task dev`      | `task d`  | 启动开发模式 |
| `task dev:fast` | `task df` | 快速开发模式 |
| `task build`    | `task b`  | 生产构建     |
| `task clean`    | `task c`  | 清理构建产物 |

### 前端任务

| 命令                    | 说明               |
| ----------------------- | ------------------ |
| `task frontend:install` | 安装前端依赖       |
| `task frontend:build`   | 构建前端           |
| `task frontend:dev`     | 启动前端开发服务器 |
| `task frontend:lint`    | 前端代码检查       |

### Go 任务

| 命令           | 说明             |
| -------------- | ---------------- |
| `task go:tidy` | 运行 go mod tidy |
| `task go:fmt`  | 格式化 Go 代码   |
| `task go:vet`  | 运行 go vet      |
| `task go:test` | 运行 Go 测试     |

### Wails 任务

| 命令                  | 别名        | 说明            |
| --------------------- | ----------- | --------------- |
| `task wails:bindings` | `task bind` | 生成 Wails 绑定 |
| `task wails:assets`   | -           | 生成构建资源    |
| `task wails:update`   | -           | 更新 Wails      |

### 平台特定构建

| 命令                         | 说明                   |
| ---------------------------- | ---------------------- |
| `task build:macos`           | macOS 构建             |
| `task build:macos:universal` | macOS Universal Binary |
| `task build:windows`         | Windows 构建           |
| `task build:linux`           | Linux 构建             |

### 工具任务

| 命令           | 说明         |
| -------------- | ------------ |
| `task help`    | 显示所有任务 |
| `task check`   | 检查开发环境 |
| `task init`    | 初始化项目   |
| `task rebuild` | 完整重建     |

---

## 💡 使用示例

### 日常开发流程

```bash
# 1. 首次克隆项目后初始化
task init

# 2. 启动开发模式
task dev

# 3. 修改代码，自动热重载...

# 4. 准备发布时构建
task build

# 5. 查看构建产物
ls -lh build/bin/
```

### 跨平台构建

```bash
# 构建所有平台
task build:macos
task build:windows
task build:linux

# 或创建 macOS Universal Binary
task build:macos:universal
```

### 代码质量检查

```bash
# Go 代码格式化和检查
task go:fmt
task go:vet

# 前端代码检查
task frontend:lint

# 运行测试
task go:test
```

### 维护任务

```bash
# 更新 Wails 到最新版本
task wails:update

# 重新生成绑定
task wails:bindings

# 清理并重新构建
task rebuild
```

---

## 🔧 自定义配置

### 修改变量

在 `Taskfile.yml` 顶部可以修改变量：

```yaml
vars:
  APP_NAME: hao-code # 应用名称
  BUILD_DIR: build/bin # 构建输出目录
  FRONTEND_DIR: frontend # 前端目录
```

### 添加新任务

在 `tasks:` 部分添加新任务：

```yaml
tasks:
  my-custom-task:
    desc: My custom task
    cmds:
      - echo "Hello, World!"
```

---

## 🐛 故障排除

### Task 找不到？

确保已安装 Task：

```bash
which task
```

如果未安装，参考上方的安装指南。

### 权限问题？

给 dev.sh 添加执行权限：

```bash
chmod +x dev.sh
```

### 绑定生成失败？

重新生成绑定：

```bash
task wails:bindings
```

### 前端依赖问题？

重新安装依赖：

```bash
task frontend:install
```

---

## 📖 更多信息

- [Task 官方文档](https://taskfile.dev/)
- [Wails v3 文档](https://v3.wails.io/)
- [Go 官方文档](https://golang.org/doc/)

---

## 🎯 最佳实践

1. **始终使用 Task 运行命令**
   - 保持一致性
   - 自动化依赖管理
   - 简化工作流

2. **定期清理和重建**

   ```bash
   task clean
   task build
   ```

3. **提交前检查代码**

   ```bash
   task go:fmt
   task go:vet
   task frontend:lint
   ```

4. **使用别名提高效率**
   ```bash
   task d    # 代替 task dev
   task b    # 代替 task build
   task c    # 代替 task clean
   ```

---

**最后更新**: 2026-04-15  
**维护者**: Hao-Code Team
