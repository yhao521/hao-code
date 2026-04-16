<template>
  <div class="file-history-panel">
    <div class="panel-header">
      <span class="title">TIMELINE</span>
    </div>
    <div class="history-list" v-if="history.length > 0">
      <div
        v-for="(commit, index) in history"
        :key="commit.hash"
        class="history-item"
      >
        <div
          class="commit-dot"
          :style="{ backgroundColor: getColor(index) }"
        ></div>
        <div class="commit-content">
          <div class="commit-msg">{{ commit.message }}</div>
          <div class="commit-meta">
            <span class="author">{{ commit.author }}</span>
            <span class="time">{{ commit.timestamp }}</span>
          </div>
          <div class="commit-hash">{{ commit.shortHash }}</div>
        </div>
      </div>
    </div>
    <div v-else class="empty-state">
      <p>No history found for this file.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { useEditorStore } from "@/stores/editor";
import { GetProjectRoot, GetFileHistory } from "@wails/backend/appservice";

const editorStore = useEditorStore();
const history = ref<any[]>([]);

const colors = ["#E57373", "#64B5F6", "#81C784", "#FFD54F", "#BA68C8"];

function getColor(index: number) {
  return colors[index % colors.length];
}

async function loadHistory(filePath: string) {
  if (!filePath) {
    history.value = [];
    return;
  }
  try {
    const root = await GetProjectRoot();
    // 计算相对路径
    const relativePath = filePath.replace(root + "/", "");
    history.value = await GetFileHistory(root, relativePath);
  } catch (error) {
    console.error("Failed to load history:", error);
  }
}

// 监听当前活动标签页的变化
watch(
  () => editorStore.activeTab?.path,
  (newPath) => {
    loadHistory(newPath || "");
  },
  { immediate: true },
);
</script>

<style scoped>
.file-history-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #1e1e1e;
}

.panel-header {
  padding: 8px 12px;
  border-bottom: 1px solid #333;
}

.title {
  font-size: 11px;
  font-weight: bold;
  color: #bbbbbb;
}

.history-list {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.history-item {
  display: flex;
  margin-bottom: 16px;
  position: relative;
}

.commit-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  margin-top: 5px;
  margin-right: 12px;
  flex-shrink: 0;
}

.commit-content {
  flex: 1;
}

.commit-msg {
  font-size: 13px;
  color: #cccccc;
  margin-bottom: 4px;
}

.commit-meta {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: #858585;
}

.commit-hash {
  font-family: monospace;
  font-size: 10px;
  color: #5c6370;
  margin-top: 2px;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
}
</style>
