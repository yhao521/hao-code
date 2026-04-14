# Go 接口解耦架构设计文档

## 📋 概述

本文档详细说明了 Hao-Code Editor 后端的接口解耦架构设计，该设计确保了项目可以轻松从 Wails v2 迁移到 Wails v3，而无需修改前端代码或业务逻辑。

---

## 🏗️ 架构设计原则

### 1. **依赖倒置原则 (DIP)**
- 高层模块（前端）不依赖于低层模块（Wails 框架）
- 两者都依赖于抽象（接口）
- 抽象不依赖于细节，细节依赖于抽象

### 2. **单一职责原则 (SRP)**
- 每个服务只负责一个功能领域
- FileSystemService 负责文件操作
- GitService 负责 Git 操作
- AppService 负责组合协调

### 3. **开闭原则 (OCP)**
- 对扩展开放，对修改关闭
- 新增功能时添加新服务，不修改现有代码
- 切换 Wails 版本只需更换适配器

---

## 📐 架构分层

```
┌─────────────────────────────────────────────────────┐
│                  前端层 (Vue 3)                      │
│              (完全 unaware 后端实现)                   │
└──────────────────┬──────────────────────────────────┘
                   │ Wails Bridge (IPC)
┌──────────────────┴──────────────────────────────────┐
│              适配器层 (Adapter Layer)                 │
│  ┌──────────────────┬──────────────────────────────┐ │
│  │ WailsV2Adapter   │ WailsV3Adapter (预留)        │ │
│  └──────────────────┴──────────────────────────────┘ │
└──────────────────┬──────────────────────────────────┘
                   │ 接口契约
┌──────────────────┴──────────────────────────────────┐
│              应用服务层 (AppService)                  │
│         (组合协调各个子服务，对外提供统一API)           │
└──────────────────┬──────────────────────────────────┘
                   │ 接口契约
┌──────────────────┴──────────────────────────────────┐
│            领域服务层 (Domain Services)               │
│  ┌──────────────────┬──────────────────────────────┐ │
│  │ FileSystemService │ GitService                   │ │
│  └──────────────────┴──────────────────────────────┘ │
└─────────────────────────────────────────────────────┘
```

---

## 🔌 核心接口定义

### IFileSystemService

```go
type IFileSystemService interface {
    ReadFile(path string) (string, error)
    WriteFile(path string, content string) error
    ListDir(path string) ([]FileInfo, error)
    GetProjectRoot() string
}
```

**职责：** 所有文件系统相关操作

**实现：** `FileSystemService`

---

### IGitService

```go
type IGitService interface {
    OpenRepository(path string) (*RepoInfo, error)
    GetGitStatus(path string) (*GitStatus, error)
    GitCommit(path, message string) (string, error)
    GitGetBranches(path string) (*BranchInfo, error)
    GitGetLog(path string, maxCommits int) ([]CommitInfo, error)
}
```

**职责：** 所有 Git 版本控制操作

**实现：** `GitService`

---

### IAppService

```go
type IAppService interface {
    Greet(name string) string
    IFileSystemService
    IGitService
}
```

**职责：** 组合所有服务，对外提供统一 API

**实现：** `AppService`

---

## 🔄 适配器模式

### WailsV2Adapter

```go
type WailsV2Adapter struct {
    ctx      context.Context
    services *ServiceContainer
}

// 实现与 app.go 相同的方法签名
func (w *WailsV2Adapter) ReadFile(path string) (string, error) {
    return w.services.App.ReadFile(path)
}
```

**作用：**
- 桥接 Wails v2 框架和业务逻辑
- 保持方法签名不变
- 处理框架特定的生命周期钩子

---

### WailsV3Adapter (预留)

```go
type WailsV3Adapter struct {
    ctx      context.Context
    services *ServiceContainer
}

// Wails v3 可能有不同的初始化方式
func (w *WailsV3Adapter) Initialize(ctx context.Context) error {
    w.ctx = ctx
    return nil
}

// 但暴露给前端的方法保持一致
func (w *WailsV3Adapter) ReadFile(path string) (string, error) {
    return w.services.App.ReadFile(path)
}
```

