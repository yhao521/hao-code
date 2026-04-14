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
    >
      <NTabPane
        v-for="tab in editorStore.tabs"
        :key="tab.id"
        :name="tab.id"
        :tab="tab.name + (tab.dirty ? ' •' : '')"
      >
        <!-- Monaco Editor 容器 -->
        <div ref="editorContainer" class="monaco-container"></div>
      </NTabPane>
    </NTabs>
    
    <!-- 空状态 -->
    <div v-else class="empty-state">
      <div class="empty-content">
        <h2>Hao-Code Editor</h2>
        <p>打开文件开始编辑</p>
        <div class="shortcuts">
          <div><kbd>Ctrl+P</kbd> 快速打开</div>
          <div><kbd>Ctrl+Shift+P</kbd> 命令面板</div>
          <div><kbd>Ctrl+B</kbd> 切换侧边栏</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { NTabs, NTabPane } from 'naive-ui'
import * as monaco from 'monaco-editor'
import { useEditorStore } from '@/stores/editor'

const editorStore = useEditorStore()
const editorContainer = ref<HTMLElement | null>(null)
let editor: monaco.editor.IStandaloneCodeEditor | null = null

function handleTabChange(tabId: string) {
  editorStore.activeEditor = tabId
  // TODO: 更新编辑器内容
}

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
      }
    })
  }
})

// 监听活动标签页变化
watch(() => editorStore.activeEditor, (newTabId) => {
  if (newTabId && editor) {
    const tab = editorStore.tabs.find(t => t.id === newTabId)
    if (tab && tab.content) {
      editor.setValue(tab.content)
    }
  }
})
</script>

<style scoped>
.editor-area {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.monaco-container {
  height: calc(100vh - 90px);
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #1E1E1E;
}

.empty-content {
  text-align: center;
  color: #888;
}

.empty-content h2 {
  color: #CCCCCC;
  margin-bottom: 16px;
}

.shortcuts {
  margin-top: 24px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.shortcuts div {
  font-size: 13px;
}

kbd {
  background-color: #3C3C3C;
  padding: 2px 6px;
  border-radius: 3px;
  border: 1px solid #555;
  font-family: monospace;
  font-size: 12px;
}
</style>
