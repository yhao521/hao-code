# Wails v3 构建资源说明

**生成时间**: 2026年4月15日  
**命令**: `wails3 generate build-assets`  
**Wails 版本**: v3.0.0-alpha.74

---

## 📁 生成的目录结构

```
build/
├── README.md                    # 构建说明文档
├── Taskfile.yml                 # 主构建任务配置（基于 Task）
├── config.yml                   # Wails v3 配置文件
├── appicon.png                  # 应用图标源文件 (129.5KB)
├── appicon.icon/                # macOS 图标资源
│   ├── icon_16x16.png
│   └── icon_32x32.png
├── bin/                         # 输出目录
│   └── hao-code.app/            # macOS 应用包
├── darwin/                      # macOS 特定文件
│   ├── Assets.car               # macOS 资源文件 (1.5MB)
│   ├── Info.plist               # 发布版配置
│   ├── Info.dev.plist           # 开发版配置
│   ├── Taskfile.yml             # macOS 构建任务
│   └── icons.icns               # macOS 图标 (45.7KB)
├── windows/                     # Windows 特定文件
│   ├── Taskfile.yml             # Windows 构建任务
│   ├── icon.ico                 # Windows 图标 (21.2KB)
│   ├── info.json                # 应用信息
│   ├── wails.exe.manifest       # Windows 清单文件
│   ├── installer/               # 安装程序配置
│   ├── msix/                    # MSIX 包配置
│   └── nsis/                    # NSIS 安装程序配置
├── linux/                       # Linux 特定文件
│   ├── Taskfile.yml             # Linux 构建任务
│   ├── desktop                  # Desktop 文件
│   ├── appimage/                # AppImage 配置
│   └── nfpm/                    # DEB/RPM 包配置
├── android/                     # Android 特定文件 (10 items)
├── ios/                         # iOS 特定文件 (14 items)
└── docker/                      # Docker 配置 (2 items)
```

---

## 🔧 主要配置文件

### 1. Taskfile.yml (主构建配置)

