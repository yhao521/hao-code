<template>
  <div class="call-hierarchy-panel">
    <div class="panel-header">
      <span>调用层级</span>
      <n-button text size="tiny" @click="refreshHierarchy" :loading="loading">
        <n-icon><RefreshOutline /></n-icon>
      </n-button>
    </div>

    <div class="hierarchy-content" v-if="rootItem">
      <CallHierarchyNode
        :item="rootItem"
        :level="0"
        @jump-to-location="handleJumpToLocation"
      />
    </div>

    <div class="empty-state" v-else-if="!loading">
      <p>在编辑器中选择一个函数或方法以查看其调用层级。</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { NButton, NIcon } from "naive-ui";
import { RefreshOutline } from "@vicons/ionicons5";
import { useEditorStore } from "@/stores/editor";
import { lspManager } from "@/utils/lspManager";
import CallHierarchyNode from "./CallHierarchyNode.vue";

const editorStore = useEditorStore();
const rootItem = ref<any>(null);
const loading = ref(false);

const refreshHierarchy = async () => {
  if (!editorStore.activeTab || !lspManager) return;
  
  // 注意：这里需要从编辑器实例获取光标位置，目前简化处理，假设用户已选中
  // 在实际集成中，可能需要通过 window.wails 调用后端或监听编辑器事件
  const position = { lineNumber: 1, column: 1 }; // 占位符

  loading.value = true;
  try {
    const langId = getLanguage(editorStore.activeTab.path);
    const uri = `file://${editorStore.activeTab.path}`;
    
    // 1. 准备调用层级入口
    const items = await lspManager.prepareCallHierarchy(
      langId, 
      uri, 
      position.lineNumber - 1, 
      position.column - 1
    );

    if (items && items.length > 0) {
      rootItem.value = items[0];
    } else {
      rootItem.value = null;
    }
  } catch (e) {
    console.error("Failed to load call hierarchy:", e);
  } finally {
    loading.value = false;
  }
};

const handleJumpToLocation = (location: any) => {
  // 触发编辑器跳转逻辑
  window.dispatchEvent(
    new CustomEvent("editor:jump-to-location", { detail: location }),
  );
};

function getLanguage(path: string): string {
  const ext = path.split(".").pop()?.toLowerCase();
  const map: Record<string, string> = {
    go: "go",
    ts: "typescript",
    js: "javascript",
    py: "python",
    java: "java",
  };
  return map[ext || ""] || "plaintext";
}

onMounted(() => {
  // 监听光标变化，自动刷新（可选）
});
</script>

<style scoped>
.call-hierarchy-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #252526;
  color: #cccccc;
}

.panel-header {
  padding: 8px 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #3e3e42;
  font-size: 11px;
  font-weight: bold;
  text-transform: uppercase;
}

.hierarchy-content {
  flex: 1;
  overflow-y: auto;
  padding: 4px 0;
}

.empty-state {
  padding: 20px;
  text-align: center;
  color: #888;
  font-size: 12px;
}
</style>
