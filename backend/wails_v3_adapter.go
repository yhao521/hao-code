package main

import (
	"context"
)

// ==================== Wails v3 适配器（预留）====================
// 当 Wails v3 正式发布后，可以启用此适配器
// 前端代码无需任何改动！

// WailsV3Adapter Wails v3 适配器示例
// 注意：这是前瞻性的设计，实际实现需要根据 v3 的 API 调整
type WailsV3Adapter struct {
	ctx      context.Context
	services *ServiceContainer
}

// NewWailsV3Adapter 创建 Wails v3 适配器
func NewWailsV3Adapter() *WailsV3Adapter {
	return &WailsV3Adapter{
		services: NewServiceContainer(),
	}
}

// Initialize Wails v3 初始化方法（假设 v3 使用不同的方法名）
func (w *WailsV3Adapter) Initialize(ctx context.Context) error {
	w.ctx = ctx
	return nil
}

// ==================== 暴露给前端的方法 ====================
// 与 WailsV2Adapter 完全相同的接口
// 这样前端代码无需修改

// Greet 问候方法
func (w *WailsV3Adapter) Greet(name string) string {
	return w.services.App.Greet(name)
}

// ReadFile 读取文件
func (w *WailsV3Adapter) ReadFile(path string) (string, error) {
	return w.services.App.ReadFile(path)
}

// WriteFile 写入文件
func (w *WailsV3Adapter) WriteFile(path string, content string) error {
	return w.services.App.WriteFile(path, content)
}

// ListDir 列出目录
func (w *WailsV3Adapter) ListDir(path string) ([]FileInfo, error) {
	return w.services.App.ListDir(path)
}

// GetProjectRoot 获取项目根目录
func (w *WailsV3Adapter) GetProjectRoot() string {
	return w.services.App.GetProjectRoot()
}

// OpenRepository 打开 Git 仓库
func (w *WailsV3Adapter) OpenRepository(path string) (*RepoInfo, error) {
	return w.services.App.OpenRepository(path)
}

// GetGitStatus 获取 Git 状态
func (w *WailsV3Adapter) GetGitStatus(path string) (*GitStatus, error) {
	return w.services.App.GetGitStatus(path)
}

// GitCommit 提交 Git 更改
func (w *WailsV3Adapter) GitCommit(path, message string) (string, error) {
	return w.services.App.GitCommit(path, message)
}

// GitGetBranches 获取 Git 分支
func (w *WailsV3Adapter) GitGetBranches(path string) (*BranchInfo, error) {
	return w.services.App.GitGetBranches(path)
}

// GitGetLog 获取 Git 日志
func (w *WailsV3Adapter) GitGetLog(path string, maxCommits int) ([]CommitInfo, error) {
	return w.services.App.GitGetLog(path, maxCommits)
}

// ==================== 迁移指南 ====================

/*
从 Wails v2 迁移到 v3 的步骤：

1. 在 main.go 中切换适配器：
   
   // 当前使用 v2
   app := NewApp() // 内部使用 WailsV2Adapter
   
   // 迁移到 v3
   app := NewAppWithV3() // 内部使用 WailsV3Adapter
   

2. 修改 NewApp 函数：
   
   func NewApp() *App {
       return &App{
           adapter: NewWailsV3Adapter(), // 切换到 v3 适配器
       }
   }
   

3. 调整启动钩子（如果 v3 有不同的生命周期方法）：
   
   func (a *App) startup(ctx context.Context) {
       a.ctx = ctx
       if v3Adapter, ok := a.adapter.(*WailsV3Adapter); ok {
           v3Adapter.Initialize(ctx) // v3 特有的初始化
       } else if v2Adapter, ok := a.adapter.(*WailsV2Adapter); ok {
           v2Adapter.Startup(ctx) // v2 的启动方式
       }
   }
   

4. 前端代码完全不需要修改！
   - TypeScript 绑定会自动重新生成
   - 方法签名保持一致
   - 调用方式完全相同

优势：
✅ 业务逻辑层（services）完全独立于 Wails 框架
✅ 可以轻松切换 Wails 版本
✅ 可以为不同平台提供不同适配器（Windows/macOS/Linux）
✅ 便于单元测试（可以 mock 服务层）
✅ 代码组织更清晰，职责分离
*/
