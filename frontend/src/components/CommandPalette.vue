<template>
  <NModal
    v-model:show="visible"
    preset="card"
    :style="{ width: '600px', maxWidth: '90vw' }"
    :closable="false"
    class="command-palette-modal"
  >
    <template #header>
      <div class="modal-header">
        <span>命令面板</span>
        <span class="hint">按 ↑↓ 选择，Enter 执行，Esc 关闭</span>
      </div>
    </template>
    
    <!-- 命令输入框 -->
    <NInput
      ref="inputRef"
      v-model:value="commandInput"
      placeholder="输入命令..."
      size="large"
      @input="filterCommands"
      @keydown="handleKeydown"
    >
      <template #prefix>
        <n-icon :component="SearchIcon" />
      </template>
    </NInput>

    <!-- 命令列表 -->
    <div class="command-list" v-if="filteredCommands.length > 0">
      <div 
        v-for="(cmd, index) in filteredCommands" 
        :key="cmd.id"
        :class="['command-item', { selected: selectedIndex === index }]"
        @click="executeCommand(cmd)"
        @mouseenter="selectedIndex = index"
      >
        <div class="command-info">
          <span class="command-label">{{ cmd.label }}</span>
          <span class="command-category" v-if="cmd.category">{{ cmd.category }}</span>
        </div>
        <div class="command-shortcut" v-if="cmd.shortcut">
          {{ formatShortcut(cmd.shortcut) }}
        </div>
      </div>
    </div>

    <!-- 无结果提示 -->
    <div class="no-commands" v-else>
      <p>未找到匹配的命令</p>
    </div>
  </NModal>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { NModal, NInput, NIcon } from 'naive-ui'
import { Search as SearchIcon } from '@vicons/ionicons5'
import { useEditorStore } from '@/stores/editor'
import { useMessage } from 'naive-ui'

const editorStore = useEditorStore()
const message = useMessage()

const visible = ref(false)
const commandInput = ref('')
const selectedIndex = ref(0)
const inputRef = ref<any>(null)

interface Command {
  id: string
  label: string
  category?: string
  shortcut?: string
  action: () => void | Promise<void>
}

// 命令列表
const commands: Command[] = [
  // 文件操作
  {
    id: 'file:new-text',
    label: '新建文本文件',
    category: '文件',
    shortcut: 'Ctrl+N',
    action: () => {
      const untitledPath = `untitled-${Date.now()}.txt`
      editorStore.openFile(untitledPath, '')
      message.success('已创建新文本文件')
    }
  },
  {
    id: 'file:open',
    label: '打开文件...',
    category: '文件',
    shortcut: 'Ctrl+O',
    action: () => {
      window.dispatchEvent(new CustomEvent('menu:open-file'))
    }
  },
  {
    id: 'file:open-folder',
    label: '打开文件夹...',
    category: '文件',
    shortcut: 'Ctrl+K Ctrl+O',
    action: () => {
      window.dispatchEvent(new CustomEvent('menu:open-folder'))
    }
  },
  {
    id: 'file:save',
    label: '保存',
    category: '文件',
    shortcut: 'Ctrl+S',
    action: () => {
      window.dispatchEvent(new CustomEvent('menu:save'))
    }
  },
  {
    id: 'file:save-as',
    label: '另存为...',
    category: '文件',
    shortcut: 'Ctrl+Shift+S',
    action: () => {
      window.dispatchEvent(new CustomEvent('menu:save-as'))
    }
  },
  {
    id: 'file:close',
    label: '关闭编辑器',
    category: '文件',
    action: () => {
      window.dispatchEvent(new CustomEvent('menu:close-editor'))
    }
  },
  
  // 编辑操作
  {
    id: 'edit:undo',
    label: '撤销',
    category: '编辑',
    shortcut: 'Ctrl+Z',
    action: () => message.info('使用编辑器内置撤销功能')
  },
  {
    id: 'edit:redo',
    label: '重做',
    category: '编辑',
    shortcut: 'Ctrl+Shift+Z',
    action: () => message.info('使用编辑器内置重做功能')
  },
  {
    id: 'edit:find',
    label: '查找',
    category: '编辑',
    shortcut: 'Ctrl+F',
    action: () => message.info('使用编辑器内置查找功能 (Ctrl+F)')
  },
  {
    id: 'edit:replace',
    label: '替换',
    category: '编辑',
    shortcut: 'Ctrl+H',
    action: () => message.info('使用编辑器内置替换功能 (Ctrl+H)')
  },
  
  // 视图操作
  {
    id: 'view:toggle-sidebar',
    label: '切换侧边栏',
    category: '视图',
    shortcut: 'Ctrl+B',
    action: () => {
      editorStore.toggleSidebar()
      message.success(editorStore.sidebarVisible ? '侧边栏已显示' : '侧边栏已隐藏')
    }
  },
  {
    id: 'view:toggle-autosave',
    label: '切换自动保存',
    category: '视图',
    action: () => {
      editorStore.toggleAutoSave()
      message.success(editorStore.autoSave ? '自动保存已启用' : '自动保存已禁用')
    }
  },
  {
    id: 'view:explorer',
    label: '显示资源管理器',
    category: '视图',
    action: () => {
      editorStore.setSidebarView('explorer')
    }
  },
  {
    id: 'view:search',
    label: '显示搜索',
    category: '视图',
    action: () => {
      editorStore.setSidebarView('search')
    }
  },
  {
    id: 'view:git',
    label: '显示源代码管理',
    category: '视图',
    action: () => {
      editorStore.setSidebarView('git')
    }
  },
  
  // Git 操作
  {
    id: 'git:commit',
    label: 'Git 提交',
    category: 'Git',
    action: () => {
      message.info('Git 提交功能开发中...')
    }
  },
  {
    id: 'git:push',
    label: 'Git 推送',
    category: 'Git',
    action: () => {
      message.info('Git 推送功能开发中...')
    }
  },
  {
    id: 'git:pull',
    label: 'Git 拉取',
    category: 'Git',
    action: () => {
      message.info('Git 拉取功能开发中...')
    }
  },
  
  // 帮助
  {
    id: 'help:welcome',
    label: '欢迎',
    category: '帮助',
    action: () => {
      window.dispatchEvent(new CustomEvent('menu:welcome'))
    }
  },
  {
    id: 'help:shortcuts',
    label: '键盘快捷方式参考',
    category: '帮助',
    action: () => {
      window.dispatchEvent(new CustomEvent('menu:keyboard-shortcuts'))
    }
  },
  {
    id: 'help:devtools',
    label: '切换开发人员工具',
    category: '帮助',
    action: () => {
      window.dispatchEvent(new CustomEvent('menu:toggle-devtools'))
    }
  }
]

