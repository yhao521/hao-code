<template>
  <div class="titlebar" :class="{ macos: isMacOS }">
    <!-- 工作区名称 -->
    <div class="titlebar-center">
      <span class="workspace-name">{{ workspaceName }}</span>
    </div>

    <!-- 功能按钮 -->
    <div class="titlebar-right">
      <NButton
        text
        circle
        size="tiny"
        @click="handleOpenFolder"
        title="打开文件夹"
      >
        <template #icon>
          <NIcon><FolderOpenOutline /></NIcon>
        </template>
      </NButton>

      <!-- macOS: 添加最近文件下拉按钮 -->
      <NButton
        v-if="isMacOS"
        text
        circle
        size="tiny"
        @click="handleShowRecentFiles"
        title="打开最近的文件"
        ref="recentFilesBtnRef"
      >
        <template #icon>
          <NIcon><TimeOutline /></NIcon>
        </template>
      </NButton>

      <!-- Windows: 自定义窗口控制按钮 -->
      <div v-if="!isMacOS" class="window-controls">
        <div
          class="control-btn minimize"
          @click="minimizeWindow"
          title="最小化"
        >
          <NIcon size="14"><RemoveOutline /></NIcon>
        </div>
        <div
          class="control-btn maximize"
          @click="maximizeWindow"
          title="最大化"
        >
          <NIcon size="14"><SquareOutline /></NIcon>
        </div>
        <div class="control-btn close" @click="closeWindow" title="关闭">
          <NIcon size="12"><CloseOutline /></NIcon>
        </div>
      </div>
    </div>

    <!-- 最近文件下拉面板 -->
    <RecentFilesDropdown
      v-model:visible="showRecentFiles"
      :position="dropdownPosition"
      ref="recentFilesDropdownRef"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, nextTick } from "vue";
import { NButton, NIcon, useMessage } from "naive-ui";
import {
  FolderOpenOutline,
  RemoveOutline,
  SquareOutline,
  CloseOutline,
  TimeOutline,
} from "@vicons/ionicons5";
import { useEditorStore } from "@/stores/editor";
import * as wailsRuntime from "@wails/runtime/runtime";
import { ListDir, OpenFolderDialog } from "@wails/go/backend/App";
import RecentFilesDropdown from "./RecentFilesDropdown.vue";

const editorStore = useEditorStore();
const message = useMessage();
const isMaximized = ref(false);

// 最近文件相关
const showRecentFiles = ref(false);
const dropdownPosition = ref({ x: 0, y: 0 });
const recentFilesBtnRef = ref<HTMLElement | null>(null);
const recentFilesDropdownRef = ref<InstanceType<
  typeof RecentFilesDropdown
> | null>(null);

// 检测是否为 macOS
const isMacOS = computed(() => {
  const platform = navigator.platform.toLowerCase();
  return platform.includes("mac");
});

const workspaceName = computed(
  () => editorStore.workspace?.name || "Hao-Code Editor",
);

onMounted(() => {
  // 可以在这里添加窗口状态监听

  // 监听菜单事件：打开最近文件
  wailsRuntime.EventsOn("menu:open-recent", () => {
    handleShowRecentFiles();
  });
});

// 显示最近文件下拉菜单
async function handleShowRecentFiles() {
  if (!recentFilesBtnRef.value) return;

  // 获取按钮位置
  const rect = recentFilesBtnRef.value.getBoundingClientRect();
  dropdownPosition.value = {
    x: rect.left,
    y: rect.bottom + 4,
  };

  showRecentFiles.value = true;

  // 等待 DOM 更新后加载数据
  await nextTick();
  recentFilesDropdownRef.value?.loadData();
}

// 打开文件夹
async function handleOpenFolder() {
  try {
    message.loading("正在打开文件夹选择对话框...", { duration: 0 });

    // 调用后端打开文件夹对话框
    const selectedPath = await OpenFolderDialog();

    message.destroyAll();

    if (!selectedPath) {
      message.info("已取消选择");
      return;
    }

    message.loading("正在加载文件夹...", { duration: 0 });

    // 验证文件夹
    try {
      await ListDir(selectedPath);
    } catch (error) {
      message.destroyAll();
      message.error("无法访问该文件夹");
      return;
    }

    // 设置工作区
    editorStore.setWorkspace(selectedPath);

    message.destroyAll();
    message.success(`已打开: ${selectedPath.split("/").pop()}`);
  } catch (error) {
    message.destroyAll();
    console.error("Failed to open folder:", error);
    const errorMsg = error instanceof Error ? error.message : String(error);
    if (!errorMsg.includes("cancelled")) {
      message.error(`打开文件夹失败: ${errorMsg}`);
    } else {
      message.info("已取消选择");
    }
  }
}

// 窗口控制函数
function minimizeWindow() {
  wailsRuntime.WindowMinimise();
}

function maximizeWindow() {
  wailsRuntime.WindowToggleMaximise();
  isMaximized.value = !isMaximized.value;
}

function closeWindow() {
  wailsRuntime.Quit();
}
</script>

<style scoped>
.titlebar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 38px;
  background-color: #323233;
  color: #cccccc;
  user-select: none;
  border-bottom: 1px solid #252526;
  padding: 0 12px;
}

/* Windows 无边框窗口，标题栏可拖拽 */
.titlebar:not(.macos) {
  -webkit-app-region: drag;
}

.titlebar-center {
  flex: 1;
  text-align: center;
}

.titlebar-right {
  display: flex;
  align-items: center;
  gap: 8px;
  -webkit-app-region: no-drag;
}

.workspace-name {
  font-size: 12px;
  font-weight: 500;
}

/* Windows 窗口控制按钮 */
.window-controls {
  display: flex;
  margin-left: 8px;
  border-left: 1px solid #3e3e42;
  padding-left: 8px;
}

.control-btn {
  width: 36px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background-color 0.2s;
}

.control-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.control-btn.close:hover {
  background-color: #e81123;
}

.control-btn.close:hover .n-icon {
  color: white;
}
</style>
