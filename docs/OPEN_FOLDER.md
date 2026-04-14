# 打开文件夹功能实现说明

## 🎯 功能概述

实现了类似 VSCode 的"打开文件夹"功能，允许用户通过系统原生对话框选择项目目录，并自动加载该目录下的文件树。

---

## ✨ 功能特性

### 1. **多种打开方式**
- ✅ **文件菜单** - 点击顶部菜单栏"文件" > "打开文件夹..."
- ✅ **快捷按钮** - 标题栏右侧的文件夹图标按钮
- ✅ **空状态按钮** - 文件树显示"未打开文件夹"时点击按钮

### 2. **跨平台支持**
- ✅ **macOS** - 使用 `osascript` (AppleScript) 调用系统对话框
- ✅ **Windows** - 使用 PowerShell + WinForms 对话框
- ✅ **Linux** - 使用 `zenity` 对话框

### 3. **完整工作流**
1. 用户选择文件夹
2. 验证文件夹可访问性
3. 设置工作区路径
4. 自动加载文件树
5. 显示成功提示

### 4. **状态管理**
- ✅ 工作区信息存储在 Pinia store
- ✅ 切换工作区时自动清空旧 tabs
- ✅ 文件树根据工作区状态显示/隐藏

---

## 🔧 技术实现

### 后端实现

#### 1. OpenFolderDialog 方法
文件: `/backend/file_service.go`

```go
func (f *FileSystemService) OpenFolderDialog() (string, error) {
    var cmd *exec.Cmd
    
    switch runtime.GOOS {
    case "darwin": // macOS
        cmd = exec.Command("osascript", "-e", `
            tell application "System Events"
                set theFolder to choose folder with prompt "选择项目文件夹"
                return POSIX path of theFolder
            end tell
        `)
    case "windows":
        // PowerShell 实现
    case "linux":
        // zenity 实现
    }
    
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("dialog cancelled or failed: %v", err)
    }
    
    path := strings.TrimSpace(string(output))
    
    // 验证路径
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return "", fmt.Errorf("selected path does not exist")
    }
    
    return path, nil
}
```

#### 2. SetProjectRoot 方法
```go
func (f *FileSystemService) SetProjectRoot(path string) error {
    // 验证路径是否存在
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return fmt.Errorf("path does not exist: %s", path)
    }
    
    // 切换到该目录
    return os.Chdir(path)
}
```

#### 3. 接口定义
文件: `/backend/interfaces.go`

```go
type IFileSystemService interface {
    // ... 其他方法
    OpenFolderDialog() (string, error)
    SetProjectRoot(path string) error
}
```

### 前端实现

#### 1. 状态管理
文件: `/frontend/src/stores/editor.ts`

```typescript
export interface Workspace {
  path: string
  name: string
}

// Store 状态
const workspace = ref<Workspace | null>(null)

// Actions
function setWorkspace(path: string) {
  const name = path.split('/').filter(Boolean).pop() || path
  workspace.value = { path, name }
  tabs.value = []  // 清空旧标签
  activeEditor.value = null
}

function clearWorkspace() {
  workspace.value = null
  tabs.value = []
  activeEditor.value = null
}
```

#### 2. TitleBar 菜单
文件: `/frontend/src/components/layout/TitleBar.vue`

```vue
<!-- 文件菜单 -->
<NDropdown
  :options="fileMenuOptions"
  @select="handleFileMenuSelect"
  trigger="click"
>
  <div class="menu-item">文件</div>
</NDropdown>

<!-- 快捷按钮 -->
<NButton text circle size="tiny" @click="handleOpenFolder" title="打开文件夹">
  <template #icon>
    <NIcon><FolderOpenOutline /></NIcon>
  </template>
</NButton>
```

#### 3. FileExplorer 空状态
文件: `/frontend/src/components/layout/FileExplorer.vue`

```vue
<!-- 未打开文件夹时显示提示 -->
<div v-if="!hasWorkspace" class="no-workspace">
  <NEmpty description="未打开文件夹">
    <template #extra>
      <NButton size="small" type="primary" @click="openFolderDialog">
        <template #icon>
          <NIcon><FolderOpenOutline /></NIcon>
        </template>
        打开文件夹
      </NButton>
    </template>
  </NEmpty>
</div>
```

---

## 📝 使用流程

### 方式一：文件菜单

1. 点击顶部菜单栏的"文件"
2. 选择"打开文件夹..."
3. 在系统对话框中选择目标文件夹
4. 点击"选择"或"打开"
5. 等待加载完成，查看文件树

### 方式二：快捷按钮

1. 点击标题栏右侧的文件夹图标
2. 选择目标文件夹
3. 自动加载

### 方式三：空状态

1. 打开应用，看到"未打开文件夹"提示
2. 点击"打开文件夹"按钮
3. 选择目标文件夹
4. 开始使用

---

## 🎨 UI/UX 设计

### 1. 标题栏改进
- **macOS 风格交通灯按钮** - 红、黄、绿三个圆形按钮
- **菜单栏** - 文件、编辑、选择、查看等 8 个菜单项
- **快捷按钮** - 右侧的文件夹打开按钮

