# 文件管理系统功能说明

## 📁 概述

Hao-Code Editor 已实现完整的文件管理系统，提供类似 VSCode 的文件浏览和操作体验。

---

## ✨ 已实现功能

### 1. 文件浏览
- ✅ **文件树显示**：使用 Naive UI 的 NTree 组件展示项目文件结构
- ✅ **懒加载**：点击展开目录时动态加载子文件，提升性能
- ✅ **文件选择**：点击文件在编辑器中打开
- ✅ **刷新功能**：一键刷新文件树

### 2. 文件创建
- ✅ **新建文件**：通过工具栏按钮创建新文件
- ✅ **新建文件夹**：通过工具栏按钮创建新目录
- ✅ **自动创建父目录**：创建文件时自动创建不存在的父目录

### 3. 文件操作
- ✅ **重命名**：修改文件或文件夹名称
- ✅ **删除**：删除文件或目录（支持递归删除）
- ✅ **复制**：复制文件或目录（自动添加 -copy 后缀）
- ✅ **移动**：移动文件或目录到指定位置

### 4. 文件信息
- ✅ **文件统计**：获取文件大小、修改时间等信息
- ✅ **文件扩展名**：获取文件扩展名
- ✅ **文本文件检测**：判断文件是否为文本文件
- ✅ **目录树获取**：获取指定深度的完整目录树

### 5. 搜索功能
- ✅ **文件搜索**：在目录中搜索包含关键词的文件
- ✅ **结果限制**：限制搜索结果数量，避免性能问题
- ✅ **智能过滤**：自动跳过隐藏文件和 node_modules

### 6. 高级功能
- ✅ **文件备份**：生成 .bak 备份文件
- ✅ **Touch 操作**：创建空文件或更新时间戳
- ✅ **编码处理**：自动处理 UTF-8 BOM

---

## 🎨 用户界面

### 工具栏按钮

文件管理器顶部提供三个快捷按钮：

1. **📄 新建文件**（+ 图标）
   - 点击后弹出对话框
   - 输入文件名（如：`index.ts`）
   - 点击"创建"完成

2. **📁 新建文件夹**（文件夹图标）
   - 点击后弹出对话框
   - 输入文件夹名（如：`components`）
   - 点击"创建"完成

3. **🔄 刷新**（刷新图标）
   - 重新加载文件树
   - 显示成功提示

### 文件树交互

- **单击文件**：在编辑器中打开
- **点击文件夹箭头**：展开/折叠目录
- **懒加载**：首次展开时加载子文件

---

## 🔧 技术实现

### 后端 (Go)

#### 核心服务
```
backend/
├── file_service.go      # 文件系统服务实现
├── interfaces.go        # IFileSystemService 接口定义
├── app_service.go       # 应用服务层（委托调用）
├── interfaces.go        # WailsV2Adapter 适配器
└── app.go              # Wails 绑定入口
```

#### 主要方法

| 方法 | 功能 | 参数 | 返回值 |
|------|------|------|--------|
| `ReadFile` | 读取文件内容 | path: string | content: string, error |
| `WriteFile` | 写入文件内容 | path: string, content: string | error |
| `ListDir` | 列出目录内容 | path: string | []FileInfo, error |
| `CreateFile` | 创建新文件 | path: string | error |
| `CreateDirectory` | 创建新目录 | path: string | error |
| `DeleteFileOrDirectory` | 删除文件或目录 | path: string | error |
| `RenameFileOrDirectory` | 重命名 | oldPath, newPath: string | error |
| `CopyFileOrDirectory` | 复制 | sourcePath, targetPath: string | error |
| `SearchFiles` | 搜索文件 | rootPath, keyword: string, maxResults: int | []FileInfo, error |

### 前端 (Vue 3 + TypeScript)

#### 组件结构
```
frontend/src/components/layout/FileExplorer.vue
```

#### 状态管理
```typescript
// 模态框状态
const showNewFileModal = ref(false)
const showNewFolderModal = ref(false)
const showRenameModal = ref(false)
const showDeleteModal = ref(false)

// 表单数据
const newFileForm = ref({ name: '' })
const newFolderForm = ref({ name: '' })
const renameForm = ref({ name: '' })
```

#### Wails 绑定调用
```typescript
import { 
  CreateFile,
  CreateDirectory,
  DeleteFileOrDirectory,
  RenameFileOrDirectory,
  CopyFileOrDirectory
} from '@wails/go/backend/App'
```

---

## 📝 使用示例

### 创建文件
```typescript
// 用户操作：
// 1. 点击工具栏 "新建文件" 按钮
// 2. 输入文件名 "test.ts"
// 3. 点击"创建"

// 后端调用：
await CreateFile('/path/to/project/test.ts')
```

### 重命名文件
```typescript
// 用户操作：
// 1. 在表单中输入新名称
// 2. 点击"确定"

// 后端调用：
await RenameFileOrDirectory('/path/to/old.ts', '/path/to/new.ts')
```

### 删除文件
```typescript
// 用户操作：
// 1. 点击删除（显示确认对话框）
// 2. 确认删除

// 后端调用：
await DeleteFileOrDirectory('/path/to/file.ts')
```

---

## 🔒 安全特性

1. **删除确认**：删除操作需要二次确认
2. **路径检查**：所有操作前检查路径是否存在
3. **权限验证**：自动处理文件权限
4. **错误处理**：完善的错误提示和用户反馈

---

## ⚡ 性能优化

1. **懒加载**：仅在展开目录时加载子文件
2. **过滤规则**：自动跳过隐藏文件和 node_modules
3. **结果限制**：搜索功能限制最大结果数
4. **异步操作**：所有文件操作使用异步模式

---

## 🚀 未来计划

### 待实现功能
- [ ] **右键菜单**：完整的右键上下文菜单
- [ ] **拖拽操作**：拖拽移动文件和文件夹
- [ ] **多选操作**：批量删除、移动、复制
- [ ] **文件预览**：图片、PDF 等文件预览
- [ ] **撤销操作**：支持撤销删除等操作
- [ ] **文件监视**：实时监听文件系统变化
- [ ] **最近文件**：显示最近打开的文件列表
- [ ] **书签功能**：为常用文件添加书签

### 高级功能
- [ ] **Git 集成**：在文件树上显示 Git 状态
- [ ] **文件图标**：根据文件类型显示不同图标
- [ ] **排序功能**：按名称、大小、时间排序
- [ ] **过滤功能**：按文件类型过滤显示
- [ ] **压缩/解压**：支持 ZIP 等格式
- [ ] **文件比较**：并排比较两个文件差异

---

## 📚 相关文件

- **设计文档**：[DESIGN.md](../DESIGN.md)
- **架构文档**：[backend/ARCHITECTURE.md](../backend/ARCHITECTURE.md)
- **实现报告**：[backend/IMPLEMENTATION_SUMMARY.md](../backend/IMPLEMENTATION_SUMMARY.md)

---

## 🎯 总结

文件管理系统已实现核心功能，提供完整的 CRUD 操作，满足日常开发需求。后续将逐步完善用户体验和高级功能。

**当前版本：v1.0**
**最后更新：2026-04-14**
