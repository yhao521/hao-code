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
        draggable="true"
        :class="{ active: editorStore.activeEditor === tab.id }"
        @click="handleTabChange(tab.id)"
        @dragstart="handleDragStart($event, tab.id)"
        @dragover.prevent
        @drop="handleDrop($event, tab.id)"
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
  GetProjectRoot,
} from "@wails/backend/appservice.js";
import Breadcrumb from "../Breadcrumb.vue";
import { LSPManager } from "@/utils/lspManager";
import { DiagnosticsManager } from "@/utils/diagnosticsManager";

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

    // 确保容器有正确的尺寸
    await new Promise<void>((resolve) => {
      requestAnimationFrame(() => {
        requestAnimationFrame(() => {
          resolve();
        });
      });
    });

    // 调试信息：检查容器尺寸
    console.log("Editor container size before create:", {
      width: editorContainer.value.offsetWidth,
      height: editorContainer.value.offsetHeight,
      clientWidth: editorContainer.value.clientWidth,
      clientHeight: editorContainer.value.clientHeight,
      computedStyle: window.getComputedStyle(editorContainer.value),
    });

    // 如果容器尺寸为 0，等待更长时间
    if (
      editorContainer.value.offsetWidth === 0 ||
      editorContainer.value.offsetHeight === 0
    ) {
      console.warn("Editor container has zero size, waiting longer...");
      await new Promise((resolve) => setTimeout(resolve, 500));
      console.log("Editor container size after waiting:", {
        width: editorContainer.value.offsetWidth,
        height: editorContainer.value.offsetHeight,
      });
    }

    editor = monaco.editor.create(editorContainer.value, {
      value: "// Welcome to Hao-Code Editor\n",
      language: "typescript",
      theme: "vs-dark",
      automaticLayout: true, // 启用自动布局，让编辑器自动适应容器大小
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

    // 多次调用 layout 确保尺寸正确
    editor.layout();
    setTimeout(() => {
      editor?.layout();
    }, 100);
    setTimeout(() => {
      editor?.layout();
    }, 300);

    // 初始化 LSP
    await initLSP();

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

    // 监听断点点击事件
    editor.onMouseDown((e) => {
      if (e.target.type === monaco.editor.MouseTargetType.GUTTER_GLYPH_MARGIN) {
        const lineNumber = e.target.position?.lineNumber;
        if (lineNumber && editorStore.activeTab) {
          editorStore.toggleBreakpoint(editorStore.activeTab.path, lineNumber);
          updateBreakpointDecorations();
        }
      }
    });
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

// 拖拽排序逻辑
let draggedTabId: string | null = null;

function handleDragStart(e: DragEvent, tabId: string) {
  draggedTabId = tabId;
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = "move";
  }
}

function handleDrop(e: DragEvent, targetTabId: string) {
  e.preventDefault();
  if (!draggedTabId || draggedTabId === targetTabId) return;

  const fromIndex = editorStore.tabs.findIndex((t) => t.id === draggedTabId);
  const toIndex = editorStore.tabs.findIndex((t) => t.id === targetTabId);

  if (fromIndex !== -1 && toIndex !== -1) {
    const newTabs = [...editorStore.tabs];
    const [movedTab] = newTabs.splice(fromIndex, 1);
    newTabs.splice(toIndex, 0, movedTab);
    editorStore.tabs = newTabs;
  }
  draggedTabId = null;
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

// 更新断点装饰器
let breakpointDecorations: string[] = [];
function updateBreakpointDecorations() {
  if (!editor || !editorStore.activeTab) return;

  const path = editorStore.activeTab.path;
  const lines = editorStore.breakpoints.get(path) || new Set();

  const newDecorations = Array.from(lines).map((line) => ({
    range: new monaco.Range(line, 1, line, 1),
    options: {
      isWholeLine: true,
      className: "breakpoint-line",
      glyphMarginClassName: "breakpoint-glyph",
    },
  }));

  // 添加当前调试行高亮
  if (editorStore.currentDebugLine?.path === path) {
    newDecorations.push({
      range: new monaco.Range(
        editorStore.currentDebugLine.line,
        1,
        editorStore.currentDebugLine.line,
        1,
      ),
      options: {
        isWholeLine: true,
        className: "debug-current-line",
        glyphMarginClassName: "debug-current-glyph",
      },
    });
  }

  breakpointDecorations = editor.deltaDecorations(
    breakpointDecorations,
    newDecorations,
  );
}

// 监听断点和调试行变化
watch(
  () => [editorStore.breakpoints, editorStore.currentDebugLine],
  () => updateBreakpointDecorations(),
  { deep: true },
);

// LSP 自动补全集成
let lspManager: LSPManager | null = null;
let diagnosticsManager: DiagnosticsManager | null = null;

async function initLSP() {
  if (!lspManager) {
    lspManager = LSPManager.getInstance();
    diagnosticsManager = new DiagnosticsManager();
  }

  if (editorStore.activeTab) {
    const rootPath = await GetProjectRoot();
    const langId = getLanguage(editorStore.activeTab.path);
    // 映射 Monaco 语言 ID 到 LSP 语言 ID
    const lspLangId =
      langId === "typescript" || langId === "javascript"
        ? langId
        : langId === "go"
          ? "go"
          : "";
    if (lspLangId) {
      await lspManager.ensureInitialized(lspLangId, rootPath);
    }
  }
}

// 注册自定义补全提供者
monaco.languages.registerCompletionItemProvider("*", {
  provideCompletionItems: async (model, position) => {
    if (!editorStore.activeTab || !lspManager) return { suggestions: [] };

    const langId = getLanguage(editorStore.activeTab.path);
    const uri = model.uri.toString();
    const line = position.lineNumber - 1;
    const col = position.column - 1;

    try {
      const items = await lspManager.getCompletions(langId, uri, line, col);

      return {
        suggestions: items.map((item: any) => ({
          label: item.label || item.insertText || "Unknown",
          kind: monaco.languages.CompletionItemKind.Text,
          insertText: item.insertText || item.label,
          detail: item.detail,
          range: new monaco.Range(
            position.lineNumber,
            position.column,
            position.lineNumber,
            position.column,
          ),
        })),
      };
    } catch (e) {
      return { suggestions: [] };
    }
  },
});

// 注册定义跳转提供者
monaco.languages.registerDefinitionProvider("*", {
  provideDefinition: async (model, position) => {
    if (!editorStore.activeTab || !lspManager) return null;

    const langId = getLanguage(editorStore.activeTab.path);
    const uri = model.uri.toString();
    const line = position.lineNumber - 1;
    const col = position.column - 1;

    try {
      const result = await lspManager.getDefinition(langId, uri, line, col);
      if (result && result.uri) {
        // 解析 LSP 返回的位置信息
        const targetUri = result.uri;
        const range = result.range || result.selectionRange;
        if (range) {
          return {
            uri: monaco.Uri.parse(targetUri),
            range: new monaco.Range(
              range.start.line + 1,
              range.start.character + 1,
              range.end.line + 1,
              range.end.character + 1,
            ),
          };
        }
      }
    } catch (e) {
      console.error("Definition provider error:", e);
    }
    return null;
  },
});

// 注册引用查找提供者
monaco.languages.registerReferenceProvider("*", {
  provideReferences: async (model, position) => {
    if (!editorStore.activeTab || !lspManager) return null;

    const langId = getLanguage(editorStore.activeTab.path);
    const uri = model.uri.toString();
    const line = position.lineNumber - 1;
    const col = position.column - 1;

    try {
      const refs = await lspManager.findReferences(langId, uri, line, col);
      if (refs && refs.length > 0) {
        return refs.map((ref: any) => {
          const range = ref.range;
          return {
            uri: monaco.Uri.parse(ref.uri),
            range: new monaco.Range(
              range.start.line + 1,
              range.start.character + 1,
              range.end.line + 1,
              range.end.character + 1,
            ),
          };
        });
      }
    } catch (e) {
      console.error("References provider error:", e);
    }
    return null;
  },
});

// 注册重命名提供者
monaco.languages.registerRenameProvider("*", {
  provideRenameEdits: async (model, position, newName) => {
    if (!editorStore.activeTab || !lspManager) return null;

    const langId = getLanguage(editorStore.activeTab.path);
    const uri = model.uri.toString();
    const line = position.lineNumber - 1;
    const col = position.column - 1;

    try {
      const result = await lspManager.renameSymbol(
        langId,
        uri,
        line,
        col,
        newName,
      );
      if (result && result.changes) {
        const edit: monaco.languages.WorkspaceEdit = {
          edits: [],
        };

        for (const [uriStr, textEdits] of Object.entries(result.changes)) {
          const edits = (textEdits as any[]).map((te) => ({
            resource: monaco.Uri.parse(uriStr),
            versionId: undefined,
            textEdit: {
              range: new monaco.Range(
                te.range.start.line + 1,
                te.range.start.character + 1,
                te.range.end.line + 1,
                te.range.end.character + 1,
              ),
              text: te.newText,
            },
          }));
          edit.edits.push(...edits);
        }
        return edit;
      }
    } catch (e) {
      console.error("Rename provider error:", e);
    }
    return null;
  },
});

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
  width: 100%;
  min-width: 0;
  min-height: 0;
  background-color: #1e1e1e;
  overflow: hidden;
  position: relative;
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
  position: relative;
  width: 100%;
  min-height: 0;
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

/* Monaco Editor 内部样式调整 - 让 Monaco 自己管理布局 */
.monaco-container :deep(.monaco-editor) {
  /* Monaco Editor 会自动计算布局，不需要手动干预 */
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

/* 断点样式 */
.breakpoint-glyph {
  background-color: #e51400;
  border-radius: 50%;
  width: 8px !important;
  height: 8px !important;
  margin: 4px !important;
}

.breakpoint-line {
  background-color: rgba(229, 20, 0, 0.1);
}

.debug-current-glyph {
  background-color: #ffcc00;
  border-radius: 50%;
  width: 8px !important;
  height: 8px !important;
  margin: 4px !important;
}

.debug-current-line {
  background-color: rgba(255, 204, 0, 0.2);
}
</style>