### 2. 加载状态
```
用户操作
  ↓
显示 "正在打开文件夹选择对话框..." (loading)
  ↓
系统对话框出现
  ↓
用户选择文件夹
  ↓
显示 "正在加载文件夹..." (loading)
  ↓
验证文件夹
  ↓
加载文件树
  ↓
显示成功提示 "已打开: 文件夹名"
```

### 3. 错误处理
- **取消操作** - 显示"已取消选择"
- **无法访问** - 显示"无法访问该文件夹"
- **其他错误** - 显示具体错误信息

---

## 🔍 核心代码

### 打开文件夹完整流程

```typescript
// TitleBar.vue 中的处理函数
async function handleOpenFolder() {
  try {
    message.loading('正在打开文件夹选择对话框...', { duration: 0 })
    
    // 1. 调用后端打开文件夹对话框
    const selectedPath = await OpenFolderDialog()
    
    message.destroyAll()
    
    if (!selectedPath) {
      message.info('已取消选择')
      return
    }
    
    message.loading('正在加载文件夹...', { duration: 0 })
    
    // 2. 验证文件夹
    try {
      await ListDir(selectedPath)
    } catch (error) {
      message.destroyAll()
      message.error('无法访问该文件夹')
      return
    }
    
    // 3. 设置工作区
    editorStore.setWorkspace(selectedPath)
    
    message.destroyAll()
    message.success(`已打开: ${selectedPath.split('/').pop()}`)
    
  } catch (error) {
    message.destroyAll()
    console.error('Failed to open folder:', error)
    const errorMsg = error instanceof Error ? error.message : String(error)
    if (!errorMsg.includes('cancelled')) {
      message.error(`打开文件夹失败: ${errorMsg}`)
    } else {
      message.info('已取消选择')
    }
  }
}
```

### FileExplorer 加载文件

```typescript
async function loadFiles(rootPath?: string) {
  loading.value = true
  try {
    // 使用指定路径或当前工作区
    const root = rootPath || await GetProjectRoot()
    projectRoot.value = root
    
    // 读取目录内容
    const files = await ListDir(root)
    treeData.value = convertToTree(files)
  } catch (error) {
    console.error('Failed to load files:', error)
    message.error('加载文件失败')
  } finally {
    loading.value = false
  }
}
```

---

## 📊 文件修改清单

### 后端文件
1. ✏️ `/backend/file_service.go` - 添加 OpenFolderDialog 和 SetProjectRoot
2. ✏️ `/backend/interfaces.go` - 更新 IFileSystemService 接口
3. ✏️ `/backend/app_service.go` - 添加方法委托
4. ✏️ `/backend/app.go` - 暴露方法给前端

### 前端文件
1. ✏️ `/frontend/src/stores/editor.ts` - 添加 Workspace 状态管理
2. ✏️ `/frontend/src/components/layout/TitleBar.vue` - 添加文件菜单和打开按钮
3. ✏️ `/frontend/src/components/layout/FileExplorer.vue` - 支持工作区状态

---

## 🚀 未来增强

### 计划功能

- [ ] **最近打开的文件夹** - 记录并显示最近打开的文件夹列表
- [ ] **多工作区** - 支持同时打开多个文件夹
- [ ] **工作区设置** - 保存工作区特定配置
- [ ] **拖拽打开** - 支持拖拽文件夹到应用窗口打开
- [ ] **快捷键** - 添加 Cmd+O / Ctrl+O 快捷键
- [ ] **自动恢复** - 记住上次打开的文件夹，启动时自动恢复

---

## 🎯 对比 VSCode

| 特性 | VSCode | Hao-Code | 状态 |
|------|--------|----------|------|
| 打开文件夹 | ✅ | ✅ | 已实现 |
| 系统原生对话框 | ✅ | ✅ | 已实现 |
| 跨平台支持 | ✅ | ✅ | 已实现 |
| 文件菜单 | ✅ | ✅ | 已实现 |
| 快捷按钮 | ✅ | ✅ | 已实现 |
| 空状态提示 | ✅ | ✅ | 已实现 |
| 最近文件夹 | ✅ | ⏳ | 待实现 |
| 多工作区 | ✅ | ⏳ | 待实现 |
| 快捷键 | ✅ | ⏳ | 待实现 |

---

## 💡 使用技巧

### 快速切换文件夹

1. 使用菜单栏 "文件" > "打开文件夹..."
2. 或者点击标题栏右侧的文件夹图标
3. 选择新的文件夹后会自动替换当前工作区

### 重新加载文件夹

如果文件夹内容有变化：
1. 点击菜单栏 "文件" > "重新加载文件夹"
2. 或者点击文件树工具栏的刷新按钮

### 注意事项

- 首次打开时会请求 macOS 的"文件和文件夹"权限
- 隐藏的 `node_modules` 和 `.` 开头的文件会被自动过滤
- 大文件夹加载可能需要几秒时间

---

## 🎊 总结

"打开文件夹"功能已经完整实现，提供了：

1. ✅ **三种打开方式** - 菜单、按钮、空状态
2. ✅ **跨平台支持** - macOS/Windows/Linux
3. ✅ **完整工作流** - 选择 → 验证 → 加载 → 提示
4. ✅ **状态管理** - 工作区信息集中管理
5. ✅ **用户体验** - 加载状态、错误提示、成功反馈

现在你可以像使用 VSCode 一样，自由选择任何文件夹进行代码编辑！

---

**更新时间**：2026-04-14  
**版本**：v1.3