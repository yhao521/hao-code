# Wails v3 绑定和构建问题完整修复报告

**日期**: 2026年4月15日  
**状态**: ✅ 完全修复

---

## 🐛 遇到的问题

### 问题 1: `window.go.backend` is undefined

```
TypeError: undefined is not an object (evaluating 'window.go.backend')
```

**原因**: Wails v3 alpha.74 改变了绑定生成机制，新绑定使用了不同的 API。

### 问题 2: `npm run build:dev` 脚本缺失

```
npm error Missing script: "build:dev"
```

**原因**: package.json 中缺少开发模式构建脚本。

### 问题 3: TypeScript 类型错误（52个错误）

```
Module '"@wails/go/backend/App"' has no exported member 'GetProjectRoot'.
Module '"@wails/go/backend/App"' has no exported member 'ListDir'.
...
```

**原因**: 新的绑定结构与旧的 TypeScript 声明不兼容。

---

## ✅ 解决方案

### 修复 1: 添加 build:dev 脚本

**文件**: `frontend/package.json`

```json
{
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc --noEmit && vite build",
    "build:dev": "vite build --mode development", // ✅ 新增
    "preview": "vite preview"
  }
}
```

### 修复 2: 安装 @wailsio/runtime

```bash
cd frontend
npm install @wailsio/runtime
```

新的 Wails v3 绑定依赖这个包。

### 修复 3: 创建环境配置文件

**文件**: `frontend/.env.development`

```env
NODE_ENV=development
VITE_APP_MODE=development
VITE_DEBUG=true
```

### 修复 4: 重新生成绑定

```bash
wails3 generate bindings -clean=true
```

这会生成新的绑定到 `frontend/bindings/` 目录。

### 修复 5: 复制绑定并创建兼容层

```bash
# 复制新生成的绑定
cp -r frontend/bindings/hao-code/backend/* frontend/wailsjs/go/backend/
```

**创建兼容层文件**:

#### App.js

```javascript
// Compatibility layer for Wails v3 bindings
export * from "./appservice.js";

export {
  BranchInfo,
  Change,
  CommitInfo,
  FileInfo,
  GitStatus,
  RecentItem,
  RepoInfo,
  SearchResult,
} from "./models.js";
```

#### App.d.ts

```typescript
// Type declarations for all backend methods
export function AddRecentFile(path: string): Promise<void>;
export function AddRecentFolder(path: string): Promise<void>;
export function GetProjectRoot(): Promise<string>;
export function ListDir(path: string): Promise<any[]>;
export function ReadFile(path: string): Promise<string>;
export function WriteFile(path: string, content: string): Promise<void>;
export function OpenFileDialog(): Promise<string>;
export function OpenFolderDialog(): Promise<string>;
export function SaveFileDialog(): Promise<string>;
export function SearchInFiles(
  rootPath: string,
  searchText: string,
  caseSensitive: boolean,
  maxResults: number,
): Promise<any[]>;
// ... 其他 30+ 个方法

export {
  BranchInfo,
  Change,
  CommitInfo,
  FileInfo,
  GitStatus,
  RecentItem,
  RepoInfo,
  SearchResult,
} from "./models.js";
```

---

## 📊 修复统计

| 项目                   | 数量 |
| ---------------------- | ---- |
| 修改文件               | 3    |
| 新增文件               | 3    |
| 修复的 TypeScript 错误 | 52   |
| 导出的后端方法         | 36   |
| 数据模型               | 8    |

---

## 🔍 技术细节

### Wails v3 绑定变化

#### 旧版本 (alpha.74 之前)

```javascript
// 直接通过 window.go 访问
export function OpenFolderDialog() {
  return window["go"]["backend"]["App"]["OpenFolderDialog"]();
}
```

#### 新版本 (alpha.74+)

```javascript
// 使用 @wailsio/runtime
import { Call as $Call } from "@wailsio/runtime";

export function OpenFolderDialog() {
  return $Call.ByID(123456789);
}
```

### 兼容层设计原理

为了保持现有代码不变，创建了双层兼容：

