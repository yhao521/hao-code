<template>
  <div
    v-if="visible"
    ref="dropdownRef"
    class="recent-files-dropdown"
    :style="dropdownStyle"
  >
    <div class="dropdown-header">
      <span class="dropdown-title">最近的文件</span>
      <NButton text size="tiny" @click="handleClearAll">清空</NButton>
    </div>

    <div class="dropdown-content">
      <!-- 最近文件列表 -->
      <div v-if="recentFiles.length > 0" class="section">
        <div class="section-title">最近文件</div>
        <div
          v-for="item in recentFiles.slice(0, 10)"
          :key="'file-' + item.path"
          class="dropdown-item"
          @click="handleOpenFile(item.path)"
        >
          <div class="item-icon">📄</div>
          <div class="item-info">
            <div class="item-name">{{ item.name }}</div>
            <div class="item-path">{{ item.path }}</div>
          </div>
          <div class="item-time">{{ formatTime(item.openedAt) }}</div>
        </div>
      </div>

      <!-- 最近文件夹列表 -->
      <div v-if="recentFolders.length > 0" class="section">
        <div class="section-title">最近文件夹</div>
        <div
          v-for="item in recentFolders.slice(0, 5)"
          :key="'folder-' + item.path"
          class="dropdown-item"
          @click="handleOpenFolder(item.path)"
        >
          <div class="item-icon">📁</div>
          <div class="item-info">
            <div class="item-name">{{ item.name }}</div>
            <div class="item-path">{{ item.path }}</div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <NEmpty
        v-if="recentFiles.length === 0 && recentFolders.length === 0"
        description="暂无最近项目"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { NButton, NEmpty, useMessage } from "naive-ui";
import {
  GetRecentFiles,
  GetRecentFolders,
  ClearRecentFiles,
  ClearRecentFolders,
  ReadFile,
  ListDir,
} from "@wails/backend/appservice";
import { useEditorStore } from "@/stores/editor";

interface Props {
  visible: boolean;
  position?: { x: number; y: number };
}

const props = withDefaults(defineProps<Props>(), {
  position: () => ({ x: 0, y: 0 }),
});

const emit = defineEmits<{
  (e: "update:visible", value: boolean): void;
}>();

const editorStore = useEditorStore();
const message = useMessage();
const dropdownRef = ref<HTMLElement | null>(null);

const recentFiles = ref<any[]>([]);
const recentFolders = ref<any[]>([]);

const dropdownStyle = computed<Record<string, string | number>>(() => ({
  position: "fixed",
  left: props.position.x + "px",
  top: props.position.y + "px",
  zIndex: 9999,
}));

// 加载最近文件和文件夹
async function loadData() {
  try {
    recentFiles.value = await GetRecentFiles();
    recentFolders.value = await GetRecentFolders();
  } catch (error) {
    console.error("Failed to load recent files:", error);
  }
}

// 打开文件
async function handleOpenFile(path: string) {
  try {
    message.loading("正在读取文件...", { duration: 0 });
    const content = await ReadFile(path);
    message.destroyAll();

    editorStore.openFile(path, content);
    emit("update:visible", false);
    message.success(`已打开: ${path.split("/").pop()}`);
  } catch (error) {
    message.destroyAll();
    const errorMsg = error instanceof Error ? error.message : String(error);
    message.error(`打开文件失败: ${errorMsg}`);
  }
}

// 打开文件夹
async function handleOpenFolder(path: string) {
  try {
    message.loading("正在加载文件夹...", { duration: 0 });

    // 验证文件夹
    try {
      await ListDir(path);
    } catch (error) {
      message.destroyAll();
      message.error("无法访问该文件夹");
      return;
    }

    editorStore.setWorkspace(path);
    message.destroyAll();
    emit("update:visible", false);
    message.success(`已打开: ${path.split("/").pop()}`);
  } catch (error) {
    message.destroyAll();
    console.error("Failed to open folder:", error);
    const errorMsg = error instanceof Error ? error.message : String(error);
    message.error(`打开文件夹失败: ${errorMsg}`);
  }
}

// 清空所有最近记录
function handleClearAll() {
  ClearRecentFiles();
  ClearRecentFolders();
  recentFiles.value = [];
  recentFolders.value = [];
  message.success("已清空最近记录");
}

// 格式化时间
function formatTime(timeStr: string): string {
  if (!timeStr) return "";
  const date = new Date(timeStr);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const minutes = Math.floor(diff / 60000);
  const hours = Math.floor(diff / 3600000);
  const days = Math.floor(diff / 86400000);

  if (minutes < 1) return "刚刚";
  if (minutes < 60) return `${minutes} 分钟前`;
  if (hours < 24) return `${hours} 小时前`;
  if (days < 7) return `${days} 天前`;
  return date.toLocaleDateString();
}

// 点击外部关闭
function handleClickOutside(event: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    emit("update:visible", false);
  }
}

// 暴露方法供父组件调用
defineExpose({
  loadData,
});

onMounted(() => {
  document.addEventListener("click", handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener("click", handleClickOutside);
});
</script>

<style scoped>
.recent-files-dropdown {
  width: 400px;
  max-height: 500px;
  background: #252526;
  border: 1px solid #3e3e42;
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.5);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.dropdown-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background: #2d2d2d;
  border-bottom: 1px solid #3e3e42;
}

.dropdown-title {
  font-size: 13px;
  font-weight: 600;
  color: #cccccc;
}

.dropdown-content {
  overflow-y: auto;
  flex: 1;
}

.section {
  padding: 8px 0;
}

.section-title {
  padding: 4px 12px;
  font-size: 11px;
  font-weight: 600;
  color: #858585;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.dropdown-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  cursor: pointer;
  transition: background-color 0.2s;
  gap: 8px;
}

.dropdown-item:hover {
  background-color: #2a2d2e;
}

.item-icon {
  font-size: 16px;
  width: 20px;
  text-align: center;
  flex-shrink: 0;
}

.item-info {
  flex: 1;
  min-width: 0;
}

.item-name {
  font-size: 13px;
  color: #cccccc;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-path {
  font-size: 11px;
  color: #858585;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-top: 2px;
}

.item-time {
  font-size: 11px;
  color: #858585;
  flex-shrink: 0;
}

.section + .section {
  border-top: 1px solid #3e3e42;
}
</style>
