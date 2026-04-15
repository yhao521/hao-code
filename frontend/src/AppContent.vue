<script lang="ts" setup>
import { ref, onMounted, computed } from "vue";
import {
  NModal,
  NForm,
  NFormItem,
  NInput,
  NSpace,
  NButton,
  useMessage,
  useDialog,
} from "naive-ui";
import { Events } from "@wailsio/runtime";
const EventsOn = Events.On;

// 导入组件
import TitleBar from "./components/layout/TitleBar.vue";
import SideBar from "./components/layout/SideBar.vue";
import EditorArea from "./components/editor/EditorArea.vue";
import StatusBar from "./components/layout/StatusBar.vue";
import LayoutToolbar from "./components/layout/LayoutToolbar.vue";
import RightPanel from "./components/layout/RightPanel.vue";
import BottomPanel from "./components/layout/BottomPanel.vue";
import ResizableSplit from "./components/layout/ResizableSplit.vue";
import RecentFilesModal from "./components/RecentFilesModal.vue";
import CommandPalette from "./components/CommandPalette.vue";
import { useEditorStore } from "./stores/editor";
import { useLayoutStore } from "./stores/layout";
import {
  GetProjectRoot,
  ListDir,
  ReadFile,
  CreateFile,
  CreateDirectory,
  DeleteFileOrDirectory,
  RenameFileOrDirectory,
  CopyFileOrDirectory,
  OpenFileDialog,
  OpenFolderDialog,
  SaveFileDialog,
  WriteFile,
  AddRecentFile,
  AddRecentFolder,
} from "@wails/backend/appservice";

const editorStore = useEditorStore();
const layoutStore = useLayoutStore();
const message = useMessage();
const dialog = useDialog();

// 检测是否为 macOS
const isMacOS = computed(() => {
  const platform = navigator.platform.toLowerCase();
  return platform.includes("mac");
});

// 新建文件对话框
const showNewFileModal = ref(false);
const newFileName = ref("");

// 最近文件模态框引用
const recentFilesModalRef = ref<any>(null);

// 命令面板引用
const commandPaletteRef = ref<any>(null);

