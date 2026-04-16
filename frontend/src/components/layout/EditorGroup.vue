<template>
  <div class="editor-group" @click="activateGroup">
    <div class="tabs-header">
      <div
        v-for="tab in group.tabs"
        :key="tab.id"
        class="tab-item"
        :class="{ active: tab.id === group.activeTabId }"
        @click.stop="selectTab(tab.id)"
      >
        <span class="tab-name">{{ tab.name }}</span>
        <span class="close-icon" @click.stop="handleCloseTab(tab.id)">×</span>
      </div>
      <button
        class="theme-btn"
        @click.stop="handleImportTheme"
        title="导入 VSCode 主题"
      >
        🎨
      </button>
    </div>
    <div class="monaco-container" ref="monacoRef"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from "vue";
import * as monaco from "monaco-editor";
import { useEditorStore, type EditorGroup } from "@/stores/editor";
import {
  GetGhostText,
  GetProjectRoot,
  GetFileBlame,
  GetDiagnostics,
  OpenFileDialog,
  GetThemeDefinition,
} from "@wails/backend/appservice";

const props = defineProps<{
  group: EditorGroup;
}>();

const store = useEditorStore();
const monacoRef = ref<HTMLDivElement>();
let editorInstance: monaco.editor.IStandaloneCodeEditor | null = null;
const models = new Map<string, monaco.editor.ITextModel>();
let ghostTextDisposable: monaco.IDisposable | null = null;

const isActive = computed(() => store.activeGroupId === props.group.id);

function activateGroup() {
  store.setActiveGroup(props.group.id);
}

function selectTab(id: string) {
  props.group.activeTabId = id;
  activateGroup();
}

function handleCloseTab(id: string) {
  // 关闭时销毁对应的 model 以释放内存
  const model = models.get(id);
  if (model) {
    model.dispose();
    models.delete(id);
  }
  store.closeTab(id);
}

function getModelForTab(tab: any): monaco.editor.ITextModel {
  let model = models.get(tab.id);
  if (!model) {
    model = monaco.editor.createModel(
      tab.content || "",
      tab.language || "plaintext",
      monaco.Uri.parse(`file://${tab.path}`),
    );
    models.set(tab.id, model);
  }
  return model;
}

// 注册 AI Ghost Text 提供者
function registerAIGhostText() {
  if (ghostTextDisposable) {
    ghostTextDisposable.dispose();
  }

  ghostTextDisposable = monaco.languages.registerInlineCompletionsProvider(
    "*",
    {
      provideInlineCompletions: async (model, position) => {
        try {
          const text = model.getValue();
          const offset = model.getOffsetAt(position);
          const prefix = text.substring(0, offset);
          const suffix = text.substring(offset);

          // 调用后端 AI 服务
          const result = await GetGhostText(
            prefix,
            suffix,
            model.getLanguageId(),
            model.uri.toString(),
          );

          if (result && result.text) {
            return {
              items: [{ insertText: result.text }],
            };
          }
        } catch (error) {
          console.error("AI completion error:", error);
        }
        return { items: [] };
      },
      disposeInlineCompletions: () => {},
    },
  );
}

// 注册 Error Lens (错误透镜)
async function registerErrorLens(model: monaco.editor.ITextModel) {
  const uri = model.uri.toString();
  const langId = model.getLanguageId();

  try {
    // 获取诊断信息
    const diagnostics = await GetDiagnostics(langId, uri);

    const decorations: monaco.editor.IModelDeltaDecoration[] = [];

    diagnostics.forEach((d: any) => {
      if (d.range && d.message) {
        const range = new monaco.Range(
          d.range.start.line + 1,
          d.range.start.character + 1,
          d.range.end.line + 1,
          d.range.end.character + 1,
        );

        decorations.push({
          range: range,
          options: {
            isWholeLine: false,
            inlineClassName: "error-lens-inline",
            after: {
              content: ` ⚠ ${d.message}`,
              inlineClassName: "error-lens-text",
            },
            glyphMarginClassName:
              d.severity === 1 ? "glyph-error" : "glyph-warning",
            hoverMessage: { value: d.message },
          },
        });
      }
    });

    model.deltaDecorations([], decorations);
  } catch (e) {
    console.error("Failed to load diagnostics", e);
  }
}

// 主题导入与应用
async function handleImportTheme() {
  try {
    const path = await OpenFileDialog();
    if (!path) return;

    // 1. 获取 Monaco 兼容的定义
    const themeDef = await GetThemeDefinition(path);

    // 2. 定义主题
    const themeName = "custom-imported-theme";
    monaco.editor.defineTheme(themeName, themeDef as any);

    // 3. 应用主题
    monaco.editor.setTheme(themeName);
    console.log("Theme applied successfully");
  } catch (e) {
    console.error("Failed to import theme", e);
  }
}

