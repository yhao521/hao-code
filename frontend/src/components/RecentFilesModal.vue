<template>
  <NModal
    v-model:show="showModal"
    preset="card"
    title="打开最近的项目"
    style="width: 600px; max-width: 90vw"
    :bordered="false"
  >
    <div class="recent-files-modal">
      <!-- 最近文件 -->
      <div class="section" v-if="recentFiles.length > 0">
        <div class="section-header">
          <span class="section-title">最近文件</span>
          <NButton text size="tiny" @click="handleClearFiles">
            清除列表
          </NButton>
        </div>
        <div class="item-list">
          <div
            v-for="item in recentFiles"
            :key="item.path"
            class="item"
            @click="handleOpenFile(item.path)"
          >
            <div class="item-icon">📄</div>
            <div class="item-content">
              <div class="item-name">{{ item.name }}</div>
              <div class="item-path">{{ item.path }}</div>
            </div>
            <div class="item-time">{{ formatTime(item.openedAt) }}</div>
          </div>
        </div>
      </div>

      <!-- 最近文件夹 -->
      <div class="section" v-if="recentFolders.length > 0">
        <div class="section-header">
          <span class="section-title">最近文件夹</span>
          <NButton text size="tiny" @click="handleClearFolders">
            清除列表
          </NButton>
        </div>
        <div class="item-list">
          <div
            v-for="item in recentFolders"
            :key="item.path"
            class="item"
            @click="handleOpenFolder(item.path)"
          >
            <div class="item-icon">📁</div>
            <div class="item-content">
              <div class="item-name">{{ item.name }}</div>
              <div class="item-path">{{ item.path }}</div>
            </div>
            <div class="item-time">{{ formatTime(item.openedAt) }}</div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <NEmpty
        v-if="recentFiles.length === 0 && recentFolders.length === 0"
        description="暂无最近打开的项目"
        style="padding: 40px 0"
      />
    </div>

    <template #footer>
      <NSpace justify="end">
        <NButton @click="showModal = false">关闭</NButton>
      </NSpace>
    </template>
  </NModal>
</template>

<script setup lang="ts">
import { ref } from "vue";
import {
  NModal,
  NButton,
  NSpace,
  NEmpty,
  useMessage,
  useDialog,
} from "naive-ui";
import {
  GetRecentFiles,
  GetRecentFolders,
  ClearRecentFiles,
  ClearRecentFolders,
  ReadFile,
  ListDir,
} from "@wails/go/backend/App";
import { useEditorStore } from "@/stores/editor";

const editorStore = useEditorStore();
const message = useMessage();
const dialog = useDialog();

const showModal = ref(false);
const recentFiles = ref<any[]>([]);
const recentFolders = ref<any[]>([]);

// 显示模态框并加载数据
async function show() {
  showModal.value = true;
  await loadData();
}

// 加载最近文件和文件夹
async function loadData() {
  try {
    recentFiles.value = await GetRecentFiles();
    recentFolders.value = await GetRecentFolders();
  } catch (error) {
    console.error("Failed to load recent files:", error);
    message.error("加载最近文件失败");
  }
}

// 打开文件
async function handleOpenFile(path: string) {
  try {
    message.loading("正在读取文件...", { duration: 0 });
    const content = await ReadFile(path);
    message.destroyAll();

    editorStore.openFile(path, content);
    showModal.value = false;
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
    showModal.value = false;
    message.success(`已打开: ${path.split("/").pop()}`);
  } catch (error) {
    message.destroyAll();
    console.error("Failed to open folder:", error);
    const errorMsg = error instanceof Error ? error.message : String(error);
    message.error(`打开文件夹失败: ${errorMsg}`);
  }
}

// 清除最近文件
function handleClearFiles() {
  dialog.warning({
    title: "确认清除",
    content: "确定要清除最近文件列表吗？",
    positiveText: "确定",
    negativeText: "取消",
    onPositiveClick: async () => {
      try {
        await ClearRecentFiles();
        recentFiles.value = [];
        message.success("已清除最近文件列表");
      } catch (error) {
        message.error("清除失败");
      }
    },
  });
}

// 清除最近文件夹
function handleClearFolders() {
  dialog.warning({
    title: "确认清除",
    content: "确定要清除最近文件夹列表吗？",
    positiveText: "确定",
    negativeText: "取消",
    onPositiveClick: async () => {
      try {
        await ClearRecentFolders();
        recentFolders.value = [];
        message.success("已清除最近文件夹列表");
      } catch (error) {
        message.error("清除失败");
      }
    },
  });
}

// 格式化时间
function formatTime(timeStr: string): string {
  try {
    const date = new Date(timeStr);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    
    // 小于1小时
    if (diff < 3600000) {
      const minutes = Math.floor(diff / 60000);
      return `${minutes}分钟前`;
    }
    
    // 小于24小时
    if (diff < 86400000) {
      const hours = Math.floor(diff / 3600000);
      return `${hours}小时前`;
    }
    
    // 小于7天
    if (diff < 604800000) {
      const days = Math.floor(diff / 86400000);
      return `${days}天前`;
    }
    
    // 显示具体日期
    return date.toLocaleDateString("zh-CN");
  } catch {
    return timeStr;
  }
}

// 暴露方法供外部调用
defineExpose({
  show,
});
</script>

<style scoped>
.recent-files-modal {
  max-height: 60vh;
  overflow-y: auto;
}

.section {
  margin-bottom: 24px;
}

.section:last-child {
  margin-bottom: 0;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #3E3E42;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #BBBBBB;
}

.item-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.item:hover {
  background-color: #2A2D2E;
}

.item-icon {
  font-size: 20px;
  flex-shrink: 0;
}

.item-content {
  flex: 1;
  min-width: 0;
}

.item-name {
  font-size: 13px;
  font-weight: 500;
  color: #CCCCCC;
  margin-bottom: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-path {
  font-size: 11px;
  color: #858585;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item-time {
  font-size: 11px;
  color: #858585;
  flex-shrink: 0;
}
</style>
