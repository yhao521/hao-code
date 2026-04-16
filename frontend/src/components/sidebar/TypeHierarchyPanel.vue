<template>
  <div class="type-hierarchy-panel">
    <div class="panel-header">
      <span>类型层次结构</span>
      <div class="actions">
        <NButton text @click="refresh" :loading="loading">
          <template #icon
            ><NIcon><RefreshOutline /></NIcon
          ></template>
        </NButton>
      </div>
    </div>

    <div class="content" v-if="rootItem">
      <TypeHierarchyNode
        :item="rootItem"
        :level="0"
        @navigate="handleNavigate"
      />
    </div>
    <div class="empty-state" v-else-if="!loading">
      在编辑器中选择一个类型以查看其层次结构
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { NButton, NIcon } from "naive-ui";
import { RefreshOutline } from "@vicons/ionicons5";
import { useEditorStore } from "@/stores/editor";
import { lspManager } from "@/utils/lspManager";
import TypeHierarchyNode from "./TypeHierarchyNode.vue";

const editorStore = useEditorStore();
const rootItem = ref<any>(null);
const loading = ref(false);

async function refresh() {
  if (!editorStore.activeTab) return;

  const model = editorStore.activeMonacoModel;
  if (!model) return;

  const position = model.getPositionAt(
    model.getOffsetAt({ lineNumber: 1, column: 1 }),
  );
  // 实际应获取当前光标位置，这里简化处理
  const cursorPosition = editorStore.activeCursor || {
    lineNumber: 1,
    column: 1,
  };

  loading.value = true;
  try {
    const langId = getLanguage(editorStore.activeTab.path);
    const uri = model.uri.toString();

    const items = await lspManager.getTypeHierarchy(
      langId,
      uri,
      cursorPosition.lineNumber - 1,
      cursorPosition.column - 1,
    );

    if (items && items.length > 0) {
      rootItem.value = items[0];
    } else {
      rootItem.value = null;
    }
  } catch (e) {
    console.error("Failed to fetch type hierarchy:", e);
  } finally {
    loading.value = false;
  }
}

function handleNavigate(item: any) {
  const uri = item.uri || item.location?.uri;
  const range = item.range || item.location?.range;
  if (uri && range) {
    window.dispatchEvent(
      new CustomEvent("editor:jump-to-location", {
        detail: {
          uri,
          line: range.start.line + 1,
          col: range.start.character + 1,
        },
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
    java: "java",
  };
  return map[ext || ""] || "";
}
</script>

<style scoped>
.type-hierarchy-panel {
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
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #3c3c3c;
}

.content {
  flex: 1;
  overflow-y: auto;
  padding: 4px 0;
}

.empty-state {
  padding: 20px;
  text-align: center;
  color: #858585;
  font-size: 13px;
}
</style>
