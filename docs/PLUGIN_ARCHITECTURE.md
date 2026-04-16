# Hao-Code 插件系统架构设计

## 1. 核心目标

构建一个对标 VSCode 的高性能、可扩展插件系统，支持前端 UI 扩展和后端逻辑增强。

## 2. 技术选型

- **前端沙箱**: `iframe` 或 `Web Worker` + `Comlink` (实现 JS 插件隔离)
- **后端加载**: Go `plugin` 包 (`.so` 动态库) 或 `gRPC` 微服务
- **通信协议**: JSON-RPC 2.0 (通过 Wails Events 或 WebSocket 桥接)
- **包管理**: 基于 npm 风格的 `.hcext` 压缩包格式

## 3. 插件结构 (Manifest)

```json
{
  "name": "my-plugin",
  "version": "1.0.0",
  "main": "./dist/extension.js",
  "contributes": {
    "commands": [{ "command": "myPlugin.hello", "title": "Hello World" }],
    "languages": [{ "id": "go", "extensions": [".go"] }]
  }
}
```

## 4. API 设计 (Extension API)

- **vscode.window**: 操作 UI (创建视图、显示消息)
- **vscode.workspace**: 操作文件系统 (读取、监听变化)
- **vscode.languages**: 注册 LSP 提供者 (补全、诊断)
- **vscode.commands**: 注册和执行命令

## 5. 生命周期管理

1. **发现**: 扫描 `~/.hao-code/extensions` 目录
2. **激活**: 根据 `activationEvents` (如 `onLanguage:go`) 懒加载
3. **运行**: 在沙箱中执行 `activate(context)`
4. **销毁**: 调用 `deactivate()` 并清理资源

## 6. 安全机制

- **权限控制**: 插件需声明所需权限 (网络、文件读写)
- **代码签名**: 验证插件来源，防止恶意代码注入
- **资源限制**: 限制插件 CPU 和内存占用