**优势：**
- ✅ 前端代码无需修改
- ✅ 业务逻辑保持不变
- ✅ 只需切换适配器

---

## 📦 服务容器

### ServiceContainer

```go
type ServiceContainer struct {
    FileSystem IFileSystemService
    Git        IGitService
    App        IAppService
}

func NewServiceContainer() *ServiceContainer {
    container := &ServiceContainer{}
    
    // 依赖注入
    fileSystemService := NewFileSystemService()
    gitService := NewGitService()
    appService := NewAppService(fileSystemService, gitService)
    
    container.FileSystem = fileSystemService
    container.Git = gitService
    container.App = appService
    
    return container
}
```

**职责：**
- 管理所有服务实例
- 处理依赖关系
- 提供统一的访问入口

---

## 📁 文件组织

```
backend/
├── interfaces.go          # 核心接口定义
├── types.go              # 共享数据类型
├── file_service.go       # 文件系统服务实现
├── git_service.go        # Git 服务实现
├── app_service.go        # 应用服务层
├── wails_v3_adapter.go   # Wails v3 适配器（预留）
└── README.md             # 本架构文档
```

---

## 🚀 迁移指南

### 从 Wails v2 到 v3

#### 步骤 1: 更新 main.go

```go
// 当前 (v2)
func main() {
    app := NewApp() // 内部使用 WailsV2Adapter
    
    err := wails.Run(&options.App{
        // ... 配置
    })
}

// 迁移到 v3
func main() {
    app := NewAppWithV3() // 内部使用 WailsV3Adapter
    
    err := wails.Run(&options.App{
        // ... v3 配置
    })
}
```

#### 步骤 2: 修改 NewApp

```go
func NewApp() *App {
    return &App{
        adapter: NewWailsV3Adapter(), // 切换到 v3
    }
}
```

#### 步骤 3: 调整启动钩子

```go
func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    
    // 根据适配器类型调用相应的初始化方法
    switch adapter := a.adapter.(type) {
    case *WailsV2Adapter:
        adapter.Startup(ctx)
    case *WailsV3Adapter:
        adapter.Initialize(ctx)
    }
}
```

#### 步骤 4: 重新生成 TypeScript 绑定

```bash
wails dev
# Wails 会自动重新生成 frontend/wailsjs/ 目录
```

---

## ✅ 迁移优势

### 1. **前端零改动**
- Vue 组件完全不需要修改
- TypeScript 类型定义自动适配
- 调用方式保持一致

### 2. **业务逻辑不变**
- FileSystemService 独立于框架
- GitService 独立于框架
- 所有业务规则保持不变

### 3. **测试友好**
```go
// 可以轻松 mock 服务进行单元测试
type MockFileService struct{}

func (m *MockFileService) ReadFile(path string) (string, error) {
    return "mock content", nil
}

func TestAppService(t *testing.T) {
    mockFS := &MockFileService{}
    app := NewAppService(mockFS, nil)
    
    content, _ := app.ReadFile("test.txt")
    if content != "mock content" {
        t.Error("Expected mock content")
    }
}
```

### 4. **多平台支持**
```go
// 可以为不同平台提供不同的实现
type WindowsFileService struct { /* Windows 特定实现 */ }
type MacOSFileService struct { /* macOS 特定实现 */ }
type LinuxFileService struct { /* Linux 特定实现 */ }

// 根据平台选择实现
func NewFileService() IFileSystemService {
    switch runtime.GOOS {
    case "windows":
        return &WindowsFileService{}
    case "darwin":
        return &MacOSFileService{}
    default:
        return &LinuxFileService{}
    }
}
```

---

## 🎯 最佳实践

### 1. **不要直接访问具体服务**

❌ **错误做法：**
```go
func (a *App) SomeMethod() {
    fs := &FileSystemService{} // 硬编码依赖
    content, _ := fs.ReadFile("test.txt")
}
```

✅ **正确做法：**
```go
func (a *App) SomeMethod() {
    content, _ := a.adapter.services.FileSystem.ReadFile("test.txt")
}
```

### 2. **通过接口交互**