这是 Wails v3 的核心构建配置文件，使用 [Task](https://taskfile.dev/) 工具。

**主要任务**:

- `go:mod:tidy` - 运行 `go mod tidy`
- `install:frontend:deps` - 安装前端依赖
- `build:frontend` - 构建前端项目
- `generate:bindings` - 生成 Go-Wails 绑定
- `build` - 完整构建应用

**使用方法**:

```bash
# 安装 Task 工具
go install github.com/go-task/task/v3/cmd/task@latest

# 查看所有可用任务
task --list

# 构建应用
task build

# 开发模式
task dev
```

### 2. config.yml (Wails v3 配置)

包含应用信息和开发模式配置。

**关键配置**:

```yaml
version: "3"
info:
  companyName: "Hao-Code Team"
  productName: "Hao-Code Editor"
  productIdentifier: "com.haocode.editor"
  version: "0.1.0"

dev_mode:
  root_path: .
  log_level: info
  debounce: 1000
```

### 3. 平台特定配置

#### macOS (darwin/)

- **Info.plist**: 应用元数据（名称、版本、权限等）
- **Info.dev.plist**: 开发模式配置（启用调试等）
- **icons.icns**: macOS 多尺寸图标
- **Assets.car**: 编译后的资源文件

#### Windows (windows/)

- **icon.ico**: Windows 应用图标
- **info.json**: 应用详细信息（用于安装程序和属性）
- **wails.exe.manifest**: Windows 清单文件（UAC、DPI 等）
- **installer/**: Inno Setup 安装程序脚本
- **nsis/**: NSIS 安装程序脚本
- **msix/**: MSIX 应用包配置

#### Linux (linux/)

- **desktop**: Desktop Entry 文件
- **appimage/**: AppImage 打包配置
- **nfpm/**:
  - deb/ - Debian/Ubuntu 包配置
  - rpm/ - RedHat/Fedora 包配置

---

## 🚀 构建应用

### 方法 1: 使用 Task（推荐）

```bash
# 安装 Task
go install github.com/go-task/task/v3/cmd/task@latest

# 查看所有任务
task --list

# 开发模式（热重载）
task dev

# 生产构建
task build

# 仅构建前端
task build:frontend

# 清理构建产物
task clean
```

### 方法 2: 使用 wails3 命令

```bash
# 开发模式
wails3 dev

# 生产构建
wails3 build

# 指定平台
wails3 build -platform darwin/universal
wails3 build -platform windows/amd64
wails3 build -platform linux/amd64
```

### 方法 3: 使用自定义脚本

```bash
# 使用我们创建的 dev.sh
./dev.sh

# 或者手动构建
cd frontend && npm run build && cd ..
go build -o build/bin/hao-code main.go
```

---

## 📦 跨平台构建

### macOS

```bash
# Universal Binary (Intel + Apple Silicon)
wails3 build -platform darwin/universal

# 仅 Intel
wails3 build -platform darwin/amd64

# 仅 Apple Silicon
wails3 build -platform darwin/arm64
```

**输出**: `build/bin/hao-code.app`

### Windows

```bash
# Windows x64
wails3 build -platform windows/amd64

# 生成安装程序
wails3 build -platform windows/amd64 -nsis
```

**输出**:

- `build/bin/hao-code.exe` - 可执行文件
- `build/bin/hao-code-installer.exe` - 安装程序（如果使用 -nsis）

### Linux

```bash
# Linux x64
wails3 build -platform linux/amd64

# 生成 DEB 包
wails3 build -platform linux/amd64 -deb

# 生成 RPM 包
wails3 build -platform linux/amd64 -rpm

# 生成 AppImage
wails3 build -platform linux/amd64 -appimage
```

**输出**:

- `build/bin/hao-code` - 可执行文件
- `build/bin/hao-code.deb` - Debian 包
- `build/bin/hao-code.rpm` - RedHat 包
- `build/bin/hao-code.AppImage` - AppImage

---

## 🎨 自定义应用图标

### 替换图标

1. **准备图标文件**:
   - PNG 格式
   - 建议尺寸: 1024x1024 或更大
   - 透明背景

2. **替换文件**:

   ```bash
   cp your-icon.png build/appicon.png
   ```

3. **重新生成图标资源**:
   ```bash
   cd build
   wails3 generate build-assets
   ```

### 平台特定图标

- **macOS**: 自动从 `appicon.png` 生成 `.icns` 文件
- **Windows**: 自动从 `appicon.png` 生成 `.ico` 文件
- **Linux**: 使用 `appicon.png` 作为应用图标

---

## ⚙️ 自定义构建配置

### macOS Info.plist

编辑 `build/darwin/Info.plist` 添加自定义配置：

```xml
<!-- 添加相机权限 -->
<key>NSCameraUsageDescription</key>
<string>需要访问相机以扫描二维码</string>

<!-- 添加麦克风权限 -->
<key>NSMicrophoneUsageDescription</key>
<string>需要访问麦克风以录制音频</string>
```

### Windows info.json

编辑 `build/windows/info.json`:

```json
{
  "fixed": {
    "product_version": "1.0.0.0",
    "file_description": "Hao-Code Editor - 跨平台代码编辑器"
  },
  "info": {
    "CompanyName": "Hao-Code Team",
    "LegalCopyright": "© 2026 Hao-Code Team"
  }
}
```

### Linux Desktop 文件

编辑 `build/linux/desktop`:

```ini
[Desktop Entry]
Name=Hao-Code Editor
Comment=Cross-platform code editor built with Wails v3
Exec=hao-code
Icon=hao-code
Terminal=false
Type=Application
Categories=Development;IDE;
```

---

## 🔍 常见问题

### Q1: Task 命令找不到？

**A**: 安装 Task 工具：

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

### Q2: 构建失败，提示缺少依赖？

**A**: 确保已安装所有依赖：

```bash
# Go 依赖
go mod tidy

# 前端依赖
cd frontend && npm install
```

### Q3: 如何减小最终二进制文件大小？

**A**: 使用 UPX 压缩：

```bash
# 安装 UPX
brew install upx        # macOS
choco install upx       # Windows
sudo apt install upx    # Linux

# 构建时启用压缩
wails3 build -upx
```

### Q4: 如何在构建时嵌入环境变量？

**A**: 在 Taskfile.yml 中添加：

```yaml
env:
  VERSION: "1.0.0"
  BUILD_TIME: '{{now | date "2006-01-02T15:04:05Z07:00"}}'
```

然后在 Go 代码中使用：

```go
var Version = os.Getenv("VERSION")
var BuildTime = os.Getenv("BUILD_TIME")
```

---

## 📊 构建产物大小参考

根据当前配置，预期的构建产物大小：

| 平台              | 未压缩 | UPX 压缩 | 减少比例 |
| ----------------- | ------ | -------- | -------- |
| macOS (Universal) | ~25 MB | ~8 MB    | 68%      |
| Windows x64       | ~22 MB | ~7 MB    | 68%      |
| Linux x64         | ~20 MB | ~6 MB    | 70%      |

_注: 实际大小取决于应用功能和依赖_

---

## 🎯 下一步

1. **测试构建**: 在各个平台上测试构建的应用
2. **自定义图标**: 替换为 Hao-Code 的品牌图标
3. **配置签名**: 为 macOS 和 Windows 应用配置代码签名
4. **自动化构建**: 设置 CI/CD 流程自动构建和发布
5. **创建安装包**: 为各平台创建专业的安装程序

---

## 📚 相关资源

- [Wails v3 官方文档](https://v3.wails.io/)
- [Task 工具文档](https://taskfile.dev/)
- [Wails 构建指南](https://v3.wails.io/docs/reference/building)
- [跨平台打包最佳实践](https://v3.wails.io/docs/guides/cross-compilation)

---

**最后更新**: 2026-04-15 23:00  
**维护者**: Hao-Code Team