onMounted(() => {
  // ==================== 文件菜单事件监听 ====================

  // 新建文本文件 - 直接创建 untitled 文件
  EventsOn("menu:new-text-file", () => {
    const untitledPath = `untitled-${Date.now()}.txt`;
    editorStore.openFile(untitledPath, "");
    message.success("已创建新文本文件");
  });

  // 新建文件... - 显示对话框
  EventsOn("menu:new-file", () => {
    showNewFileModal.value = true;
    newFileName.value = "";
  });

  // 打开文件...
  EventsOn("menu:open-file", async () => {
    try {
      message.loading("正在打开文件选择对话框...", { duration: 0 });
      const selectedPath = await OpenFileDialog();
      message.destroyAll();

      if (!selectedPath) {
        message.info("已取消选择");
        return;
      }

      message.loading("正在读取文件...", { duration: 0 });
      const content = await ReadFile(selectedPath);
      message.destroyAll();

      editorStore.openFile(selectedPath, content);
      message.success(`已打开: ${selectedPath.split("/").pop()}`);

      // 记录到最近文件
      try {
        await AddRecentFile(selectedPath);
      } catch (error) {
        console.error("Failed to add recent file:", error);
      }
    } catch (error) {
      message.destroyAll();
      const errorMsg = error instanceof Error ? error.message : String(error);
      if (!errorMsg.includes("cancelled")) {
        message.error(`打开文件失败: ${errorMsg}`);
      } else {
        message.info("已取消选择");
      }
    }
  });

  // 打开文件夹
  EventsOn("menu:open-folder", async () => {
    try {
      message.loading("正在打开文件夹选择对话框...", { duration: 0 });
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

      editorStore.setWorkspace(selectedPath);
      message.destroyAll();
      message.success(`已打开: ${selectedPath.split("/").pop()}`);

      // 记录到最近文件夹
      try {
        await AddRecentFolder(selectedPath);
      } catch (error) {
        console.error("Failed to add recent folder:", error);
      }
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
  });

  // 打开最近的文件
  EventsOn("menu:open-recent", () => {
    recentFilesModalRef.value?.show();
  });

  // 保存
  EventsOn("menu:save", async () => {
    if (!editorStore.activeTab) {
      message.warning("没有打开的文件");
      return;
    }

    await editorStore.saveFile(editorStore.activeTab.id);
    message.success("文件已保存");
  });

  // 另存为...
  EventsOn("menu:save-as", async () => {
    if (!editorStore.activeTab) {
      message.warning("没有打开的文件");
      return;
    }

    // 确保 content 不为 undefined
    const content = editorStore.activeTab.content ?? "";

    try {
      const savePath = await SaveFileDialog();
      if (!savePath) {
        message.info("已取消保存");
        return;
      }

      await WriteFile(savePath, content);
      message.success(`文件已保存到: ${savePath}`);
    } catch (error) {
      const errorMsg = error instanceof Error ? error.message : String(error);
      message.error(`保存文件失败: ${errorMsg}`);
    }
  });

  // 切换自动保存
  EventsOn("menu:toggle-auto-save", () => {
    editorStore.toggleAutoSave();
    message.success(editorStore.autoSave ? "自动保存已启用" : "自动保存已禁用");
  });

  // 关闭当前编辑器标签
  EventsOn("menu:close-editor", () => {
    if (!editorStore.activeTab) {
      message.warning("没有打开的文件");
      return;
    }

    const tab = editorStore.activeTab;
    if (tab.dirty) {
      dialog.warning({
        title: "确认关闭",
        content: `"${tab.name}" 有未保存的更改,是否保存?`,
        positiveText: "保存",
        negativeText: "不保存",
        onPositiveClick: async () => {
          await editorStore.saveFile(tab.id);
          editorStore.closeTab(tab.id);
        },
        onNegativeClick: () => {
          editorStore.closeTab(tab.id);
        },
      });
    } else {
      editorStore.closeTab(tab.id);
    }
  });

  // 关闭文件夹
  EventsOn("menu:close-folder", () => {
    if (!editorStore.workspace) {
      message.warning("没有打开的文件夹");
      return;
    }

    if (editorStore.dirtyTabs.length > 0) {
      dialog.warning({
        title: "确认关闭",
        content: `有 ${editorStore.dirtyTabs.length} 个文件未保存,是否保存后关闭?`,
        positiveText: "保存并关闭",
        negativeText: "不保存直接关闭",
        onPositiveClick: async () => {
          await editorStore.saveAllFiles();
          editorStore.clearWorkspace();
          message.success("已关闭文件夹");
        },
        onNegativeClick: () => {
          editorStore.clearWorkspace();
          message.success("已关闭文件夹");
        },
      });
    } else {
      editorStore.clearWorkspace();
      message.success("已关闭文件夹");
    }
  });

  // 新建窗口 - 目前暂不支持
  EventsOn("menu:new-window", () => {
    message.info("多窗口功能开发中...");
  });

  // 欢迎
  EventsOn("menu:welcome", () => {
    message.info("欢迎使用 Hao-Code Editor!");
  });

  // 显示所有命令
  EventsOn("menu:show-all-commands", () => {
    commandPaletteRef.value?.show();
  });

  // 文档
  EventsOn("menu:documentation", () => {
    message.info("文档功能开发中...");
  });

  // 视频教程
  EventsOn("menu:video-tutorials", () => {
    message.info("视频教程功能开发中...");
  });

  // 键盘快捷方式参考
  EventsOn("menu:keyboard-shortcuts", () => {
    message.info("快捷键参考功能开发中...");
  });

  // 搜索功能请求
  EventsOn("menu:search-feature-requests", () => {
    message.info("功能请求搜索功能开发中...");
  });

  // 使用英文报告问题
  EventsOn("menu:report-issues", () => {
    message.info("问题报告功能开发中...");
  });

  // 查看许可证
  EventsOn("menu:view-license", () => {
    message.info("许可证查看功能开发中...");
  });

  // 隐私声明
  EventsOn("menu:privacy-statement", () => {
    message.info("隐私声明功能开发中...");
  });

  // 切换开发人员工具
  EventsOn("menu:toggle-devtools", () => {
    message.info("开发人员工具切换功能开发中...");
  });

  // 打开进程资源管理器
  EventsOn("menu:open-process-explorer", () => {
    message.info("进程资源管理器功能开发中...");
  });

  // ==================== 快捷键绑定 ====================
  window.addEventListener("keydown", (e) => {
    // Ctrl+Shift+F: 全局搜索
    if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key === "F") {
      e.preventDefault();
      editorStore.setSidebarView("search");
    }
  });
});

// 处理新建文件
async function handleCreateFile() {
  if (!newFileName.value.trim()) {
    message.warning("请输入文件名");
    return;
  }

  try {
    const workspacePath = editorStore.workspace?.path || "";
    const fullPath = workspacePath
      ? `${workspacePath}/${newFileName.value}`
      : newFileName.value;

    await CreateFile(fullPath);
    message.success(`文件已创建: ${newFileName.value}`);
    showNewFileModal.value = false;

    // 读取并打开新文件
    const content = await ReadFile(fullPath);
    editorStore.openFile(fullPath, content);
  } catch (error) {
    const errorMsg = error instanceof Error ? error.message : String(error);
    message.error(`创建文件失败: ${errorMsg}`);
  }
}
</script>

