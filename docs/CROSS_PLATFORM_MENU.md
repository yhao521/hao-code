# 跨平台菜单栏实现方案

## 🎯 设计理念

根据不同操作系统的 UI 规范，采用统一和差异化相结合的菜单栏实现方案：

| 平台 | 菜单栏 | 窗口样式 | 窗口控制 |
|------|--------|----------|----------|
| **macOS** | Wails 原生菜单（系统顶部） | 有边框窗口 | 系统交通灯按钮 |
| **Windows** | Wails 原生菜单（系统顶部） | 无边框窗口 | 前端自定义按钮（右上角） |

---

## 📐 架构设计

### 1. 后端配置（main.go）

```go
// 检测操作系统
isMacOS := goruntime.GOOS == "darwin"

// 统一使用系统菜单栏（所有平台）
// 通过禁用 Frameless 来启用系统菜单

// macOS: 有边框窗口
Mac: &mac.Options{
    TitleBar: mac.TitleBarDefault(), // 使用默认标题栏（有边框）
    Frameless: false,                // 有边框
}

// Windows: 无边框窗口（前端自定义控制按钮）
Windows: &windows.Options{
    DisableFramelessWindowDecorations: true, // 无边框
}
Frameless: !isMacOS, // macOS: false, Windows: true
```

### 2. 前端适配（TitleBar.vue）

```vue
<!-- 所有平台都显示标题栏 -->
<div class="titlebar" :class="{ 'macos': isMacOS }">
  <!-- 工作区名称 -->
  <div class="titlebar-center">
    <span>{{ workspaceName }}</span>
  </div>
  
  <!-- 功能按钮 -->
  <div class="titlebar-right">
    <!-- 打开文件夹按钮 -->
    <NButton @click="handleOpenFolder">
      <FolderOpenOutline />
    </NButton>
    
    <!-- Windows: 自定义窗口控制按钮 -->
    <div v-if="!isMacOS" class="window-controls">
      <div class="control-btn minimize" @click="minimizeWindow">—</div>
      <div class="control-btn maximize" @click="maximizeWindow">□</div>
      <div class="control-btn close" @click="closeWindow">✕</div>
    </div>
  </div>
</div>
```

### 3. 窗口控制实现

```typescript
import * as wailsRuntime from '@wails/runtime/runtime'

// Windows 窗口控制
function minimizeWindow() {
  wailsRuntime.WindowMinimise()
}

function maximizeWindow() {
  wailsRuntime.WindowToggleMaximise()
}

function closeWindow() {
  wailsRuntime.WindowClose()
}
```

---

## 🎨 UI 对比

### macOS
```
┌─────────────────────────────────────────────┐
│ ● ○ ●  Hao-Code Editor              [📁]  │ ← 系统菜单栏 + 交通灯
├─────────────────────────────────────────────┤
│  侧边栏  │  编辑器区域                     │
│         │                                  │
└─────────┴──────────────────────────────────┘
```

### Windows
```
┌─────────────────────────────────────────────┐
│  Hao-Code Editor             [📁] [—] [□] [✕]│ ← 自定义标题栏
├─────────────────────────────────────────────┤
│  侧边栏  │  编辑器区域                     │
│         │                                  │
└─────────┴──────────────────────────────────┘
```

---

## 🔧 实现细节

### macOS 方案

#### 特点
- ✅ **有边框窗口** - `Frameless: false` + `mac.TitleBarDefault()`
- ✅ **系统交通灯** - 红、黄、绿三个按钮（左上角）
- ✅ **系统菜单栏** - Wails 原生菜单显示在屏幕顶部
- ✅ **标准窗口行为** - 拖拽、缩放、全屏由系统管理

#### 配置
```go
Mac: &mac.Options{
    TitleBar: mac.TitleBarDefault(), // 默认标题栏
}
Frameless: false // 有边框
```

#### 前端适配
```vue
<!-- 所有平台都显示标题栏 -->
<TitleBar />

<!-- 标题栏内容 -->
<div class="titlebar macos">
  <div class="titlebar-center">
    <span>{{ workspaceName }}</span>
  </div>
  <div class="titlebar-right">
    <!-- 只显示功能按钮，不显示窗口控制 -->
    <NButton @click="handleOpenFolder">
      <FolderOpenOutline />
    </NButton>
  </div>
</div>
```

---

### Windows 方案

#### 特点
- ✅ **无边框窗口** - `Frameless: true`
- ✅ **系统菜单栏** - Wails 原生菜单显示在窗口顶部
- ✅ **自定义窗口控制** - 前端实现最小化、最大化、关闭按钮
- ✅ **可拖拽标题栏** - `-webkit-app-region: drag`

#### 配置
```go
Windows: &windows.Options{
    DisableFramelessWindowDecorations: true, // 无边框
}
Frameless: true // 无边框
```