const filteredCommands = ref<Command[]>(commands)

// 过滤命令
function filterCommands() {
  if (!commandInput.value.trim()) {
    filteredCommands.value = commands
  } else {
    const query = commandInput.value.toLowerCase()
    filteredCommands.value = commands.filter(cmd =>
      cmd.label.toLowerCase().includes(query) ||
      cmd.category?.toLowerCase().includes(query) ||
      cmd.id.toLowerCase().includes(query)
    )
  }
  selectedIndex.value = 0
}

// 执行命令
async function executeCommand(cmd: Command) {
  try {
    await cmd.action()
    visible.value = false
    commandInput.value = ''
  } catch (error) {
    console.error('Command execution failed:', error)
    message.error('命令执行失败')
  }
}

// 键盘导航
function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    selectedIndex.value = Math.min(selectedIndex.value + 1, filteredCommands.value.length - 1)
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    selectedIndex.value = Math.max(selectedIndex.value - 1, 0)
  } else if (e.key === 'Enter') {
    e.preventDefault()
    if (filteredCommands.value[selectedIndex.value]) {
      executeCommand(filteredCommands.value[selectedIndex.value])
    }
  } else if (e.key === 'Escape') {
    e.preventDefault()
    visible.value = false
  }
}

// 格式化快捷键显示
function formatShortcut(shortcut: string): string {
  return shortcut
    .replace('Ctrl', '⌃')
    .replace('Shift', '⇧')
    .replace('Alt', '⌥')
    .replace('Meta', '⌘')
}

// 显示命令面板
function show() {
  visible.value = true
  commandInput.value = ''
  filteredCommands.value = commands
  selectedIndex.value = 0
  
  // 聚焦输入框
  nextTick(() => {
    inputRef.value?.focus()
  })
}

// 监听全局快捷键
function handleGlobalKeydown(e: KeyboardEvent) {
  // Ctrl+Shift+P 或 Cmd+Shift+P
  if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key === 'p') {
    e.preventDefault()
    show()
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleGlobalKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleGlobalKeydown)
})

// 暴露方法供外部调用
defineExpose({ show })
</script>

<style scoped>
.command-palette-modal :deep(.n-card__content) {
  padding: 0;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.hint {
  font-size: 12px;
  color: #858585;
}

.command-list {
  max-height: 400px;
  overflow-y: auto;
  margin-top: 12px;
}

.command-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  cursor: pointer;
  border-radius: 6px;
  margin-bottom: 4px;
  transition: all 0.2s;
}

.command-item:hover,
.command-item.selected {
  background-color: #094771;
}

.command-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.command-label {
  font-size: 14px;
  color: #CCCCCC;
}

.command-category {
  font-size: 11px;
  color: #858585;
  background-color: #3C3C3C;
  padding: 2px 6px;
  border-radius: 3px;
}

.command-shortcut {
  font-size: 12px;
  color: #858585;
  font-family: monospace;
}

.no-commands {
  text-align: center;
  padding: 40px 20px;
  color: #858585;
}
</style>
