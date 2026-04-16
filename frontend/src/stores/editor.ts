import {
  WriteFile,
  AddRecentFile,
  AddRecentFolder,
} from "@wails/backend/appservice";
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

export interface DiffInfo {
  path: string;
  oldContent: string;
  newContent: string;
}

export interface Breakpoint {
  id: number;
  path: string;
  line: number;
  verified: boolean;
}

export type SidebarView = "explorer" | "search" | "git" | "extensions";

export interface EditorGroup {
  id: string;
  tabs: Tab[];
  activeTabId: string | null;
}

export const useEditorStore = defineStore("editor", () => {
  // State
  const editorGroups = ref<EditorGroup[]>([
    { id: "1", tabs: [], activeTabId: null },
  ]);
  const activeGroupId = ref<string>("1");
  const sidebarVisible = ref(true);
  const sidebarView = ref<SidebarView>("explorer");
  const workspace = ref<Workspace | null>(null);
  const autoSave = ref(true); // 自动保存状态
  const autoSaveDelay = ref(1000); // 自动保存延迟（毫秒）
  const recentFiles = ref<string[]>([]); // 最近打开的文件
  const isDiffMode = ref(false);
  const diffInfo = ref<DiffInfo | null>(null);
  const breakpoints = ref<Map<string, Set<number>>>(new Map());
  const currentDebugLine = ref<{ path: string; line: number } | null>(null);
  const activeMonacoModel = ref<any>(null);
  const activeCursor = ref<{ lineNumber: number; column: number } | null>(null);

  // 自动保存定时器
  let autoSaveTimer: NodeJS.Timeout | null = null;

  // Getters
  const activeGroup = computed(
    () =>
      editorGroups.value.find((g) => g.id === activeGroupId.value) ||
      editorGroups.value[0],
  );
  const activeTab = computed(() =>
    activeGroup.value.tabs.find((t) => t.id === activeGroup.value.activeTabId),
  );
  const tabs = computed(() => activeGroup.value.tabs);
  const dirtyTabs = computed(() =>
    activeGroup.value.tabs.filter((t) => t.dirty),
  );

  // Actions
  function setWorkspace(path: string) {
    const name = path.split("/").filter(Boolean).pop() || path;
    workspace.value = { path, name };
    // Reset groups to a single empty group
    editorGroups.value = [
      { id: Date.now().toString(), tabs: [], activeTabId: null },
    ];
    activeGroupId.value = editorGroups.value[0].id;

    // 调用后端 API 记录最近文件夹
    AddRecentFolder(path).catch((error) => {
      console.error("Failed to add recent folder:", error);
    });

    // 触发工作区变更事件，通知 FileExplorer 重新加载
    window.dispatchEvent(
      new CustomEvent("workspace-changed", { detail: path }),
    );
  }

  function clearWorkspace() {
    workspace.value = null;
    editorGroups.value = [
      { id: Date.now().toString(), tabs: [], activeTabId: null },
    ];
    activeGroupId.value = editorGroups.value[0].id;
  }

  function setActiveGroup(id: string) {
    activeGroupId.value = id;
  }

  function splitEditor(direction: "right" | "down") {
    const currentGroup = activeGroup.value;
    if (!currentGroup.activeTabId) return;

    const newGroupId = Date.now().toString();
    const newGroup: EditorGroup = {
      id: newGroupId,
      tabs: [...currentGroup.tabs.map((t) => ({ ...t }))], // Clone tabs for simplicity or share them
      activeTabId: currentGroup.activeTabId,
    };

    // For now, let's just add a new empty group and focus it
    // A real implementation would handle layout tree (split-view)
    // Here we simulate by adding a group. In UI we need a SplitView component.
    editorGroups.value.push(newGroup);
    activeGroupId.value = newGroupId;
  }

  function openFile(path: string, content: string = "") {
    const group = activeGroup.value;
    const existingTab = group.tabs.find((t) => t.path === path);

    if (existingTab) {
      group.activeTabId = existingTab.id;
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

    group.tabs.push(tab);
    group.activeTabId = tab.id;

    // 添加到最近文件（后端持久化）
    addToRecentFiles(path);

    // 调用后端 API 记录最近文件
    AddRecentFile(path).catch((error) => {
      console.error("Failed to add recent file:", error);
    });
  }

  function closeTab(id: string) {
    const group = activeGroup.value;
    const index = group.tabs.findIndex((t) => t.id === id);
    if (index === -1) return;

    const tab = group.tabs[index];
    if (tab.dirty) {
      // TODO: 提示保存
      console.warn("File has unsaved changes:", tab.path);
      return;
    }

    group.tabs.splice(index, 1);

    if (group.activeTabId === id) {
      group.activeTabId = group.tabs[Math.max(0, index - 1)]?.id || null;
    }
  }

  function closeAllTabs() {
    activeGroup.value.tabs = [];
    activeGroup.value.activeTabId = null;
  }

  function updateContent(id: string, content: string) {
    const tab = activeGroup.value.tabs.find((t) => t.id === id);
    if (tab) {
      tab.content = content;
      tab.dirty = true;
    }
  }

  async function saveFile(id: string) {
    const tab = activeGroup.value.tabs.find((t) => t.id === id);
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

  function setDiffMode(enabled: boolean, info?: DiffInfo) {
    isDiffMode.value = enabled;
    if (info) {
      diffInfo.value = info;
    }
  }

  function toggleBreakpoint(path: string, line: number) {
    if (!breakpoints.value.has(path)) {
      breakpoints.value.set(path, new Set());
    }
    const lines = breakpoints.value.get(path)!;
    if (lines.has(line)) {
      lines.delete(line);
    } else {
      lines.add(line);
    }
  }

  function hasBreakpoint(path: string, line: number): boolean {
    return breakpoints.value.get(path)?.has(line) || false;
  }

  function setCurrentDebugLine(path: string, line: number) {
    currentDebugLine.value = { path, line };
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
    editorGroups,
    activeGroupId,
    sidebarVisible,
    sidebarView,
    workspace,
    autoSave,
    autoSaveDelay,
    recentFiles,
    isDiffMode,
    diffInfo,
    // Getters
    activeGroup,
    activeTab,
    tabs,
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
    setDiffMode,
    toggleBreakpoint,
    hasBreakpoint,
    setCurrentDebugLine,
    breakpoints,
    currentDebugLine,
    activeMonacoModel,
    activeCursor,
    setActiveGroup,
    splitEditor,
  };
});
