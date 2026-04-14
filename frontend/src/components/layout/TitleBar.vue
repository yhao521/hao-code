<template>
  <div class="titlebar" :class="{ 'macos': isMacOS }">
    <!-- 工作区名称 -->
    <div class="titlebar-center">
      <span class="workspace-name">{{ workspaceName }}</span>
    </div>
    
    <!-- 功能按钮 -->
    <div class="titlebar-right">
      <NButton text circle size="tiny" @click="handleOpenFolder" title="打开文件夹">
        <template #icon>
          <NIcon><FolderOpenOutline /></NIcon>
        </template>
      </NButton>
      
      <!-- Windows: 自定义窗口控制按钮 -->
      <div v-if="!isMacOS" class="window-controls">
        <div class="control-btn minimize" @click="minimizeWindow" title="最小化">
          <NIcon size="14"><RemoveOutline /></NIcon>
        </div>
        <div class="control-btn maximize" @click="maximizeWindow" title="最大化">
          <NIcon size="14"><SquareOutline /></NIcon>
        </div>
        <div class="control-btn close" @click="closeWindow" title="关闭">
          <NIcon size="12"><CloseOutline /></NIcon>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { NButton, NIcon, useMessage } from 'naive-ui'
import { 
  FolderOpenOutline,
  RemoveOutline,
  SquareOutline,
  CloseOutline
} from '@vicons/ionicons5'
import { useEditorStore } from '@/stores/editor'
import * as wailsRuntime from '@wails/runtime/runtime'

const editorStore = useEditorStore()
const message = useMessage()
const isMaximized = ref(false)

// 检测是否为 macOS
const isMacOS = computed(() => {
  const platform = navigator.platform.toLowerCase()
  return platform.includes('mac')
})

const workspaceName = computed(() => 
  editorStore.workspace?.name || 'Hao-Code Editor'
)

onMounted(() => {
  // 可以在这里添加窗口状态监听
})

// 打开文件夹
async function handleOpenFolder() {
  try {
    message.loading('正在打开文件夹选择对话框...', { duration: 0 })
    
    // 调用后端打开文件夹对话框
    const selectedPath = await window.go.backend.App.OpenFolderDialog()
    
    message.destroyAll()
    
    if (!selectedPath) {
      message.info('已取消选择')
      return
    }
    
    message.loading('正在加载文件夹...', { duration: 0 })
    
    // 验证文件夹
    try {
      await window.go.backend.App.ListDir(selectedPath)
    } catch (error) {
      message.destroyAll()
      message.error('无法访问该文件夹')
      return
    }
    
    // 设置工作区
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

// 窗口控制函数
function minimizeWindow() {
  wailsRuntime.WindowMinimise()
}

function maximizeWindow() {
  wailsRuntime.WindowToggleMaximise()
  isMaximized.value = !isMaximized.value
}

function closeWindow() {
  wailsRuntime.Quit()
}
</script>

<style scoped>
.titlebar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 38px;
  background-color: #323233;
  color: #CCCCCC;
  user-select: none;
  border-bottom: 1px solid #252526;
  padding: 0 12px;
}

/* Windows 无边框窗口，标题栏可拖拽 */
.titlebar:not(.macos) {
  -webkit-app-region: drag;
}

.titlebar-center {
  flex: 1;
  text-align: center;
}

.titlebar-right {
  display: flex;
  align-items: center;
  gap: 8px;
  -webkit-app-region: no-drag;
}

.workspace-name {
  font-size: 12px;
  font-weight: 500;
}

/* Windows 窗口控制按钮 */
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
  background-color: #E81123;
}

.control-btn.close:hover .n-icon {
  color: white;
}
</style>