#### 前端实现
```vue
<div class="titlebar">
  <!-- 可拖拽区域 -->
  <div class="titlebar-center" style="-webkit-app-region: drag">
    <span>{{ workspaceName }}</span>
  </div>
  
  <!-- 窗口控制按钮（不可拖拽） -->
  <div class="window-controls" style="-webkit-app-region: no-drag">
    <div class="control-btn minimize" @click="minimizeWindow">—</div>
    <div class="control-btn maximize" @click="maximizeWindow">□</div>
    <div class="control-btn close" @click="closeWindow">✕</div>
  </div>
</div>
```

#### 窗口控制样式
```
.window-controls {
  display: flex;
  margin-left: 8px;
  border-left: 1px solid #3E3E42;
  padding-left: 8px;
}

.control-btn {
  width: 36px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background-color 0.2s;
}

.control-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.control-btn.close:hover {
  background-color: #E81123; /* Windows 风格红色 */
}
```

---

## 📝 实现清单

### 后端配置
- [x] `main.go` - 跨平台窗口配置
  - macOS: `Frameless: false`（有边框）
  - Windows: `Frameless: true`（无边框）
  - 所有平台: 使用系统菜单栏

### 前端实现
- [x] `TitleBar.vue` - 统一的标题栏组件
  - macOS: 显示工作区名称 + 功能按钮
  - Windows: 工作区名称 + 功能按钮 + 窗口控制
- [x] `TitleBar.vue` - Windows 窗口控制按钮
  - 最小化: `WindowMinimise()`
  - 最大化: `WindowToggleMaximise()`
  - 关闭: `WindowClose()`
- [x] `TitleBar.vue` - 可拖拽区域
  - 标题栏: `-webkit-app-region: drag`
  - 按钮: `-webkit-app-region: no-drag`

---

## 🚀 窗口控制功能

### Windows 按钮功能

#### 最小化
``typescript
function minimizeWindow() {
  wailsRuntime.WindowMinimise()
}
```
- 效果：窗口最小化到任务栏
- API：`WindowMinimise()`

#### 最大化/还原
``typescript
function maximizeWindow() {
  wailsRuntime.WindowToggleMaximise()
}
```
- 效果：切换最大化和窗口模式
- API：`WindowToggleMaximise()`
- 双击标题栏也可触发

#### 关闭（退出应用）
``typescript
function closeWindow() {
  wailsRuntime.Quit()
}
```
- 效果：完全退出应用程序
- API：`Quit()`
- 悬停效果：红色背景（#E81123）

---

## 🎯 菜单使用

### macOS
1. 菜单栏显示在屏幕顶部
2. 点击"文件"、"编辑"等菜单项
3. 使用系统交通灯控制窗口

### Windows
1. 菜单栏显示在窗口顶部
2. 点击菜单项操作
3. 使用自定义按钮控制窗口

---

## 🔮 未来增强

### macOS
- [ ] 自定义系统菜单项（文件、编辑等）
- [ ] 添加快捷键支持（Cmd+O, Cmd+S）
- [ ] 集成 macOS 原生对话框

### Windows
- [ ] 优化窗口控制按钮样式
- [ ] 添加最大化/还原状态图标切换
- [ ] 实现双击标题栏最大化

### 通用
- [ ] 菜单与前端状态联动
- [ ] 最近文件列表
- [ ] 动态菜单项
- [ ] 菜单栏键盘导航

---

## 💡 技术要点

### 1. 窗口配置
```go
// 根据平台设置 Frameless
Frameless: !isMacOS // macOS: false, Windows: true
```

### 2. 可拖拽区域
```css
/* 整个标题栏可拖拽 */
.titlebar:not(.macos) {
  -webkit-app-region: drag;
}

/* 按钮区域不可拖拽 */
.window-controls {
  -webkit-app-region: no-drag;
}
```

### 3. 窗口控制 API
```typescript
import * as wailsRuntime from '@wails/runtime/runtime'

wailsRuntime.WindowMinimise()      // 最小化
wailsRuntime.WindowToggleMaximise() // 最大化/还原
wailsRuntime.WindowClose()          // 关闭
```

### 4. 平台检测
```typescript
const isMacOS = computed(() => {
  return navigator.platform.toLowerCase().includes('mac')
})
```

---

## 📊 对比其他编辑器

| 特性 | VSCode | Sublime | Hao-Code |
|------|--------|---------|----------|
| macOS 有边框 | ✅ | ✅ | ✅ |
| macOS 系统菜单 | ✅ | ✅ | ✅ |
| Windows 无边框 | ✅ | ❌ | ✅ |
| Windows 自定义按钮 | ✅ | ❌ | ✅ |
| Windows 系统菜单 | ✅ | ✅ | ✅ |
| 窗口拖拽 | ✅ | ✅ | ✅ |

---

## ⚠️ 注意事项

1. **Windows 无边框** - 使用 `Frameless: true` 时需自定义窗口控制
2. **可拖拽区域** - 按钮必须设置 `-webkit-app-region: no-drag`
3. **macOS 交通灯** - 系统自动管理，前端无需实现
4. **菜单统一性** - 所有平台使用 Wails 原生菜单
5. **窗口状态** - 最大化状态需要前端跟踪

---

**更新时间**：2026-04-14  
**版本**：v1.5