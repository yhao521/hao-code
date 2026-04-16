<template>
  <div class="tasks-panel">
    <div class="panel-header">
      <span>TASKS</span>
      <button @click="refreshTasks" class="refresh-btn" title="刷新任务列表">
        ↻
      </button>
    </div>

    <div v-if="loading" class="loading-state">Loading tasks...</div>
    <div v-else-if="tasks.length === 0" class="empty-state">
      No tasks found. Try adding scripts to package.json or a Taskfile.yml.
    </div>
    <div v-else class="task-list">
      <div
        v-for="(task, index) in tasks"
        :key="index"
        class="task-item"
        @click="runTask(task)"
      >
        <span class="task-icon">▶</span>
        <div class="task-info">
          <span class="task-name">{{ task.name }}</span>
          <span class="task-source">{{ task.source }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useMessage } from "naive-ui";
import { GetTasks, RunTask } from "@wails/backend/appservice";
import { useEditorStore } from "@/stores/editor";

const store = useEditorStore();
const message = useMessage();
const tasks = ref<any[]>([]);
const loading = ref(false);

async function refreshTasks() {
  if (!store.workspace?.path) return;

  loading.value = true;
  try {
    tasks.value = await GetTasks(store.workspace.path);
  } catch (error) {
    message.error("Failed to load tasks");
  } finally {
    loading.value = false;
  }
}

async function runTask(task: any) {
  if (!store.workspace?.path) return;

  message.info(`Running: ${task.name}`);
  // 在实际实现中，这里应该触发终端面板并发送命令
  // 目前我们简单调用后端执行
  try {
    await RunTask(store.workspace.path, task.command);
    message.success(`Task '${task.name}' finished.`);
  } catch (error) {
    message.error(`Task '${task.name}' failed.`);
  }
}

onMounted(() => {
  refreshTasks();
});
</script>

<style scoped>
.tasks-panel {
  padding: 10px;
  height: 100%;
  overflow-y: auto;
  color: #cccccc;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  font-weight: bold;
  font-size: 11px;
  letter-spacing: 1px;
}

.refresh-btn {
  background: none;
  border: none;
  color: #cccccc;
  cursor: pointer;
  font-size: 16px;
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.task-item {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  cursor: pointer;
  border-radius: 3px;
}

.task-item:hover {
  background-color: #2a2d2e;
}

.task-icon {
  color: #89d185;
  margin-right: 8px;
  font-size: 10px;
}

.task-info {
  display: flex;
  flex-direction: column;
}

.task-name {
  font-size: 13px;
}

.task-source {
  font-size: 11px;
  color: #858585;
}

.empty-state,
.loading-state {
  font-size: 13px;
  color: #858585;
  text-align: center;
  margin-top: 20px;
}
</style>
