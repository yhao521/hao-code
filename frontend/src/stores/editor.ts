import {
  WriteFile,
  AddRecentFile,
  AddRecentFolder,
} from "@wails/go/backend/App";
import { defineStore } from "pinia";
import { ref, computed } from "vue";

export interface Workspace {
  path: string;
  name: string;
}

export interface Tab {
  id: string;
  path: string;
  name: string;
  content?: string;
  dirty: boolean;
  language?: string;
}

export type SidebarView = "explorer" | "search" | "git" | "extensions";

export const useEditorStore = defineStore("editor", () => {
  // State
  const activeEditor = ref<string | null>(null);
  const tabs = ref<Tab[]>([]);
  const sidebarVisible = ref(true);
  const sidebarView = ref<SidebarView>("explorer");
  const workspace = ref<Workspace | null>(null);
  const autoSave = ref(true); // 自动保存状态
  const autoSaveDelay = ref(1000); // 自动保存延迟（毫秒）
  const recentFiles = ref<string[]>([]); // 最近打开的文件

  // 自动保存定时器
  let autoSaveTimer: NodeJS.Timeout | null = null;

  // Getters
  const activeTab = computed(() =>
    tabs.value.find((t) => t.id === activeEditor.value),
  );

  const dirtyTabs = computed(() => tabs.value.filter((t) => t.dirty));

  // Actions
  function setWorkspace(path: string) {
    const name = path.split("/").filter(Boolean).pop() || path;
    workspace.value = { path, name };
    // 清空 tabs
    tabs.value = [];
    activeEditor.value = null;

    // 调用后端 API 记录最近文件夹
    AddRecentFolder(path).catch((error) => {
      console.error("Failed to add recent folder:", error);
    });
  }

  function clearWorkspace() {
    workspace.value = null;
    tabs.value = [];
    activeEditor.value = null;
  }

  function openFile(path: string, content: string = "") {
    const existingTab = tabs.value.find((t) => t.path === path);

    if (existingTab) {
      activeEditor.value = existingTab.id;
      return;
    }

    const tab: Tab = {
      id: Date.now().toString(),
      path,
      name: path.split("/").pop() || path,
      content,
      dirty: false,
      language: getLanguageFromPath(path),
    };

    tabs.value.push(tab);
    activeEditor.value = tab.id;

    // 添加到最近文件（后端持久化）
    addToRecentFiles(path);

    // 调用后端 API 记录最近文件
    AddRecentFile(path).catch((error) => {
      console.error("Failed to add recent file:", error);
    });
  }

  function closeTab(id: string) {
    const index = tabs.value.findIndex((t) => t.id === id);
    if (index === -1) return;

    const tab = tabs.value[index];
    if (tab.dirty) {
      // TODO: 提示保存
      console.warn("File has unsaved changes:", tab.path);
      return;
    }

    tabs.value.splice(index, 1);

    if (activeEditor.value === id) {
      activeEditor.value = tabs.value[Math.max(0, index - 1)]?.id || null;
    }
  }

  function closeAllTabs() {
    tabs.value = [];
    activeEditor.value = null;
  }

  function updateContent(id: string, content: string) {
    const tab = tabs.value.find((t) => t.id === id);
    if (tab) {
      tab.content = content;
      tab.dirty = true;
    }
  }

  async function saveFile(id: string) {
    const tab = tabs.value.find((t) => t.id === id);
    if (tab && tab.content !== undefined) {
      try {
        await WriteFile(tab.path, tab.content);
        tab.dirty = false;
        return true;
      } catch (error) {
        console.error("Failed to save file:", error);
        return false;
      }
    }
    return false;
  }

  async function saveAllFiles() {
    const dirtyTabsList = dirtyTabs.value;
    let successCount = 0;

    for (const tab of dirtyTabsList) {
      if (await saveFile(tab.id)) {
        successCount++;
      }
    }

    return successCount;
  }

  function toggleAutoSave() {
    autoSave.value = !autoSave.value;
    if (autoSave.value) {
      startAutoSave();
    } else {
      stopAutoSave();
    }
  }

  function startAutoSave() {
    // 清除现有定时器
    if (autoSaveTimer) {
      clearInterval(autoSaveTimer);
    }

    // 启动新的定时器
    autoSaveTimer = setInterval(async () => {
      const dirtyTabsList = dirtyTabs.value;
      if (dirtyTabsList.length > 0) {
        console.log(`Auto-saving ${dirtyTabsList.length} file(s)...`);
        for (const tab of dirtyTabsList) {
          await saveFile(tab.id);
        }
      }
    }, autoSaveDelay.value);
  }

  function stopAutoSave() {
    if (autoSaveTimer) {
      clearInterval(autoSaveTimer);
      autoSaveTimer = null;
    }
  }

  // 初始化时启动自动保存
  if (autoSave.value) {
    startAutoSave();
  }

  function addToRecentFiles(path: string) {
    // 移除已存在的路径
    recentFiles.value = recentFiles.value.filter((p) => p !== path);
    // 添加到开头
    recentFiles.value.unshift(path);
    // 最多保留 10 个最近文件
    if (recentFiles.value.length > 10) {
      recentFiles.value = recentFiles.value.slice(0, 10);
    }
  }

  function toggleSidebar() {
    sidebarVisible.value = !sidebarVisible.value;
  }

  function setSidebarView(view: SidebarView) {
    sidebarView.value = view;
    sidebarVisible.value = true;
  }

  function getLanguageFromPath(path: string): string {
    const ext = path.split(".").pop()?.toLowerCase();
    const langMap: Record<string, string> = {
      ts: "typescript",
      js: "javascript",
      tsx: "typescriptreact",
      jsx: "javascriptreact",
      py: "python",
      go: "go",
      java: "java",
      cpp: "cpp",
      c: "c",
      rs: "rust",
      rb: "ruby",
      php: "php",
      html: "html",
      css: "css",
      json: "json",
      md: "plaintext",
      yaml: "yaml",
      xml: "xml",
    };
    return langMap[ext || ""] || "plaintext";
  }

  return {
    // State
    activeEditor,
    tabs,
    sidebarVisible,
    sidebarView,
    workspace,
    autoSave,
    autoSaveDelay,
    recentFiles,
    // Getters
    activeTab,
    dirtyTabs,
    // Actions
    openFile,
    closeTab,
    closeAllTabs,
    updateContent,
    saveFile,
    saveAllFiles,
    toggleAutoSave,
    addToRecentFiles,
    toggleSidebar,
    setSidebarView,
    setWorkspace,
    clearWorkspace,
  };
});
