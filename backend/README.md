# 接口解耦架构 - 快速参考

## 📦 文件结构

```
backend/
├── interfaces.go          # 核心接口定义 (IFileSystemService, IGitService, IAppService)
├── types.go              # 共享数据类型 (FileInfo, RepoInfo, GitStatus 等)
├── file_service.go       # 文件系统服务实现
├── git_service.go        # Git 服务实现
├── app_service.go        # 应用服务层（组合所有服务）
└── wails_v3_adapter.go   # Wails v3 适配器示例和迁移指南
```

## 🔄 调用流程

```
前端调用 → Wails Bridge → WailsV2Adapter → AppService → Domain Service
                                                    ↓
                                              FileSystemService
                                              GitService
```

## 🎯 关键概念

### 1. 接口定义 (`interfaces.go`)

```go
// 文件系统服务接口
type IFileSystemService interface {
    ReadFile(path string) (string, error)
    WriteFile(path string, content string) error
    ListDir(path string) ([]FileInfo, error)
    GetProjectRoot() string
}

// Git 服务接口
type IGitService interface {
    OpenRepository(path string) (*RepoInfo, error)
    GetGitStatus(path string) (*GitStatus, error)
    GitCommit(path, message string) (string, error)
    GitGetBranches(path string) (*BranchInfo, error)
    GitGetLog(path string, maxCommits int) ([]CommitInfo, error)
}

// 应用服务接口（组合）
type IAppService interface {
    Greet(name string) string
    IFileSystemService
    IGitService
}
```

### 2. 服务容器 (`interfaces.go`)

```go
type ServiceContainer struct {
    FileSystem IFileSystemService
    Git        IGitService
    App        IAppService
}

func NewServiceContainer() *ServiceContainer {
    // 创建依赖并注入
    fs := NewFileSystemService()
    git := NewGitService()
    app := NewAppService(fs, git)
    
    return &ServiceContainer{
        FileSystem: fs,
        Git:        git,
        App:        app,
    }
}
```

### 3. 适配器模式

```go
// Wails v2 适配器
type WailsV2Adapter struct {
    ctx      context.Context
    services *ServiceContainer
}

// 暴露给前端的方法
func (w *WailsV2Adapter) ReadFile(path string) (string, error) {
    return w.services.App.ReadFile(path) // 委托给服务层
}

// Wails v3 适配器（预留）
type WailsV3Adapter struct {
    // ... v3 特定字段
}

// 相同的方法签名！
func (w *WailsV3Adapter) ReadFile(path string) (string, error) {
    return w.services.App.ReadFile(path) // 同样的委托
}
```

## 🚀 从 Wails v2 迁移到 v3

### 步骤概览

1. **修改 `app.go`** - 切换适配器
2. **调整生命周期** - 更新启动钩子
3. **重新生成绑定** - `wails dev`
4. **测试验证** - 确保功能正常

### 具体操作

#### 1. 修改 NewApp

```go
// 当前 (v2)
func NewApp() *App {
    return &App{
        adapter: NewWailsV2Adapter(),
    }
}

// 迁移后 (v3)
func NewApp() *App {
    return &App{
        adapter: NewWailsV3Adapter(),
    }
}
```

#### 2. 更新启动钩子

```go
func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    
    // 根据适配器类型调用相应方法
    switch adapter := a.adapter.(type) {
    case *WailsV2Adapter:
        adapter.Startup(ctx)
    case *WailsV3Adapter:
        adapter.Initialize(ctx)
    }
}
```

#### 3. 重新生成 TypeScript 绑定

```bash
wails dev
# 自动重新生成 frontend/wailsjs/ 目录
```

## ✅ 优势总结

| 优势 | 说明 |
|------|------|
| **前端零改动** | Vue 组件完全不需要修改 |
| **业务逻辑不变** | 所有服务实现保持不变 |
| **快速切换** | 只需修改几行代码 |
| **易于测试** | 可以 mock 服务进行单元测试 |
| **多平台支持** | 可以为不同平台提供不同实现 |

## 🧪 测试示例

```go
// Mock 文件系统服务
type MockFileService struct{}

func (m *MockFileService) ReadFile(path string) (string, error) {
    return "mock content", nil
}

// 使用 mock 进行测试
func TestReadFile(t *testing.T) {
    mockFS := &MockFileService{}
    app := NewAppService(mockFS, nil)
    
    content, _ := app.ReadFile("test.txt")
    if content != "mock content" {
        t.Error("Expected mock content")
    }
}
```

## 📚 更多信息

详细架构设计请查看 [ARCHITECTURE.md](ARCHITECTURE.md)
