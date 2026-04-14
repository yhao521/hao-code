# 可拖拽调整宽度功能说明

## 🎯 功能概述

实现了类似 VSCode 的侧边栏宽度可拖拽调整功能，用户可以通过鼠标拖拽侧边栏和编辑器区域之间的分割线来调整宽度。

---

## ✨ 功能特性

### 1. **拖拽调整宽度**
- ✅ 鼠标悬停在分割线上时显示高亮效果
- ✅ 拖拽时显示拖拽状态
- ✅ 拖拽结束后保持新宽度
- ✅ 光标自动变为调整大小光标

### 2. **宽度限制**
- ✅ **最小宽度**: 180px（防止侧边栏过窄）
- ✅ **最大宽度**: 500px（防止编辑器区域过小）
- ✅ 自动限制拖拽范围

### 3. **视觉反馈**
- ✅ **默认状态**: 灰色细线（#3E3E42）
- ✅ **悬停状态**: 蓝色高亮（#0E639C）
- ✅ **拖拽状态**: 深蓝色背景 + 加粗线条
- ✅ **光标变化**: col-resize 光标

### 4. **用户体验**
- ✅ 拖拽时禁止文本选择
- ✅ 拖拽结束后恢复文本选择
- ✅ 平滑过渡动画
- ✅ 防止误操作

---

## 🔧 技术实现

### 组件结构

```
App.vue
└── ResizableSplit.vue (自定义组件)
    ├── Panel 1 (侧边栏)
    ├── Resize Handle (拖拽条)
    └── Panel 2 (编辑器区域)
```

### ResizableSplit 组件 API

#### Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `min` | number | 150 | 最小宽度（像素） |
| `max` | number | 600 | 最大宽度（像素） |
| `horizontal` | boolean | true | 是否水平分割（true=左右，false=上下） |

#### Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| `update:size` | `value: number` | 宽度变化时触发 |

#### Expose

| 方法/属性 | 类型 | 说明 |
|-----------|------|------|
| `size` | `ref<number>` | 当前宽度值 |

### 核心实现

#### 拖拽逻辑

```typescript
// 1. 开始拖拽
function startResize(e: MouseEvent) {
  isDragging.value = true
  startPos.value = e.clientX
  startSize.value = size.value
  
  document.addEventListener('mousemove', handleResize)
  document.addEventListener('mouseup', stopResize)
  document.body.style.cursor = 'col-resize'
  document.body.style.userSelect = 'none'
}

// 2. 处理拖拽
function handleResize(e: MouseEvent) {
  const delta = e.clientX - startPos.value
  let newSize = startSize.value + delta
  
  // 限制范围
  newSize = Math.max(props.min, Math.min(props.max, newSize))
  size.value = newSize
}

// 3. 停止拖拽
function stopResize() {
  isDragging.value = false
  document.removeEventListener('mousemove', handleResize)
  document.removeEventListener('mouseup', stopResize)
  document.body.style.cursor = ''
  document.body.style.userSelect = ''
}
```

#### 样式设计

```css
/* 拖拽条样式 */
.resize-handle {
  position: relative;
  width: 4px;
  cursor: col-resize;
  transition: background-color 0.2s;
}

/* 悬停效果 */
.resize-handle--hover {
  background-color: rgba(14, 99, 156, 0.2);
}

/* 拖拽效果 */
.resize-handle--dragging {
  background-color: rgba(14, 99, 156, 0.3);
}

/* 分割线 */
.resize-handle-line {
  width: 1px;
  height: 40px;
  background-color: #3E3E42;
  transition: background-color 0.2s, width 0.2s;
}

.resize-handle--hover .resize-handle-line {
  background-color: #0E639C;
  width: 2px;
}
```

---

## 📝 使用示例

### 基础用法

```vue
<template>
  <ResizableSplit :min="180" :max="500">
    <template #1>
      <div>左侧面板内容</div>
    </template>
    <template #2>
      <div>右侧面板内容</div>
    </template>
  </ResizableSplit>
</template>

<script setup>
import ResizableSplit from './components/layout/ResizableSplit.vue'
</script>
```

### 垂直分割（上下布局）

```vue
<ResizableSplit 
  :min="100" 
  :max="400" 
  :horizontal="false"
>
  <template #1>
    <div>上方面板</div>
  </template>
  <template #2>
    <div>下方面板</div>
  </template>
</ResizableSplit>
```

### 监听宽度变化

```vue
<template>
  <ResizableSplit 
    :min="180" 
    :max="500"
    @update:size="handleSizeChange"
  >
    <!-- ... -->
  </ResizableSplit>
</template>

<script setup>
function handleSizeChange(newSize: number) {
  console.log('新宽度:', newSize)
  // 可以在这里保存用户偏好设置
}
</script>
```

---

## 🎨 样式定制

### 自定义拖拽条颜色

```css
:deep(.resize-handle--hover) {
  background-color: rgba(你的颜色, 0.2);
}

:deep(.resize-handle-line) {
  background-color: 你的颜色;
}
```

### 自定义拖拽条宽度

```vue
<ResizableSplit :min="180" :max="500">
  <!-- 在模板中 -->
</ResizableSplit>

<style>
:deep(.resize-handle) {
  width: 8px !important; /* 更宽的拖拽区域 */
}
</style>
```

---

## 🚀 未来增强

### 计划功能

- [ ] **记忆用户设置** - 保存用户调整的宽度到本地存储
- [ ] **双击重置** - 双击分割线恢复默认宽度
- [ ] **键盘快捷键** - 支持键盘调整宽度
- [ ] **动画过渡** - 宽度变化时的平滑动画
- [ ] **响应式适配** - 根据窗口大小自动调整
- [ ] **多面板支持** - 支持多个可拖拽面板

### 扩展场景

1. **底部面板** - 终端/输出面板的上下拖拽
2. **多栏布局** - 支持三栏或多栏拖拽
3. **嵌套拖拽** - 支持嵌套的可拖拽面板
4. **触摸支持** - 移动端触摸拖拽

---

##  对比 VSCode

| 特性 | VSCode | Hao-Code | 状态 |
|------|--------|----------|------|
| 拖拽调整 | ✅ | ✅ | 已实现 |
| 宽度限制 | ✅ | ✅ | 已实现 |
| 视觉反馈 | ✅ | ✅ | 已实现 |
| 光标变化 | ✅ | ✅ | 已实现 |
| 记忆设置 | ✅ | ⏳ | 待实现 |
| 双击重置 | ✅ | ⏳ | 待实现 |
| 键盘快捷键 | ✅ | ⏳ | 待实现 |

---

## 🎯 总结

可拖拽调整宽度功能已经完整实现，提供了：

1. ✅ **流畅的拖拽体验** - 鼠标拖拽调整侧边栏宽度
2. ✅ **合理的限制** - 最小 180px，最大 500px
3. ✅ **视觉反馈** - 悬停和拖拽时的颜色变化
4. ✅ **用户体验** - 光标变化和防误操作

用户可以像使用 VSCode 一样自由调整侧边栏宽度，提升使用体验！

---

**更新时间**：2026-04-14  
**版本**：v1.2