// 注册 Blame Hover Provider
async function registerBlameHover(model: monaco.editor.ITextModel) {
  const path = model.uri.fsPath;
  if (!path) return;

  try {
    const root = await GetProjectRoot();
    const relativePath = path.replace(root + "/", "");
    const blames = await GetFileBlame(root, relativePath);

    // 创建一个按行号映射的 Blame 信息 Map
    const blameMap = new Map<number, any>();
    // 生成一个颜色 Map，为不同的提交者分配不同颜色
    const authorColors = new Map<string, string>();
    let hue = 0;

    // 为每行创建 Blame 信息并分配颜色
    blames.forEach((b: any) => {
      blameMap.set(b.line, b);
      if (!authorColors.has(b.author)) {
        // 使用 HSL 生成不同色调的颜色
        authorColors.set(b.author, `hsl(${hue}, 70%, 60%)`);
        hue = (hue + 137.5) % 360; // 使用黄金角度确保颜色分散
      }
    });

    // 创建 Glyph Margin 装饰器
    const glyphDecorations: monaco.editor.IModelDeltaDecoration[] = [];
    blames.forEach((b: any) => {
      glyphDecorations.push({
        range: new monaco.Range(b.line + 1, 1, b.line + 1, 1),
        options: {
          glyphMarginClassName: `blame-glyph-margin ${b.author.replace(/[^a-zA-Z0-9]/g, "-")}`,
          glyphMarginHoverMessage: {
            value: `**${b.author}** (${b.timestamp})\n${b.message}`,
          },
          isWholeLine: true,
          linesDecorationsClassName: "blame-line-decoration",
        },
      });
    });

    // 应用 Glyph Margin 装饰器
    model.deltaDecorations([], glyphDecorations);

    // 注册 Glyph Margin 点击事件 (通过监听容器点击并判断位置实现)
    const domNode = editorInstance!.getDomNode();
    if (domNode) {
      domNode.addEventListener("click", (e: MouseEvent) => {
        const target = e.target as HTMLElement;
        if (
          target.classList.contains("blame-glyph-margin") ||
          target.parentElement?.classList.contains("blame-glyph-margin")
        ) {
          // 获取当前鼠标位置对应的行号
          const position = editorInstance!.getPosition();
          if (position) {
            const blame = blameMap.get(position.lineNumber - 1);
            if (blame) {
              window.dispatchEvent(
                new CustomEvent("editor:show-blame", { detail: blame }),
              );
            }
          }
        }
      });
    }

    // 注册 Hover Provider
    monaco.languages.registerHoverProvider("*", {
      provideHover: (model, position) => {
        const blame = blameMap.get(position.lineNumber - 1); // Blame 行号从0开始
        if (blame) {
          return {
            contents: [
              { value: `**${blame.author}** (${blame.timestamp})` },
              { value: blame.message },
            ],
          };
        }
        return null;
      },
    });
  } catch (e) {
    console.error("Failed to load blame info", e);
  }
}

onMounted(() => {
  if (monacoRef.value) {
    editorInstance = monaco.editor.create(monacoRef.value, {
      theme: "vs-dark",
      automaticLayout: true,
      minimap: { enabled: false },
      scrollBeyondLastLine: false,
      inlineSuggest: { enabled: true },
    });

    // 注册 AI 补全
    registerAIGhostText();

    // 监听保存事件进行自动格式化
    editorInstance.onDidChangeModelContent(() => {
      // 这里可以触发防抖后的 LSP 诊断请求
    });
  }
});

watch(
  () => props.group.activeTabId,
  (newId) => {
    const tab = props.group.tabs.find((t) => t.id === newId);
    if (tab && editorInstance) {
      const model = getModelForTab(tab);
      editorInstance.setModel(model);
      // 为当前模型加载 Blame 信息和 Error Lens
      registerBlameHover(model);
      registerErrorLens(model);
    }
  },
  { immediate: true },
);

watch(
  () => isActive.value,
  (active) => {
    if (active && editorInstance) {
      editorInstance.layout();
      editorInstance.focus();
    }
  },
);
</script>

<style scoped>
.editor-group {
  display: flex;
  flex-direction: column;
  height: 100%;
  border-right: 1px solid #2b2b2b;
}

.tabs-header {
  display: flex;
  background-color: #252526;
  overflow-x: auto;
  align-items: center;
}

.theme-btn {
  background: transparent;
  border: none;
  color: #858585;
  cursor: pointer;
  padding: 4px 8px;
  font-size: 14px;
  margin-left: auto;
}

.theme-btn:hover {
  color: #ffffff;
  background-color: #2a2d2e;
}

.tab-item {
  padding: 8px 12px;
  font-size: 13px;
  cursor: pointer;
  background-color: #2d2d2d;
  border-right: 1px solid #252526;
  display: flex;
  align-items: center;
  min-width: 100px;
}

.tab-item.active {
  background-color: #1e1e1e;
  color: #ffffff;
  border-top: 1px solid #007acc;
}

.tab-name {
  margin-right: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.close-icon {
  opacity: 0;
  transition: opacity 0.2s;
}

.tab-item:hover .close-icon {
  opacity: 1;
}

.monaco-container {
  flex: 1;
  width: 100%;
}

/* Error Lens Styles */
.error-lens-text {
  color: #888;
  font-style: italic;
  margin-left: 8px;
  opacity: 0.7;
}

.glyph-error {
  background: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16'><circle cx='8' cy='8' r='7' fill='%23f48771'/></svg>")
    no-repeat center center;
}

.glyph-warning {
  background: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16'><circle cx='8' cy='8' r='7' fill='%23cca700'/></svg>")
    no-repeat center center;
}
</style>
