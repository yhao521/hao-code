<template>
  <div class="editor-area">
    <!-- 面包屑导航 -->
    <Breadcrumb
      v-if="editorStore.activeTab"
      :path="editorStore.activeTab.path"
    />

    <!-- 标签页 -->
    <div class="tabs-container" v-if="editorStore.tabs.length > 0">
      <div
        v-for="tab in editorStore.tabs"
        :key="tab.id"
        class="tab"
        :class="{ active: editorStore.activeEditor === tab.id }"
        @click="handleTabChange(tab.id)"
      >
        <span class="tab-icon">{{ getFileIcon(tab.name) }}</span>
        <span class="tab-name" :class="{ dirty: tab.dirty }">{{
          tab.name
        }}</span>
        <span v-if="tab.dirty" class="dirty-indicator">●</span>
        <span
          class="tab-close"
          @click.stop="editorStore.closeTab(tab.id)"
          title="关闭"
        >
          ×
        </span>
      </div>
    </div>

    <!-- Monaco Editor 容器 -->
    <div
      ref="editorContainer"
      class="monaco-container"
      v-if="!editorStore.isDiffMode"
    ></div>

    <!-- Diff Editor 容器 -->
    <div ref="diffContainer" class="monaco-container diff-mode" v-else>
      <div class="diff-toolbar">
        <span class="diff-title">{{ editorStore.diffInfo?.path }}</span>
        <div class="diff-actions">
          <button @click="toggleDiffMode" class="diff-btn">
            切换到普通视图
          </button>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="editorStore.tabs.length === 0" class="empty-state">
      <div class="empty-content">
        <h2>Hao-Code Editor</h2>
        <p class="subtitle">轻量级跨平台代码编辑器</p>
        <div class="shortcuts">
          <div class="shortcut-item">
            <kbd>Ctrl+P</kbd>
            <span>快速打开文件</span>
          </div>
          <div class="shortcut-item">
            <kbd>Ctrl+Shift+P</kbd>
            <span>显示所有命令</span>
          </div>
          <div class="shortcut-item">
            <kbd>Ctrl+B</kbd>
            <span>切换侧边栏可见性</span>
          </div>
          <div class="shortcut-item">
            <kbd>Ctrl+S</kbd>
            <span>保存文件</span>
          </div>
        </div>
        <div class="start-actions">
          <button @click="triggerOpenFile" class="action-button">
            <span>📄</span> 打开文件
          </button>
          <button @click="triggerOpenFolder" class="action-button">
            <span>📁</span> 打开文件夹
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from "vue";
import { useMessage } from "naive-ui";
import * as monaco from "monaco-editor";
import { useEditorStore } from "@/stores/editor";
import {
  WriteFile,
  OpenFileDialog,
  OpenFolderDialog,
} from "@wails/backend/appservice";
import Breadcrumb from "../Breadcrumb.vue";

const editorStore = useEditorStore();
const message = useMessage();
const editorContainer = ref<HTMLElement | null>(null);
const diffContainer = ref<HTMLElement | null>(null);
let editor: monaco.editor.IStandaloneCodeEditor | null = null;
let diffEditor: monaco.editor.IStandaloneDiffEditor | null = null;
let resizeObserver: ResizeObserver | null = null;

// 监听跳转事件
onMounted(() => {
  window.addEventListener("editor:jump-to-line", handleJumpToLine as any);
});

onUnmounted(() => {
  window.removeEventListener("editor:jump-to-line", handleJumpToLine as any);
});

function handleJumpToLine(event: CustomEvent) {
  const { path, line } = event.detail;
  if (editor && editorStore.activeTab?.path === path) {
    editor.revealLineInCenter(line);
    editor.setPosition({ lineNumber: line, column: 1 });
    editor.focus();
  }
}

// 获取文件图标（简单实现）
function getFileIcon(filename: string): string {
  const ext = filename.split(".").pop()?.toLowerCase();
  const iconMap: Record<string, string> = {
    ts: "🔷",
    js: "📜",
    vue: "💚",
    go: "🔵",
    py: "🐍",
    java: "☕",
    html: "🌐",
    css: "🎨",
    json: "📋",
    md: "📝",
    txt: "📄",
  };
  return iconMap[ext || ""] || "📄";
}

// 触发打开文件
async function triggerOpenFile() {
  try {
    const selectedPath = await OpenFileDialog();
    if (selectedPath) {
      window.dispatchEvent(new CustomEvent("menu:open-file"));
    }
  } catch (error) {
    console.error("Failed to open file:", error);
  }
}

// 触发打开文件夹
async function triggerOpenFolder() {
  try {
    const selectedPath = await OpenFolderDialog();
    if (selectedPath) {
      window.dispatchEvent(new CustomEvent("menu:open-folder"));
    }
  } catch (error) {
    console.error("Failed to open folder:", error);
  }
}

