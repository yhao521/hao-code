<template>
  <div class="outline-panel">
    <div class="panel-header">
      <span>大纲</span>
    </div>
    <div class="outline-content" v-if="symbols.length > 0">
      <div
        v-for="symbol in symbols"
        :key="symbol.name + symbol.range.start.line"
        class="outline-item"
        :style="{ paddingLeft: (symbol.level || 0) * 12 + 'px' }"
        @click="jumpToSymbol(symbol)"
      >
        <span class="symbol-icon">{{ getSymbolIcon(symbol.kind) }}</span>
        <span class="symbol-name">{{ symbol.name }}</span>
      </div>
    </div>
    <div class="empty-outline" v-else>
      <p>当前文件没有可显示的符号</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { useEditorStore } from "@/stores/editor";
import { LSPManager } from "@/utils/lspManager";
import * as monaco from "monaco-editor";

const editorStore = useEditorStore();
const lspManager = LSPManager.getInstance();
const symbols = ref<any[]>([]);

async function loadSymbols() {
  if (!editorStore.activeTab) {
    symbols.value = [];
    return;
  }

  const path = editorStore.activeTab.path;
  const langId = getLanguage(path);
  const uri = monaco.Uri.file(path).toString();

  // 简单的语言映射
  const lspLangId =
    langId === "typescript" || langId === "javascript"
      ? langId
      : langId === "go"
        ? "go"
        : "";

  if (lspLangId) {
    const result = await lspManager.getDocumentSymbols(lspLangId, uri);
    symbols.value = flattenSymbols(result || []);
  } else {
    symbols.value = [];
  }
}

function flattenSymbols(items: any[], level = 0): any[] {
  let result: any[] = [];
  for (const item of items) {
    result.push({ ...item, level });
    if (item.children) {
      result = result.concat(flattenSymbols(item.children, level + 1));
    }
  }
  return result;
}

function getSymbolIcon(kind: number): string {
  // LSP SymbolKind 映射到简单图标
  const icons: Record<number, string> = {
    1: "📦", // File
    2: "📦", // Module
    3: "📦", // Namespace
    4: "📦", // Package
    5: "🏗️", // Class
    6: "⚡", // Method
    7: "📋", // Property
    8: "🔧", // Field
    9: "🔢", // Constructor
    10: "🏗️", // Enum
    11: "🔗", // Interface
    12: "🛠️", // Function
    13: "🔄", // Variable
    14: "🔄", // Constant
    15: "📝", // String
    16: "🔢", // Number
    17: "✅", // Boolean
    18: "📊", // Array
  };
  return icons[kind] || "📄";
}

function jumpToSymbol(symbol: any) {
  if (symbol.range) {
    const line = symbol.range.start.line + 1;
    window.dispatchEvent(
      new CustomEvent("editor:jump-to-line", {
        detail: { line },
      }),
    );
  }
}

function getLanguage(path: string): string {
  const ext = path.split(".").pop()?.toLowerCase();
  const map: Record<string, string> = {
    ts: "typescript",
    js: "javascript",
    go: "go",
    py: "python",
  };
  return map[ext || ""] || "plaintext";
}

watch(
  () => editorStore.activeTab?.path,
  () => {
    loadSymbols();
  },
);

// 初始加载
loadSymbols();
</script>

<style scoped>
.outline-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #252526;
  color: #cccccc;
}

.panel-header {
  padding: 8px 12px;
  font-size: 11px;
  font-weight: bold;
  text-transform: uppercase;
  border-bottom: 1px solid #3c3c3c;
}

.outline-content {
  flex: 1;
  overflow-y: auto;
  padding: 4px 0;
}

.outline-item {
  padding: 4px 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}

.outline-item:hover {
  background-color: #2a2d2e;
}

.symbol-icon {
  width: 16px;
  text-align: center;
}

.empty-outline {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #858585;
  font-size: 13px;
}
</style>