❌ **错误做法：**
```go
func ProcessFile(fs *FileSystemService) {
    // 依赖具体实现
}
```

✅ **正确做法：**
```go
func ProcessFile(fs IFileSystemService) {
    // 依赖抽象接口
}
```

### 3. **在构造函数中注入依赖**

❌ **错误做法：**
```go
type MyService struct {
    fs *FileSystemService // 硬编码
}
```

✅ **正确做法：**
```go
type MyService struct {
    fs IFileSystemService // 接口
}

func NewMyService(fs IFileSystemService) *MyService {
    return &MyService{fs: fs}
}
```

---

## 🧪 测试策略

### 单元测试

```go
func TestGitService(t *testing.T) {
    service := NewGitService()
    
    // 测试打开仓库
    repo, err := service.OpenRepository("/path/to/repo")
    if err != nil {
        t.Fatalf("Failed to open repo: %v", err)
    }
    
    assert.Equal(t, "main", repo.CurrentBranch)
}
```

### 集成测试

```go
func TestAppIntegration(t *testing.T) {
    container := NewServiceContainer()
    app := container.App
    
    // 测试完整的工作流
    content, err := app.ReadFile("test.txt")
    assert.NoError(t, err)
    assert.NotEmpty(t, content)
}
```

---

## 📊 性能考量

### 内存占用

```
服务容器启动时一次性创建所有服务：
- FileSystemService: ~1KB
- GitService: ~2KB
- AppService: ~1KB
- 适配器层: ~1KB
总计: ~5KB (可忽略不计)
```

### 方法调用开销

```
调用链：
Frontend → Adapter → AppService → DomainService

每次调用增加 2-3 层函数调用栈
性能影响: < 0.1ms (可忽略)
```

---

## 🔮 未来扩展

### 1. **微服务化**

如果应用变得非常复杂，可以将服务拆分为独立的微服务：

```go
// 远程服务调用
type RemoteGitService struct {
    client *grpc.Client
}

func (r *RemoteGitService) GetGitStatus(path string) (*GitStatus, error) {
    resp, err := r.client.GetGitStatus(context.Background(), &Request{Path: path})
    return parseResponse(resp)
}
```

### 2. **插件化服务**

```go
// 允许插件注册自己的服务
type ServiceRegistry struct {
    services map[string]interface{}
}

func (r *ServiceRegistry) Register(name string, service interface{}) {
    r.services[name] = service
}

func (r *ServiceRegistry) Get(name string) interface{} {
    return r.services[name]
}
```

### 3. **缓存层**

```go
type CachedFileService struct {
    wrapped IFileSystemService
    cache   *LRUCache
}

func (c *CachedFileService) ReadFile(path string) (string, error) {
    if cached, ok := c.cache.Get(path); ok {
        return cached.(string), nil
    }
    
    content, err := c.wrapped.ReadFile(path)
    if err == nil {
        c.cache.Set(path, content)
    }
    return content, err
}
```

---

## 📝 总结

### 核心优势

1. ✅ **框架无关** - 业务逻辑不依赖 Wails 版本
2. ✅ **易于测试** - 可以通过 mock 进行单元测试
3. ✅ **灵活扩展** - 轻松添加新服务或替换实现
4. ✅ **平滑迁移** - 从 v2 到 v3 只需切换适配器
5. ✅ **代码清晰** - 职责分离，易于维护

### 关键原则

1. **面向接口编程** - 始终依赖抽象而非具体实现
2. **单一职责** - 每个服务只做一件事
3. **依赖注入** - 通过构造函数传递依赖
4. **适配器模式** - 隔离框架变化

### 迁移成本

| 项目 | 成本 |
|------|------|
| 前端代码修改 | **零成本** ✅ |
| 业务逻辑修改 | **零成本** ✅ |
| 适配器切换 | **< 1小时** ⚡ |
| 测试验证 | **2-4小时** 🧪 |
| **总计** | **~半天** 🎉 |

---

<div align="center">

**🎊 这就是面向未来的架构设计！**

现在可以安心使用 Wails v2，随时准备迁移到 v3！

</div>