// 初始化编辑器
onMounted(async () => {
  if (editorContainer.value && !editor) {
    // 等待容器完全渲染
    await nextTick();

    editor = monaco.editor.create(editorContainer.value, {
      value: "// Welcome to Hao-Code Editor\n",
      language: "typescript",
      theme: "vs-dark",
      automaticLayout: false, // 禁用自动布局，改为手动控制
      fontSize: 14,
      lineHeight: 22, // 增加行高，提升可读性
      fontFamily:
        "'Fira Code', 'Cascadia Code', 'Source Code Pro', Consolas, 'Courier New', monospace",
      fontLigatures: true, // 启用字体连字
      minimap: {
        enabled: true,
        scale: 1,
        showSlider: "mouseover",
        renderCharacters: false,
      },
      scrollBeyondLastLine: false,
      renderWhitespace: "selection",
      bracketPairColorization: {
        enabled: true,
      },
      guides: {
        bracketPairs: true,
        indentation: true,
        highlightActiveBracketPair: true,
        highlightActiveIndentation: true,
      },
      tabSize: 2,
      insertSpaces: true,
      wordWrap: "on",
      padding: {
        top: 10, // 顶部内边距
        bottom: 10, // 底部内边距
      },
      cursorBlinking: "smooth", // 光标平滑闪烁
      cursorSmoothCaretAnimation: "on", // 平滑光标动画
      smoothScrolling: true, // 平滑滚动
      renderLineHighlight: "all", // 高亮当前行
      lineNumbers: "on",
      lineNumbersMinChars: 2, // 减小行号最小宽度，从3改为2，节省左侧空间
      glyphMargin: false, // 禁用字形边距，减少左侧空白
      folding: true, // 启用代码折叠
      foldingStrategy: "indentation",
      showFoldingControls: "mouseover",
      links: true, // 启用链接检测
      colorDecorators: true, // 颜色装饰器
      formatOnPaste: true, // 粘贴时格式化
      formatOnType: true, // 输入时格式化
      scrollbar: {
        verticalScrollbarSize: 10, // 垂直滚动条宽度
        horizontalScrollbarSize: 10, // 水平滚动条宽度
        verticalSliderSize: 10,
        horizontalSliderSize: 10,
      },
      // 关键配置：禁用多余的边距
      roundedSelection: false,
      suggestOnTriggerCharacters: true,
      acceptSuggestionOnEnter: "on",
      tabCompletion: "on",
      wordBasedSuggestions: "currentDocument",
    });

    // 初始布局
    await nextTick();
    editor.layout();

    // 使用 ResizeObserver 监听容器尺寸变化
    resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        // 使用 requestAnimationFrame 确保在下一帧更新
        requestAnimationFrame(() => {
          if (editor) {
            editor.layout();
          }
        });
      }
    });

    if (editorContainer.value) {
      resizeObserver.observe(editorContainer.value);
    }

    // 监听窗口大小变化
    const handleResize = () => {
      requestAnimationFrame(() => {
        if (editor) {
          editor.layout();
        }
      });
    };
    window.addEventListener("resize", handleResize);

    // 存储清理函数
    (editor as any)._cleanupResize = () => {
      window.removeEventListener("resize", handleResize);
    };

    // 监听内容变化
    editor.onDidChangeModelContent(() => {
      if (editorStore.activeEditor) {
        const content = editor!.getValue();
        editorStore.updateContent(editorStore.activeEditor, content);
      }
    });

    // 注册保存快捷键
    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS, () => {
      if (editorStore.activeEditor) {
        handleSave(editorStore.activeEditor);
      }
    });

    // 监听跳转事件（确保编辑器实例存在后也能响应）
    window.addEventListener("editor:jump-to-line", handleJumpToLine as any);
  }
});

// 监听标签页变化
watch(
  () => editorStore.activeEditor,
  async (newTabId) => {
    if (newTabId && editor && !editorStore.isDiffMode) {
      const tab = editorStore.tabs.find((t) => t.id === newTabId);
      if (tab) {
        await loadFileIntoEditor(tab);
      }
    }
  },
);

// 监听 Diff 模式变化
watch(
  () => editorStore.isDiffMode,
  async (isDiff) => {
    if (isDiff) {
      await initDiffEditor();
    } else {
      disposeDiffEditor();
    }
  },
);

// 将文件加载到编辑器
async function loadFileIntoEditor(tab: any) {
  const language = tab.language || "plaintext";

  // 设置语言
  monaco.editor.setModelLanguage(editor!.getModel()!, language);

  // 设置内容
  if (tab.content !== undefined) {
    editor!.setValue(tab.content);
  }

  // 滚动到顶部
  editor!.setScrollPosition({ scrollTop: 0 });

  // 等待内容渲染后重新布局
  await nextTick();
  requestAnimationFrame(() => {
    editor?.layout();
  });
}

// 处理标签页切换
function handleTabChange(tabId: string) {
  editorStore.activeEditor = tabId;
}

// 保存文件
async function handleSave(tabId: string) {
  const tab = editorStore.tabs.find((t) => t.id === tabId);
  if (!tab || !tab.content) return;

  try {
    await WriteFile(tab.path, tab.content);
    editorStore.saveFile(tabId);
    message.success(`已保存: ${tab.name}`);
  } catch (error) {
    message.error(`保存失败: ${error}`);
  }
}

