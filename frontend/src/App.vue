<script lang="ts" setup>
import { ref, onMounted, computed } from "vue";
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

// 检测是否为 macOS
const isMacOS = computed(() => {
  const platform = navigator.platform.toLowerCase();
  return platform.includes("mac");
});

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

  // ==================== 帮助菜单事件监听 ====================

  // 欢迎
  EventsOn("menu:welcome", () => {
    message.info("欢迎使用 Hao-Code Editor！🎉");
  });

  // 显示所有命令
  EventsOn("menu:show-all-commands", () => {
    message.info("命令面板功能开发中... (⇧⌘P)");
  });

  // 文档
  EventsOn("menu:documentation", () => {
    window.open("https://github.com/your-repo/hao-code/wiki", "_blank");
    message.info("正在打开文档...");
  });

  // 视频教程
  EventsOn("menu:video-tutorials", () => {
    window.open("https://www.youtube.com/", "_blank");
    message.info("正在打开视频教程...");
  });

  // 键盘快捷方式参考
  EventsOn("menu:keyboard-shortcuts", () => {
    message.info("键盘快捷方式：\n\n保存: ⌘S\n打开: ⌘O\n打开文件夹: ⇧⌘O\n新建: ⌘N\n关闭: ⌘W");
  });

  // 搜索功能请求
  EventsOn("menu:search-feature-requests", () => {
    window.open("https://github.com/your-repo/hao-code/issues", "_blank");
    message.info("正在打开功能请求页面...");
  });

  // 使用英文报告问题
  EventsOn("menu:report-issues", () => {
    window.open("https://github.com/your-repo/hao-code/issues/new", "_blank");
    message.info("正在打开问题报告页面...");
  });

  // 查看许可证
  EventsOn("menu:view-license", () => {
    dialog.info({
      title: "开源许可证",
      content: "本项目使用 GNU Affero General Public License v3.0 (AGPL-3.0) 开源许可证。这是一个强 copyleft 许可证，要求网络使用也必须公开源代码。",
      positiveText: "确定",
    });
  });

  // 隐私声明
  EventsOn("menu:privacy-statement", () => {
    dialog.info({
      title: "隐私声明",
      content: "Hao-Code Editor 尊重您的隐私。所有数据都存储在本地，不会收集或传输任何个人信息。",
      positiveText: "确定",
    });
  });

  // 切换开发人员工具
  EventsOn("menu:toggle-devtools", () => {
    // 使用浏览器原生的开发者工具
    // @ts-ignore
    if (window.devToolsExtension) {
      // @ts-ignore
      window.devToolsExtension.open();
    } else {
      message.info("请使用浏览器快捷键: ⌘I (macOS) 或 Ctrl+Shift+I (Windows/Linux)");
    }
  });

  // 打开进程资源管理器
  EventsOn("menu:open-process-explorer", () => {
    message.info("进程资源管理器功能开发中...");
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
            <!-- 标题栏（仅 Windows/Linux 显示，macOS 使用系统标题栏） -->
            <TitleBar v-if="!isMacOS" />

            <!-- 主内容区 - 可拖拽分割 -->
            <div class="main-content" :class="{ 'macos-content': isMacOS }">
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

/* macOS: 没有自定义标题栏，内容区域需要额外的顶部内边距以避免被交通灯按钮遮挡 */
.macos-content {
  padding-top: 20px;
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
