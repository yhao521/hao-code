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
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { NTabs, NTabPane, useMessage } from 'naive-ui'
import * as monaco from 'monaco-editor'
import { useEditorStore } from '@/stores/editor'
import { WriteFile } from '@wails/go/backend/App'

const editorStore = useEditorStore()
const message = useMessage()
const editorContainer = ref<HTMLElement | null>(null)
let editor: monaco.editor.IStandaloneCodeEditor | null = null
let resizeObserver: ResizeObserver | null = null

// 初始化编辑器
onMounted(async () => {
  if (editorContainer.value && !editor) {
    // 等待容器完全渲染
    await nextTick()
    
    editor = monaco.editor.create(editorContainer.value, {
      value: '// Welcome to Hao-Code Editor\n',
      language: 'typescript',
      theme: 'vs-dark',
      automaticLayout: false,  // 禁用自动布局，改为手动控制
      fontSize: 14,
      lineHeight: 22,  // 增加行高，提升可读性
      fontFamily: "'Fira Code', 'Cascadia Code', 'Source Code Pro', Consolas, 'Courier New', monospace",
      fontLigatures: true,  // 启用字体连字
      minimap: {
        enabled: true,
        scale: 1,
        showSlider: 'mouseover',
        renderCharacters: false
      },
      scrollBeyondLastLine: false,
      renderWhitespace: 'selection',
      bracketPairColorization: {
        enabled: true
      },
      guides: {
        bracketPairs: true,
        indentation: true,
        highlightActiveBracketPair: true,
        highlightActiveIndentation: true
      },
      tabSize: 2,
      insertSpaces: true,
      wordWrap: 'on',
      padding: {
        top: 10,  // 顶部内边距
        bottom: 10  // 底部内边距
      },
      cursorBlinking: 'smooth',  // 光标平滑闪烁
      cursorSmoothCaretAnimation: 'on',  // 平滑光标动画
      smoothScrolling: true,  // 平滑滚动
      renderLineHighlight: 'all',  // 高亮当前行
      lineNumbers: 'on',
      lineNumbersMinChars: 2,  // 减小行号最小宽度，从3改为2，节省左侧空间
      glyphMargin: false,  // 禁用字形边距，减少左侧空白
      folding: true,  // 启用代码折叠
      foldingStrategy: 'indentation',
      showFoldingControls: 'mouseover',
      links: true,  // 启用链接检测
      colorDecorators: true,  // 颜色装饰器
      formatOnPaste: true,  // 粘贴时格式化
      formatOnType: true,  // 输入时格式化
      scrollbar: {
        verticalScrollbarSize: 10,  // 垂直滚动条宽度
        horizontalScrollbarSize: 10,  // 水平滚动条宽度
        verticalSliderSize: 10,
        horizontalSliderSize: 10
      },
      // 关键配置：禁用多余的边距
      roundedSelection: false,
      suggestOnTriggerCharacters: true,
      acceptSuggestionOnEnter: 'on',
      tabCompletion: 'on',
      wordBasedSuggestions: 'currentDocument'
    })
    
    // 初始布局
    await nextTick()
    editor.layout()
    
    // 使用 ResizeObserver 监听容器尺寸变化
    resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        // 使用 requestAnimationFrame 确保在下一帧更新
        requestAnimationFrame(() => {
          editor?.layout()
        })
      }
    })
    
    if (editorContainer.value) {
      resizeObserver.observe(editorContainer.value)
    }
    
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
watch(() => editorStore.activeEditor, async (newTabId) => {
  if (newTabId && editor) {
    const tab = editorStore.tabs.find(t => t.id === newTabId)
    if (tab) {
      await loadFileIntoEditor(tab)
    }
  }
})

// 将文件加载到编辑器
async function loadFileIntoEditor(tab: any) {
  const language = tab.language || 'plaintext'
  
  // 设置语言
  monaco.editor.setModelLanguage(editor!.getModel()!, language)
  
  // 设置内容
  if (tab.content !== undefined) {
    editor!.setValue(tab.content)
  }
  
  // 滚动到顶部
  editor!.setScrollPosition({ scrollTop: 0 })
  
  // 等待内容渲染后重新布局
  await nextTick()
  requestAnimationFrame(() => {
    editor?.layout()
  })
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
  if (resizeObserver) {
    resizeObserver.disconnect()
  }
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
  flex-shrink: 0;  /* 防止标签页被压缩 */
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
  border-radius: 0 !important;  /* 移除圆角，更像 VSCode */
  border-top: 2px solid transparent !important;  /* 活动标签顶部高亮 */
}

:deep(.n-tabs .n-tabs-tab:hover) {
  background-color: #2D2D2D !important;
  color: #CCCCCC !important;
}

:deep(.n-tabs .n-tabs-tab--active) {
  background-color: #1E1E1E !important;
  color: #FFFFFF !important;
  border-top: 2px solid #007ACC !important;  /* 活动标签蓝色顶部边框 */
}

:deep(.n-tabs .n-tabs-pane-wrapper) {
  height: 100% !important;
}

:deep(.n-tabs .n-tabs-pane) {
  height: 100% !important;
  padding: 0 !important;
}

/* Monaco 编辑器容器 */
.monaco-container {
  flex: 1;
  overflow: hidden;
  background-color: #1E1E1E;
}

/* Monaco Editor 内部样式调整 - 强制内容靠左 */
.monaco-container :deep(.monaco-editor) {
  /* 减少编辑器左侧的额外边距 */
  .margin {
    width: auto !important;
    min-width: fit-content !important;
  }
  
  /* 调整行号区域 */
  .line-numbers {
    padding-right: 4px !important;  /* 进一步减少行号右侧内边距 */
  }
  
  /* 强制内容区域靠左 */
  .overflow-guard {
    .view-lines {
      left: 0 !important;  /* 强制代码内容从左侧开始 */
    }
    
    .view-overlays {
      .current-line {
        border-width: 0 !important;  /* 移除当前行左边框 */
      }
    }
  }
  
  /* 调整内容边距 */
  .monaco-scrollable-element {
    .scrollbar {
      &.vertical {
        right: 0 !important;
      }
    }
  }
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