<template>
  <div class="app-container" :class="{ 'macos-content': isMacOS }">
    <!-- 标题栏（仅 Windows/Linux 显示，macOS 使用系统标题栏） -->
    <TitleBar v-if="!isMacOS" />

    <!-- 布局切换工具栏 -->
    <LayoutToolbar />

    <!-- 主内容区 - 可拖拽分割 -->
    <div class="main-content">
      <!-- 主分割：左侧边栏 + 中间区域 -->
      <ResizableSplit
        v-if="layoutStore.sidebarVisible"
        :min="180"
        :max="500"
        :horizontal="true"
        :default-size="layoutStore.sidebarWidth"
        @update:size="layoutStore.sidebarWidth = $event"
        class="main-split"
      >
        <template #1>
          <div class="sidebar-container">
            <SideBar />
          </div>
        </template>
        <template #2>
          <!-- 中间区域：编辑器 + 右侧面板 -->
          <div class="center-area">
            <ResizableSplit
              v-if="layoutStore.rightPanelVisible"
              :min="200"
              :max="500"
              :horizontal="true"
              :default-size="layoutStore.rightPanelWidth"
              @update:size="layoutStore.rightPanelWidth = $event"
              class="editor-right-split"
              direction="rtl"
            >
              <template #1>
                <div class="right-panel-container">
                  <RightPanel />
                </div>
              </template>
              <template #2>
                <div class="editor-container">
                  <EditorArea />
                </div>
              </template>
            </ResizableSplit>

            <!-- 无右侧面板时的编辑器容器 -->
            <div v-else class="editor-container-full">
              <EditorArea />
            </div>
          </div>
        </template>
      </ResizableSplit>

      <!-- 隐藏侧边栏时的布局 -->
      <div v-else class="main-area-without-sidebar">
        <div class="center-area">
          <ResizableSplit
            v-if="layoutStore.rightPanelVisible"
            :min="200"
            :max="500"
            :horizontal="true"
            :default-size="layoutStore.rightPanelWidth"
            @update:size="layoutStore.rightPanelWidth = $event"
            class="editor-right-split"
            direction="rtl"
          >
            <template #1>
              <div class="right-panel-container">
                <RightPanel />
              </div>
            </template>
            <template #2>
              <div class="editor-container">
                <EditorArea />
              </div>
            </template>
          </ResizableSplit>

          <div v-else class="editor-container-full">
            <EditorArea />
          </div>
        </div>
      </div>
    </div>

    <!-- 下面板 -->
    <ResizableSplit
      v-if="layoutStore.bottomPanelVisible"
      :min="100"
      :max="600"
      :horizontal="false"
      :default-size="layoutStore.bottomPanelHeight"
      @update:size="layoutStore.bottomPanelHeight = $event"
      class="bottom-split"
    >
      <template #1>
        <div class="bottom-panel-container">
          <BottomPanel />
        </div>
      </template>
      <template #2>
        <div class="bottom-placeholder"></div>
      </template>
    </ResizableSplit>

    <!-- 状态栏 -->
    <StatusBar />

    <!-- 最近文件模态框 -->
    <RecentFilesModal ref="recentFilesModalRef" />

    <!-- 命令面板 -->
    <CommandPalette ref="commandPaletteRef" />

    <!-- 新建文件对话框 -->
    <NModal v-model:show="showNewFileModal" preset="dialog" title="新建文件">
      <NForm
        :model="{ name: newFileName }"
        label-placement="left"
        label-width="80"
      >
        <NFormItem label="文件名" path="name">
          <NInput
            v-model:value="newFileName"
            placeholder="例如: index.ts"
            @keyup.enter="handleCreateFile"
          />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace>
          <NButton @click="showNewFileModal = false">取消</NButton>
          <NButton type="primary" @click="handleCreateFile">创建</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>

<style>
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100vw;
  background-color: #1e1e1e;
  overflow: hidden;
}

.main-content {
  flex: 1;
  overflow: hidden;
  background-color: #1e1e1e;
  display: flex;
  flex-direction: column;
}

/* macOS: 没有自定义标题栏，内容区域需要额外的顶部内边距以避免被交通灯按钮遮挡 */
.macos-content {
  padding-top: 20px;
  box-sizing: border-box;
  height: calc(100vh - 20px);
}

.main-split {
  flex: 1;
  width: 100%;
}

.bottom-split {
  height: auto;
  width: 100%;
}

.sidebar-container {
  height: 100%;
  background-color: #252526;
  border-right: 1px solid #3e3e42;
  overflow: hidden;
}

.center-area {
  display: flex;
  flex: 1;
  height: 100%;
  width: 100%;
  overflow: hidden;
}

.editor-container {
  flex: 1;
  height: 100%;
  background-color: #1e1e1e;
  overflow: hidden;
}

.editor-container-full {
  flex: 1;
  height: 100%;
  background-color: #1e1e1e;
  overflow: hidden;
}

.right-panel-container {
  height: 100%;
  background-color: #252526;
  overflow: hidden;
}

.bottom-panel-container {
  height: 100%;
  background-color: #1e1e1e;
  overflow: hidden;
}

.bottom-placeholder {
  flex: 1;
  background-color: #1e1e1e;
}

.main-area-without-sidebar {
  flex: 1;
  display: flex;
  width: 100%;
  overflow: hidden;
}

.editor-right-split {
  height: 100%;
  width: 100%;
}
</style>
