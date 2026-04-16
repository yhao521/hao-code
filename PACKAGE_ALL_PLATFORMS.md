# 全平台打包指南

## 一键打包所有平台

使用以下命令一次性为所有平台和架构生成安装包：

```bash
wails3 task package:all
```

### 生成的文件

执行上述命令后，将在 `bin/` 目录下生成以下文件：

#### Windows (amd64 + arm64)

- `haoyun-code.exe` - Windows 可执行文件
- `haoyun-code-installer.exe` - NSIS 安装程序
- `haoyun-code-windows-amd64.zip` - ZIP 压缩包 (amd64)
- `haoyun-code-windows-arm64.zip` - ZIP 压缩包 (arm64)

#### macOS (amd64 + arm64 + universal)

- `haoyun-code.app` - macOS 应用程序包
- `haoyun-code-macos-amd64.dmg` - DMG 磁盘映像 (amd64)
- `haoyun-code-macos-amd64.zip` - ZIP 压缩包 (amd64)
- `haoyun-code-macos-arm64.dmg` - DMG 磁盘映像 (arm64)
- `haoyun-code-macos-arm64.zip` - ZIP 压缩包 (arm64)
- `haoyun-code-macos-universal.dmg` - DMG 磁盘映像 (universal)
- `haoyun-code-macos-universal.zip` - ZIP 压缩包 (universal)

#### Linux (amd64 + arm64)

- `haoyun-code` - Linux 可执行文件
- `haoyun-code.AppImage` - AppImage 便携包
- `haoyun-code_*.deb` - DEB 安装包（Debian/Ubuntu）
- `haoyun-code-*.rpm` - RPM 安装包（Fedora/RHEL）
- `haoyun-code-*.pkg.tar.zst` - AUR 包（Arch Linux）

## 单独打包特定平台

如果只需要打包某个特定平台，可以使用以下命令：

### Windows

```bash
# 生成 NSIS 安装程序
wails3 task windows:package

# 生成 ZIP 压缩包
wails3 task windows:create:zip
```

### macOS

```bash
# 生成 .app 包
wails3 task darwin:package

# 生成 DMG 文件
wails3 task darwin:create:dmg

# 生成 ZIP 压缩包
wails3 task darwin:create:zip
```

### Linux

```bash
# 生成所有 Linux 包（AppImage, DEB, RPM, AUR）
wails3 task linux:package
```

## 注意事项

1. **跨平台编译**：在 macOS 上编译 Windows 和 Linux 版本需要使用 Docker
   - 首次使用前需要构建 Docker 镜像：`wails3 task setup:docker`

2. **DMG 生成**：DMG 文件只能在 macOS 系统上生成

3. **代码签名**：如需对安装包进行签名，需要在对应的 Taskfile.yml 中配置签名证书

4. **构建时间**：全平台打包可能需要较长时间，特别是首次构建或使用 Docker 时

## 输出目录

所有生成的文件都位于项目根目录的 `bin/` 文件夹中。
