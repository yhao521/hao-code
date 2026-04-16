<template>
  <div class="todo-tree-panel">
    <div class="panel-header">
      <span class="title">TODO TREE</span>
      <NIcon class="refresh-icon" @click="loadTodos"><RefreshOutline /></NIcon>
    </div>
    <div class="todo-list" v-if="todos.length > 0">
      <div
        v-for="(item, index) in todos"
        :key="index"
        class="todo-item"
        @click="jumpToLine(item)"
      >
        <div class="todo-tag" :class="getTagClass(item.lineContent)">
          {{ getTag(item.lineContent) }}
        </div>
        <div class="todo-content">
          <div class="file-path">{{ getRelativePath(item.filePath) }}</div>
          <div class="line-text">{{ item.lineContent }}</div>
        </div>
      </div>
    </div>
    <div v-else class="empty-state">
      <p>No TODOs found.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { NIcon } from "naive-ui";
import { RefreshOutline } from "@vicons/ionicons5";
import { GetProjectRoot, ScanTodos } from "@wails/backend/appservice";
import { useEditorStore } from "@/stores/editor";

interface TodoItem {
  filePath: string;
  lineNumber: number;
  lineContent: string;
}

const todos = ref<TodoItem[]>([]);
const editorStore = useEditorStore();

async function loadTodos() {
  try {
    const root = await GetProjectRoot();
    todos.value = await ScanTodos(root);
  } catch (error) {
    console.error("Failed to scan TODOs:", error);
  }
}

function getTag(content: string): string {
  const match = content.match(/(TODO|FIXME|HACK|XXX|BUG):/);
  return match ? match[1] : "NOTE";
}

function getTagClass(content: string): string {
  const tag = getTag(content).toLowerCase();
  return `tag-${tag}`;
}

function getRelativePath(fullPath: string): string {
  // 简化处理，实际应根据项目根目录计算
  const parts = fullPath.split("/");
  return parts.slice(-3).join("/");
}

function jumpToLine(item: TodoItem) {
  // 打开文件（需要先读取文件内容）
  // 这里简化处理，直接传递空内容，实际应该由 EditorArea 监听事件后加载
  editorStore.openFile(item.filePath, "");

  // 触发跳转到指定行的事件
  window.dispatchEvent(
    new CustomEvent("editor:jump-to-line", {
      detail: { path: item.filePath, line: item.lineNumber },
    }),
  );
}

onMounted(() => {
  loadTodos();
});
</script>

<style scoped>
.todo-tree-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #1e1e1e;
}

.panel-header {
  padding: 8px 12px;
  border-bottom: 1px solid #333;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 11px;
  font-weight: bold;
  color: #bbbbbb;
}

.refresh-icon {
  cursor: pointer;
  color: #cccccc;
}

.todo-list {
  flex: 1;
  overflow-y: auto;
}

.todo-item {
  padding: 6px 12px;
  cursor: pointer;
  display: flex;
  gap: 8px;
  border-bottom: 1px solid #2b2b2b;
}

.todo-item:hover {
  background-color: #2a2d2e;
}

.todo-tag {
  font-size: 10px;
  padding: 2px 4px;
  border-radius: 2px;
  font-weight: bold;
  height: fit-content;
}

.tag-todo {
  background: #0e639c;
  color: white;
}
.tag-fixme {
  background: #f48771;
  color: white;
}
.tag-hack {
  background: #cca700;
  color: black;
}
.tag-bug {
  background: #f14c4c;
  color: white;
}

.todo-content {
  flex: 1;
  overflow: hidden;
}

.file-path {
  font-size: 11px;
  color: #858585;
  margin-bottom: 2px;
}

.line-text {
  font-size: 12px;
  color: #cccccc;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
}
</style>
