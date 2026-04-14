<template>
  <div class="editor-area">
    <!-- 标签页 -->
    <NTabs
      v-if="editorStore.tabs.length > 0"
      :value="editorStore.activeEditor || undefined"
      type="card"
      closable
      @close="(tabId: string) => editorStore.closeTab(tabId)"
      @update:value="handleTabChange"
      class="editor-tabs"
    >
      <NTabPane
        v-for="tab in editorStore.tabs"
        :key="tab.id"
        :name="tab.id"
      >
        <template #tab>
          <div class="tab-label">
            <span :class="{ dirty: tab.dirty }">{{ tab.name }}</span>
            <span v-if="tab.dirty" class="dirty-indicator">•</span>
          </div>
        </template>
      </NTabPane>
    </NTabs>
    
    <!-- Monaco Editor 容器 -->
    <div ref="editorContainer" class="monaco-container"></div>
    
    <!-- 空状态 -->
    <div v-if="editorStore.tabs.length === 0" class="empty-state">
      <div class="empty-content">
        <h2>Hao-Code Editor</h2>
        <p>打开文件开始编辑</p>
        <div class="shortcuts">
          <div><kbd>Ctrl+P</kbd> 快速打开</div>
          <div><kbd>Ctrl+S</kbd> 保存文件</div>
          <div><kbd>Ctrl+B</kbd> 切换侧边栏</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { NTabs, NTabPane, useMessage } from 'naive-ui'
import * as monaco from 'monaco-editor'
import { useEditorStore } from '@/stores/editor'
import { WriteFile } from '@wails/go/backend/App'

const editorStore = useEditorStore()
const message = useMessage()
const editorContainer = ref<HTMLElement | null>(null)
let editor: monaco.editor.IStandaloneCodeEditor | null = null
let currentModel: monaco.editor.ITextModel | null = null

// 初始化编辑器
onMounted(() => {
  if (editorContainer.value && !editor) {
    editor = monaco.editor.create(editorContainer.value, {
      value: '// Welcome to Hao-Code Editor\n',
      language: 'typescript',
      theme: 'vs-dark',
      automaticLayout: true,
      fontSize: 14,
      fontFamily: "'Fira Code', 'Cascadia Code', 'Source Code Pro', Consolas, 'Courier New', monospace",
      minimap: {
        enabled: true
      },
      scrollBeyondLastLine: false,
      renderWhitespace: 'selection',
      bracketPairColorization: {
        enabled: true
      },
      tabSize: 2,
      insertSpaces: true,
      wordWrap: 'on'
    })
    
    // 监听内容变化
    editor.onDidChangeModelContent(() => {
      if (editorStore.activeEditor) {
        const content = editor!.getValue()
        editorStore.updateContent(editorStore.activeEditor, content)
      }
    })
    
    // 注册保存快捷键
    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS, () => {
      if (editorStore.activeEditor) {
        handleSave(editorStore.activeEditor)
      }
    })
  }
})

// 监听标签页变化
watch(() => editorStore.activeEditor, (newTabId) => {
  if (newTabId && editor) {
    const tab = editorStore.tabs.find(t => t.id === newTabId)
    if (tab) {
      loadFileIntoEditor(tab)
    }
  }
})

// 将文件加载到编辑器
function loadFileIntoEditor(tab: any) {
  const language = tab.language || 'plaintext'
  
  // 设置语言
  monaco.editor.setModelLanguage(editor!.getModel()!, language)
  
  // 设置内容
  if (tab.content !== undefined) {
    editor!.setValue(tab.content)
  }
  
  // 滚动到顶部
  editor!.setScrollPosition({ scrollTop: 0 })
}

// 处理标签页切换
function handleTabChange(tabId: string) {
  editorStore.activeEditor = tabId
}

// 保存文件
async function handleSave(tabId: string) {
  const tab = editorStore.tabs.find(t => t.id === tabId)
  if (!tab || !tab.content) return
  
  try {
    await WriteFile(tab.path, tab.content)
    editorStore.saveFile(tabId)
    message.success(`已保存: ${tab.name}`)
  } catch (error) {
    message.error(`保存失败: ${error}`)
  }
}

// 清理
onUnmounted(() => {
  if (editor) {
    editor.dispose()
  }
})
</script>

<style scoped>
.editor-area {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #1E1E1E;
}

/* 标签页样式优化 */
:deep(.editor-tabs) {
  background-color: #252526;
  border-bottom: 1px solid #3E3E42;
  min-height: 35px;
}

:deep(.n-tabs .n-tabs-nav) {
  background-color: #252526 !important;
}

:deep(.n-tabs .n-tabs-tab) {
  background-color: #2D2D2D !important;
  color: #969696 !important;
  border-color: #252526 !important;
  padding: 6px 12px !important;
  font-size: 13px !important;
  margin-right: 1px !important;
}

:deep(.n-tabs .n-tabs-tab:hover) {
  background-color: #2D2D2D !important;
  color: #CCCCCC !important;
}

:deep(.n-tabs .n-tabs-tab--active) {
  background-color: #1E1E1E !important;
  color: #FFFFFF !important;
}

/* Monaco 编辑器容器 */
.monaco-container {
  flex: 1;
  overflow: hidden;
  background-color: #1E1E1E;
}

/* 空状态 */
.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #1E1E1E;
}

.empty-content {
  text-align: center;
  color: #858585;
}

.empty-content h2 {
  font-size: 28px;
  font-weight: 300;
  color: #CCCCCC;
  margin-bottom: 16px;
}

.empty-content p {
  font-size: 14px;
  color: #858585;
}

.shortcuts {
  margin-top: 24px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
}

.shortcuts div {
  font-size: 13px;
  color: #858585;
}

kbd {
  background-color: #3C3C3C;
  padding: 2px 6px;
  border-radius: 3px;
  border: 1px solid #555;
  font-family: monospace;
  font-size: 12px;
  margin: 0 2px;
  color: #CCCCCC;
}

.tab-label {
  display: flex;
  align-items: center;
  gap: 4px;
}

.dirty-indicator {
  color: #4EC9B0;
  font-size: 16px;
  line-height: 1;
}
</style>