// 初始化 Diff 编辑器
async function initDiffEditor() {
  if (!diffContainer.value || diffEditor) return;

  await nextTick();

  diffEditor = monaco.editor.createDiffEditor(diffContainer.value, {
    theme: "vs-dark",
    automaticLayout: true,
    renderSideBySide: true,
    readOnly: true,
  });

  const info = editorStore.diffInfo;
  if (info) {
    const originalModel = monaco.editor.createModel(
      info.oldContent,
      getLanguage(info.path),
    );
    const modifiedModel = monaco.editor.createModel(
      info.newContent,
      getLanguage(info.path),
    );
    diffEditor.setModel({
      original: originalModel,
      modified: modifiedModel,
    });
  }
}

function disposeDiffEditor() {
  if (diffEditor) {
    diffEditor.dispose();
    diffEditor = null;
  }
}

function toggleDiffMode() {
  editorStore.setDiffMode(false);
}

function getLanguage(path: string): string {
  const ext = path.split(".").pop()?.toLowerCase();
  const map: Record<string, string> = {
    ts: "typescript",
    js: "javascript",
    vue: "vue",
    go: "go",
    py: "python",
    java: "java",
    html: "html",
    css: "css",
    json: "json",
    md: "markdown",
  };
  return map[ext || ""] || "plaintext";
}

// 清理
onUnmounted(() => {
  window.removeEventListener("editor:jump-to-line", handleJumpToLine as any);
  if (resizeObserver) {
    resizeObserver.disconnect();
  }
  if (editor) {
    editor.dispose();
  }
});
</script>

<style scoped>
.editor-area {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #1e1e1e;
}

/* VSCode 风格标签页 */
.tabs-container {
  display: flex;
  background-color: #2d2d2d;
  border-bottom: 1px solid #252526;
  overflow-x: auto;
  overflow-y: hidden;
  min-height: 35px;
  flex-shrink: 0;
}

.tabs-container::-webkit-scrollbar {
  height: 3px;
}

.tabs-container::-webkit-scrollbar-thumb {
  background-color: #424242;
}

.tab {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background-color: #2d2d2d;
  color: #969696;
  border-right: 1px solid #252526;
  cursor: pointer;
  user-select: none;
  font-size: 13px;
  min-width: 120px;
  max-width: 200px;
  transition: background-color 0.15s;
  position: relative;
}

.tab:hover {
  background-color: #2a2a2a;
}

.tab.active {
  background-color: #1e1e1e;
  color: #ffffff;
  border-top: 1px solid #007acc;
}

.tab-icon {
  font-size: 14px;
  flex-shrink: 0;
}

.tab-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tab-name.dirty {
  font-style: italic;
}

.dirty-indicator {
  color: #e8c27a;
  font-size: 10px;
  flex-shrink: 0;
}

.tab-close {
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 3px;
  font-size: 16px;
  line-height: 1;
  opacity: 0;
  transition: all 0.15s;
  flex-shrink: 0;
}

.tab:hover .tab-close,
.tab.active .tab-close {
  opacity: 1;
}

.tab-close:hover {
  background-color: #5a5a5a;
  color: #ffffff;
}

/* Monaco 编辑器容器 */
.monaco-container {
  flex: 1;
  overflow: hidden;
  background-color: #1e1e1e;
}

.diff-mode {
  display: flex;
  flex-direction: column;
}

.diff-toolbar {
  height: 30px;
  background-color: #2d2d2d;
  border-bottom: 1px solid #252526;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 10px;
  font-size: 12px;
  color: #cccccc;
}

.diff-btn {
  background: transparent;
  border: 1px solid #454545;
  color: #cccccc;
  padding: 2px 8px;
  cursor: pointer;
  border-radius: 2px;
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
    padding-right: 4px !important; /* 进一步减少行号右侧内边距 */
  }

  /* 强制内容区域靠左 */
  .overflow-guard {
    .view-lines {
      left: 0 !important; /* 强制代码内容从左侧开始 */
    }

    .view-overlays {
      .current-line {
        border-width: 0 !important; /* 移除当前行左边框 */
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
  background-color: #1e1e1e;
}

.empty-content {
  text-align: center;
  color: #858585;
  max-width: 500px;
  padding: 40px;
}

.empty-content h2 {
  font-size: 32px;
  font-weight: 300;
  color: #cccccc;
  margin-bottom: 8px;
  letter-spacing: -0.5px;
}

.subtitle {
  font-size: 14px;
  color: #858585;
  margin-bottom: 32px;
}

.shortcuts {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 32px;
  text-align: left;
}

.shortcut-item {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 13px;
  color: #858585;
}

kbd {
  background-color: #3c3c3c;
  padding: 4px 8px;
  border-radius: 3px;
  border: 1px solid #555;
  font-family: "Consolas", "Monaco", monospace;
  font-size: 12px;
  color: #cccccc;
  min-width: 80px;
  text-align: center;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}

.start-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.action-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background-color: #0e639c;
  color: #ffffff;
  border: none;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.action-button:hover {
  background-color: #1177bb;
}

.action-button span {
  font-size: 16px;
}
</style>
