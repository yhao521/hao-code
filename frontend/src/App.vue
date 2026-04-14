<script lang="ts" setup>
import { ref, onMounted } from "vue";
import {
  NConfigProvider,
  NMessageProvider,
  NNotificationProvider,
  NDialogProvider,
  NModal,
  NForm,
  NFormItem,
  NInput,
  NSpace,
  NButton,
  darkTheme,
  zhCN,
  dateZhCN,
  useMessage,
  useDialog,
} from "naive-ui";
import { EventsOn } from "@wails/runtime/runtime";

// 导入组件
import TitleBar from "./components/layout/TitleBar.vue";
import SideBar from "./components/layout/SideBar.vue";
import EditorArea from "./components/editor/EditorArea.vue";
import StatusBar from "./components/layout/StatusBar.vue";
import ResizableSplit from "./components/layout/ResizableSplit.vue";
import { useEditorStore } from "./stores/editor";
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
} from "@wails/go/backend/App";

const editorStore = useEditorStore();
const message = useMessage();
const dialog = useDialog();

// 新建文件对话框
const showNewFileModal = ref(false);
const newFileName = ref("");

// 主题配置 - VSCode 风格深色主题
const theme = ref(darkTheme);
const themeOverrides = {
  common: {
    primaryColor: "#0E639C",
    primaryColorHover: "#1177BB",
    bodyColor: "#1E1E1E",
    cardColor: "#252526",
    borderColor: "#3E3E42",
    textColor: "#CCCCCC",
  },
  Menu: {
    itemColor: "#CCCCCC",
    itemColorActive: "#FFFFFF",
    itemColorHover: "#2A2D2E",
  },
  Tree: {
    nodeColor: "#CCCCCC",
    nodeColorHover: "#2A2D2E",
  },
};

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

  // 打开文件夹... - 触发标题栏的打开文件夹功能
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
      await ListDir(selectedPath);
      message.destroyAll();

      editorStore.setWorkspace(selectedPath);
      message.success(`已打开: ${selectedPath.split("/").pop()}`);
    } catch (error) {
      message.destroyAll();
      const errorMsg = error instanceof Error ? error.message : String(error);
      if (!errorMsg.includes("cancelled")) {
        message.error(`打开文件夹失败: ${errorMsg}`);
      } else {
        message.info("已取消选择");
      }
    }
  });

  // 保存当前文件
  EventsOn("menu:save", async () => {
    if (!editorStore.activeTab) {
      message.warning("没有打开的文件");
      return;
    }

    const tab = editorStore.activeTab;
    if (!tab.dirty) {
      message.info("文件已保存");
      return;
    }

    const success = await editorStore.saveFile(tab.id);
    if (success) {
      message.success(`已保存: ${tab.name}`);
    } else {
      message.error("保存失败");
    }
  });

  // 另存为...
  EventsOn("menu:save-as", async () => {
    if (!editorStore.activeTab) {
      message.warning("没有打开的文件");
      return;
    }

    try {
      message.loading("正在打开保存对话框...", { duration: 0 });
      const selectedPath = await SaveFileDialog();
      message.destroyAll();

      if (!selectedPath) {
        message.info("已取消保存");
        return;
      }

      // 更新当前标签页的路径和内容
      const tab = editorStore.activeTab;
      if (tab && tab.content !== undefined) {
        await WriteFile(selectedPath, tab.content);
        tab.path = selectedPath;
        tab.name = selectedPath.split("/").pop() || selectedPath;
        tab.dirty = false;
        message.success(`已另存为: ${tab.name}`);
      }
    } catch (error) {
      message.destroyAll();
      const errorMsg = error instanceof Error ? error.message : String(error);
      if (!errorMsg.includes("cancelled")) {
        message.error(`另存为失败: ${errorMsg}`);
      } else {
        message.info("已取消保存");
      }
    }
  });

  // 全部保存
  EventsOn("menu:save-all", async () => {
    if (editorStore.dirtyTabs.length === 0) {
      message.info("所有文件已保存");
      return;
    }

    const count = await editorStore.saveAllFiles();
    message.success(`已保存 ${count} 个文件`);
  });

  // 切换自动保存
  EventsOn("menu:toggle-auto-save", () => {
    editorStore.toggleAutoSave();
    const status = editorStore.autoSave ? "已启用" : "已禁用";
    message.info(`自动保存${status}`);
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
        content: `"${tab.name}" 有未保存的更改，是否保存？`,
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
        content: `有 ${editorStore.dirtyTabs.length} 个文件未保存，是否保存后关闭？`,
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
});
</script>

<template>
  <NConfigProvider
    :theme="theme"
    :theme-overrides="themeOverrides"
    :locale="zhCN"
    :date-locale="dateZhCN"
  >
    <NMessageProvider>
      <NNotificationProvider>
        <NDialogProvider>
          <div class="app-container">
            <!-- 标题栏（所有平台显示） -->
            <TitleBar />

            <!-- 主内容区 - 可拖拽分割 -->
            <div class="main-content">
              <ResizableSplit
                :min="180"
                :max="500"
                :horizontal="true"
                class="main-split"
              >
                <template #1>
                  <div class="sidebar-container">
                    <SideBar />
                  </div>
                </template>
                <template #2>
                  <div class="editor-container">
                    <EditorArea />
                  </div>
                </template>
              </ResizableSplit>
            </div>

            <!-- 状态栏 -->
            <StatusBar />
          </div>
        </NDialogProvider>
      </NNotificationProvider>
    </NMessageProvider>
  </NConfigProvider>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family:
    -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue",
    Arial, sans-serif;
  overflow: hidden;
}

#app {
  width: 100vw;
  height: 100vh;
}

.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #1e1e1e;
}

.main-content {
  flex: 1;
  overflow: hidden;
  background-color: #1e1e1e;
}

.main-split {
  width: 100%;
  height: 100%;
}

.sidebar-container {
  height: 100%;
  background-color: #252526;
  border-right: 1px solid #3e3e42;
  overflow: hidden;
}

.editor-container {
  height: 100%;
  background-color: #1e1e1e;
  overflow: hidden;
}
</style>