1. **JavaScript 层** (`App.js`)
   - 重新导出 `appservice.js` 的所有方法
   - 重新导出 `models.js` 的所有类型

2. **TypeScript 层** (`App.d.ts`)
   - 手动声明所有方法的类型
   - 确保 IDE 智能提示正常工作

这样现有的导入语句仍然有效：

```typescript
import { OpenFolderDialog, SearchInFiles } from "@wails/go/backend/App";
```

---

## 📋 验证清单

- [x] 添加 build:dev 脚本
- [x] 安装 @wailsio/runtime
- [x] 创建 .env.development
- [x] 重新生成绑定
- [x] 复制绑定文件
- [x] 创建 App.js 兼容层
- [x] 创建 App.d.ts 类型声明
- [x] 前端构建成功（无 TS 错误）
- [x] 后端编译成功
- [x] 应用正常启动
- [x] 打开文件夹功能正常
- [x] SearchInFiles 方法可用

---

## 🚀 使用方法

### 开发模式

```bash
# 方法 1: 使用 dev.sh（推荐）
./dev.sh

# 方法 2: 使用 Task
task dev
```

### 生产构建

```bash
# 方法 1: 使用 dev.sh
./dev.sh  # 会自动构建前端

# 方法 2: 手动构建
cd frontend && npm run build && cd ..
go build -o build/bin/hao-code main.go
```

### 重新生成绑定

每次修改 Go 后端代码后：

```bash
# 重新生成绑定
wails3 generate bindings -clean=true

# 或使用 Task
task wails:bindings

# 复制到新位置
cp -r frontend/bindings/hao-code/backend/* frontend/wailsjs/go/backend/

# 重启应用
./dev.sh
```

---

## 🐛 故障排除

### Q1: 仍然看到 TypeScript 错误？

**A**: 确保 App.d.ts 包含所有需要的方法。检查错误信息中提到的方法是否在文件中声明。

### Q2: 运行时仍然报错 `window.go.backend`？

**A**:

1. 确保已安装 `@wailsio/runtime`
2. 完全重启应用（不是热重载）
3. 清除浏览器缓存（Cmd+Shift+R）

### Q3: build:dev 失败？

**A**: 检查是否有 `.env.development` 文件，并确保 `package.json` 中有 `build:dev` 脚本。

### Q4: 绑定生成很慢？

**A**: 这是正常的，Wails v3 需要处理所有 Go 包。首次生成可能需要 10-15 秒。

---

## 📚 相关文档

- [WAILS_BINDINGS_FIX.md](./WAILS_BINDINGS_FIX.md) - 详细的绑定修复指南
- [TASKFILE_GUIDE.md](./TASKFILE_GUIDE.md) - Task 使用指南
- [VSCODE_UI_REFACTOR.md](./VSCODE_UI_REFACTOR.md) - UI 重构报告

---

## 💡 最佳实践

### 1. 定期更新绑定

```bash
# 每次修改 Go 代码后
task wails:bindings
cp -r frontend/bindings/hao-code/backend/* frontend/wailsjs/go/backend/
```

### 2. 提交前验证

```bash
# 确保没有 TypeScript 错误
cd frontend && npm run build

# 确保应用能启动
./dev.sh
```

### 3. 备份绑定

在重新生成之前：

```bash
cp -r frontend/wailsjs frontend/wailsjs.backup
```

如果新绑定有问题：

```bash
rm -rf frontend/wailsjs
mv frontend/wailsjs.backup frontend/wailsjs
```

---

## 🎯 成果总结

✅ **完全修复了所有绑定和构建问题**

- 52 个 TypeScript 错误全部解决
- 前端构建成功
- 后端编译成功
- 应用正常启动和运行

✅ **创建了完整的兼容层**

- JavaScript 兼容层（App.js）
- TypeScript 类型声明（App.d.ts）
- 支持所有 36 个后端方法

✅ **改进了开发工作流**

- 添加了 build:dev 脚本
- 创建了环境配置
- 文档化了修复过程

---

**报告生成时间**: 2026-04-15 23:45  
**维护者**: Hao-Code Team  
**版本**: 1.0
