# 多架构打包支持

## 📋 概述

Hao-Code 项目现已支持为多个 CPU 架构生成安装包，包括：

- **amd64** (x86_64) - 传统的 Intel/AMD 64位处理器
- **arm64** (aarch64) - ARM 64位处理器（如 Apple Silicon M1/M2/M3）
- **universal** (macOS only) - 通用二进制，同时包含 amd64 和 arm64

## 🎯 使用方法

### 一键打包所有架构

```bash
wails3 task package:all
```

此命令会为所有平台和所有支持的架构生成安装包。

### 单独打包特定架构

#### Windows

```bash
# amd64 架构
wails3 task windows:package ARCH=amd64
wails3 task windows:create:zip ARCH=amd64

# arm64 架构
wails3 task windows:package ARCH=arm64
wails3 task windows:create:zip ARCH=arm64
```

#### macOS

```bash
# amd64 架构
wails3 task darwin:package ARCH=amd64
wails3 task darwin:create:dmg ARCH=amd64
wails3 task darwin:create:zip ARCH=amd64

# arm64 架构
wails3 task darwin:package ARCH=arm64
wails3 task darwin:create:dmg ARCH=arm64
wails3 task darwin:create:zip ARCH=arm64

# universal 架构 (amd64 + arm64)
wails3 task darwin:package:universal
wails3 task darwin:create:dmg ARCH=universal
wails3 task darwin:create:zip ARCH=universal
```

#### Linux

```bash
# amd64 架构
wails3 task linux:package ARCH=amd64

# arm64 架构
wails3 task linux:package ARCH=arm64
```

## 📦 生成的文件命名规范

### Windows

- `haoyun-code-windows-amd64.zip` - amd64 架构
- `haoyun-code-windows-arm64.zip` - arm64 架构
- `haoyun-code-installer.exe` - NSIS 安装程序（当前架构）

### macOS

- `haoyun-code-macos-amd64.dmg` - amd64 架构 DMG
- `haoyun-code-macos-amd64.zip` - amd64 架构 ZIP
- `haoyun-code-macos-arm64.dmg` - arm64 架构 DMG
- `haoyun-code-macos-arm64.zip` - arm64 架构 ZIP
- `haoyun-code-macos-universal.dmg` - universal 架构 DMG
- `haoyun-code-macos-universal.zip` - universal 架构 ZIP

### Linux

Linux 包的文件名中会包含架构信息（由 nfpm 工具自动添加）：

- `haoyun-code_amd64.deb` / `haoyun-code_arm64.deb`
- `haoyun-code.x86_64.rpm` / `haoyun-code.aarch64.rpm`
- `haoyun-code-amd64.AppImage` / `haoyun-code-arm64.AppImage`

## 🔧 技术实现

### Taskfile 配置

主 Taskfile (`Taskfile.yml`) 中的 `package:all` 任务会为每个平台分别调用不同架构的打包任务：

```yaml
package:all:
  cmds:
    # Windows amd64
    - task: windows:package
      vars:
        ARCH: amd64
    - task: windows:create:zip
      vars:
        ARCH: amd64

    # Windows arm64
    - task: windows:package
      vars:
        ARCH: arm64
    - task: windows:create:zip
      vars:
        ARCH: arm64

    # macOS amd64, arm64, universal
    # ... (类似配置)

    # Linux amd64, arm64
    # ... (类似配置)
```

### 架构变量传递

各个平台的 Taskfile 都支持 `ARCH` 变量：

- **Windows**: `build/windows/Taskfile.yml`
  - `generate:syso` 任务使用 ARCH 生成对应的 .syso 文件
  - `build:docker` 任务使用 DOCKER_ARCH 传递给 Docker 容器
  - `create:zip` 任务使用 ARCH 命名输出文件

- **macOS**: `build/darwin/Taskfile.yml`
  - `build:native` 任务使用 ARCH 设置 GOARCH
  - `create:dmg` 和 `create:zip` 任务使用 ARCH 命名输出文件
  - `package:universal` 任务创建通用二进制

- **Linux**: `build/linux/Taskfile.yml`
  - `build:docker` 任务使用 DOCKER_ARCH 传递给 Docker 容器
  - nfpm 配置文件会根据架构生成不同的包

## ⚠️ 注意事项

1. **构建时间**: 打包所有架构会显著增加构建时间，因为需要为每个架构单独编译

2. **磁盘空间**: 确保有足够的磁盘空间存储所有架构的构建产物

3. **Docker 要求**: 跨平台编译需要使用 Docker 容器

   ```bash
   wails3 task setup:docker
   ```

4. **Universal Binary**: macOS 的 universal binary 只能在 macOS 系统上创建（需要 lipo 工具）

5. **测试建议**: 在发布前，建议在目标架构的设备上测试生成的安装包

## 🚀 最佳实践

### 开发阶段

只构建当前开发机器的架构：

```bash
wails3 task build
```

### 测试阶段

构建需要测试的特定架构：

```bash
wails3 task darwin:package ARCH=arm64  # 测试 Apple Silicon
wails3 task windows:package ARCH=amd64  # 测试传统 Windows
```

### 发布阶段

使用 `package:all` 生成所有架构的安装包：

```bash
wails3 task package:all
```

## 📊 架构选择指南

| 用户设备                     | 推荐架构           | 说明                                        |
| ---------------------------- | ------------------ | ------------------------------------------- |
| Intel Mac (2020及之前)       | amd64              | 传统 Intel 处理器                           |
| Apple Silicon Mac (M1/M2/M3) | arm64 或 universal | native arm64 性能更好，universal 兼容性更好 |
| Windows PC (大多数)          | amd64              | 传统 x86_64 处理器                          |
| Windows on ARM               | arm64              | Surface Pro X 等 ARM 设备                   |
| Linux x86_64                 | amd64              | 大多数桌面和服务器                          |
| Linux ARM64                  | arm64              | Raspberry Pi 4+, ARM 服务器                 |

---

**更新时间**: 2026-04-16
**状态**: ✅ 已实现并测试
